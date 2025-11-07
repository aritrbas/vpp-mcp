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
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// ExecutePodVPPCommand runs a VPP command directly on a specified Kubernetes pod
func ExecutePodVPPCommand(ctx context.Context, podName, command string) (map[string]interface{}, error) {
	// Use hardcoded defaults
	namespace := "calico-vpp-dataplane"
	containerName := "vpp"

	// Build kubectl exec command
	cmdArgs := []string{
		"exec",
		"-n", namespace,
		podName,
		"-c", containerName,
	}

	// Add the vppctl command
	cmdArgs = append(cmdArgs, "--", "vppctl")

	// Add the specific VPP command arguments
	cmdArgs = append(cmdArgs, strings.Fields(command)...)

	// Execute the command with a timeout
	log.Printf("Executing command: kubectl %s", strings.Join(cmdArgs, " "))

	// Set a timeout for the command
	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

	err := execErr

	if err != nil {
		errorMsg := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			errorMsg = string(exitErr.Stderr)
		}
		return map[string]interface{}{
			"success":   false,
			"error":     fmt.Sprintf("%v - %s", err, errorMsg),
			"pod":       podName,
			"namespace": namespace,
			"command":   command,
		}, err
	}
	return map[string]interface{}{
		"success":   true,
		"output":    string(output),
		"command":   command,
		"pod":       podName,
		"namespace": namespace,
		"container": containerName,
	}, nil
}

const kubeClientTimeout = 30 * time.Second

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

// VPPCommandInput represents the generic input for VPP command tools
type VPPCommandInput struct {
	// PodName specifies the name of the Kubernetes pod running VPP
	PodName string `json:"pod_name"`
}

// VPPCaptureInput represents the input for VPP packet capture tools (trace, pcap, dispatch)
type VPPCaptureInput struct {
	// PodName specifies the name of the Kubernetes pod running VPP
	PodName string `json:"pod_name"`
	// Count specifies the number of packets to capture (default: run for 30 seconds)
	Count int `json:"count,omitempty"`
	// Interface specifies the interface type or name to capture from
	Interface string `json:"interface,omitempty"`
}

// VPPFIBInput represents the input for VPP FIB tools requiring fib_index
type VPPFIBInput struct {
	// PodName specifies the name of the Kubernetes pod running VPP
	PodName string `json:"pod_name"`
	// FibIndex specifies the FIB table index
	FibIndex string `json:"fib_index"`
}

// VPPFIBPrefixInput represents the input for VPP FIB tools requiring fib_index and prefix
type VPPFIBPrefixInput struct {
	// PodName specifies the name of the Kubernetes pod running VPP
	PodName string `json:"pod_name"`
	// FibIndex specifies the FIB table index
	FibIndex string `json:"fib_index"`
	// Prefix specifies the IP prefix to query
	Prefix string `json:"prefix"`
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

// VPPMCPServer implements the MCP server for VPP debugging
type VPPMCPServer struct {
	server *mcp.Server
}

// NewVPPMCPServer creates a new VPP MCP server
func NewVPPMCPServer() *VPPMCPServer {
	return &VPPMCPServer{}
}

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

// handleVPPCommand is a generic handler for VPP commands
func (s *VPPMCPServer) handleVPPCommand(ctx context.Context, input VPPCommandInput, command, commandDescription string) (*mcp.CallToolResult, any, error) {
	// Log the request details
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received %s request with input: %s", commandDescription, string(inputJSON))
	log.Printf("Executing vppctl %s command on pod: %s", command, input.PodName)

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: PodName is required. Please specify the Kubernetes pod name running VPP.",
				},
			},
		}, nil, fmt.Errorf("PodName is required")
	}

	// Execute the VPP command on the Kubernetes pod
	log.Printf("About to execute pod VPP command...")
	result, err := ExecutePodVPPCommand(ctx, input.PodName, command)

	log.Printf("Command execution completed, processing results...")
	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		pod := result["pod"].(string)

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: vppctl %s\nPod: %s (container: vpp)",
						commandDescription, output, cmd, pod),
				},
			},
		}

		log.Println("Successfully executed VPP command, returning result")
		return response, nil, nil
	} else {
		errorMsg := result["error"].(string)
		cmd := result["command"].(string)
		pod := result["pod"].(string)

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing VPP command on pod %s: %s\nCommand attempted: vppctl %s",
						pod, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing VPP command on pod %s: %s", pod, errorMsg)
		return errorResponse, nil, nil
	}
}

// handleVPPFIBCommand is a handler for VPP FIB commands that require fib_index
func (s *VPPMCPServer) handleVPPFIBCommand(ctx context.Context, input VPPFIBInput, commandTemplate, commandDescription string) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received %s request with input: %s", commandDescription, string(inputJSON))

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: PodName is required. Please specify the Kubernetes pod name running VPP.",
				},
			},
		}, nil, fmt.Errorf("PodName is required")
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
	log.Printf("Executing vppctl %s command on pod: %s", command, input.PodName)

	result, err := ExecutePodVPPCommand(ctx, input.PodName, command)

	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		pod := result["pod"].(string)

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: vppctl %s\nPod: %s (container: vpp)",
						commandDescription, output, cmd, pod),
				},
			},
		}

		log.Println("Successfully executed VPP FIB command, returning result")
		return response, nil, nil
	} else {
		errorMsg := result["error"].(string)
		cmd := result["command"].(string)
		pod := result["pod"].(string)

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing VPP command on pod %s: %s\nCommand attempted: vppctl %s",
						pod, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing VPP FIB command on pod %s: %s", pod, errorMsg)
		return errorResponse, nil, nil
	}
}

// handleVPPFIBPrefixCommand is a handler for VPP FIB commands that require fib_index and prefix
func (s *VPPMCPServer) handleVPPFIBPrefixCommand(ctx context.Context, input VPPFIBPrefixInput, commandTemplate, commandDescription string) (*mcp.CallToolResult, any, error) {
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received %s request with input: %s", commandDescription, string(inputJSON))

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: PodName is required. Please specify the Kubernetes pod name running VPP.",
				},
			},
		}, nil, fmt.Errorf("PodName is required")
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
	log.Printf("Executing vppctl %s command on pod: %s", command, input.PodName)

	result, err := ExecutePodVPPCommand(ctx, input.PodName, command)

	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		pod := result["pod"].(string)

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: vppctl %s\nPod: %s (container: vpp)",
						commandDescription, output, cmd, pod),
				},
			},
		}

		log.Println("Successfully executed VPP FIB prefix command, returning result")
		return response, nil, nil
	} else {
		errorMsg := result["error"].(string)
		cmd := result["command"].(string)
		pod := result["pod"].(string)

		errorResponse := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error executing VPP command on pod %s: %s\nCommand attempted: vppctl %s",
						pod, errorMsg, cmd),
				},
			},
		}
		log.Printf("Error executing VPP FIB prefix command on pod %s: %s", pod, errorMsg)
		return errorResponse, nil, nil
	}
}

// handleTraceCapture implements VPP trace capture
func (s *VPPMCPServer) handleTraceCapture(ctx context.Context, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
	log.Printf("Received trace capture request for pod: %s", input.PodName)

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: PodName is required. Please specify the Kubernetes pod name running VPP.",
				},
			},
		}, nil, fmt.Errorf("PodName is required")
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

	// Map interface type to VPP input node
	vppInputNode, _, err := mapInterfaceTypeToVppInputNode(k8sClient, input.Interface)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error mapping interface: %v", err),
				},
			},
		}, nil, err
	}

	// Determine count (default 500 if not specified)
	count := input.Count
	if count == 0 {
		count = 500
	}

	// Step 1: Clear trace to ensure clean state
	log.Printf("Clearing trace on pod %s", input.PodName)
	_, err = ExecutePodVPPCommand(ctx, input.PodName, "clear trace")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error clearing trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 2: Start trace capture
	traceCmd := fmt.Sprintf("trace add %s %d", vppInputNode, count)
	log.Printf("Starting trace: %s", traceCmd)
	_, err = ExecutePodVPPCommand(ctx, input.PodName, traceCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture (30 seconds or until count is reached)
	log.Printf("Capturing packets for 30 seconds or until %d packets captured...", count)
	time.Sleep(30 * time.Second)

	// Step 4: Get trace results
	traceCmd = fmt.Sprintf("show trace max %d", count)
	log.Printf("Retrieving trace results...")
	result, err := ExecutePodVPPCommand(ctx, input.PodName, traceCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error retrieving trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 5: Clear trace after retrieval
	_, _ = ExecutePodVPPCommand(ctx, input.PodName, "clear trace")

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP Trace Capture Results:\n\n%s\n\nCapture Parameters:\n- VPP Input Node: %s\n- Count: %d\n- Capture Duration: 30 seconds\n- Pod: %s\n\n**Important**: Trace is not saved to any file\n\n",
						output, vppInputNode, count, input.PodName),
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

// handlePcapCapture implements VPP pcap capture
func (s *VPPMCPServer) handlePcapCapture(ctx context.Context, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
	log.Printf("Received pcap capture request for pod: %s", input.PodName)

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: PodName is required. Please specify the Kubernetes pod name running VPP.",
				},
			},
		}, nil, fmt.Errorf("PodName is required")
	}

	// Get list of available interfaces
	interfaceResult, err := ExecutePodVPPCommand(ctx, input.PodName, "show int")
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
	if len(availableInterfaces) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: No up interfaces found in VPP",
				},
			},
		}, nil, fmt.Errorf("no up interfaces found")
	}

	// Validate interface if provided
	interfaceName := input.Interface
	if interfaceName == "" {
		// Default to 'any' interface
		interfaceName = "any"
	} else if interfaceName != "any" {
		// Validate provided interface (skip validation for 'any' since it's special)
		found := false
		for _, iface := range availableInterfaces {
			if iface == interfaceName {
				found = true
				break
			}
		}
		if !found {
			var ifaceList strings.Builder
			ifaceList.WriteString("\nAvailable interfaces:")
			for i, iface := range availableInterfaces {
				ifaceList.WriteString(fmt.Sprintf("\n%d. %s", i+1, iface))
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error: Interface '%s' not found.%s", interfaceName, ifaceList.String()),
					},
				},
			}, nil, fmt.Errorf("interface not found")
		}
	}

	// Determine count (default 500 if not specified)
	count := input.Count
	if count == 0 {
		count = 500
	}

	// Step 1: Stop any existing pcap capture
	log.Printf("Stopping any existing pcap capture on pod %s", input.PodName)
	_, _ = ExecutePodVPPCommand(ctx, input.PodName, "pcap trace off")

	// Step 2: Start pcap capture
	pcapCmd := fmt.Sprintf("pcap trace tx rx max %d intfc %s file trace.pcap", count, interfaceName)
	log.Printf("Starting pcap: %s", pcapCmd)
	_, err = ExecutePodVPPCommand(ctx, input.PodName, pcapCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting pcap: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture (30 seconds or until count is reached)
	log.Printf("Capturing packets for 30 seconds or until %d packets captured...", count)
	time.Sleep(30 * time.Second)

	// Step 4: Stop pcap capture
	log.Printf("Stopping pcap capture...")
	result, err := ExecutePodVPPCommand(ctx, input.PodName, "pcap trace off")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error stopping pcap: %v", err),
				},
			},
		}, nil, err
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP PCAP Capture Results:\n\n%s\n\nCapture Parameters:\n- Interface: %s\n- Count: %d\n- Capture Duration: 30 seconds\n- Pod: %s\n\n**Important**: PCAP file saved at /tmp/trace.pcap\n\n",
						output, interfaceName, count, input.PodName),
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

// handleDispatchCapture implements VPP dispatch trace capture
func (s *VPPMCPServer) handleDispatchCapture(ctx context.Context, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
	log.Printf("Received dispatch capture request for pod: %s", input.PodName)

	if input.PodName == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: PodName is required. Please specify the Kubernetes pod name running VPP.",
				},
			},
		}, nil, fmt.Errorf("PodName is required")
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

	// Map interface type to VPP input node
	vppInputNode, _, err := mapInterfaceTypeToVppInputNode(k8sClient, input.Interface)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error mapping interface: %v", err),
				},
			},
		}, nil, err
	}

	// Determine count (default 500 if not specified)
	count := input.Count
	if count == 0 {
		count = 500
	}

	// Step 1: Stop any existing dispatch trace
	log.Printf("Stopping any existing dispatch trace on pod %s", input.PodName)
	_, _ = ExecutePodVPPCommand(ctx, input.PodName, "pcap dispatch trace off")

	// Step 2: Start dispatch trace capture
	dispatchCmd := fmt.Sprintf("pcap dispatch trace on max %d buffer-trace %s %d", count, vppInputNode, count)
	log.Printf("Starting dispatch trace: %s", dispatchCmd)
	_, err = ExecutePodVPPCommand(ctx, input.PodName, dispatchCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting dispatch trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture (30 seconds or until count is reached)
	log.Printf("Capturing packets for 30 seconds or until %d packets captured...", count)
	time.Sleep(30 * time.Second)

	// Step 4: Stop dispatch trace
	log.Printf("Stopping dispatch trace...")
	result, err := ExecutePodVPPCommand(ctx, input.PodName, "pcap dispatch trace off")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error stopping dispatch trace: %v", err),
				},
			},
		}, nil, err
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP Dispatch Trace Results:\n\n%s\n\nCapture Parameters:\n- VPP Input Node: %s\n- Count: %d\n- Capture Duration: 30 seconds\n- Pod: %s\n\n**Important**: Dispatch PCAP file saved at /tmp/dispatch.pcap\n\n",
						output, vppInputNode, count, input.PodName),
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

	// Define the vpp_show_version tool with a better description
	tool := &mcp.Tool{
		Name: "vpp_show_version",
		Description: "Get VPP version information by running 'vppctl show version' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}

	// Add the tool to the server
	mcp.AddTool(vppServer.server, tool, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show version", "VPP Version Information")
	})

	// Define vpp_show_int tool
	toolShowInt := &mcp.Tool{
		Name: "vpp_show_int",
		Description: "Get VPP interface information by running 'vppctl show int' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowInt, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show int", "VPP Interface Information")
	})

	// Define vpp_show_int_addr tool
	toolShowIntAddr := &mcp.Tool{
		Name: "vpp_show_int_addr",
		Description: "Get VPP interface address information by running 'vppctl show int addr' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowIntAddr, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show int addr", "VPP Interface Address Information")
	})

	// Define vpp_show_errors tool
	toolShowErrors := &mcp.Tool{
		Name: "vpp_show_errors",
		Description: "Get VPP error counters by running 'vppctl show errors' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowErrors, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show errors", "VPP Error Counters")
	})

	// Define vpp_show_session_verbose tool
	toolShowSession := &mcp.Tool{
		Name: "vpp_show_session_verbose",
		Description: "Get VPP session information by running 'vppctl show session verbose 2' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowSession, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show session verbose 2", "VPP Session Information (Verbose)")
	})

	// Define vpp_show_npol_rules tool
	toolShowNpolRules := &mcp.Tool{
		Name: "vpp_show_npol_rules",
		Description: "List rules that are referenced by policies by running 'vppctl show npol rules' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowNpolRules, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol rules", "VPP NPOL Rules")
	})

	// Define vpp_show_npol_policies tool
	toolShowNpolPolicies := &mcp.Tool{
		Name: "vpp_show_npol_policies",
		Description: "List all the policies that are referenced on interfaces by running 'vppctl show npol policies' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowNpolPolicies, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol policies", "VPP NPOL Policies")
	})

	// Define vpp_show_npol_ipset tool
	toolShowNpolIpset := &mcp.Tool{
		Name: "vpp_show_npol_ipset",
		Description: "List ipsets that are referenced by rules (IPsets are just list of IPs) by running 'vppctl show npol ipset' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowNpolIpset, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol ipset", "VPP NPOL IPset")
	})

	// Define vpp_show_npol_interfaces tool
	toolShowNpolInterfaces := &mcp.Tool{
		Name: "vpp_show_npol_interfaces",
		Description: "Show the resulting policies configured for every interface in VPP by running 'vppctl show npol interfaces' in a Kubernetes VPP container.\n\n" +
			"The first IPv4 address of every pod is provided to help identify which pod and interface belongs to.\n\n" +
			"Output interpretation:\n" +
			"- tx: contains rules that are applied on packets that LEAVE VPP on a given interface. Rules are applied top to bottom.\n" +
			"- rx: contains rules that are applied on packets that ENTER VPP on a given interface. Rules are applied top to bottom.\n" +
			"- profiles: are specific rules that are enforced when a matched rule action is PASS or when no policies are configured.\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowNpolInterfaces, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show npol interfaces", "VPP NPOL Interfaces")
	})

	// Define vpp_trace tool
	toolTrace := &mcp.Tool{
		Name: "vpp_trace",
		Description: "Capture VPP packet traces by running 'vppctl trace add' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n\n" +
			"Optional parameters:\n" +
			"- count: Number of packets to capture (default: 500)\n" +
			"- interface: Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)\n\n" +
			"The tool will:\n" +
			"1. Clear existing traces\n" +
			"2. Start packet capture\n" +
			"3. Wait 30 seconds or until count is reached\n" +
			"4. Display captured traces",
	}
	mcp.AddTool(vppServer.server, toolTrace, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleTraceCapture(ctx, input)
	})

	// Define vpp_pcap tool
	toolPcap := &mcp.Tool{
		Name: "vpp_pcap",
		Description: "Capture VPP packets to pcap file by running 'vppctl pcap trace' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n\n" +
			"Optional parameters:\n" +
			"- count: Number of packets to capture (default: 500)\n" +
			"- interface: Interface name (e.g., host-eth0) or 'any' (default: first available interface)\n\n" +
			"The tool will:\n" +
			"1. Validate the interface exists\n" +
			"2. Start pcap capture on tx/rx\n" +
			"3. Wait 30 seconds or until count is reached\n" +
			"4. Stop capture and save to /tmp/vpp-capture-<timestamp>.pcap\n" +
			"5. Display capture status",
	}
	mcp.AddTool(vppServer.server, toolPcap, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handlePcapCapture(ctx, input)
	})

	// Define vpp_dispatch tool
	toolDispatch := &mcp.Tool{
		Name: "vpp_dispatch",
		Description: "Capture VPP dispatch trace to pcap file by running 'vppctl pcap dispatch trace' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n\n" +
			"Optional parameters:\n" +
			"- count: Number of packets to capture (default: 500)\n" +
			"- interface: Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)\n\n" +
			"The tool will:\n" +
			"1. Start dispatch trace with buffer trace\n" +
			"2. Wait 30 seconds or until count is reached\n" +
			"3. Stop capture and save to /tmp/vpp-dispatch-<timestamp>.pcap\n" +
			"4. Display capture status",
	}
	mcp.AddTool(vppServer.server, toolDispatch, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCaptureInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleDispatchCapture(ctx, input)
	})

	// Define vpp_get_pods tool
	toolGetPods := &mcp.Tool{
		Name: "vpp_get_pods",
		Description: "List all calico-vpp pods along with their IP addresses and the node on which they are running\n\n" +
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

	// Define vpp_clear_errors tool
	toolClearErrors := &mcp.Tool{
		Name: "vpp_clear_errors",
		Description: "Reset the error counters by running 'vppctl clear errors' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolClearErrors, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "clear errors", "VPP Clear Error Counters")
	})

	// Define vpp_tcp_stats tool
	toolTcpStats := &mcp.Tool{
		Name: "vpp_tcp_stats",
		Description: "Display global statistics reported by TCP by running 'vppctl show tcp stats' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolTcpStats, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show tcp stats", "VPP TCP Statistics")
	})

	// Define vpp_session_stats tool
	toolSessionStats := &mcp.Tool{
		Name: "vpp_session_stats",
		Description: "Display global statistics reported by the session layer by running 'vppctl show session stats' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolSessionStats, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show session stats", "VPP Session Statistics")
	})

	// Define vpp_get_logs tool
	toolGetLogs := &mcp.Tool{
		Name: "vpp_get_logs",
		Description: "Display VPP logs by running 'vppctl show logging' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolGetLogs, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show logging", "VPP Logs")
	})

	// Define vpp_show_cnat_translation tool
	toolShowCnatTranslation := &mcp.Tool{
		Name: "vpp_show_cnat_translation",
		Description: "Shows the active CNAT translations by running 'vppctl show cnat translation' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowCnatTranslation, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show cnat translation", "VPP CNAT Translation")
	})

	// Define vpp_show_cnat_session tool
	toolShowCnatSession := &mcp.Tool{
		Name: "vpp_show_cnat_session",
		Description: "Lists the active CNAT sessions from the established five tuple to the five tuple rewrites by running 'vppctl show cnat session' in a Kubernetes VPP container\n\n" +
			"Output interpretation:\n" +
			"The output shows the `incoming 5-tuple` first that is used to match packets along with the `protocol`. " +
			"Then it displays the `5-tuple after dNAT & sNAT`, followed by the `direction` and finally the `age` in seconds. " +
			"`direction` being input for the PRE-ROUTING sessions and output is the POST-ROUTING sessions\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowCnatSession, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show cnat session", "VPP CNAT Session")
	})

	// Define vpp_clear_run tool
	toolClearRun := &mcp.Tool{
		Name: "vpp_clear_run",
		Description: "Clears live running error stats in VPP by running 'vppctl clear run' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolClearRun, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "clear run", "VPP Clear Runtime Statistics")
	})

	// Define vpp_show_run tool
	toolShowRun := &mcp.Tool{
		Name: "vpp_show_run",
		Description: "Shows live running error stats in VPP by running 'vppctl show run' in a Kubernetes VPP container\n\n" +
			"Debugging workflow:\n" +
			"Sometimes to debug an issue, you might need to run `vpp_clear_run` to erase historic stats and then wait for a few seconds in the issue state / run some tests " +
			"so that the error stats are repopulated and then run `vpp_show_run` in order to diagnose what is going on in the system\n\n" +
			"Output interpretation:\n" +
			"A loaded VPP will typically have (1) a high Vectors/Call maxing out at 256 (2) a low loops/sec struggling around 10000. " +
			"The Clocks column tells you the consumption in cycles per node on average. Beyond 1e3 is expensive.\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowRun, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show run", "VPP Runtime Statistics")
	})

	// Define vpp_show_ip_table tool
	toolShowIpTable := &mcp.Tool{
		Name: "vpp_show_ip_table",
		Description: "Prints all available IPv4 VRFs by running 'vppctl show ip table' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowIpTable, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show ip table", "VPP IPv4 VRF Tables")
	})

	// Define vpp_show_ip6_table tool
	toolShowIp6Table := &mcp.Tool{
		Name: "vpp_show_ip6_table",
		Description: "Prints all available IPv6 VRFs by running 'vppctl show ip6 table' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowIp6Table, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "show ip6 table", "VPP IPv6 VRF Tables")
	})

	// Define vpp_show_ip_fib tool
	toolShowIpFib := &mcp.Tool{
		Name: "vpp_show_ip_fib",
		Description: "Prints all routes in a given pod IPv4 VRF by running 'vppctl show ip fib index <idx>' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index",
	}
	mcp.AddTool(vppServer.server, toolShowIpFib, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBCommand(ctx, input, "show ip fib index %s", "VPP IPv4 FIB Routes")
	})

	// Define vpp_show_ip6_fib tool
	toolShowIp6Fib := &mcp.Tool{
		Name: "vpp_show_ip6_fib",
		Description: "Prints all routes in a given pod IPv6 VRF by running 'vppctl show ip6 fib index <idx>' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index",
	}
	mcp.AddTool(vppServer.server, toolShowIp6Fib, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBCommand(ctx, input, "show ip6 fib index %s", "VPP IPv6 FIB Routes")
	})

	// Define vpp_show_ip_fib_prefix tool
	toolShowIpFibPrefix := &mcp.Tool{
		Name: "vpp_show_ip_fib_prefix",
		Description: "Prints information about a specific prefix in a given pod IPv4 VRF by running 'vppctl show ip fib index <idx> <prefix>' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index\n" +
			"- prefix: The IP prefix to query (e.g., 10.0.0.0/24)",
	}
	mcp.AddTool(vppServer.server, toolShowIpFibPrefix, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBPrefixInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBPrefixCommand(ctx, input, "show ip fib index %s %s", "VPP IPv4 FIB Prefix Information")
	})

	// Define vpp_show_ip6_fib_prefix tool
	toolShowIp6FibPrefix := &mcp.Tool{
		Name: "vpp_show_ip6_fib_prefix",
		Description: "Prints information about a specific prefix in a given pod IPv6 VRF by running 'vppctl show ip6 fib index <idx> <prefix>' in a Kubernetes VPP container\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP\n" +
			"- fib_index: The FIB table index\n" +
			"- prefix: The IPv6 prefix to query (e.g., 2001:db8::/32)",
	}
	mcp.AddTool(vppServer.server, toolShowIp6FibPrefix, func(ctx context.Context, req *mcp.CallToolRequest, input VPPFIBPrefixInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPFIBPrefixCommand(ctx, input, "show ip6 fib index %s %s", "VPP IPv6 FIB Prefix Information")
	})

	// Define bgp_show_neighbors tool
	toolBgpShowNeighbors := &mcp.Tool{
		Name: "bgp_show_neighbors",
		Description: "Show BGP peers by running 'gobgp neighbor' in the agent container of a calico-vpp pod\n\n" +
			"Required parameters:\n" +
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
		Description: "Show BGP global information by running 'gobgp global' in the agent container of a calico-vpp pod\n\n" +
			"Required parameters:\n" +
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
		Description: "Show BGP IPv4 RIB information by running 'gobgp global rib -a 4' in the agent container of a calico-vpp pod\n\n" +
			"Required parameters:\n" +
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
		Description: "Show BGP IPv6 RIB information by running 'gobgp global rib -a 6' in the agent container of a calico-vpp pod\n\n" +
			"Required parameters:\n" +
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
		Description: "Show BGP RIB entry for a specific IP by running 'gobgp global rib <ip>' in the agent container of a calico-vpp pod\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n" +
			"- ip: The IP address to query\n\n" +
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
		Description: "Show BGP RIB entry for a specific prefix by running 'gobgp global rib <prefix>' in the agent container of a calico-vpp pod\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n" +
			"- prefix: The prefix to query (e.g., 10.0.0.0/24)\n\n" +
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
		Description: "Show detailed information for a specific BGP neighbor by running 'gobgp neighbor <neighborIP>' in the agent container of a calico-vpp pod\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running the agent container with gobgp\n" +
			"- neighbor_ip: The IP address of the BGP neighbor\n\n" +
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
