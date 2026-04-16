---
name: pods-troubleshoot
description: CalicoVPP pod interface troubleshooting workflow mapped to upstream pods/troubleshooting.md
---

# CalicoVPP Pod Interface Troubleshooting Workflow

This runbook maps directly to the upstream guide:
- https://github.com/projectcalico/vpp-dataplane/blob/master/docs/pods/troubleshooting.md

It is tailored for this MCP server and for CalicoVPP where VPP runs in Kubernetes pods.

## Scope

Use this workflow when a specific pod has connectivity issues — tuntap misconfiguration, memif connection failures, VCL/hoststack problems, or incorrect interface queue placement.

## Step 1: Discover target pods

Call `cluster` with `command=get_pods` to list all CalicoVPP pods.

Identify which node the affected workload pod is running on, then target the corresponding CalicoVPP pod.

---

## Tuntap Interface Troubleshooting

### Step 2a: Find the tuntap interface by IP

Call `vppctl` with:
- `command=show int addr`, `pod_name=<pod-name>`

Search the output for the pod's IP address. This gives you the `tun<N>` interface name and its `fib-idx`.

### Step 3a: Inspect tuntap interface details

Call `vppctl` with:
- `command=show tun`, `pod_name=<pod-name>`

Grep for the interface name (e.g. `tun10`) to see:
- `host-ns`: network namespace path for the pod
- `host-mtu-size`: MTU configured in the pod
- `gso-enabled`, `csum-enabled`: offload settings
- Number of RX/TX virtqueues and queue sizes
- `avail.idx` vs `used.idx`: large gaps indicate queue stalls

Also check the hardware interface for queue placement and mode:

Call `vppctl` with:
- `command=show int rx-placement`, `pod_name=<pod-name>`

Validate:
- Each `tun<N>` is assigned to a worker thread
- Mode is `adaptive` (not stuck in `polling` if traffic is low, not `interrupt` under load)

---

## Memif Interface Troubleshooting

### Step 2b: Check pod-side memif socket (run on the workload pod)

The memif abstract socket should be present in the pod's network namespace. Verify using `lsof`, `ss`, or `netstat` inside the pod for `@vpp/memif-eth0`.

Expected states:
- `LISTEN` only → pod is waiting for VPP to connect
- `CONNECTED` → memif session is established

### Step 3b: Find the memif interface in VPP by IP

Call `vppctl` with:
- `command=show int addr`, `pod_name=<pod-name>`

Search for `memif` in the output. The memif shares the pod IP with the tuntap interface.

### Step 4b: Validate memif admin state

Call `vppctl` with:
- `command=show memif`, `pod_name=<pod-name>`

Validate:
- The socket entry for the pod's netns is present
- `flags admin-up` is set
- `conn-fd` is non-zero (fd=0 means not yet connected)

### Step 5b: Check memif hardware interface details

Call `vppctl` with:
- `command=show hardware-interface`, `pod_name=<pod-name>`

Grep for the memif interface name. Check:
- RX/TX queue assignment to worker threads
- `packets received` / `packets sent` counters are incrementing as expected

### Step 6b: Check PBL port-based load balancing (if memif uses port-based split)

If the pod annotation `cni.projectcalico.org/vppExtraMemifPorts` is set, validate the PBL classifier:

Call `vppctl` with:
- `command=show pbl client`, `pod_name=<pod-name>`

Validate:
- `matched dpo` points to the memif interface (not `dpo-drop`)
- `default dpo` points to the tuntap interface (not `dpo-drop`)
- A `dpo-drop` in either path means the corresponding interface is down or disconnected

### Step 7b: Check memif connection/disconnection events

Call `vppctl` with:
- `command=show log`, `pod_name=<pod-name>`

Look for `memif_plugin` lines with `disconnected` — these indicate memif session teardowns and can reveal flapping.

---

## VCL / Hoststack Troubleshooting

### Step 2c: Validate application namespace setup

Call `vppctl` with:
- `command=show app ns`, `pod_name=<pod-name>`

Each pod requesting VCL support should have an entry with:
- A `loop<N>` interface
- An abstract socket path matching the pod's network namespace (`abstract:vpp/session,netns_name=...`)

If the pod's namespace is absent, VCL is not configured or the agent failed to program it.

### Step 3c: Confirm application attachment

Call `vppctl` with:
- `command=show app`, `pod_name=<pod-name>`

The workload application should appear with its name and namespace. If absent, the application has not attached to VPP over the `@vpp/session` socket.

### Step 4c: Inspect active sessions

Call `vppctl` with:
- `command=show session verbose`, `pod_name=<pod-name>`

Key fields per session:
- `[thread:session-index][proto]` prefix
- State: `LISTEN`, `ESTABLISHED`, `OPENED`, `CLOSED`
- `Rx-f` / `Tx-f`: pending bytes in fifos — non-zero values indicate application is not consuming data

For a specific session, use:

Call `vppctl` with:
- `command=show session thread <T> index <I>`, `pod_name=<pod-name>`

This shows full session detail including fifo sizes, packet/byte counts, and transport flags.

---

## Common Interface Check (All Interface Types)

Call `vppctl` with:
- `command=show int rx-placement`, `pod_name=<pod-name>`

Validate that all pod interfaces (tun, memif, loop) are assigned to worker threads with appropriate modes (`adaptive` for tuntap, `polling` for memif).

## Fast Triage Order

**Tuntap issues:**
1. `cluster` → `get_pods`
2. `vppctl` → `show int addr` (find tun interface and fib-idx)
3. `vppctl` → `show tun` (queue and offload details)
4. `vppctl` → `show int rx-placement` (thread assignment and mode)

**Memif issues:**
1. `cluster` → `get_pods`
2. `vppctl` → `show int addr` (find memif interface)
3. `vppctl` → `show memif` (admin state and socket connectivity)
4. `vppctl` → `show pbl client` (port-based split validation, if applicable)
5. `vppctl` → `show log` (disconnection events)
6. `vppctl` → `show hardware-interface` (counter and queue details)

**VCL/Hoststack issues:**
1. `cluster` → `get_pods`
2. `vppctl` → `show app ns` (namespace registration)
3. `vppctl` → `show app` (application attachment)
4. `vppctl` → `show session verbose` (active sessions and fifo state)
