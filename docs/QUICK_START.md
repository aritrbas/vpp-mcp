# Quick Start Guide

## Run Tests

```bash
# Test everything
../tests/demo_test.sh calico-vpp-node-hnk97

# Test one tool
../tests/test_tool.sh vpp_show_int calico-vpp-node-hnk97
```

## Available Tools (34 Total)

### VPP Tools (27)
1. **vpp_show_version** - VPP version info
2. **vpp_show_int** - Interface statistics  
3. **vpp_show_int_addr** - Interface addresses
4. **vpp_show_errors** - Error counters
5. **vpp_show_session_verbose** - Session info
6. **vpp_show_npol_rules** - Network policy rules
7. **vpp_show_npol_policies** - Network policy summaries
8. **vpp_show_npol_ipset** - IPsets referenced by rules
9. **vpp_show_npol_interfaces** - Policies on interfaces
10. **vpp_trace** - Packet trace capture
11. **vpp_pcap** - PCAP capture
12. **vpp_dispatch** - Dispatch trace capture
13. **vpp_get_pods** - List all CalicoVPP pods
14. **vpp_clear_errors** - Reset error counters
15. **vpp_tcp_stats** - TCP statistics
16. **vpp_session_stats** - Session layer statistics
17. **vpp_get_logs** - VPP logs
18. **vpp_show_cnat_translation** - CNAT translations
19. **vpp_show_cnat_session** - CNAT sessions
20. **vpp_clear_run** - Clear runtime statistics
21. **vpp_show_run** - Runtime statistics
22. **vpp_show_ip_table** - IPv4 VRF tables
23. **vpp_show_ip6_table** - IPv6 VRF tables
24. **vpp_show_ip_fib** - IPv4 FIB routes (requires fib_index)
25. **vpp_show_ip6_fib** - IPv6 FIB routes (requires fib_index)
26. **vpp_show_ip_fib_prefix** - IPv4 FIB prefix info (requires fib_index & prefix)
27. **vpp_show_ip6_fib_prefix** - IPv6 FIB prefix info (requires fib_index & prefix)

### BGP Tools (7)
28. **bgp_show_neighbors** - BGP peers
29. **bgp_show_global_info** - BGP global information
30. **bgp_show_global_rib4** - BGP IPv4 RIB
31. **bgp_show_global_rib6** - BGP IPv6 RIB
32. **bgp_show_ip** - BGP RIB entry for IP (requires IP address)
33. **bgp_show_prefix** - BGP RIB entry for prefix (requires prefix)
34. **bgp_show_neighbor** - BGP neighbor details (requires neighbor IP)

**Note**: All tools use namespace `calico-vpp-dataplane` and container `vpp` by default.

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
# Run from project root
cd ..
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
# Run from project root
cd ..
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

- See `TEST_SUMMARY.md` for test results (same directory)
- See `../README.md` for full documentation
