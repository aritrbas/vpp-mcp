# Makefile for VPP MCP Server

# Variables
BINARY_NAME=vpp-mcp-server
MAIN_FILE=main.go
BUILD_DIR=build

# Build the server
.PHONY: build
build:
	@echo "Building VPP MCP Server..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BINARY_NAME)"

# Build for current platform (detects OS automatically)
.PHONY: build-all
build-all:
	@echo "Detecting platform..."
	@OS=$$(uname -s); \
	if [ "$$OS" = "Linux" ]; then \
		$(MAKE) build-linux; \
	elif [ "$$OS" = "Darwin" ]; then \
		$(MAKE) build-darwin; \
	elif echo "$$OS" | grep -q "MINGW\|MSYS\|CYGWIN"; then \
		$(MAKE) build-windows; \
	else \
		echo "Unknown platform: $$OS"; \
		exit 1; \
	fi

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_FILE)

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos $(MAIN_FILE)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_FILE)

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod tidy
	go mod download

# Run the server
.PHONY: run
run: build
	@echo "Starting VPP MCP Server..."
	./$(BINARY_NAME)

# Test the setup
.PHONY: test
test:
	@echo "Running setup tests..."
	./test_setup.sh

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the VPP MCP server binary"
	@echo "  build-all    - Build for current platform (auto-detects OS)"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-darwin - Build for macOS"
	@echo "  build-windows- Build for Windows"
	@echo "  deps         - Download Go dependencies"
	@echo "  run          - Build and run the server"
	@echo "  test         - Run setup tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"
