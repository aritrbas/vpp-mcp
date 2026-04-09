---
name: cluster-healthcheck
description: Daily VPP cluster health check for IPv6 troubleshooting
---

# VPP Cluster Daily Health Check

This workflow performs comprehensive health checks on your VPP cluster, focusing on IPv6 connectivity and BGP peering issues.

## Steps

### 1. Get list of VPP pods
Use `cluster` with `command=get_pods` to retrieve all calico-vpp pods with their IPs and node assignments.

### 2. Check VPP version and daemonset image
- Use `vppctl` with `command=show version` on each pod to verify VPP version consistency
- Use `cluster` with `command=get_daemonset`, `resource_name=calico-vpp-node` to check the DaemonSet YAML (includes container images)

### 3. Check IPv6 enablement on interfaces
Use `vppctl` with `command=show int addr` on each pod to verify:
- `host-eth0` has IPv6 addresses (fdab:504:82d8::/64)
- Loop interfaces have fd20::*/128 addresses
- Tunnel interfaces (tun1-N) have IPv6 configured

### 4. Check BGP peering status
Use `gobgp` with `command=neighbor` on each pod to verify:
- All expected peers show "Establ" state
- Both IPv4 and IPv6 peerings are established
- Routes are being received and accepted (#Received > 0, Accepted > 0)

### 5. Check for UNRESOLVED FIB entries

#### IPv4 FIB check:
- Use `vppctl` with `command=show ip fib index 0` on each pod
- Search output for "UNRESOLVED" entries
- All routes should have proper next-hops (via, glean, receive, drop)

#### IPv6 FIB check:
- Use `vppctl` with `command=show ip6 fib index 0` on each pod
- Search output for "UNRESOLVED" entries
- All routes should have proper next-hops resolved

### 6. Check IPIP tunnel status
Use `vppctl` with `command=show ipip tunnel` on each pod to verify:
- Tunnel instances are configured
- Source/destination IPs match node IPs
- Table IDs and sw-if-idx are correct
- No error flags present

### 7. Check VXLAN tunnel status
Use `vppctl` with `command=show vxlan tunnel` on each pod to verify:
- VXLAN tunnels use IPv6 addresses (fdab:504:82d8::*)
- VNI is consistent (typically 4096)
- Source/destination ports are 4789
- FIB indices are correct

## Manual Debug Commands Reference

If issues are found, you can exec into pods for detailed debugging:

**Agent container (BGP):**
```bash
kubectl exec -n calico-vpp-dataplane <pod-name> -c agent -- gobgp neighbor
```

**VPP container:**
```bash
kubectl exec -n calico-vpp-dataplane <pod-name> -c vpp -- bash -c "ip a && vppctl sh ipip tunnel && vppctl sh ip6 fib | grep UNRESOLVED"
```

## Expected Healthy Output

- **Pods**: All pods Running 2/2
- **VPP Version**: Consistent across all nodes
- **IPv6**: All interfaces have IPv6 addresses assigned
- **BGP**: Full mesh established (N-1 peers per node for N nodes, both IPv4+IPv6)
- **FIB**: No UNRESOLVED entries in either IPv4 or IPv6 tables
- **IPIP Tunnels**: One tunnel per remote node, proper src/dst IPs
- **VXLAN Tunnels**: IPv6-based, VNI 4096, port 4789

## Common Issues

- **BGP peers in "Active" or "Idle"**: Check kubernetes API access, firewall rules
- **UNRESOLVED FIB entries**: Indicates routing issues, check BGP route exchange
- **Missing IPv6 addresses**: Check node configuration, CNI settings
- **Tunnel mismatches**: Verify node IPs, check for MTU issues
