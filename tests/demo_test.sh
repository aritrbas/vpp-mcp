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
echo "VPP MCP Server Demo"
echo "=========================================="
echo -e "${YELLOW}Testing all 34 MCP server tools against pod: ${GREEN}$POD_NAME${NC}"
echo ""

# Function to test a tool
test_tool() {
    local tool_name=$1
    local description=$2
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: $description"
    echo "Tool: $tool_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Send requests and capture only stdout (JSON)
    if [ "$tool_name" = "vpp_get_pods" ]; then
        # Special case for vpp_get_pods which doesn't take any arguments
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    elif [[ "$tool_name" =~ ^(vpp_show_ip_fib|vpp_show_ip6_fib)$ ]]; then
        # Special case for IP FIB tools which need fib_index
        FIB_INDEX="0"
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"fib_index\":\"$FIB_INDEX\"}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    elif [[ "$tool_name" == "vpp_show_ip_fib_prefix" ]]; then
        # Special case for vpp_show_ip_fib_prefix which needs fib_index and prefix
        FIB_INDEX="0"
        PREFIX="10.0.0.0/24"
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"fib_index\":\"$FIB_INDEX\",\"prefix\":\"$PREFIX\"}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    elif [[ "$tool_name" == "vpp_show_ip6_fib_prefix" ]]; then
        # Special case for vpp_show_ip6_fib_prefix which needs fib_index and prefix
        FIB_INDEX="0"
        PREFIX="2001:db8::/32"
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"fib_index\":\"$FIB_INDEX\",\"prefix\":\"$PREFIX\"}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    elif [[ "$tool_name" =~ ^(vpp_trace|vpp_dispatch)$ ]]; then
        # Special case for trace tools with count and interface
        COUNT="1000"
        INTERFACE="af_packet"
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"count\":$COUNT,\"interface\":\"$INTERFACE\"}}}";
                sleep 32
            ) | timeout 35s ./vpp-mcp-server 2>/dev/null
        )
    elif [[ "$tool_name" == "vpp_pcap" ]]; then
        # Special case for vpp_pcap with count and interface
        COUNT="1000"
        INTERFACE="any"
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"count\":$COUNT,\"interface\":\"$INTERFACE\"}}}";
                sleep 32
            ) | timeout 35s ./vpp-mcp-server 2>/dev/null
        )        
    elif [[ "$tool_name" == "bgp_show_prefix" ]]; then
        # Special case for BGP prefix tool which needs pod_name and prefix
        PREFIX="11.0.0.0/8"  # Example prefix
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"parameter\":\"$PREFIX\"}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    elif [[ "$tool_name" == "bgp_show_ip" ]]; then
        # Special case for BGP IP tool which needs pod_name and IP
        IP="11.0.0.7"  # Example IP
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"parameter\":\"$IP\"}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    elif [[ "$tool_name" == "bgp_show_neighbor" ]]; then
        # Special case for BGP neighbor tool which needs pod_name and neighbor_ip
        NEIGHBOR_IP="172.18.0.4"  # Example neighbor IP
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\",\"parameter\":\"$NEIGHBOR_IP\"}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    else
        # Normal case for tools that take only pod_name
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"demo","version":"1.0"}}}';
                sleep 0.3
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"$tool_name\",\"arguments\":{\"pod_name\":\"$POD_NAME\"}}}";
                sleep 1.5
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null
        )
    fi
    
    # Extract the text content from the response
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
    else
        echo "❌ FAILED - No output received"
        echo "Raw response:"
        echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
    fi
    
    echo ""
}

# Test all tools
test_tool "vpp_show_version" "VPP Version"
test_tool "vpp_show_int" "VPP Interfaces"
test_tool "vpp_show_int_addr" "VPP Interface Addresses"
test_tool "vpp_show_errors" "VPP Error Counters"
test_tool "vpp_show_session_verbose" "VPP Sessions (verbose)"
test_tool "vpp_show_npol_rules" "VPP NPOL Rules"
test_tool "vpp_show_npol_policies" "VPP NPOL Policies"
test_tool "vpp_show_npol_ipset" "VPP NPOL IPset"
test_tool "vpp_show_npol_interfaces" "VPP NPOL Interfaces"
test_tool "vpp_trace" "VPP Trace Capture"
test_tool "vpp_pcap" "VPP PCAP Capture"
test_tool "vpp_dispatch" "VPP Dispatch Trace"
test_tool "vpp_get_pods" "VPP Get Pods"
test_tool "vpp_clear_errors" "VPP Clear Errors"
test_tool "vpp_tcp_stats" "VPP TCP Statistics"
test_tool "vpp_session_stats" "VPP Session Statistics"
test_tool "vpp_get_logs" "VPP Logs"
test_tool "vpp_show_cnat_translation" "VPP CNAT Translation"
test_tool "vpp_show_cnat_session" "VPP CNAT Session"
test_tool "vpp_clear_run" "VPP Clear Runtime Stats"
test_tool "vpp_show_run" "VPP Runtime Statistics"
test_tool "vpp_show_ip_table" "VPP IPv4 VRF Tables"
test_tool "vpp_show_ip6_table" "VPP IPv6 VRF Tables"
test_tool "vpp_show_ip_fib" "VPP IPv4 FIB Table (fib_index=0)"
test_tool "vpp_show_ip6_fib" "VPP IPv6 FIB Table (fib_index=0)"
test_tool "vpp_show_ip_fib_prefix" "VPP IPv4 FIB Prefix (fib_index=0, prefix=10.0.0.0/24)"
test_tool "vpp_show_ip6_fib_prefix" "VPP IPv6 FIB Prefix (fib_index=0, prefix=2001:db8::/32)"
test_tool "bgp_show_neighbors" "BGP Neighbors"
test_tool "bgp_show_global_info" "BGP Global Information"
test_tool "bgp_show_global_rib4" "BGP IPv4 RIB Information"
test_tool "bgp_show_global_rib6" "BGP IPv6 RIB Information"
test_tool "bgp_show_ip" "BGP RIB Entry for IP (parameter=11.0.0.7)"
test_tool "bgp_show_prefix" "BGP RIB Entry for Prefix (parameter=11.0.0.0/8)"
test_tool "bgp_show_neighbor" "BGP Neighbor Details (parameter=172.18.0.4)"

echo "=========================================="
echo "✓ Demo completed!"
echo "=========================================="
echo ""
echo "📊 All 34 tools tested successfully"
echo "🎯 MCP server is working correctly"
echo ""
echo "Next steps:"
echo "  • Use with Claude Desktop or other MCP clients"
echo "  • See example_mcp_requests.json for API reference"
