# Quick Start Guide

## 1) Build

```bash
make build
```

## 2) Smoke test (stdio)

```bash
./tests/test_mcp_server.sh
```

## 3) Test a single tool

```bash
# pod-scoped tools
./tests/test_tool.sh <pod-name> vpp_show_int

# no-pod tools
./tests/test_tool.sh vpp_get_pods
./tests/test_tool.sh vpp_show_daemonset_image
./tests/test_tool.sh bgp_cluster_show_neighbors
```

## 4) Run full demo

```bash
./tests/demo_test.sh <pod-name>
```

## Tool Inventory (Current Branch)

Total: **44 tools**
- **35 VPP tools**
- **9 BGP tools**

### VPP tools (35)
1. `vpp_show_version`
2. `vpp_show_int`
3. `vpp_show_int_addr`
4. `vpp_show_hardware_interfaces`
5. `vpp_show_hardware_interface`
6. `vpp_show_errors`
7. `vpp_show_session_verbose`
8. `vpp_show_npol_rules`
9. `vpp_show_npol_policies`
10. `vpp_show_npol_ipset`
11. `vpp_show_npol_interfaces`
12. `vpp_clear_errors`
13. `vpp_tcp_stats`
14. `vpp_session_stats`
15. `vpp_get_logs`
16. `vpp_show_cnat_translation`
17. `vpp_show_cnat_session`
18. `vpp_clear_run`
19. `vpp_show_run`
20. `vpp_show_ipip_tunnel`
21. `vpp_show_vxlan_tunnel`
22. `vpp_show_tun_all`
23. `vpp_show_tun_interface`
24. `vpp_show_ip_table`
25. `vpp_show_ip6_table`
26. `vpp_show_ip_fib`
27. `vpp_show_ip6_fib`
28. `vpp_show_ip_fib_prefix`
29. `vpp_show_ip6_fib_prefix`
30. `vpp_trace`
31. `vpp_pcap`
32. `vpp_dispatch`
33. `vpp_capture_cleanup`
34. `vpp_get_pods`
35. `vpp_show_daemonset_image`

### BGP tools (9)
1. `bgp_show_neighbors`
2. `bgp_show_global_info`
3. `bgp_show_global_rib4`
4. `bgp_show_global_rib6`
5. `bgp_show_ip`
6. `bgp_show_prefix`
7. `bgp_show_neighbor`
8. `bgp_cluster_show_neighbors`
9. `bgp_get_agent_logs`

## Mode Notes

- Most `vpp_*` command tools support:
  - Kubernetes mode (default): `pod_name` required
  - Standalone mode: set `mode="standalone"` (optional `sock_path`)
- Kubernetes-specific tools:
  - `vpp_get_pods`
  - `vpp_show_daemonset_image`
  - `vpp_show_hardware_interface`
  - `vpp_show_tun_interface`
  - `bgp_cluster_show_neighbors`
  - `bgp_get_agent_logs`

## Manual JSON-RPC example

```bash
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'
  sleep 0.5
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"vpp_show_daemonset_image","arguments":{}}}'
  sleep 1
) | ./vpp-mcp-server
```
