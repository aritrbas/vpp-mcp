# VPP MCP Server Test Summary

## Current Branch Status

- MCP server starts in `stdio` and `http` modes.
- Tool inventory is **4 broad tools**: `list`, `cluster`, `vppctl`, `gobgp`
- Each tool accepts a `command` parameter; the agent/LLM decides which commands to run
- Cluster tool covers: get_pods, get_nodes, get_configmap, get_daemonset, get_events, get_deployments, get_services, get_endpoints, get_replicaset, get_ippool, get_namespaces, top_pods, top_nodes, logs, describe_pod, describe_node, exec (read-only)

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
./tests/test_tool.sh <pod-name> vppctl "show version"
./tests/test_tool.sh cluster get_pods
./tests/test_tool.sh list
./tests/test_tool.sh <pod-name> gobgp neighbor

# Full demo pass across all tools
./tests/demo_test.sh <pod-name>
```

## Notes

- Capture commands (`trace`, `pcap`, `dispatch`, `capture_cleanup`) are mutating — pass via `vppctl` tool.
- Stats-reset commands (`clear errors`, `clear run`) are mutating — pass via `vppctl` tool.
- `gobgp` tool executes in the `agent` container of the calico-vpp pod.
- Kubernetes connectivity (`kubectl` context + RBAC) is required for `cluster` and `gobgp` tools.

## Expected Outputs

- `./tests/test_mcp_server.sh` prints discovered tool count from `tools/list`.
- `./tests/test_http_server.sh` verifies `/`, `/health`, and `/sse` endpoint availability.
- `./tests/test_tool.sh` prints `=== Tool Result ===` with command output.
