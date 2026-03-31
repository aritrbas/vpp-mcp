package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
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

// VPPFIBInput represents the input for VPP FIB tools requiring fib_index
type VPPFIBInput struct {
	// Mode specifies how VPP is running: "kubernetes" (default) or "standalone"
	Mode string `json:"mode,omitempty"`
	// PodName specifies the name of the Kubernetes pod running VPP (required for kubernetes mode)
	PodName string `json:"pod_name,omitempty"`
	// SockPath specifies the VPP socket path for standalone mode (default: /var/run/vpp/cli.sock)
	SockPath string `json:"sock_path,omitempty"`
	// FibIndex specifies the FIB table index
	FibIndex string `json:"fib_index"`
}

// GetTarget creates a VPPTarget from the input
func (v *VPPFIBInput) GetTarget() *VPPTarget {
	mode := strings.ToLower(v.Mode)
	if mode == "standalone" || mode == "local" || mode == "daemon" {
		return NewStandaloneTarget(v.SockPath)
	}
	return NewKubernetesTarget(v.PodName)
}

// VPPFIBPrefixInput represents the input for VPP FIB tools requiring fib_index and prefix
type VPPFIBPrefixInput struct {
	// Mode specifies how VPP is running: "kubernetes" (default) or "standalone"
	Mode string `json:"mode,omitempty"`
	// PodName specifies the name of the Kubernetes pod running VPP (required for kubernetes mode)
	PodName string `json:"pod_name,omitempty"`
	// SockPath specifies the VPP socket path for standalone mode (default: /var/run/vpp/cli.sock)
	SockPath string `json:"sock_path,omitempty"`
	// FibIndex specifies the FIB table index
	FibIndex string `json:"fib_index"`
	// Prefix specifies the IP prefix to query
	Prefix string `json:"prefix"`
}

// GetTarget creates a VPPTarget from the input
func (v *VPPFIBPrefixInput) GetTarget() *VPPTarget {
	mode := strings.ToLower(v.Mode)
	if mode == "standalone" || mode == "local" || mode == "daemon" {
		return NewStandaloneTarget(v.SockPath)
	}
	return NewKubernetesTarget(v.PodName)
}

// BGPCommandInput represents the input for BGP command tools
type BGPCommandInput struct {
	// PodName specifies the name of the Kubernetes pod running the agent container with gobgp
	PodName string `json:"pod_name"`
}

// BGPParameterCommandInput represents the input for BGP command tools that require a parameter (IP, prefix, or neighbor IP)
type BGPParameterCommandInput struct {
	// PodName specifies the name of the Kubernetes pod running the agent container with gobgp
	PodName string `json:"pod_name"`
	// Parameter specifies the parameter value (IP address, prefix, or neighbor IP)
	Parameter string `json:"parameter"`
}

// EmptyInput represents tools that don't require any input parameters
type EmptyInput struct{}

// VPPDaemonsetImageInput represents input parameters for daemonset image lookup
type VPPDaemonsetImageInput struct {
	// Namespace specifies the Kubernetes namespace (default: calico-vpp-dataplane)
	Namespace string `json:"namespace,omitempty"`
	// DaemonsetName specifies the daemonset name (default: calico-vpp-node)
	DaemonsetName string `json:"daemonset_name,omitempty"`
	// ContainerName specifies the container name in the daemonset pod spec (default: vpp)
	ContainerName string `json:"container_name,omitempty"`
}

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

// handleTunnelInterface emulates `vppctl show tun | grep <tunX> -A 40` by locating the
// requested tunnel interface in the `show tun` output and returning the next 40 lines.
func (s *VPPMCPServer) handleTunnelInterface(ctx context.Context, input VPPTunnelInterfaceInput) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received tunnel interface inspection request: %s", string(inputJSON))

	podName := strings.TrimSpace(input.PodName)
	if podName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: pod_name is required. Please specify the Kubernetes pod running VPP."},
			},
		}, nil, fmt.Errorf("pod_name is required")
	}

	interfaceName := strings.TrimSpace(input.InterfaceName)
	if interfaceName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: interface_name is required. Please provide the tunnel interface name (for example, tun1)."},
			},
		}, nil, fmt.Errorf("interface_name is required")
	}

	target := NewKubernetesTarget(podName)
	result, err := ExecuteVPPCommand(ctx, target, "show tun")
	if err != nil {
		log.Printf("Error executing 'show tun' on pod %s: %v", podName, err)
	}

	if result == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to retrieve tunnel information for pod %s.", podName)},
			},
		}, nil, err
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		errorMsg, _ := result["error"].(string)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error executing VPP command on pod %s: %s", podName, errorMsg)},
			},
		}, nil, err
	}

	output := result["output"].(string)
	lines := strings.Split(output, "\n")
	needle := fmt.Sprintf("Interface: %s", interfaceName)
	start := -1
	for idx, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), needle) {
			start = idx
			break
		}
	}

	if start == -1 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("No tunnel interface named %s found in pod %s.", interfaceName, podName)},
			},
		}, nil, nil
	}

	end := start + 41 // include matching line plus the next 40 lines (grep -A 40)
	if end > len(lines) {
		end = len(lines)
	}
	snippet := strings.Join(lines[start:end], "\n")

	responseText := fmt.Sprintf(
		"VPP Tunnel Interface %s:\n\n%s\n\nCommand executed: vppctl show tun | grep '%s' -A 40 (emulated)\nPod: %s (container: vpp)",
		interfaceName,
		snippet,
		interfaceName,
		podName,
	)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

// handleHardwareInterface runs `vppctl show hardware-interfaces <iface>` to inspect a specific interface
func (s *VPPMCPServer) handleHardwareInterface(ctx context.Context, input VPPInterfaceInput) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received hardware interface inspection request: %s", string(inputJSON))

	podName := strings.TrimSpace(input.PodName)
	if podName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: pod_name is required. Please specify the Kubernetes pod running VPP."},
			},
		}, nil, fmt.Errorf("pod_name is required")
	}

	interfaceName := strings.TrimSpace(input.InterfaceName)
	if interfaceName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: interface_name is required. Please provide the interface name (for example, tun1)."},
			},
		}, nil, fmt.Errorf("interface_name is required")
	}

	command := fmt.Sprintf("show hardware-interfaces %s", interfaceName)
	target := NewKubernetesTarget(podName)
	result, err := ExecuteVPPCommand(ctx, target, command)
	if err != nil {
		log.Printf("Error executing '%s' on pod %s: %v", command, podName, err)
	}

	if result == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to retrieve hardware interface information for pod %s.", podName)},
			},
		}, nil, err
	}

	if success, ok := result["success"].(bool); ok && success {
		output, _ := result["output"].(string)
		responseText := fmt.Sprintf(
			"VPP Hardware Interface %s:\n\n%s\n\nCommand executed: vppctl %s\nPod: %s (container: vpp)",
			interfaceName,
			output,
			command,
			podName,
		)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
		}, nil, nil
	}

	errorMsg, _ := result["error"].(string)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Error executing VPP command on pod %s: %s\nCommand attempted: vppctl %s", podName, errorMsg, command)},
		},
	}, nil, err
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

// VPPTunnelInterfaceInput represents the input for VPP tunnel interface tools
type VPPTunnelInterfaceInput struct {
	// PodName specifies the name of the Kubernetes pod running VPP
	PodName string `json:"pod_name"`
	// InterfaceName specifies the tunnel interface to inspect (e.g., tun1)
	InterfaceName string `json:"interface_name"`
}

// VPPInterfaceInput represents input for tools that target a specific interface
type VPPInterfaceInput struct {
	// PodName specifies the name of the Kubernetes pod running VPP
	PodName string `json:"pod_name"`
	// InterfaceName specifies the interface to inspect (for example, GigabitEthernet0/8/0 or tun1)
	InterfaceName string `json:"interface_name"`
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
// VPP COMMAND HANDLERS
// =============================================================================

// handleVPPCommand is a generic handler for VPP commands (supports both Kubernetes and Standalone modes)
func (s *VPPMCPServer) handleVPPCommand(ctx context.Context, input VPPCommandInput, command, commandDescription string) (*mcp.CallToolResult, any, error) {
	// Log the request details
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received %s request with input: %s", commandDescription, string(inputJSON))

	// Get the target based on mode
	target := input.GetTarget()

	// Validate input based on mode
	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Please specify the Kubernetes pod name running VPP, or set mode='standalone' for local VPP daemon.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	log.Printf("Executing vppctl %s command (mode: %s)", command, target.Mode)

	// Execute the VPP command
	result, err := ExecuteVPPCommand(ctx, target, command)

	log.Printf("Command execution completed, processing results...")
	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		mode := result["mode"].(string)

		// Build target info based on mode
		targetInfo := ""
		if mode == string(ModeKubernetes) {
			pod := result["pod"].(string)
			targetInfo = fmt.Sprintf("Pod: %s (container: vpp)", pod)
		} else {
			sockPath := result["sock_path"].(string)
			targetInfo = fmt.Sprintf("Mode: standalone, Socket: %s", sockPath)
		}

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: vppctl %s\n%s",
						commandDescription, output, cmd, targetInfo),
				},
			},
		}

		log.Println("Successfully executed VPP command, returning result")
		return response, nil, nil
	} else {
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

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing VPP command on %s: %s\nCommand attempted: vppctl %s",
						targetInfo, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing VPP command on %s: %s", targetInfo, errorMsg)
		return errorResponse, nil, nil
	}
}

// handleVPPFIBCommand is a handler for VPP FIB commands that require fib_index (supports both modes)
func (s *VPPMCPServer) handleVPPFIBCommand(ctx context.Context, input VPPFIBInput, commandTemplate, commandDescription string) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received %s request with input: %s", commandDescription, string(inputJSON))

	target := input.GetTarget()

	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Set mode='standalone' for local VPP.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	if input.FibIndex == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: fib_index is required. Please specify the FIB table index.",
				},
			},
		}, nil, fmt.Errorf("fib_index is required")
	}

	// Build the command with fib_index
	command := fmt.Sprintf(commandTemplate, input.FibIndex)
	log.Printf("Executing vppctl %s command (mode: %s)", command, target.Mode)

	result, err := ExecuteVPPCommand(ctx, target, command)

	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
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

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: vppctl %s\n%s",
						commandDescription, output, cmd, targetInfo),
				},
			},
		}

		log.Println("Successfully executed VPP FIB command, returning result")
		return response, nil, nil
	} else {
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

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing VPP command on %s: %s\nCommand attempted: vppctl %s",
						targetInfo, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing VPP FIB command on %s: %s", targetInfo, errorMsg)
		return errorResponse, nil, nil
	}
}

// handleVPPFIBPrefixCommand is a handler for VPP FIB commands that require fib_index and prefix (supports both modes)
func (s *VPPMCPServer) handleVPPFIBPrefixCommand(ctx context.Context, input VPPFIBPrefixInput, commandTemplate, commandDescription string) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received %s request with input: %s", commandDescription, string(inputJSON))

	target := input.GetTarget()

	if target.Mode == ModeKubernetes && target.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: pod_name is required for Kubernetes mode. Set mode='standalone' for local VPP.",
				},
			},
		}, nil, fmt.Errorf("pod_name is required for Kubernetes mode")
	}

	if input.FibIndex == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: fib_index is required. Please specify the FIB table index.",
				},
			},
		}, nil, fmt.Errorf("fib_index is required")
	}

	if input.Prefix == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: prefix is required. Please specify the IP prefix.",
				},
			},
		}, nil, fmt.Errorf("prefix is required")
	}

	// Build the command with fib_index and prefix
	command := fmt.Sprintf(commandTemplate, input.FibIndex, input.Prefix)
	log.Printf("Executing vppctl %s command (mode: %s)", command, target.Mode)

	result, err := ExecuteVPPCommand(ctx, target, command)

	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
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

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: vppctl %s\n%s",
						commandDescription, output, cmd, targetInfo),
				},
			},
		}

		log.Println("Successfully executed VPP FIB prefix command, returning result")
		return response, nil, nil
	} else {
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

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing VPP command on %s: %s\nCommand attempted: vppctl %s",
						targetInfo, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing VPP FIB prefix command on %s: %s", targetInfo, errorMsg)
		return errorResponse, nil, nil
	}
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
// KUBERNETES-SPECIFIC HANDLERS
// =============================================================================

// handleGetPods implements listing all calico-vpp pods with IPs and nodes
func (s *VPPMCPServer) handleGetPods(ctx context.Context, input EmptyInput) (*mcp.CallToolResult, any, error) {
	log.Printf("Received vpp_get_pods request")

	// Execute kubectl command to get pods with wide output
	cmdArgs := []string{
		"get", "pods",
		"-n", "calico-vpp-dataplane",
		"-owide",
	}

	log.Printf("Executing command: kubectl %s", strings.Join(cmdArgs, " "))

	// Set a timeout for the command
	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "kubectl", cmdArgs...)

	// Capture stdout and stderr separately
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	execErr := cmd.Run()

	// Get the output
	output := stdout.String()
	errOutput := stderr.String()

	if errOutput != "" {
		log.Printf("Command stderr: %s", errOutput)
	}

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

	response := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Calico VPP Pods:\n\n%s\n\nCommand executed: kubectl %s",
					output, strings.Join(cmdArgs, " ")),
			},
		},
	}

	log.Println("Successfully executed kubectl command, returning result")
	return response, nil, nil
}

// handleShowDaemonsetImage returns the image configured for a container in a daemonset
func (s *VPPMCPServer) handleShowDaemonsetImage(ctx context.Context, input VPPDaemonsetImageInput) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received vpp_show_daemonset_image request: %s", string(inputJSON))

	namespace := strings.TrimSpace(input.Namespace)
	if namespace == "" {
		namespace = "calico-vpp-dataplane"
	}
	daemonsetName := strings.TrimSpace(input.DaemonsetName)
	if daemonsetName == "" {
		daemonsetName = "calico-vpp-node"
	}
	containerName := strings.TrimSpace(input.ContainerName)
	if containerName == "" {
		containerName = "vpp"
	}

	jsonPath := fmt.Sprintf(`{.spec.template.spec.containers[?(@.name=="%s")].image}`, containerName)
	cmdArgs := []string{
		"get", "daemonset",
		"-n", namespace,
		daemonsetName,
		"-o", "jsonpath=" + jsonPath,
	}

	log.Printf("Executing command: kubectl %s", strings.Join(cmdArgs, " "))

	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "kubectl", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	execErr := cmd.Run()
	output := strings.TrimSpace(stdout.String())
	errOutput := strings.TrimSpace(stderr.String())

	if errOutput != "" {
		log.Printf("Command stderr: %s", errOutput)
	}

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
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("No image found for container '%s' in daemonset '%s' (namespace: %s).",
						containerName, daemonsetName, namespace),
				},
			},
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Daemonset image:\n\n%s\n\nNamespace: %s\nDaemonset: %s\nContainer: %s\nCommand executed: kubectl %s",
					output, namespace, daemonsetName, containerName, strings.Join(cmdArgs, " ")),
			},
		},
	}, nil, nil
}

// =============================================================================
// BGP HANDLERS - KUBERNETES-ONLY, RUNS IN AGENT CONTAINER
// =============================================================================

// HandleGoBGPCommand is a generic handler for gobgp commands
func (s *VPPMCPServer) HandleGoBGPCommand(ctx context.Context, input BGPCommandInput, command, commandDescription string) (*mcp.CallToolResult, any, error) {
	// Log the request details
	log.Printf("Received %s request for pod: %s", commandDescription, input.PodName)
	log.Printf("Executing gobgp %s command on pod: %s", command, input.PodName)

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: Pod name is required. Please specify the Kubernetes pod name.",
				},
			},
		}, nil, fmt.Errorf("pod name is required")
	}

	// Initialize Kubernetes client for validation
	k8sClient, err := newKubeClient()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: Failed to create Kubernetes client: %v", err),
				},
			},
		}, nil, err
	}

	namespace := "calico-vpp-dataplane"

	// Validate pod exists
	_, err = k8sClient.CoreV1().Pods(namespace).Get(ctx, input.PodName, metav1.GetOptions{})
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error validating pod: %v", err),
				},
			},
		}, nil, err
	}

	// Execute the gobgp command on the Kubernetes pod
	result, err := ExecutePodGoBGPCommand(ctx, input.PodName, command)

	if err != nil {
		log.Printf("Error executing gobgp command: %v", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		node := result["node"].(string)
		pod := result["pod"].(string)

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: gobgp %s\nNode: %s\nPod: %s (container: agent)",
						commandDescription, output, cmd, node, pod),
				},
			},
		}

		log.Println("Successfully executed gobgp command, returning result")
		return response, nil, nil
	} else {
		errorMsg := result["error"].(string)
		cmd := result["command"].(string)
		node := result["node"].(string)
		pod, _ := result["pod"].(string)

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing gobgp command on node %s (pod: %s): %s\nCommand attempted: gobgp %s",
						node, pod, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing gobgp command on node %s (pod: %s): %s", node, pod, errorMsg)
		return errorResponse, nil, nil
	}
}

// HandleGoBGPParameterCommand is a consolidated handler for gobgp commands that require a parameter (IP, prefix, or neighbor)
func (s *VPPMCPServer) HandleGoBGPParameterCommand(ctx context.Context, input BGPParameterCommandInput, commandTemplate, commandDescription string) (*mcp.CallToolResult, any, error) {
	log.Printf("Received %s request for pod: %s, parameter: %s", commandDescription, input.PodName, input.Parameter)

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: Pod name is required. Please specify the Kubernetes pod name.",
				},
			},
		}, nil, fmt.Errorf("pod name is required")
	}

	if input.Parameter == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: Parameter is required. Please specify the IP address, prefix, or neighbor IP.",
				},
			},
		}, nil, fmt.Errorf("parameter is required")
	}

	namespace := "calico-vpp-dataplane"

	// Initialize Kubernetes client for validation
	k8sClient, err := newKubeClient()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: Failed to create Kubernetes client: %v", err),
				},
			},
		}, nil, err
	}

	// Validate pod exists
	_, err = k8sClient.CoreV1().Pods(namespace).Get(ctx, input.PodName, metav1.GetOptions{})
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error validating pod: %v", err),
				},
			},
		}, nil, err
	}

	// Build the command with parameter
	command := fmt.Sprintf(commandTemplate, input.Parameter)
	log.Printf("Executing gobgp %s command on pod: %s", command, input.PodName)

	// Execute the gobgp command on the Kubernetes pod
	result, err := ExecutePodGoBGPCommand(ctx, input.PodName, command)

	if err != nil {
		log.Printf("Error executing gobgp command: %v", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		node := result["node"].(string)
		pod := result["pod"].(string)

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: gobgp %s\nNode: %s\nPod: %s (container: agent)",
						commandDescription, output, cmd, node, pod),
				},
			},
		}

		log.Println("Successfully executed gobgp command, returning result")
		return response, nil, nil
	} else {
		errorMsg := result["error"].(string)
		cmd := result["command"].(string)
		node := result["node"].(string)
		pod, _ := result["pod"].(string)

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing gobgp command on node %s (pod: %s): %s\nCommand attempted: gobgp %s",
						node, pod, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing gobgp command on node %s (pod: %s): %s", node, pod, errorMsg)
		return errorResponse, nil, nil
	}
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
		Version: "1.0.0",
	}

	vppServer.server = mcp.NewServer(impl, nil)

	// Common mode description for tool documentation
	modeDesc := "\n\nExecution Modes:\n" +
		"- Kubernetes (default): Set pod_name to run on a VPP pod in Kubernetes\n" +
		"- Standalone: Set mode='standalone' to run on local VPP daemon (sock_path optional, default: /var/run/vpp/cli.sock)"

	// =========================================================================
	// SECTION 1: VPP CORE TOOLS - General VPP commands
	// =========================================================================

	// Define the vpp_show_version tool
	tool := &mcp.Tool{
		Name: "vpp_show_version",
		Description: "Get VPP version information by running 'vppctl show version'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, tool, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show version", "VPP Version Information")
	})

	// Define vpp_show_int tool
	toolShowInt := &mcp.Tool{
		Name: "vpp_show_int",
		Description: "Get VPP interface information by running 'vppctl show int'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowInt, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show int", "VPP Interface Information")
	})

	// Define vpp_show_int_addr tool
	toolShowIntAddr := &mcp.Tool{
		Name: "vpp_show_int_addr",
		Description: "Get VPP interface address information by running 'vppctl show int addr'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIntAddr, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show int addr", "VPP Interface Address Information")
	})

	// Define vpp_show_hardware_interfaces tool
	toolShowHardwareInterfaces := &mcp.Tool{
		Name: "vpp_show_hardware_interfaces",
		Description: "Shows detailed hardware interface information with VIRTIO queue statistics for all interfaces by running 'vppctl show hardware-interfaces' in a Kubernetes VPP container\n\n" +
			"This command provides comprehensive hardware-level details including:\n" +
			"- Hardware interface names and indices\n" +
			"- Link state and speed\n" +
			"- MAC addresses\n" +
			"- Driver information\n" +
			"- VIRTIO queue statistics and depths\n" +
			"- Hardware offload capabilities\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowHardwareInterfaces, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show hardware-interfaces", "VPP Hardware Interface Information")
	})

	// Define vpp_show_hardware_interface tool
	toolShowHardwareInterface := &mcp.Tool{
		Name: "vpp_show_hardware_interface",
		Description: "Show hardware information for a specific interface by running 'vppctl show hardware-interfaces <interface>' in a Kubernetes VPP container\n\n" +
			"This command returns hardware-level details (link state, speed, MAC, driver, queue stats, offload capabilities) for the selected interface only.\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n" +
			"- interface_name: The interface to inspect (for example, tun1)",
	}
	mcp.AddTool(vppServer.server, toolShowHardwareInterface, func(ctx context.Context, req *mcp.CallToolRequest, input VPPInterfaceInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleHardwareInterface(ctx, input)
	})

	// Define vpp_show_errors tool
	toolShowErrors := &mcp.Tool{
		Name: "vpp_show_errors",
		Description: "Get VPP error counters by running 'vppctl show errors'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowErrors, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show errors", "VPP Error Counters")
	})

	// Define vpp_show_session_verbose tool
	toolShowSession := &mcp.Tool{
		Name: "vpp_show_session_verbose",
		Description: "Get VPP session information by running 'vppctl show session verbose 2'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowSession, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show session verbose 2", "VPP Session Information (Verbose)")
	})

	// Define vpp_show_npol_rules tool
	toolShowNpolRules := &mcp.Tool{
		Name: "vpp_show_npol_rules",
		Description: "List rules that are referenced by policies by running 'vppctl show npol rules'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowNpolRules, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol rules", "VPP NPOL Rules")
	})

	// Define vpp_show_npol_policies tool
	toolShowNpolPolicies := &mcp.Tool{
		Name: "vpp_show_npol_policies",
		Description: "List all the policies that are referenced on interfaces by running 'vppctl show npol policies'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowNpolPolicies, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol policies", "VPP NPOL Policies")
	})

	// Define vpp_show_npol_ipset tool
	toolShowNpolIpset := &mcp.Tool{
		Name: "vpp_show_npol_ipset",
		Description: "List ipsets that are referenced by rules (IPsets are just list of IPs) by running 'vppctl show npol ipset'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowNpolIpset, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol ipset", "VPP NPOL IPset")
	})

	// Define vpp_show_npol_interfaces tool
	toolShowNpolInterfaces := &mcp.Tool{
		Name: "vpp_show_npol_interfaces",
		Description: "Show the resulting policies configured for every interface in VPP by running 'vppctl show npol interfaces'.\n\n" +
			"The first IPv4 address of every pod is provided to help identify which pod and interface belongs to.\n\n" +
			"Output interpretation:\n" +
			"- tx: contains rules that are applied on packets that LEAVE VPP on a given interface. Rules are applied top to bottom.\n" +
			"- rx: contains rules that are applied on packets that ENTER VPP on a given interface. Rules are applied top to bottom.\n" +
			"- profiles: are specific rules that are enforced when a matched rule action is PASS or when no policies are configured." + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowNpolInterfaces, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol interfaces", "VPP NPOL Interfaces")
	})

	// Define vpp_clear_errors tool
	toolClearErrors := &mcp.Tool{
		Name: "vpp_clear_errors",
		Description: "Reset the error counters by running 'vppctl clear errors'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolClearErrors, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "clear errors", "VPP Clear Error Counters")
	})

	// Define vpp_tcp_stats tool
	toolTcpStats := &mcp.Tool{
		Name: "vpp_tcp_stats",
		Description: "Display global statistics reported by TCP by running 'vppctl show tcp stats'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolTcpStats, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show tcp stats", "VPP TCP Statistics")
	})

	// Define vpp_session_stats tool
	toolSessionStats := &mcp.Tool{
		Name: "vpp_session_stats",
		Description: "Display global statistics reported by the session layer by running 'vppctl show session stats'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolSessionStats, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show session stats", "VPP Session Statistics")
	})

	// Define vpp_get_logs tool
	toolGetLogs := &mcp.Tool{
		Name: "vpp_get_logs",
		Description: "Display VPP logs by running 'vppctl show logging'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolGetLogs, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show logging", "VPP Logs")
	})

	// Define vpp_show_cnat_translation tool
	toolShowCnatTranslation := &mcp.Tool{
		Name: "vpp_show_cnat_translation",
		Description: "Shows the active CNAT translations by running 'vppctl show cnat translation'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowCnatTranslation, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show cnat translation", "VPP CNAT Translation")
	})

	// Define vpp_show_cnat_session tool
	toolShowCnatSession := &mcp.Tool{
		Name: "vpp_show_cnat_session",
		Description: "Lists the active CNAT sessions from the established five tuple to the five tuple rewrites by running 'vppctl show cnat session'.\n\n" +
			"Output interpretation:\n" +
			"The output shows the `incoming 5-tuple` first that is used to match packets along with the `protocol`. " +
			"Then it displays the `5-tuple after dNAT & sNAT`, followed by the `direction` and finally the `age` in seconds. " +
			"`direction` being input for the PRE-ROUTING sessions and output is the POST-ROUTING sessions" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowCnatSession, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show cnat session", "VPP CNAT Session")
	})

	// Define vpp_clear_run tool
	toolClearRun := &mcp.Tool{
		Name: "vpp_clear_run",
		Description: "Clears live running error stats in VPP by running 'vppctl clear run'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolClearRun, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "clear run", "VPP Clear Runtime Statistics")
	})

	// Define vpp_show_run tool
	toolShowRun := &mcp.Tool{
		Name: "vpp_show_run",
		Description: "Shows live running error stats in VPP by running 'vppctl show run'.\n\n" +
			"Debugging workflow:\n" +
			"Sometimes to debug an issue, you might need to run `vpp_clear_run` to erase historic stats and then wait for a few seconds in the issue state / run some tests " +
			"so that the error stats are repopulated and then run `vpp_show_run` in order to diagnose what is going on in the system\n\n" +
			"Output interpretation:\n" +
			"A loaded VPP will typically have (1) a high Vectors/Call maxing out at 256 (2) a low loops/sec struggling around 10000. " +
			"The Clocks column tells you the consumption in cycles per node on average. Beyond 1e3 is expensive." + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowRun, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show run", "VPP Runtime Statistics")
	})

	// Define vpp_show_ipip_tunnel tool
	toolShowIpipTunnel := &mcp.Tool{
		Name: "vpp_show_ipip_tunnel",
		Description: "Display IPIP tunnel status by running 'vppctl show ipip tunnel'" + modeDesc +
			"\n\nThis command shows IPIP tunnel configuration and status including:\n" +
			"- Tunnel instance number\n" +
			"- Source and destination IP addresses\n" +
			"- Table ID and software interface index\n" +
			"- Tunnel flags and DSCP settings\n\n" +
			"Parameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIpipTunnel, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show ipip tunnel", "VPP IPIP Tunnel Status")
	})

	// Define vpp_show_vxlan_tunnel tool
	toolShowVxlanTunnel := &mcp.Tool{
		Name: "vpp_show_vxlan_tunnel",
		Description: "Display VXLAN tunnel status by running 'vppctl show vxlan tunnel'" + modeDesc +
			"\n\nThis command shows VXLAN tunnel configuration and status including:\n" +
			"- Tunnel instance number\n" +
			"- Source and destination IPv6 addresses\n" +
			"- Source/destination ports and VNI\n" +
			"- FIB index, software interface index\n" +
			"- Encap/decap indices\n\n" +
			"Parameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowVxlanTunnel, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show vxlan tunnel", "VPP VXLAN Tunnel Status")
	})

	// Define vpp_show_tun_all tool
	toolShowTunAll := &mcp.Tool{
		Name: "vpp_show_tun_all",
		Description: "Display all tunnel interfaces in VPP by running 'vppctl show tun' in a Kubernetes VPP container\n\n" +
			"This command shows detailed information about all tunnel interfaces configured in VPP, including:\n" +
			"- Tunnel interface names and indices\n" +
			"- Tunnel types (GRE, VXLAN, IPSec, etc.)\n" +
			"- Source and destination addresses\n" +
			"- Tunnel state and configuration\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowTunAll, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show tun", "VPP Tunnel Interfaces")
	})

	toolShowTunInterface := &mcp.Tool{
		Name: "vpp_show_tun_interface",
		Description: "Inspect a specific tunnel interface in VPP by running 'vppctl show tun' and filtering to the requested interface (emulating 'grep <name> -A 40').\n\n" +
			"This command provides detailed statistics for the selected tunnel interface, including queue depths, buffer usage, and offload capabilities.\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n" +
			"- interface_name: The tunnel interface to inspect (for example, tun1)",
	}
	mcp.AddTool(vppServer.server, toolShowTunInterface, func(ctx context.Context, req *mcp.CallToolRequest, input VPPTunnelInterfaceInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleTunnelInterface(ctx, input)
	})

	// Define vpp_show_ip_table tool
	toolShowIpTable := &mcp.Tool{
		Name: "vpp_show_ip_table",
		Description: "Prints all available IPv4 VRFs by running 'vppctl show ip table'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIpTable, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show ip table", "VPP IPv4 VRF Tables")
	})

	// Define vpp_show_ip6_table tool
	toolShowIp6Table := &mcp.Tool{
		Name: "vpp_show_ip6_table",
		Description: "Prints all available IPv6 VRFs by running 'vppctl show ip6 table'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIp6Table, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show ip6 table", "VPP IPv6 VRF Tables")
	})

	// Define vpp_show_ip_fib tool
	toolShowIpFib := &mcp.Tool{
		Name: "vpp_show_ip_fib",
		Description: "Prints all routes in a given pod IPv4 VRF by running 'vppctl show ip fib index <idx>'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIpFib, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBCommand(ctx, input, "show ip fib index %s", "VPP IPv4 FIB Routes")
	})

	// Define vpp_show_ip6_fib tool
	toolShowIp6Fib := &mcp.Tool{
		Name: "vpp_show_ip6_fib",
		Description: "Prints all routes in a given pod IPv6 VRF by running 'vppctl show ip6 fib index <idx>'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIp6Fib, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBCommand(ctx, input, "show ip6 fib index %s", "VPP IPv6 FIB Routes")
	})

	// Define vpp_show_ip_fib_prefix tool
	toolShowIpFibPrefix := &mcp.Tool{
		Name: "vpp_show_ip_fib_prefix",
		Description: "Prints information about a specific prefix in a given pod IPv4 VRF by running 'vppctl show ip fib index <idx> <prefix>'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index\n" +
			"- prefix: The IP prefix to query (e.g., 10.0.0.0/24)\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIpFibPrefix, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBPrefixInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBPrefixCommand(ctx, input, "show ip fib index %s %s", "VPP IPv4 FIB Prefix Information")
	})

	// Define vpp_show_ip6_fib_prefix tool
	toolShowIp6FibPrefix := &mcp.Tool{
		Name: "vpp_show_ip6_fib_prefix",
		Description: "Prints information about a specific prefix in a given pod IPv6 VRF by running 'vppctl show ip6 fib index <idx> <prefix>'" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index\n" +
			"- prefix: The IPv6 prefix to query (e.g., 2001:db8::/32)\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path",
	}
	mcp.AddTool(vppServer.server, toolShowIp6FibPrefix, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBPrefixInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBPrefixCommand(ctx, input, "show ip6 fib index %s %s", "VPP IPv6 FIB Prefix Information")
	})

	// =========================================================================
	// SECTION 2: VPP CAPTURE TOOLS - Packet capture with locking
	// =========================================================================

	captureDesc := "\n\nExecution Modes:\n" +
		"- Kubernetes (default): Set pod_name to run on a VPP pod in Kubernetes\n" +
		"- Standalone: Set mode='standalone' to run on local VPP daemon\n\n" +
		"**Note**: Capture operations use locking to prevent parallel execution. " +
		"Use vpp_capture_cleanup to force cleanup if a previous capture failed."

	// Define vpp_trace tool
	toolTrace := &mcp.Tool{
		Name: "vpp_trace",
		Description: "Capture VPP packet traces by running 'vppctl trace add'" + captureDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) Pod name running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path\n" +
			"- count: (optional) Number of packets to capture (default: 500)\n" +
			"- timeout: (optional) Capture duration in seconds (default: 30)\n" +
			"- interface: (optional) Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl\n\n" +
			"The tool will: Clear traces → Start capture → Wait → Return traces → Cleanup",
	}
	mcp.AddTool(vppServer.server, toolTrace, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleTraceCapture(ctx, input)
	})

	// Define vpp_pcap tool
	toolPcap := &mcp.Tool{
		Name: "vpp_pcap",
		Description: "Capture VPP packets to pcap file by running 'vppctl pcap trace'" + captureDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) Pod name running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path\n" +
			"- count: (optional) Number of packets to capture (default: 500)\n" +
			"- timeout: (optional) Capture duration in seconds (default: 30)\n" +
			"- interface: (optional) Interface name or 'any' (default: any)\n\n" +
			"The tool will: Validate interface → Start pcap → Wait → Stop → Save to /tmp/trace.pcap",
	}
	mcp.AddTool(vppServer.server, toolPcap, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handlePcapCapture(ctx, input)
	})

	// Define vpp_dispatch tool
	toolDispatch := &mcp.Tool{
		Name: "vpp_dispatch",
		Description: "Capture VPP dispatch trace to pcap file by running 'vppctl pcap dispatch trace'" + captureDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) Pod name running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path\n" +
			"- count: (optional) Number of packets to capture (default: 500)\n" +
			"- timeout: (optional) Capture duration in seconds (default: 30)\n" +
			"- interface: (optional) Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl\n\n" +
			"The tool will: Start dispatch trace → Wait → Stop → Save to /tmp/dispatch.pcap",
	}
	mcp.AddTool(vppServer.server, toolDispatch, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleDispatchCapture(ctx, input)
	})

	// Define vpp_capture_cleanup tool
	toolCaptureCleanup := &mcp.Tool{
		Name: "vpp_capture_cleanup",
		Description: "Force cleanup of all VPP capture operations (trace, pcap, dispatch)" + modeDesc +
			"\n\nParameters:\n" +
			"- pod_name: (Kubernetes mode) Pod name running VPP\n" +
			"- mode: (optional) 'kubernetes' or 'standalone'\n" +
			"- sock_path: (Standalone mode, optional) VPP socket path\n\n" +
			"Use this tool to:\n" +
			"- Stop all active captures\n" +
			"- Remove capture lock files\n" +
			"- Clean up temporary capture files\n" +
			"- Restore system to clean state after a failed capture",
	}
	mcp.AddTool(vppServer.server, toolCaptureCleanup, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleCaptureCleanup(ctx, input)
	})

	// =========================================================================
	// SECTION 3: CALICOVPP-SPECIFIC TOOLS - CalicoVPP/Kubernetes specific
	// =========================================================================

	// Define vpp_get_pods tool
	toolGetPods := &mcp.Tool{
		Name: "vpp_get_pods",
		Description: "List all calico-vpp pods along with their IP addresses and the node on which they are running\n\n" +
			"**Note**: This tool is specific to CalicoVPP/Kubernetes environments.\n\n" +
			"This tool runs 'kubectl get pods -n calico-vpp-dataplane -owide' to display:\n" +
			"- Pod names\n" +
			"- Pod status\n" +
			"- Pod IP addresses\n" +
			"- Node names\n" +
			"- Age and other metadata\n\n" +
			"No parameters required.",
	}
	mcp.AddTool(vppServer.server, toolGetPods, func(ctx context.Context, req *mcp.CallToolRequest, input EmptyInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleGetPods(ctx, input)
	})

	// Define vpp_show_daemonset_image tool
	toolShowDaemonsetImage := &mcp.Tool{
		Name: "vpp_show_daemonset_image",
		Description: "Show the image configured for a daemonset container (defaults match CalicoVPP: namespace=calico-vpp-dataplane, daemonset=calico-vpp-node, container=vpp).\n\n" +
			"This tool runs the equivalent of:\n" +
			"kubectl get daemonset -n <namespace> <daemonset_name> -o jsonpath='{.spec.template.spec.containers[?(@.name==\"<container_name>\")].image}'\n\n" +
			"Optional parameters:\n" +
			"- namespace: Kubernetes namespace (default: calico-vpp-dataplane)\n" +
			"- daemonset_name: Daemonset name (default: calico-vpp-node)\n" +
			"- container_name: Container name in daemonset spec (default: vpp)",
	}
	mcp.AddTool(vppServer.server, toolShowDaemonsetImage, func(ctx context.Context, req *mcp.CallToolRequest, input VPPDaemonsetImageInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleShowDaemonsetImage(ctx, input)
	})

	// =========================================================================
	// SECTION 4: BGP TOOLS - GoBGP commands (Kubernetes-only, runs in agent container)
	// =========================================================================

	bgpNote := "\n\n**Note**: This tool is specific to CalicoVPP/Kubernetes environments. " +
		"It runs in the 'agent' container (not 'vpp') of the calico-vpp pod in the 'calico-vpp-dataplane' namespace."

	// Define bgp_show_neighbors tool
	toolBgpShowNeighbors := &mcp.Tool{
		Name: "bgp_show_neighbors",
		Description: "Show BGP peers by running 'gobgp neighbor' in the agent container of a calico-vpp pod" + bgpNote +
			"\n\nRequired parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n\n" +
			"Output interpretation:\n" +
			"- Established peerings will show up as Establ\n" +
			"- Unsuccessful connections will show up as Opened with 0 in #Received Accepted\n" +
			"- CalicoVPP learns about new peers using the kubernetes API. If peers are missing from this list, there might be an issue accessing this API",
	}
	mcp.AddTool(vppServer.server, toolBgpShowNeighbors, func(ctx context.Context, req *mcp.CallToolRequest, input BGPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.HandleGoBGPCommand(ctx, input, "neighbor", "BGP Neighbor Information")
	})

	// Define bgp_show_global_info tool
	toolBgpShowGlobalInfo := &mcp.Tool{
		Name: "bgp_show_global_info",
		Description: "Show BGP global information by running 'gobgp global' in the agent container of a calico-vpp pod" + bgpNote +
			"\n\nRequired parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n\n" +
			"Output interpretation:\n" +
			"- Shows the information goBGP advertises to peers",
	}
	mcp.AddTool(vppServer.server, toolBgpShowGlobalInfo, func(ctx context.Context, req *mcp.CallToolRequest, input BGPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.HandleGoBGPCommand(ctx, input, "global", "BGP Global Information")
	})

	// Define bgp_show_global_rib4 tool
	toolBgpShowGlobalRib4 := &mcp.Tool{
		Name: "bgp_show_global_rib4",
		Description: "Show BGP IPv4 RIB information by running 'gobgp global rib -a 4' in the agent container of a calico-vpp pod" + bgpNote +
			"\n\nRequired parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n\n" +
			"Output interpretation:\n" +
			"- Prints out the IPv4 prefixes advertised by peers\n" +
			"- Next Hop being the peer's IP\n" +
			"- Shows all route information",
	}
	mcp.AddTool(vppServer.server, toolBgpShowGlobalRib4, func(ctx context.Context, req *mcp.CallToolRequest, input BGPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.HandleGoBGPCommand(ctx, input, "global rib -a 4", "BGP IPv4 RIB Information")
	})

	// Define bgp_show_global_rib6 tool
	toolBgpShowGlobalRib6 := &mcp.Tool{
		Name: "bgp_show_global_rib6",
		Description: "Show BGP IPv6 RIB information by running 'gobgp global rib -a 6' in the agent container of a calico-vpp pod" + bgpNote +
			"\n\nRequired parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n\n" +
			"Output interpretation:\n" +
			"- Prints out the IPv6 prefixes advertised by peers\n" +
			"- Next Hop being the peer's IP\n" +
			"- Shows all route information",
	}
	mcp.AddTool(vppServer.server, toolBgpShowGlobalRib6, func(ctx context.Context, req *mcp.CallToolRequest, input BGPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.HandleGoBGPCommand(ctx, input, "global rib -a 6", "BGP IPv6 RIB Information")
	})

	// Define bgp_show_ip tool
	toolBgpShowIp := &mcp.Tool{
		Name: "bgp_show_ip",
		Description: "Show BGP RIB entry for a specific IP by running 'gobgp global rib <ip>' in the agent container of a calico-vpp pod" + bgpNote +
			"\n\nRequired parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n" +
			"- parameter: The IP address to query\n\n" +
			"Output interpretation:\n" +
			"- Prints the RIB entry for that specific IP\n" +
			"- Shows specific route information",
	}
	mcp.AddTool(vppServer.server, toolBgpShowIp, func(ctx context.Context, req *mcp.CallToolRequest, input BGPParameterCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.HandleGoBGPParameterCommand(ctx, input, "global rib %s", "BGP RIB Entry for IP")
	})

	// Define bgp_show_prefix tool
	toolBgpShowPrefix := &mcp.Tool{
		Name: "bgp_show_prefix",
		Description: "Show BGP RIB entry for a specific prefix by running 'gobgp global rib <prefix>' in the agent container of a calico-vpp pod" + bgpNote +
			"\n\nRequired parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n" +
			"- parameter: The prefix to query (e.g., 10.0.0.0/24)\n\n" +
			"Output interpretation:\n" +
			"- Prints the RIB entry for that specific prefix\n" +
			"- Shows specific route information",
	}
	mcp.AddTool(vppServer.server, toolBgpShowPrefix, func(ctx context.Context, req *mcp.CallToolRequest, input BGPParameterCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.HandleGoBGPParameterCommand(ctx, input, "global rib %s", "BGP RIB Entry for Prefix")
	})

	// Define bgp_show_neighbor tool
	toolBgpShowNeighbor := &mcp.Tool{
		Name: "bgp_show_neighbor",
		Description: "Show detailed information for a specific BGP neighbor by running 'gobgp neighbor <neighborIP>' in the agent container of a calico-vpp pod" + bgpNote +
			"\n\nRequired parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n" +
			"- parameter: The IP address of the BGP neighbor\n\n" +
			"Output interpretation:\n" +
			"- Prints detailed status information for the specified BGP peer",
	}
	mcp.AddTool(vppServer.server, toolBgpShowNeighbor, func(ctx context.Context, req *mcp.CallToolRequest, input BGPParameterCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.HandleGoBGPParameterCommand(ctx, input, "neighbor %s", "BGP Neighbor Details")
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
