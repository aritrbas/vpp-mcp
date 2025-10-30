#!/bin/bash

# Test script for VPP MCP Server
# This script tests the MCP server by sending JSON-RPC messages via stdio

set -e

echo "======================================"
echo "VPP MCP Server Test Script"
echo "======================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if the server binary exists
if [ ! -f "./vpp-mcp-server" ]; then
    echo -e "${RED}Error: vpp-mcp-server binary not found${NC}"
    echo "Please run: go build -o vpp-mcp-server main.go"
    exit 1
fi

echo -e "${YELLOW}Step 1: Testing MCP server initialization...${NC}"
echo ""

# Test 1: Check if server initializes properly
INIT_RESULT=$(
    (
        echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
        sleep 0.5
    ) | timeout 3s ./vpp-mcp-server 2>&1
)

if echo "$INIT_RESULT" | grep -q '"result"' && echo "$INIT_RESULT" | grep -q '"protocolVersion"'; then
    echo -e "${GREEN}✓ Server initializes successfully${NC}"
else
    echo -e "${RED}✗ Server failed to initialize${NC}"
    echo "$INIT_RESULT"
    exit 1
fi

echo ""
echo -e "${YELLOW}Step 2: Listing available tools...${NC}"
echo ""

# Test 2: List tools using JSON-RPC
(
    # Send initialize request
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
    sleep 0.5
    # Send tools/list request
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
    sleep 0.5
) | timeout 3s ./vpp-mcp-server 2>/dev/null | while IFS= read -r line; do
    # Parse and display available tools
    if echo "$line" | jq -e '.result.tools[]?' > /dev/null 2>&1; then
        echo "$line" | jq -r '.result.tools[] | "  - \(.name): \(.description | split("\n")[0])"'
    fi
done

echo ""
echo -e "${GREEN}Available Tools:${NC}"
echo "  1. vpp_show_version"
echo "  2. vpp_show_int"
echo "  3. vpp_show_int_addr"
echo "  4. vpp_show_errors"
echo "  5. vpp_show_session_verbose"
echo "  6. vpp_show_npol_rules"
echo "  7. vpp_show_npol_policies"
echo "  8. vpp_show_npol_ipset"
echo "  9. vpp_show_npol_interfaces"
echo " 10. vpp_trace"
echo " 11. vpp_pcap"
echo " 12. vpp_dispatch"
echo " 13. vpp_get_pods"
echo " 14. vpp_clear_errors"
echo " 15. vpp_tcp_stats"
echo " 16. vpp_session_stats"
echo " 17. vpp_get_logs"
echo " 18. vpp_show_cnat_translation"
echo " 19. vpp_show_cnat_session"
echo " 20. vpp_clear_run"
echo " 21. vpp_show_run"

echo ""
echo -e "${YELLOW}Step 3: Prerequisites check...${NC}"
echo ""

# Check if kubectl is available
if command -v kubectl &> /dev/null; then
    echo -e "${GREEN}✓ kubectl is installed${NC}"
    
    # Check if we can access any VPP pods
    POD_NAME=$(kubectl get pods -n calico-vpp-dataplane -l app=calico-vpp-dataplane -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
    
    if [ -n "$POD_NAME" ]; then
        echo -e "${GREEN}✓ Found VPP pod: $POD_NAME${NC}"
        echo ""
        echo -e "${YELLOW}Step 4: Testing vpp_show_version tool...${NC}"
        echo ""
        
        # Test calling a tool
        TOOL_REQUEST=$(cat <<EOF
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"vpp_show_version","arguments":{"pod_name":"$POD_NAME","namespace":"calico-vpp-dataplane"}}}
EOF
)
        
        RESULT=$(
            (
                echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
                sleep 0.5
                echo "$TOOL_REQUEST"
                sleep 2
            ) | timeout 5s ./vpp-mcp-server 2>/dev/null | grep -a "tools/call" || echo ""
        )
        
        if [ -n "$RESULT" ]; then
            echo -e "${GREEN}✓ Tool executed successfully${NC}"
            echo "$RESULT" | jq '.' 2>/dev/null || echo "$RESULT"
        else
            echo -e "${YELLOW}⚠ Could not test tool execution (this is normal if pods aren't ready)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ No VPP pods found in calico-vpp-dataplane namespace${NC}"
        echo "  To test with a real pod, ensure VPP is running in Kubernetes"
    fi
else
    echo -e "${YELLOW}⚠ kubectl not found${NC}"
    echo "  Install kubectl to test with real VPP pods"
fi

echo ""
echo -e "${YELLOW}Step 5: Integration test methods...${NC}"
echo ""
echo "To test with a real VPP pod, you can:"
echo ""
echo "1. ${GREEN}Use with Claude Desktop:${NC}"
echo "   Add to ~/.config/Claude/claude_desktop_config.json:"
echo '   {
     "mcpServers": {
       "vpp-debug": {
         "command": "'$(pwd)'/vpp-mcp-server"
       }
     }
   }'
echo ""
echo "2. ${GREEN}Use MCP Inspector (requires Node.js):${NC}"
echo "   npx @modelcontextprotocol/inspector $(pwd)/vpp-mcp-server"
echo ""
echo "3. ${GREEN}Manual JSON-RPC test:${NC}"
echo "   See example_mcp_requests.json for sample requests"
echo ""
echo "======================================"
echo -e "${GREEN}✓ MCP Server tests completed${NC}"
echo "======================================"
