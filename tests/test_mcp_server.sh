#!/bin/bash

# Test script for VPP MCP Server
# This script tests the MCP server by sending JSON-RPC messages via stdio

set -e

echo "==================================="
echo "VPP MCP Server Test stdio Transport"
echo "==================================="
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
echo -e "${YELLOW}Step 3: Integration test methods...${NC}"
echo ""
echo "To test with a real VPP pod, you can:"
echo ""
echo "1. Use with VS Code (Copilot Agent) / Windsurf (Cascade Agent)"
echo ""
echo "   Add to <project-root>/.codeium/windsurf/mcp_config.json (Windsurf)"
echo '
{
  "mcpServers": {
    "vpp-debug": {
      "command": "<path to vpp-mcp-server>",
      "disabledTools": [],
      "disabled": false
    }
  }
}'
echo ""
echo "   OR if the server is hosted remotely over HTTP"
echo ""
echo "   Add to <project-root>/.vscode/mcp.json (VSCode)"
echo '
{
  "servers": {
    "vpp-debug-remote": {
      "type": "sse",
      "url": "http://<remote-IP>:<remote-port>/sse"
    }
  }
}'
echo ""
echo "2. Use MCP Inspector (requires Node.js):"
echo "   npx @modelcontextprotocol/inspector <path to vpp-mcp-server>"
echo ""
echo "3. Manual JSON-RPC test:"
echo "   See example_mcp_requests.json for sample requests"
echo ""
echo "======================================"
echo -e "${GREEN}✓ MCP Server tests completed${NC}"
echo "======================================"
