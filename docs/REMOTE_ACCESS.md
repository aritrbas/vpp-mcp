# Remote Access Guide

This guide explains how to expose the VPP MCP Server over HTTP/SSE so an agent on another machine can use it.

## Start server on host machine

```bash
cd /home/aritrbas/vpp/vpp-mcp
make build
./vpp-mcp-server --transport=http --port=8080
```

Server endpoints:
- `/sse` (MCP SSE transport)
- `/health` (returns `OK`)
- `/` (basic info page)

## Verify connectivity from client machine

```bash
curl http://<server-ip>:8080/health
```

Expected response: `OK`

## Client config examples

### Windsurf

Add to `.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "vpp-mcp-remote": {
      "serverUrl": "http://<server-ip>:8080/sse"
    }
  }
}
```

### VS Code MCP config

```json
{
  "servers": {
    "vpp-mcp-remote": {
      "type": "sse",
      "url": "http://<server-ip>:8080/sse"
    }
  }
}
```

## Security recommendations

1. Prefer SSH tunnel or VPN over direct exposure.
2. If exposed, restrict source IPs via firewall.
3. Add TLS termination with reverse proxy for non-lab environments.

## SSH tunnel option (recommended)

On client machine:

```bash
ssh -L 8080:localhost:8080 user@<server-ip>
```

Then point MCP client to:

```text
http://localhost:8080/sse
```

## Tool availability over remote HTTP

Remote clients can use the full current toolset (**42 tools**):
- 35 VPP tools (including `vpp_show_daemonset_image`)
- 7 BGP tools

## Troubleshooting

- Connection refused:
  - Verify server process is running.
  - Check firewall rules and listening port.
- SSE session fails:
  - Confirm client points to `/sse`.
  - Check server logs for connection errors.
- Kubernetes tool failures:
  - Verify `kubectl` context and permissions on the server machine.
