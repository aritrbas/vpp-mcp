#!/bin/bash

# Quick tool testing script for VPP MCP Server
# Usage: ./test_tool.sh <tool_name> <pod_name> [namespace]

TOOL_NAME="$1"
POD_NAME="$2"
NAMESPACE="${3:-calico-vpp-dataplane}"

if [ -z "$TOOL_NAME" ] || [ -z "$POD_NAME" ]; then
    echo "Usage: $0 <tool_name> <pod_name> [namespace]"
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
    echo ""
    echo "Example:"
    echo "  $0 vpp_show_version calico-vpp-node-xxxxx"
    exit 1
fi

echo "Testing tool: $TOOL_NAME"
echo "Pod: $POD_NAME"
echo "Namespace: $NAMESPACE"
echo ""

# Create temporary file for requests
TEMP_REQUESTS=$(mktemp)

# Write requests to temp file
cat > "$TEMP_REQUESTS" << EOF
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"$TOOL_NAME","arguments":{"pod_name":"$POD_NAME","namespace":"$NAMESPACE"}}}
EOF

# Execute the server with requests
echo "Executing MCP server..."
echo ""

timeout 10s ./vpp-mcp-server < "$TEMP_REQUESTS" 2>&1 | while IFS= read -r line; do
    # Try to parse as JSON and pretty print
    if echo "$line" | jq -e . >/dev/null 2>&1; then
        # Check if it's a tool result
        if echo "$line" | jq -e '.result.content[]?.text' >/dev/null 2>&1; then
            echo "=== Tool Result ==="
            echo "$line" | jq -r '.result.content[].text'
            echo ""
            echo "=== Raw JSON ==="
            echo "$line" | jq '.'
        else
            # Other JSON responses
            echo "$line" | jq '.'
        fi
    else
        # Non-JSON output (logs)
        echo "$line"
    fi
done

# Clean up
rm -f "$TEMP_REQUESTS"

echo ""
echo "Test completed."
