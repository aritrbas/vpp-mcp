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
	"k8s.io/client-go/tools/clientcmd"
)

// ExecutePodVPPCommand runs a VPP command directly on a specified Kubernetes pod
func ExecutePodVPPCommand(ctx context.Context, podName, namespace, containerName, command string) (map[string]interface{}, error) {
	// Apply default values if not provided
	if namespace == "" {
		namespace = "calico-vpp-dataplane" // Default namespace
	}

	if containerName == "" {
		containerName = "vpp" // Default container name
	}

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

	var err error = execErr

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

// getAvailableNodeNames retrieves all node names from the cluster
func (k *KubeClient) getAvailableNodeNames() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), k.timeout)
	defer cancel()

	nodes, err := k.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeNames []string
	for _, node := range nodes.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	return nodeNames, nil
}

// validateNodeName validates that the provided node name exists in the cluster
func validateNodeName(k *KubeClient, nodeName string) (string, error) {
	nodeNames, err := k.getAvailableNodeNames()
	if err != nil {
		return "", err
	}

	if len(nodeNames) == 0 {
		return "", fmt.Errorf("no nodes found. Is cluster running?")
	}

	if nodeName == "" && len(nodeNames) == 1 {
		return nodeNames[0], nil
	}

	for _, n := range nodeNames {
		if n == nodeName {
			return nodeName, nil
		}
	}

	var nodeList strings.Builder
	nodeList.WriteString("\nAvailable nodes:")
	for i, n := range nodeNames {
		nodeList.WriteString(fmt.Sprintf("\n%d. %s", i+1, n))
	}

	return "", fmt.Errorf("node '%s' not found.%s", nodeName, nodeList.String())
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
	// Namespace specifies the Kubernetes namespace where the pod is running (default: calico-vpp-dataplane)
	Namespace string `json:"namespace,omitempty"`
	// ContainerName specifies the container name within the VPP pod (default: vpp)
	ContainerName string `json:"container_name,omitempty"`
}

// VPPGetPodsInput represents the input for getting calico-vpp pods
type VPPGetPodsInput struct {
	// No arguments needed for this tool
}

// VPPCaptureInput represents the input for VPP packet capture tools (trace, pcap, dispatch)
type VPPCaptureInput struct {
	// PodName specifies the name of the Kubernetes pod running VPP
	PodName string `json:"pod_name"`
	// Namespace specifies the Kubernetes namespace where the pod is running (default: calico-vpp-dataplane)
	Namespace string `json:"namespace,omitempty"`
	// ContainerName specifies the container name within the VPP pod (default: vpp)
	ContainerName string `json:"container_name,omitempty"`
	// NodeName specifies the Kubernetes node name (optional, validated against cluster)
	NodeName string `json:"node_name,omitempty"`
	// Count specifies the number of packets to capture (default: run for 15 seconds)
	Count int `json:"count,omitempty"`
	// Interface specifies the interface type or name to capture from
	Interface string `json:"interface,omitempty"`
}

// VPPMCPServer implements the MCP server for VPP debugging
type VPPMCPServer struct {
	server *mcp.Server
}

// NewVPPMCPServer creates a new VPP MCP server
func NewVPPMCPServer() *VPPMCPServer {
	return &VPPMCPServer{}
}

// handleGetPods implements listing all calico-vpp pods with IPs and nodes
func (s *VPPMCPServer) handleGetPods(ctx context.Context, input VPPGetPodsInput) (*mcp.CallToolResult, any, error) {
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
	result, err := ExecutePodVPPCommand(ctx,
		input.PodName,
		input.Namespace,
		input.ContainerName,
		command)

	log.Printf("Command execution completed, processing results...")
	if err != nil {
		log.Printf("Error executing VPP command: %v", err)
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		cmd := result["command"].(string)
		pod := result["pod"].(string)
		namespace := result["namespace"].(string)

		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("%s:\n\n%s\n\nCommand executed: vppctl %s\nPod: %s\nNamespace: %s",
						commandDescription, output, cmd, pod, namespace),
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

	// Validate node name if provided
	if input.NodeName != "" {
		_, err := validateNodeName(k8sClient, input.NodeName)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error validating node: %v", err),
					},
				},
			}, nil, err
		}
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
	_, err = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "clear trace")
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
	_, err = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, traceCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture (15 seconds or until count is reached)
	log.Printf("Capturing packets for 15 seconds or until %d packets captured...", count)
	time.Sleep(15 * time.Second)

	// Step 4: Get trace results
	log.Printf("Retrieving trace results...")
	result, err := ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "show trace")
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
	_, _ = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "clear trace")

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP Trace Capture Results:\n\n%s\n\nCapture Parameters:\n- Node: %s\n- Count: %d\n- Pod: %s\n- Namespace: %s",
						output, vppInputNode, count, input.PodName, input.Namespace),
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

	// Validate node name if provided
	if input.NodeName != "" {
		_, err := validateNodeName(k8sClient, input.NodeName)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error validating node: %v", err),
					},
				},
			}, nil, err
		}
	}

	// Get list of available interfaces
	interfaceResult, err := ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "show int")
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
	if interfaceName == "" || interfaceName == "any" {
		// Use first available interface
		interfaceName = availableInterfaces[0]
	} else {
		// Validate provided interface
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
	_, _ = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "pcap trace off")

	// Step 2: Start pcap capture
	pcapFile := fmt.Sprintf("/tmp/vpp-capture-%d.pcap", time.Now().Unix())
	pcapCmd := fmt.Sprintf("pcap trace tx rx max %d intfc %s file %s", count, interfaceName, pcapFile)
	log.Printf("Starting pcap: %s", pcapCmd)
	_, err = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, pcapCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting pcap: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture (15 seconds or until count is reached)
	log.Printf("Capturing packets for 15 seconds or until %d packets captured...", count)
	time.Sleep(15 * time.Second)

	// Step 4: Stop pcap capture
	log.Printf("Stopping pcap capture...")
	_, err = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "pcap trace off")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error stopping pcap: %v", err),
				},
			},
		}, nil, err
	}

	// Step 5: Get pcap trace status
	result, err := ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "show pcap")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error retrieving pcap status: %v", err),
				},
			},
		}, nil, err
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP PCAP Capture Results:\n\n%s\n\nCapture Parameters:\n- Interface: %s\n- Count: %d\n- File: %s\n- Pod: %s\n- Namespace: %s\n\nNote: Capture file saved at %s on the pod",
						output, interfaceName, count, pcapFile, input.PodName, input.Namespace, pcapFile),
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

	// Validate node name if provided
	if input.NodeName != "" {
		_, err := validateNodeName(k8sClient, input.NodeName)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error validating node: %v", err),
					},
				},
			}, nil, err
		}
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
	_, _ = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "pcap dispatch trace off")

	// Step 2: Start dispatch trace capture
	pcapFile := fmt.Sprintf("/tmp/vpp-dispatch-%d.pcap", time.Now().Unix())
	dispatchCmd := fmt.Sprintf("pcap dispatch trace on max %d buffer-trace %s %d file %s", count, vppInputNode, count, pcapFile)
	log.Printf("Starting dispatch trace: %s", dispatchCmd)
	_, err = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, dispatchCmd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error starting dispatch trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 3: Wait for capture (15 seconds or until count is reached)
	log.Printf("Capturing packets for 15 seconds or until %d packets captured...", count)
	time.Sleep(15 * time.Second)

	// Step 4: Stop dispatch trace
	log.Printf("Stopping dispatch trace...")
	_, err = ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "pcap dispatch trace off")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error stopping dispatch trace: %v", err),
				},
			},
		}, nil, err
	}

	// Step 5: Get dispatch trace status
	result, err := ExecutePodVPPCommand(ctx, input.PodName, input.Namespace, input.ContainerName, "show pcap")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error retrieving dispatch trace status: %v", err),
				},
			},
		}, nil, err
	}

	if success, ok := result["success"].(bool); ok && success {
		output := result["output"].(string)
		response := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("VPP Dispatch Trace Results:\n\n%s\n\nCapture Parameters:\n- Node: %s\n- Count: %d\n- File: %s\n- Pod: %s\n- Namespace: %s\n\nNote: Capture file saved at %s on the pod",
						output, vppInputNode, count, pcapFile, input.PodName, input.Namespace, pcapFile),
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
			"- node_name: Kubernetes node name (validated against cluster)\n" +
			"- count: Number of packets to capture (default: 500)\n" +
			"- interface: Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)\n\n" +
			"The tool will:\n" +
			"1. Clear existing traces\n" +
			"2. Start packet capture\n" +
			"3. Wait 15 seconds or until count is reached\n" +
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
			"- node_name: Kubernetes node name (validated against cluster)\n" +
			"- count: Number of packets to capture (default: 500)\n" +
			"- interface: Interface name (e.g., host-eth0) or 'any' (default: first available interface)\n\n" +
			"The tool will:\n" +
			"1. Validate the interface exists\n" +
			"2. Start pcap capture on tx/rx\n" +
			"3. Wait 15 seconds or until count is reached\n" +
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
			"- node_name: Kubernetes node name (validated against cluster)\n" +
			"- count: Number of packets to capture (default: 500)\n" +
			"- interface: Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)\n\n" +
			"The tool will:\n" +
			"1. Start dispatch trace with buffer trace\n" +
			"2. Wait 15 seconds or until count is reached\n" +
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
	mcp.AddTool(vppServer.server, toolGetPods, func(ctx context.Context, req *mcp.CallToolRequest, input VPPGetPodsInput) (*mcp.CallToolResult, any, error) {
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
		Description: "Lists the active CNAT sessions from the established five tuple to the five tuple rewrites by running 'vppctl cnat session' in a Kubernetes VPP container\n\n" +
			"Output interpretation:\n" +
			"The output shows the `incoming 5-tuple` first that is used to match packets along with the `protocol`. " +
			"Then it displays the `5-tuple after dNAT & sNAT`, followed by the `direction` and finally the `age` in seconds. " +
			"`direction` being input for the PRE-ROUTING sessions and output is the POST-ROUTING sessions\n\n" +
			"Required parameters:\n" +
			"- pod_name: The name of the Kubernetes pod running VPP",
	}
	mcp.AddTool(vppServer.server, toolShowCnatSession, func(ctx context.Context, req *mcp.CallToolRequest, input VPPCommandInput) (*mcp.CallToolResult, any, error) {
		return vppServer.handleVPPCommand(ctx, input, "cnat session", "VPP CNAT Session")
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
	defer session.Close()

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
		w.Write([]byte("OK"))
	})

	// Root endpoint with info
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
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
</html>`))
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
