# VPP `vppctl test` Command Catalog

## Purpose
- This catalog groups all `test ...` commands available in the VPP CLI (`vppctl`).
- `test` commands exercise internal subsystems, simulate failure scenarios, and validate correctness without affecting production traffic.
- These are primarily intended for development, CI, and debugging use.

## Category Index
| Category | Command Count |
|---|---:|
| Infrastructure & Memory | 6 |
| Networking & Protocols | 7 |
| L2 & Forwarding | 4 |
| Transport & Application | 5 |
| NAT & Translation | 2 |
| Tunnels & Overlay | 4 |
| Logging & Telemetry | 2 |
| Service Chaining & Tenant Frameworks | 6 |
| HTTP | 1 |

## Commands by Category

### Infrastructure & Memory
Test DMA engines, DPDK buffers, frame-queue sizing, and hardware driver internals.

| Command | Description |
|---|---|
| `test dma` | Exercise the DMA subsystem by performing test copy operations and verifying completion. |
| `test dpdk buffer` | Validate DPDK buffer pool allocation and free-list consistency. |
| `test frame-queue nelts` | Change the number of elements in inter-worker frame queues for tuning or stress testing. |
| `test frame-queue threshold` | Adjust the frame-queue congestion threshold used for inter-worker handoff back-pressure. |
| `test rdma dump` | Dump internal state and descriptor rings of an RDMA (mlx5) interface for driver debugging. |
| `test vmxnet3` | Exercise vmxnet3 driver internals: descriptor rings, completion queues, and interrupt coalescing. |

### Networking & Protocols
Test IP checksum, IPv6 link/ND, DNS formatting, and IGMP timers.

| Command | Description |
|---|---|
| `test dns expire` | Force-expire one or all DNS cache entries to exercise the cache-eviction code path. |
| `test dns format` | Encode a DNS name into wire format and print the result for protocol debugging. |
| `test dns unformat` | Decode a wire-format DNS name back into dotted notation for validation. |
| `test igmp timers` | Override IGMP query/report timer values for accelerated protocol-state-machine testing. |
| `test ip checksum` | Compute and verify IPv4 header checksums using the VPP incremental-checksum routines. |
| `test ip6 connection-tracker` | Exercise the IPv6 connection-tracker (ct6) plugin by injecting synthetic state entries. |
| `test ip6 link` | Validate IPv6 link-local address generation and ND state-machine transitions on an interface. |

### L2 & Forwarding
Stress-test L2 FIB learning, L2 cross-connect patching, FIB walks, and load-balancer tables.

| Command | Description |
|---|---|
| `test fib-walk-process` | Trigger a manual FIB convergence walk across all FIB entries for dependency-resolution testing. |
| `test l2fib` | Insert, delete, or look up L2 FIB entries in bulk for hash-table stress testing. |
| `test l2patch` | Create or remove L2 patch (interface cross-connect) entries for testing direct L2 forwarding. |
| `test lb flowtable flush` | Flush the load-balancer per-flow consistent-hashing table to force flow reassignment. |

### Transport & Application
Test echo client/server, TLS client/server, and the TCP proxy server.

| Command | Description |
|---|---|
| `test echo clients` | Launch built-in echo clients that open sessions and transmit/receive data for throughput testing. |
| `test echo server` | Start a built-in echo server that reflects received data back to connected clients. |
| `test proxy server` | Start the built-in TCP proxy that forwards sessions between two application namespaces. |
| `test tls client` | Launch a built-in TLS client to test encrypted session establishment and data exchange. |
| `test tls server` | Start a built-in TLS server to validate certificate handling and encrypted data serving. |

### NAT & Translation
Test CNAT Maglev hashing and session scanner internals.

| Command | Description |
|---|---|
| `test cnat maglev` | Validate the Maglev consistent-hash backend-selection algorithm used by the CNAT plugin. |
| `test cnat scanner` | Trigger a manual run of the CNAT session-scanner to age out idle translation entries. |

### Tunnels & Overlay
Test L2TP counters, ONE/NSH overlay entries, and placeholder decap nodes.

| Command | Description |
|---|---|
| `test lt2p counters` | Validate L2TPv3 per-session counter increment and rollover behaviour. |
| `test one nsh` | Exercise the ONE (Overlay Network Engine) NSH encapsulation and decapsulation paths. |
| `test one nsh add-placeholder-decap-node` | Register a placeholder graph node for NSH decapsulation to test ONE service-function chaining. |
| `test-url-handler enable` | Enable a test HTTP URL handler for validating the HTTP static-server plugin dispatch logic. |

### Logging & Telemetry
Test log emission and syslog formatting.

| Command | Description |
|---|---|
| `test log` | Emit a synthetic log message at a specified severity level to validate the logging subsystem. |
| `test syslog` | Generate a test syslog (RFC 5424) message to verify exporter connectivity and formatting. |

### Service Chaining & Tenant Frameworks
Test SASC service-chain execution and SFDP session expiry.

| Command | Description |
|---|---|
| `test sasc list` | List available SASC test-case identifiers for selective test execution. |
| `test sasc run` | Execute a single SASC integration test case by name or index. |
| `test sasc run-all` | Execute all registered SASC integration test cases sequentially. |
| `test sasc run-id` | Execute a specific SASC test case identified by its numeric run-id. |
| `test sfdp expiry disable` | Disable the SFDP session-expiry timer to prevent automatic session cleanup during debugging. |
| `test sfdp expiry enable` | Re-enable the SFDP session-expiry timer after it was disabled for debugging. |
