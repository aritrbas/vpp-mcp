# VPP `vppctl create` Command Catalog

## Purpose
- This catalog groups all `create ...` commands available in the VPP CLI (`vppctl`).
- `create` commands instantiate new data-plane objects: interfaces, tunnels, bridge domains, and protocol entries.

## Category Index
| Category | Command Count |
|---|---:|
| Physical & Virtual Interfaces | 13 |
| L2, Bonding & Bridge Domain | 3 |
| Tunnels & Encapsulation | 9 |
| Overlay & Protocol Entries | 3 |

## Commands by Category

### Physical & Virtual Interfaces
Create host, kernel, memory-shared, and software-loopback interfaces.

| Command | Description |
|---|---|
| `create host-interface` | Create an AF_PACKET interface attached to a Linux host netdev for kernel-bypass bridging. |
| `create interface af_xdp` | Create an AF_XDP (Express Data Path) interface bound to a Linux netdev using eBPF/XDP. |
| `create interface memif` | Create a memif shared-memory interface for high-speed inter-process or inter-VM communication. |
| `create interface rdma` | Create an RDMA (ibverb/mlx5) interface for direct NIC queue access on Mellanox/ConnectX devices. |
| `create interface virtio` | Create a virtio-backed network interface for VM or container connectivity. |
| `create interface vmxnet3` | Create a VMware vmxnet3 paravirtualized network interface backed by the native VPP driver. |
| `create loopback interface` | Create a software loopback interface with an auto-assigned MAC for use as a BVI or anchor. |
| `create memif socket` | Create a named memif socket file to serve as the rendezvous point for memif connections. |
| `create packet-generator interface` | Create a virtual packet-generator (pg) interface for injecting synthetic traffic into the graph. |
| `create sub-interfaces` | Create one or more VLAN sub-interfaces (dot1q/dot1ad) on an existing parent interface. |
| `create tap` | Create a TUN/TAP interface connecting VPP to the Linux kernel network stack. |
| `create vhost-user` | Create a vhost-user interface for high-performance virtio-based guest VM networking. |
| `create teib` | Create a Tunnel Endpoint Information Base entry mapping a tunnel next-hop to an underlay peer. |

### L2, Bonding & Bridge Domain
Create link-aggregation groups and L2 broadcast domains.

| Command | Description |
|---|---|
| `create bond` | Create a bonding master interface supporting LACP, round-robin, or active-backup modes. |
| `create bridge-domain` | Create an L2 bridge domain with configurable learning, flooding, ARP termination, and MAC aging. |
| `create pppoe cp` | Create a PPPoE control-plane session for managing PPP negotiation and address assignment. |

### Tunnels & Encapsulation
Instantiate IP-in-IP, GRE, VXLAN, GENEVE, GTPU, and other tunnel endpoints.

| Command | Description |
|---|---|
| `create 6rd tunnel` | Create an IPv6 Rapid Deployment (6rd) tunnel for automatic IPv6-over-IPv4 encapsulation. |
| `create geneve tunnel` | Create a GENEVE tunnel endpoint for flexible network virtualization overlays. |
| `create gre tunnel` | Create a GRE tunnel (IPv4/IPv6, ERSPAN, or TEB mode) between two endpoints. |
| `create gtpu forward` | Create a GTPU forwarding rule to direct encapsulated packets to a specified next-hop or interface. |
| `create gtpu tunnel` | Create a GTP-U tunnel endpoint for mobile backhaul user-plane encapsulation. |
| `create ipip tunnel` | Create an IP-in-IP (IPIP or 6-in-4/4-in-6) tunnel for simple IP encapsulation. |
| `create l2tpv3 tunnel` | Create an L2TPv3 tunnel session for Ethernet-over-IP pseudowire transport. |
| `create vxlan tunnel` | Create a VXLAN tunnel endpoint for L2 overlay networking with a 24-bit VNI. |
| `create vxlan-gpe tunnel` | Create a VXLAN-GPE tunnel supporting multi-protocol encapsulation via a next-protocol field. |

### Overlay & Protocol Entries
Create NSH service-function-chaining entries and PPPoE data-plane sessions.

| Command | Description |
|---|---|
| `create nsh entry` | Create a Network Service Header entry defining SPI/SI, next-protocol, and metadata fields. |
| `create nsh map` | Create an NSH forwarding map that binds an NSH SPI/SI pair to a next-hop action or interface. |
| `create pppoe session` | Create a PPPoE data-plane session binding a session-id and peer MAC to a decap/encap path. |
