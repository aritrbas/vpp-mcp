---
name: calicovpp-troubleshooting
description: Dedicated CalicoVPP troubleshooting workflow aligned with upstream docs/troubleshooting.md
---

# CalicoVPP General Troubleshooting Workflow

This workflow maps to the upstream guide:
- https://github.com/projectcalico/vpp-dataplane/blob/master/docs/troubleshooting.md

It is tailored for CalicoVPP container deployments where VPP runs in the `vpp` container of pods in namespace `calico-vpp-dataplane`.

## Step 1: Discover target pods

Use:
- `vpp_get_pods`

Pick the affected pod(s) first, then expand to all nodes if needed.

## Step 2: Run the baseline sequence using existing tools

For each target pod, run these tools in order:

1. `vpp_clear_run`
2. Wait 3-10 seconds while the issue condition exists
3. `vpp_show_run`
4. `vpp_show_errors`
5. `vpp_tcp_stats`
6. `vpp_session_stats`
7. `vpp_get_logs`

Use:
- `vpp_clear_errors` only when you intentionally want to reset error counters before a new measurement window

## Step 3: Process output into a troubleshooting summary

From `vpp_show_run`:
- Extract `loops/sec`
- Extract highest `Vectors/Call`
- List nodes where `Clocks > 1e3`
- Mark pod as "load-sensitive" if vectors are high and loops/sec is low

From `vpp_show_errors`:
- List top non-zero counters (top 10)
- On repeated runs, compare counter deltas and keep only increasing counters as active signals
- Ignore one-time static counters unless they keep increasing

From `vpp_tcp_stats` and `vpp_session_stats`:
- Keep only non-zero lines
- Highlight timeout, reset, and listen-conflict lines

From `vpp_get_logs`:
- Keep lines matching `error|fatal|panic|failed|timeout|watch|kubernetes|bgp|session`
- Deduplicate repeating lines so recurring loops are visible

## Suggested triage order

1. `vpp_get_pods`
2. Run Step 2 sequence on one affected pod
3. Run Step 2 sequence on peers for comparison
4. Compare processed summaries to isolate node-specific vs cluster-wide issues

## Output template

Use this per pod:

- Pod:
- Time window:
- `show run`: loops/sec=?, max vectors/call=?, high-clocks nodes=?
- `show errors`: top increasing counters=?
- `show tcp stats`: non-zero signals=?
- `show session stats`: non-zero signals=?
- `show logging`: critical/repeating lines=?
- Verdict: healthy / degraded / critical
