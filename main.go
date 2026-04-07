package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// =============================================================================
// EMBEDDED COMMAND CATALOGS
// =============================================================================

//go:embed vppctl-cmds/vppctl-commands-summary.md
var vppctlCommandCatalog string

// =============================================================================
// CONFIGURATION & CONSTANTS
// =============================================================================

// VPPExecutionMode defines how VPP is running
type VPPExecutionMode string

const (
	// ModeKubernetes - VPP runs as a container in a Kubernetes pod (CalicoVPP)
	ModeKubernetes VPPExecutionMode = "kubernetes"
	// ModeStandalone - VPP runs as a daemon/service on baremetal/VM
	ModeStandalone VPPExecutionMode = "standalone"
)

// VPP socket and lock paths
const (
	defaultVppSockPath = "/var/run/vpp/cli.sock"
	vppctlPath         = "/usr/bin/vppctl"
	// Lock file path - on host for standalone, in container for kubernetes
	captureLockFile = "/tmp/vpp-mcp-capture.lock"
)

const kubeClientTimeout = 30 * time.Second

// CaptureLockInfo represents information stored in the lock file
type CaptureLockInfo struct {
	Operation   string    `json:"operation"`
	StartedAt   time.Time `json:"started_at"`
	Hostname    string    `json:"hostname"`
	Target      string    `json:"target"` // pod name or "localhost" for standalone
	CaptureType string    `json:"capture_type"`
}

// Global lock to prevent concurrent capture operations within the same MCP server instance
var captureMutex sync.Mutex

// =============================================================================
// TYPES & DATA STRUCTURES
// =============================================================================

// VPPTarget represents the target where VPP is running
type VPPTarget struct {
	Mode          VPPExecutionMode
	PodName       string // For Kubernetes mode
	Namespace     string // For Kubernetes mode (default: calico-vpp-dataplane)
	ContainerName string // For Kubernetes mode (default: vpp)
	VppSockPath   string // For Standalone mode (default: /var/run/vpp/cli.sock)
}

// NewKubernetesTarget creates a VPP target for Kubernetes/CalicoVPP mode
func NewKubernetesTarget(podName string) *VPPTarget {
	return &VPPTarget{
		Mode:          ModeKubernetes,
		PodName:       podName,
		Namespace:     "calico-vpp-dataplane",
		ContainerName: "vpp",
	}
}

// NewStandaloneTarget creates a VPP target for standalone daemon mode
func NewStandaloneTarget(sockPath string) *VPPTarget {
	if sockPath == "" {
		sockPath = defaultVppSockPath
	}
	return &VPPTarget{
		Mode:        ModeStandalone,
		VppSockPath: sockPath,
	}
}

// VPPMCPServer implements the MCP server for VPP debugging
type VPPMCPServer struct {
	server *mcp.Server
}

// NewVPPMCPServer creates a new VPP MCP server
func NewVPPMCPServer() *VPPMCPServer {
	return &VPPMCPServer{}
}

// KubeClient wraps Kubernetes client for VPP operations
type KubeClient struct {
	clientset *kubernetes.Clientset
	timeout   time.Duration
}

// CoreV1 returns the CoreV1 client
func (k *KubeClient) CoreV1() corev1client.CoreV1Interface {
	return k.clientset.CoreV1()
}

// newKubeClient creates a new Kubernetes client
func newKubeClient() (*KubeClient, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return &KubeClient{clientset: clientset, timeout: kubeClientTimeout}, nil
}

// =============================================================================
// INPUT TYPES - TOOL PARAMETER STRUCTURES
// =============================================================================

// VPPCommandInput represents the generic input for VPP command tools
type VPPCommandInput struct {
	// Mode specifies how VPP is running: "kubernetes" (default) or "standalone"
	Mode string `json:"mode,omitempty"`
	// PodName specifies the name of the Kubernetes pod running VPP (required for kubernetes mode)
	PodName string `json:"pod_name,omitempty"`
	// SockPath specifies the VPP socket path for standalone mode (default: /var/run/vpp/cli.sock)
	SockPath string `json:"sock_path,omitempty"`
}

// GetTarget creates a VPPTarget from the input
func (v *VPPCommandInput) GetTarget() *VPPTarget {
	mode := strings.ToLower(v.Mode)
	if mode == "standalone" || mode == "local" || mode == "daemon" {
		return NewStandaloneTarget(v.SockPath)
	}
	// Default to Kubernetes mode
	return NewKubernetesTarget(v.PodName)
}

// VPPCaptureInput represents the input for VPP packet capture tools (trace, pcap, dispatch)
type VPPCaptureInput struct {
	// Mode specifies how VPP is running: "kubernetes" (default) or "standalone"
	Mode string `json:"mode,omitempty"`
	// PodName specifies the name of the Kubernetes pod running VPP (required for kubernetes mode)
	PodName string `json:"pod_name,omitempty"`
	// SockPath specifies the VPP socket path for standalone mode (default: /var/run/vpp/cli.sock)
	SockPath string `json:"sock_path,omitempty"`
	// Count specifies the number of packets to capture (default: 500)
	Count int `json:"count,omitempty"`
	// Timeout specifies the capture timeout in seconds (default: 30)
	Timeout int `json:"timeout,omitempty"`
	// Interface specifies the interface type or name to capture from
	Interface string `json:"interface,omitempty"`
}

// GetTarget creates a VPPTarget from the input
func (v *VPPCaptureInput) GetTarget() *VPPTarget {
	mode := strings.ToLower(v.Mode)
	if mode == "standalone" || mode == "local" || mode == "daemon" {
		return NewStandaloneTarget(v.SockPath)
	}
	// Default to Kubernetes mode
	return NewKubernetesTarget(v.PodName)
}

// VppctlInput represents input for the vppctl tool
type VppctlInput struct {
	// Command is the vppctl command to run (e.g., "show version", "show int addr", "show ip fib index 0")
	Command string `json:"command"`
	// PodName specifies the Kubernetes pod running VPP (required for kubernetes mode)
	PodName string `json:"pod_name,omitempty"`
	// Mode specifies how VPP is running: "kubernetes" (default) or "standalone"
	Mode string `json:"mode,omitempty"`
	// SockPath specifies the VPP socket path for standalone mode (default: /var/run/vpp/cli.sock)
	SockPath string `json:"sock_path,omitempty"`
	// Count specifies packet count for capture commands (default: 500)
	Count int `json:"count,omitempty"`
	// Timeout specifies capture timeout in seconds (default: 30)
	Timeout int `json:"timeout,omitempty"`
	// Interface specifies interface for capture commands
	Interface string `json:"interface,omitempty"`
}

// GetTarget creates a VPPTarget from the VppctlInput
func (v *VppctlInput) GetTarget() *VPPTarget {
	mode := strings.ToLower(v.Mode)
	if mode == "standalone" || mode == "local" || mode == "daemon" {
		return NewStandaloneTarget(v.SockPath)
	}
	return NewKubernetesTarget(v.PodName)
}

// GobgpInput represents input for the gobgp tool
type GobgpInput struct {
	// Command is the gobgp command to run (e.g., "neighbor", "global rib -a ipv4", "neighbor 10.0.0.1 adj-in")
	Command string `json:"command"`
	// PodName specifies the Kubernetes pod running the agent container with gobgp
	PodName string `json:"pod_name"`
}

// ClusterInput represents input for the cluster tool
type ClusterInput struct {
	// Command is the cluster operation (e.g., "get_pods", "get_nodes", "get_configmap", "logs", "describe_pod")
	Command string `json:"command"`
	// Namespace specifies the Kubernetes namespace (default: calico-vpp-dataplane)
	Namespace string `json:"namespace,omitempty"`
	// PodName specifies a pod name (for logs, describe_pod)
	PodName string `json:"pod_name,omitempty"`
	// Container specifies the container name (for logs; default: agent)
	Container string `json:"container,omitempty"`
	// ResourceName specifies a resource name (for get_configmap, get_daemonset, describe_node)
	ResourceName string `json:"resource_name,omitempty"`
	// TailLines specifies number of log lines (default: 200, max: 5000)
	TailLines int `json:"tail_lines,omitempty"`
}

// EmptyInput represents tools that don't require any input parameters
type EmptyInput struct{}

// =============================================================================
// VPP COMMAND EXECUTION - SUPPORTS BOTH KUBERNETES AND STANDALONE MODES
// =============================================================================

// ExecuteVPPCommand runs a VPP command on the target (either Kubernetes pod or local daemon)
func ExecuteVPPCommand(ctx context.Context, target *VPPTarget, command string) (map[string]interface{}, error) {
	switch target.Mode {
	case ModeKubernetes:
		return executeKubernetesVPPCommand(ctx, target, command)
	case ModeStandalone:
		return executeStandaloneVPPCommand(ctx, target, command)
	default:
		return nil, fmt.Errorf("unknown VPP execution mode: %s", target.Mode)
	}
}

// executeKubernetesVPPCommand runs a VPP command on a Kubernetes pod
func executeKubernetesVPPCommand(ctx context.Context, target *VPPTarget, command string) (map[string]interface{}, error) {
	if target.PodName == "" {
		return nil, fmt.Errorf("pod name is required for Kubernetes mode")
	}

	namespace := target.Namespace
	if namespace == "" {
		namespace = "calico-vpp-dataplane"
	}
	containerName := target.ContainerName
	if containerName == "" {
		containerName = "vpp"
	}

	// Build kubectl exec command
	cmdArgs := []string{
		"exec",
		"-n", namespace,
		target.PodName,
		"-c", containerName,
	}

	// Add the vppctl command
	cmdArgs = append(cmdArgs, "--", "vppctl")

	// Add the specific VPP command arguments
	cmdArgs = append(cmdArgs, strings.Fields(command)...)

	log.Printf("[K8s] Executing: kubectl %s", strings.Join(cmdArgs, " "))

	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "kubectl", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	execErr := cmd.Run()

	output := stdout.Bytes()
	errOutput := stderr.String()

	if errOutput != "" {
		log.Printf("[K8s] stderr: %s", errOutput)
	}

	if execErr != nil {
		errorMsg := ""
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			errorMsg = string(exitErr.Stderr)
		}
		return map[string]interface{}{
			"success":   false,
			"error":     fmt.Sprintf("%v - %s", execErr, errorMsg),
			"mode":      string(ModeKubernetes),
			"pod":       target.PodName,
			"namespace": namespace,
			"command":   command,
		}, execErr
	}

	return map[string]interface{}{
		"success":   true,
		"output":    string(output),
		"command":   command,
		"mode":      string(ModeKubernetes),
		"pod":       target.PodName,
		"namespace": namespace,
		"container": containerName,
	}, nil
}

// executeStandaloneVPPCommand runs a VPP command on a local VPP daemon
func executeStandaloneVPPCommand(ctx context.Context, target *VPPTarget, command string) (map[string]interface{}, error) {
	sockPath := target.VppSockPath
	if sockPath == "" {
		sockPath = defaultVppSockPath
	}

	// Build vppctl command
	cmdArgs := []string{"-s", sockPath}
	cmdArgs = append(cmdArgs, strings.Fields(command)...)

	log.Printf("[Standalone] Executing: vppctl %s", strings.Join(cmdArgs, " "))

	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, vppctlPath, cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	execErr := cmd.Run()

	output := stdout.Bytes()
	errOutput := stderr.String()

	if errOutput != "" {
		log.Printf("[Standalone] stderr: %s", errOutput)
	}

	if execErr != nil {
		errorMsg := ""
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			errorMsg = string(exitErr.Stderr)
		}
		return map[string]interface{}{
			"success":   false,
			"error":     fmt.Sprintf("%v - %s", execErr, errorMsg),
			"mode":      string(ModeStandalone),
			"sock_path": sockPath,
			"command":   command,
		}, execErr
	}

	return map[string]interface{}{
		"success":   true,
		"output":    string(output),
		"command":   command,
		"mode":      string(ModeStandalone),
		"sock_path": sockPath,
	}, nil
}

// ExecuteShellCommand executes a shell command on the target (used for file operations, cleanup, etc.)
func ExecuteShellCommand(ctx context.Context, target *VPPTarget, shellCmd string) (string, error) {
	switch target.Mode {
	case ModeKubernetes:
		return executeKubernetesShellCommand(ctx, target, shellCmd)
	case ModeStandalone:
		return executeStandaloneShellCommand(ctx, shellCmd)
	default:
		return "", fmt.Errorf("unknown VPP execution mode: %s", target.Mode)
	}
}

// executeKubernetesShellCommand runs a shell command inside a Kubernetes pod
func executeKubernetesShellCommand(ctx context.Context, target *VPPTarget, shellCmd string) (string, error) {
	namespace := target.Namespace
	if namespace == "" {
		namespace = "calico-vpp-dataplane"
	}
	containerName := target.ContainerName
	if containerName == "" {
		containerName = "vpp"
	}

	cmdArgs := []string{
		"exec",
		"-n", namespace,
		target.PodName,
		"-c", containerName,
		"--",
		"sh", "-c", shellCmd,
	}

	log.Printf("[K8s Shell] Executing: kubectl %s", strings.Join(cmdArgs, " "))

	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "kubectl", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// executeStandaloneShellCommand runs a shell command on the local host
func executeStandaloneShellCommand(ctx context.Context, shellCmd string) (string, error) {
	log.Printf("[Standalone Shell] Executing: sh -c %s", shellCmd)

	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "sh", "-c", shellCmd)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// =============================================================================
// BGP COMMAND EXECUTION - KUBERNETES-ONLY, RUNS IN AGENT CONTAINER
// =============================================================================

// ExecutePodGoBGPCommand runs a gobgp command directly on a specified Kubernetes pod
func ExecutePodGoBGPCommand(ctx context.Context, podName, command string) (map[string]interface{}, error) {
	if podName == "" {
		return nil, fmt.Errorf("pod name is required")
	}

	namespace := "calico-vpp-dataplane"

	// Get the node name for the pod
	nodeName := ""
	k8sClient, err := newKubeClient()
	if err == nil {
		pod, err := k8sClient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err == nil {
			nodeName = pod.Spec.NodeName
		}
	}

	// Build kubectl command to execute in the agent container
	cmdArgs := []string{
		"exec",
		"-n", namespace,
		"-c", "agent", // Use the agent container
		podName,
		"--",
		"gobgp",
	}

	// Add the specific gobgp command arguments
	cmdArgs = append(cmdArgs, strings.Fields(command)...)

	// Execute the command with a timeout
	log.Printf("Executing command: kubectl %s", strings.Join(cmdArgs, " "))

	// Set a timeout for the command
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "kubectl", cmdArgs...)

	// Capture stdout and stderr separately
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("Starting command execution...")
	execErr := cmd.Run()
	log.Printf("Command completed with status: %v", execErr == nil)

	// Get the output
	output := stdout.Bytes()
	errOutput := stderr.String()

	if errOutput != "" {
		log.Printf("Command stderr: %s", errOutput)
	}

	if execErr != nil {
		errorMsg := ""
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			errorMsg = string(exitErr.Stderr)
		}
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("%v - %s", execErr, errorMsg),
			"node":    nodeName,
			"pod":     podName,
			"command": command,
		}, execErr
	}
	return map[string]interface{}{
		"success": true,
		"output":  string(output),
		"command": command,
		"node":    nodeName,
		"pod":     podName,
	}, nil
}

// =============================================================================
// UTILITY FUNCTIONS - INTERFACE MAPPING & VPP DRIVER DETECTION
// =============================================================================

// getVppDriverFromConfigMap retrieves the vppDriver from the calico-vpp-config ConfigMap
func getVppDriverFromConfigMap(k *KubeClient) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), k.timeout)
	defer cancel()

	configMap, err := k.clientset.CoreV1().ConfigMaps("calico-vpp-dataplane").Get(ctx, "calico-vpp-config", metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get calico-vpp-config ConfigMap: %v", err)
	}

	interfacesData, exists := configMap.Data["CALICOVPP_INTERFACES"]
	if !exists {
		return "", fmt.Errorf("CALICOVPP_INTERFACES not found in ConfigMap")
	}

	// Parse the JSON directly instead of using kubectl + jq
	var interfacesConfig struct {
		UplinkInterfaces []struct {
			VppDriver string `json:"vppDriver"`
		} `json:"uplinkInterfaces"`
	}

	err = json.Unmarshal([]byte(interfacesData), &interfacesConfig)
	if err != nil {
		return "", fmt.Errorf("failed to parse CALICOVPP_INTERFACES JSON: %v", err)
	}

	if len(interfacesConfig.UplinkInterfaces) == 0 {
		return "", fmt.Errorf("no uplink interfaces found in configuration")
	}

	driver := strings.TrimSpace(interfacesConfig.UplinkInterfaces[0].VppDriver)
	if driver == "" {
		return "", fmt.Errorf("vppDriver not found or is empty")
	}

	return driver, nil
}

// mapInterfaceTypeToVppInputNode maps interface types to VPP graph input nodes
func mapInterfaceTypeToVppInputNode(k *KubeClient, interfaceType string) (string, string, error) {
	switch interfaceType {
	case "phy":
		// Get the actual VPP driver from the ConfigMap
		actualDriver, err := getVppDriverFromConfigMap(k)
		if err != nil {
			return "", "", fmt.Errorf("failed to get VPP driver from ConfigMap: %v", err)
		}
		// Recursively call with the actual driver
		return mapInterfaceTypeToVppInputNode(k, actualDriver)
	case "af_xdp":
		return "af-xdp-input", "af_xdp", nil
	case "af_packet":
		return "af-packet-input", "af_packet", nil
	case "avf":
		return "avf-input", "avf", nil
	case "vmxnet3":
		return "vmxnet3-input", "vmxnet3", nil
	case "virtio", "tuntap":
		return "virtio-input", "virtio", nil
	case "rdma":
		return "rdma-input", "rdma", nil
	case "dpdk":
		return "dpdk-input", "dpdk", nil
	case "memif":
		return "memif-input", "memif", nil
	case "vcl":
		return "session-queue", "vcl", nil
	case "":
		return "virtio-input", "virtio", nil // default to tuntap (virtio)
	default:
		errorMsg := fmt.Sprintf("Invalid interface type: %s\n\nSupported interface types:\n", interfaceType)
		errorMsg += "  phy       : use the physical interface driver configured in calico-vpp-config\n"
		errorMsg += "  af_xdp    : use an AF_XDP socket to drive the interface\n"
		errorMsg += "  af_packet : use an AF_PACKET socket to drive the interface\n"
		errorMsg += "  avf       : use the VPP native driver for Intel 700-Series and 800-Series interfaces\n"
		errorMsg += "  vmxnet3   : use the VPP native driver for VMware virtual interfaces\n"
		errorMsg += "  virtio    : use the VPP native driver for Virtio virtual interfaces\n"
		errorMsg += "  tuntap    : alias for virtio (default)\n"
		errorMsg += "  rdma      : use the VPP native driver for Mellanox CX-4 and CX-5 interfaces\n"
		errorMsg += "  dpdk      : use the DPDK interface drivers with VPP\n"
		errorMsg += "  memif     : use shared memory interfaces (memif)\n"
		errorMsg += "  vcl       : capture packets at the session layer\n"
		errorMsg += "\nDefault: virtio (if no interface type is specified)"
		return "", "", fmt.Errorf("%s", errorMsg)
	}
}

// parseVppInterfaces parses the output of "vppctl show interface" and returns a list of up interfaces
func parseVppInterfaces(output string) []string {
	var upInterfaces []string
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// Skip empty lines and header lines
		if strings.TrimSpace(line) == "" || strings.Contains(line, "Name") || strings.Contains(line, "Counter") || strings.Contains(line, "Count") {
			continue
		}

		// Skip lines that don't start with an interface name (statistics lines, etc.)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "rx ") || strings.HasPrefix(trimmed, "tx ") ||
			strings.HasPrefix(trimmed, "drops") || strings.HasPrefix(trimmed, "punt") ||
			strings.HasPrefix(trimmed, "ip4") || strings.HasPrefix(trimmed, "ip6") {
			continue
		}

		// Look for interface lines (they start with interface name)
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			// Check if the line contains interface information
			// Format: "interface_name    idx    state    mtu"
			interfaceName := fields[0]
			state := fields[2]

			// Only add interfaces that are "up"
			if state == "up" && interfaceName != "" {
				upInterfaces = append(upInterfaces, interfaceName)
			}
		}
	}

	return upInterfaces
}

// mapInterfaceTypeToVppInputNodeStandalone maps interface types to VPP input nodes (for standalone mode)
func mapInterfaceTypeToVppInputNodeStandalone(interfaceType string) (string, string, error) {
	switch interfaceType {
	case "af_xdp":
		return "af-xdp-input", "af_xdp", nil
	case "af_packet":
		return "af-packet-input", "af_packet", nil
	case "avf":
		return "avf-input", "avf", nil
	case "vmxnet3":
		return "vmxnet3-input", "vmxnet3", nil
	case "virtio", "tuntap":
		return "virtio-input", "virtio", nil
	case "rdma":
		return "rdma-input", "rdma", nil
	case "dpdk":
		return "dpdk-input", "dpdk", nil
	case "memif":
		return "memif-input", "memif", nil
	case "vcl":
		return "session-queue", "vcl", nil
	case "":
		return "dpdk-input", "dpdk", nil // default to dpdk for standalone
	default:
		errorMsg := fmt.Sprintf("Invalid interface type: %s\n\nSupported interface types:\n", interfaceType)
		errorMsg += "  af_xdp    : use an AF_XDP socket to drive the interface\n"
		errorMsg += "  af_packet : use an AF_PACKET socket to drive the interface\n"
		errorMsg += "  avf       : use the VPP native driver for Intel 700-Series and 800-Series interfaces\n"
		errorMsg += "  vmxnet3   : use the VPP native driver for VMware virtual interfaces\n"
		errorMsg += "  virtio    : use the VPP native driver for Virtio virtual interfaces\n"
		errorMsg += "  tuntap    : alias for virtio\n"
		errorMsg += "  rdma      : use the VPP native driver for Mellanox CX-4 and CX-5 interfaces\n"
		errorMsg += "  dpdk      : use the DPDK interface drivers with VPP (default for standalone)\n"
		errorMsg += "  memif     : use shared memory interfaces (memif)\n"
		errorMsg += "  vcl       : capture packets at the session layer\n"
		return "", "", fmt.Errorf("%s", errorMsg)
	}
}

// getVppInputNode returns the VPP input node for a given interface type and target
func getVppInputNode(target *VPPTarget, interfaceType string, k8sClient *KubeClient) (string, string, error) {
	if target.Mode == ModeStandalone {
		return mapInterfaceTypeToVppInputNodeStandalone(interfaceType)
	}
	// For Kubernetes mode, use the original function that can look up ConfigMap
	return mapInterfaceTypeToVppInputNode(k8sClient, interfaceType)
}

// =============================================================================
// CAPTURE LOCK MECHANISM - PREVENTS PARALLEL CAPTURE OPERATIONS
// =============================================================================

// checkCaptureLock checks if there's an active capture lock and returns conflict info
func checkCaptureLock(ctx context.Context, target *VPPTarget) (*CaptureLockInfo, error) {
	checkCmd := fmt.Sprintf("test -f %s && cat %s", captureLockFile, captureLockFile)
	output, err := ExecuteShellCommand(ctx, target, checkCmd)
	if err != nil {
		// Lock file doesn't exist - no active capture
		return nil, nil
	}

	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	var lockInfo CaptureLockInfo
	err = json.Unmarshal([]byte(output), &lockInfo)
	if err != nil {
		// Invalid lock file - consider it stale
		log.Printf("Warning: Invalid lock file contents, treating as stale: %v", err)
		return nil, nil
	}

	// Check if lock is stale (older than 10 minutes)
	if time.Since(lockInfo.StartedAt) > 10*time.Minute {
		log.Printf("Warning: Found stale lock file (started %v ago), ignoring", time.Since(lockInfo.StartedAt))
		return nil, nil
	}

	return &lockInfo, nil
}

// createCaptureLock creates a lock file for the capture operation
func createCaptureLock(ctx context.Context, target *VPPTarget, operation, captureType string) error {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	targetName := "localhost"
	if target.Mode == ModeKubernetes {
		targetName = target.PodName
	}

	lockInfo := CaptureLockInfo{
		Operation:   operation,
		StartedAt:   time.Now(),
		Hostname:    hostname,
		Target:      targetName,
		CaptureType: captureType,
	}

	lockData, err := json.MarshalIndent(lockInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock info: %v", err)
	}

	// Create lock file
	createCmd := fmt.Sprintf("echo '%s' > %s", string(lockData), captureLockFile)
	_, err = ExecuteShellCommand(ctx, target, createCmd)
	if err != nil {
		return fmt.Errorf("failed to create lock file: %v", err)
	}

	return nil
}

// removeCaptureLock removes the capture lock file
func removeCaptureLock(ctx context.Context, target *VPPTarget) error {
	removeCmd := fmt.Sprintf("rm -f %s", captureLockFile)
	_, err := ExecuteShellCommand(ctx, target, removeCmd)
	if err != nil {
		return fmt.Errorf("failed to remove lock file: %v", err)
	}
	return nil
}

// cleanupCaptureFiles removes temporary capture files
func cleanupCaptureFiles(ctx context.Context, target *VPPTarget) error {
	cleanupCmd := "rm -f /tmp/trace.txt /tmp/trace.txt.gz /tmp/trace.pcap /tmp/trace.pcap.gz /tmp/dispatch.pcap /tmp/dispatch.pcap.gz"
	_, err := ExecuteShellCommand(ctx, target, cleanupCmd)
	return err
}

// =============================================================================
// CAPTURE HANDLERS - TRACE, PCAP, DISPATCH (SUPPORTS BOTH MODES WITH LOCKING)
// =============================================================================

// handleTraceCapture implements VPP trace capture with locking (supports both modes)
func (s *VPPMCPServer) handleTraceCapture(ctx context.Context, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received trace capture request: %s", string(inputJSON))

	target := input.GetTarget()

	// Validate input based on mode
	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Set mode='standalone' for local VPP daemon.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	// Acquire in-process lock
	captureMutex.Lock()
	defer captureMutex.Unlock()

	// Check for existing capture lock
	lockInfo, err := checkCaptureLock(ctx, target)
	if err != nil {
		log.Printf("Warning: Failed to check capture lock: %v", err)
	}
	if lockInfo != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: A capture operation is already running.\n\n"+
						"Active capture details:\n"+
						"- Operation: %s\n"+
						"- Type: %s\n"+
						"- Started: %s\n"+
						"- Started by: %s\n"+
						"- Target: %s\n\n"+
						"Use 'vpp_capture_cleanup' to force cleanup if the previous operation failed.",
						lockInfo.Operation, lockInfo.CaptureType,
						lockInfo.StartedAt.Format("2006-01-02 15:04:05"),
						lockInfo.Hostname, lockInfo.Target),
				},
			},
		}, nil, fmt.Errorf("capture already in progress")
	}

	// Create capture lock
	if err := createCaptureLock(ctx, target, "trace", "packet_trace"); err != nil {
		log.Printf("Warning: Failed to create capture lock: %v", err)
	}

	// Ensure cleanup on exit
	defer func() {
		if err := removeCaptureLock(ctx, target); err != nil {
			log.Printf("Warning: Failed to remove capture lock: %v", err)
		}
	}()

	// Get Kubernetes client for interface mapping (only needed for K8s mode with 'phy' interface)
	var k8sClient *KubeClient
	if target.Mode == ModeKubernetes {
		k8sClient, _ = newKubeClient()
	}

	// Map interface type to VPP input node
	vppInputNode, _, err := getVppInputNode(target, input.Interface, k8sClient)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error mapping interface: %v", err),
				},
			},
		}, nil, err
	}

	// Determine count and timeout
	count := input.Count
	if count == 0 {
		count = 500
	}
	timeout := input.Timeout
	if timeout == 0 {
		timeout = 30
	}

	// Build target info for display
	targetInfo := "localhost (standalone)"
	if target.Mode == ModeKubernetes {
		targetInfo = fmt.Sprintf("pod %s", target.PodName)
	}

	log.Printf("Starting trace capture on %s (count=%d, timeout=%ds, node=%s)", targetInfo, count, timeout, vppInputNode)

	// Step 1: Clear trace to ensure clean state
	_, _ = ExecuteVPPCommand(ctx, target, "clear trace")

	// Step 2: Start trace capture
	traceCmd := fmt.Sprintf("trace add %s %d", vppInputNode, count)
	_, err = ExecuteVPPCommand(ctx, target, traceCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture
	log.Printf("Capturing packets for %d seconds...", timeout)
	time.Sleep(time.Duration(timeout) * time.Second)

	// Step 4: Get trace results
	showTraceCmd := fmt.Sprintf("show trace max %d", count)
	result, err := ExecuteVPPCommand(ctx, target, showTraceCmd)
	if err != nil {
		// Still try to cleanup
		_, _ = ExecuteVPPCommand(ctx, target, "clear trace")
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error retrieving trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 5: Clear trace after retrieval
	_, _ = ExecuteVPPCommand(ctx, target, "clear trace")

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP Trace Capture Results:\n\n%s\n\n"+
						"Capture Parameters:\n"+
						"- VPP Input Node: %s\n"+
						"- Packet Count: %d\n"+
						"- Capture Duration: %d seconds\n"+
						"- Target: %s\n"+
						"- Mode: %s\n\n"+
						"**Note**: Trace output is returned directly (not saved to file)",
						output, vppInputNode, count, timeout, targetInfo, target.Mode),
				},
			},
		}
		return response, nil, nil
	}

	errorMsg := result["error"].(string)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Error executing trace capture: %s", errorMsg),
			},
		},
	}, nil, nil
}

// handlePcapCapture implements VPP pcap capture with locking (supports both modes)
func (s *VPPMCPServer) handlePcapCapture(ctx context.Context, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received pcap capture request: %s", string(inputJSON))

	target := input.GetTarget()

	// Validate input based on mode
	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Set mode='standalone' for local VPP daemon.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	// Acquire in-process lock
	captureMutex.Lock()
	defer captureMutex.Unlock()

	// Check for existing capture lock
	lockInfo, err := checkCaptureLock(ctx, target)
	if err != nil {
		log.Printf("Warning: Failed to check capture lock: %v", err)
	}
	if lockInfo != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: A capture operation is already running.\n\n"+
						"Active capture details:\n"+
						"- Operation: %s\n"+
						"- Type: %s\n"+
						"- Started: %s\n"+
						"- Started by: %s\n"+
						"- Target: %s\n\n"+
						"Use 'vpp_capture_cleanup' to force cleanup if the previous operation failed.",
						lockInfo.Operation, lockInfo.CaptureType,
						lockInfo.StartedAt.Format("2006-01-02 15:04:05"),
						lockInfo.Hostname, lockInfo.Target),
				},
			},
		}, nil, fmt.Errorf("capture already in progress")
	}

	// Create capture lock
	if err := createCaptureLock(ctx, target, "pcap", "packet_capture"); err != nil {
		log.Printf("Warning: Failed to create capture lock: %v", err)
	}

	// Ensure cleanup on exit
	defer func() {
		if err := removeCaptureLock(ctx, target); err != nil {
			log.Printf("Warning: Failed to remove capture lock: %v", err)
		}
	}()

	// Get list of available interfaces
	interfaceResult, err := ExecuteVPPCommand(ctx, target, "show int")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error getting interfaces: %v", err),
				},
			},
		}, nil, err
	}

	// Parse interfaces
	availableInterfaces := parseVppInterfaces(interfaceResult["output"].(string))

	// Validate interface if provided
	interfaceName := input.Interface
	if interfaceName == "" {
		interfaceName = "any"
	} else if interfaceName != "any" {
		found := false
		for _, iface := range availableInterfaces {
			if iface == interfaceName {
				found = true
				break
			}
		}
		if !found {
			var ifaceList strings.Builder
			ifaceList.WriteString("\nAvailable UP interfaces:")
			for i, iface := range availableInterfaces {
				ifaceList.WriteString(fmt.Sprintf("\n%d. %s", i+1, iface))
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error: Interface '%s' not found or is down.%s", interfaceName, ifaceList.String()),
					},
				},
			}, nil, fmt.Errorf("interface not found")
		}
	}

	// Determine count and timeout
	count := input.Count
	if count == 0 {
		count = 500
	}
	timeout := input.Timeout
	if timeout == 0 {
		timeout = 30
	}

	// Build target info for display
	targetInfo := "localhost (standalone)"
	if target.Mode == ModeKubernetes {
		targetInfo = fmt.Sprintf("pod %s", target.PodName)
	}

	// Clean up any pre-existing capture file
	_, _ = ExecuteShellCommand(ctx, target, "rm -f /tmp/trace.pcap /tmp/trace.pcap.gz")

	log.Printf("Starting pcap capture on %s (count=%d, timeout=%ds, interface=%s)", targetInfo, count, timeout, interfaceName)

	// Step 1: Stop any existing pcap capture
	_, _ = ExecuteVPPCommand(ctx, target, "pcap trace off")

	// Step 2: Start pcap capture
	pcapCmd := fmt.Sprintf("pcap trace tx rx max %d intfc %s file trace.pcap", count, interfaceName)
	_, err = ExecuteVPPCommand(ctx, target, pcapCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting pcap: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture
	log.Printf("Capturing packets for %d seconds...", timeout)
	time.Sleep(time.Duration(timeout) * time.Second)

	// Step 4: Stop pcap capture
	result, err := ExecuteVPPCommand(ctx, target, "pcap trace off")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error stopping pcap: %v", err),
				},
			},
		}, nil, err
	}

	// Determine pcap file location
	pcapFile := "/tmp/trace.pcap"

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP PCAP Capture Results:\n\n%s\n\n"+
						"Capture Parameters:\n"+
						"- Interface: %s\n"+
						"- Packet Count: %d\n"+
						"- Capture Duration: %d seconds\n"+
						"- Target: %s\n"+
						"- Mode: %s\n\n"+
						"**PCAP file saved at**: %s\n\n"+
						"To retrieve the file:\n"+
						"- Kubernetes: kubectl cp <namespace>/<pod>:%s ./capture.pcap -c vpp\n"+
						"- Standalone: The file is on the local filesystem",
						output, interfaceName, count, timeout, targetInfo, target.Mode, pcapFile, pcapFile),
				},
			},
		}
		return response, nil, nil
	}

	errorMsg := result["error"].(string)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Error executing pcap capture: %s", errorMsg),
			},
		},
	}, nil, nil
}

// handleDispatchCapture implements VPP dispatch trace capture with locking (supports both modes)
func (s *VPPMCPServer) handleDispatchCapture(ctx context.Context, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received dispatch capture request: %s", string(inputJSON))

	target := input.GetTarget()

	// Validate input based on mode
	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Set mode='standalone' for local VPP daemon.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	// Acquire in-process lock
	captureMutex.Lock()
	defer captureMutex.Unlock()

	// Check for existing capture lock
	lockInfo, err := checkCaptureLock(ctx, target)
	if err != nil {
		log.Printf("Warning: Failed to check capture lock: %v", err)
	}
	if lockInfo != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: A capture operation is already running.\n\n"+
						"Active capture details:\n"+
						"- Operation: %s\n"+
						"- Type: %s\n"+
						"- Started: %s\n"+
						"- Started by: %s\n"+
						"- Target: %s\n\n"+
						"Use 'vpp_capture_cleanup' to force cleanup if the previous operation failed.",
						lockInfo.Operation, lockInfo.CaptureType,
						lockInfo.StartedAt.Format("2006-01-02 15:04:05"),
						lockInfo.Hostname, lockInfo.Target),
				},
			},
		}, nil, fmt.Errorf("capture already in progress")
	}

	// Create capture lock
	if err := createCaptureLock(ctx, target, "dispatch", "dispatch_trace"); err != nil {
		log.Printf("Warning: Failed to create capture lock: %v", err)
	}

	// Ensure cleanup on exit
	defer func() {
		if err := removeCaptureLock(ctx, target); err != nil {
			log.Printf("Warning: Failed to remove capture lock: %v", err)
		}
	}()

	// Get Kubernetes client for interface mapping (only needed for K8s mode with 'phy' interface)
	var k8sClient *KubeClient
	if target.Mode == ModeKubernetes {
		k8sClient, _ = newKubeClient()
	}

	// Map interface type to VPP input node
	vppInputNode, _, err := getVppInputNode(target, input.Interface, k8sClient)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error mapping interface: %v", err),
				},
			},
		}, nil, err
	}

	// Determine count and timeout
	count := input.Count
	if count == 0 {
		count = 500
	}
	timeout := input.Timeout
	if timeout == 0 {
		timeout = 30
	}

	// Build target info for display
	targetInfo := "localhost (standalone)"
	if target.Mode == ModeKubernetes {
		targetInfo = fmt.Sprintf("pod %s", target.PodName)
	}

	// Clean up any pre-existing capture file
	_, _ = ExecuteShellCommand(ctx, target, "rm -f /tmp/dispatch.pcap /tmp/dispatch.pcap.gz")

	log.Printf("Starting dispatch capture on %s (count=%d, timeout=%ds, node=%s)", targetInfo, count, timeout, vppInputNode)

	// Step 1: Stop any existing dispatch trace
	_, _ = ExecuteVPPCommand(ctx, target, "pcap dispatch trace off")

	// Step 2: Start dispatch trace capture
	dispatchCmd := fmt.Sprintf("pcap dispatch trace on max %d buffer-trace %s %d file dispatch.pcap", count, vppInputNode, count)
	_, err = ExecuteVPPCommand(ctx, target, dispatchCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting dispatch trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture
	log.Printf("Capturing packets for %d seconds...", timeout)
	time.Sleep(time.Duration(timeout) * time.Second)

	// Step 4: Stop dispatch trace
	result, err := ExecuteVPPCommand(ctx, target, "pcap dispatch trace off")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error stopping dispatch trace: %v", err),
				},
			},
		}, nil, err
	}

	// Determine dispatch file location
	dispatchFile := "/tmp/dispatch.pcap"

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP Dispatch Trace Results:\n\n%s\n\n"+
						"Capture Parameters:\n"+
						"- VPP Input Node: %s\n"+
						"- Packet Count: %d\n"+
						"- Capture Duration: %d seconds\n"+
						"- Target: %s\n"+
						"- Mode: %s\n\n"+
						"**Dispatch PCAP file saved at**: %s\n\n"+
						"To retrieve the file:\n"+
						"- Kubernetes: kubectl cp <namespace>/<pod>:%s ./dispatch.pcap -c vpp\n"+
						"- Standalone: The file is on the local filesystem",
						output, vppInputNode, count, timeout, targetInfo, target.Mode, dispatchFile, dispatchFile),
				},
			},
		}
		return response, nil, nil
	}

	errorMsg := result["error"].(string)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Error executing dispatch trace: %s", errorMsg),
			},
		},
	}, nil, nil
}

// handleCaptureCleanup performs forced cleanup of all capture operations (supports both modes)
func (s *VPPMCPServer) handleCaptureCleanup(ctx context.Context, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received capture cleanup request: %s", string(inputJSON))

	target := input.GetTarget()

	// Validate input based on mode
	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Set mode='standalone' for local VPP daemon.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	// Build target info for display
	targetInfo := "localhost (standalone)"
	if target.Mode == ModeKubernetes {
		targetInfo = fmt.Sprintf("pod %s", target.PodName)
	}

	var results []string

	// Stop all captures
	log.Printf("Stopping all captures on %s", targetInfo)

	// Clear trace
	_, err := ExecuteVPPCommand(ctx, target, "clear trace")
	if err != nil {
		results = append(results, fmt.Sprintf("- Clear trace: FAILED (%v)", err))
	} else {
		results = append(results, "- Clear trace: OK")
	}

	// Stop PCAP
	_, err = ExecuteVPPCommand(ctx, target, "pcap trace off")
	if err != nil {
		results = append(results, fmt.Sprintf("- Stop PCAP trace: FAILED (%v)", err))
	} else {
		results = append(results, "- Stop PCAP trace: OK")
	}

	// Stop dispatch trace
	_, err = ExecuteVPPCommand(ctx, target, "pcap dispatch trace off")
	if err != nil {
		results = append(results, fmt.Sprintf("- Stop dispatch trace: FAILED (%v)", err))
	} else {
		results = append(results, "- Stop dispatch trace: OK")
	}

	// Clean up capture files
	err = cleanupCaptureFiles(ctx, target)
	if err != nil {
		results = append(results, fmt.Sprintf("- Clean up capture files: FAILED (%v)", err))
	} else {
		results = append(results, "- Clean up capture files: OK")
	}

	// Remove lock file
	err = removeCaptureLock(ctx, target)
	if err != nil {
		results = append(results, fmt.Sprintf("- Remove lock file: FAILED (%v)", err))
	} else {
		results = append(results, "- Remove lock file: OK")
	}

	response := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("VPP Capture Cleanup Results:\n\nTarget: %s\nMode: %s\n\nCleanup Actions:\n%s\n\n"+
					"The system is now ready for new capture operations.",
					targetInfo, target.Mode, strings.Join(results, "\n")),
			},
		},
	}
	return response, nil, nil
}

// =============================================================================
// NEW TOOL HANDLERS - 4-TOOL ARCHITECTURE
// =============================================================================

// executeKubectl runs a kubectl command and returns the result as an MCP response
func (s *VPPMCPServer) executeKubectl(ctx context.Context, cmdArgs []string) (*mcp.CallToolResult, any, error) {
	log.Printf("Executing: kubectl %s", strings.Join(cmdArgs, " "))
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "kubectl", cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	execErr := cmd.Run()
	output := strings.TrimSpace(stdout.String())
	errOutput := strings.TrimSpace(stderr.String())

	if execErr != nil {
		errorMsg := errOutput
		if errorMsg == "" {
			errorMsg = execErr.Error()
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing kubectl command: %s\nCommand: kubectl %s",
						errorMsg, strings.Join(cmdArgs, " ")),
				},
			},
		}, nil, nil
	}

	if output == "" {
		output = "(no output)"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%s\n\nCommand executed: kubectl %s", output, strings.Join(cmdArgs, " ")),
			},
		},
	}, nil, nil
}

// handleVppctl is the unified handler for all vppctl commands
func (s *VPPMCPServer) handleVppctl(ctx context.Context, input VppctlInput) (*mcp.CallToolResult, any, error) {
	command := strings.TrimSpace(input.Command)
	if command == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: 'command' is required. Specify the vppctl command to run (e.g., 'show version', 'show int addr').",
				},
			},
		}, nil, fmt.Errorf("command is required")
	}

	log.Printf("Received vppctl request: command=%q, mode=%s, pod=%s", command, input.Mode, input.PodName)

	// Route special capture commands to dedicated handlers
	lowerCmd := strings.ToLower(command)
	switch {
	case lowerCmd == "trace" || strings.HasPrefix(lowerCmd, "trace "):
		captureInput := VPPCaptureInput{
			Mode: input.Mode, PodName: input.PodName, SockPath: input.SockPath,
			Count: input.Count, Timeout: input.Timeout, Interface: input.Interface,
		}
		return s.handleTraceCapture(ctx, captureInput)

	case lowerCmd == "pcap" || strings.HasPrefix(lowerCmd, "pcap "):
		captureInput := VPPCaptureInput{
			Mode: input.Mode, PodName: input.PodName, SockPath: input.SockPath,
			Count: input.Count, Timeout: input.Timeout, Interface: input.Interface,
		}
		return s.handlePcapCapture(ctx, captureInput)

	case lowerCmd == "dispatch" || strings.HasPrefix(lowerCmd, "dispatch "):
		captureInput := VPPCaptureInput{
			Mode: input.Mode, PodName: input.PodName, SockPath: input.SockPath,
			Count: input.Count, Timeout: input.Timeout, Interface: input.Interface,
		}
		return s.handleDispatchCapture(ctx, captureInput)

	case lowerCmd == "capture_cleanup":
		cmdInput := VPPCommandInput{Mode: input.Mode, PodName: input.PodName, SockPath: input.SockPath}
		return s.handleCaptureCleanup(ctx, cmdInput)
	}

	// Normal vppctl command - pass through to VPP
	target := input.GetTarget()

	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Set mode='standalone' for local VPP daemon.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	log.Printf("Executing vppctl %s (mode: %s)", command, target.Mode)

	result, err := ExecuteVPPCommand(ctx, target, command)
	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
	}

	if result == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: no result returned for command: vppctl %s", command),
				},
			},
		}, nil, fmt.Errorf("no result from VPP command")
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		mode := result["mode"].(string)

		targetInfo := ""
		if mode == string(ModeKubernetes) {
			pod := result["pod"].(string)
			targetInfo = fmt.Sprintf("Pod: %s (container: vpp)", pod)
		} else {
			sockPath := result["sock_path"].(string)
			targetInfo = fmt.Sprintf("Mode: standalone, Socket: %s", sockPath)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("vppctl %s:\n\n%s\n\nCommand executed: vppctl %s\n%s",
						command, output, cmd, targetInfo),
				},
			},
		}, nil, nil
	}

	errorMsg := result["error"].(string)
	cmd := result["command"].(string)
	mode := result["mode"].(string)
	targetInfo := ""
	if mode == string(ModeKubernetes) {
		pod, _ := result["pod"].(string)
		targetInfo = fmt.Sprintf("pod %s", pod)
	} else {
		targetInfo = "standalone VPP"
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Error executing VPP command on %s: %s\nCommand attempted: vppctl %s",
					targetInfo, errorMsg, cmd),
			},
		},
	}, nil, nil
}

// handleGobgp is the unified handler for all gobgp commands
func (s *VPPMCPServer) handleGobgp(ctx context.Context, input GobgpInput) (*mcp.CallToolResult, any, error) {
	command := strings.TrimSpace(input.Command)
	podName := strings.TrimSpace(input.PodName)

	if command == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: 'command' is required. Specify the gobgp command to run (e.g., 'neighbor', 'global rib -a ipv4').",
				},
			},
		}, nil, fmt.Errorf("command is required")
	}

	if podName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: 'pod_name' is required. Specify the Kubernetes pod running the agent container with gobgp.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required")
	}

	log.Printf("Received gobgp request: command=%q, pod=%s", command, podName)

	result, err := ExecutePodGoBGPCommand(ctx, podName, command)
	if err != nil {
		log.Printf("Error executing gobgp command: %v", err)
	}

	if result == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: no result returned for command: gobgp %s", command),
				},
			},
		}, nil, fmt.Errorf("no result from gobgp command")
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		node := result["node"].(string)
		pod := result["pod"].(string)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("gobgp %s:\n\n%s\n\nCommand executed: gobgp %s\nNode: %s\nPod: %s (container: agent)",
						command, output, cmd, node, pod),
				},
			},
		}, nil, nil
	}

	errorMsg := result["error"].(string)
	cmd := result["command"].(string)
	node := result["node"].(string)
	pod, _ := result["pod"].(string)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Error executing gobgp command on node %s (pod: %s): %s\nCommand attempted: gobgp %s",
					node, pod, errorMsg, cmd),
			},
		},
	}, nil, nil
}

// handleCluster is the unified handler for Kubernetes cluster commands
func (s *VPPMCPServer) handleCluster(ctx context.Context, input ClusterInput) (*mcp.CallToolResult, any, error) {
	command := strings.TrimSpace(input.Command)
	if command == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: 'command' is required. Supported: get_pods, get_nodes, get_configmap, get_daemonset, get_events, get_deployments, get_services, get_endpoints, get_replicaset, get_ippool, get_namespaces, top_pods, top_nodes, logs, describe_pod, describe_node, exec.",
				},
			},
		}, nil, fmt.Errorf("command is required")
	}

	namespace := strings.TrimSpace(input.Namespace)
	if namespace == "" {
		namespace = "calico-vpp-dataplane"
	}

	log.Printf("Received cluster request: command=%q, namespace=%s", command, namespace)

	switch command {
	case "get_pods":
		return s.executeKubectl(ctx, []string{"get", "pods", "-n", namespace, "-owide"})

	case "get_nodes":
		return s.executeKubectl(ctx, []string{"get", "nodes", "-owide"})

	case "get_configmap":
		name := strings.TrimSpace(input.ResourceName)
		if name == "" {
			name = "calico-vpp-config"
		}
		return s.executeKubectl(ctx, []string{"get", "configmap", "-n", namespace, name, "-oyaml"})

	case "get_daemonset":
		name := strings.TrimSpace(input.ResourceName)
		if name == "" {
			name = "calico-vpp-node"
		}
		return s.executeKubectl(ctx, []string{"get", "daemonset", "-n", namespace, name, "-oyaml"})

	case "get_events":
		return s.executeKubectl(ctx, []string{"get", "events", "-n", namespace, "--sort-by=.metadata.creationTimestamp"})

	case "get_deployments":
		return s.executeKubectl(ctx, []string{"get", "deployments", "-n", namespace, "-owide"})

	case "get_services":
		return s.executeKubectl(ctx, []string{"get", "svc", "-n", namespace, "-owide"})

	case "get_endpoints":
		name := strings.TrimSpace(input.ResourceName)
		if name == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: 'resource_name' (service name) is required for get_endpoints."},
				},
			}, nil, fmt.Errorf("resource_name is required for get_endpoints")
		}
		return s.executeKubectl(ctx, []string{"get", "endpoints", "-n", namespace, name, "-oyaml"})

	case "get_replicaset":
		return s.executeKubectl(ctx, []string{"get", "replicaset", "-n", namespace, "-owide"})

	case "get_ippool":
		return s.executeKubectl(ctx, []string{"get", "ippool", "-oyaml"})

	case "get_namespaces":
		return s.executeKubectl(ctx, []string{"get", "ns"})

	case "top_pods":
		return s.executeKubectl(ctx, []string{"top", "pod", "-n", namespace})

	case "top_nodes":
		return s.executeKubectl(ctx, []string{"top", "node"})

	case "logs":
		podName := strings.TrimSpace(input.PodName)
		if podName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: 'pod_name' is required for the 'logs' command."},
				},
			}, nil, fmt.Errorf("pod_name is required for logs")
		}
		container := strings.TrimSpace(input.Container)
		if container == "" {
			container = "agent"
		}
		tailLines := input.TailLines
		if tailLines <= 0 {
			tailLines = 200
		}
		if tailLines > 5000 {
			tailLines = 5000
		}
		return s.executeKubectl(ctx, []string{"logs", "-n", namespace, podName, "-c", container, "--tail", strconv.Itoa(tailLines)})

	case "describe_pod":
		podName := strings.TrimSpace(input.PodName)
		if podName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: 'pod_name' is required for the 'describe_pod' command."},
				},
			}, nil, fmt.Errorf("pod_name is required for describe_pod")
		}
		return s.executeKubectl(ctx, []string{"describe", "pod", "-n", namespace, podName})

	case "describe_node":
		name := strings.TrimSpace(input.ResourceName)
		if name == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: 'resource_name' is required for the 'describe_node' command."},
				},
			}, nil, fmt.Errorf("resource_name is required for describe_node")
		}
		return s.executeKubectl(ctx, []string{"describe", "node", name})

	case "exec":
		podName := strings.TrimSpace(input.PodName)
		if podName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: 'pod_name' is required for the 'exec' command."},
				},
			}, nil, fmt.Errorf("pod_name is required for exec")
		}
		container := strings.TrimSpace(input.Container)
		if container == "" {
			container = "vpp"
		}
		execCommand := strings.TrimSpace(input.ResourceName)
		if execCommand == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: 'resource_name' (command to execute) is required for exec. Only read-only commands allowed (e.g., 'ls', 'cat /etc/config', 'ip a')."},
				},
			}, nil, fmt.Errorf("resource_name (command) is required for exec")
		}

		// Security check: block obviously mutating commands
		lowerCmd := strings.ToLower(execCommand)
		blockedPrefixes := []string{"rm", "delete", "edit", "apply", "create", "patch", "replace", "set", "update", "mv", "cp", "write", "echo >", ">", ">>", "|", "&", ";"}
		for _, prefix := range blockedPrefixes {
			if strings.HasPrefix(lowerCmd, prefix) || strings.Contains(lowerCmd, " "+prefix) {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Error: Command '%s' appears to be mutating and is blocked. Only read-only commands are allowed via exec.", execCommand)},
					},
				}, nil, fmt.Errorf("mutating command blocked: %s", execCommand)
			}
		}

		// Build exec command - note: we don't use -it in non-interactive contexts
		return s.executeKubectl(ctx, []string{"exec", "-n", namespace, podName, "-c", container, "--", "sh", "-c", execCommand})

	default:
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: unsupported cluster command '%s'.\nSupported commands: get_pods, get_nodes, get_configmap, get_daemonset, get_events, get_deployments, get_services, get_endpoints, get_replicaset, get_ippool, get_namespaces, top_pods, top_nodes, logs, describe_pod, describe_node, exec.", command),
				},
			},
		}, nil, fmt.Errorf("unsupported cluster command: %s", command)
	}
}

// handleList returns a comprehensive command reference for all tool categories
func (s *VPPMCPServer) handleList(ctx context.Context, input EmptyInput) (*mcp.CallToolResult, any, error) {
	log.Printf("Received list request")

	var sb strings.Builder

	sb.WriteString(`VPP MCP Server - Complete Command Reference
=============================================

This server provides 4 tools for debugging VPP in Kubernetes (CalicoVPP) or standalone environments.
Each tool accepts a 'command' parameter. The agent/LLM decides which commands to run based on the
troubleshooting context.

== TOOL: vppctl ==
Run any vppctl command on VPP. Pass the full command string as 'command'.
Supports both Kubernetes (pod_name required) and standalone (mode='standalone') modes.

Parameters:
  command   (required) - The vppctl command to run
  pod_name  (required for Kubernetes mode) - Pod running VPP
  mode      (optional) - 'kubernetes' (default) or 'standalone'
  sock_path (optional) - VPP socket path for standalone mode (default: /var/run/vpp/cli.sock)

For packet capture commands (trace, pcap, dispatch, capture_cleanup):
  count     (optional) - Number of packets to capture (default: 500)
  timeout   (optional) - Capture duration in seconds (default: 30)
  interface (optional) - Interface type or name for captures

--- COMPLETE VPPCTL COMMAND CATALOG ---
`)
	sb.WriteString(vppctlCommandCatalog)

	sb.WriteString(`
== TOOL: gobgp ==
Run any gobgp command in the agent container of a CalicoVPP pod.
Kubernetes-only. Requires pod_name.

Parameters:
  command  (required) - The gobgp command to run
  pod_name (required) - Pod running the agent container

Complete gobgp command reference:

Neighbor commands:
  neighbor                        - List all BGP peers with state summary
  neighbor <ip>                   - Detailed info for a specific peer
  neighbor <ip> adj-in            - Routes received from a peer
  neighbor <ip> adj-out           - Routes advertised to a peer
  neighbor <ip> reset             - Reset peer session
  neighbor <ip> softreset         - Soft reset peer
  neighbor <ip> softresetin       - Soft reset inbound routes
  neighbor <ip> softresetout      - Soft reset outbound routes

Global / RIB commands:
  global                          - Show global BGP configuration (AS, router-id)
  global rib                      - Show full global RIB
  global rib -a ipv4              - IPv4 unicast RIB
  global rib -a ipv6              - IPv6 unicast RIB
  global rib -a ipv4 -l           - Longer prefixes in IPv4 RIB
  global rib -a ipv6 -l           - Longer prefixes in IPv6 RIB
  global rib <prefix>             - Lookup specific prefix in RIB
  global rib -a ipv4 summary      - IPv4 RIB summary
  global rib -a ipv6 summary      - IPv6 RIB summary

VRF commands:
  vrf                             - Show all VRFs
  vrf <name> rib                  - Show VRF-specific RIB
  vrf <name> rib -a ipv4          - Show VRF IPv4 RIB
  vrf <name> rib -a ipv6          - Show VRF IPv6 RIB

Policy commands:
  policy prefix                   - Show prefix sets
  policy neighbor                 - Show neighbor sets
  policy community                - Show community sets
  policy ext-community            - Show extended community sets
  policy as-path                  - Show AS-path sets
  policy statement                - Show policy statements
  policy                          - Show all policies

Monitor commands:
  monitor neighbor                - Watch BGP neighbor state changes in real-time
  monitor global rib              - Watch global RIB changes in real-time

Output interpretation:
  - Established peerings show as 'Establ'
  - Unsuccessful connections show as 'Opened' with 0 in Received/Accepted
  - 'Idle' means peer is configured but not attempting connection

== TOOL: cluster ==
Run Kubernetes cluster commands for CalicoVPP troubleshooting.
Kubernetes-only.

Parameters:
  command       (required) - The cluster operation to run
  namespace     (optional) - Kubernetes namespace (default: calico-vpp-dataplane)
  pod_name      (optional) - Pod name (for logs, describe_pod, exec)
  container     (optional) - Container name (for logs, exec; default varies by command)
  resource_name (optional) - Resource name or command to execute (for exec)
  tail_lines    (optional) - Number of log lines to fetch (default: 200, max: 5000)

Complete cluster command reference:

  get_pods              - List CalicoVPP pods (status, IPs, nodes)
                          → kubectl get pods -n <namespace> -owide
  get_nodes             - List all cluster nodes
                          → kubectl get nodes -owide
  get_configmap         - Get a ConfigMap (resource_name default: calico-vpp-config)
                          → kubectl get configmap -n <namespace> <name> -oyaml
  get_daemonset         - Get DaemonSet details in YAML (resource_name default: calico-vpp-node)
                          → kubectl get daemonset -n <namespace> <name> -oyaml
  get_events            - Get events sorted by creation timestamp
                          → kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
  get_deployments       - List deployments
                          → kubectl get deployments -n <namespace> -owide
  get_services          - List services
                          → kubectl get svc -n <namespace> -owide
  get_endpoints         - Get endpoints for a service (requires resource_name)
                          → kubectl get endpoints -n <namespace> <service> -oyaml
  get_replicaset        - List replicasets
                          → kubectl get replicaset -n <namespace> -owide
  get_ippool            - Get IPPool resources (Calico CRD)
                          → kubectl get ippool -oyaml
  get_namespaces        - List all namespaces
                          → kubectl get ns
  top_pods              - Resource usage for pods
                          → kubectl top pod -n <namespace>
  top_nodes             - Resource usage for nodes
                          → kubectl top node
  logs                  - Get container logs (requires pod_name; container default: agent)
                          → kubectl logs -n <namespace> <pod> -c <container> --tail <n>
  describe_pod          - Describe a pod (requires pod_name)
                          → kubectl describe pod -n <namespace> <pod>
  describe_node         - Describe a node (requires resource_name)
                          → kubectl describe node <name>
  exec                  - Execute read-only command in pod (requires pod_name, resource_name=command)
                          → kubectl exec -n <namespace> <pod> -c <container> -- sh -c <command>
                          NOTE: Only read-only commands allowed. Mutating commands are blocked.
`)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: sb.String()},
		},
	}, nil, nil
}

// =============================================================================
// TOOL REGISTRATION & MAIN
// =============================================================================

func main() {
	// Parse command-line flags
	transportMode := flag.String("transport", "stdio", "Transport mode: stdio or http")
	port := flag.String("port", "8080", "HTTP port (only used when transport=http)")
	flag.Parse()

	log.Printf("Starting VPP MCP Server with transport=%s...", *transportMode)

	// Create the VPP MCP server instance
	vppServer := NewVPPMCPServer()

	// Create MCP server with implementation info
	impl := &mcp.Implementation{
		Name:    "vpp-mcp-server",
		Version: "2.0.0",
	}

	vppServer.server = mcp.NewServer(impl, nil)

	// =========================================================================
	// TOOL 1: list - Returns a comprehensive command reference
	// =========================================================================

	toolList := &mcp.Tool{
		Name:        "list",
		Description: "Returns a comprehensive reference of all available commands across all tool categories (vppctl, gobgp, cluster). Use this tool to discover what commands you can run.",
	}
	mcp.AddTool(vppServer.server, toolList, func(ctx context.Context, req *mcp.CallToolRequest, input EmptyInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleList(ctx, input)
	})

	// =========================================================================
	// TOOL 2: cluster - Kubernetes cluster commands
	// =========================================================================

	toolCluster := &mcp.Tool{
		Name: "cluster",
		Description: "Run Kubernetes cluster commands for CalicoVPP troubleshooting. Kubernetes-only.\n\n" +
			"Required parameters:\n" +
			"- command: The cluster operation to run\n\n" +
			"Supported commands:\n" +
			"- get_pods: List CalicoVPP pods with status, IPs, and node placement\n" +
			"- get_nodes: List all cluster nodes\n" +
			"- get_configmap: Get a ConfigMap (resource_name default: calico-vpp-config)\n" +
			"- get_daemonset: Get DaemonSet YAML (resource_name default: calico-vpp-node)\n" +
			"- get_events: Get events sorted by creation timestamp\n" +
			"- get_deployments: List deployments\n" +
			"- get_services: List services\n" +
			"- get_endpoints: Get endpoints for a service (requires resource_name)\n" +
			"- get_replicaset: List replicasets\n" +
			"- get_ippool: Get IPPool resources (Calico CRD)\n" +
			"- get_namespaces: List all namespaces\n" +
			"- top_pods: Resource usage for pods\n" +
			"- top_nodes: Resource usage for nodes\n" +
			"- logs: Get container logs (requires pod_name; container default: agent; tail_lines default: 200)\n" +
			"- describe_pod: Describe a pod (requires pod_name)\n" +
			"- describe_node: Describe a node (requires resource_name)\n" +
			"- exec: Execute read-only command in pod (requires pod_name, resource_name=command; container default: vpp)\n\n" +
			"Optional parameters:\n" +
			"- namespace: Kubernetes namespace (default: calico-vpp-dataplane)\n" +
			"- pod_name: Pod name (for logs, describe_pod, exec)\n" +
			"- container: Container name (for logs, exec; default varies by command)\n" +
			"- resource_name: Resource name (for get_configmap, get_daemonset, get_endpoints, describe_node) or command to execute (for exec)\n" +
			"- tail_lines: Number of log lines to fetch (default: 200, max: 5000)",
	}
	mcp.AddTool(vppServer.server, toolCluster, func(ctx context.Context, req *mcp.CallToolRequest, input ClusterInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleCluster(ctx, input)
	})

	// =========================================================================
	// TOOL 3: vppctl - Run any vppctl command on VPP
	// =========================================================================

	toolVppctl := &mcp.Tool{
		Name: "vppctl",
		Description: "Run any vppctl command on VPP. Passes the command string directly to the VPP CLI.\n\n" +
			"Supports both Kubernetes (CalicoVPP pods) and standalone (local VPP daemon) modes.\n\n" +
			"Required parameters:\n" +
			"- command: The vppctl command to run (e.g., 'show version', 'show int addr', 'show ip fib index 0 10.0.0.0/24')\n\n" +
			"Optional parameters:\n" +
			"- pod_name: (Kubernetes mode, required) The name of the Kubernetes pod running VPP\n" +
			"- mode: 'kubernetes' (default) or 'standalone'\n" +
			"- sock_path: (Standalone mode) VPP socket path (default: /var/run/vpp/cli.sock)\n\n" +
			"For packet capture commands (trace, pcap, dispatch, capture_cleanup), additional optional parameters:\n" +
			"- count: Number of packets to capture (default: 500)\n" +
			"- timeout: Capture duration in seconds (default: 30)\n" +
			"- interface: Interface type or name for captures\n\n" +
			"Common commands:\n" +
			"  show version, show int, show int addr, show hardware-interfaces,\n" +
			"  show errors, show run, clear errors, clear run, show logging,\n" +
			"  show ip table, show ip6 table, show ip fib [index <n>] [<prefix>],\n" +
			"  show ip6 fib [index <n>] [<prefix>], show ip neighbor,\n" +
			"  show session verbose 2, show session stats, show tcp stats,\n" +
			"  show npol rules/policies/ipset/interfaces,\n" +
			"  show cnat translation, show cnat session,\n" +
			"  show ipip tunnel, show vxlan tunnel, show tun,\n" +
			"  trace, pcap, dispatch, capture_cleanup\n\n" +
			"Use the 'list' tool for a complete command reference.",
	}
	mcp.AddTool(vppServer.server, toolVppctl, func(ctx context.Context, req *mcp.CallToolRequest, input VppctlInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVppctl(ctx, input)
	})

	// =========================================================================
	// TOOL 4: gobgp - Run any gobgp command
	// =========================================================================

	toolGobgp := &mcp.Tool{
		Name: "gobgp",
		Description: "Run any gobgp command in the agent container of a CalicoVPP pod. Kubernetes-only.\n\n" +
			"Required parameters:\n" +
			"- command: The gobgp command to run (e.g., 'neighbor', 'global rib -a ipv4')\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container\n\n" +
			"Common commands:\n" +
			"  neighbor                  - List all BGP peers\n" +
			"  neighbor <ip>             - Detailed peer info\n" +
			"  neighbor <ip> adj-in      - Routes received from peer\n" +
			"  neighbor <ip> adj-out     - Routes advertised to peer\n" +
			"  global                    - Global BGP config (AS, router-id)\n" +
			"  global rib -a ipv4        - IPv4 unicast RIB\n" +
			"  global rib -a ipv6        - IPv6 unicast RIB\n" +
			"  global rib <prefix>       - Lookup specific prefix\n" +
			"  vrf                       - Show all VRFs\n" +
			"  policy prefix/community/as-path/statement\n\n" +
			"Output interpretation:\n" +
			"- Established peerings show as 'Establ'\n" +
			"- Unsuccessful connections show as 'Opened' with 0 in Received/Accepted\n\n" +
			"Use the 'list' tool for a complete command reference.",
	}
	mcp.AddTool(vppServer.server, toolGobgp, func(ctx context.Context, req *mcp.CallToolRequest, input GobgpInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleGobgp(ctx, input)
	})

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Choose transport based on flag
	switch *transportMode {
	case "stdio":
		log.Println("Using stdio transport...")
		runStdioTransport(ctx, vppServer)

	case "http":
		log.Printf("Using HTTP transport on port %s...", *port)
		runHTTPTransport(ctx, vppServer, *port, sigChan)

	default:
		log.Fatalf("Invalid transport mode: %s. Use 'stdio' or 'http'", *transportMode)
	}
}

// runStdioTransport runs the server with stdio transport
func runStdioTransport(ctx context.Context, vppServer *VPPMCPServer) {
	// Create stdio transport and connect
	transport := &mcp.StdioTransport{}

	// Connect the server
	log.Println("Connecting MCP server...")
	session, err := vppServer.server.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect server: %v", err)
	}
	log.Println("MCP server connected successfully")
	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}()

	// Wait for the session to complete
	log.Println("Waiting for session to complete...")
	if err := session.Wait(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	log.Println("Session completed")
}

// runHTTPTransport runs the server with HTTP/SSE transport
func runHTTPTransport(ctx context.Context, vppServer *VPPMCPServer, port string, sigChan chan os.Signal) {
	// Create HTTP server with SSE handler
	mux := http.NewServeMux()

	// MCP SSE endpoint - use NewSSEHandler for automatic session management
	sseHandler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		log.Printf("New SSE connection from %s", r.RemoteAddr)
		return vppServer.server
	}, &mcp.SSEOptions{})

	mux.Handle("/sse", sseHandler)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})

	// Root endpoint with info
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head><title>VPP MCP Server</title></head>
<body>
	<h1>VPP MCP Server</h1>
	<p>This is a Model Context Protocol (MCP) server for VPP debugging.</p>
	<h2>Endpoints:</h2>
	<ul>
		<li><strong>/sse</strong> - MCP SSE endpoint for client connections</li>
		<li><strong>/health</strong> - Health check endpoint</li>
	</ul>
	<p>Use an MCP client to connect to the /sse endpoint.</p>
</body>
</html>`
		_, err := w.Write([]byte(html))
		if err != nil {
			log.Printf("Error writing HTML response: %v", err)
		}
	})

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("HTTP server listening on port %s", port)
		log.Printf("MCP SSE endpoint: http://localhost:%s/sse", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, gracefully shutting down...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	log.Println("Server shutdown complete")
}
