---
name: network-troubleshoot
description: CalicoVPP network/FIB troubleshooting workflow mapped to upstream network/troubleshooting.md
---

# CalicoVPP Network Troubleshooting Workflow

This runbook maps directly to the upstream guide:
- https://github.com/projectcalico/vpp-dataplane/blob/master/docs/network/troubleshooting.md

It is tailored for this MCP server and for CalicoVPP where VPP runs in Kubernetes pods.

## Scope

Use this workflow when there are suspected routing or forwarding issues — missing routes, incorrect next-hops, unresolved FIB entries, or pod-to-pod connectivity failures.

## Step 1: Discover target pods

Call `cluster` with `command=get_pods` to list all CalicoVPP pods across nodes.

Pick the affected node(s) first, then expand to all nodes if the issue appears cluster-wide.

## Step 2: List available VRFs

Call `vppctl` with:
- `command=show ip table`, `pod_name=<pod-name>`
- `command=show ip6 table`, `pod_name=<pod-name>`

Validate the expected VRF structure:
- `table_id=0` → default VRF (uplink and tunnels to other nodes)
- `punt-table-ip4/6` → PUNT table (traffic to host or VCL pods)
- `calico-pods-ip4/6` → indirection table preventing asymmetric pod→nodeIP traffic
- Per-pod VRFs (one per pod) and their RPF (uRPF check) counterparts

## Step 3: Locate a pod VRF by IP address

If a specific pod IP is under investigation, find its fib-index using:

Call `vppctl` with:
- `command=show int addr`, `pod_name=<pod-name>`

Then grep the output for the target IP to get its `fib-idx`. This index is used in Step 4 to query that pod's VRF directly.

## Step 4: Check pod VRF routes

Call `vppctl` with:
- `command=show ip fib index <fib-idx>`, `pod_name=<pod-name>`
- `command=show ip6 fib index <fib-idx>`, `pod_name=<pod-name>`

Validate the expected routes in a pod VRF:
- `0.0.0.0/0` → default lookup in `calico-pods-ip4` (not a drop)
- `dpo-drop` entries for `0.0.0.0/32`, multicast, and broadcast ranges are **expected**

## Step 5: Check the default VRF (index 0)

Call `vppctl` with:
- `command=show ip fib index 0`, `pod_name=<pod-name>`
- `command=show ip6 fib index 0`, `pod_name=<pod-name>`

Expected route categories in VRF 0:
- `0.0.0.0/0` → default route via uplink (e.g. `host-eth0`)
- Local pod routes (e.g. `11.0.0.1/32`) → via `tun<N>` (local tuntap)
- Remote node pod CIDRs (e.g. `11.0.0.64/26`) → via `ipip<N>` with stacked node IP
- Service VIPs → via `cnat-client` DPO
- Remote node IPs → via `host-eth0` adjacency
- Glean routes for attached prefixes
- Local addresses → `dpo-receive`
- Expected drop routes for multicast/broadcast

Flags to watch for:
- `UNRESOLVED` next-hops → routing misconfiguration
- Routes missing that should exist → BGP not advertising/receiving correctly
- `dpo-drop` where a forwarding path is expected → broken route

## Step 6: Targeted prefix lookup

To drill into a specific prefix or IP in a given VRF:

Call `vppctl` with:
- `command=show ip fib index 0 <prefix-or-ip>`, `pod_name=<pod-name>`

This shows the full path-list, adjacency details, and forwarding chain for that entry.

## Step 7: Correlate with BGP if routes are missing

If expected routes are absent from the FIB:

1. Call `gobgp` with `command=global rib -a ipv4`, `pod_name=<pod-name>` — check if BGP has the route
2. Call `gobgp` with `command=neighbor`, `pod_name=<pod-name>` — check if peers are `Establ`
3. If BGP has the route but FIB does not, suspect a programming failure in the dataplane

## Fast Triage Order

1. `cluster` → `get_pods` (discover all pods)
2. `vppctl` → `show ip table` + `show ip6 table` (confirm VRF layout)
3. `vppctl` → `show int addr` (find fib-idx for affected pod IP)
4. `vppctl` → `show ip fib index 0` + `show ip6 fib index 0` (default VRF check)
5. `vppctl` → `show ip fib index <pod-fib-idx>` (pod VRF check)
6. `vppctl` → `show ip fib index 0 <specific-prefix>` (targeted lookup)
7. `gobgp` → `global rib -a ipv4/ipv6` (if routes missing from FIB)
