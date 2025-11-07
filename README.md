# VPP MCP Server (Go Implementation)

A Model Context Protocol (MCP) server written in Go for debugging Vector Packet Processing (VPP) instances running in Kubernetes pods.

## Overview

This MCP server provides tools to interact with VPP instances for debugging purposes. It executes VPP commands on Kubernetes pods and exposes VPP functionality through MCP tools that can be used by AI agents and other MCP clients.

## Features

- **Kubernetes Integration**: Executes VPP commands on Kubernetes pods running VPP
- **Multiple Transport Modes**: 
  - **Stdio** for local client-server communication
  - **HTTP/SSE** for remote network access between machines
- **34 Debugging Tools**: Comprehensive toolset for VPP and BGP debugging
  - Pod management (list all CalicoVPP pods)
  - Version information
  - Interface statistics and addresses
  - Error counters and error clearing
  - Session information and statistics
  - TCP statistics
  - NPOL rules and policies
  - CNAT translations and sessions
  - Runtime statistics
  - IP routing tables and FIBs
  - VPP logs
  - Packet trace, PCAP, and dispatch trace capture
  - BGP neighbors and global information
  - BGP RIB queries (IPv4/IPv6, IPs, prefixes)
- **Official MCP Go SDK**: Uses the official Model Context Protocol Go SDK maintained by Google
- **Go Implementation**: Fast, efficient, and easy to deploy
- **Extensible Architecture**: Easy to add more VPP debugging tools
- **Remote Access**: Connect from any machine to debug VPP instances on remote servers

## Prerequisites

- Go 1.24+
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

**Note**: All VPP tools use namespace `calico-vpp-dataplane` and container `vpp`.

#### `vpp_show_version`
- **Description**: Get VPP version information
- **Command**: `vppctl show version`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_int`
- **Description**: Get VPP interface information
- **Command**: `vppctl show int`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_int_addr`
- **Description**: Get VPP interface address information
- **Command**: `vppctl show int addr`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_errors`
- **Description**: Get VPP error counters
- **Command**: `vppctl show errors`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_session_verbose`
- **Description**: Get VPP session information with verbose output
- **Command**: `vppctl show session verbose 2`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_npol_rules`
- **Description**: List rules that are referenced by policies
- **Command**: `vppctl show npol rules`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_npol_policies`
- **Description**: List all the policies that are referenced on interfaces
- **Command**: `vppctl show npol policies`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_npol_ipset`
- **Description**: List ipsets that are referenced by rules (IPsets are just list of IPs)
- **Command**: `vppctl show npol ipset`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_npol_interfaces`
- **Description**: Show the resulting policies configured for every interface in VPP. The first IPv4 address of every pod is provided to help identify which pod and interface belongs to.
- **Command**: `vppctl show npol interfaces`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
- **Output interpretation**:
  - `tx`: contains rules that are applied on packets that LEAVE VPP on a given interface. Rules are applied top to bottom.
  - `rx`: contains rules that are applied on packets that ENTER VPP on a given interface. Rules are applied top to bottom.
  - `profiles`: are specific rules that are enforced when a matched rule action is PASS or when no policies are configured.

#### `vpp_trace`
- **Description**: Capture VPP packet traces
- **Command**: `vppctl trace add`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `count` (optional): Number of packets to capture (default: 500)
  - `interface` (optional): Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)

#### `vpp_pcap`
- **Description**: Capture VPP packets to pcap file
- **Command**: `vppctl pcap trace`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `count` (optional): Number of packets to capture (default: 500)
  - `interface` (optional): Interface name (e.g., host-eth0) or 'any' (default: 'any')

#### `vpp_dispatch`
- **Description**: Capture VPP dispatch trace to pcap file
- **Command**: `vppctl pcap dispatch trace`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `count` (optional): Number of packets to capture (default: 500)
  - `interface` (optional): Interface type - phy|af_xdp|af_packet|avf|vmxnet3|virtio|rdma|dpdk|memif|vcl (default: virtio)

#### `vpp_get_pods`
- **Description**: List all CalicoVPP pods with their IPs and nodes on which they are running
- **Command**: `kubectl get pods -n calico-vpp-dataplane -owide`
- **Parameters**: None required

#### `vpp_clear_errors`
- **Description**: Reset the error counters
- **Command**: `vppctl clear errors`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_tcp_stats`
- **Description**: Display global statistics reported by TCP
- **Command**: `vppctl show tcp stats`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_session_stats`
- **Description**: Display global statistics reported by the session layer
- **Command**: `vppctl show session stats`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_get_logs`
- **Description**: Display VPP logs
- **Command**: `vppctl show logging`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_cnat_translation`
- **Description**: Shows the active CNAT translations
- **Command**: `vppctl show cnat translation`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_cnat_session`
- **Description**: Lists the active CNAT sessions from the established five tuple to the five tuple rewrites
- **Command**: `vppctl show cnat session`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
- **Output interpretation**: The output shows the `incoming 5-tuple` first that is used to match packets along with the `protocol`. Then it displays the `5-tuple after dNAT & sNAT`, followed by the `direction` and finally the `age` in seconds. `direction` being input for the PRE-ROUTING sessions and output is the POST-ROUTING sessions

#### `vpp_clear_run`
- **Description**: Clears live running error stats in VPP
- **Command**: `vppctl clear run`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_run`
- **Description**: Shows live running error stats in VPP
- **Command**: `vppctl show run`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
- **Debugging workflow**: Sometimes to debug an issue, you might need to run `vpp_clear_run` to erase historic stats and then wait for a few seconds in the issue state / run some tests so that the error stats are repopulated and then run `vpp_show_run` in order to diagnose what is going on in the system
- **Output interpretation**: A loaded VPP will typically have (1) a high Vectors/Call maxing out at 256 (2) a low loops/sec struggling around 10000. The Clocks column tells you the consumption in cycles per node on average. Beyond 1e3 is expensive.

#### `vpp_show_ip_table`
- **Description**: Prints all available IPv4 VRFs
- **Command**: `vppctl show ip table`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_ip6_table`
- **Description**: Prints all available IPv6 VRFs
- **Command**: `vppctl show ip6 table`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP

#### `vpp_show_ip_fib`
- **Description**: Prints all routes in a given pod IPv4 VRF
- **Command**: `vppctl show ip fib index <idx>`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `fib_index` (required): The FIB table index

#### `vpp_show_ip6_fib`
- **Description**: Prints all routes in a given pod IPv6 VRF
- **Command**: `vppctl show ip6 fib index <idx>`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `fib_index` (required): The FIB table index

#### `vpp_show_ip_fib_prefix`
- **Description**: Prints information about a specific prefix in a given pod IPv4 VRF
- **Command**: `vppctl show ip fib index <idx> <prefix>`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `fib_index` (required): The FIB table index
  - `prefix` (required): The IP prefix to query (e.g., 10.0.0.0/24)

#### `vpp_show_ip6_fib_prefix`
- **Description**: Prints information about a specific prefix in a given pod IPv6 VRF
- **Command**: `vppctl show ip6 fib index <idx> <prefix>`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running VPP
  - `fib_index` (required): The FIB table index
  - `prefix` (required): The IPv6 prefix to query (e.g., 2001:db8::/32)

#### `bgp_show_neighbors`
- **Description**: Show BGP peers
- **Command**: `gobgp neighbor`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running the agent container with gobgp

#### `bgp_show_global_info`
- **Description**: Show BGP global information
- **Command**: `gobgp global`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running the agent container with gobgp

#### `bgp_show_global_rib4`
- **Description**: Show BGP IPv4 RIB information
- **Command**: `gobgp global rib -a 4`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running the agent container with gobgp

#### `bgp_show_global_rib6`
- **Description**: Show BGP IPv6 RIB information
- **Command**: `gobgp global rib -a 6`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running the agent container with gobgp

#### `bgp_show_ip`
- **Description**: Show BGP RIB entry for a specific IP
- **Command**: `gobgp global rib <ip>`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running the agent container with gobgp
  - `parameter` (required): The IP address to query

#### `bgp_show_prefix`
- **Description**: Show BGP RIB entry for a specific prefix
- **Command**: `gobgp global rib <prefix>`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running the agent container with gobgp
  - `parameter` (required): The prefix to query (e.g., 10.0.0.0/24)

#### `bgp_show_neighbor`
- **Description**: Show detailed information for a specific BGP neighbor
- **Command**: `gobgp neighbor <neighborIP>`
- **Parameters**:
  - `pod_name` (required): Name of the Kubernetes pod running the agent container with gobgp
  - `parameter` (required): The neighbor IP address  to query

### Kubernetes Pod Management

The server executes VPP commands on existing Kubernetes pods:
1. Connects to specified pods via kubectl
2. Executes vppctl commands in the VPP container
3. Executes gobhp commands in the agent container
4. Returns results via MCP protocol

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
│   ├── test_mcp_server.sh      # Test MCP server setup in stdio transport
│   ├── demo_test.sh            # Demo all tools
│   ├── test_tool.sh            # Test individual tools
│   └── test_http_server.sh     # Test MCP server setup in HTTP transport
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


1. Test server functionality:
```bash
# stdio transport
./tests/test_mcp_server.sh
# HTTP transport
./tests/test_http_server.sh
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
   - Check if the pod name is correct

3. **MCP connection issues**:
   - Verify the binary is built correctly (`go build`)
   - Check MCP client configuration
   - Review server logs for errors

4. **Build issues**:
   - Ensure Go 1.24+ is installed
   - Run `make deps` to download dependencies
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
- More VPP debugging tools
- Configuration management tools
- Log analysis capabilities
- Performance monitoring tools
- Configuration file support
- Workflow support
- Workflow visualization tools
- Automated workflow execution engine