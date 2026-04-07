# VPP MCP Server (Go Implementation)

A Model Context Protocol (MCP) server written in Go for debugging Vector Packet Processing (VPP) instances running in Kubernetes pods or as standalone daemons.

## Overview

This MCP server provides tools to interact with VPP instances for debugging purposes. It supports dual execution modes: **Kubernetes** (CalicoVPP pods) and **Standalone** (local VPP daemon). It exposes VPP functionality through MCP tools that can be used by AI agents and other MCP clients.

## Features

- **Dual Execution Modes**:
  - **Kubernetes**: Executes VPP commands on Kubernetes pods running CalicoVPP
  - **Standalone**: Executes VPP commands on a local VPP daemon via vppctl socket
- **Multiple Transport Modes**: 
  - **Stdio** for local client-server communication
  - **HTTP/SSE** for remote network access between machines
- **4 Broad Tools** — lightweight MCP surface that leverages agent/LLM intelligence:
  - **`list`**: Discover all available commands across categories
  - **`cluster`**: Kubernetes cluster operations (pods, nodes, configmaps, daemonsets, logs)
  - **`vppctl`**: Run any vppctl command on VPP (1000++ commands supported)
  - **`gobgp`**: Run any gobgp command in the CalicoVPP agent container
- **Official MCP Go SDK**: Uses the official Model Context Protocol Go SDK maintained by Google
- **Go Implementation**: Fast, efficient, and easy to deploy
- **Extensible Architecture**: Easy to add more VPP debugging tools
- **Remote Access**: Connect from any machine to debug VPP instances on remote servers

## Prerequisites

- Go 1.24+
- **Kubernetes mode**: kubectl installed and configured with access to your Kubernetes cluster; VPP running in Kubernetes pods (e.g., Calico VPP dataplane)
- **Standalone mode**: VPP running as a local daemon with vppctl accessible (default socket: `/var/run/vpp/cli.sock`)
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

The server exposes **4 broad tools** that leverage agent/LLM intelligence to decide which specific commands to run. Use the `list` tool to discover all available commands.

#### `list`
- **Description**: Returns a comprehensive command reference for all tool categories, including the complete 1000++ vppctl command catalog (embedded at compile time via `//go:embed`), full gobgp command reference, and all cluster commands with kubectl mappings.
- **Parameters**: None
- **Output size**: ~75KB (includes full vppctl command catalog from `vppctl-cmds/vppctl-commands-summary.md`)

#### `cluster`
- **Description**: Run Kubernetes cluster commands for CalicoVPP troubleshooting
- **Parameters**:
  - `command` (required): Operation to run — `get_pods`, `get_nodes`, `get_configmap`, `get_daemonset`, `get_events`, `get_deployments`, `get_services`, `get_endpoints`, `get_replicaset`, `get_ippool`, `get_namespaces`, `top_pods`, `top_nodes`, `logs`, `describe_pod`, `describe_node`, `exec`
  - `namespace` (optional): Kubernetes namespace (default: `calico-vpp-dataplane`)
  - `pod_name` (optional): Pod name (for `logs`, `describe_pod`, `exec`)
  - `container` (optional): Container name (for `logs`, `exec`; defaults: `agent` for logs, `vpp` for exec)
  - `resource_name` (optional): Resource name (for `get_configmap`, `get_daemonset`, `get_endpoints`, `describe_node`) or command to execute (for `exec`)
  - `tail_lines` (optional): Number of log lines (default: 200, max: 5000)
- **Note**: The `exec` command only allows read-only operations (e.g., `ls`, `cat`, `ip a`). Mutating commands are blocked.

#### `vppctl`
- **Description**: Run any vppctl command on VPP. Passes the command string directly to the VPP CLI.
- **Parameters**:
  - `command` (required): The vppctl command (e.g., `show version`, `show int addr`, `show ip fib index 0 10.0.0.0/24`)
  - `pod_name` (Kubernetes mode, required): Pod running VPP
  - `mode` (optional): `kubernetes` (default) or `standalone`
  - `sock_path` (Standalone mode, optional): VPP socket path
  - `count` / `timeout` / `interface` (optional): For capture commands (`trace`, `pcap`, `dispatch`)
- **Common commands**: `show version`, `show int`, `show int addr`, `show hardware-interfaces`, `show errors`, `show run`, `clear errors`, `clear run`, `show logging`, `show ip table`, `show ip fib`, `show session verbose 2`, `show tcp stats`, `show npol rules/policies/ipset/interfaces`, `show cnat translation/session`, `show ipip tunnel`, `show vxlan tunnel`, `show tun`, `trace`, `pcap`, `dispatch`, `capture_cleanup`

#### `gobgp`
- **Description**: Run any gobgp command in the agent container of a CalicoVPP pod
- **Parameters**:
  - `command` (required): The gobgp command (e.g., `neighbor`, `global rib -a ipv4`, `neighbor 10.0.0.1 adj-in`)
  - `pod_name` (required): Pod running the agent container
- **Common commands**: `neighbor`, `neighbor <ip>`, `neighbor <ip> adj-in/adj-out`, `global`, `global rib -a ipv4/ipv6`, `global rib <prefix>`, `vrf`, `policy prefix/community/as-path/statement`

### CalicoVPP BGP Troubleshooting Workflow

This workflow mirrors the upstream troubleshooting guide and is adapted to this MCP server's toolset:
- Upstream reference: https://github.com/projectcalico/vpp-dataplane/blob/master/docs/bgp/troubleshooting.md
- Local runbook: [skills/BGP_TROUBLESHOOTING_WORKFLOW.md](skills/BGP_TROUBLESHOOTING_WORKFLOW.md)

### Kubernetes Pod Management

The server executes VPP commands on existing Kubernetes pods:
1. Connects to specified pods via kubectl
2. Executes vppctl commands in the VPP container
3. Executes gobgp commands in the agent container
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
├── skills/                      # LLM workflow skills
│   └── bgp-troubleshoot.md      # BGP triage runbook
│   └── cluster-healthcheck.md   # Cluster health check
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

### Adding New Commands

With the 4-tool architecture, adding support for new commands is straightforward:

- **New vppctl commands**: Simply pass any valid vppctl command string to the `vppctl` tool. No code changes needed.
- **New gobgp commands**: Simply pass any valid gobgp command string to the `gobgp` tool. No code changes needed.
- **New cluster operations**: Add a new `case` in the `handleCluster` switch statement in `main.go`.
- **Update the list**: Add the new command to `handleList` so agents can discover it.

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
./tests/test_tool.sh <pod-name> vppctl "show int"
```

### Dependencies

This project uses:
- `github.com/modelcontextprotocol/go-sdk` - Official Model Context Protocol Go SDK maintained in collaboration with Google
- Standard Go libraries for command execution
- `k8s.io/client-go` - Kubernetes client library for pod discovery and ConfigMap access

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
