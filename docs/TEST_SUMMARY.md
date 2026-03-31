# VPP MCP Server Test Summary

## Current Branch Status

- MCP server starts in `stdio` and `http` modes.
- Tool inventory is **44 tools** total:
  - **35 VPP tools**
  - **9 BGP tools**
- Includes Kubernetes health-check support tools:
  - `vpp_get_pods`
  - `vpp_show_daemonset_image`
  - `bgp_cluster_show_neighbors`
  - `bgp_get_agent_logs`

## Verification Commands

Run from project root:

```bash
# Build
make build

# stdio initialization + tools/list
./tests/test_mcp_server.sh

# HTTP endpoint checks
./tests/test_http_server.sh

# Single tool checks
./tests/test_tool.sh <pod-name> vpp_show_version
./tests/test_tool.sh vpp_get_pods
./tests/test_tool.sh vpp_show_daemonset_image
./tests/test_tool.sh bgp_cluster_show_neighbors
./tests/test_tool.sh <pod-name> bgp_get_agent_logs

# Full demo pass across all tools
./tests/demo_test.sh <pod-name>
```

## Notes

- `vpp_trace`, `vpp_pcap`, `vpp_dispatch`, `vpp_capture_cleanup` are mutating capture tools.
- `vpp_clear_errors` and `vpp_clear_run` are mutating stats-reset tools.
- BGP tools execute in the `agent` container of the calico-vpp pod.
- Kubernetes connectivity (`kubectl` context + RBAC) is required for Kubernetes mode tools.

## Expected Outputs

- `./tests/test_mcp_server.sh` prints discovered tool count from `tools/list`.
- `./tests/test_http_server.sh` verifies `/`, `/health`, and `/sse` endpoint availability.
- `./tests/test_tool.sh` prints `=== Tool Result ===` with command output.
