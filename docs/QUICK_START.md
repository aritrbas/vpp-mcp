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
# vppctl tool
./tests/test_tool.sh <pod-name> vppctl "show int"

# cluster tool (no pod needed)
./tests/test_tool.sh cluster get_pods

# list tool
./tests/test_tool.sh list
```

## 4) Run full demo

```bash
./tests/demo_test.sh <pod-name>
```

## Tool Inventory (Current Branch)

Total: **4 broad tools** that leverage agent/LLM intelligence:

1. **`list`** — Returns a comprehensive command reference for all categories
2. **`cluster`** — Kubernetes cluster operations (get_pods, get_nodes, get_configmap, get_daemonset, get_events, get_deployments, get_services, get_endpoints, get_replicaset, get_ippool, get_namespaces, top_pods, top_nodes, logs, describe_pod, describe_node, exec)
3. **`vppctl`** — Run any vppctl command on VPP (1000++ commands, including captures)
4. **`gobgp`** — Run any gobgp command in the CalicoVPP agent container

## Mode Notes

- `vppctl` supports:
  - Kubernetes mode (default): `pod_name` required
  - Standalone mode: set `mode="standalone"` (optional `sock_path`)
- `gobgp` and `cluster` are Kubernetes-only
- Use `list` tool to discover all available commands

## Manual JSON-RPC example

```bash
(
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'
  sleep 0.5
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"cluster","arguments":{"command":"get_pods"}}}'
  sleep 1
) | ./vpp-mcp-server
```
