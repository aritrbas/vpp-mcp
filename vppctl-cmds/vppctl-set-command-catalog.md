# VPP `vppctl set` Command Catalog

## Purpose
- This catalog groups all `set ...` commands available in the VPP CLI (`vppctl`).
- `set` commands modify runtime configuration: interface properties, protocol parameters, security policies, and platform-wide tunables.

## Category Index
| Category | Command Count |
|---|---:|
| CLI & Terminal | 3 |
| Logging, Trace & Telemetry | 9 |
| Interface — General Properties | 21 |
| Interface — IP/IPv4 Configuration | 9 |
| Interface — IPv6 Configuration | 8 |
| Interface — L2 & Bridge Domain | 14 |
| Interface — Security, NAT & Firewall | 7 |
| Interface — Protocol Bindings | 3 |
| Bridge Domain Properties | 11 |
| IP & IPv6 Routing Configuration | 11 |
| FIB Tuning | 2 |
| DHCP & DHCPv6 | 5 |
| Classification, QoS & Policer | 4 |
| IPFIX & Flow Export | 2 |
| NAT Configuration | 8 |
| ACL Plugin | 5 |
| CNAT Source-NAT Policy | 4 |
| Crypto & IPsec | 4 |
| Host Interface & Device Tuning | 4 |
| L2 FIB & L2TP | 2 |
| LLDP | 1 |
| iOAM & Telemetry Overlay | 6 |
| VXLAN-GPE iOAM | 4 |
| NSH iOAM | 2 |
| Segment Routing | 2 |
| PNAT | 1 |
| uRPF | 2 |
| Syslog | 2 |
| UDP-Ping | 2 |
| Network Delay Simulator | 1 |
| Proof-of-Transit | 2 |
| Node & Scheduler | 2 |
| Session & Transport | 2 |
| QUIC | 2 |
| WireGuard | 1 |
| Service Chaining — SASC | 4 |
| Service Chaining — SFDP | 10 |
| Virtio & CT6 | 2 |
| IKEv2 | 1 |

## Commands by Category

### CLI & Terminal
Configure terminal display, history depth, and pager behaviour.

| Command | Description |
|---|---|
| `set terminal ansi` | Enable or disable ANSI color/escape-code output for the current CLI session. |
| `set terminal history` | Set the maximum number of commands retained in the CLI history ring. |
| `set terminal pager` | Enable, disable, or set the page size for the interactive CLI output pager. |

### Logging, Trace & Telemetry
Control log levels, trace formatting, buffer monitoring, and pcap filter functions.

| Command | Description |
|---|---|
| `set bpf trace filter` | Attach a compiled BPF program as a trace-capture filter to select which packets are traced. |
| `set buffer traces` | Enable or disable per-buffer trace tracking via the bufmon plugin. |
| `set clock adjust` | Apply a one-time or ongoing adjustment to the VPP wall-clock for time-sync compensation. |
| `set logging class` | Set the log severity threshold for a specific log class (e.g., `vlib`, `vnet`). |
| `set logging size` | Set the maximum number of entries retained in the in-memory log ring buffer. |
| `set logging unthrottle-time` | Set the minimum interval (seconds) between repeated identical log messages. |
| `set pcap filter function` | Select a compiled classifier function used to filter packets for pcap capture. |
| `set trace filter function` | Select a compiled classifier function used to filter packets for trace capture. |
| `set trace timestamp-format` | Choose the timestamp format (elapsed, calendar, or both) shown in packet trace output. |

### Interface — General Properties
Set interface admin state, naming, MTU, MAC addresses, queue placement, and feature arcs.

| Command | Description |
|---|---|
| `set interface state` | Administratively bring an interface up or down. |
| `set interface name` | Assign a user-friendly name alias to an interface (referenced elsewhere by this name). |
| `set interface tag` | Attach a free-form descriptive tag string to an interface for operational labeling. |
| `set interface mac address` | Override the hardware MAC address of an interface. |
| `set interface secondary-mac-address` | Add or remove secondary (alias) MAC addresses on an interface for multi-MAC reception. |
| `set interface mtu` | Set the maximum transmission unit (L3 payload, IP, or MPLS MTU) on an interface. |
| `set interface promiscuous` | Enable or disable promiscuous mode (receive all frames regardless of destination MAC). |
| `set interface unnumbered` | Borrow the IP address of another interface so this interface can forward without its own address. |
| `set interface hw-class` | Change the hardware-class (e.g., Ethernet to IP) of a loopback or tunnel interface. |
| `set interface handoff` | Pin an interface's input traffic to specific worker threads via RSS or static assignment. |
| `set interface rx-mode` | Set the receive mode (polling, interrupt, or adaptive) for an interface's input queues. |
| `set interface rx-placement` | Pin a specific interface rx-queue to a particular worker thread. |
| `set interface rss queues` | Configure the number and mapping of receive-side-scaling queues on an interface. |
| `set interface tx-hash` | Select the transmit hashing algorithm (L2, L3, L4 fields) for multi-queue TX. |
| `set interface tx-queue` | Configure transmit queue parameters (size, thread affinity) for an interface. |
| `set interface span` | Mirror (SPAN) ingress and/or egress traffic from one interface to another for monitoring. |
| `set interface feature` | Enable or disable a named feature-arc node on an interface's input or output path. |
| `set interface feature gso` | Enable or disable Generic Segmentation Offload on an interface. |
| `set interface reassembly` | Enable or disable IP reassembly (full or shallow-virtual) on an interface. |
| `set interface input acl` | Apply an input classify-and-act ACL (L2/IP4/IP6) to an interface. |
| `set interface output acl` | Apply an output classify-and-act ACL (L2/IP4/IP6) to an interface. |

### Interface — IP/IPv4 Configuration
Assign IP addresses, bind to VRF tables, and control per-interface IPv4 forwarding behaviour.

| Command | Description |
|---|---|
| `set interface ip address` | Add or delete an IPv4 or IPv6 address (with prefix length) on an interface. |
| `set interface ip directed-broadcast` | Enable or disable directed-broadcast forwarding on an interface. |
| `set interface ip table` | Bind an interface to a specific IPv4 VRF / FIB table for route lookups. |
| `set interface ip source-and-port-range-check` | Attach source-address and port-range verification policy to an interface. |
| `set interface ip geneve-bypass` | Enable or disable the IPv4 Geneve decap-bypass optimization on an interface. |
| `set interface ip gtpu-bypass` | Enable or disable the IPv4 GTP-U decap-bypass optimization on an interface. |
| `set interface ip pvti-bypass` | Enable or disable the IPv4 PVTI tunnel decap-bypass optimization on an interface. |
| `set interface ip vxlan-bypass` | Enable or disable the IPv4 VXLAN decap-bypass optimization on an interface. |
| `set interface ip vxlan-gpe-bypass` | Enable or disable the IPv4 VXLAN-GPE decap-bypass optimization on an interface. |

### Interface — IPv6 Configuration
Bind interfaces to IPv6 VRF tables and control ND proxy, tunnel bypass, and L2TPv3 bindings.

| Command | Description |
|---|---|
| `set interface ip6 table` | Bind an interface to a specific IPv6 VRF / FIB table for route lookups. |
| `set interface ip6 l2tpv3` | Associate an L2TPv3 tunnel session with an interface for IPv6-transported pseudowires. |
| `set interface ip6 geneve-bypass` | Enable or disable the IPv6 Geneve decap-bypass optimization on an interface. |
| `set interface ip6 gtpu-bypass` | Enable or disable the IPv6 GTP-U decap-bypass optimization on an interface. |
| `set interface ip6 pvti-bypass` | Enable or disable the IPv6 PVTI tunnel decap-bypass optimization on an interface. |
| `set interface ip6 vxlan-bypass` | Enable or disable the IPv6 VXLAN decap-bypass optimization on an interface. |
| `set interface ip6 vxlan-gpe-bypass` | Enable or disable the IPv6 VXLAN-GPE decap-bypass optimization on an interface. |
| `set interface ip6-nd proxy` | Enable or disable IPv6 Neighbor Discovery proxy on an interface. |

### Interface — L2 & Bridge Domain
Bridge-port membership, L2 learning, flooding, tag-rewrite, cross-connect, and classification.

| Command | Description |
|---|---|
| `set interface l2 bridge` | Place an interface into an L2 bridge domain as a normal port, BVI, or UU-flood port. |
| `set interface l2 efp-filter` | Enable or disable Ethernet Flow Point (EFP) ingress filtering on a bridge-domain member. |
| `set interface l2 flood` | Override per-interface L2 flood behaviour (enable/disable unknown-unicast/multicast/broadcast). |
| `set interface l2 forward` | Enable or disable L2 unicast forwarding on a bridge-domain member port. |
| `set interface l2 input classify` | Attach an L2 input classifier table to an interface for fine-grained ingress policy. |
| `set interface l2 learn` | Enable or disable MAC-address learning on a bridge-domain member port. |
| `set interface l2 output classify` | Attach an L2 output classifier table to an interface for fine-grained egress policy. |
| `set interface l2 pbb-tag-rewrite` | Configure Provider Backbone Bridge (802.1ah) tag push/pop/translate on an interface. |
| `set interface l2 rewrite` | Set Ethernet header rewrite (MAC src/dst replacement) rules on an L2 output interface. |
| `set interface l2 tag-rewrite` | Configure VLAN tag push, pop, or translate operations on an L2 sub-interface. |
| `set interface l2 xconnect` | Create a Layer-2 cross-connect (point-to-point wire) between two interfaces. |
| `set interface l2 xcrw` | Set the cross-connect rewrite (header prepend/strip) on an L2-to-L3 cross-connect path. |
| `set interface l3` | Remove an interface from its L2 bridge domain and return it to L3 (routed) mode. |
| `set bridge-domain rewrite` | Set the default Ethernet rewrite template applied to packets leaving a bridge domain. |

### Interface — Security, NAT & Firewall
Attach IPsec SPDs, NAT inside/outside roles, proxy-ARP, and TCP MSS clamping to interfaces.

| Command | Description |
|---|---|
| `set interface ipsec spd` | Bind an IPsec Security Policy Database to an interface for inbound/outbound policy lookup. |
| `set interface nat44` | Designate an interface as NAT44-ED inside, outside, or both for endpoint-dependent NAT. |
| `set interface nat44 ei` | Designate an interface as NAT44-EI inside, outside, or both for endpoint-independent NAT. |
| `set interface nat64` | Designate an interface as NAT64 inside or outside for IPv6-to-IPv4 stateful translation. |
| `set interface nat66` | Designate an interface as NAT66 inside or outside for IPv6-to-IPv6 stateful translation. |
| `set interface proxy-arp` | Enable or disable proxy-ARP replies on an interface within configured address ranges. |
| `set interface tcp-mss-clamp` | Clamp the TCP Maximum Segment Size in SYN packets traversing an interface to a specified value. |

### Interface — Protocol Bindings
Bind LLDP, MPLS, and DET44 to specific interfaces.

| Command | Description |
|---|---|
| `set interface lldp` | Enable LLDP (Link Layer Discovery Protocol) transmission and reception on an interface. |
| `set interface mpls` | Enable or disable MPLS forwarding (label imposition/disposition) on an interface. |
| `set interface det44` | Designate an interface as DET44 (deterministic NAT44) inside or outside. |

### Bridge Domain Properties
Configure learning, flooding, ARP termination, MAC aging, and forwarding within a bridge domain.

| Command | Description |
|---|---|
| `set bridge-domain arp entry` | Add or remove a static ARP entry in a bridge domain's ARP-termination table. |
| `set bridge-domain arp term` | Enable or disable ARP termination (proxy-ARP within the BD) for a bridge domain. |
| `set bridge-domain arp-ufwd` | Enable or disable ARP unicast-forwarding for a bridge domain when ARP termination is active. |
| `set bridge-domain default-learn-limit` | Set the global default MAC learning limit applied to newly created bridge domains. |
| `set bridge-domain flood` | Enable or disable unknown-unicast/multicast/broadcast flooding for a bridge domain. |
| `set bridge-domain forward` | Enable or disable L2 unicast forwarding for a bridge domain. |
| `set bridge-domain learn-limit` | Set the maximum number of dynamically learned MAC entries for a specific bridge domain. |
| `set bridge-domain learn` | Enable or disable MAC-address learning for a bridge domain. |
| `set bridge-domain mac-age` | Set the MAC aging timeout (minutes) for dynamically learned entries in a bridge domain. |
| `set bridge-domain uu-flood` | Set or clear the unknown-unicast flood interface for a bridge domain. |

### IP & IPv6 Routing Configuration
Set IP/IPv6 classification, flow-hashing, neighbor entries, DAD, ND proxy, and addresses.

| Command | Description |
|---|---|
| `set ip classify` | Attach a classify table to IPv4 input for policy-based forwarding or filtering. |
| `set ip flow-hash` | Configure the fields (src, dst, sport, dport, proto) used for IPv4 ECMP flow hashing. |
| `set ip neighbor-config` | Set IPv4 neighbor (ARP) table parameters: max entries, age limit, and recycle threshold. |
| `set ip neighbor` | Add, delete, or modify a static IPv4 ARP neighbor entry on an interface. |
| `set ip source-and-port-range-check` | Configure source-address and port-range allowlists for an IPv4 VRF table. |
| `set ip6 address` | Add or delete an IPv6 address on an interface (alternative to `set interface ip address`). |
| `set ip6 classify` | Attach a classify table to IPv6 input for policy-based forwarding or filtering. |
| `set ip6 dad disable` | Disable IPv6 Duplicate Address Detection on an interface. |
| `set ip6 dad enable` | Enable IPv6 Duplicate Address Detection on an interface. |
| `set ip6 flow-hash` | Configure the fields used for IPv6 ECMP flow hashing. |
| `set ip6 nd proxy` | Configure IPv6 Neighbor Discovery proxy for a prefix on an interface. |

### FIB Tuning
Adjust FIB walk quotas and histogram bucket sizes.

| Command | Description |
|---|---|
| `set fib walk histogram elements size` | Set the bucket size for the FIB-walk duration histogram. |
| `set fib walk quota` | Set the maximum number of FIB entries processed per walk iteration before yielding. |

### DHCP & DHCPv6
Configure DHCP client/proxy, Option-82 VSS, and DHCPv6 proxy.

| Command | Description |
|---|---|
| `set dhcp client` | Start or stop a DHCPv4 client on an interface to acquire an address via DHCP. |
| `set dhcp option-82 vss` | Set the DHCP Option-82 Virtual Subnet Selection (VSS) parameters for relay agent info. |
| `set dhcp proxy` | Configure a DHCPv4 relay-proxy server address and source address for a VRF. |
| `set dhcpv6 proxy` | Configure a DHCPv6 relay-proxy server address and source address for a VRF. |
| `set dhcpv6 vss` | Set the DHCPv6 Virtual Subnet Selection (VSS) parameters for relay agent info. |

### Classification, QoS & Policer
Attach classifiers for flow-based, policer-based, or IPFIX-based classification.

| Command | Description |
|---|---|
| `set flow classify` | Attach a flow-classifier table to an interface for per-flow traffic identification. |
| `set policer classify` | Attach a policer-classifier table (ip4/ip6/l2) to an interface for rate-limiting. |
| `set ipfix classify stream` | Bind an IPFIX export stream to a classifier table for per-flow telemetry export. |
| `set ipfix exporter` | Configure the IPFIX/Netflow exporter: collector IP, source IP, port, template interval, and path-MTU. |

### NAT Configuration
Set NAT timeouts, worker affinity, frame-queue sizing, and session limits.

| Command | Description |
|---|---|
| `set nat frame-queue-nelts` | Set the number of elements in the NAT inter-worker frame queue. |
| `set nat timeout` | Set global NAT session idle timeouts (TCP established, TCP transitory, UDP, ICMP). |
| `set nat workers` | Assign specific worker threads to handle NAT translation processing. |
| `set nat44 ei timeout` | Set NAT44 endpoint-independent session idle timeouts per protocol. |
| `set nat44 ei workers` | Assign specific worker threads for NAT44-EI translation processing. |
| `set nat44 session limit` | Set the maximum number of concurrent NAT44-ED translation sessions per user or globally. |
| `set det44 timeouts` | Set DET44 (deterministic NAT) session timeouts per protocol. |
| `set node function` | Select an alternative implementation (e.g., SIMD variant) for a graph node's dispatch function. |

### ACL Plugin
Configure ACL plugin parameters, create/update ACL rules, and bind ACLs to interfaces.

| Command | Description |
|---|---|
| `set acl-plugin` | Set global ACL plugin parameters (connection-table sizing, timeout tunables). |
| `set acl-plugin acl` | Create or replace an L3/L4 ACL rule set (permit/deny with 5-tuple match). |
| `set acl-plugin interface` | Bind or unbind an L3/L4 ACL to an interface's input or output path. |
| `set acl-plugin macip acl` | Create or replace a MAC+IP ACL rule set for combined L2+L3 filtering. |
| `set acl-plugin macip interface` | Bind or unbind a MACIP ACL to an interface. |

### CNAT Source-NAT Policy
Configure cloud-NAT source-NAT policy rules, address pools, and interface/prefix scoping.

| Command | Description |
|---|---|
| `set cnat snat-policy` | Set the global CNAT source-NAT policy (none, if-pfx, or k8s). |
| `set cnat snat-policy addr` | Add or remove addresses in the CNAT source-NAT pool. |
| `set cnat snat-policy if` | Mark an interface for CNAT source-NAT policy inclusion or exclusion. |
| `set cnat snat-policy prefix` | Add or remove IP prefixes that bypass CNAT source-NAT (no-SNAT subnets). |

### Crypto & IPsec
Select crypto handlers, DPDK cryptodev assignment, and async dispatch modes.

| Command | Description |
|---|---|
| `set crypto handler` | Select the engine (openssl, native, ipsecmb, DPDK) for a specific crypto algorithm. |
| `set dpdk cryptodev assignment` | Assign DPDK crypto devices to worker threads for hardware-accelerated encryption. |
| `set ipsec async mode` | Enable or disable asynchronous crypto dispatch for IPsec SA processing. |
| `set wireguard async mode` | Enable or disable asynchronous crypto dispatch for WireGuard tunnel processing. |

### Host Interface & Device Tuning
Configure AF_PACKET checksum offload, qdisc bypass, and DPDK descriptor ring sizes.

| Command | Description |
|---|---|
| `set host-interface l4-cksum-offload` | Enable or disable L4 (TCP/UDP) checksum offload on an AF_PACKET host interface. |
| `set host-interface qdisc-bypass` | Enable or disable Linux qdisc bypass on an AF_PACKET host interface for lower TX latency. |
| `set dpdk interface descriptors` | Set the number of RX/TX descriptors (ring size) for a DPDK-managed interface. |
| `set sw_scheduler` | Configure the software crypto-scheduler worker thread assignment and queue parameters. |

### L2 FIB & L2TP
Tune L2 FIB scan interval and L2TPv3 tunnel cookies.

| Command | Description |
|---|---|
| `set l2fib scan-delay` | Set the interval (seconds) between L2 FIB aging scans for dynamically learned MACs. |
| `set l2tpv3 tunnel cookie` | Set or change the authentication cookie value on an existing L2TPv3 tunnel session. |

### LLDP
Configure global LLDP parameters.

| Command | Description |
|---|---|
| `set lldp` | Set global LLDP parameters: system-name, tx-hold multiplier, and tx-interval. |

### iOAM & Telemetry Overlay
Configure in-band OAM trace profiles, analysis, caching, and IPFIX export.

| Command | Description |
|---|---|
| `set ioam analyse` | Enable or disable iOAM trace-data analysis on received packets. |
| `set ioam export ipfix` | Configure IPFIX export of iOAM hop-by-hop telemetry data to a collector. |
| `set ioam ip6 cache` | Enable or disable the iOAM IPv6 reassembly/analysis cache for aggregated trace data. |
| `set ioam ip6 sr-tunnel-select` | Select which SRv6 tunnel policy iOAM telemetry data should be associated with. |
| `set ioam rewrite` | Configure the iOAM hop-by-hop extension header rewrite (trace, PoT, edge-to-edge). |
| `set ioam-trace profile` | Define an iOAM trace profile: trace-type, node-id, app-data, and number of hops. |

### VXLAN-GPE iOAM
Configure iOAM telemetry within VXLAN-GPE tunnels.

| Command | Description |
|---|---|
| `set vxlan-gpe-ioam` | Enable or disable iOAM data insertion for a VXLAN-GPE tunnel. |
| `set vxlan-gpe-ioam export ipfix` | Configure IPFIX export of iOAM telemetry collected within VXLAN-GPE tunnels. |
| `set vxlan-gpe-ioam rewrite` | Configure the iOAM header rewrite for VXLAN-GPE encapsulated packets. |
| `set vxlan-gpe-ioam-transit` | Enable iOAM processing on a transit node for VXLAN-GPE encapsulated traffic. |

### NSH iOAM
Configure iOAM telemetry within NSH service-function-chaining paths.

| Command | Description |
|---|---|
| `set nsh-md2-ioam export ipfix` | Configure IPFIX export of iOAM telemetry collected within NSH MD-type-2 headers. |
| `set nsh-md2-ioam-transit` | Enable iOAM processing on a transit node for NSH MD-type-2 service paths. |

### Segment Routing
Set SRv6 encapsulation source address and hop-limit.

| Command | Description |
|---|---|
| `set sr encaps hop-limit` | Set the IPv6 hop-limit value used in SRv6 encapsulation headers. |
| `set sr encaps source` | Set the source IPv6 address used in SRv6 encapsulation headers. |

### PNAT
Configure Policy-NAT translation rules.

| Command | Description |
|---|---|
| `set pnat translation` | Add, modify, or delete a PNAT (Policy NAT) 1:1 address/port translation rule on an interface. |

### uRPF
Configure Unicast Reverse Path Forwarding checks.

| Command | Description |
|---|---|
| `set urpf` | Enable strict or loose uRPF (Unicast Reverse Path Forwarding) on an interface. |
| `set urpf-accept` | Configure a uRPF accept-list to exempt specific prefixes from RPF drop. |

### Syslog
Configure syslog (RFC 5424) export parameters and severity filters.

| Command | Description |
|---|---|
| `set syslog filter` | Set the minimum severity level for messages exported via the syslog sender. |
| `set syslog sender` | Configure the syslog exporter: collector IP/port, source IP, and maximum message size. |

### UDP-Ping
Configure UDP-ping active probing and its IPFIX export.

| Command | Description |
|---|---|
| `set udp-ping` | Configure UDP-ping probes: target IP/port, interval, and number of packets. |
| `set udp-ping export-ipfix` | Enable or disable IPFIX export of UDP-ping round-trip-time measurement results. |

### Network Delay Simulator
Configure the network simulator (nsim) plugin parameters.

| Command | Description |
|---|---|
| `set nsim` | Set nsim parameters: delay, bandwidth, packet-size, loss-ratio, and reorder-rate. |

### Proof-of-Transit
Configure PoT verification profiles.

| Command | Description |
|---|---|
| `set pot profile` | Configure a Proof-of-Transit profile with secret shares, prime, and validator fields. |
| `set pot profile-active` | Select which PoT profile index is currently active for validation. |

### Session & Transport
Configure punt behaviour and session-layer scheduling.

| Command | Description |
|---|---|
| `set punt` | Configure punt-to-host behaviour: protocol, port, and exception path. |
| `set node function` | Select an alternative implementation (e.g., SIMD variant) for a graph node's dispatch function. |

### QUIC
Tune QUIC transport parameters.

| Command | Description |
|---|---|
| `set quic cc` | Select the QUIC congestion-control algorithm (newreno, cubic). |
| `set quic max_packets_per_key` | Set the maximum number of packets encrypted with a single QUIC key before rotation. |

### WireGuard
Configure WireGuard async crypto mode.

| Command | Description |
|---|---|
| `set wireguard async mode` | Enable or disable asynchronous crypto dispatch for WireGuard tunnel processing. |

### Service Chaining — SASC
Configure SASC (Service-Aware Service Chaining) tenant, service, and interface bindings.

| Command | Description |
|---|---|
| `set sasc ingress interface` | Designate an interface as an SASC ingress point for service-chain traffic steering. |
| `set sasc services` | Register or update the set of services available within an SASC tenant's chain. |
| `set sasc tenant` | Create, modify, or remove an SASC tenant definition (ID, name, service-chain binding). |
| `set sasc timeout` | Set the idle-session timeout for SASC-tracked flows. |

### Service Chaining — SFDP
Configure SFDP (Stateful Datapath) services, NAT pools, ACLs, and session parameters.

| Command | Description |
|---|---|
| `set sfdp acl` | Attach or detach an ACL to the SFDP classification stage for traffic selection. |
| `set sfdp eviction sessions-margin` | Set the session-count margin that triggers proactive SFDP session eviction. |
| `set sfdp gateway geneve-output` | Configure the GENEVE output tunnel endpoint for SFDP gateway-mode forwarding. |
| `set sfdp icmp-error-node` | Set the graph node used to generate ICMP error responses for SFDP-dropped packets. |
| `set sfdp interface-input` | Designate an interface as an SFDP input point for stateful datapath processing. |
| `set sfdp nat external-interface` | Set the external (public-facing) interface for SFDP integrated NAT. |
| `set sfdp nat snat` | Configure the SFDP source-NAT address or pool for outbound translation. |
| `set sfdp services` | Register or update services available within the SFDP framework. |
| `set sfdp sp-node` | Set the graph node used for SFDP service-policy enforcement. |
| `set sfdp timeout` | Set the idle-session timeout for SFDP-tracked flows. |

### Virtio & CT6
Configure VirtIO PCI parameters and the IPv6 connection tracker.

| Command | Description |
|---|---|
| `set virtio pci` | Set VirtIO PCI device parameters (features, ring sizes) after device creation. |
| `set ct6` | Enable or disable the IPv6 connection tracker (ct6) on an interface. |

### IKEv2
Configure IKEv2 authentication key.

| Command | Description |
|---|---|
| `set ikev2 local key` | Load or set the local private key (file path or inline PEM) used for IKEv2 authentication. |
