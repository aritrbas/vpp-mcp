# VPP `vppctl clear` Command Catalog

## Purpose
- This catalog groups all `clear ...` commands available in the VPP CLI (`vppctl`).
- `clear` commands reset counters, flush caches, and remove transient state without altering persistent configuration.

## Category Index
| Category | Command Count |
|---|---:|
| Core Platform & Telemetry | 7 |
| Interfaces & L2 | 5 |
| IP Routing & FIB | 1 |
| Transport, Session & Application | 6 |
| Security, NAT & ACL | 5 |
| Tunnels, Encapsulation & Overlay | 5 |
| Service Chaining & Tenant Frameworks | 1 |

## Commands by Category

### Core Platform & Telemetry
Reset runtime counters, error tallies, trace buffers, and logging state.

| Command | Description |
|---|---|
| `clear api histogram` | Reset the API message processing latency histogram to zero. |
| `clear buffer traces` | Clear buffer-monitoring trace data collected by the bufmon plugin. |
| `clear errors` | Zero all per-node packet-processing error counters. |
| `clear logging` | Flush the in-memory log buffer and reset log sequence numbers. |
| `clear node counters` | Reset all per-node vector-rate and packet counters to zero. |
| `clear runtime` | Zero the per-node runtime statistics (clocks, vectors, calls, suspends). |
| `clear trace` | Discard all packets currently held in the packet trace buffer. |

### Interfaces & L2
Clear interface statistics, L2 FIB entries, and link-layer protocol counters.

| Command | Description |
|---|---|
| `clear hardware-interfaces` | Reset hardware-interface counters (rx/tx bytes, packets, drops, errors) on all or selected hardware interfaces. |
| `clear interface tag` | Remove the user-defined descriptive tag string from an interface. |
| `clear interfaces` | Reset all software-interface counters (rx/tx packets, bytes, drops, errors) to zero. |
| `clear l2fib` | Remove all dynamically-learned entries from the L2 forwarding information base. |
| `clear l2tp counters` | Zero the per-session packet and byte counters for all L2TPv3 tunnels. |

### IP Routing & FIB
Reset FIB walk statistics.

| Command | Description |
|---|---|
| `clear fib walk` | Reset the FIB convergence-walk statistics (walk counts, times, and histograms). |

### Transport, Session & Application
Clear session tables, transport statistics, and HTTP/proxy caches.

| Command | Description |
|---|---|
| `clear http connect proxy client stats` | Reset connection and request counters for the HTTP CONNECT proxy client. |
| `clear http static cache` | Invalidate all entries in the HTTP static file-server response cache. |
| `clear http stats` | Zero the HTTP layer statistics (requests, responses, errors). |
| `clear session` | Forcefully clean up one or all host-stack sessions and release associated FIFOs. |
| `clear session stats` | Reset the session-layer aggregate statistics (accepts, connects, resets). |
| `clear tcp stats` | Zero all TCP transport-layer counters (segments sent/received, retransmits, errors). |

### Security, NAT & ACL
Flush NAT sessions, IPsec counters, ACL connection-tracking state, and MAC-time entries.

| Command | Description |
|---|---|
| `clear acl-plugin sessions` | Purge all connection-tracking sessions maintained by the ACL stateful-firewall plugin. |
| `clear ipsec counters` | Zero per-SA and per-SPD packet and byte counters for IPsec. |
| `clear ipsec sa` | Remove a specific IPsec Security Association and release its crypto resources. |
| `clear mactime` | Reset the MAC-time allow/drop counters for all configured MAC addresses. |
| `clear nat44 ed sessions` | Flush all active NAT44 endpoint-dependent translation sessions. |
| `clear nat44 ei sessions` | Flush all active NAT44 endpoint-independent translation sessions. |

### Tunnels, Encapsulation & Overlay
Clear iOAM, PoT, and SRv6 overlay telemetry state.

| Command | Description |
|---|---|
| `clear igmp` | Reset IGMP protocol state and group membership records on all interfaces. |
| `clear ioam rewrite` | Remove the active iOAM hop-by-hop header rewrite configuration. |
| `clear ioam-trace profile` | Delete the currently configured iOAM trace profile (node-id, app-data, trace type). |
| `clear pot profile` | Remove all Proof-of-Transit profiles and reset the active profile index. |
| `clear sr localsid-counters` | Zero the packet/byte counters on all SRv6 local-SID entries. |
| `clear vxlan-gpe-ioam rewrite` | Remove the VXLAN-GPE iOAM header rewrite configuration. |

### Service Chaining & Tenant Frameworks
Clear SASC session state.

| Command | Description |
|---|---|
| `clear sasc sessions` | Flush all sessions tracked by the Service-Aware Service Chaining (SASC) framework. |
