# VPP `vppctl delete` Command Catalog

## Purpose
- This catalog groups all `delete ...` commands available in the VPP CLI (`vppctl`).
- `delete` commands destroy previously created data-plane objects: interfaces, tunnels, ACLs, and sockets.

## Category Index
| Category | Command Count |
|---|---:|
| Physical & Virtual Interfaces | 12 |
| Tunnels & Encapsulation | 2 |
| Security & ACL | 2 |
| Other Resources | 2 |

## Commands by Category

### Physical & Virtual Interfaces
Remove host, kernel, memory-shared, loopback, and paravirtualized interfaces.

| Command | Description |
|---|---|
| `delete host-interface` | Destroy an AF_PACKET host interface and release its Linux netdev binding. |
| `delete interface af_xdp` | Destroy an AF_XDP interface and detach its eBPF/XDP program from the netdev. |
| `delete interface memif` | Destroy a memif shared-memory interface and disconnect the peer. |
| `delete interface rdma` | Destroy an RDMA (ibverb/mlx5) interface and release its hardware queue resources. |
| `delete interface virtio` | Destroy a virtio-backed network interface and free associated vring descriptors. |
| `delete interface vmxnet3` | Destroy a vmxnet3 paravirtualized interface and release PCI BAR mappings. |
| `delete loopback interface` | Remove a software loopback interface and release its sw_if_index. |
| `delete memif socket` | Remove a memif socket file entry (only when no interfaces reference it). |
| `delete packet-generator interface` | Remove a virtual packet-generator (pg) interface from the forwarding graph. |
| `delete sub-interface` | Remove a VLAN sub-interface from its parent and release its sw_if_index. |
| `delete tap` | Destroy a TUN/TAP interface and remove the associated Linux netdev. |
| `delete vhost-user` | Destroy a vhost-user interface and close the Unix domain socket connection. |

### Tunnels & Encapsulation
Remove IP-in-IP and 6rd tunnel endpoints.

| Command | Description |
|---|---|
| `delete 6rd tunnel` | Tear down a 6rd (IPv6 Rapid Deployment) tunnel and remove its encap/decap paths. |
| `delete ipip tunnel` | Tear down an IP-in-IP tunnel and remove its adjacency and FIB entries. |

### Security & ACL
Delete ACL and MACIP ACL rule sets.

| Command | Description |
|---|---|
| `delete acl-plugin acl` | Delete an L3/L4 ACL rule set by index and remove it from all interface bindings. |
| `delete acl-plugin macip acl` | Delete a MAC+IP ACL rule set by index and detach it from bound interfaces. |

### Other Resources
Remove TEIB entries and bonding groups.

| Command | Description |
|---|---|
| `delete bond` | Delete a bonding master interface and release all enslaved member interfaces. |
| `delete teib` | Remove a Tunnel Endpoint Information Base entry and its associated underlay next-hop. |
