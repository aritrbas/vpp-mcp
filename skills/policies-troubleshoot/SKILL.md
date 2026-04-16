---
name: policies-troubleshoot
description: CalicoVPP network policy troubleshooting workflow mapped to upstream policies/troubleshooting.md
---

# CalicoVPP Policies Troubleshooting Workflow

This runbook maps directly to the upstream guide:
- https://github.com/projectcalico/vpp-dataplane/blob/master/docs/policies/troubleshooting.md

It is tailored for this MCP server and for CalicoVPP where VPP runs in Kubernetes pods.

## Scope

Use this workflow when traffic is unexpectedly blocked or permitted — suspected policy misconfiguration, missing rules, or incorrect ipset membership.

## Step 1: Discover target pods

Call `cluster` with `command=get_pods` to list all CalicoVPP pods.

Pick the node where the affected workload is running.

## Step 2: Inspect per-interface policy configuration

Call `vppctl` with:
- `command=show npol interfaces`, `pod_name=<pod-name>`

This lists every interface with policies applied and their configured rules. Key fields:
- Interface name + `sw_if_index` + first IPv4 address (to identify which pod the interface belongs to)
- `tx` section: rules applied to packets **leaving** VPP on this interface (top to bottom)
- `rx` section: rules applied to packets **entering** VPP on this interface (top to bottom)
- `profiles` section: rules enforced when a matched rule action is `PASS` or when no policies are configured

Validate:
- Expected interfaces are listed (tuntap `tun<N>`, `tap0` for host)
- Correct policies are attached to the right interfaces
- `tx`/`rx` rule ordering is as intended

## Step 3: Inspect all active policies

Call `vppctl` with:
- `command=show npol policies verbose`, `pod_name=<pod-name>`

This lists all policies referenced across interfaces, with their full rule sets. Use this to:
- Confirm the content of each `policy#N` seen in Step 2
- Verify that expected allow/deny rules are present
- Check for unexpected rules that may be causing unintended behaviour

## Step 4: Inspect individual rules

Call `vppctl` with:
- `command=show npol rules`, `pod_name=<pod-name>`

This lists all rules referenced by policies. Validate:
- Rule actions are correct (`allow` vs `deny`)
- Match criteria (protocol, port, source/destination) match intent
- Rules referencing ipsets (`ipset#N`) reference the expected sets (cross-check with Step 5)

## Step 5: Inspect ipsets

Call `vppctl` with:
- `command=show npol ipset`, `pod_name=<pod-name>`

IPsets are lists of IP addresses referenced by rules. Validate:
- The expected pod or node IPs are members of each ipset
- No stale or incorrect IPs are present

Common issue: a pod IP that should be allowed is missing from its ipset — this causes traffic to hit a default-deny rule instead.

## Step 6: Correlate with container logs if policies look wrong

If the VPP policy state does not match the expected Kubernetes NetworkPolicy intent:

Call `cluster` with:
- `command=logs`, `pod_name=<pod-name>`, `container=agent`, `tail_lines=300`

Look for:
- Policy programming errors or reconciliation failures
- Kubernetes watch/list errors preventing policy updates

## Fast Triage Order

1. `cluster` → `get_pods` (discover all pods)
2. `vppctl` → `show npol interfaces` (identify which policies apply to the affected interface)
3. `vppctl` → `show npol policies verbose` (inspect the content of those policies)
4. `vppctl` → `show npol rules` (verify rule match criteria and actions)
5. `vppctl` → `show npol ipset` (verify ipset membership for rules that reference ipsets)
6. `cluster` → `logs` with `container=agent` (if VPP state diverges from Kubernetes policy intent)
