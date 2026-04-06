# VPP `vppctl sh ?` Command Catalog

## Purpose
- This catalog groups all source-registered `show ...` commands so you can map them into MCP tools for troubleshooting workflows.
- `vppctl sh ...` is an alias of `vppctl show ...`; every command below can be invoked with `sh`.

## Scope and Extraction Method
- Snapshot generated: 2026-02-27 08:17:19 UTC.
- Source scan: all `VLIB_CLI_COMMAND` registrations in `src/**/*.c` with `.path` beginning with `"show "` (excluding `src/examples/**`).
- Total unique commands found: **413**.
- Availability note: commands under `src/plugins/<plugin>/...` are only available when that plugin is built and loaded.

## Category Index
| Category | Command Count |
|---|---:|
| Core Platform & CLI Infrastructure | 30 |
| Telemetry, Trace & Performance | 25 |
| Interfaces, L2 & Device Drivers | 36 |
| IP Routing, FIB, Forwarding & SR/MPLS | 89 |
| Transport, Session & Application Plane | 41 |
| Security, NAT, ACL & Crypto | 82 |
| Tunnels, Encapsulation & Overlay Networking | 70 |
| QoS, Classification & Traffic Policy | 25 |
| Service Chaining & Tenant Frameworks | 15 |

## Commands by Category
### Core Platform & CLI Infrastructure
Global process/runtime state, command infrastructure, memory/process visibility, and platform-wide control-plane diagnostics.

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show api` | Displays api state. Syntax: `Show API information`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlibmemory/vlib_api_cli.c:175` |
| `show api clients` | Client information. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlibmemory/vlib_api_cli.c:184` |
| `show api dump` | Displays api dump state. Syntax: `show api dump file <filename> [numeric \\| compare-current]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlibmemory/vlib_api_cli.c:1461` |
| `show api histogram` | Displays api histogram state. Use this for quick health checks and baseline comparisons. | `core` | `src/vlibmemory/vlib_api_cli.c:51` |
| `show api message-table` | Message Table. Use this to validate route resolution and forwarding decisions. | `core` | `src/vlibmemory/vlib_api_cli.c:242` |
| `show api plugin` | Displays api plugin state. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlibmemory/vlib_api_cli.c:307` |
| `show api ring-stats` | Message ring statistics. Use this for quick health checks and baseline comparisons. | `core` | `src/vlibmemory/memory_api.c:1119` |
| `show api trace-status` | Display API trace status. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlibmemory/vlib_api_cli.c:194` |
| `show bihash` | Displays bihash state. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vpp/vnet/main.c:532` |
| `show cli` | Displays cli state. Syntax: `show cli [mp-safe][not-mp-safe][hit][clear-hit]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/cli.c:1876` |
| `show cli-sessions` | Displays cli sessions state. Syntax: `Show current CLI sessions`. Use this when debugging live flow/session state. | `core` | `src/vlib/unix/cli.c:3717` |
| `show clock` | Displays clock state. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/threads.c:1813` |
| `show cpu` | Displays cpu state. Syntax: `Show cpu information`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/cli.c:941` |
| `show dma` | Displays dma state. Syntax: `show dma [config <x>]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/dma/cli.c:156` |
| `show dma backends` | Displays dma backends state. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/dma/cli.c:27` |
| `show files` | Displays files state. Syntax: `Show files in use`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/file.c:301` |
| `show gdb` | Describe functions which can be called from gdb. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vnet/unix/gdb_funcs.c:312` |
| `show hash` | Displays hash state. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vnet/hash/cli.c:29` |
| `show init-function` | Displays init function state. Syntax: `show init-function [init \\| enter \\| exit][verbose [nn]]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/init.c:646` |
| `show macro` | Displays macro state. Syntax: `show macro [noevaluate]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/unix/cli.c:4023` |
| `show memory` | Displays memory state. Syntax: `show memory [api-segment][stats-segment][verbose] [map][main-heap]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/cli.c:909` |
| `show physmem` | Displays physmem state. Syntax: `show physmem [verbose \\| detail \\| map]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/physmem.c:155` |
| `show plugins` | Displays plugins state. Syntax: `show loaded plugins`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/unix/plugin.c:1111` |
| `show terminal` | Displays terminal state. Syntax: `Show current session terminal settings`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/unix/cli.c:3636` |
| `show threads` | Displays threads state. Syntax: `Show threads`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/threads_cli.c:91` |
| `show unix errors` | Displays unix errors state. Syntax: `Show Unix system call error history`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/unix/cli.c:3530` |
| `show version` | Displays version state. Syntax: `show version [verbose] [cmdline]`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vpp/app/version.c:128` |
| `show vlib frame-allocation` | Displays vlib frame allocation state. Syntax: `Show node dispatch frame statistics`. Use this for global process/runtime diagnostics and CLI inventorying. | `core` | `src/vlib/main.c:210` |
| `show vlib graph` | Displays vlib graph state. Syntax: `Show packet processing node graph`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/node_cli.c:56` |
| `show vlib graphviz` | Dump packet processing node graph as a graphviz dotfile. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/node_cli.c:293` |

### Telemetry, Trace & Performance
Dataplane introspection commands for trace buffers, node/runtime counters, perf monitoring, and debug-time graph visibility.

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show buffer metadata` | Displays buffer metadata state. Use this for instrumentation-driven troubleshooting and performance triage. | `plugin:mdata` | `src/plugins/mdata/mdata.c:466` |
| `show buffer traces` | Displays buffer traces state. Syntax: `show buffer traces [status\\|verbose]`. Use this for dataplane pipeline debugging and hotspot triage. | `plugin:bufmon` | `src/plugins/bufmon/bufmon.c:284` |
| `show buffers` | Displays buffers state. Syntax: `show buffers [detail] - show packet buffer allocation`. Use this for instrumentation-driven troubleshooting and performance triage. | `core` | `src/vlib/buffer.c:630` |
| `show errors` | Displays errors state. Syntax: `Show error counts`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/error.c:333` |
| `show event-logger` | Displays event logger state. Syntax: `Show event logger info`. Use this for instrumentation-driven troubleshooting and performance triage. | `core` | `src/vlib/main.c:774` |
| `show frame-queue` | Displays frame queue state. Syntax: `show frame-queue trace`. Use this for instrumentation-driven troubleshooting and performance triage. | `core` | `src/vlib/threads_cli.c:365` |
| `show frame-queue histogram` | Displays frame queue histogram state. Use this for quick health checks and baseline comparisons. | `core` | `src/vlib/threads_cli.c:371` |
| `show graph` | Displays graph state. Syntax: `show graph [node <index>\\|<name>] [want_arcs] [input\\|trace_supported] [drop] [output] [punt] [handoff] [no_free] [polling] [interrupt]`. Use this for dataplane pipeline debugging and hotspot triage. | `plugin:tracedump` | `src/plugins/tracedump/graph_cli.c:126` |
| `show logging` | Displays logging state. Use this for instrumentation-driven troubleshooting and performance triage. | `core` | `src/vlib/log.c:443` |
| `show logging configuration` | Displays logging configuration state. Use this for instrumentation-driven troubleshooting and performance triage. | `core` | `src/vlib/log.c:488` |
| `show node` | Displays node state. Syntax: `show node [index] <node-name \\| node-index>`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/node_cli.c:826` |
| `show node counters` | Displays node counters state. Syntax: `Show node counters`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/error.c:339` |
| `show packet-generator` | Displays packet generator state. Syntax: `show packet-generator [verbose]`. Use this for instrumentation-driven troubleshooting and performance triage. | `core` | `src/vnet/pg/cli.c:214` |
| `show packet-generator interface` | Displays packet generator interface state. Syntax: `show packet-generator interface {<interface name> \\| sw_if_index <sw_idx>}`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/pg/cli.c:813` |
| `show pcap filter function` | Displays pcap filter function state. Use this for instrumentation-driven troubleshooting and performance triage. | `core` | `src/vnet/interface_cli.c:2447` |
| `show perfmon active-bundle` | Displays perfmon active bundle state. Use this for instrumentation-driven troubleshooting and performance triage. | `plugin:perfmon` | `src/plugins/perfmon/cli.c:334` |
| `show perfmon bundle` | Displays perfmon bundle state. Syntax: `show perfmon bundle [<bundle-name>] [verbose]`. Use this for instrumentation-driven troubleshooting and performance triage. | `plugin:perfmon` | `src/plugins/perfmon/cli.c:215` |
| `show perfmon source` | Displays perfmon source state. Syntax: `show perfmon source [<source-name>] [verbose]`. Use this for instrumentation-driven troubleshooting and performance triage. | `plugin:perfmon` | `src/plugins/perfmon/cli.c:316` |
| `show perfmon statistics` | Displays perfmon statistics state. Syntax: `show perfmon statistics [raw]`. Use this for quick health checks and baseline comparisons. | `plugin:perfmon` | `src/plugins/perfmon/cli.c:493` |
| `show runtime` | Displays runtime state. Syntax: `show runtime [time] [verbose] [max] [summary]`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/node_cli.c:604` |
| `show statistics hash` | Displays statistics hash state. Use this for quick health checks and baseline comparisons. | `core` | `src/vlib/stats/cli.c:220` |
| `show statistics segment` | Displays statistics segment state. Syntax: `show statistics segment [counter-name] [verbose]`. Use this for quick health checks and baseline comparisons. | `core` | `src/vlib/stats/cli.c:226` |
| `show trace` | Displays trace state. Syntax: `Show trace buffer [max COUNT]`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/trace.c:351` |
| `show trace filter function` | Displays trace filter function state. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/trace.c:667` |
| `show trace timestamp-format` | Displays trace timestamp format state. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vlib/trace.c:797` |

### Interfaces, L2 & Device Drivers
Interface and driver-level operational state (NICs, virtual interfaces, L2 constructs, host interface integration).

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show bond` | Displays bond state. Syntax: `show bond [details]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/bonding/cli.c:1075` |
| `show bridge-domain` | Displays bridge domain state. Syntax: `show bridge-domain [bridge-domain-id [detail\\|int\\|arp\\|bd-tag]]`. Use this for link-layer and NIC/driver diagnostics. | `core` | `src/vnet/l2/l2_bd.c:1321` |
| `show cdp` | Displays cdp state. Syntax: `Show cdp command`. Use this for link-layer and NIC/driver diagnostics. | `plugin:cdp` | `src/plugins/cdp/cdp_input.c:455` |
| `show device` | Displays device state. Syntax: `show device [counters] [zero-counters] [debug] [debug-level <n>] [<device-id> ...]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/dev/cli.c:377` |
| `show device drivers` | Displays device drivers state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/dev/cli.c:445` |
| `show dpdk buffer` | Displays dpdk buffer state. Use this for link-layer and NIC/driver diagnostics. | `plugin:dpdk` | `src/plugins/dpdk/device/cli.c:70` |
| `show dpdk physmem` | Displays dpdk physmem state. Use this for link-layer and NIC/driver diagnostics. | `plugin:dpdk` | `src/plugins/dpdk/device/cli.c:160` |
| `show dpdk version` | Displays dpdk version state. Use this for link-layer and NIC/driver diagnostics. | `plugin:dpdk` | `src/plugins/dpdk/device/cli.c:360` |
| `show hardware-interfaces` | Displays hardware interfaces state. Syntax: `show hardware-interfaces [brief\\|verbose\\|detail] [bond] [<interface> [<interface> [..]]] [<sw_idx> [<sw_idx> [..]]]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/interface_cli.c:221` |
| `show interface` | Displays interface state. Syntax: `show interface [address\\|addr\\|features\\|feat\\|vtr\\|tag] [<interface> [<interface> [..]]] [verbose]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/interface_cli.c:478` |
| `show interface rx-placement` | Displays interface rx placement state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/interface_cli.c:1660` |
| `show interface secondary-mac-address` | Displays interface secondary mac address state. Syntax: `show interface secondary-mac-address [<interface>]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/interface_cli.c:1227` |
| `show interface span` | Shows SPAN mirror table. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/span/span.c:226` |
| `show interface transceiver` | Displays interface transceiver state. Syntax: `show interface transceiver [<interface>] [module] [diag] [eeprom] [verbose]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/interface_cli.c:2782` |
| `show interface tx-hash` | Displays interface tx hash state. Syntax: `show interface tx-hash [interface]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/interface_cli.c:2628` |
| `show l2 rewrite entries` | Displays l2 rewrite entries state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/l2/l2_rw.c:554` |
| `show l2 rewrite interfaces` | Displays l2 rewrite interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/l2/l2_rw.c:521` |
| `show l2fib` | Displays l2fib state. Syntax: `show l2fib [all] \\| [bd_id <nn> \\| bd_index <nn>] [learn \\| add] \\| [raw]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/l2/l2_fib.c:343` |
| `show l2patch` | Displays l2patch state. Syntax: `Show l2 interface cross-connect entries`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/l2/l2_patch.c:415` |
| `show l2xcrw` | Displays l2xcrw state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/l2/l2_xcrw.c:571` |
| `show lacp` | Displays lacp state. Syntax: `show lacp [<interface>] [details]`. Use this for link-layer and NIC/driver diagnostics. | `plugin:lacp` | `src/plugins/lacp/cli.c:300` |
| `show lcp` | Displays lcp state. Syntax: `show lcp [phy <interface>]`. Use this for link-layer and NIC/driver diagnostics. | `plugin:linux-cp` | `src/plugins/linux-cp/lcp_cli.c:348` |
| `show lcp adj` | Displays lcp adj state. Use this to validate route resolution and forwarding decisions. | `plugin:linux-cp` | `src/plugins/linux-cp/lcp_adj.c:184` |
| `show lcp ethertype` | Displays lcp ethertype state. Use this for link-layer and NIC/driver diagnostics. | `plugin:linux-cp` | `src/plugins/linux-cp/lcp_cli.c:405` |
| `show lldp` | Displays lldp state. Syntax: `show lldp [detail]`. Use this for link-layer and NIC/driver diagnostics. | `plugin:lldp` | `src/plugins/lldp/lldp_cli.c:697` |
| `show mactime` | Displays mactime state. Syntax: `show mactime [verbose]`. Use this for link-layer and NIC/driver diagnostics. | `plugin:mactime` | `src/plugins/mactime/mactime.c:668` |
| `show memif` | Displays memif state. Syntax: `show memif [<interface>] [descriptors]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:memif` | `src/plugins/memif/cli.c:462` |
| `show mode` | Displays mode state. Syntax: `show mode [<if-name1> <if-name2> ...]`. Use this for link-layer and NIC/driver diagnostics. | `core` | `src/vnet/l2/l2_input.c:862` |
| `show pci` | Displays pci state. Syntax: `show pci [all]`. Use this for link-layer and NIC/driver diagnostics. | `core` | `src/vlib/pci/pci.c:378` |
| `show soft-rss` | Displays soft rss state. Syntax: `show soft-rss [<interface>]`. Use this for link-layer and NIC/driver diagnostics. | `plugin:soft-rss` | `src/plugins/soft-rss/cli.c:236` |
| `show tap` | Displays tap state. Syntax: `show tap {<interface>] [descriptors]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:tap` | `src/plugins/tap/cli.c:378` |
| `show tun` | Displays tun state. Syntax: `show tun {<interface>] [descriptors]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:tap` | `src/plugins/tap/cli.c:423` |
| `show vhost-user` | Displays vhost user state. Syntax: `show vhost-user [<interface> [<interface> [..]]] [[descriptors] [verbose]]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:vhost` | `src/plugins/vhost/vhost_user.c:2531` |
| `show virtio pci` | Displays virtio pci state. Syntax: `show virtio pci [<interface>] [descriptors \\| desc] [debug-device]`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/devices/virtio/cli.c:240` |
| `show vmxnet3` | Displays vmxnet3 state. Syntax: `show vmxnet3 [[<interface>] ([desc] \\| ([rx-comp] \\| [rx-desc-0] \\| [rx-desc-1] \\| [tx-comp] \\| [tx-desc]) [<slot>])]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:vmxnet3` | `src/plugins/vmxnet3/cli.c:567` |
| `show vrrp vr` | Displays vrrp vr state. Syntax: `show vrrp vr [(<intf_name>\\|sw_if_index <n>)]`. Use this for link-layer and NIC/driver diagnostics. | `plugin:vrrp` | `src/plugins/vrrp/vrrp_cli.c:207` |

### IP Routing, FIB, Forwarding & SR/MPLS
Routing/forwarding state for IPv4/IPv6/MPLS/SRv6/BIER and core forwarding graph objects (FIB, adj, DPO, punt, uRPF).

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show adj` | Displays adj state. Syntax: `show adj [<adj_index>] [interface] [summary]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/adj/adj.c:728` |
| `show adj nbr` | Displays adj nbr state. Syntax: `show adj nbr [<adj_index>] [interface]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/adj/adj_nbr.c:1082` |
| `show arp proxy` | Displays arp proxy state. Syntax: `show ip arp`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/arp/arp_proxy.c:424` |
| `show bier bift` | Displays bier bift state. Syntax: `show bier bift [set <value>] [sd <value>] [bsl <value>]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/bier/bier_bift_table.c:277` |
| `show bier disp entry` | Displays bier disp entry state. Syntax: `show bier disp entry index`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/bier/bier_disp_entry.c:382` |
| `show bier disp table` | Displays bier disp table state. Syntax: `show bier disp table [index]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/bier/bier_disp_table.c:391` |
| `show bier fib` | Displays bier fib state. Syntax: `show bier fib [table-index] [bit-position]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/bier/bier_update.c:182` |
| `show bier fmask` | Displays bier fmask state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/bier/bier_fmask.c:564` |
| `show bier imp` | Displays bier imp state. Syntax: `show bier imp [index]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/bier/bier_imp.c:282` |
| `show dpo memory` | Displays dpo memory state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/dpo/dpo.c:648` |
| `show fib entry` | Displays fib entry state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_entry.c:1900` |
| `show fib entry-delegate` | Displays fib entry delegate state. Syntax: `show fib entry delegate`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_entry_delegate.c:337` |
| `show fib memory` | Displays fib memory state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_node.c:288` |
| `show fib path-lists` | Displays fib path lists state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_path_list.c:1467` |
| `show fib paths` | Displays fib paths state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_path.c:2875` |
| `show fib source` | Displays fib source state. Syntax: `show fib source [prio]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_source.c:201` |
| `show fib uRPF` | Displays fib uRPF state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_urpf_list.c:233` |
| `show fib walk` | Displays fib walk state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/fib_walk.c:1111` |
| `show inacl` | Displays inacl state. Syntax: `show inacl type [ip4\\|ip6\\|l2]`. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/classify/in_out_acl.c:382` |
| `show ip` | Internet protocol (IP) show commands. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/lookup.c:542` |
| `show ip container` | Displays ip container state. Syntax: `show ip container <address> <interface>`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip_container_proxy.c:242` |
| `show ip fib` | Displays ip fib state. Syntax: `show ip fib [summary] [table <table-id>] [index <fib-id>] [<ip4-addr>[/<mask>]] [mtrie] [detail] [memory]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/ip4_fib.c:610` |
| `show ip local` | Displays ip local state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip4_forward.c:1953` |
| `show ip mfib` | Displays ip mfib state. Syntax: `show ip mfib [summary] [table <table-id>] [index <fib-id>] [<grp-addr>[/<mask>]] [<grp-addr>] [<src-addr> <grp-addr>]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/mfib/ip4_mfib.c:635` |
| `show ip neighbor` | Displays ip neighbor state. Syntax: `show ip neighbor [interface]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1066` |
| `show ip neighbor-config` | Displays ip neighbor config state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1873` |
| `show ip neighbor-stats` | Displays ip neighbor stats state. Syntax: `show ip neighbor-stats [interface]`. Use this for quick health checks and baseline comparisons. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1884` |
| `show ip neighbor-watcher` | Displays ip neighbor watcher state. Syntax: `show ip neighbors-watcher`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor_watch.c:237` |
| `show ip neighbors` | Displays ip neighbors state. Syntax: `show ip neighbors [interface]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1051` |
| `show ip pmtu` | Displays ip pmtu state. Syntax: `show ip path MTU`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip_path_mtu.c:872` |
| `show ip punt redirect` | Displays ip punt redirect state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip4_punt_drop.c:260` |
| `show ip session redirect` | Displays ip session redirect state. Syntax: `show ip session redirect [all\\|[table <table-index>] [punt\\|acl] [ip4\\|ip6] [match]]`. Use this when debugging live flow/session state. | `plugin:ip_session_redirect` | `src/plugins/ip_session_redirect/redirect.c:277` |
| `show ip source-and-port-range-check` | Displays ip source and port range check state. Syntax: `show ip source-and-port-range-check vrf <table-id> <ip-addr> [port <n>]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip4_source_and_port_range_check.c:1375` |
| `show ip table` | Displays ip table state. Syntax: `show ip table <table-id>`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/ip/lookup.c:618` |
| `show ip4 neighbor` | Displays ip4 neighbor state. Syntax: `show ip4 neighbor [interface]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1071` |
| `show ip4 neighbor-sorted` | Displays ip4 neighbor sorted state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1081` |
| `show ip4 neighbors` | Displays ip4 neighbors state. Syntax: `show ip4 neighbors [interface]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1056` |
| `show ip4-full-reassembly` | Displays ip4 full reassembly state. Syntax: `show ip4-full-reassembly [details]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/reass/ip4_full_reass.c:1829` |
| `show ip4-sv-reassembly` | Displays ip4 sv reassembly state. Syntax: `show ip4-sv-reassembly [details]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/reass/ip4_sv_reass.c:1644` |
| `show ip6` | Internet protocol version 6 (IPv6) show commands. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/lookup.c:547` |
| `show ip6 addresses` | Displays ip6 addresses state. Use this for L3 forwarding correctness and path selection troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_pd_client_cp.c:1100` |
| `show ip6 connection-tracker` | Displays ip6 connection tracker state. Use this when debugging live flow/session state. | `plugin:ct6` | `src/plugins/ct6/ct6.c:314` |
| `show ip6 fib` | Displays ip6 fib state. Syntax: `show ip6 fib [summary] [table <table-id>] [index <fib-id>] [<ip6-addr>[/<width>]] [detail]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/ip6_fib.c:725` |
| `show ip6 hbh` | Displays ip6 hbh state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip6_hop_by_hop.c:1047` |
| `show ip6 interface` | Displays ip6 interface state. Syntax: `show ip6 interface <interface>`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/ip/ip6_link.c:721` |
| `show ip6 local` | Displays ip6 local state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip6_forward.c:2969` |
| `show ip6 mfib` | Displays ip6 mfib state. Syntax: `show ip mfib [summary] [table <table-id>] [index <fib-id>] [<grp-addr>[/<mask>]] [<grp-addr>] [<src-addr> <grp-addr>]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/mfib/ip6_mfib.c:769` |
| `show ip6 neighbor` | Displays ip6 neighbor state. Syntax: `show ip6 neighbor [interface]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1076` |
| `show ip6 neighbor-sorted` | Displays ip6 neighbor sorted state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1086` |
| `show ip6 neighbors` | Displays ip6 neighbors state. Syntax: `show ip6 neighbors [interface]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip-neighbor/ip_neighbor.c:1061` |
| `show ip6 pd clients` | Displays ip6 pd clients state. Use this for L3 forwarding correctness and path selection troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_pd_client_cp.c:1199` |
| `show ip6 prefixes` | Displays ip6 prefixes state. Use this for L3 forwarding correctness and path selection troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_pd_client_cp.c:1134` |
| `show ip6 punt redirect` | Displays ip6 punt redirect state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip6_punt_drop.c:263` |
| `show ip6 table` | Displays ip6 table state. Syntax: `show ip6 table <table-id>`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/ip/lookup.c:624` |
| `show ip6-full-reassembly` | Displays ip6 full reassembly state. Syntax: `show ip6-full-reassembly [details]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/reass/ip6_full_reass.c:1952` |
| `show ip6-ll` | Displays ip6 ll state. Syntax: `show ip6-ll [summary] [interface] [<ip6-addr>[/<width>]] [detail]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/ip6_ll_table.c:282` |
| `show ip6-sv-reassembly` | Displays ip6 sv reassembly state. Syntax: `show ip6-sv-reassembly [details]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/ip/reass/ip6_sv_reass.c:1359` |
| `show load-balance` | Displays load balance state. Syntax: `show load-balance [<index>]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/dpo/load_balance.c:1102` |
| `show load-balance-map` | Displays load balance map state. Syntax: `show load-balance-map [<index>]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/dpo/load_balance_map.c:579` |
| `show lookup-dpo` | Displays lookup dpo state. Syntax: `show lookup-dpo [<index>]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/dpo/lookup_dpo.c:1440` |
| `show mfib entry` | Displays mfib entry state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/mfib/mfib_entry.c:1660` |
| `show mfib interface` | Displays mfib interface state. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/mfib/mfib_itf.c:255` |
| `show mfib itf flags` | Flags applicable to an MFIB interfaces. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/mfib/mfib_types.c:272` |
| `show mfib route flags` | Flags applicable to an MFIB route. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/mfib/mfib_types.c:245` |
| `show mpls fib` | Displays mpls fib state. Syntax: `show mpls fib [summary] [table <n>]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/fib/mpls_fib.c:469` |
| `show mpls interface` | Displays mpls interface state. Syntax: `Show MPLS interface forwarding`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/mpls/interface.c:217` |
| `show mpls tunnel` | Displays mpls tunnel state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/mpls/mpls_tunnel.c:1025` |
| `show outacl` | Displays outacl state. Syntax: `show outacl type [ip4\\|ip6\\|l2]`. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/classify/in_out_acl.c:387` |
| `show punt client` | Displays punt client state. Syntax: `show client[s] registered with the punt infra`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vlib/punt.c:595` |
| `show punt db` | Displays punt db state. Syntax: `show the punt DB`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vlib/punt.c:655` |
| `show punt reasons` | Displays punt reasons state. Syntax: `show all punt reasons`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vlib/punt.c:616` |
| `show punt socket registrations` | Displays punt socket registrations state. Syntax: `show punt socket registrations [l4\\|exception]`. Use this when debugging live flow/session state. | `core` | `src/vnet/ip/punt.c:808` |
| `show punt stats` | Displays punt stats state. Syntax: `show the punt stats`. Use this for quick health checks and baseline comparisons. | `core` | `src/vlib/punt.c:680` |
| `show replicate` | Displays replicate state. Syntax: `show replicate [<index>]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/dpo/replicate_dpo.c:715` |
| `show sr encaps hop-limit` | Displays sr encaps hop limit state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/srv6/sr_policy_rewrite.c:1249` |
| `show sr encaps source addr` | Displays sr encaps source addr state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/srv6/sr_policy_rewrite.c:1230` |
| `show sr localsids` | Displays sr localsids state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/srv6/sr_localsid.c:697` |
| `show sr localsids behaviors` | Displays sr localsids behaviors state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/srv6/sr_localsid.c:2464` |
| `show sr mpls policies` | Displays sr mpls policies state. Use this for L3 forwarding correctness and path selection troubleshooting. | `plugin:srmpls` | `src/plugins/srmpls/sr_mpls_policy.c:626` |
| `show sr mpls steering policies` | Displays sr mpls steering policies state. Use this for L3 forwarding correctness and path selection troubleshooting. | `plugin:srmpls` | `src/plugins/srmpls/sr_mpls_steering.c:860` |
| `show sr policies` | Displays sr policies state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/srv6/sr_policy_rewrite.c:1211` |
| `show sr policy behaviors` | Displays sr policy behaviors state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/srv6/sr_policy_rewrite.c:3536` |
| `show sr steering-policies` | Displays sr steering policies state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/srv6/sr_steering.c:511` |
| `show svs` | Source VRF select show. Use this for L3 forwarding correctness and path selection troubleshooting. | `plugin:svs` | `src/plugins/svs/svs.c:578` |
| `show teib` | Displays teib state. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/teib/teib_cli.c:164` |
| `show udp encap` | Displays udp encap state. Syntax: `show udp encap [ID]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/udp/udp_encap.c:559` |
| `show udp ports` | Displays udp ports state. Syntax: `show udp ports [ip4\\|ip6] [bind\\|all\\|<port>]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/udp/udp_cli.c:306` |
| `show udp punt` | Displays udp punt state. Syntax: `show udp punt [ipv4\\|ipv6]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/udp/udp_cli.c:211` |
| `show udp transport ports` | Displays udp transport ports state. Syntax: `show udp transport ports [ip4\\|ip6] [<port>]`. Use this for L3 forwarding correctness and path selection troubleshooting. | `core` | `src/vnet/udp/udp_cli.c:398` |

### Transport, Session & Application Plane
Session-layer and host-stack state (TCP/UDP/session tables, app namespaces, DHCP/DNS/HTTP/QUIC operational views).

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show app` | Displays app state. Syntax: `show app [index] [listeners\\|client] [mq] [verbose] [transports]`. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/session/application.c:2060` |
| `show app certificate` | List app certs and keys present in store. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/session/application_crypto.c:303` |
| `show app evt-collector` | Displays app evt collector state. Syntax: `show app evt-collector [app <app> listeners-filter]`. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/session/application_eventing.c:751` |
| `show app ns` | Displays app ns state. Syntax: `show app ns [id <id> [api-clients]]`. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/session/application_namespace.c:644` |
| `show auto-sdl` | Displays auto sdl state. Syntax: `show auto-sdl [appns <id>] [<rmt-ip>]\\|[summary]`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:auto_sdl` | `src/plugins/auto_sdl/auto_sdl.c:640` |
| `show bfd` | Displays bfd state. Syntax: `show bfd [keys\\|sessions\\|echo-source]`. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/bfd/bfd_cli.c:232` |
| `show dhcp client` | Displays dhcp client state. Syntax: `show dhcp client [intfc <intfc>][verbose]`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/client.c:942` |
| `show dhcp option-82-address interface` | Displays dhcp option 82 address interface state. Syntax: `show dhcp option-82-address interface <interface>`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:dhcp` | `src/plugins/dhcp/dhcp4_proxy_node.c:1113` |
| `show dhcp proxy` | Display dhcp proxy server info. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp4_proxy_node.c:1002` |
| `show dhcp vss` | Displays dhcp vss state. Syntax: `show dhcp VSS`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp4_proxy_node.c:1066` |
| `show dhcp6 addresses` | Displays dhcp6 addresses state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_ia_na_client_cp.c:522` |
| `show dhcp6 clients` | Displays dhcp6 clients state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_ia_na_client_cp.c:582` |
| `show dhcpv6 link-address interface` | Displays dhcpv6 link address interface state. Syntax: `show dhcpv6 link-address interface <interface>`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_proxy_node.c:1173` |
| `show dhcpv6 proxy` | Display dhcpv6 proxy info. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_proxy_node.c:1064` |
| `show dhcpv6 vss` | Displays dhcpv6 vss state. Syntax: `show dhcpv6 VSS`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dhcp` | `src/plugins/dhcp/dhcp6_proxy_node.c:1127` |
| `show dns cache` | Displays dns cache state. Syntax: `show dns cache [verbose [nn]]`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dns` | `src/plugins/dns/dns.c:2154` |
| `show dns servers` | Displays dns servers state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:dns` | `src/plugins/dns/dns.c:2189` |
| `show http connect proxy client` | Displays http connect proxy client state. Syntax: `show http connect proxy [listeners] [sessions] [stats]`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:hs_apps` | `src/plugins/hs_apps/http_connect_proxy_client.c:2282` |
| `show http static server` | Displays http static server state. Syntax: `show http static server [sessions] [cache] [listeners] [verbose [<nn>]]`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:http_static` | `src/plugins/http_static/static_server.c:1525` |
| `show http stats` | Displays http stats state. Use this for quick health checks and baseline comparisons. | `plugin:http` | `src/plugins/http/http.c:1847` |
| `show http tps` | Displays http tps state. Syntax: `http tps [listeners]`. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:hs_apps` | `src/plugins/hs_apps/http_tps.c:959` |
| `show igmp config` | Displays igmp config state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:igmp` | `src/plugins/igmp/igmp_cli.c:343` |
| `show igmp ssm-ranges` | Displays igmp ssm ranges state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:igmp` | `src/plugins/igmp/igmp_ssm_range.c:117` |
| `show igmp timers` | Displays igmp timers state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:igmp` | `src/plugins/igmp/igmp_cli.c:360` |
| `show quic` | Displays quic state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:quic` | `src/plugins/quic/quic_cli.c:482` |
| `show quic crypto context` | List quic crypto contextes. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:quic` | `src/plugins/quic/quic_cli.c:487` |
| `show segment-manager` | Displays segment manager state. Syntax: `show segment-manager [segments][verbose][index <nn>]`. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/session/segment_manager.c:1217` |
| `show session` | Displays session state. Syntax: `show session [protos][states][rt-backend][verbose [n]] [transport][events][listeners <proto>] [<session-id>][thread <n> [[proto <p>] index <n>]][elog] [thread <n>][proto <proto>][state <state>][range <min> [<max>]] [lcl\\|rmt\\|ep <ip>[:<port>]][tree][force-print]`. Use this when debugging live flow/session state. | `core` | `src/vnet/session/session_cli.c:961` |
| `show session dbg clock_cycles` | Displays session dbg clock_cycles state. Use this when debugging live flow/session state. | `core` | `src/vnet/session/session_debug.c:45` |
| `show session fifo trace` | Displays session fifo trace state. Syntax: `show session fifo trace <session>`. Use this when debugging live flow/session state. | `core` | `src/vnet/session/session_cli.c:1079` |
| `show session lookup` | Displays session lookup state. Syntax: `show session lookup [table <fib-index>]`. Use this when debugging live flow/session state. | `core` | `src/vnet/session/session_lookup.c:1982` |
| `show session rules` | Displays session rules state. Syntax: `show session rules [<proto> appns <id> <lcl-ip/plen> <lcl-port> <rmt-ip/plen> <rmt-port> scope <scope>]`. Use this when debugging live flow/session state. | `core` | `src/vnet/session/session_lookup.c:1901` |
| `show session sdl` | Displays session sdl state. Syntax: `show session sdl [appns <id>] [table <fib-index>] [<rmt-ip>]`. Use this when debugging live flow/session state. | `core` | `src/vnet/session/session_sdl.c:793` |
| `show session stats` | Displays session stats state. Use this for quick health checks and baseline comparisons. | `core` | `src/vnet/session/session_cli.c:1204` |
| `show sflow` | Displays sflow state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:sflow` | `src/plugins/sflow/sflow.c:1170` |
| `show tcp config` | Displays tcp config state. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/tcp/tcp_cli.c:895` |
| `show tcp punt` | Displays tcp punt state. Use this for host-stack and application-facing datapath troubleshooting. | `core` | `src/vnet/tcp/tcp_cli.c:826` |
| `show tcp scoreboard trace` | Displays tcp scoreboard trace state. Syntax: `show tcp scoreboard trace <connection>`. Use this for dataplane pipeline debugging and hotspot triage. | `core` | `src/vnet/tcp/tcp_cli.c:682` |
| `show tcp stats` | Displays tcp stats state. Use this for quick health checks and baseline comparisons. | `core` | `src/vnet/tcp/tcp_cli.c:931` |
| `show test alpn server` | Displays test alpn server state. Use this for host-stack and application-facing datapath troubleshooting. | `plugin:hs_apps` | `src/plugins/hs_apps/alpn_server.c:271` |
| `show udp-ping summary` | Summary of udp-ping. Use this for quick health checks and baseline comparisons. | `plugin:ioam` | `src/plugins/ioam/udp-ping/udp_ping_node.c:364` |

### Security, NAT, ACL & Crypto
ACL/NAT/IPsec/WireGuard/crypto observability, translation/session state, and security policy troubleshooting commands.

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show abf attach` | Displays abf attach state. Syntax: `show abf attach <interface>`. Use this for security policy and translation/crypto state validation. | `plugin:abf` | `src/plugins/abf/abf_itf_attach.c:446` |
| `show abf policy` | Displays abf policy state. Syntax: `show abf policy <value>`. Use this for security policy and translation/crypto state validation. | `plugin:abf` | `src/plugins/abf/abf_policy.c:388` |
| `show acl-plugin acl` | Displays acl plugin acl state. Syntax: `show acl-plugin acl [index N]`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:acl` | `src/plugins/acl/acl.c:3679` |
| `show acl-plugin decode 5tuple` | Displays acl plugin decode 5tuple state. Syntax: `show acl-plugin decode 5tuple XXXX XXXX XXXX XXXX XXXX XXXX`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:acl` | `src/plugins/acl/acl.c:3697` |
| `show acl-plugin interface` | Displays acl plugin interface state. Syntax: `show acl-plugin interface [sw_if_index N] [acl]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:acl` | `src/plugins/acl/acl.c:3703` |
| `show acl-plugin lookup context` | Displays acl plugin lookup context state. Syntax: `show acl-plugin lookup context [index N]`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:acl` | `src/plugins/acl/acl.c:3685` |
| `show acl-plugin lookup user` | Displays acl plugin lookup user state. Syntax: `show acl-plugin lookup user [index N]`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:acl` | `src/plugins/acl/acl.c:3691` |
| `show acl-plugin macip acl` | Displays acl plugin macip acl state. Syntax: `show acl-plugin macip acl [index N]`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:acl` | `src/plugins/acl/acl.c:3727` |
| `show acl-plugin macip interface` | Displays acl plugin macip interface state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:acl` | `src/plugins/acl/acl.c:3733` |
| `show acl-plugin memory` | Displays acl plugin memory state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:acl` | `src/plugins/acl/acl.c:3709` |
| `show acl-plugin sessions` | Displays acl plugin sessions state. Use this when debugging live flow/session state. | `plugin:acl` | `src/plugins/acl/acl.c:3715` |
| `show acl-plugin tables` | Displays acl plugin tables state. Syntax: `show acl-plugin tables [ acl [index N] \\| applied [ lc_index N ] \\| mask \\| hash [verbose N] ]`. Use this to validate route resolution and forwarding decisions. | `plugin:acl` | `src/plugins/acl/acl.c:3721` |
| `show cnat client` | Displays cnat client state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:cnat` | `src/plugins/cnat/cnat_client.c:303` |
| `show cnat session` | Displays cnat session state. Use this when debugging live flow/session state. | `plugin:cnat` | `src/plugins/cnat/cnat_session.c:127` |
| `show cnat snat-policy` | Displays cnat snat policy state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:cnat` | `src/plugins/cnat/cnat_snat_policy.c:470` |
| `show cnat timestamp` | Displays cnat timestamp state. Syntax: `show cnat timestamp [verbose]`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:cnat` | `src/plugins/cnat/cnat_session.c:317` |
| `show cnat translation` | Displays cnat translation state. Syntax: `show cnat translation <VIP>`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:cnat` | `src/plugins/cnat/cnat_translation.c:477` |
| `show crypto async status` | Displays crypto async status state. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/crypto/cli.c:215` |
| `show crypto engines` | Displays crypto engines state. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/crypto/cli.c:35` |
| `show crypto handlers` | Displays crypto handlers state. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/crypto/cli.c:93` |
| `show cryptodev assignment` | Displays cryptodev assignment state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:dpdk` | `src/plugins/dpdk/cryptodev/cryptodev.c:649` |
| `show cryptodev cache status` | Displays cryptodev cache status state. Syntax: `show status of all cryptodev cache rings`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:dpdk` | `src/plugins/dpdk/cryptodev/cryptodev.c:733` |
| `show det44 interfaces` | Displays det44 interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/det44/det44_cli.c:674` |
| `show det44 mappings` | Displays det44 mappings state. Use this for security policy and translation/crypto state validation. | `plugin:nat` | `src/plugins/nat/det44/det44_cli.c:516` |
| `show det44 sessions` | Displays det44 sessions state. Use this when debugging live flow/session state. | `plugin:nat` | `src/plugins/nat/det44/det44_cli.c:563` |
| `show det44 timeouts` | Displays det44 timeouts state. Use this for security policy and translation/crypto state validation. | `plugin:nat` | `src/plugins/nat/det44/det44_cli.c:627` |
| `show dslite aftr-tunnel-endpoint-address` | Displays dslite aftr tunnel endpoint address state. Use this to validate route resolution and forwarding decisions. | `plugin:nat` | `src/plugins/nat/dslite/dslite_cli.c:328` |
| `show dslite b4-tunnel-endpoint-address` | Displays dslite b4 tunnel endpoint address state. Use this to validate route resolution and forwarding decisions. | `plugin:nat` | `src/plugins/nat/dslite/dslite_cli.c:356` |
| `show dslite pool` | Displays dslite pool state. Use this for security policy and translation/crypto state validation. | `plugin:nat` | `src/plugins/nat/dslite/dslite_cli.c:300` |
| `show dslite sessions` | Displays dslite sessions state. Use this when debugging live flow/session state. | `plugin:nat` | `src/plugins/nat/dslite/dslite_cli.c:377` |
| `show ikev2 profile` | Displays ikev2 profile state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:ikev2` | `src/plugins/ikev2/ikev2_cli.c:666` |
| `show ikev2 sa` | Displays ikev2 sa state. Syntax: `show ikev2 sa [rspi <rspi>] [details]`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:ikev2` | `src/plugins/ikev2/ikev2_cli.c:271` |
| `show ikev2 sleep interval` | Displays ikev2 sleep interval state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:ikev2` | `src/plugins/ikev2/ikev2_cli.c:757` |
| `show interface tcp-mss-clamp` | Displays interface tcp mss clamp state. Syntax: `show interface tcp-mss-clamp [interface-name]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:mss_clamp` | `src/plugins/mss_clamp/mss_clamp.c:260` |
| `show ipsec all` | Displays ipsec all state. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/ipsec/ipsec_cli.c:539` |
| `show ipsec backends` | Displays ipsec backends state. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/ipsec/ipsec_cli.c:722` |
| `show ipsec interface` | Displays ipsec interface state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/ipsec/ipsec_itf.c:496` |
| `show ipsec protect` | Displays ipsec protect state. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/ipsec/ipsec_cli.c:883` |
| `show ipsec protect-hash` | Displays ipsec protect hash state. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/ipsec/ipsec_cli.c:937` |
| `show ipsec sa` | Displays ipsec sa state. Syntax: `show ipsec sa [index]`. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/ipsec/ipsec_cli.c:605` |
| `show ipsec spd` | Displays ipsec spd state. Syntax: `show ipsec spd [index]`. Use this when troubleshooting policy, translation, or crypto paths. | `core` | `src/vnet/ipsec/ipsec_cli.c:645` |
| `show ipsec tunnel` | Displays ipsec tunnel state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/ipsec/ipsec_cli.c:661` |
| `show nat mss-clamping` | Displays nat mss clamping state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2004` |
| `show nat timeouts` | Displays nat timeouts state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:1935` |
| `show nat workers` | Displays nat workers state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:1901` |
| `show nat44 addresses` | Displays nat44 addresses state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2077` |
| `show nat44 ei addr-port-assignment-alg` | Displays nat44 ei addr port assignment alg state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1721` |
| `show nat44 ei addresses` | Displays nat44 ei addresses state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1865` |
| `show nat44 ei ha` | Displays nat44 ei ha state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1787` |
| `show nat44 ei hash tables` | Displays nat44 ei hash tables state. Syntax: `show nat44 ei hash tables [detail\\|verbose]`. Use this to validate route resolution and forwarding decisions. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1823` |
| `show nat44 ei interface address` | Displays nat44 ei interface address state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1996` |
| `show nat44 ei interfaces` | Displays nat44 ei interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1899` |
| `show nat44 ei mss-clamping` | Displays nat44 ei mss clamping state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1749` |
| `show nat44 ei sessions` | Displays nat44 ei sessions state. Syntax: `show nat44 ei sessions [detail] [filter saddr <ip>]`. Use this when debugging live flow/session state. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:2008` |
| `show nat44 ei static mappings` | Displays nat44 ei static mappings state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1967` |
| `show nat44 ei timeouts` | Displays nat44 ei timeouts state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1661` |
| `show nat44 ei workers` | Displays nat44 ei workers state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ei/nat44_ei_cli.c:1627` |
| `show nat44 hash tables` | Displays nat44 hash tables state. Syntax: `show nat44 hash tables [detail\\|verbose]`. Use this to validate route resolution and forwarding decisions. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2016` |
| `show nat44 interface address` | Displays nat44 interface address state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2287` |
| `show nat44 interfaces` | Displays nat44 interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2110` |
| `show nat44 sessions` | Displays nat44 sessions state. Syntax: `show nat44 sessions [filter {i2o \\| o2i} {saddr <ip4-addr> \\| sport <n> \\| daddr <ip4-addr> \\| dport <n> \\| proto <proto>} [filter .. [..]]]`. Use this when debugging live flow/session state. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2299` |
| `show nat44 static mappings` | Displays nat44 static mappings state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2217` |
| `show nat44 summary` | Displays nat44 summary state. Use this for quick health checks and baseline comparisons. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2047` |
| `show nat44 vrf tables` | Displays nat44 vrf tables state. Use this to validate route resolution and forwarding decisions. | `plugin:nat` | `src/plugins/nat/nat44-ed/nat44_ed_cli.c:2270` |
| `show nat64 bib` | Displays nat64 bib state. Syntax: `show nat64 bib all\\|tcp\\|udp\\|icmp\\|unknown`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat64/nat64_cli.c:890` |
| `show nat64 interfaces` | Displays nat64 interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/nat64/nat64_cli.c:848` |
| `show nat64 pool` | Displays nat64 pool state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat64/nat64_cli.c:815` |
| `show nat64 prefix` | Displays nat64 prefix state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat64/nat64_cli.c:952` |
| `show nat64 session table` | Displays nat64 session table state. Syntax: `show nat64 session table all\\|tcp\\|udp\\|icmp\\|unknown`. Use this when debugging live flow/session state. | `plugin:nat` | `src/plugins/nat/nat64/nat64_cli.c:918` |
| `show nat66 interfaces` | Displays nat66 interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/nat66/nat66_cli.c:364` |
| `show nat66 static mappings` | Displays nat66 static mappings state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/nat66/nat66_cli.c:397` |
| `show npt66 bindings` | Displays npt66 bindings state. Use this for security policy and translation/crypto state validation. | `plugin:npt66` | `src/plugins/npt66/npt66_cli.c:117` |
| `show pnat interfaces` | Displays pnat interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:nat` | `src/plugins/nat/pnat/pnat_cli.c:325` |
| `show pnat translations` | Displays pnat translations state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:nat` | `src/plugins/nat/pnat/pnat_cli.c:305` |
| `show policer` | Displays policer state. Syntax: `show policer [name <name> \\| index <index>]`. Use this for security policy and translation/crypto state validation. | `plugin:policer` | `src/plugins/policer/policer_cli.c:744` |
| `show policer pools` | Displays policer pools state. Use this for security policy and translation/crypto state validation. | `plugin:policer` | `src/plugins/policer/policer_cli.c:759` |
| `show snort clients` | Displays snort clients state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:snort` | `src/plugins/snort/cli.c:505` |
| `show snort instances` | Displays snort instances state. Syntax: `show snort instances [verbose]`. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:snort` | `src/plugins/snort/cli.c:424` |
| `show snort interfaces` | Displays snort interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:snort` | `src/plugins/snort/cli.c:469` |
| `show wireguard interface` | Displays wireguard interface state. Syntax: `show wireguard`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:wireguard` | `src/plugins/wireguard/wireguard_cli.c:337` |
| `show wireguard mode` | Displays wireguard mode state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:wireguard` | `src/plugins/wireguard/wireguard_cli.c:392` |
| `show wireguard peer` | Displays wireguard peer state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:wireguard` | `src/plugins/wireguard/wireguard_cli.c:309` |

### Tunnels, Encapsulation & Overlay Networking
Overlay/tunnel datapath state (VXLAN/GENEVE/GTPU/LISP/NSH/PPPoE/L2TP/IPIP and related encapsulation features).

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show geneve tunnel` | Displays geneve tunnel state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:geneve` | `src/plugins/geneve/geneve.c:904` |
| `show gpe adjacency` | Displays gpe adjacency state. Use this to validate route resolution and forwarding decisions. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe_adjacency.c:563` |
| `show gpe encap` | Displays gpe encap state. Syntax: `show GPE encapulation mode`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe.c:297` |
| `show gpe entry` | Displays gpe entry state. Syntax: `show gpe entry vni <vni> vrf <vrf> [leid <leid>] reid <reid>`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe_fwd_entry.c:1478` |
| `show gpe interface` | Displays gpe interface state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe.c:371` |
| `show gpe native-forward` | Displays gpe native forward state. Use this when troubleshooting policy, translation, or crypto paths. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe.c:408` |
| `show gpe sub-interface` | Displays gpe sub interface state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe_sub_interface.c:245` |
| `show gpe tenant` | Displays gpe tenant state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe_tenant.c:304` |
| `show gpe tunnel` | Displays gpe tunnel state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:lisp` | `src/plugins/lisp/lisp-gpe/lisp_gpe_tunnel.c:256` |
| `show gre tunnel` | Displays gre tunnel state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:gre` | `src/plugins/gre/interface.c:811` |
| `show gtpu tunnel` | Displays gtpu tunnel state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:gtpu` | `src/plugins/gtpu/gtpu.c:1130` |
| `show ila entries` | Displays ila entries state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:ila` | `src/plugins/ila/ila.c:1063` |
| `show ioam analyse` | Displays ioam analyse state. Syntax: `show ioam analyser information`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:ioam` | `src/plugins/ioam/analyse/ip6/ip6_ioam_analyse.c:121` |
| `show ioam e2e` | Displays ioam e2e state. Syntax: `show ioam e2e information`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:ioam` | `src/plugins/ioam/encap/ip6_ioam_e2e.c:150` |
| `show ioam ip6 cache` | Displays ioam ip6 cache state. Syntax: `show ioam ip6 cache [verbose]`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:ioam` | `src/plugins/ioam/ip6/ioam_cache.c:293` |
| `show ioam nsh-lisp-gpe trace` | IOAM trace statistics. Use this for dataplane pipeline debugging and hotspot triage. | `plugin:nsh` | `src/plugins/nsh/nsh-md2-ioam/nsh_md2_ioam_trace.c:325` |
| `show ioam pot` | IOAM pot statistics. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:ioam` | `src/plugins/ioam/encap/ip6_ioam_pot.c:217` |
| `show ioam summary` | Displays ioam summary state. Use this for quick health checks and baseline comparisons. | `core` | `src/vnet/ip/ip6_hop_by_hop.c:1435` |
| `show ioam trace` | IOAM trace statistics. Use this for dataplane pipeline debugging and hotspot triage. | `plugin:ioam` | `src/plugins/ioam/encap/ip6_ioam_trace.c:393` |
| `show ioam vxlan-gpe trace` | IOAM trace statistics. Use this for dataplane pipeline debugging and hotspot triage. | `plugin:ioam` | `src/plugins/ioam/lib-vxlan-gpe/vxlan_gpe_ioam_trace.c:414` |
| `show ioam-trace profile` | Displays ioam trace profile state. Use this for dataplane pipeline debugging and hotspot triage. | `plugin:ioam` | `src/plugins/ioam/lib-trace/trace_util.c:179` |
| `show ipip tunnel` | Displays ipip tunnel state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/ipip/ipip_cli.c:279` |
| `show ipip tunnel-hash` | Displays ipip tunnel hash state. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/ipip/ipip_cli.c:316` |
| `show l2tpv3` | Displays l2tpv3 state. Syntax: `show l2tpv3 [verbose]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:l2tp` | `src/plugins/l2tp/l2tp.c:128` |
| `show l3xc` | Displays l3xc state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:l3xc` | `src/plugins/l3xc/l3xc.c:317` |
| `show lb` | Displays lb state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lb` | `src/plugins/lb/cli.c:362` |
| `show lb vips` | Displays lb vips state. Syntax: `show lb vips [verbose]`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lb` | `src/plugins/lb/cli.c:395` |
| `show lisp adjacencies` | Displays lisp adjacencies state. Use this to validate route resolution and forwarding decisions. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:57` |
| `show lisp eid-table` | Displays lisp eid table state. Syntax: `show lisp eid-table [local\\|remote\\|eid <eid>]`. Use this to validate route resolution and forwarding decisions. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:853` |
| `show lisp eid-table map` | Displays lisp eid table map state. Syntax: `show lisp eid-table map l2\\|l3`. Use this to validate route resolution and forwarding decisions. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:1080` |
| `show lisp locator-set` | Shows locator-sets. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:1282` |
| `show lisp map-request itr-rlocs` | Shows map-request itr-rlocs. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:1422` |
| `show lisp map-request mode` | Displays lisp map request mode state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:588` |
| `show lisp map-resolvers` | Displays lisp map resolvers state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:609` |
| `show lisp petr` | Displays lisp petr state. Syntax: `Show petr`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:1524` |
| `show lisp pitr` | Displays lisp pitr state. Syntax: `Show pitr`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:717` |
| `show lisp status` | Displays lisp status state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/lisp_cli.c:1019` |
| `show map domain` | Displays map domain state. Syntax: `show map domain index <n> [counters]`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:map` | `src/plugins/map/map.c:1427` |
| `show map stats` | Displays map stats state. Use this for quick health checks and baseline comparisons. | `plugin:map` | `src/plugins/map/map.c:1440` |
| `show nsh entry` | Displays nsh entry state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:nsh` | `src/plugins/nsh/nsh_cli.c:605` |
| `show nsh map` | Displays nsh map state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:nsh` | `src/plugins/nsh/nsh_cli.c:310` |
| `show nsim` | Display network delay simulator configuration. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:nsim` | `src/plugins/nsim/nsim.c:821` |
| `show one adjacencies` | Displays one adjacencies state. Use this to validate route resolution and forwarding decisions. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:55` |
| `show one eid-table` | Displays one eid table state. Syntax: `show one eid-table [local\\|remote\\|eid <eid>]`. Use this to validate route resolution and forwarding decisions. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:1151` |
| `show one eid-table map` | Displays one eid table map state. Syntax: `show one eid-table map l2\\|l3`. Use this to validate route resolution and forwarding decisions. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:1594` |
| `show one l2 arp entries` | Displays one l2 arp entries state. Syntax: `Show ONE L2 ARP entries`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:414` |
| `show one locator-set` | Shows locator-sets. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:1795` |
| `show one map-register fallback-threshold` | Displays one map register fallback threshold state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:865` |
| `show one map-register state` | Displays one map register state state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:2076` |
| `show one map-register ttl` | Displays one map register ttl state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:1405` |
| `show one map-request itr-rlocs` | Shows map-request itr-rlocs. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:1934` |
| `show one map-request mode` | Displays one map request mode state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:771` |
| `show one map-resolvers` | Displays one map resolvers state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:792` |
| `show one map-servers` | Displays one map servers state. Syntax: `show one map servers`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:2056` |
| `show one modes` | Displays one modes state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:2217` |
| `show one ndp entries` | Displays one ndp entries state. Syntax: `Show ONE NDP entries`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:446` |
| `show one petr` | Displays one petr state. Syntax: `Show petr`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:2035` |
| `show one pitr` | Displays one pitr state. Syntax: `Show pitr`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:1014` |
| `show one rloc state` | Displays one rloc state state. Syntax: `show one RLOC state`. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:2096` |
| `show one statistics details` | Displays one statistics details state. Syntax: `show ONE statistics`. Use this for quick health checks and baseline comparisons. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:2144` |
| `show one statistics status` | Displays one statistics status state. Syntax: `show ONE statistics enable/disable status`. Use this for quick health checks and baseline comparisons. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:2112` |
| `show one status` | Displays one status state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:lisp` | `src/plugins/lisp/lisp-cp/one_cli.c:1533` |
| `show pot profile` | Displays pot profile state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:ioam` | `src/plugins/ioam/lib-pot/pot_util.c:430` |
| `show pppoe fib` | Displays pppoe fib state. Use this to validate route resolution and forwarding decisions. | `plugin:pppoe` | `src/plugins/pppoe/pppoe.c:718` |
| `show pppoe session` | Displays pppoe session state. Use this when debugging live flow/session state. | `plugin:pppoe` | `src/plugins/pppoe/pppoe.c:637` |
| `show pvti interface` | Displays pvti interface state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:pvti` | `src/plugins/pvti/pvti.c:386` |
| `show pvti rx peers` | Displays pvti rx peers state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:pvti` | `src/plugins/pvti/pvti.c:398` |
| `show pvti tx peers` | Displays pvti tx peers state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:pvti` | `src/plugins/pvti/pvti.c:392` |
| `show vxlan tunnel` | Displays vxlan tunnel state. Syntax: `show vxlan tunnel [raw]`. Use this to validate interface provisioning and link/datapath bindings. | `plugin:vxlan` | `src/plugins/vxlan/vxlan.c:963` |
| `show vxlan-gpe` | Displays vxlan gpe state. Use this for tunnel/overlay state validation and endpoint diagnostics. | `plugin:vxlan-gpe` | `src/plugins/vxlan-gpe/vxlan_gpe.c:1011` |

### QoS, Classification & Traffic Policy
Classification, QoS markings/maps, policing, flow policy/steering, feature arcs, and policy-action troubleshooting.

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show bpf trace filter` | Displays bpf trace filter state. Use this for dataplane pipeline debugging and hotspot triage. | `plugin:bpf_trace_filter` | `src/plugins/bpf_trace_filter/cli.c:75` |
| `show classify filter` | Displays classify filter state. Syntax: `show classify filter [verbose [nn]]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/classify/vnet_classify.c:2202` |
| `show classify flow` | Displays classify flow state. Syntax: `show classify flow type [ip4\\|ip6]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/classify/flow_classify.c:206` |
| `show classify policer` | Displays classify policer state. Syntax: `show classify policer type [ip4\\|ip6\\|l2]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/classify/policer_classify.c:223` |
| `show classify tables` | Displays classify tables state. Syntax: `show classify tables [index <nn>]`. Use this to validate route resolution and forwarding decisions. | `core` | `src/vnet/classify/vnet_classify.c:2294` |
| `show features` | Displays features state. Syntax: `show features [verbose]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/feature/feature.c:526` |
| `show flow entry` | Displays flow entry state. Syntax: `show flow entry [index <index>]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/flow/flow_cli.c:247` |
| `show flow interface` | Displays flow interface state. Syntax: `show flow interface <interface name>`. Use this to validate interface provisioning and link/datapath bindings. | `core` | `src/vnet/flow/flow_cli.c:311` |
| `show flow ranges` | Displays flow ranges state. Use this for policy-chain and classification debugging. | `core` | `src/vnet/flow/flow_cli.c:269` |
| `show flowprobe feature` | Displays flowprobe feature state. Use this for policy-chain and classification debugging. | `plugin:flowprobe` | `src/plugins/flowprobe/flowprobe.c:1323` |
| `show flowprobe params` | Displays flowprobe params state. Use this for policy-chain and classification debugging. | `plugin:flowprobe` | `src/plugins/flowprobe/flowprobe.c:1329` |
| `show flowprobe statistics` | Displays flowprobe statistics state. Use this for quick health checks and baseline comparisons. | `plugin:flowprobe` | `src/plugins/flowprobe/flowprobe.c:1340` |
| `show flowprobe table` | Displays flowprobe table state. Use this to validate route resolution and forwarding decisions. | `plugin:flowprobe` | `src/plugins/flowprobe/flowprobe.c:1335` |
| `show npol interfaces` | Displays npol interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:npol` | `src/plugins/npol/npol_interface.c:108` |
| `show npol ipsets` | Displays npol ipsets state. Use this for policy-chain and classification debugging. | `plugin:npol` | `src/plugins/npol/npol_ipset.c:126` |
| `show npol policies` | Displays npol policies state. Syntax: `show npol policies [verbose]`. Use this for policy-chain and classification debugging. | `plugin:npol` | `src/plugins/npol/npol_policy.c:96` |
| `show npol rules` | Displays npol rules state. Use this for policy-chain and classification debugging. | `plugin:npol` | `src/plugins/npol/npol_rule.c:160` |
| `show qos egress map` | Displays qos egress map state. Syntax: `show qos egress map id %d`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/qos/qos_egress_map.c:256` |
| `show qos mark` | Displays qos mark state. Syntax: `show qos mark [interface]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/qos/qos_mark.c:262` |
| `show qos record` | Displays qos record state. Syntax: `show qos record [interface]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/qos/qos_record.c:276` |
| `show qos store` | Displays qos store state. Syntax: `show qos store [interface]`. Use this for policy-chain and classification debugging. | `core` | `src/vnet/qos/qos_store.c:286` |
| `show stn rules` | Displays stn rules state. Use this for policy-chain and classification debugging. | `plugin:stn` | `src/plugins/stn/stn.c:375` |
| `show sw_scheduler workers` | Displays sw_scheduler workers state. Use this for policy-chain and classification debugging. | `plugin:crypto_sw_scheduler` | `src/plugins/crypto_sw_scheduler/main.c:632` |
| `show syslog filter` | Displays syslog filter state. Use this for policy-chain and classification debugging. | `core` | `src/vnet/syslog/syslog.c:587` |
| `show syslog sender` | Displays syslog sender state. Use this for policy-chain and classification debugging. | `core` | `src/vnet/syslog/syslog.c:531` |

### Service Chaining & Tenant Frameworks
SFDP/SASC tenant and service-chain oriented state (services, tenants, sessions, chain schemas, and summaries).

| Command (`vppctl sh ...`) | Detailed Description | Availability | Source |
|---|---|---|---|
| `show sasc ingress interfaces` | Displays sasc ingress interfaces state. Use this to validate interface provisioning and link/datapath bindings. | `plugin:sasc` | `src/plugins/sasc/ingress/cli.c:108` |
| `show sasc next-indices` | Displays sasc next indices state. Use this for tenant/service-chain operational visibility. | `plugin:sasc` | `src/plugins/sasc/cli.c:354` |
| `show sasc pcap` | Displays sasc pcap state. Syntax: `show sasc pcap - Show PCAP service status`. Use this for tenant/service-chain operational visibility. | `plugin:sasc` | `src/plugins/sasc/services/pcap/cli.c:146` |
| `show sasc schema` | Displays sasc schema state. Use this for tenant/service-chain operational visibility. | `plugin:sasc` | `src/plugins/sasc/export.c:540` |
| `show sasc service-chains` | Displays sasc service chains state. Syntax: `Show service chains`. Use this for tenant/service-chain operational visibility. | `plugin:sasc` | `src/plugins/sasc/cli.c:105` |
| `show sasc services` | Displays sasc services state. Syntax: `show sasc services [verbose]`. Use this for tenant/service-chain operational visibility. | `plugin:sasc` | `src/plugins/sasc/cli.c:28` |
| `show sasc session` | Displays sasc session state. Syntax: `show sasc session [session index] [thread <n>] [tenant <tenant-id>] [0x<session-id>] [detail\\|compact]`. Use this when debugging live flow/session state. | `plugin:sasc` | `src/plugins/sasc/cli.c:265` |
| `show sasc summary` | Displays sasc summary state. Use this for quick health checks and baseline comparisons. | `plugin:sasc` | `src/plugins/sasc/cli.c:429` |
| `show sasc tenant` | Displays sasc tenant state. Syntax: `show sasc tenant [<tenant-index> [detail]]`. Use this for tenant/service-chain operational visibility. | `plugin:sasc` | `src/plugins/sasc/cli.c:196` |
| `show sfdp services` | Displays sfdp services state. Use this for tenant/service-chain operational visibility. | `core` | `src/vnet/sfdp/cli.c:546` |
| `show sfdp session-detail` | Displays sfdp session detail state. Syntax: `show sfdp session-detail 0x<session-id>`. Use this when debugging live flow/session state. | `core` | `src/vnet/sfdp/cli.c:559` |
| `show sfdp session-table` | Displays sfdp session table state. Syntax: `show sfdp session-table [tenant <tenant-id>] [max <max_value>] [unsafe-show-all]`. Use this when debugging live flow/session state. | `core` | `src/vnet/sfdp/cli.c:552` |
| `show sfdp status` | Displays sfdp status state. Use this for tenant/service-chain operational visibility. | `core` | `src/vnet/sfdp/cli.c:577` |
| `show sfdp tcp session-table` | Displays sfdp tcp session table state. Syntax: `show sfdp tcp session-table [tenant <tenant-id>]`. Use this when debugging live flow/session state. | `plugin:sfdp_services` | `src/plugins/sfdp_services/base/tcp-check/cli.c:63` |
| `show sfdp tenant` | Displays sfdp tenant state. Syntax: `show sfdp tenant [<tenant-id> [detail]]`. Use this for tenant/service-chain operational visibility. | `core` | `src/vnet/sfdp/cli.c:571` |

