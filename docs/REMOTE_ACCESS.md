# Remote Access Guide

This guide explains how to expose the VPP MCP Server over HTTP to enable remote access from an agent running on a different machine.

## Scenario

- **Machine X**: Server running VPP in Kubernetes (where MCP server will run)
- **Machine Y**: Client machine with an MCP agent (e.g., Claude Desktop, custom agent)

## Quick Start

### Step 1: Start Server on Machine X

```bash
# Navigate to the project directory
cd /home/aritrbas/vpp/vpp-mcp

# Build the server (if not already built)
go build -o vpp-mcp-server main.go

# Start with HTTP transport
./vpp-mcp-server --transport=http --port=8080
```

You should see output like:
```
Starting VPP MCP Server with transport=http...
HTTP server listening on port 8080
MCP SSE endpoint: http://localhost:8080/sse
```

### Step 2: Configure Firewall on Machine X

Ensure the port is accessible from Machine Y:

```bash
# Ubuntu/Debian
sudo ufw allow 8080/tcp
sudo ufw status

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

### Step 3: Test Connectivity

From Machine Y, verify you can reach the server:

```bash
# Test health endpoint
curl http://<machine-x-ip>:8080/health

# Expected response: OK
```

### Step 4: Configure MCP Client on Machine Y

#### For Claude Desktop

Edit your `claude_desktop_config.json`:

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

#### For Custom MCP Clients

Connect to the SSE endpoint:
```
http://<machine-x-ip>:8080/sse
```

The server uses Server-Sent Events (SSE) for the MCP protocol transport.

## Command-Line Options

```bash
./vpp-mcp-server --help

Options:
  --transport string
        Transport mode: stdio or http (default "stdio")
  --port string
        HTTP port (only used when transport=http) (default "8080")
```

## Available Endpoints

When running in HTTP mode, the server exposes:

- **`/sse`** - MCP Server-Sent Events endpoint (main MCP connection point)
- **`/health`** - Health check endpoint (returns "OK")
- **`/`** - Information page (HTML)

## Security Best Practices

### 1. SSH Tunnel (Recommended for Production)

Instead of exposing the port directly, use SSH tunneling:

**On Machine Y:**
```bash
# Create tunnel
ssh -L 8080:localhost:8080 user@machine-x

# Keep this terminal open
```

**Configure client to use localhost:**
```json
{
  "mcpServers": {
    "vpp-debug": {
      "url": "http://localhost:8080/sse",
      "transport": "sse"
    }
  }
}
```

### 2. Reverse Proxy with TLS

Use nginx or Apache as a reverse proxy with TLS:

**nginx example:**
```nginx
server {
    listen 443 ssl;
    server_name vpp-mcp.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /sse {
        proxy_pass http://localhost:8080/sse;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 3. Network Restrictions

Use firewall rules to limit access:

```bash
# Allow only specific IP
sudo ufw allow from <machine-y-ip> to any port 8080

# Or allow a subnet
sudo ufw allow from 192.168.1.0/24 to any port 8080
```

### 4. VPN

Run both machines on a VPN for encrypted communication:
- WireGuard
- OpenVPN
- Tailscale

## Troubleshooting

### Connection Refused

**Problem:** Client cannot connect to server

**Solutions:**
1. Verify server is running: `curl http://<machine-x-ip>:8080/health`
2. Check firewall rules: `sudo ufw status`
3. Verify port is listening: `netstat -tlnp | grep 8080`
4. Check network connectivity: `ping <machine-x-ip>`

### Session Errors

**Problem:** MCP session fails to establish

**Solutions:**
1. Check server logs on Machine X
2. Ensure kubectl is configured on Machine X
3. Verify VPP pods are accessible: `kubectl get pods -n calico-vpp-dataplane`

### Timeout Issues

**Problem:** Connection times out

**Solutions:**
1. Increase timeout in client configuration
2. Check for network proxies or middleware
3. Verify no load balancers are interfering with SSE connections

## Example Usage

Once connected, your agent on Machine Y can use all VPP debugging tools:

```javascript
// Example: Get VPP version from remote pod
{
  "tool": "vpp_show_version",
  "parameters": {
    "pod_name": "calico-vpp-node-abcd",
    "namespace": "calico-vpp-dataplane"
  }
}

// Example: Capture packet trace
{
  "tool": "vpp_trace",
  "parameters": {
    "pod_name": "calico-vpp-node-abcd",
    "count": 100,
    "interface": "virtio"
  }
}
```

## Performance Considerations

- **Latency**: Network latency between machines will affect response times
- **Bandwidth**: Large outputs (traces, pcaps) may take longer over network
- **Concurrent Connections**: Server supports multiple simultaneous SSE connections
- **Keepalive**: SSE connections maintain persistent HTTP connections

## Monitoring

### Server Logs

Monitor server activity:
```bash
./vpp-mcp-server --transport=http --port=8080 2>&1 | tee vpp-mcp.log
```

### Health Checks

Set up monitoring:
```bash
# Simple monitoring script
while true; do
  curl -sf http://localhost:8080/health || echo "Server down!"
  sleep 30
done
```

### systemd Service (Optional)

Create a systemd service for automatic startup:

```ini
[Unit]
Description=VPP MCP Server
After=network.target

[Service]
Type=simple
User=vpp-user
WorkingDirectory=/home/aritrbas/vpp/vpp-mcp
ExecStart=/home/aritrbas/vpp/vpp-mcp/vpp-mcp-server --transport=http --port=8080
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable vpp-mcp-server
sudo systemctl start vpp-mcp-server
sudo systemctl status vpp-mcp-server
```

## Next Steps

1. Test the connection with a simple health check
2. Configure your MCP client on Machine Y
3. Start debugging VPP instances remotely
4. Implement security measures appropriate for your environment

For more information, see the main [README.md](../README.md).
