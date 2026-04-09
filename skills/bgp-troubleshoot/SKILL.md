---
name: bgp-troubleshoot
description: CalicoVPP BGP troubleshooting workflow mapped to upstream bgp/troubleshooting.md
---

# CalicoVPP BGP Troubleshooting Workflow

This runbook maps directly to the upstream guide:
- https://github.com/projectcalico/vpp-dataplane/blob/master/docs/bgp/troubleshooting.md

It is tailored for this MCP server and for CalicoVPP where VPP runs in Kubernetes pods.

## Scope

Use this workflow when pods have missing BGP peers, peers stuck outside `Establ`, or routes not showing up in VPP/BGP tables.

## Step 1: Cluster-wide peering snapshot

1. Call `cluster` with `command=get_pods` to discover CalicoVPP pods.

2. For **each** pod returned, call `gobgp` with `command=neighbor`, `pod_name=<pod-name>` to collect per-node BGP peering status.

3. Aggregate the results across all pods to build a cluster-wide health picture.

Healthy expectation:
- Every peer on every pod is `Establ`
- `Accepted` is greater than `0`
- Each node should have N-1 IPv4 peers and N-1 IPv6 peers (for N nodes)

## Step 2: Per-node peering details (`gobgp neigh`)

For each affected pod, call `gobgp` with:
- `command=neighbor`
- `pod_name=<pod-name>`

Interpretation (same as upstream guide):
- `Establ`: session is up
- `Opened`/`Idle`/`Active` or `Accepted=0`: peering issue or no accepted routes
- Missing peers: likely control-plane or watcher/API path problem

## Step 3: Node BGP identity (`gobgp global`)

Call `gobgp` with:
- `command=global`
- `pod_name=<pod-name>`

Validate:
- Correct `AS`
- Correct `Router-ID`
- Listening addresses include expected IPv4/IPv6 values

## Step 4: Route visibility (`gobgp global rib -a 4/-a 6`)

Call `gobgp` with:
- `command=global rib -a ipv4`, `pod_name=<pod-name>`
- `command=global rib -a ipv6`, `pod_name=<pod-name>`

Validate:
- Expected node/pod prefixes are present
- Next hops match expected peers

## Step 5: Targeted prefix/IP lookups

Call `gobgp` with:
- `command=global rib 11.0.0.7`, `pod_name=<pod-name>`
- `command=global rib 11.0.0.0/8`, `pod_name=<pod-name>`
- `command=neighbor 172.18.0.4`, `pod_name=<pod-name>`

Use these when one route or one peer is suspected.

## Step 6: If peers are missing or unstable, inspect container logs

Call `cluster` with:
- `command=logs`, `pod_name=<pod-name>`, `container=agent`, `tail_lines=300`

Look for:
- Kubernetes API/watch/list errors
- Repeated reconnect loops
- Neighbor bring-up or policy errors

If VPP itself is crashing or misbehaving, also call `cluster` with:
- `command=logs`, `pod_name=<pod-name>`, `container=vpp`, `tail_lines=300`

Look for:
- VPP startup failures or panics
- vpp-manager daemon errors
- Interface or driver initialization issues

## Step 7: Correlate with VPP dataplane tables (optional but recommended)

If BGP looks correct but forwarding still fails:
- Call `vppctl` with `command=show ip fib index 0`, `pod_name=<pod-name>`
- Call `vppctl` with `command=show ip6 fib index 0`, `pod_name=<pod-name>`
- Call `vppctl` with `command=show int addr`, `pod_name=<pod-name>`

Look for unresolved/missing routes or incorrect interface addressing.

## Fast Triage Order

1. `cluster` → `get_pods` (discover all pods)
2. `gobgp` → `neighbor` on each pod (build cluster-wide view)
3. `gobgp` → `global` (on degraded pods)
4. `gobgp` → `global rib -a ipv4` + `global rib -a ipv6`
5. `cluster` → `logs` with `container=agent` (when peers are missing or unstable)
6. `cluster` → `logs` with `container=vpp` (when VPP itself is crashing)
7. `vppctl` → `show ip fib index 0` / `show ip6 fib index 0` if forwarding mismatch persists
