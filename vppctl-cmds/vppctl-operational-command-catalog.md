# VPP `vppctl` Operational Command Catalog

## Purpose
- This catalog groups all operational commands that do not fall under the `show`, `set`, `clear`, `create`, `delete`, or `test` verb prefixes.
- These commands perform protocol configuration, runtime actions, traffic generation, diagnostics, and lifecycle management.

## Category Index
| Category | Command Count |
|---|---:|
| CLI Infrastructure & Runtime Control | 15 |
| Event Logger | 6 |
| Packet Generator & Capture | 11 |
| Trace & Pcap | 5 |
| Performance Monitoring & Prometheus | 4 |
| Interface Lifecycle & Device Management | 18 |
| IP & IPv6 Routing, Punt & Forwarding | 18 |
| MPLS | 3 |
| BIER | 2 |
| Segment Routing (SRv6 & SR-MPLS) | 10 |
| L2 FIB, Classification & Bridge Domain | 10 |
| ABF, ADL & Network Policy (npol) | 15 |
| BFD (Bidirectional Forwarding Detection) | 10 |
| NAT44, NAT44-EI, NAT64, NAT66 & DET44 | 42 |
| CNAT | 3 |
| Policer | 7 |
| IPsec & IKEv2 | 12 |
| Snort & Intrusion Detection | 5 |
| LISP & ONE Overlay | 34 |
| GPE (Generic Protocol Extension) | 5 |
| MAP (Mapping of Address and Port) | 11 |
| DHCP & DNS | 5 |
| IGMP | 4 |
| QoS, Flow & Classification | 7 |
| SVS, STN & L3XC | 4 |
| Load Balancer | 6 |
| Transport, Session & Application | 21 |
| Tunnels: PVTI, WireGuard, VRRP & Pipe | 11 |
| Network Simulator (nsim) | 2 |
| sFlow, Soft-RSS & Misc | 11 |
| Service Chaining — SASC & SFDP | 13 |

## Commands by Category

### CLI Infrastructure & Runtime Control
Session management, command execution, macro definition, and process lifecycle.

| Command | Description |
|---|---|
| `api trace` | Enable, disable, or dump binary API message tracing for debugging API interactions. |
| `binary-api` | Send a raw binary API message from the CLI for low-level API testing. |
| `define` | Define a CLI macro (alias) that expands to a sequence of commands. |
| `echo` | Print a literal string to the CLI output (useful in scripts and exec files). |
| `exec` | Execute VPP CLI commands from an external file, one command per line. |
| `exit` | Close the current CLI session (same as `quit`). |
| `history` | Display the command history for the current CLI session. |
| `memory-trace` | Enable or disable per-allocation memory tracing for leak detection in a heap. |
| `save memory-trace` | Dump the current memory-trace allocation log to a file for offline analysis. |
| `q` | Alias for `quit`; close the current CLI session. |
| `quit` | Close the current CLI session and release its resources. |
| `restart` | Restart the VPP process (re-exec the binary with the same arguments). |
| `suspend` | Suspend the current CLI process context and yield to the main event loop. |
| `undefine` | Remove a previously defined CLI macro (alias). |
| `wait` | Pause CLI command processing for a specified number of seconds (useful in scripts). |

### Event Logger
Control the circular event-logger used for lightweight in-memory tracing.

| Command | Description |
|---|---|
| `event-logger clear` | Discard all events in the circular event-log buffer. |
| `event-logger resize` | Change the number of entries in the circular event-log buffer. |
| `event-logger restart` | Re-enable event logging after it was stopped. |
| `event-logger save` | Write the current event-log buffer contents to a file for offline analysis. |
| `event-logger stop` | Stop recording events into the event-log buffer. |
| `event-logger trace` | Configure which event classes are captured by the event logger. |

### Packet Generator & Capture
Generate synthetic traffic, control streams, and capture packets.

| Command | Description |
|---|---|
| `packet-generator capture` | Enable or disable pcap capture on a packet-generator interface. |
| `packet-generator configure` | Modify parameters of an existing packet-generator stream (rate, size, count). |
| `packet-generator delete` | Remove a named packet-generator stream definition. |
| `packet-generator disable-stream` | Stop a named packet-generator stream from producing packets. |
| `packet-generator enable-stream` | Start (or restart) a named packet-generator stream. |
| `packet-generator mac-filter` | Restrict packet-generator capture to frames matching specified MAC addresses. |
| `packet-generator new` | Define a new packet-generator stream with packet template, rate, count, and interface. |
| `pcap dispatch trace` | Enable per-node pcap capture at a specific point in the dispatch graph. |
| `pcap trace` | Start or stop pcap packet capture on an interface with an optional filter and max count. |
| `sasc pcap start` | Start pcap capture on SASC-processed traffic for service-chain debugging. |
| `sasc pcap stop` | Stop an active SASC pcap capture session. |

### Trace & Pcap
Add trace entries, configure filters, and control frame-queue tracing.

| Command | Description |
|---|---|
| `trace add` | Add packet-trace entries for a specific graph node (capture next N packets through that node). |
| `trace filter` | Set a classify-based filter so only matching packets are retained in the trace buffer. |
| `trace frame-queue` | Enable or disable tracing of inter-worker frame-queue handoff events. |
| `tracenode feature` | Enable the tracenode feature-arc on an interface to trace all packets entering a feature arc. |
| `selog emit-elog` | Emit a manually-triggered structured event-log entry for correlation with other telemetry. |

### Performance Monitoring & Prometheus
Control hardware performance counters and Prometheus metric export.

| Command | Description |
|---|---|
| `perfmon reset` | Reset all active performance-monitoring counters and discard collected statistics. |
| `perfmon start` | Start collecting hardware performance counters using the configured perfmon bundle. |
| `perfmon stop` | Stop collecting hardware performance counters. |
| `prom` | Enable or configure the Prometheus metrics exporter endpoint. |
| `prom patterns` | Set glob patterns to filter which VPP stats are exposed via the Prometheus exporter. |

### Interface Lifecycle & Device Management
Enable/disable interfaces, manage device drivers, BVI, loopback, pipe, and RSS.

| Command | Description |
|---|---|
| `adjacency counters` | Enable or disable per-adjacency packet/byte counters (may increase memory usage). |
| `buffer metadata tracking` | Enable or disable per-buffer metadata change tracking via the mdata plugin. |
| `bvi create` | Create a Bridge Virtual Interface (BVI) for routing between bridge domains. |
| `bvi delete` | Delete a previously created BVI interface. |
| `device attach` | Attach a physical or virtual device to the VPP device-driver framework. |
| `device create-interface` | Create a primary network interface on an attached device. |
| `device create-secondary-interface` | Create a secondary (additional queue/function) interface on an attached device. |
| `device detach` | Detach a device from the VPP device-driver framework and release its resources. |
| `device remove-interface` | Remove a primary interface from an attached device. |
| `device remove-secondary-interface` | Remove a secondary interface from an attached device. |
| `device reset` | Reset a device to its initial state (re-initialize queues and firmware). |
| `device set-rss-key` | Set the RSS hash key on a device for receive-side traffic distribution. |
| `disable ip6 interface` | Disable IPv6 processing (link-local address, ND) on an interface. |
| `enable ip4 interface` | Enable IPv4 processing on an interface (if not enabled by default). |
| `enable ip6 interface` | Enable IPv6 processing and generate a link-local address on an interface. |
| `interface collect detailed-stats` | Enable collection of detailed per-interface statistics (extended counters). |
| `loopback create-interface` | Create a software loopback interface (alternative syntax to `create loopback interface`). |
| `loopback delete-interface` | Delete a software loopback interface (alternative syntax to `delete loopback interface`). |
| `monitor interface` | Start real-time counter display for an interface (updates periodically in the CLI). |
| `p2p_ethernet` | Create or delete a point-to-point Ethernet sub-interface keyed by a remote MAC address. |
| `pipe create` | Create a software pipe interface pair (two connected interfaces for internal plumbing). |
| `pipe delete` | Delete a previously created software pipe interface pair. |
| `renumber interface` | Change the sw_if_index of an interface (useful for deterministic interface numbering). |
| `soft-rss clear` | Clear software RSS statistics on an interface. |
| `soft-rss config` | Configure software RSS hash parameters on an interface. |
| `soft-rss disable` | Disable software-based RSS on an interface. |
| `soft-rss enable` | Enable software-based RSS for multi-queue packet distribution on an interface. |

### IP & IPv6 Routing, Punt & Forwarding
Add/delete routes, neighbors, punt redirects, and forwarding policy entries.

| Command | Description |
|---|---|
| `ip container` | Add or remove an IP container-proxy entry binding an IP address to an interface. |
| `ip mroute` | Add or delete an IPv4 multicast route entry in a VRF table. |
| `ip neighbor` | Add, delete, or flush an IPv4/IPv6 neighbor (ARP/ND) entry on an interface. |
| `ip punt policer` | Apply a policer to punted IPv4 packets to rate-limit traffic sent to the host stack. |
| `ip punt redirect` | Redirect punted IPv4 packets to a specified next-hop and interface instead of the host. |
| `ip route` | Add, delete, or modify an IPv4 unicast route in a VRF table (next-hop, weight, preference). |
| `ip session redirect` | Add or delete a session-redirect rule that steers matching flows to an alternate path. |
| `ip syn filter` | Enable or disable TCP SYN flood filtering on IPv4 punt path. |
| `ip table` | Create or delete an IPv4 VRF / FIB table by table-id. |
| `ip virtual` | Set a virtual (floating) IPv4 address for ARP/punt responses (e.g., gateway VIP). |
| `ip6 nd` | Configure IPv6 Neighbor Discovery parameters (RA interval, prefix advertisement) on an interface. |
| `ip6 nd address autoconfig` | Enable or disable SLAAC (Stateless Address Autoconfiguration) on an interface. |
| `ip6 punt policer` | Apply a policer to punted IPv6 packets to rate-limit traffic sent to the host stack. |
| `ip6 punt redirect` | Redirect punted IPv6 packets to a specified next-hop and interface. |
| `ip6 table` | Create or delete an IPv6 VRF / FIB table by table-id. |
| `punt socket register` | Register a Unix domain socket to receive punted packets of a specified type. |
| `punt socket deregister` | Remove a previously registered punt socket binding. |
| `arping` | Send a gratuitous or solicited ARP/NDP request from an interface for address resolution testing. |

### MPLS
Configure MPLS labels, tables, and tunnels.

| Command | Description |
|---|---|
| `mpls local-label` | Bind a local MPLS label to an IP prefix for label imposition or disposition. |
| `mpls table` | Create or delete an MPLS label table by table-id. |
| `mpls tunnel` | Create, modify, or delete an MPLS tunnel (label stack, next-hop, and path). |

### BIER
Configure Bit Index Explicit Replication tables and routes.

| Command | Description |
|---|---|
| `bier route` | Add or delete a BIER forwarding entry (bit-position to next-hop mapping) in a BIER table. |
| `bier table` | Create or delete a BIER table (set, sub-domain, bit-string-length). |

### Segment Routing (SRv6 & SR-MPLS)
Configure SRv6 local-SIDs, policies, steering rules, and SR-MPLS policies.

| Command | Description |
|---|---|
| `sr localsid` | Create or delete an SRv6 local-SID with a specified behaviour (End, End.X, End.DX4, etc.). |
| `sr policy` | Create, modify, or delete an SRv6 policy (BSID, segment lists, encap/insert mode). |
| `sr steer` | Create or delete an SRv6 traffic-steering rule that directs matching prefixes into an SR policy. |
| `sr mpls policy` | Create or delete an SR-MPLS policy (BSID, label stack segments). |
| `sr mpls policy te` | Create or delete an SR-MPLS Traffic-Engineering policy with explicit path constraints. |
| `sr mpls steer` | Create or delete an SR-MPLS steering rule to direct traffic into an SR-MPLS policy. |
| `sr pt add iface` | Add an interface to SRv6 Path Tracing for per-hop latency measurement. |
| `sr pt del iface` | Remove an interface from SRv6 Path Tracing. |
| `sr pt show iface` | Display SRv6 Path Tracing configuration for an interface. |

### L2 FIB, Classification & Bridge Domain
Manage L2 forwarding entries, classification tables, and MAC-time access control.

| Command | Description |
|---|---|
| `bond add` | Add a member (slave) interface to an existing bonding group. |
| `bond del` | Remove a member (slave) interface from a bonding group. |
| `cdp` | Enable or disable Cisco Discovery Protocol processing globally. |
| `classify filter` | Attach a classify table as a pcap or trace capture filter. |
| `classify session` | Add or delete a session (match rule + action) in a classify table. |
| `classify table` | Create or delete a classify table with mask, match-n-vectors, and memory-size parameters. |
| `l2 rewrite entry` | Add or delete an L2 output-rewrite entry (destination MAC replacement) for an interface. |
| `l2fib add` | Add a static entry (MAC + bridge-domain → interface) to the L2 FIB. |
| `l2fib del` | Delete a specific entry from the L2 FIB by MAC address and bridge-domain. |
| `l2fib flush-mac all` | Flush all dynamically-learned MAC entries from the L2 FIB across all bridge domains. |
| `l2fib flush-mac bridge-domain` | Flush all dynamically-learned MAC entries in a specific bridge domain. |
| `l2fib flush-mac interface` | Flush all dynamically-learned MAC entries associated with a specific interface. |
| `mactime enable-disable` | Enable or disable MAC-time-based access control (time-of-day allow/deny by MAC) on an interface. |
| `debug lacp` | Enable or disable LACP protocol debug logging on a bonded interface. |

### ABF, ADL & Network Policy (npol)
Configure ACL-Based Forwarding, Allow/Deny Lists, and network-policy rule sets.

| Command | Description |
|---|---|
| `abf attach` | Attach an ABF (ACL-Based Forwarding) policy to an interface for policy-based routing. |
| `abf policy` | Create or delete an ABF policy binding an ACL index to a set of forwarding paths. |
| `adl allowlist` | Add or remove entries in the ADL (Allow/Deny List) allowlist for an interface. |
| `adl interface` | Enable or disable ADL (Allow/Deny List) processing on an interface. |
| `auto-sdl` | Enable or disable automatic Session Descriptor List generation for DDoS mitigation. |
| `npol interface clear` | Remove all network-policy bindings from an interface. |
| `npol interface configure` | Attach a network-policy to an interface's input or output path. |
| `npol ipset add` | Create a named IP-set for use in network-policy match rules. |
| `npol ipset add member` | Add an IP prefix or address to an existing named IP-set. |
| `npol ipset del` | Delete a named IP-set and all its members. |
| `npol ipset del member` | Remove an IP prefix or address from a named IP-set. |
| `npol match` | Test a packet tuple against configured network-policy rules and display the result. |
| `npol policy add` | Create a network-policy binding a set of rules to a named policy. |
| `npol policy del` | Delete a named network-policy. |
| `npol rule add` | Add a match rule (IP-set, protocol, port, action) to a network-policy. |
| `npol rule del` | Delete a rule from a network-policy by index. |

### BFD (Bidirectional Forwarding Detection)
Manage BFD sessions, authentication keys, and echo-source for fast failure detection.

| Command | Description |
|---|---|
| `bfd key del` | Delete a BFD authentication key by key-id. |
| `bfd key set` | Create or update a BFD authentication key (SHA1 or meticulous-SHA1). |
| `bfd udp echo-source del` | Remove the interface used as the BFD echo source. |
| `bfd udp echo-source set` | Set the interface whose address is used as the BFD echo-function source. |
| `bfd udp session add` | Create a new BFD-over-UDP session to a remote peer for liveliness detection. |
| `bfd udp session auth activate` | Activate authentication on an existing BFD session using a configured key. |
| `bfd udp session auth deactivate` | Deactivate authentication on a BFD session (revert to unauthenticated). |
| `bfd udp session del` | Delete a BFD-over-UDP session. |
| `bfd udp session mod` | Modify the desired-min-tx or required-min-rx intervals of an existing BFD session. |
| `bfd udp session set-flags` | Set administrative flags (e.g., admin-down) on a BFD session. |

### NAT44, NAT44-EI, NAT64, NAT66 & DET44
Add/remove NAT pools, static mappings, identity mappings, forwarding toggles, HA, and deterministic NAT.

| Command | Description |
|---|---|
| `configure policer` | Create or update a hierarchical policer (single/dual-rate, color-aware) by name. |
| `det44 add` | Add a DET44 (deterministic NAT) inside-to-outside address mapping range. |
| `det44 close session in` | Close a DET44 session by specifying the inside (pre-NAT) 5-tuple. |
| `det44 close session out` | Close a DET44 session by specifying the outside (post-NAT) 5-tuple. |
| `det44 forward` | Display the outside address and port range for a given DET44 inside address. |
| `det44 plugin` | Enable or disable the DET44 deterministic-NAT plugin. |
| `det44 reverse` | Display the inside address for a given DET44 outside address and port. |
| `dslite add pool address` | Add an IPv4 address to the DS-Lite AFTR NAT pool. |
| `dslite set aftr-tunnel-endpoint-address` | Set the IPv6 address of the DS-Lite AFTR tunnel endpoint. |
| `dslite set b4-tunnel-endpoint-address` | Set the IPv6 address of the DS-Lite B4 tunnel endpoint. |
| `nat ipfix logging` | Enable or disable IPFIX/Netflow logging of NAT translation events. |
| `nat mss-clamping` | Set TCP MSS clamping value applied to packets traversing NAT. |
| `nat set logging level` | Set the syslog severity level for NAT event messages. |
| `nat44 add address` | Add one or more IPv4 addresses to the NAT44-ED external address pool. |
| `nat44 add identity mapping` | Create a NAT44-ED identity (no-translate) mapping for specific traffic. |
| `nat44 add interface address` | Use an interface's dynamic address as a NAT44-ED external pool address. |
| `nat44 add load-balancing back-end` | Add a back-end server (IP:port + weight) to a NAT44-ED load-balanced static mapping. |
| `nat44 add load-balancing static mapping` | Create a NAT44-ED static mapping with load-balanced back-end servers. |
| `nat44 add static mapping` | Create a NAT44-ED 1:1 or port-restricted static mapping (inside ↔ outside). |
| `nat44 del session` | Delete a specific active NAT44-ED session by its inside or outside 5-tuple. |
| `nat44 ei add address` | Add one or more IPv4 addresses to the NAT44-EI external address pool. |
| `nat44 ei add identity mapping` | Create a NAT44-EI identity (no-translate) mapping for specific traffic. |
| `nat44 ei add interface address` | Use an interface's dynamic address as a NAT44-EI external pool address. |
| `nat44 ei add static mapping` | Create a NAT44-EI 1:1 or port-restricted static mapping (inside ↔ outside). |
| `nat44 ei addr-port-assignment-alg` | Select the NAT44-EI address/port assignment algorithm (default, MAP-E, port-range). |
| `nat44 ei del session` | Delete a specific active NAT44-EI session by its 5-tuple. |
| `nat44 ei del user` | Delete all NAT44-EI sessions belonging to a specific inside user address. |
| `nat44 ei forwarding` | Enable or disable NAT44-EI forwarding (pass-through) of non-translated packets. |
| `nat44 ei ha failover` | Configure NAT44-EI high-availability failover parameters (peer address, port). |
| `nat44 ei ha flush` | Flush all NAT44-EI HA session-sync state. |
| `nat44 ei ha listener` | Configure NAT44-EI HA listener to receive session-sync updates from a peer. |
| `nat44 ei ha resync` | Trigger a full NAT44-EI HA session resynchronization with the peer. |
| `nat44 ei ipfix logging` | Enable or disable IPFIX logging for NAT44-EI translation events. |
| `nat44 ei mss-clamping` | Set TCP MSS clamping value for NAT44-EI translated connections. |
| `nat44 ei plugin` | Enable or disable the NAT44-EI (endpoint-independent) plugin. |
| `nat44 ei set logging level` | Set the syslog severity level for NAT44-EI event messages. |
| `nat44 forwarding` | Enable or disable NAT44-ED forwarding of non-translated packets. |
| `nat44 plugin` | Enable or disable the NAT44-ED (endpoint-dependent) plugin with pool/mode parameters. |
| `nat44 vrf route` | Add or remove an inter-VRF NAT44-ED routing entry. |
| `nat44 vrf table` | Create or delete a NAT44-ED VRF routing table for multi-tenant deployments. |
| `nat64 add interface address` | Use an interface's address as a NAT64 external pool address. |
| `nat64 add pool address` | Add an IPv4 address or range to the NAT64 address pool. |
| `nat64 add prefix` | Add a NAT64 well-known or network-specific IPv6 prefix (e.g., 64:ff9b::/96). |
| `nat64 add static bib` | Add a static NAT64 Binding Information Base entry (IPv6 ↔ IPv4 address + port). |
| `nat64 plugin` | Enable or disable the NAT64 stateful translation plugin. |
| `nat66 add static mapping` | Add a static NAT66 1:1 IPv6 address mapping. |
| `nat66 plugin` | Enable or disable the NAT66 stateful IPv6-to-IPv6 translation plugin. |

### CNAT
Configure Cloud-NAT translation rules and clients.

| Command | Description |
|---|---|
| `cnat client add` | Register a CNAT client (source endpoint) for translation tracking. |
| `cnat log` | Enable or disable CNAT session logging to syslog or IPFIX. |
| `cnat translation` | Create, modify, or delete a CNAT translation rule (VIP → backends with health/weights). |

### Policer
Create, bind, and apply traffic policers.

| Command | Description |
|---|---|
| `policer add` | Create a named policer with CIR/PIR rates, burst sizes, and conform/exceed/violate actions. |
| `policer bind` | Bind a named policer to a worker thread for thread-local rate enforcement. |
| `policer del` | Delete a named policer and release its token-bucket state. |
| `policer input` | Apply a named policer to an interface's input path. |
| `policer output` | Apply a named policer to an interface's output path. |
| `policer reset` | Reset a policer's token buckets and counters to their initial state. |

### IPsec & IKEv2
Create/delete IPsec SAs, SPDs, tunnel-protect bindings, interfaces, and IKEv2 profiles.

| Command | Description |
|---|---|
| `ikev2 dpd disable` | Disable IKEv2 Dead Peer Detection for a profile. |
| `ikev2 initiate` | Manually initiate an IKEv2 SA negotiation for a configured profile. |
| `ikev2 profile` | Create, modify, or delete an IKEv2 profile (authentication, lifetime, proposals). |
| `ikev2 set liveness` | Set the IKEv2 liveness-check interval and threshold for peer reachability. |
| `ikev2 set sleep interval` | Set the IKEv2 background process sleep interval between periodic tasks. |
| `ipsec itf create` | Create an IPsec tunnel interface for route-based VPN. |
| `ipsec itf delete` | Delete an IPsec tunnel interface. |
| `ipsec policy` | Add or delete an IPsec Security Policy (inbound/outbound) in an SPD. |
| `ipsec sa` | Add or delete an IPsec Security Association (algorithm, key, SPI, tunnel endpoints). |
| `ipsec sa bind` | Bind an IPsec SA to a specific worker thread for deterministic processing. |
| `ipsec spd` | Create or delete an IPsec Security Policy Database. |
| `ipsec tunnel protect` | Apply IPsec tunnel-protect (SA binding) to an interface or tunnel. |

### Snort & Intrusion Detection
Attach/detach Snort IDS/IPS instances to interfaces.

| Command | Description |
|---|---|
| `snort attach` | Attach a Snort instance to an interface for inline IDS/IPS inspection. |
| `snort create-instance` | Create a named Snort instance with queue-size and drop-on-disconnect parameters. |
| `snort delete instance` | Delete a named Snort instance and release its shared-memory queues. |
| `snort detach` | Detach a Snort instance from an interface. |
| `snort disconnect client` | Force-disconnect a Snort client from its shared-memory queue. |

### LISP & ONE Overlay
Configure LISP and ONE (Overlay Network Engine) EID tables, locators, map-servers, and overlays.

| Command | Description |
|---|---|
| `lisp adjacency` | Add or delete a LISP adjacency (EID-pair → RLOC forwarding entry). |
| `lisp disable` | Disable the LISP control plane. |
| `lisp eid-table` | Add or delete a LISP EID-table mapping (EID-prefix → locator-set). |
| `lisp eid-table map` | Map a LISP EID-table to an L2 bridge-domain or L3 VRF. |
| `lisp enable` | Enable the LISP control plane. |
| `lisp locator` | Add or delete a locator (RLOC interface + priority + weight) in a locator-set. |
| `lisp locator-set` | Create or delete a named LISP locator-set. |
| `lisp map-register` | Enable or disable LISP map-register message sending. |
| `lisp map-request itr-rlocs` | Set the locator-set used as source RLOCs in outgoing LISP map-request messages. |
| `lisp map-request mode` | Set the LISP map-request mode (source-destination or destination-only). |
| `lisp map-resolver` | Add or delete a LISP map-resolver address. |
| `lisp map-server` | Add or delete a LISP map-server address. |
| `lisp pitr` | Enable or disable LISP Proxy Ingress Tunnel Router using a specified locator-set. |
| `lisp remote-mapping` | Add or delete a remote EID-to-RLOC mapping in the LISP map-cache. |
| `lisp rloc-probe` | Enable or disable LISP RLOC probing for locator reachability verification. |
| `lisp use-petr` | Enable or disable use of a Proxy ETR (PETR) for encapsulating to non-LISP destinations. |
| `one adjacency` | Add or delete a ONE overlay adjacency. |
| `one disable` | Disable the ONE control plane. |
| `one eid-table` | Add or delete a ONE EID-table mapping. |
| `one eid-table map` | Map a ONE EID-table to an L2 bridge-domain or L3 VRF. |
| `one enable` | Enable the ONE control plane. |
| `one l2 arp` | Add or delete a static L2 ARP entry in the ONE overlay. |
| `one locator` | Add or delete a locator in a ONE locator-set. |
| `one locator-set` | Create or delete a named ONE locator-set. |
| `one map-register` | Enable or disable ONE map-register message sending. |
| `one map-register fallback-threshold` | Set the ONE map-register fallback threshold for map-server failover. |
| `one map-register ttl` | Set the TTL value used in ONE map-register messages. |
| `one map-request itr-rlocs` | Set the locator-set for ONE map-request source RLOCs. |
| `one map-request mode` | Set the ONE map-request mode. |
| `one map-resolver` | Add or delete a ONE map-resolver address. |
| `one map-server` | Add or delete a ONE map-server address. |
| `one ndp` | Add or delete a static NDP entry in the ONE overlay. |
| `one nsh-mapping` | Add or delete a ONE NSH-to-locator mapping. |
| `one petr mode` | Enable or disable ONE Proxy ETR mode. |
| `one pitr` | Enable or disable ONE Proxy Ingress Tunnel Router. |
| `one pitr mode` | Set ONE PITR mode parameters. |
| `one remote-mapping` | Add or delete a remote mapping in the ONE map-cache. |
| `one rloc-probe` | Enable or disable ONE RLOC probing. |
| `one statistics` | Enable or disable ONE per-EID traffic statistics collection. |
| `one statistics flush` | Flush all collected ONE traffic statistics. |
| `one use-petr` | Enable or disable ONE Proxy ETR usage. |
| `one xtr mode` | Enable or disable ONE xTR (Ingress/Egress Tunnel Router) mode. |

### GPE (Generic Protocol Extension)
Configure LISP-GPE encapsulation, entries, interfaces, and native forwarding.

| Command | Description |
|---|---|
| `gpe` | Enable or disable GPE (Generic Protocol Extension) datapath. |
| `gpe encap` | Set the GPE encapsulation mode (LISP, VXLAN, or NSH). |
| `gpe entry` | Add or delete a GPE forwarding entry (VNI + EID → tunnel). |
| `gpe iface` | Create or delete a GPE interface for a given VNI. |
| `gpe native-forward` | Add or delete a native-forward (non-encapsulated) entry in GPE. |

### MAP (Mapping of Address and Port)
Configure MAP-E/MAP-T domains, rules, and translation parameters.

| Command | Description |
|---|---|
| `map add domain` | Create a MAP-E or MAP-T domain with IPv4/IPv6 prefixes, EA-bits, and PSID parameters. |
| `map add rule` | Add a forwarding rule (user IPv6 prefix → shared IPv4 address) to a MAP domain. |
| `map del domain` | Delete a MAP domain by its domain index. |
| `map interface` | Enable or disable MAP-E/MAP-T processing on an interface. |
| `map params fragment` | Configure MAP fragmentation parameters (inner/outer, DF-bit handling). |
| `map params icmp source-address` | Set the source IPv4 address used in MAP-generated ICMP error messages. |
| `map params icmp6 unreachables` | Enable or disable generation of ICMPv6 unreachable messages for MAP errors. |
| `map params pre-resolve` | Configure MAP pre-resolved next-hop addresses for IPv4 and IPv6. |
| `map params security-check` | Enable or disable MAP security-check validation of source addresses. |
| `map params tcp-mss` | Set the TCP MSS clamping value for MAP-translated segments. |
| `map params traffic-class` | Set the IPv6 traffic-class / IPv4 DSCP handling for MAP encapsulation. |

### DHCP & DNS
Configure DHCP clients, DNS resolution, and name-server addresses.

| Command | Description |
|---|---|
| `dhcp6 client` | Start or stop a DHCPv6 (IA_NA) client on an interface for address acquisition. |
| `dhcp6 pd client` | Start or stop a DHCPv6 Prefix Delegation client on an interface. |
| `dns` | Enable or disable the VPP DNS name-resolution cache and resolver. |
| `dns cache` | Add or remove static entries in the DNS resolver cache. |
| `dns name-server` | Add or remove an upstream DNS name-server address for recursive resolution. |

### IGMP
Configure IGMP group membership, proxy devices, and listeners.

| Command | Description |
|---|---|
| `igmp` | Enable or disable IGMP protocol processing globally. |
| `igmp listen` | Add or remove a static IGMP group membership (join/leave) on an interface. |
| `igmp proxy-dev` | Create or delete an IGMP proxy device for upstream multicast signaling. |
| `igmp proxy-dev itf` | Add or remove a downstream interface from an IGMP proxy device. |

### QoS, Flow & Classification
Configure QoS maps/marking/recording, hardware flow rules, and IPFIX.

| Command | Description |
|---|---|
| `flow` | Add, delete, or enable/disable a hardware-offloaded flow classification rule. |
| `flowprobe feature add-del` | Enable or disable the flowprobe (IPFIX flow-export) feature on an interface. |
| `flowprobe params` | Set flowprobe parameters: active/passive timeout and record format (L2/L3/L4). |
| `ipfix classify table` | Bind a classify table to the IPFIX flow-export for per-flow telemetry. |
| `ipfix flush` | Force an immediate flush of all buffered IPFIX template and data records. |
| `qos egress map` | Create or delete a QoS egress map translating internal QoS values to header markings. |
| `qos mark` | Enable or disable QoS marking (DSCP/MPLS-EXP/VLAN-PCP write-back) on an interface. |
| `qos record` | Enable or disable QoS recording (read DSCP/MPLS-EXP/VLAN-PCP into internal QoS) on an interface. |
| `qos store` | Set a fixed QoS value to store on all packets entering an interface. |

### SVS, STN & L3XC
Configure Source VRF Select, Steal-The-NIC, and L3 cross-connect rules.

| Command | Description |
|---|---|
| `l3xc` | Create or delete an L3 cross-connect rule (steer all traffic from an interface to a next-hop path). |
| `svs enable` | Enable or disable Source VRF Select on an interface. |
| `svs route` | Add or delete a Source VRF Select route (source-prefix → VRF lookup). |
| `svs table` | Create or delete an SVS source-address routing table. |
| `stn rule` | Add or delete a Steal-The-NIC rule to punt specific destination traffic to the host. |

### Load Balancer
Configure the LB plugin: VIPs, application servers, and interface NAT bindings.

| Command | Description |
|---|---|
| `lb as` | Add or delete an application server (back-end) address for a load-balancer VIP. |
| `lb conf` | Set global load-balancer parameters (flow-timeout, buckets, GRE/encap options). |
| `lb flush vip` | Flush per-flow state for a VIP, forcing new flows to be rehashed to back-ends. |
| `lb set interface nat4` | Enable or disable LB-plugin NAT44 on an interface (for DSR return traffic). |
| `lb set interface nat6` | Enable or disable LB-plugin NAT66 on an interface (for DSR return traffic). |
| `lb vip` | Create or delete a load-balancer Virtual IP with protocol, port, and encap type. |

### Transport, Session & Application
Configure session rules, TCP parameters, HTTP servers/clients, QUIC, TLS, UDP, and VRRP.

| Command | Description |
|---|---|
| `app crypto add tls-profile` | Register a named TLS profile (certificate, key, CA) for application use. |
| `app crypto del tls-profile` | Remove a named TLS profile from the certificate store. |
| `app evt-collector` | Configure an application event collector for session-layer telemetry. |
| `app ns` | Create or delete an application namespace with secret, FIB table, and scope. |
| `hsi` | Start the built-in HTTP server index (HSI) test application. |
| `http cli client` | Send an HTTP request to VPP's own CLI-over-HTTP endpoint for remote command execution. |
| `http cli server` | Enable the HTTP server that exposes VPP CLI commands via HTTP GET/POST. |
| `http client` | Launch a one-shot HTTP client request (GET/POST) from VPP to an external server. |
| `http connect proxy client enable` | Enable the HTTP CONNECT proxy client for tunneling TCP over HTTP. |
| `http connect proxy client listener` | Configure the HTTP CONNECT proxy client listener (bind address, port). |
| `http static listener` | Configure the HTTP static file-server listener (bind address, port, root directory). |
| `http static server` | Enable the HTTP static file server plugin with cache-size and fifo parameters. |
| `http tps` | Start the HTTP transactions-per-second (TPS) benchmark application. |
| `ping` | Send ICMP echo requests from VPP to a destination address for reachability testing. |
| `session` | Enable or disable the session layer (host-stack) and set global session parameters. |
| `session replay fifo` | Replay a previously captured session FIFO trace for debugging. |
| `session rule` | Add or delete a session-layer rule (permit/deny/redirect by 5-tuple and app namespace). |
| `session sdl` | Add or delete a Session Descriptor List entry for traffic classification. |
| `tcp debug` | Set TCP debug verbosity level or toggle TCP-specific debug features. |
| `tcp replay scoreboard` | Replay a captured TCP SACK scoreboard trace for debugging retransmission logic. |
| `tcp src-address` | Set a preferred source address for outgoing TCP connections from VPP. |
| `tls openssl set` | Set OpenSSL engine parameters for TLS processing. |
| `tls openssl set-tls` | Set the default TLS version or cipher-suite list for the OpenSSL TLS engine. |
| `quic set crypto api` | Select the crypto API (vpp, openssl, picotls) used by the QUIC transport. |
| `quic set fifo-size` | Set the default FIFO size for QUIC transport sessions. |
| `udp decap` | Configure UDP decapsulation (strip outer UDP/IP headers) on a port for tunnel termination. |
| `udp encap` | Create or delete a UDP encapsulation entry (src/dst IP + port) for tunnel origination. |
| `udp-echo` | Start a built-in UDP echo client or server for basic connectivity testing. |

### Tunnels: PVTI, WireGuard, VRRP & Pipe
Manage PVTI tunnel interfaces, WireGuard tunnels/peers, VRRP virtual routers.

| Command | Description |
|---|---|
| `pvti interface create` | Create a Packet Vector Tunnel Interface (PVTI) for batched-packet tunnel transport. |
| `pvti interface delete` | Delete a PVTI tunnel interface. |
| `wireguard create` | Create a WireGuard tunnel interface with local key, port, and source address. |
| `wireguard delete` | Delete a WireGuard tunnel interface. |
| `wireguard peer add` | Add a WireGuard peer (public key, endpoint, allowed-IPs) to a tunnel interface. |
| `wireguard peer remove` | Remove a WireGuard peer from a tunnel interface. |
| `vrrp peers` | Configure VRRP peer addresses for unicast advertisement mode. |
| `vrrp proto` | Set VRRP protocol parameters (advertisement interval, priority, preemption). |
| `vrrp vr add` | Create a VRRP virtual router instance on an interface with VR-ID and virtual IPs. |
| `vrrp vr del` | Delete a VRRP virtual router instance. |
| `vrrp vr track-if` | Add or remove an interface-tracking entry that adjusts VRRP priority on link failure. |

### Network Simulator (nsim)
Enable the network delay/loss simulator on interfaces.

| Command | Description |
|---|---|
| `nsim cross-connect enable-disable` | Enable or disable nsim as a cross-connect between two interfaces (bidirectional delay). |
| `nsim output-feature enable-disable` | Enable or disable nsim as an output feature on an interface (unidirectional delay). |

### sFlow, Soft-RSS & Misc
Configure sFlow sampling, ILA, and miscellaneous operational commands.

| Command | Description |
|---|---|
| `ila entry` | Add or delete an Identifier-Locator Addressing (ILA) translation entry. |
| `ila interface` | Enable or disable ILA processing on an interface. |
| `kill sfdp session` | Forcefully terminate a specific SFDP-tracked session by session-id. |
| `sflow direction` | Set the sFlow sampling direction (ingress, egress, or both) on an interface. |
| `sflow drop-monitoring` | Enable or disable sFlow drop-reason monitoring for packet-loss visibility. |
| `sflow enable-disable` | Enable or disable sFlow sampling globally or on a specific interface. |
| `sflow header-bytes` | Set the number of packet header bytes captured in sFlow samples. |
| `sflow polling-interval` | Set the sFlow counter-polling interval (seconds) for interface statistics export. |
| `sflow sampling-rate` | Set the sFlow packet sampling rate (1-in-N packets). |
| `dump sasc session` | Dump detailed SASC session state for a specific session index or 5-tuple. |
| `dump sasc session ring` | Dump the SASC session ring buffer for debugging session lifecycle events. |

### Service Chaining — SASC & SFDP
Operational commands for SASC and SFDP service-chaining frameworks.

| Command | Description |
|---|---|
| `sfdp gateway geneve-input` | Configure the SFDP GENEVE input interface for gateway-mode decapsulation. |
| `sfdp nat alloc-pool` | Allocate an address pool for SFDP-integrated NAT translations. |
| `sfdp snort create-instance` | Create a Snort instance managed by the SFDP service-chaining framework. |
| `sfdp snort delete-instance` | Delete an SFDP-managed Snort instance. |
| `sfdp tenant` | Create, modify, or delete an SFDP tenant definition with service-chain bindings. |
