#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Use command line argument if provided, otherwise try to find a pod automatically
if [ -n "$1" ]; then
    POD_NAME="$1"
    echo -e "${YELLOW}Using specified pod: $POD_NAME${NC}"
else
    # Check if kubectl is available
    if command -v kubectl &> /dev/null; then
        echo -e "${GREEN}✓ kubectl is installed${NC}"
        
        # Check if we can access any VPP pods
        POD_NAME=$(kubectl get pods -n calico-vpp-dataplane -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
        
        if [ -n "$POD_NAME" ]; then
            echo -e "${GREEN}✓ Found VPP pod: $POD_NAME${NC}"
        else
            echo -e "${RED}✗ No VPP pods found in calico-vpp-dataplane namespace${NC}"
            echo "Please specify a pod name as an argument:"
            echo "  $0 <pod-name>"
            exit 1
        fi
    else
        echo -e "${RED}✗ kubectl not found${NC}"
        echo "Please specify a pod name as an argument:"
        echo "  $0 <pod-name>"
        exit 1
    fi
fi

echo "=========================================="
echo "VPP MCP Server Demo (4-Tool Architecture)"
echo "=========================================="
echo -e "${YELLOW}Testing all 4 MCP tools with representative commands against pod: ${GREEN}$POD_NAME${NC}"
echo ""

PASS_COUNT=0
FAIL_COUNT=0
TEST_COUNT=0

# Function to call a tool with given arguments JSON
call_tool() {
    local tool_name=$1
    local arguments=$2
    local description=$3
    local timeout_val=${4:-5s}
    local sleep_val=${5:-1.5}
    
    TEST_COUNT=$((TEST_COUNT + 1))
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Test $TEST_COUNT: $description"
    echo "Tool: $tool_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    RESULT=$(
        (
            echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
            sleep 0.3
            echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":$arguments}}";
            sleep $sleep_val
        ) | timeout $timeout_val ./vpp-mcp-server 2>/dev/null
    )
    
    OUTPUT=$(echo "$RESULT" | jq -r 'select(.id==2) | .result.content[].text' 2>/dev/null || echo "")
    
    if [ -n "$OUTPUT" ]; then
        echo "$OUTPUT" | head -30
        LINE_COUNT=$(echo "$OUTPUT" | wc -l)
        if [ "$LINE_COUNT" -gt 30 ]; then
            echo ""
            echo "... ($(($LINE_COUNT - 30)) more lines)"
        fi
        echo ""
        echo "✅ SUCCESS"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo "❌ FAILED - No output received"
        echo "Raw response:"
        echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    echo ""
}

# =========================================================================
# TOOL 1: list
# =========================================================================
echo -e "${YELLOW}=== Testing: list tool ===${NC}"
call_tool "list" "{}" "List all available commands"

# =========================================================================
# TOOL 2: cluster
# =========================================================================
echo -e "${YELLOW}=== Testing: cluster tool ===${NC}"
call_tool "cluster" "{\"command\":\"get_pods\"}" "cluster: get_pods"
call_tool "cluster" "{\"command\":\"get_nodes\"}" "cluster: get_nodes"
call_tool "cluster" "{\"command\":\"get_namespaces\"}" "cluster: get_namespaces"
call_tool "cluster" "{\"command\":\"get_configmap\"}" "cluster: get_configmap (default: calico-vpp-config)"
call_tool "cluster" "{\"command\":\"get_daemonset\",\"resource_name\":\"calico-vpp-node\"}" "cluster: get_daemonset"
call_tool "cluster" "{\"command\":\"get_events\"}" "cluster: get_events"
call_tool "cluster" "{\"command\":\"get_services\"}" "cluster: get_services"
call_tool "cluster" "{\"command\":\"top_nodes\"}" "cluster: top_nodes"
call_tool "cluster" "{\"command\":\"top_pods\"}" "cluster: top_pods"
call_tool "cluster" "{\"command\":\"logs\",\"pod_name\":\"$POD_NAME\",\"container\":\"agent\",\"tail_lines\":50}" "cluster: logs (agent container)"
call_tool "cluster" "{\"command\":\"logs\",\"pod_name\":\"$POD_NAME\",\"container\":\"vpp\",\"tail_lines\":50}" "cluster: logs (vpp container)"
call_tool "cluster" "{\"command\":\"describe_pod\",\"pod_name\":\"$POD_NAME\"}" "cluster: describe_pod"
call_tool "cluster" "{\"command\":\"exec\",\"pod_name\":\"$POD_NAME\",\"container\":\"vpp\",\"resource_name\":\"ip a\"}" "cluster: exec (read-only command)"

# =========================================================================
# TOOL 3: vppctl
# =========================================================================
echo -e "${YELLOW}=== Testing: vppctl tool ===${NC}"
call_tool "vppctl" "{\"command\":\"show version\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show version"
call_tool "vppctl" "{\"command\":\"show int\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show int"
call_tool "vppctl" "{\"command\":\"show int addr\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show int addr"
call_tool "vppctl" "{\"command\":\"show hardware-interfaces\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show hardware-interfaces"
call_tool "vppctl" "{\"command\":\"show errors\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show errors"
call_tool "vppctl" "{\"command\":\"show session verbose 2\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show session verbose 2"
call_tool "vppctl" "{\"command\":\"show npol rules\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show npol rules"
call_tool "vppctl" "{\"command\":\"show npol policies\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show npol policies"
call_tool "vppctl" "{\"command\":\"show npol ipset\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show npol ipset"
call_tool "vppctl" "{\"command\":\"show npol interfaces\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show npol interfaces"
call_tool "vppctl" "{\"command\":\"show tcp stats\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show tcp stats"
call_tool "vppctl" "{\"command\":\"show session stats\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show session stats"
call_tool "vppctl" "{\"command\":\"show logging\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show logging"
call_tool "vppctl" "{\"command\":\"show cnat translation\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show cnat translation"
call_tool "vppctl" "{\"command\":\"show cnat session\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show cnat session"
call_tool "vppctl" "{\"command\":\"show run\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show run"
call_tool "vppctl" "{\"command\":\"show ipip tunnel\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show ipip tunnel"
call_tool "vppctl" "{\"command\":\"show vxlan tunnel\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show vxlan tunnel"
call_tool "vppctl" "{\"command\":\"show tun\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show tun"
call_tool "vppctl" "{\"command\":\"show ip table\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show ip table"
call_tool "vppctl" "{\"command\":\"show ip6 table\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show ip6 table"
call_tool "vppctl" "{\"command\":\"show ip fib index 0\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show ip fib index 0"
call_tool "vppctl" "{\"command\":\"show ip6 fib index 0\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show ip6 fib index 0"
call_tool "vppctl" "{\"command\":\"show ip fib index 0 10.0.0.0/24\",\"pod_name\":\"$POD_NAME\"}" "vppctl: show ip fib index 0 10.0.0.0/24"
call_tool "vppctl" "{\"command\":\"capture_cleanup\",\"pod_name\":\"$POD_NAME\"}" "vppctl: capture_cleanup"

# Capture tools (longer timeout)
call_tool "vppctl" "{\"command\":\"trace\",\"pod_name\":\"$POD_NAME\",\"count\":1000,\"interface\":\"af_packet\"}" "vppctl: trace capture" "35s" "32"
call_tool "vppctl" "{\"command\":\"pcap\",\"pod_name\":\"$POD_NAME\",\"count\":1000,\"interface\":\"any\"}" "vppctl: pcap capture" "35s" "32"
call_tool "vppctl" "{\"command\":\"dispatch\",\"pod_name\":\"$POD_NAME\",\"count\":1000,\"interface\":\"af_packet\"}" "vppctl: dispatch capture" "35s" "32"
call_tool "vppctl" "{\"command\":\"capture_cleanup\",\"pod_name\":\"$POD_NAME\"}" "vppctl: capture_cleanup (post-capture)"

# =========================================================================
# TOOL 4: gobgp
# =========================================================================
echo -e "${YELLOW}=== Testing: gobgp tool ===${NC}"
call_tool "gobgp" "{\"command\":\"neighbor\",\"pod_name\":\"$POD_NAME\"}" "gobgp: neighbor"
call_tool "gobgp" "{\"command\":\"global\",\"pod_name\":\"$POD_NAME\"}" "gobgp: global"
call_tool "gobgp" "{\"command\":\"global rib -a ipv4\",\"pod_name\":\"$POD_NAME\"}" "gobgp: global rib -a ipv4"
call_tool "gobgp" "{\"command\":\"global rib -a ipv6\",\"pod_name\":\"$POD_NAME\"}" "gobgp: global rib -a ipv6"
call_tool "gobgp" "{\"command\":\"global rib 11.0.0.0/8\",\"pod_name\":\"$POD_NAME\"}" "gobgp: global rib 11.0.0.0/8"

echo "=========================================="
echo "✓ Demo completed!"
echo "=========================================="
echo ""
echo "Results: $PASS_COUNT passed, $FAIL_COUNT failed out of $TEST_COUNT tests"
echo ""
echo "Next steps:"
echo "  - Use with Claude Desktop, Windsurf, or other MCP clients"
echo "  - See examples/example_mcp_requests.json for API reference"
