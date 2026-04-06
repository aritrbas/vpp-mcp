# VPP `vppctl` Complete Command Summary

## Purpose
This document provides a single flat reference of **every** VPP CLI command with a concise one-liner description.
Commands are listed alphabetically. For categorized views with deeper context, see the per-verb catalogs:

| Catalog | File |
|---|---|
| `show` commands | `vppctl-show-command-catalog.md` |
| `set` commands | `vppctl-set-command-catalog.md` |
| `clear` commands | `vppctl-clear-command-catalog.md` |
| `create` commands | `vppctl-create-command-catalog.md` |
| `delete` commands | `vppctl-delete-command-catalog.md` |
| `test` commands | `vppctl-test-command-catalog.md` |
| Operational commands | `vppctl-operational-command-catalog.md` |

---

## All Commands (A–Z)

| Command | Description |
|---|---|
| `abf attach` | Attach an ACL-Based Forwarding policy to an interface for policy-based routing. |
| `abf policy` | Create or delete an ABF policy binding an ACL index to forwarding paths. |
| `adjacency counters` | Enable or disable per-adjacency packet/byte counters. |
| `adl allowlist` | Add or remove entries in the Allow/Deny List for an interface. |
| `adl interface` | Enable or disable ADL processing on an interface. |
| `api trace` | Enable, disable, or dump binary API message tracing. |
| `app crypto add tls-profile` | Register a named TLS profile (certificate, key, CA) for application use. |
| `app crypto del tls-profile` | Remove a named TLS profile from the certificate store. |
| `app evt-collector` | Configure an application event collector for session-layer telemetry. |
| `app ns` | Create or delete an application namespace with secret, FIB table, and scope. |
| `arping` | Send gratuitous or solicited ARP/NDP requests from an interface. |
| `auto-sdl` | Enable or disable automatic Session Descriptor List generation for DDoS mitigation. |
| `bfd key del` | Delete a BFD authentication key by key-id. |
| `bfd key set` | Create or update a BFD authentication key (SHA1 or meticulous-SHA1). |
| `bfd udp echo-source del` | Remove the interface used as the BFD echo source. |
| `bfd udp echo-source set` | Set the interface whose address is used as the BFD echo-function source. |
| `bfd udp session add` | Create a BFD-over-UDP session to a remote peer for liveliness detection. |
| `bfd udp session auth activate` | Activate authentication on an existing BFD session. |
| `bfd udp session auth deactivate` | Deactivate authentication on a BFD session. |
| `bfd udp session del` | Delete a BFD-over-UDP session. |
| `bfd udp session mod` | Modify the timing intervals of an existing BFD session. |
| `bfd udp session set-flags` | Set administrative flags on a BFD session. |
| `bier route` | Add or delete a BIER forwarding entry in a BIER table. |
| `bier table` | Create or delete a BIER table (set, sub-domain, bit-string-length). |
| `binary-api` | Send a raw binary API message from the CLI. |
| `bond add` | Add a member (slave) interface to a bonding group. |
| `bond del` | Remove a member (slave) interface from a bonding group. |
| `buffer metadata tracking` | Enable or disable per-buffer metadata change tracking. |
| `bvi create` | Create a Bridge Virtual Interface for routing between bridge domains. |
| `bvi delete` | Delete a previously created BVI interface. |
| `cdp` | Enable or disable Cisco Discovery Protocol processing globally. |
| `classify filter` | Attach a classify table as a pcap or trace capture filter. |
| `classify session` | Add or delete a session (match rule + action) in a classify table. |
| `classify table` | Create or delete a classify table with mask and memory-size parameters. |
| `clear acl-plugin sessions` | Purge all ACL stateful-firewall connection-tracking sessions. |
| `clear api histogram` | Reset the API message processing latency histogram. |
| `clear buffer traces` | Clear buffer-monitoring trace data from the bufmon plugin. |
| `clear errors` | Zero all per-node packet-processing error counters. |
| `clear fib walk` | Reset FIB convergence-walk statistics. |
| `clear hardware-interfaces` | Reset hardware-interface counters on all or selected interfaces. |
| `clear http connect proxy client stats` | Reset HTTP CONNECT proxy client counters. |
| `clear http static cache` | Invalidate all HTTP static file-server cache entries. |
| `clear http stats` | Zero HTTP layer statistics. |
| `clear igmp` | Reset IGMP protocol state and group membership records. |
| `clear interface tag` | Remove the descriptive tag string from an interface. |
| `clear interfaces` | Reset all software-interface counters to zero. |
| `clear ioam rewrite` | Remove the active iOAM hop-by-hop header rewrite configuration. |
| `clear ioam-trace profile` | Delete the currently configured iOAM trace profile. |
| `clear ipsec counters` | Zero per-SA and per-SPD IPsec counters. |
| `clear ipsec sa` | Remove a specific IPsec Security Association. |
| `clear l2fib` | Remove all dynamically-learned L2 FIB entries. |
| `clear l2tp counters` | Zero L2TPv3 per-session packet and byte counters. |
| `clear logging` | Flush the in-memory log buffer. |
| `clear mactime` | Reset MAC-time allow/drop counters. |
| `clear nat44 ed sessions` | Flush all active NAT44 endpoint-dependent sessions. |
| `clear nat44 ei sessions` | Flush all active NAT44 endpoint-independent sessions. |
| `clear node counters` | Reset all per-node packet counters to zero. |
| `clear pot profile` | Remove all Proof-of-Transit profiles. |
| `clear runtime` | Zero per-node runtime statistics (clocks, vectors, calls). |
| `clear sasc sessions` | Flush all SASC-tracked sessions. |
| `clear session` | Forcefully clean up host-stack sessions and release FIFOs. |
| `clear session stats` | Reset session-layer aggregate statistics. |
| `clear sr localsid-counters` | Zero counters on all SRv6 local-SID entries. |
| `clear tcp stats` | Zero all TCP transport-layer counters. |
| `clear trace` | Discard all packets in the packet trace buffer. |
| `clear vxlan-gpe-ioam rewrite` | Remove the VXLAN-GPE iOAM header rewrite configuration. |
| `cnat client add` | Register a CNAT client endpoint for translation tracking. |
| `cnat log` | Enable or disable CNAT session logging. |
| `cnat translation` | Create, modify, or delete a CNAT translation rule (VIP → backends). |
| `configure policer` | Create or update a hierarchical policer by name. |
| `create 6rd tunnel` | Create a 6rd (IPv6 Rapid Deployment) tunnel for IPv6-over-IPv4 encapsulation. |
| `create bond` | Create a bonding master interface (LACP, round-robin, active-backup). |
| `create bridge-domain` | Create an L2 bridge domain with learning, flooding, and ARP termination. |
| `create geneve tunnel` | Create a GENEVE tunnel endpoint for network virtualization. |
| `create gre tunnel` | Create a GRE tunnel (IPv4/IPv6, ERSPAN, or TEB mode). |
| `create gtpu forward` | Create a GTPU forwarding rule to a next-hop or interface. |
| `create gtpu tunnel` | Create a GTP-U tunnel for mobile backhaul user-plane encapsulation. |
| `create host-interface` | Create an AF_PACKET interface attached to a Linux netdev. |
| `create interface af_xdp` | Create an AF_XDP interface bound to a Linux netdev via eBPF/XDP. |
| `create interface memif` | Create a memif shared-memory interface for inter-process communication. |
| `create interface rdma` | Create an RDMA (mlx5) interface for direct NIC queue access. |
| `create interface virtio` | Create a virtio-backed network interface. |
| `create interface vmxnet3` | Create a VMware vmxnet3 paravirtualized interface. |
| `create ipip tunnel` | Create an IP-in-IP tunnel for simple IP encapsulation. |
| `create l2tpv3 tunnel` | Create an L2TPv3 tunnel session for Ethernet-over-IP pseudowire. |
| `create loopback interface` | Create a software loopback interface. |
| `create memif socket` | Create a named memif socket file for memif connections. |
| `create nsh entry` | Create a Network Service Header entry (SPI/SI, metadata). |
| `create nsh map` | Create an NSH forwarding map binding SPI/SI to a next-hop action. |
| `create packet-generator interface` | Create a virtual packet-generator interface. |
| `create pppoe cp` | Create a PPPoE control-plane session. |
| `create pppoe session` | Create a PPPoE data-plane session. |
| `create sub-interfaces` | Create VLAN sub-interfaces (dot1q/dot1ad) on a parent interface. |
| `create tap` | Create a TUN/TAP interface connecting VPP to the Linux kernel. |
| `create teib` | Create a Tunnel Endpoint Information Base entry. |
| `create vhost-user` | Create a vhost-user interface for virtio-based VM networking. |
| `create vxlan tunnel` | Create a VXLAN tunnel endpoint with a 24-bit VNI. |
| `create vxlan-gpe tunnel` | Create a VXLAN-GPE tunnel with multi-protocol encapsulation. |
| `debug lacp` | Enable or disable LACP protocol debug logging. |
| `define` | Define a CLI macro (alias) expanding to a command sequence. |
| `delete 6rd tunnel` | Tear down a 6rd tunnel. |
| `delete acl-plugin acl` | Delete an L3/L4 ACL rule set by index. |
| `delete acl-plugin macip acl` | Delete a MAC+IP ACL rule set by index. |
| `delete bond` | Delete a bonding master interface and release members. |
| `delete host-interface` | Destroy an AF_PACKET host interface. |
| `delete interface af_xdp` | Destroy an AF_XDP interface. |
| `delete interface memif` | Destroy a memif shared-memory interface. |
| `delete interface rdma` | Destroy an RDMA (mlx5) interface. |
| `delete interface virtio` | Destroy a virtio-backed interface. |
| `delete interface vmxnet3` | Destroy a vmxnet3 paravirtualized interface. |
| `delete ipip tunnel` | Tear down an IP-in-IP tunnel. |
| `delete loopback interface` | Remove a software loopback interface. |
| `delete memif socket` | Remove a memif socket file entry. |
| `delete packet-generator interface` | Remove a packet-generator interface. |
| `delete sub-interface` | Remove a VLAN sub-interface. |
| `delete tap` | Destroy a TUN/TAP interface. |
| `delete teib` | Remove a TEIB entry. |
| `delete vhost-user` | Destroy a vhost-user interface. |
| `det44 add` | Add a DET44 inside-to-outside address mapping range. |
| `det44 close session in` | Close a DET44 session by inside 5-tuple. |
| `det44 close session out` | Close a DET44 session by outside 5-tuple. |
| `det44 forward` | Display the outside address/port for a DET44 inside address. |
| `det44 plugin` | Enable or disable the DET44 deterministic-NAT plugin. |
| `det44 reverse` | Display the inside address for a DET44 outside address and port. |
| `device attach` | Attach a device to the VPP device-driver framework. |
| `device create-interface` | Create a primary network interface on an attached device. |
| `device create-secondary-interface` | Create a secondary interface on an attached device. |
| `device detach` | Detach a device from the VPP device-driver framework. |
| `device remove-interface` | Remove a primary interface from a device. |
| `device remove-secondary-interface` | Remove a secondary interface from a device. |
| `device reset` | Reset a device to its initial state. |
| `device set-rss-key` | Set the RSS hash key on a device. |
| `dhcp6 client` | Start or stop a DHCPv6 (IA_NA) client on an interface. |
| `dhcp6 pd client` | Start or stop a DHCPv6 Prefix Delegation client. |
| `disable ip6 interface` | Disable IPv6 processing on an interface. |
| `dns` | Enable or disable the VPP DNS resolver cache. |
| `dns cache` | Add or remove static DNS cache entries. |
| `dns name-server` | Add or remove an upstream DNS name-server address. |
| `dslite add pool address` | Add an IPv4 address to the DS-Lite AFTR NAT pool. |
| `dslite set aftr-tunnel-endpoint-address` | Set the DS-Lite AFTR tunnel endpoint IPv6 address. |
| `dslite set b4-tunnel-endpoint-address` | Set the DS-Lite B4 tunnel endpoint IPv6 address. |
| `dump sasc session` | Dump detailed SASC session state. |
| `dump sasc session ring` | Dump the SASC session ring buffer. |
| `echo` | Print a literal string to CLI output. |
| `enable ip4 interface` | Enable IPv4 processing on an interface. |
| `enable ip6 interface` | Enable IPv6 processing and link-local address on an interface. |
| `event-logger clear` | Discard all events in the event-log buffer. |
| `event-logger resize` | Change the event-log buffer size. |
| `event-logger restart` | Re-enable event logging after stop. |
| `event-logger save` | Write event-log buffer to a file. |
| `event-logger stop` | Stop recording events. |
| `event-logger trace` | Configure which event classes are captured. |
| `exec` | Execute CLI commands from an external file. |
| `exit` | Close the current CLI session. |
| `flow` | Add, delete, or enable/disable a hardware-offloaded flow rule. |
| `flowprobe feature add-del` | Enable or disable the flowprobe (IPFIX) feature on an interface. |
| `flowprobe params` | Set flowprobe active/passive timeout and record format. |
| `gpe` | Enable or disable GPE datapath. |
| `gpe encap` | Set the GPE encapsulation mode (LISP, VXLAN, NSH). |
| `gpe entry` | Add or delete a GPE forwarding entry. |
| `gpe iface` | Create or delete a GPE interface for a VNI. |
| `gpe native-forward` | Add or delete a GPE native-forward entry. |
| `history` | Display CLI command history. |
| `hsi` | Start the built-in HTTP server index test application. |
| `http cli client` | Send an HTTP request to VPP's CLI-over-HTTP endpoint. |
| `http cli server` | Enable CLI command access via HTTP GET/POST. |
| `http client` | Launch a one-shot HTTP client request from VPP. |
| `http connect proxy client enable` | Enable the HTTP CONNECT proxy client. |
| `http connect proxy client listener` | Configure the HTTP CONNECT proxy listener address and port. |
| `http static listener` | Configure the HTTP static file-server listener. |
| `http static server` | Enable the HTTP static file server with cache and FIFO parameters. |
| `http tps` | Start the HTTP transactions-per-second benchmark application. |
| `igmp` | Enable or disable IGMP protocol processing. |
| `igmp listen` | Add or remove a static IGMP group membership on an interface. |
| `igmp proxy-dev` | Create or delete an IGMP proxy device. |
| `igmp proxy-dev itf` | Add or remove a downstream interface from an IGMP proxy device. |
| `ikev2 dpd disable` | Disable IKEv2 Dead Peer Detection for a profile. |
| `ikev2 initiate` | Manually initiate IKEv2 SA negotiation. |
| `ikev2 profile` | Create, modify, or delete an IKEv2 profile. |
| `ikev2 set liveness` | Set IKEv2 liveness-check interval and threshold. |
| `ikev2 set sleep interval` | Set IKEv2 background process sleep interval. |
| `ila entry` | Add or delete an Identifier-Locator Addressing translation entry. |
| `ila interface` | Enable or disable ILA processing on an interface. |
| `interface collect detailed-stats` | Enable detailed per-interface statistics collection. |
| `ip container` | Add or remove an IP container-proxy entry. |
| `ip mroute` | Add or delete an IPv4 multicast route. |
| `ip neighbor` | Add, delete, or flush an IPv4/IPv6 neighbor entry. |
| `ip punt policer` | Apply a policer to punted IPv4 packets. |
| `ip punt redirect` | Redirect punted IPv4 packets to a next-hop. |
| `ip route` | Add, delete, or modify an IPv4 unicast route. |
| `ip session redirect` | Add or delete a session-redirect rule for flow steering. |
| `ip syn filter` | Enable or disable TCP SYN flood filtering on IPv4 punt. |
| `ip table` | Create or delete an IPv4 VRF / FIB table. |
| `ip virtual` | Set a virtual (floating) IPv4 address for ARP/punt. |
| `ip6 nd` | Configure IPv6 Neighbor Discovery parameters on an interface. |
| `ip6 nd address autoconfig` | Enable or disable SLAAC on an interface. |
| `ip6 punt policer` | Apply a policer to punted IPv6 packets. |
| `ip6 punt redirect` | Redirect punted IPv6 packets to a next-hop. |
| `ip6 table` | Create or delete an IPv6 VRF / FIB table. |
| `ipfix classify table` | Bind a classify table to IPFIX flow export. |
| `ipfix flush` | Force-flush all buffered IPFIX records. |
| `ipsec itf create` | Create an IPsec tunnel interface. |
| `ipsec itf delete` | Delete an IPsec tunnel interface. |
| `ipsec policy` | Add or delete an IPsec Security Policy in an SPD. |
| `ipsec sa` | Add or delete an IPsec Security Association. |
| `ipsec sa bind` | Bind an IPsec SA to a specific worker thread. |
| `ipsec spd` | Create or delete an IPsec Security Policy Database. |
| `ipsec tunnel protect` | Apply IPsec tunnel-protect (SA binding) to a tunnel. |
| `kill sfdp session` | Forcefully terminate a specific SFDP session. |
| `l2 rewrite entry` | Add or delete an L2 output-rewrite entry. |
| `l2fib add` | Add a static L2 FIB entry. |
| `l2fib del` | Delete an L2 FIB entry by MAC and bridge-domain. |
| `l2fib flush-mac all` | Flush all dynamic L2 FIB entries globally. |
| `l2fib flush-mac bridge-domain` | Flush dynamic L2 FIB entries in a bridge domain. |
| `l2fib flush-mac interface` | Flush dynamic L2 FIB entries for an interface. |
| `l3xc` | Create or delete an L3 cross-connect rule. |
| `lb as` | Add or delete a load-balancer back-end server. |
| `lb conf` | Set global load-balancer parameters. |
| `lb flush vip` | Flush per-flow state for a load-balancer VIP. |
| `lb set interface nat4` | Enable or disable LB NAT44 on an interface. |
| `lb set interface nat6` | Enable or disable LB NAT66 on an interface. |
| `lb vip` | Create or delete a load-balancer Virtual IP. |
| `lisp adjacency` | Add or delete a LISP adjacency. |
| `lisp disable` | Disable the LISP control plane. |
| `lisp eid-table` | Add or delete a LISP EID-table mapping. |
| `lisp eid-table map` | Map a LISP EID-table to a bridge-domain or VRF. |
| `lisp enable` | Enable the LISP control plane. |
| `lisp locator` | Add or delete a locator in a LISP locator-set. |
| `lisp locator-set` | Create or delete a named LISP locator-set. |
| `lisp map-register` | Enable or disable LISP map-register sending. |
| `lisp map-request itr-rlocs` | Set LISP map-request source RLOCs. |
| `lisp map-request mode` | Set LISP map-request mode (src-dst or dst-only). |
| `lisp map-resolver` | Add or delete a LISP map-resolver address. |
| `lisp map-server` | Add or delete a LISP map-server address. |
| `lisp pitr` | Enable or disable LISP Proxy Ingress Tunnel Router. |
| `lisp remote-mapping` | Add or delete a remote LISP map-cache entry. |
| `lisp rloc-probe` | Enable or disable LISP RLOC probing. |
| `lisp use-petr` | Enable or disable LISP Proxy ETR usage. |
| `loopback create-interface` | Create a software loopback interface (alternative syntax). |
| `loopback delete-interface` | Delete a software loopback interface (alternative syntax). |
| `mactime enable-disable` | Enable or disable MAC-time access control on an interface. |
| `map add domain` | Create a MAP-E or MAP-T domain. |
| `map add rule` | Add a forwarding rule to a MAP domain. |
| `map del domain` | Delete a MAP domain. |
| `map interface` | Enable or disable MAP processing on an interface. |
| `map params fragment` | Configure MAP fragmentation parameters. |
| `map params icmp source-address` | Set the source IPv4 address for MAP ICMP errors. |
| `map params icmp6 unreachables` | Enable or disable ICMPv6 unreachables for MAP errors. |
| `map params pre-resolve` | Configure MAP pre-resolved next-hop addresses. |
| `map params security-check` | Enable or disable MAP source-address security check. |
| `map params tcp-mss` | Set TCP MSS clamping for MAP-translated segments. |
| `map params traffic-class` | Set traffic-class handling for MAP encapsulation. |
| `memory-trace` | Enable or disable per-allocation memory tracing for leak detection. |
| `monitor interface` | Start real-time counter display for an interface. |
| `mpls local-label` | Bind a local MPLS label to an IP prefix. |
| `mpls table` | Create or delete an MPLS label table. |
| `mpls tunnel` | Create, modify, or delete an MPLS tunnel. |
| `nat ipfix logging` | Enable or disable IPFIX logging of NAT events. |
| `nat mss-clamping` | Set TCP MSS clamping for NAT-translated packets. |
| `nat set logging level` | Set syslog severity level for NAT events. |
| `nat44 add address` | Add addresses to the NAT44-ED external pool. |
| `nat44 add identity mapping` | Create a NAT44-ED identity (no-translate) mapping. |
| `nat44 add interface address` | Use an interface address as a NAT44-ED pool address. |
| `nat44 add load-balancing back-end` | Add a back-end to a NAT44-ED load-balanced mapping. |
| `nat44 add load-balancing static mapping` | Create a NAT44-ED load-balanced static mapping. |
| `nat44 add static mapping` | Create a NAT44-ED 1:1 or port-restricted static mapping. |
| `nat44 del session` | Delete a specific NAT44-ED session. |
| `nat44 ei add address` | Add addresses to the NAT44-EI external pool. |
| `nat44 ei add identity mapping` | Create a NAT44-EI identity mapping. |
| `nat44 ei add interface address` | Use an interface address as a NAT44-EI pool address. |
| `nat44 ei add static mapping` | Create a NAT44-EI static mapping. |
| `nat44 ei addr-port-assignment-alg` | Select the NAT44-EI address/port assignment algorithm. |
| `nat44 ei del session` | Delete a specific NAT44-EI session. |
| `nat44 ei del user` | Delete all NAT44-EI sessions for an inside user. |
| `nat44 ei forwarding` | Enable or disable NAT44-EI pass-through forwarding. |
| `nat44 ei ha failover` | Configure NAT44-EI HA failover parameters. |
| `nat44 ei ha flush` | Flush NAT44-EI HA session-sync state. |
| `nat44 ei ha listener` | Configure NAT44-EI HA listener for session-sync. |
| `nat44 ei ha resync` | Trigger NAT44-EI HA full session resync. |
| `nat44 ei ipfix logging` | Enable or disable IPFIX logging for NAT44-EI. |
| `nat44 ei mss-clamping` | Set TCP MSS clamping for NAT44-EI. |
| `nat44 ei plugin` | Enable or disable the NAT44-EI plugin. |
| `nat44 ei set logging level` | Set syslog level for NAT44-EI events. |
| `nat44 forwarding` | Enable or disable NAT44-ED forwarding. |
| `nat44 plugin` | Enable or disable the NAT44-ED plugin. |
| `nat44 vrf route` | Add or remove a NAT44-ED inter-VRF routing entry. |
| `nat44 vrf table` | Create or delete a NAT44-ED VRF routing table. |
| `nat64 add interface address` | Use an interface address as a NAT64 pool address. |
| `nat64 add pool address` | Add addresses to the NAT64 pool. |
| `nat64 add prefix` | Add a NAT64 IPv6 prefix (e.g., 64:ff9b::/96). |
| `nat64 add static bib` | Add a static NAT64 BIB entry. |
| `nat64 plugin` | Enable or disable the NAT64 plugin. |
| `nat66 add static mapping` | Add a static NAT66 1:1 IPv6 mapping. |
| `nat66 plugin` | Enable or disable the NAT66 plugin. |
| `npol interface clear` | Remove all network-policy bindings from an interface. |
| `npol interface configure` | Attach a network-policy to an interface. |
| `npol ipset add` | Create a named IP-set for policy rules. |
| `npol ipset add member` | Add an IP prefix to an IP-set. |
| `npol ipset del` | Delete a named IP-set. |
| `npol ipset del member` | Remove an IP prefix from an IP-set. |
| `npol match` | Test a packet tuple against network-policy rules. |
| `npol policy add` | Create a named network-policy. |
| `npol policy del` | Delete a named network-policy. |
| `npol rule add` | Add a match rule to a network-policy. |
| `npol rule del` | Delete a rule from a network-policy. |
| `nsim cross-connect enable-disable` | Enable or disable nsim cross-connect delay between two interfaces. |
| `nsim output-feature enable-disable` | Enable or disable nsim as an output feature on an interface. |
| `one adjacency` | Add or delete a ONE overlay adjacency. |
| `one disable` | Disable the ONE control plane. |
| `one eid-table` | Add or delete a ONE EID-table mapping. |
| `one eid-table map` | Map a ONE EID-table to a bridge-domain or VRF. |
| `one enable` | Enable the ONE control plane. |
| `one l2 arp` | Add or delete a static L2 ARP entry in the ONE overlay. |
| `one locator` | Add or delete a locator in a ONE locator-set. |
| `one locator-set` | Create or delete a named ONE locator-set. |
| `one map-register` | Enable or disable ONE map-register sending. |
| `one map-register fallback-threshold` | Set ONE map-register fallback threshold. |
| `one map-register ttl` | Set TTL for ONE map-register messages. |
| `one map-request itr-rlocs` | Set ONE map-request source RLOCs. |
| `one map-request mode` | Set ONE map-request mode. |
| `one map-resolver` | Add or delete a ONE map-resolver address. |
| `one map-server` | Add or delete a ONE map-server address. |
| `one ndp` | Add or delete a static NDP entry in the ONE overlay. |
| `one nsh-mapping` | Add or delete a ONE NSH-to-locator mapping. |
| `one petr mode` | Enable or disable ONE Proxy ETR mode. |
| `one pitr` | Enable or disable ONE Proxy ITR. |
| `one pitr mode` | Set ONE PITR mode parameters. |
| `one remote-mapping` | Add or delete a remote ONE map-cache entry. |
| `one rloc-probe` | Enable or disable ONE RLOC probing. |
| `one statistics` | Enable or disable ONE per-EID statistics. |
| `one statistics flush` | Flush all ONE traffic statistics. |
| `one use-petr` | Enable or disable ONE Proxy ETR usage. |
| `one xtr mode` | Enable or disable ONE xTR mode. |
| `p2p_ethernet` | Create or delete a point-to-point Ethernet sub-interface. |
| `packet-generator capture` | Enable or disable pcap capture on a pg interface. |
| `packet-generator configure` | Modify parameters of a packet-generator stream. |
| `packet-generator delete` | Remove a packet-generator stream. |
| `packet-generator disable-stream` | Stop a packet-generator stream. |
| `packet-generator enable-stream` | Start a packet-generator stream. |
| `packet-generator mac-filter` | Restrict pg capture to specific MAC addresses. |
| `packet-generator new` | Define a new packet-generator stream. |
| `pcap dispatch trace` | Enable per-node pcap capture in the dispatch graph. |
| `pcap trace` | Start or stop pcap capture on an interface. |
| `perfmon reset` | Reset all perfmon counters and discard statistics. |
| `perfmon start` | Start collecting hardware performance counters. |
| `perfmon stop` | Stop collecting hardware performance counters. |
| `ping` | Send ICMP echo requests from VPP. |
| `pipe create` | Create a software pipe interface pair. |
| `pipe delete` | Delete a software pipe interface pair. |
| `policer add` | Create a named policer with CIR/PIR and actions. |
| `policer bind` | Bind a policer to a worker thread. |
| `policer del` | Delete a named policer. |
| `policer input` | Apply a policer to an interface's input. |
| `policer output` | Apply a policer to an interface's output. |
| `policer reset` | Reset a policer's token buckets and counters. |
| `prom` | Enable or configure the Prometheus metrics exporter. |
| `prom patterns` | Set patterns to filter stats exposed via Prometheus. |
| `punt socket deregister` | Remove a punt socket binding. |
| `punt socket register` | Register a Unix domain socket to receive punted packets. |
| `pvti interface create` | Create a PVTI tunnel interface. |
| `pvti interface delete` | Delete a PVTI tunnel interface. |
| `q` | Alias for quit; close the CLI session. |
| `qos egress map` | Create or delete a QoS egress map. |
| `qos mark` | Enable or disable QoS marking on an interface. |
| `qos record` | Enable or disable QoS recording on an interface. |
| `qos store` | Set a fixed QoS value on packets entering an interface. |
| `quic set crypto api` | Select the crypto API for QUIC transport. |
| `quic set fifo-size` | Set default FIFO size for QUIC sessions. |
| `quit` | Close the current CLI session. |
| `renumber interface` | Change the sw_if_index of an interface. |
| `restart` | Restart the VPP process. |
| `sasc pcap start` | Start pcap capture on SASC traffic. |
| `sasc pcap stop` | Stop SASC pcap capture. |
| `save memory-trace` | Dump memory-trace log to a file. |
| `selog emit-elog` | Emit a manually-triggered structured event-log entry. |
| `session` | Enable or disable the session layer and set parameters. |
| `session replay fifo` | Replay a captured session FIFO trace. |
| `session rule` | Add or delete a session-layer rule. |
| `session sdl` | Add or delete a Session Descriptor List entry. |
| `set acl-plugin` | Set global ACL plugin parameters. |
| `set acl-plugin acl` | Create or replace an L3/L4 ACL rule set. |
| `set acl-plugin interface` | Bind or unbind an ACL to an interface. |
| `set acl-plugin macip acl` | Create or replace a MAC+IP ACL rule set. |
| `set acl-plugin macip interface` | Bind or unbind a MACIP ACL to an interface. |
| `set arp proxy` | Configure proxy-ARP address ranges for an interface. |
| `set bpf trace filter` | Attach a BPF program as a trace-capture filter. |
| `set bridge-domain arp entry` | Add or remove a static ARP entry in a bridge domain. |
| `set bridge-domain arp term` | Enable or disable ARP termination for a bridge domain. |
| `set bridge-domain arp-ufwd` | Enable or disable ARP unicast-forwarding for a bridge domain. |
| `set bridge-domain default-learn-limit` | Set default MAC learning limit for new bridge domains. |
| `set bridge-domain flood` | Enable or disable flooding for a bridge domain. |
| `set bridge-domain forward` | Enable or disable L2 forwarding for a bridge domain. |
| `set bridge-domain learn` | Enable or disable MAC learning for a bridge domain. |
| `set bridge-domain learn-limit` | Set MAC learning limit for a bridge domain. |
| `set bridge-domain mac-age` | Set MAC aging timeout for a bridge domain. |
| `set bridge-domain rewrite` | Set Ethernet rewrite template for a bridge domain. |
| `set bridge-domain uu-flood` | Set unknown-unicast flood interface for a bridge domain. |
| `set buffer traces` | Enable or disable buffer trace tracking. |
| `set clock adjust` | Adjust the VPP wall-clock. |
| `set cnat snat-policy` | Set the global CNAT source-NAT policy. |
| `set cnat snat-policy addr` | Add or remove CNAT source-NAT pool addresses. |
| `set cnat snat-policy if` | Mark an interface for CNAT SNAT inclusion/exclusion. |
| `set cnat snat-policy prefix` | Add or remove CNAT no-SNAT prefixes. |
| `set crypto handler` | Select the engine for a crypto algorithm. |
| `set ct6` | Enable or disable the IPv6 connection tracker on an interface. |
| `set det44 timeouts` | Set DET44 session timeouts per protocol. |
| `set dhcp client` | Start or stop a DHCPv4 client on an interface. |
| `set dhcp option-82 vss` | Set DHCP Option-82 VSS parameters. |
| `set dhcp proxy` | Configure a DHCPv4 relay-proxy server. |
| `set dhcpv6 proxy` | Configure a DHCPv6 relay-proxy server. |
| `set dhcpv6 vss` | Set DHCPv6 VSS parameters. |
| `set dpdk cryptodev assignment` | Assign DPDK crypto devices to workers. |
| `set dpdk interface descriptors` | Set RX/TX descriptor ring sizes for a DPDK interface. |
| `set fib walk histogram elements size` | Set FIB-walk histogram bucket size. |
| `set fib walk quota` | Set max FIB entries per walk iteration. |
| `set flow classify` | Attach a flow-classifier table to an interface. |
| `set flow-offload gtpu` | Enable or disable flow-offload for GTPU on an interface. |
| `set flow-offload vxlan` | Enable or disable flow-offload for VXLAN on an interface. |
| `set host-interface l4-cksum-offload` | Enable or disable L4 checksum offload on an AF_PACKET interface. |
| `set host-interface qdisc-bypass` | Enable or disable Linux qdisc bypass on an AF_PACKET interface. |
| `set ikev2 local key` | Set the local private key for IKEv2 authentication. |
| `set interface bond` | Configure bonding parameters on a member interface. |
| `set interface det44` | Designate an interface as DET44 inside or outside. |
| `set interface feature` | Enable or disable a feature-arc node on an interface. |
| `set interface feature gso` | Enable or disable GSO on an interface. |
| `set interface handoff` | Pin interface input to specific worker threads. |
| `set interface hw-class` | Change the hardware-class of an interface. |
| `set interface input acl` | Apply an input ACL to an interface. |
| `set interface ip address` | Add or delete an IP address on an interface. |
| `set interface ip directed-broadcast` | Enable or disable directed-broadcast on an interface. |
| `set interface ip geneve-bypass` | Enable or disable IPv4 Geneve bypass on an interface. |
| `set interface ip gtpu-bypass` | Enable or disable IPv4 GTPU bypass on an interface. |
| `set interface ip pvti-bypass` | Enable or disable IPv4 PVTI bypass on an interface. |
| `set interface ip source-and-port-range-check` | Attach source/port verification to an interface. |
| `set interface ip table` | Bind an interface to an IPv4 VRF table. |
| `set interface ip vxlan-bypass` | Enable or disable IPv4 VXLAN bypass on an interface. |
| `set interface ip vxlan-gpe-bypass` | Enable or disable IPv4 VXLAN-GPE bypass on an interface. |
| `set interface ip6 geneve-bypass` | Enable or disable IPv6 Geneve bypass on an interface. |
| `set interface ip6 gtpu-bypass` | Enable or disable IPv6 GTPU bypass on an interface. |
| `set interface ip6 l2tpv3` | Associate an L2TPv3 session with an interface. |
| `set interface ip6 pvti-bypass` | Enable or disable IPv6 PVTI bypass on an interface. |
| `set interface ip6 table` | Bind an interface to an IPv6 VRF table. |
| `set interface ip6 vxlan-bypass` | Enable or disable IPv6 VXLAN bypass on an interface. |
| `set interface ip6 vxlan-gpe-bypass` | Enable or disable IPv6 VXLAN-GPE bypass on an interface. |
| `set interface ip6-nd proxy` | Enable or disable IPv6 ND proxy on an interface. |
| `set interface ipsec spd` | Bind an IPsec SPD to an interface. |
| `set interface l2 bridge` | Place an interface into an L2 bridge domain. |
| `set interface l2 efp-filter` | Enable or disable EFP filtering on a bridge-domain member. |
| `set interface l2 flood` | Override per-interface L2 flood behaviour. |
| `set interface l2 forward` | Enable or disable L2 forwarding on a port. |
| `set interface l2 input classify` | Attach an L2 input classifier to an interface. |
| `set interface l2 learn` | Enable or disable MAC learning on a port. |
| `set interface l2 output classify` | Attach an L2 output classifier to an interface. |
| `set interface l2 pbb-tag-rewrite` | Configure 802.1ah PBB tag operations on an interface. |
| `set interface l2 rewrite` | Set Ethernet header rewrite rules on an L2 interface. |
| `set interface l2 tag-rewrite` | Configure VLAN tag push/pop/translate on a sub-interface. |
| `set interface l2 xconnect` | Create an L2 cross-connect between two interfaces. |
| `set interface l2 xcrw` | Set cross-connect rewrite on an L2-to-L3 path. |
| `set interface l3` | Return an interface from L2 bridge mode to L3 routed mode. |
| `set interface lldp` | Enable LLDP on an interface. |
| `set interface mac address` | Override the MAC address of an interface. |
| `set interface mpls` | Enable or disable MPLS on an interface. |
| `set interface mtu` | Set the MTU on an interface. |
| `set interface name` | Assign a user-friendly name to an interface. |
| `set interface nat44` | Designate an interface as NAT44-ED inside/outside. |
| `set interface nat44 ei` | Designate an interface as NAT44-EI inside/outside. |
| `set interface nat64` | Designate an interface as NAT64 inside/outside. |
| `set interface nat66` | Designate an interface as NAT66 inside/outside. |
| `set interface output acl` | Apply an output ACL to an interface. |
| `set interface promiscuous` | Enable or disable promiscuous mode on an interface. |
| `set interface proxy-arp` | Enable or disable proxy-ARP on an interface. |
| `set interface reassembly` | Enable or disable IP reassembly on an interface. |
| `set interface rss queues` | Configure RSS queue count and mapping on an interface. |
| `set interface rx-mode` | Set receive mode (polling/interrupt/adaptive) on an interface. |
| `set interface rx-placement` | Pin an interface rx-queue to a worker thread. |
| `set interface secondary-mac-address` | Add or remove secondary MAC addresses on an interface. |
| `set interface span` | Mirror traffic from one interface to another. |
| `set interface state` | Bring an interface administratively up or down. |
| `set interface tag` | Attach a descriptive tag to an interface. |
| `set interface tcp-mss-clamp` | Clamp TCP MSS in SYN packets on an interface. |
| `set interface tx-hash` | Select the TX hashing algorithm for multi-queue. |
| `set interface tx-queue` | Configure transmit queue parameters. |
| `set interface unnumbered` | Borrow the IP address of another interface. |
| `set ioam analyse` | Enable or disable iOAM trace analysis. |
| `set ioam export ipfix` | Configure IPFIX export of iOAM telemetry. |
| `set ioam ip6 cache` | Enable or disable iOAM IPv6 analysis cache. |
| `set ioam ip6 sr-tunnel-select` | Select the SRv6 tunnel for iOAM association. |
| `set ioam rewrite` | Configure iOAM hop-by-hop extension header rewrite. |
| `set ioam-trace profile` | Define an iOAM trace profile. |
| `set ip classify` | Attach a classify table to IPv4 input. |
| `set ip flow-hash` | Configure IPv4 ECMP flow-hash fields. |
| `set ip neighbor` | Add or modify a static ARP entry. |
| `set ip neighbor-config` | Set ARP table parameters. |
| `set ip source-and-port-range-check` | Configure source/port allowlists for a VRF. |
| `set ip6 address` | Add or delete an IPv6 address on an interface. |
| `set ip6 classify` | Attach a classify table to IPv6 input. |
| `set ip6 dad disable` | Disable IPv6 DAD on an interface. |
| `set ip6 dad enable` | Enable IPv6 DAD on an interface. |
| `set ip6 flow-hash` | Configure IPv6 ECMP flow-hash fields. |
| `set ip6 nd proxy` | Configure IPv6 ND proxy for a prefix. |
| `set ipfix classify stream` | Bind IPFIX export to a classifier table. |
| `set ipfix exporter` | Configure IPFIX exporter (collector, source, template). |
| `set ipsec async mode` | Enable or disable async crypto for IPsec. |
| `set l2fib scan-delay` | Set L2 FIB aging scan interval. |
| `set l2tpv3 tunnel cookie` | Set authentication cookie on an L2TPv3 session. |
| `set lldp` | Set global LLDP parameters. |
| `set logging class` | Set log severity for a log class. |
| `set logging size` | Set log ring buffer size. |
| `set logging unthrottle-time` | Set minimum interval between repeated log messages. |
| `set nat frame-queue-nelts` | Set NAT inter-worker frame-queue size. |
| `set nat timeout` | Set global NAT session timeouts. |
| `set nat workers` | Assign worker threads for NAT processing. |
| `set nat44 ei timeout` | Set NAT44-EI session timeouts. |
| `set nat44 ei workers` | Assign workers for NAT44-EI processing. |
| `set nat44 session limit` | Set max concurrent NAT44-ED sessions. |
| `set node function` | Select an alternative node dispatch function. |
| `set nsh-md2-ioam export ipfix` | Configure IPFIX export for NSH iOAM. |
| `set nsh-md2-ioam-transit` | Enable iOAM transit processing for NSH paths. |
| `set nsim` | Set network delay simulator parameters. |
| `set pcap filter function` | Select pcap capture filter function. |
| `set pnat translation` | Add or delete a PNAT translation rule. |
| `set policer classify` | Attach a policer-classifier to an interface. |
| `set pot profile` | Configure a Proof-of-Transit profile. |
| `set pot profile-active` | Select the active PoT profile index. |
| `set punt` | Configure punt-to-host behaviour. |
| `set quic cc` | Select QUIC congestion-control algorithm. |
| `set quic max_packets_per_key` | Set max packets per QUIC key before rotation. |
| `set sasc ingress interface` | Designate an SASC ingress interface. |
| `set sasc services` | Register SASC services for a tenant. |
| `set sasc tenant` | Create or modify an SASC tenant. |
| `set sasc timeout` | Set SASC session timeout. |
| `set sfdp acl` | Attach an ACL to SFDP classification. |
| `set sfdp eviction sessions-margin` | Set SFDP session eviction margin. |
| `set sfdp gateway geneve-output` | Configure SFDP GENEVE output tunnel. |
| `set sfdp icmp-error-node` | Set SFDP ICMP error graph node. |
| `set sfdp interface-input` | Designate an SFDP input interface. |
| `set sfdp nat external-interface` | Set SFDP NAT external interface. |
| `set sfdp nat snat` | Configure SFDP source-NAT. |
| `set sfdp services` | Register SFDP services. |
| `set sfdp sp-node` | Set SFDP service-policy node. |
| `set sfdp timeout` | Set SFDP session timeout. |
| `set sr encaps hop-limit` | Set SRv6 encapsulation hop-limit. |
| `set sr encaps source` | Set SRv6 encapsulation source address. |
| `set sw_scheduler` | Configure software crypto-scheduler. |
| `set syslog filter` | Set syslog export severity filter. |
| `set syslog sender` | Configure syslog exporter parameters. |
| `set terminal ansi` | Enable or disable ANSI output in CLI. |
| `set terminal history` | Set CLI history depth. |
| `set terminal pager` | Configure CLI output pager. |
| `set trace filter function` | Select trace capture filter function. |
| `set trace timestamp-format` | Choose trace timestamp format. |
| `set udp-ping` | Configure UDP-ping probes. |
| `set udp-ping export-ipfix` | Enable or disable IPFIX export of UDP-ping results. |
| `set urpf` | Enable uRPF (strict/loose) on an interface. |
| `set urpf-accept` | Configure uRPF accept-list exemptions. |
| `set virtio pci` | Set VirtIO PCI device parameters. |
| `set vxlan-gpe-ioam` | Enable or disable iOAM for VXLAN-GPE. |
| `set vxlan-gpe-ioam export ipfix` | Configure IPFIX export for VXLAN-GPE iOAM. |
| `set vxlan-gpe-ioam rewrite` | Configure iOAM rewrite for VXLAN-GPE. |
| `set vxlan-gpe-ioam-transit` | Enable VXLAN-GPE iOAM transit processing. |
| `set wireguard async mode` | Enable or disable async crypto for WireGuard. |
| `sfdp gateway geneve-input` | Configure SFDP GENEVE input for gateway decap. |
| `sfdp nat alloc-pool` | Allocate an address pool for SFDP NAT. |
| `sfdp snort create-instance` | Create an SFDP-managed Snort instance. |
| `sfdp snort delete-instance` | Delete an SFDP-managed Snort instance. |
| `sfdp tenant` | Create, modify, or delete an SFDP tenant. |
| `sflow direction` | Set sFlow sampling direction on an interface. |
| `sflow drop-monitoring` | Enable or disable sFlow drop-reason monitoring. |
| `sflow enable-disable` | Enable or disable sFlow sampling. |
| `sflow header-bytes` | Set sFlow sample header capture size. |
| `sflow polling-interval` | Set sFlow counter-polling interval. |
| `sflow sampling-rate` | Set sFlow packet sampling rate (1-in-N). |
| `show abf attach` | Display ABF policy attachments on interfaces. |
| `show abf policy` | Display configured ABF policies and their forwarding paths. |
| `show acl-plugin acl` | Display L3/L4 ACL rule sets managed by the ACL plugin. |
| `show acl-plugin decode 5tuple` | Decode a raw 5-tuple match value into human-readable ACL fields. |
| `show acl-plugin interface` | Display ACL bindings on interfaces. |
| `show acl-plugin lookup context` | Display ACL lookup context state. |
| `show acl-plugin lookup user` | Display ACL lookup user registrations. |
| `show acl-plugin macip acl` | Display MAC+IP ACL rule sets. |
| `show acl-plugin macip interface` | Display MACIP ACL interface bindings. |
| `show acl-plugin memory` | Display ACL plugin memory usage. |
| `show acl-plugin sessions` | Display ACL connection-tracking session state. |
| `show acl-plugin tables` | Display ACL plugin internal hash and match tables. |
| `show adj` | Display adjacency table entries. |
| `show adj nbr` | Display neighbor adjacency entries. |
| `show api clients` | Display connected binary API clients. |
| `show api dump` | Dump API message definitions or compare with a baseline. |
| `show api histogram` | Display API message processing latency histogram. |
| `show api message-table` | Display the binary API message name-to-ID table. |
| `show api plugin` | Display API plugin registrations. |
| `show api ring-stats` | Display shared-memory API ring statistics. |
| `show api trace-status` | Display API trace enable/disable status. |
| `show app` | Display host-stack application registrations and sessions. |
| `show app certificate` | List TLS certificates and keys in the application store. |
| `show app evt-collector` | Display application event-collector configuration. |
| `show app ns` | Display application namespace definitions. |
| `show app tls-profile` | Display registered TLS profiles. |
| `show arp proxy` | Display proxy-ARP address ranges. |
| `show auto-sdl` | Display automatic Session Descriptor List state. |
| `show bfd` | Display BFD keys, sessions, and echo-source state. |
| `show bier bift` | Display BIER Bit Index Forwarding Table entries. |
| `show bier disp entry` | Display BIER disposition entry details. |
| `show bier disp table` | Display BIER disposition tables. |
| `show bier fib` | Display BIER forwarding information base entries. |
| `show bier fmask` | Display BIER forwarding masks. |
| `show bier imp` | Display BIER imposition objects. |
| `show bihash` | Display bounded-index hash table statistics. |
| `show bond` | Display bonding group state and member interfaces. |
| `show bpf trace filter` | Display the active BPF trace filter program. |
| `show bridge-domain` | Display bridge-domain state, members, and ARP tables. |
| `show buffer metadata` | Display per-buffer metadata tracking state. |
| `show buffer traces` | Display buffer-monitoring trace data. |
| `show buffers` | Display packet buffer pool allocation statistics. |
| `show cdp` | Display Cisco Discovery Protocol neighbor information. |
| `show classify filter` | Display active classify filters for pcap/trace. |
| `show classify flow` | Display flow-classifier table bindings. |
| `show classify policer` | Display policer-classifier table bindings. |
| `show classify tables` | Display classify table definitions and match entries. |
| `show cli` | Display CLI command tree and MP-safety annotations. |
| `show cli-sessions` | Display active CLI sessions and their terminal state. |
| `show clock` | Display VPP wall-clock time. |
| `show cnat client` | Display CNAT client endpoint registrations. |
| `show cnat session` | Display active CNAT translation sessions. |
| `show cnat snat-policy` | Display CNAT source-NAT policy configuration. |
| `show cnat timestamp` | Display CNAT session timestamp state. |
| `show cnat translation` | Display CNAT translation rules and backends. |
| `show cpu` | Display CPU topology, features, and thread affinity. |
| `show crypto algorithm` | Display registered cryptographic algorithms. |
| `show crypto async status` | Display async crypto dispatch queue status. |
| `show crypto engines` | Display registered crypto engines and their priorities. |
| `show crypto handlers` | Display per-algorithm crypto handler assignments. |
| `show det44 interfaces` | Display DET44 inside/outside interface bindings. |
| `show det44 mappings` | Display DET44 deterministic address mapping ranges. |
| `show det44 sessions` | Display active DET44 sessions. |
| `show det44 timeouts` | Display DET44 session timeout values. |
| `show device` | Display device state, counters, and debug information. |
| `show device drivers` | Display registered device drivers. |
| `show dhcp client` | Display DHCPv4 client state on interfaces. |
| `show dhcp option-82-address interface` | Display DHCP Option-82 address on an interface. |
| `show dhcp proxy` | Display DHCPv4 relay-proxy server configuration. |
| `show dhcp vss` | Display DHCP Virtual Subnet Selection parameters. |
| `show dhcp6 addresses` | Display DHCPv6-acquired addresses. |
| `show dhcp6 clients` | Display DHCPv6 client state. |
| `show dhcpv6 link-address interface` | Display DHCPv6 link-address on an interface. |
| `show dhcpv6 proxy` | Display DHCPv6 relay-proxy configuration. |
| `show dhcpv6 vss` | Display DHCPv6 VSS parameters. |
| `show dma` | Display DMA engine configuration and status. |
| `show dma backends` | Display registered DMA backends. |
| `show dns cache` | Display DNS resolver cache entries. |
| `show dns servers` | Display configured upstream DNS servers. |
| `show dpdk buffer` | Display DPDK buffer pool statistics. |
| `show dpdk cryptodev assignment` | Display DPDK cryptodev-to-worker assignments. |
| `show dpdk cryptodev cache status` | Display DPDK cryptodev cache ring status. |
| `show dpdk cryptodev capabilities` | Display DPDK cryptodev algorithm capabilities. |
| `show dpdk physmem` | Display DPDK physical memory (hugepage) layout. |
| `show dpdk version` | Display DPDK library version. |
| `show dpo memory` | Display DPO (Data-Path Object) memory consumption. |
| `show dslite aftr-tunnel-endpoint-address` | Display DS-Lite AFTR tunnel endpoint address. |
| `show dslite b4-tunnel-endpoint-address` | Display DS-Lite B4 tunnel endpoint address. |
| `show dslite pool` | Display DS-Lite NAT address pool. |
| `show dslite sessions` | Display active DS-Lite translation sessions. |
| `show errors` | Display per-node packet-processing error counters. |
| `show event-logger` | Display event-logger buffer size and status. |
| `show features` | Display feature-arc registrations and ordering. |
| `show fib entry` | Display FIB entry details by index. |
| `show fib entry-delegate` | Display FIB entry delegate chains. |
| `show fib memory` | Display FIB memory consumption. |
| `show fib path-lists` | Display FIB path-list objects. |
| `show fib paths` | Display FIB path resolution state. |
| `show fib source` | Display FIB source priority table. |
| `show fib uRPF` | Display FIB uRPF (Unicast RPF) lists. |
| `show fib walk` | Display FIB walk statistics and histogram. |
| `show files` | Display open file descriptors managed by VPP. |
| `show flow entry` | Display hardware-offloaded flow entries. |
| `show flow interface` | Display flow rule bindings on an interface. |
| `show flow ranges` | Display configured flow ID ranges. |
| `show flowprobe feature` | Display flowprobe feature enablement on interfaces. |
| `show flowprobe params` | Display flowprobe timeout and record format parameters. |
| `show flowprobe statistics` | Display flowprobe packet/flow export statistics. |
| `show flowprobe table` | Display flowprobe flow table entries. |
| `show frame-queue` | Display inter-worker frame-queue state and trace. |
| `show frame-queue histogram` | Display frame-queue depth histogram. |
| `show gdb` | List GDB helper functions callable from the debugger. |
| `show geneve tunnel` | Display GENEVE tunnel endpoint state. |
| `show gpe adjacency` | Display GPE adjacency entries. |
| `show gpe encap` | Display GPE encapsulation mode. |
| `show gpe entry` | Display GPE forwarding entries. |
| `show gpe interface` | Display GPE interfaces. |
| `show gpe native-forward` | Display GPE native-forward entries. |
| `show gpe sub-interface` | Display GPE sub-interfaces. |
| `show gpe tenant` | Display GPE tenant state. |
| `show gpe tunnel` | Display GPE tunnel endpoints. |
| `show graph` | Display packet-processing graph node topology. |
| `show gre tunnel` | Display GRE tunnel endpoint state. |
| `show gtpu tunnel` | Display GTP-U tunnel endpoint state. |
| `show hardware-interfaces` | Display hardware interface state, counters, and driver info. |
| `show hash` | Display registered hash algorithms. |
| `show http connect proxy client` | Display HTTP CONNECT proxy client state. |
| `show http static server` | Display HTTP static file-server state and cache. |
| `show http stats` | Display HTTP layer statistics. |
| `show http tps` | Display HTTP TPS benchmark application state. |
| `show igmp config` | Display IGMP configuration on interfaces. |
| `show igmp ssm-ranges` | Display IGMP Source-Specific Multicast address ranges. |
| `show igmp timers` | Display IGMP protocol timer values. |
| `show ikev2 profile` | Display IKEv2 profile configuration. |
| `show ikev2 sa` | Display IKEv2 Security Association state. |
| `show ikev2 sleep interval` | Display IKEv2 background process sleep interval. |
| `show ila entries` | Display ILA translation entries. |
| `show inacl` | Display input ACL bindings by type. |
| `show init-function` | Display init/enter/exit function registration and execution state. |
| `show interface` | Display interface state, addresses, features, and counters. |
| `show interface rx-placement` | Display interface rx-queue to worker-thread placement. |
| `show interface secondary-mac-address` | Display secondary MAC addresses on interfaces. |
| `show interface span` | Display SPAN (port mirror) configuration. |
| `show interface tcp-mss-clamp` | Display TCP MSS clamping configuration on interfaces. |
| `show interface transceiver` | Display SFP/QSFP transceiver module information. |
| `show interface tx-hash` | Display transmit hash algorithm for interfaces. |
| `show ioam analyse` | Display iOAM trace analysis results. |
| `show ioam e2e` | Display iOAM edge-to-edge measurement data. |
| `show ioam ip6 cache` | Display iOAM IPv6 analysis cache entries. |
| `show ioam nsh-lisp-gpe trace` | Display iOAM trace data for NSH-LISP-GPE paths. |
| `show ioam pot` | Display iOAM Proof-of-Transit statistics. |
| `show ioam summary` | Display iOAM operational summary. |
| `show ioam trace` | Display iOAM hop-by-hop trace statistics. |
| `show ioam vxlan-gpe trace` | Display iOAM trace data for VXLAN-GPE tunnels. |
| `show ioam-trace profile` | Display configured iOAM trace profile. |
| `show ip container` | Display IP container-proxy entries. |
| `show ip fib` | Display IPv4 FIB table entries and MTRIE structure. |
| `show ip local` | Display IPv4 local (punt) delivery table. |
| `show ip mfib` | Display IPv4 multicast FIB entries. |
| `show ip neighbor` | Display IPv4/IPv6 neighbor (ARP/ND) entries. |
| `show ip neighbor-config` | Display ARP/ND table configuration parameters. |
| `show ip neighbor-stats` | Display neighbor table statistics. |
| `show ip neighbor-watcher` | Display registered neighbor-change watchers. |
| `show ip neighbors` | Display all IP neighbor entries (alias). |
| `show ip pmtu` | Display IP Path MTU cache entries. |
| `show ip punt redirect` | Display IPv4 punt redirect configuration. |
| `show ip session redirect` | Display IP session redirect rules. |
| `show ip source-and-port-range-check` | Display source/port range check configuration. |
| `show ip table` | Display IPv4 VRF/FIB table info. |
| `show ip4 neighbor` | Display IPv4 ARP neighbor entries. |
| `show ip4 neighbor-sorted` | Display IPv4 neighbors sorted by address. |
| `show ip4 neighbors` | Display IPv4 neighbor entries (alias). |
| `show ip4-full-reassembly` | Display IPv4 full reassembly state. |
| `show ip4-sv-reassembly` | Display IPv4 shallow-virtual reassembly state. |
| `show ip6 addresses` | Display DHCPv6-acquired IPv6 addresses. |
| `show ip6 connection-tracker` | Display IPv6 connection-tracker (ct6) state. |
| `show ip6 dad` | Display IPv6 Duplicate Address Detection state. |
| `show ip6 fib` | Display IPv6 FIB table entries. |
| `show ip6 hbh` | Display IPv6 hop-by-hop extension header handlers. |
| `show ip6 interface` | Display IPv6 interface state (link-local, RA). |
| `show ip6 local` | Display IPv6 local (punt) delivery table. |
| `show ip6 mfib` | Display IPv6 multicast FIB entries. |
| `show ip6 neighbor` | Display IPv6 ND neighbor entries. |
| `show ip6 neighbor-sorted` | Display IPv6 neighbors sorted by address. |
| `show ip6 neighbors` | Display IPv6 neighbor entries (alias). |
| `show ip6 pd clients` | Display IPv6 Prefix Delegation client state. |
| `show ip6 prefixes` | Display IPv6 delegated prefixes. |
| `show ip6 punt redirect` | Display IPv6 punt redirect configuration. |
| `show ip6 table` | Display IPv6 VRF/FIB table info. |
| `show ip6-full-reassembly` | Display IPv6 full reassembly state. |
| `show ip6-ll` | Display IPv6 link-local FIB table. |
| `show ip6-sv-reassembly` | Display IPv6 shallow-virtual reassembly state. |
| `show ipip tunnel` | Display IP-in-IP tunnel endpoint state. |
| `show ipip tunnel-hash` | Display IPIP tunnel hash table. |
| `show ipsec all` | Display all IPsec state (SAs, SPDs, tunnel-protect). |
| `show ipsec interface` | Display IPsec interface (ipsecN) state. |
| `show ipsec protect` | Display IPsec tunnel-protect bindings. |
| `show ipsec protect-hash` | Display IPsec protect hash table. |
| `show ipsec sa` | Display IPsec Security Association details. |
| `show ipsec spd` | Display IPsec Security Policy Database entries. |
| `show ipsec tunnel` | Display IPsec tunnel state. |
| `show l2 rewrite entries` | Display L2 output-rewrite entries. |
| `show l2 rewrite interfaces` | Display interfaces with L2 rewrite rules. |
| `show l2fib` | Display L2 FIB (MAC forwarding table) entries. |
| `show l2patch` | Display L2 patch (cross-connect) entries. |
| `show l2tpv3` | Display L2TPv3 tunnel sessions. |
| `show l2xcrw` | Display L2 cross-connect rewrite entries. |
| `show l3xc` | Display L3 cross-connect rules. |
| `show lacp` | Display LACP protocol state on bonded interfaces. |
| `show lb` | Display load-balancer global state. |
| `show lb vips` | Display load-balancer VIPs and backends. |
| `show lisp adjacencies` | Display LISP adjacency state. |
| `show lisp eid-table` | Display LISP EID-table mappings. |
| `show lisp eid-table map` | Display LISP EID-table to VRF/BD mappings. |
| `show lisp locator-set` | Display LISP locator-set definitions. |
| `show lisp map-request itr-rlocs` | Display LISP map-request ITR-RLOC configuration. |
| `show lisp map-request mode` | Display LISP map-request mode. |
| `show lisp map-resolvers` | Display LISP map-resolver addresses. |
| `show lisp petr` | Display LISP Proxy ETR configuration. |
| `show lisp pitr` | Display LISP Proxy ITR configuration. |
| `show lisp status` | Display LISP enable/disable status. |
| `show lldp` | Display LLDP neighbor and configuration state. |
| `show load-balance` | Display load-balance DPO objects. |
| `show load-balance-map` | Display load-balance map objects. |
| `show logging` | Display log buffer contents. |
| `show logging configuration` | Display per-class log severity configuration. |
| `show lookup-dpo` | Display lookup DPO objects. |
| `show macro` | Display defined CLI macros. |
| `show mactime` | Display MAC-time access-control entries and counters. |
| `show map domain` | Display MAP-E/MAP-T domain configurations. |
| `show map stats` | Display MAP translation statistics. |
| `show memif` | Display memif interface state and descriptors. |
| `show memory` | Display heap memory usage for main, API, and stats segments. |
| `show mfib entry` | Display multicast FIB entry details. |
| `show mfib interface` | Display multicast FIB interface state. |
| `show mfib itf flags` | Display MFIB interface flag definitions. |
| `show mfib route flags` | Display MFIB route flag definitions. |
| `show mode` | Display L2/L3 mode of interfaces. |
| `show mpls fib` | Display MPLS FIB label table. |
| `show mpls interface` | Display MPLS-enabled interfaces. |
| `show mpls tunnel` | Display MPLS tunnel state. |
| `show nat mss-clamping` | Display NAT TCP MSS clamping value. |
| `show nat timeouts` | Display global NAT session timeouts. |
| `show nat workers` | Display NAT worker thread assignments. |
| `show nat44 addresses` | Display NAT44-ED external address pool. |
| `show nat44 ei addr-port-assignment-alg` | Display NAT44-EI address/port assignment algorithm. |
| `show nat44 ei addresses` | Display NAT44-EI external address pool. |
| `show nat44 ei ha` | Display NAT44-EI high-availability state. |
| `show nat44 ei hash tables` | Display NAT44-EI hash table statistics. |
| `show nat44 ei interface address` | Display NAT44-EI interface-sourced pool addresses. |
| `show nat44 ei interfaces` | Display NAT44-EI inside/outside interfaces. |
| `show nat44 ei mss-clamping` | Display NAT44-EI MSS clamping value. |
| `show nat44 ei sessions` | Display active NAT44-EI sessions. |
| `show nat44 ei static mappings` | Display NAT44-EI static mapping rules. |
| `show nat44 ei timeouts` | Display NAT44-EI session timeouts. |
| `show nat44 ei workers` | Display NAT44-EI worker assignments. |
| `show nat44 hash tables` | Display NAT44-ED hash table statistics. |
| `show nat44 interface address` | Display NAT44-ED interface-sourced pool addresses. |
| `show nat44 interfaces` | Display NAT44-ED inside/outside interfaces. |
| `show nat44 sessions` | Display active NAT44-ED sessions. |
| `show nat44 static mappings` | Display NAT44-ED static mapping rules. |
| `show nat44 summary` | Display NAT44-ED summary (pool sizes, session counts). |
| `show nat44 vrf tables` | Display NAT44-ED VRF routing tables. |
| `show nat64 bib` | Display NAT64 Binding Information Base entries. |
| `show nat64 interfaces` | Display NAT64 inside/outside interfaces. |
| `show nat64 pool` | Display NAT64 address pool. |
| `show nat64 prefix` | Display NAT64 IPv6 prefixes. |
| `show nat64 session table` | Display active NAT64 sessions. |
| `show nat66 interfaces` | Display NAT66 inside/outside interfaces. |
| `show nat66 static mappings` | Display NAT66 static mapping rules. |
| `show node` | Display graph node details by name or index. |
| `show node counters` | Display per-node packet counters. |
| `show npol interfaces` | Display network-policy interface bindings. |
| `show npol ipsets` | Display configured IP-sets. |
| `show npol policies` | Display configured network-policies. |
| `show npol rules` | Display network-policy match rules. |
| `show nsh entry` | Display NSH entry definitions. |
| `show nsh map` | Display NSH forwarding map entries. |
| `show nsim` | Display network delay simulator configuration. |
| `show one adjacencies` | Display ONE overlay adjacency state. |
| `show one eid-table` | Display ONE EID-table mappings. |
| `show one eid-table map` | Display ONE EID-table to VRF/BD mappings. |
| `show one l2 arp entries` | Display ONE L2 ARP entries. |
| `show one locator-set` | Display ONE locator-set definitions. |
| `show one map-register fallback-threshold` | Display ONE map-register fallback threshold. |
| `show one map-register state` | Display ONE map-register state. |
| `show one map-register ttl` | Display ONE map-register TTL. |
| `show one map-request itr-rlocs` | Display ONE map-request ITR-RLOC configuration. |
| `show one map-request mode` | Display ONE map-request mode. |
| `show one map-resolvers` | Display ONE map-resolver addresses. |
| `show one map-servers` | Display ONE map-server addresses. |
| `show one modes` | Display ONE operational modes. |
| `show one ndp entries` | Display ONE NDP entries. |
| `show one petr` | Display ONE Proxy ETR configuration. |
| `show one pitr` | Display ONE Proxy ITR configuration. |
| `show one rloc state` | Display ONE RLOC reachability state. |
| `show one statistics details` | Display ONE per-EID traffic statistics. |
| `show one statistics status` | Display ONE statistics collection status. |
| `show one status` | Display ONE enable/disable status. |
| `show outacl` | Display output ACL bindings by type. |
| `show packet-generator` | Display packet-generator stream definitions. |
| `show packet-generator interface` | Display packet-generator interface state. |
| `show pcap filter function` | Display active pcap filter function. |
| `show pci` | Display PCI device information. |
| `show perfmon active-bundle` | Display the active perfmon counter bundle. |
| `show perfmon bundle` | Display available perfmon bundles. |
| `show perfmon source` | Display available perfmon counter sources. |
| `show perfmon statistics` | Display collected perfmon statistics. |
| `show physmem` | Display physical memory (hugepage) allocation. |
| `show plugins` | Display loaded VPP plugins. |
| `show pnat interfaces` | Display PNAT interface bindings. |
| `show pnat translations` | Display PNAT translation rules. |
| `show policer` | Display policer configuration and token-bucket state. |
| `show policer pools` | Display policer pool allocations. |
| `show pot profile` | Display Proof-of-Transit profile configuration. |
| `show pppoe fib` | Display PPPoE FIB entries. |
| `show pppoe session` | Display PPPoE session state. |
| `show punt client` | Display registered punt clients. |
| `show punt db` | Display punt database entries. |
| `show punt reasons` | Display all punt reason registrations. |
| `show punt socket registrations` | Display punt socket bindings. |
| `show punt stats` | Display punt statistics. |
| `show pvti interface` | Display PVTI interface state. |
| `show pvti rx peers` | Display PVTI receive-side peer state. |
| `show pvti tx peers` | Display PVTI transmit-side peer state. |
| `show qos egress map` | Display QoS egress map definitions. |
| `show qos mark` | Display QoS marking configuration on interfaces. |
| `show qos record` | Display QoS recording configuration on interfaces. |
| `show qos store` | Display QoS store configuration on interfaces. |
| `show quic` | Display QUIC transport state. |
| `show quic crypto context` | Display QUIC crypto context state. |
| `show replicate` | Display replicate DPO objects. |
| `show runtime` | Display per-node runtime statistics (clocks, vectors, calls). |
| `show sasc ingress interfaces` | Display SASC ingress interface bindings. |
| `show sasc next-indices` | Display SASC graph-node next-index mappings. |
| `show sasc pcap` | Display SASC pcap capture status. |
| `show sasc schema` | Display SASC export schema. |
| `show sasc service-chains` | Display SASC service-chain definitions. |
| `show sasc services` | Display registered SASC services. |
| `show sasc session` | Display SASC session state. |
| `show sasc summary` | Display SASC operational summary. |
| `show sasc tenant` | Display SASC tenant definitions. |
| `show segment-manager` | Display session segment-manager state. |
| `show session` | Display host-stack session state and statistics. |
| `show session fifo trace` | Display session FIFO trace data. |
| `show session lookup` | Display session-layer lookup table. |
| `show session rules` | Display session-layer rules (permit/deny/redirect). |
| `show session sdl` | Display Session Descriptor List entries. |
| `show session stats` | Display session-layer aggregate statistics. |
| `show sfdp services` | Display SFDP registered services. |
| `show sfdp session-detail` | Display SFDP session detail by session-id. |
| `show sfdp session-table` | Display SFDP session table entries. |
| `show sfdp status` | Display SFDP operational status. |
| `show sfdp tcp session-table` | Display SFDP TCP session table entries. |
| `show sfdp tenant` | Display SFDP tenant definitions. |
| `show sflow` | Display sFlow sampling configuration and statistics. |
| `show snort clients` | Display connected Snort clients. |
| `show snort instances` | Display Snort instance state. |
| `show snort interfaces` | Display Snort interface attachments. |
| `show soft-rss` | Display software RSS state on interfaces. |
| `show sr encaps hop-limit` | Display SRv6 encapsulation hop-limit. |
| `show sr encaps source addr` | Display SRv6 encapsulation source address. |
| `show sr localsids` | Display SRv6 local-SID entries. |
| `show sr localsids behaviors` | Display available SRv6 local-SID behaviours. |
| `show sr mpls policies` | Display SR-MPLS policy state. |
| `show sr mpls steering policies` | Display SR-MPLS steering rules. |
| `show sr policies` | Display SRv6 policy state. |
| `show sr policy behaviors` | Display available SRv6 policy behaviours. |
| `show sr steering-policies` | Display SRv6 steering rules. |
| `show statistics hash` | Display statistics segment hash table info. |
| `show statistics segment` | Display statistics segment counters. |
| `show stn rules` | Display Steal-The-NIC rules. |
| `show svs` | Display Source VRF Select state. |
| `show sw_scheduler workers` | Display software crypto-scheduler worker assignments. |
| `show syslog filter` | Display syslog severity filter. |
| `show syslog sender` | Display syslog exporter configuration. |
| `show tap` | Display TUN/TAP interface state. |
| `show tcp config` | Display TCP transport configuration. |
| `show tcp punt` | Display TCP punt configuration. |
| `show tcp scoreboard trace` | Display TCP SACK scoreboard trace. |
| `show tcp stats` | Display TCP transport statistics. |
| `show teib` | Display Tunnel Endpoint Information Base entries. |
| `show terminal` | Display current CLI terminal settings. |
| `show test tls server` | Display built-in TLS test server state. |
| `show threads` | Display VPP worker thread topology. |
| `show trace` | Display captured packet trace buffer. |
| `show trace filter function` | Display active trace filter function. |
| `show trace timestamp-format` | Display trace timestamp format setting. |
| `show tun` | Display TUN interface state. |
| `show udp encap` | Display UDP encapsulation entries. |
| `show udp ports` | Display UDP port bindings. |
| `show udp punt` | Display UDP punt configuration. |
| `show udp transport ports` | Display UDP transport port registrations. |
| `show udp-ping summary` | Display UDP-ping measurement summary. |
| `show unix errors` | Display Unix system-call error history. |
| `show version` | Display VPP version, build info, and command line. |
| `show vhost-user` | Display vhost-user interface state and descriptors. |
| `show virtio pci` | Display VirtIO PCI interface state and rings. |
| `show vlib frame-allocation` | Display node dispatch frame allocation statistics. |
| `show vlib graph` | Display packet-processing node graph. |
| `show vlib graphviz` | Dump node graph as a Graphviz dotfile. |
| `show vmxnet3` | Display vmxnet3 interface state and descriptor rings. |
| `show vrrp vr` | Display VRRP virtual router state. |
| `show vxlan tunnel` | Display VXLAN tunnel endpoint state. |
| `show vxlan-gpe` | Display VXLAN-GPE tunnel state. |
| `show wireguard interface` | Display WireGuard interface state. |
| `show wireguard mode` | Display WireGuard async/sync mode. |
| `show wireguard peer` | Display WireGuard peer state. |
| `snort attach` | Attach a Snort instance to an interface. |
| `snort create-instance` | Create a named Snort instance. |
| `snort delete instance` | Delete a named Snort instance. |
| `snort detach` | Detach a Snort instance from an interface. |
| `snort disconnect client` | Force-disconnect a Snort client. |
| `soft-rss clear` | Clear software RSS statistics. |
| `soft-rss config` | Configure software RSS parameters. |
| `soft-rss disable` | Disable software RSS on an interface. |
| `soft-rss enable` | Enable software RSS on an interface. |
| `sr localsid` | Create or delete an SRv6 local-SID. |
| `sr mpls policy` | Create or delete an SR-MPLS policy. |
| `sr mpls policy te` | Create or delete an SR-MPLS TE policy. |
| `sr mpls steer` | Create or delete an SR-MPLS steering rule. |
| `sr policy` | Create, modify, or delete an SRv6 policy. |
| `sr pt add iface` | Add an interface to SRv6 Path Tracing. |
| `sr pt del iface` | Remove an interface from SRv6 Path Tracing. |
| `sr pt show iface` | Display SRv6 Path Tracing interface config. |
| `sr steer` | Create or delete an SRv6 steering rule. |
| `stn rule` | Add or delete a Steal-The-NIC rule. |
| `suspend` | Suspend the CLI process context and yield. |
| `svs enable` | Enable or disable Source VRF Select on an interface. |
| `svs route` | Add or delete an SVS source-prefix route. |
| `svs table` | Create or delete an SVS table. |
| `tcp debug` | Set TCP debug verbosity or toggle features. |
| `tcp replay scoreboard` | Replay a TCP SACK scoreboard trace. |
| `tcp src-address` | Set preferred source address for outgoing TCP. |
| `test cnat maglev` | Validate CNAT Maglev consistent-hash algorithm. |
| `test cnat scanner` | Trigger manual CNAT session scanner run. |
| `test dma` | Exercise DMA subsystem with test operations. |
| `test dns expire` | Force-expire DNS cache entries. |
| `test dns format` | Encode a DNS name into wire format. |
| `test dns unformat` | Decode wire-format DNS name to dotted notation. |
| `test dpdk buffer` | Validate DPDK buffer pool consistency. |
| `test echo clients` | Launch built-in echo clients for throughput testing. |
| `test echo server` | Start built-in echo server for data reflection. |
| `test fib-walk-process` | Trigger a manual FIB convergence walk. |
| `test frame-queue nelts` | Change inter-worker frame-queue element count. |
| `test frame-queue threshold` | Adjust frame-queue congestion threshold. |
| `test igmp timers` | Override IGMP timer values for testing. |
| `test ip checksum` | Verify IPv4 checksum computation routines. |
| `test ip6 connection-tracker` | Exercise IPv6 connection-tracker with synthetic state. |
| `test ip6 link` | Validate IPv6 link-local address and ND state machine. |
| `test l2fib` | Stress-test L2 FIB hash table operations. |
| `test l2patch` | Test L2 patch cross-connect operations. |
| `test lb flowtable flush` | Flush load-balancer per-flow hash table. |
| `test log` | Emit a synthetic log message for validation. |
| `test lt2p counters` | Validate L2TPv3 counter behaviour. |
| `test one nsh` | Exercise ONE NSH encap/decap paths. |
| `test one nsh add-placeholder-decap-node` | Register placeholder NSH decap node for testing. |
| `test proxy server` | Start built-in TCP proxy server. |
| `test rdma dump` | Dump RDMA (mlx5) driver internal state. |
| `test sasc list` | List available SASC test cases. |
| `test sasc run` | Execute a single SASC test case. |
| `test sasc run-all` | Execute all SASC test cases. |
| `test sasc run-id` | Execute SASC test case by run-id. |
| `test sfdp expiry disable` | Disable SFDP session-expiry timer. |
| `test sfdp expiry enable` | Re-enable SFDP session-expiry timer. |
| `test syslog` | Generate a test syslog message. |
| `test tls client` | Launch built-in TLS client. |
| `test tls server` | Start built-in TLS server. |
| `test vmxnet3` | Exercise vmxnet3 driver internals. |
| `test-url-handler enable` | Enable test HTTP URL handler. |
| `tls openssl set` | Set OpenSSL engine parameters. |
| `tls openssl set-tls` | Set default TLS version or cipher suites. |
| `trace add` | Add packet-trace entries for a graph node. |
| `trace filter` | Set a classify-based filter for the trace buffer. |
| `trace frame-queue` | Enable or disable frame-queue handoff tracing. |
| `tracenode feature` | Enable tracenode feature-arc on an interface. |
| `udp decap` | Configure UDP decapsulation on a port. |
| `udp encap` | Create or delete a UDP encapsulation entry. |
| `udp-echo` | Start a built-in UDP echo client or server. |
| `undefine` | Remove a previously defined CLI macro. |
| `vrrp peers` | Configure VRRP peer addresses for unicast mode. |
| `vrrp proto` | Set VRRP protocol parameters. |
| `vrrp vr add` | Create a VRRP virtual router instance. |
| `vrrp vr del` | Delete a VRRP virtual router instance. |
| `vrrp vr track-if` | Add or remove VRRP interface-tracking entry. |
| `wait` | Pause CLI processing for a specified duration. |
| `wireguard create` | Create a WireGuard tunnel interface. |
| `wireguard delete` | Delete a WireGuard tunnel interface. |
| `wireguard peer add` | Add a WireGuard peer to a tunnel. |
| `wireguard peer remove` | Remove a WireGuard peer from a tunnel. |
