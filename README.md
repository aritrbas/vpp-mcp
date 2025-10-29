# VPP MCP Server (Go Implementation)

A Model Context Protocol (MCP) server written in Go for debugging Vector Packet Processing (VPP) instances running in Kubernetes pods.

## Overview

This MCP server provides tools to interact with VPP instances for debugging purposes. It executes VPP commands on Kubernetes pods and exposes VPP functionality through MCP tools that can be used by AI agents and other MCP clients.

## Features

- **Kubernetes Integration**: Executes VPP commands on Kubernetes pods running VPP
- **Multiple Transport Modes**: 
  - **Stdio** for local client-server communication
  - **HTTP/SSE** for remote network access between machines
- **Multiple VPP Tools**: Ten debugging tools for VPP inspection
  - Version information
  - Interface statistics
  - Interface addresses
  - Error counters
  - Session information
  - NPOL rules and policies
  - Packet trace capture
  - PCAP capture
  - Dispatch trace capture
- **Official MCP Go SDK**: Uses the official Model Context Protocol Go SDK maintained by Google
- **Go Implementation**: Fast, efficient, and easy to deploy
- **Extensible Architecture**: Easy to add more VPP debugging tools
- **Remote Access**: Connect from any machine to debug VPP instances on remote servers

## Prerequisites

- Go 1.21+
- kubectl installed and configured with access to your Kubernetes cluster
- VPP running in Kubernetes pods (e.g., Calico VPP dataplane)
- MCP client (like Claude Desktop, Cline, or other MCP-compatible tools)

## Installation

1. Clone or navigate to the project directory:
```bash
cd /home/aritrbas/vpp/vpp-mcp
```

2. Download Go dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
go build -o vpp-mcp-server main.go
```

## Usage

### Running the MCP Server

The server supports two transport modes: **stdio** (local) and **http** (network).

#### Stdio Transport (Local)

Start the server using stdio transport (default):
```bash
./vpp-mcp-server
```

Or with explicit flag:
```bash
./vpp-mcp-server --transport=stdio
```

Or run directly with Go:
```bash
go run main.go
```

#### HTTP Transport (Network)

Start the server with HTTP transport for remote access:
```bash
./vpp-mcp-server --transport=http --port=8080
```

This exposes the following endpoints:
- **`http://localhost:8080/sse`** - MCP SSE endpoint for client connections
- **`http://localhost:8080/health`** - Health check endpoint
- **`http://localhost:8080/`** - Server information page

For remote access, replace `localhost` with the server's IP address or hostname.

### Available Tools

#### `vpp_show_version`
- **Description**: Get VPP version information
- **Command**: `vppctl show version`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `namespace` (optional): Kubernetes namespace (default: `calico-vpp-dataplane`)
  - `container_name` (optional): Container name within the pod (default: `vpp`)

#### `vpp_show_int`
- **Description**: Get VPP interface information
- **Command**: `vppctl show int`
- **Parameters**: Same as `vpp_show_version`

#### `vpp_show_int_addr`
- **Description**: Get VPP interface address information
- **Command**: `vppctl show int addr`
- **Parameters**: Same as `vpp_show_version`

#### `vpp_show_errors`
- **Description**: Get VPP error counters
- **Command**: `vppctl show errors`
- **Parameters**: Same as `vpp_show_version`

#### `vpp_show_session_verbose`
- **Description**: Get VPP session information with verbose output
- **Command**: `vppctl show session verbose 2`
- **Parameters**: Same as `vpp_show_version`

#### `vpp_show_npol_rules`
- **Description**: Get VPP NPOL rules
- **Command**: `vppctl show npol rules`
- **Parameters**: Same as `vpp_show_version`

#### `vpp_show_npol_policies`
- **Description**: Get VPP NPOL policies
- **Command**: `vppctl show npol policies`
- **Parameters**: Same as `vpp_show_version`

#### `vpp_trace`
- **Description**: Capture VPP packet traces
- **Command**: `vppctl trace add`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `namespace` (optional): Kubernetes namespace (default: `calico-vpp-dataplane`)
  - `container_name` (optional): Container name within the pod (default: `vpp`)
  - `node_name` (optional): Kubernetes node name (validated against cluster)
  - `count` (optional): Number of packets to capture (default: 500)
  - `interface` (optional): Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)

#### `vpp_pcap`
- **Description**: Capture VPP packets to pcap file
- **Command**: `vppctl pcap trace`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `namespace` (optional): Kubernetes namespace (default: `calico-vpp-dataplane`)
  - `container_name` (optional): Container name within the pod (default: `vpp`)
  - `node_name` (optional): Kubernetes node name (validated against cluster)
  - `count` (optional): Number of packets to capture (default: 500)
  - `interface` (optional): Interface name (e.g., host-eth0) or 'any' (default: first available interface)

#### `vpp_dispatch`
- **Description**: Capture VPP dispatch trace to pcap file
- **Command**: `vppctl pcap dispatch trace`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `namespace` (optional): Kubernetes namespace (default: `calico-vpp-dataplane`)
  - `container_name` (optional): Container name within the pod (default: `vpp`)
  - `node_name` (optional): Kubernetes node name (validated against cluster)
  - `count` (optional): Number of packets to capture (default: 500)
  - `interface` (optional): Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)

### Kubernetes Pod Management

The server executes VPP commands on existing Kubernetes pods:
1. Connects to specified pods via kubectl
2. Executes vppctl commands in the VPP container
3. Returns results via MCP protocol

### Pod Configuration

You need to provide:
- **Pod Name**: The name of the Kubernetes pod running VPP (required)
- **Namespace**: Kubernetes namespace (default: `calico-vpp-dataplane`)
- **Container Name**: Container within the pod (default: `vpp`)

## Configuration

### MCP Client Configuration

#### Local Configuration (Stdio Transport)

To use this server with an MCP client on the same machine, add it to your client's configuration. For example, with Claude Desktop, add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "vpp-debug": {
      "command": "/home/aritrbas/vpp/vpp-mcp/vpp-mcp-server",
      "cwd": "/home/aritrbas/vpp/vpp-mcp"
    }
  }
}
```

#### Remote Configuration (HTTP Transport)

For remote access from **Machine Y** to **Machine X**:

**On Machine X (Server):**
1. Start the server with HTTP transport:
```bash
./vpp-mcp-server --transport=http --port=8080
```

2. Ensure the port is accessible (check firewall rules):
```bash
# Example: Allow port 8080 on Ubuntu/Debian
sudo ufw allow 8080/tcp
```

**On Machine Y (Client):**
Configure your MCP client to connect to the HTTP endpoint. For example, with Claude Desktop:

```json
{
  "mcpServers": {
    "vpp-debug-remote": {
      "url": "http://<machine-x-ip>:8080/sse",
      "transport": "sse"
    }
  }
}
```

Replace `<machine-x-ip>` with the actual IP address or hostname of Machine X.

**Security Considerations:**
- The HTTP transport does not include authentication by default
- For production use, consider adding:
  - Reverse proxy with TLS (nginx, Apache)
  - API authentication (API keys, OAuth)
  - Network security (VPN, SSH tunneling)
  - Firewall rules to restrict access

**Example with SSH Tunnel (Secure Alternative):**
```bash
# On Machine Y, create SSH tunnel
ssh -L 8080:localhost:8080 user@machine-x

# Then configure client to use localhost:8080
```

### Customizing Default Settings

You can modify the constants in `main.go` to:
- Change default namespace
- Change default container name
- Add additional VPP commands as tools

## Development

### Project Structure

```
vpp-mcp/
├── main.go                      # Main MCP server implementation
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Makefile                     # Build automation
├── README.md                    # This file
├── .gitignore                   # Git ignore rules
├── vpp-mcp-server               # Compiled binary
├── docs/                        # Documentation
│   ├── QUICK_START.md          # Quick reference
│   ├── REMOTE_ACCESS.md        # Remote access guide
│   └── TEST_SUMMARY.md         # Test results
├── tests/                       # Test scripts
│   ├── test_mcp_server.sh      # Automated test suite
│   ├── demo_test.sh            # Demo all tools
│   ├── test_tool.sh            # Test individual tools
│   └── test_http_server.sh     # HTTP transport tests
└── examples/                    # Example files
    └── example_mcp_requests.json # JSON-RPC examples
```

### Building

Build the server:
```bash
go build -o vpp-mcp-server main.go
```

Build for different platforms:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o vpp-mcp-server-linux main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o vpp-mcp-server-macos main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o vpp-mcp-server.exe main.go
```

### Adding New Tools

To add new VPP debugging tools:

1. Define your tool input structure:
```go
type YourToolInput struct {
    Parameter string `json:"parameter"`
}
```

2. Create a tool handler function:
```go
func (s *VPPMCPServer) handleYourTool(ctx context.Context, req *mcp.CallToolRequest, input YourToolInput) (*mcp.CallToolResult, any, error) {
    genericInput := VPPCommandInput{
        PodName:       input.PodName,
        Namespace:     input.Namespace,
        ContainerName: input.ContainerName,
    }

    return s.handleVPPCommand(ctx, genericInput, "your vppctl subcommand", "Your tool description")
}
```

3. Add the tool to the server in main():
```go
tool := &mcp.Tool{
    Name:        "your_tool_name",
    Description: "Tool description",
}
mcp.AddTool(vppServer.server, tool, vppServer.handleYourTool)
```

### Testing

Test the server functionality:

1. Run automated tests:
```bash
./tests/test_mcp_server.sh
```

2. Demo all tools:
```bash
./tests/demo_test.sh <pod-name>
```

3. Test individual tool:
```bash
./tests/test_tool.sh vpp_show_int <pod-name>
```

### Dependencies

This project uses:
- `github.com/modelcontextprotocol/go-sdk` - Official Model Context Protocol Go SDK maintained in collaboration with Google
- Standard Go libraries for Docker command execution

## Troubleshooting

### Common Issues

1. **kubectl access fails**:
   - Verify kubectl is installed and configured
   - Check you have access to the VPP namespace
   - Ensure proper RBAC permissions for pod exec

2. **vppctl commands fail**:
   - Verify VPP is running in the target pod
   - Check the pod name and namespace are correct
   - Ensure the container name is correct (default: 'vpp')

3. **MCP connection issues**:
   - Verify the binary is built correctly (`go build`)
   - Check MCP client configuration
   - Review server logs for errors

4. **Build issues**:
   - Ensure Go 1.21+ is installed
   - Run `go mod tidy` to download dependencies
   - Check for any compilation errors

### Logs

The server logs important events to help with debugging:
- Container lifecycle events
- Command execution results
- Error conditions

View logs by running the server and monitoring stdout/stderr output.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Future Enhancements

Planned features:
- More VPP debugging tools (`show ip fib`, `show adj`, etc.)
- Configuration management tools
- Log analysis capabilities
- Performance monitoring tools
- Configuration file support