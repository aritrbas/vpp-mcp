# Quick Start Guide

## Run Tests

```bash
# Test everything
./demo_test.sh calico-vpp-node-hnk97

# Test one tool
./test_tool.sh vpp_show_int calico-vpp-node-hnk97
```

## Available Tools

1. **vpp_show_version** - VPP version info
2. **vpp_show_int** - Interface statistics  
3. **vpp_show_int_addr** - Interface addresses
4. **vpp_show_errors** - Error counters
5. **vpp_show_session_verbose** - Session info
6. **vpp_show_npol_rules** - Network policy rules
7. **vpp_show_npol_policies** - Network policy summaries
8. **vpp_trace** - Packet trace capture
9. **vpp_pcap** - PCAP capture
10. **vpp_dispatch** - Dispatch trace capture

## Use with Claude Desktop

1. Edit `~/.config/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "vpp-debug": {
      "command": "/home/aritrbas/vpp/vpp-mcp/vpp-mcp-server"
    }
  }
}
```

2. Restart Claude Desktop

3. Ask Claude:
```
Show VPP interfaces for pod calico-vpp-node-hnk97
```

## Use with MCP Inspector

```bash
npx @modelcontextprotocol/inspector $(pwd)/vpp-mcp-server
```

## Your VPP Pods

```bash
kubectl get pods -n calico-vpp-dataplane
```

Found:
- calico-vpp-node-hnk97
- calico-vpp-node-mhmdq  
- calico-vpp-node-tf9x5

## Manual Test

```bash
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'
  sleep 0.5
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"vpp_show_int","arguments":{"pod_name":"calico-vpp-node-hnk97"}}}'
  sleep 2
) | ./vpp-mcp-server 2>/dev/null | jq -r 'select(.id==2) | .result.content[].text'
```

## Rebuild

```bash
go build -o vpp-mcp-server main.go
```

## Direct kubectl (Bypass MCP)

```bash
kubectl exec -n calico-vpp-dataplane calico-vpp-node-hnk97 -c vpp -- vppctl show int
```

## More Info

- See `TEST_SUMMARY.md` for test results
- See `README.md` for full documentation
