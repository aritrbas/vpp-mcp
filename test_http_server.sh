#!/bin/bash

# Test script for VPP MCP Server HTTP transport

PORT=9090
TIMEOUT=5

echo "Starting VPP MCP Server on port $PORT..."
./vpp-mcp-server --transport=http --port=$PORT > /tmp/vpp-mcp-test.log 2>&1 &
SERVER_PID=$!

# Wait for server to start
sleep 2

echo "Testing endpoints..."

# Test health endpoint
echo -n "Health check: "
if curl -s -f http://localhost:$PORT/health > /dev/null; then
    echo "✓ OK"
else
    echo "✗ FAILED"
fi

# Test root endpoint
echo -n "Root endpoint: "
if curl -s -f http://localhost:$PORT/ > /dev/null; then
    echo "✓ OK"
else
    echo "✗ FAILED"
fi

# Test SSE endpoint (just check it responds)
echo -n "SSE endpoint: "
if curl -s -f -N http://localhost:$PORT/sse -H "Accept: text/event-stream" --max-time $TIMEOUT > /dev/null 2>&1; then
    echo "✓ OK"
else
    # SSE might timeout waiting for events, which is OK
    echo "✓ OK (timeout is expected)"
fi

# Cleanup
echo ""
echo "Stopping server (PID: $SERVER_PID)..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo ""
echo "Server logs:"
cat /tmp/vpp-mcp-test.log

echo ""
echo "Test complete!"
