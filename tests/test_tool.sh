#!/bin/bash

# Quick tool testing script for VPP MCP Server (4-tool architecture)
# Usage:
#   ./test_tool.sh <pod_name> vppctl "<command>"
#   ./test_tool.sh <pod_name> gobgp "<command>"
#   ./test_tool.sh cluster <command> [resource_name]
#   ./test_tool.sh list

NO_POD_TOOLS_REGEX='^(list|cluster)$'

if [[ "$1" =~ $NO_POD_TOOLS_REGEX ]]; then
    POD_NAME=""
    TOOL_NAME="$1"
    PARAMETER1="$2"
    PARAMETER2="$3"
else
    POD_NAME="$1"
    TOOL_NAME="$2"
    PARAMETER1="$3"
    PARAMETER2="$4"
fi

if [ -z "$TOOL_NAME" ] || { [ -z "$POD_NAME" ] && ! [[ "$TOOL_NAME" =~ $NO_POD_TOOLS_REGEX ]]; }; then
    echo "Usage:"
    echo "  $0 <pod_name> vppctl \"<command>\"          # Run vppctl command"
    echo "  $0 <pod_name> gobgp \"<command>\"            # Run gobgp command"
    echo "  $0 cluster <command> [resource_name]        # Run cluster command"
    echo "  $0 list                                     # List all commands"
    echo ""
    echo "Examples:"
    echo "  $0 my-pod vppctl \"show version\""
    echo "  $0 my-pod vppctl \"show ip fib index 0\""
    echo "  $0 my-pod vppctl trace                      # capture (uses count/timeout defaults)"
    echo "  $0 my-pod gobgp neighbor"
    echo "  $0 my-pod gobgp \"global rib -a ipv4\""
    echo "  $0 cluster get_pods"
    echo "  $0 cluster logs my-pod                      # logs for pod (container default: agent)"
    echo "  $0 cluster get_configmap"
    echo "  $0 list"
    echo ""
    echo "Available tools: list, cluster, vppctl, gobgp"
    exit 1
fi

if [ -n "$POD_NAME" ]; then
    echo "Pod: $POD_NAME"
else
    echo "Pod: (not required for this tool)"
fi
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

# Build tool call based on tool type
case "$TOOL_NAME" in
    list)
        echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"list","arguments":{}}}' >> "$TEMP_REQUESTS"
        ;;
    cluster)
        COMMAND="${PARAMETER1:-get_pods}"
        if [[ "$COMMAND" == "logs" || "$COMMAND" == "describe_pod" ]]; then
            # These need pod_name from PARAMETER2
            CLUSTER_POD="${PARAMETER2:-}"
            if [ -z "$CLUSTER_POD" ]; then
                echo "Error: '$COMMAND' requires a pod name as the next argument"
                exit 1
            fi
            echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"cluster\",\"arguments\":{\"command\":\"$COMMAND\",\"pod_name\":\"$CLUSTER_POD\"}}}" >> "$TEMP_REQUESTS"
        elif [[ "$COMMAND" == "exec" ]]; then
            # exec needs pod_name and command (resource_name)
            CLUSTER_POD="${PARAMETER2:-}"
            if [ -z "$CLUSTER_POD" ]; then
                echo "Error: 'exec' requires a pod name as the next argument"
                exit 1
            fi
            # For exec, we expect a third parameter as the command
            shift 3  # Skip script name, 'cluster', 'exec', pod_name
            EXEC_COMMAND="$@"
            if [ -z "$EXEC_COMMAND" ]; then
                EXEC_COMMAND="ls"  # default
            fi
            echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"cluster\",\"arguments\":{\"command\":\"exec\",\"pod_name\":\"$CLUSTER_POD\",\"resource_name\":\"$EXEC_COMMAND\"}}}" >> "$TEMP_REQUESTS"
        elif [[ "$COMMAND" == "describe_node" ]]; then
            NODE_NAME="${PARAMETER2:-}"
            if [ -z "$NODE_NAME" ]; then
                echo "Error: 'describe_node' requires a node name as the next argument"
                exit 1
            fi
            echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"cluster\",\"arguments\":{\"command\":\"$COMMAND\",\"resource_name\":\"$NODE_NAME\"}}}" >> "$TEMP_REQUESTS"
        elif [[ "$COMMAND" == "get_endpoints" ]]; then
            # get_endpoints requires resource_name (service name)
            RESOURCE="${PARAMETER2:-}"
            if [ -z "$RESOURCE" ]; then
                echo "Error: 'get_endpoints' requires a service name as the next argument"
                exit 1
            fi
            echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"cluster\",\"arguments\":{\"command\":\"$COMMAND\",\"resource_name\":\"$RESOURCE\"}}}" >> "$TEMP_REQUESTS"
        elif [[ "$COMMAND" == "get_configmap" || "$COMMAND" == "get_daemonset" ]]; then
            RESOURCE="${PARAMETER2:-}"
            if [ -n "$RESOURCE" ]; then
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"cluster\",\"arguments\":{\"command\":\"$COMMAND\",\"resource_name\":\"$RESOURCE\"}}}" >> "$TEMP_REQUESTS"
            else
                echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"cluster\",\"arguments\":{\"command\":\"$COMMAND\"}}}" >> "$TEMP_REQUESTS"
            fi
        else
            echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"cluster\",\"arguments\":{\"command\":\"$COMMAND\"}}}" >> "$TEMP_REQUESTS"
        fi
        ;;
    vppctl)
        COMMAND="${PARAMETER1:-show version}"
        # Check if it's a capture command that needs longer timeout
        LOWER_CMD=$(echo "$COMMAND" | tr '[:upper:]' '[:lower:]')
        echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"vppctl\",\"arguments\":{\"command\":\"$COMMAND\",\"pod_name\":\"$POD_NAME\"}}}" >> "$TEMP_REQUESTS"
        ;;
    gobgp)
        COMMAND="${PARAMETER1:-neighbor}"
        echo "{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"gobgp\",\"arguments\":{\"command\":\"$COMMAND\",\"pod_name\":\"$POD_NAME\"}}}" >> "$TEMP_REQUESTS"
        ;;
    *)
        echo "Unknown tool: $TOOL_NAME"
        echo "Available tools: list, cluster, vppctl, gobgp"
        exit 1
        ;;
esac

# Execute the server with requests
echo "Executing MCP server..."
echo ""

# Create a temporary file for output
TEMP_OUTPUT=$(mktemp)

# Set timeout based on command type (capture commands need more time)
TIMEOUT="5s"
if [ "$TOOL_NAME" = "vppctl" ]; then
    LOWER_CMD=$(echo "${PARAMETER1:-}" | tr '[:upper:]' '[:lower:]')
    if [[ "$LOWER_CMD" =~ ^(trace|pcap|dispatch) ]]; then
        TIMEOUT="35s"
    fi
fi

# Stream requests with short delays so initialize and tools/call are processed reliably
{
    while IFS= read -r request_line; do
        echo "$request_line"
        sleep 0.4
    done < "$TEMP_REQUESTS"

    # Keep stdin open long enough for command execution and response emission
    if [ "$TIMEOUT" = "35s" ]; then
        sleep 32
    else
        sleep 1.5
    fi
} | timeout $TIMEOUT ./vpp-mcp-server > "$TEMP_OUTPUT" 2>&1

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
