#!/bin/bash

# Quick tool testing script for VPP MCP Server
# Usage: ./test_tool.sh <pod_name> <tool_name> [parameter1] [parameter2]

POD_NAME="$1"
TOOL_NAME="$2"
PARAMETER1="$3"
PARAMETER2="$4"

if [ -z "$POD_NAME" ] || [ -z "$TOOL_NAME" ]; then
    echo "Usage: $0 <pod_name> <tool_name> [parameter1] [parameter2]"
    echo ""
    echo "Most tools only need pod_name. Some have optional/required parameters:"
    echo "- vpp_show_ip_fib: optional fib_index (default: 0)"
    echo "- vpp_show_ip6_fib: optional fib_index (default: 0)"
    echo "- vpp_show_ip_fib_prefix: optional 'fib_index prefix' (default: '0 10.0.0.0/24')"
    echo "- vpp_show_ip6_fib_prefix: optional 'fib_index prefix' (default: '0 2001:db8::/32')"
    echo "- vpp_trace: optional 'count interface' (default: 1000 af-packet)"
    echo "- vpp_pcap: optional 'count interface' (default: 1000 any)"
    echo "- vpp_dispatch: optional 'count interface' (default: 1000 af-packet)"        
    echo "- bgp_show_ip: optional IP (default: 11.0.0.7)"
    echo "- bgp_show_prefix: optional prefix (default: 11.0.0.0/8)"
    echo "- bgp_show_neighbor: optional neighbor IP (default: 172.18.0.4)"
    echo ""
    echo "Available tools:"
    echo "  - vpp_show_version"
    echo "  - vpp_show_int"
    echo "  - vpp_show_int_addr"
    echo "  - vpp_show_errors"
    echo "  - vpp_show_session_verbose"
    echo "  - vpp_show_npol_rules"
    echo "  - vpp_show_npol_policies"
    echo "  - vpp_show_npol_ipset"
    echo "  - vpp_show_npol_interfaces"
    echo "  - vpp_trace"
    echo "  - vpp_pcap"
    echo "  - vpp_dispatch"
    echo "  - vpp_get_pods"
    echo "  - vpp_clear_errors"
    echo "  - vpp_tcp_stats"
    echo "  - vpp_session_stats"
    echo "  - vpp_get_logs"
    echo "  - vpp_show_cnat_translation"
    echo "  - vpp_show_cnat_session"
    echo "  - vpp_clear_run"
    echo "  - vpp_show_run"
    echo "  - vpp_show_ip_table"
    echo "  - vpp_show_ip6_table"
    echo "  - vpp_show_ip_fib"
    echo "  - vpp_show_ip6_fib"
    echo "  - vpp_show_ip_fib_prefix"
    echo "  - vpp_show_ip6_fib_prefix"
    echo "  - bgp_show_neighbors"
    echo "  - bgp_show_global_info"
    echo "  - bgp_show_global_rib4"
    echo "  - bgp_show_global_rib6"
    echo "  - bgp_show_ip"
    echo "  - bgp_show_prefix"
    echo "  - bgp_show_neighbor"
    exit 1
fi

echo "Pod: $POD_NAME"
echo "Tool: $TOOL_NAME"
if [ -n "$PARAMETER1" ]; then
    echo "Parameter 1: $PARAMETER1"
fi
if [ -n "$PARAMETER2" ]; then
    echo "Parameter 2: $PARAMETER2"
fi
echo ""

# Create temporary file for requests
TEMP_REQUESTS=$(mktemp)

# Write initialize request to temp file
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' > "$TEMP_REQUESTS"

# Handle different tool types for the second request
if [ "$TOOL_NAME" = "vpp_get_pods" ]; then
    # vpp_get_pods doesn't need a pod_name
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"'$TOOL_NAME'","arguments":{}}}' >> "$TEMP_REQUESTS"
elif [[ "$TOOL_NAME" =~ ^(bgp_show_ip|bgp_show_prefix|bgp_show_neighbor)$ ]]; then
    # BGP tools that need a parameter (use provided or defaults)
    if [ -n "$PARAMETER1" ]; then
        PARAM_VALUE="$PARAMETER1"
    else
        # Use defaults
        case "$TOOL_NAME" in
            "bgp_show_ip")
                PARAM_VALUE="11.0.0.7"
                ;;
            "bgp_show_prefix")
                PARAM_VALUE="11.0.0.0/8"
                ;;
            "bgp_show_neighbor")
                PARAM_VALUE="172.18.0.4"
                ;;
        esac
    fi
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"'$TOOL_NAME'","arguments":{"pod_name":"'$POD_NAME'","parameter":"'$PARAM_VALUE'"}}}' >> "$TEMP_REQUESTS"
elif [[ "$TOOL_NAME" =~ ^(vpp_show_ip_fib|vpp_show_ip6_fib)$ ]]; then
    # FIB tools that need fib_index (use parameter if provided, otherwise default to 0)
    FIB_INDEX="${PARAMETER1:-0}"
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"'$TOOL_NAME'","arguments":{"pod_name":"'$POD_NAME'","fib_index":"'$FIB_INDEX'"}}}' >> "$TEMP_REQUESTS"
elif [[ "$TOOL_NAME" =~ ^(vpp_show_ip_fib_prefix|vpp_show_ip6_fib_prefix)$ ]]; then
    # FIB prefix tools that need fib_index and prefix
    # Check if both parameters are provided
    if [ -n "$PARAMETER1" ] && [ -n "$PARAMETER2" ]; then
        FIB_INDEX="$PARAMETER1"
        PREFIX="$PARAMETER2"
    # Check if only one parameter is provided (assume it's a combined string)
    elif [ -n "$PARAMETER1" ] && [[ "$PARAMETER1" == *" "* ]]; then
        # Extract fib_index and prefix from parameter (format: "fib_index prefix")
        FIB_INDEX=$(echo "$PARAMETER1" | cut -d' ' -f1)
        PREFIX=$(echo "$PARAMETER1" | cut -d' ' -f2-)
    # Use defaults
    else
        FIB_INDEX="${PARAMETER1:-0}"  # Use PARAMETER1 if provided, otherwise 0
        if [[ "$TOOL_NAME" == "vpp_show_ip_fib_prefix" ]]; then
            PREFIX="${PARAMETER2:-10.0.0.0/24}"
        else
            PREFIX="${PARAMETER2:-2001:db8::/32}"
        fi
    fi
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"'$TOOL_NAME'","arguments":{"pod_name":"'$POD_NAME'","fib_index":"'$FIB_INDEX'","prefix":"'$PREFIX'"}}}' >> "$TEMP_REQUESTS"
elif [[ "$TOOL_NAME" =~ ^(vpp_trace|vpp_pcap|vpp_dispatch)$ ]]; then
    # Capture tools with optional count and interface
    if [ -n "$PARAMETER1" ] && [ -n "$PARAMETER2" ]; then
        # Both parameters provided
        COUNT="$PARAMETER1"
        INTERFACE="$PARAMETER2"
    elif [ -n "$PARAMETER1" ] && [[ "$PARAMETER1" == *" "* ]]; then
        # Extract count and interface from parameter (format: "count interface")
        COUNT=$(echo "$PARAMETER1" | cut -d' ' -f1)
        INTERFACE=$(echo "$PARAMETER1" | cut -d' ' -f2-)
    else
        COUNT="${PARAMETER1:-1000}"
        # Set interface based on tool
        if [[ "$TOOL_NAME" == "vpp_pcap" ]]; then
            INTERFACE="any"
        else
            INTERFACE="af_packet"
        fi
    fi
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"'$TOOL_NAME'","arguments":{"pod_name":"'$POD_NAME'","count":'$COUNT',"interface":"'$INTERFACE'"}}}' >> "$TEMP_REQUESTS"
else
    # Standard case for most tools
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"'$TOOL_NAME'","arguments":{"pod_name":"'$POD_NAME'"}}}' >> "$TEMP_REQUESTS"
fi

# Execute the server with requests
echo "Executing MCP server..."
echo ""

# Create a temporary file for output
TEMP_OUTPUT=$(mktemp)

# Set timeout based on tool type (capture tools need more time)
if [[ "$TOOL_NAME" =~ ^(vpp_trace|vpp_pcap|vpp_dispatch)$ ]]; then
    TIMEOUT="35s"
else
    TIMEOUT="5s"
fi

# Run the MCP server and capture all output
timeout $TIMEOUT ./vpp-mcp-server < "$TEMP_REQUESTS" > "$TEMP_OUTPUT" 2>&1

# First print non-JSON lines (logs)
grep -v '^{' "$TEMP_OUTPUT" || true

# Then find and process JSON responses
while read -r line; do
    if [[ $line == {* ]]; then  # Only process lines starting with {
        # Try to parse as JSON
        if echo "$line" | jq -e . >/dev/null 2>&1; then
            # Check if it's a tool result
            if echo "$line" | jq -e '.result.content[]?.text' >/dev/null 2>&1; then
                echo ""
                echo "=== Tool Result ==="
                echo "$line" | jq -r '.result.content[].text'
                echo ""
            fi
        fi
    fi
done < "$TEMP_OUTPUT"

# Clean up temporary files
rm -f "$TEMP_OUTPUT" "$TEMP_REQUESTS"

echo ""
echo "Test completed."
