# CalicoVPP BGP Troubleshooting Workflow

This runbook maps directly to the upstream guide:
- https://github.com/projectcalico/vpp-dataplane/blob/master/docs/bgp/troubleshooting.md

It is tailored for this MCP server and for CalicoVPP where VPP runs in Kubernetes pods.

## Scope

Use this workflow when pods have missing BGP peers, peers stuck outside `Establ`, or routes not showing up in VPP/BGP tables.

## Step 1: Cluster-wide peering snapshot

1. Call `vpp_get_pods` to discover CalicoVPP pods.

2. For **each** pod returned, call `bgp_show_neighbors` with `pod_name=<pod-name>` to collect per-node BGP peering status.

3. Aggregate the results across all pods to build a cluster-wide health picture.

Healthy expectation:
- Every peer on every pod is `Establ`
- `Accepted` is greater than `0`
- Each node should have N-1 IPv4 peers and N-1 IPv6 peers (for N nodes)

## Step 2: Per-node peering details (`gobgp neigh`)

For each affected pod, call `bgp_show_neighbors` with:
- `pod_name=<pod-name>`

Interpretation (same as upstream guide):
- `Establ`: session is up
- `Opened`/`Idle`/`Active` or `Accepted=0`: peering issue or no accepted routes
- Missing peers: likely control-plane or watcher/API path problem

## Step 3: Node BGP identity (`gobgp global`)

Call `bgp_show_global_info` with:
- `pod_name=<pod-name>`

Validate:
- Correct `AS`
- Correct `Router-ID`
- Listening addresses include expected IPv4/IPv6 values

## Step 4: Route visibility (`gobgp global rib -a 4/-a 6`)

Call:
- `bgp_show_global_rib4` with `pod_name=<pod-name>`
- `bgp_show_global_rib6` with `pod_name=<pod-name>`

Validate:
- Expected node/pod prefixes are present
- Next hops match expected peers

## Step 5: Targeted prefix/IP lookups

Call:
- `bgp_show_ip` with `pod_name=<pod-name>`, `parameter=11.0.0.7`
- `bgp_show_prefix` with `pod_name=<pod-name>`, `parameter=11.0.0.0/8`
- `bgp_show_neighbor` with `pod_name=<pod-name>`, `parameter=172.18.0.4`

Use these when one route or one peer is suspected.

## Step 6: If peers are missing or unstable, inspect container logs

Call `get_agent_logs` with:
- `pod_name=<pod-name>`
- `tail_lines=300` (or any suitable value)

Look for:
- Kubernetes API/watch/list errors
- Repeated reconnect loops
- Neighbor bring-up or policy errors

If VPP itself is crashing or misbehaving, also call `get_vpp_manager_logs` with:
- `pod_name=<pod-name>`
- `tail_lines=300`

Look for:
- VPP startup failures or panics
- vpp-manager daemon errors
- Interface or driver initialization issues

## Step 7: Correlate with VPP dataplane tables (optional but recommended)

If BGP looks correct but forwarding still fails:
- Call `vpp_show_ip_fib` with `pod_name=<pod-name>`, `fib_index=0`
- Call `vpp_show_ip6_fib` with `pod_name=<pod-name>`, `fib_index=0`
- Call `vpp_show_int_addr` with `pod_name=<pod-name>`

Look for unresolved/missing routes or incorrect interface addressing.

## Fast Triage Order

1. `vpp_get_pods` (discover all pods)
2. `bgp_show_neighbors` on each pod (build cluster-wide view)
3. `bgp_show_global_info` (on degraded pods)
4. `bgp_show_global_rib4` + `bgp_show_global_rib6`
5. `get_agent_logs` (when peers are missing or unstable)
6. `get_vpp_manager_logs` (when VPP itself is crashing)
7. `vpp_show_ip_fib`/`vpp_show_ip6_fib` if forwarding mismatch persists
