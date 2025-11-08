VERSION := 0.1.0
PROVIDER := keycloak
PROJECT := github.com/raushan606/pulumi-keycloak

# Build configuration
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR := bin
PROVIDER_BINARY := pulumi-resource-$(PROVIDER)

# Schema configuration
SCHEMA_FILE := schema.json
PULUMI_CONVERT := true

.PHONY: build install clean test lint schema generate all

# Default target
all: build

# Build the provider binary
build:
	@echo "Building $(PROVIDER_BINARY) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	cd provider && go build -o ../$(BUILD_DIR)/$(PROVIDER_BINARY) ./cmd/pulumi-resource-$(PROVIDER)
	@echo "✅ Built $(BUILD_DIR)/$(PROVIDER_BINARY)"

# Install the provider binary to PATH
install: build
	@echo "Installing $(PROVIDER_BINARY) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(PROVIDER_BINARY) /usr/local/bin/
	@echo "✅ Installed $(PROVIDER_BINARY)"

# Install to local path (no sudo required)
install-local: build
	@echo "Installing $(PROVIDER_BINARY) to ~/bin..."
	@mkdir -p ~/bin
	cp $(BUILD_DIR)/$(PROVIDER_BINARY) ~/bin/
	@echo "✅ Installed $(PROVIDER_BINARY) to ~/bin"
	@echo "⚠️  Make sure ~/bin is in your PATH"

# Generate SDKs
generate-sdk: build
	@echo "Generating SDKs..."
	pulumi package gen-sdk --language typescript ./$(BUILD_DIR)/$(PROVIDER_BINARY)
	pulumi package gen-sdk --language java ./$(BUILD_DIR)/$(PROVIDER_BINARY)
	@echo "✅ Generated TypeScript and Java SDKs"
	@echo "Install dependencies for generated SDKs"
	cd sdk/nodejs && npm install
	@echo "✅ Installed dependencies for TypeScript SDK"

# Generate TypeScript SDK only
generate-sdk-typescript: build
	@echo "Generating TypeScript SDK..."
	pulumi package gen-sdk --language typescript ./$(BUILD_DIR)/$(PROVIDER_BINARY)
	@echo "✅ Generated TypeScript SDK in sdk/nodejs"

# Generate Java SDK only
generate-sdk-java: build
	@echo "Generating Java SDK..."
	pulumi package gen-sdk --language java ./$(BUILD_DIR)/$(PROVIDER_BINARY)
	@echo "✅ Generated Java SDK in sdk/java"

# Do all above steps in one
all: build
	@echo "Make build.."
	@$(MAKE) build
	@echo "Make install.."
	@$(MAKE) install
	@echo "Make generate-sdk.."
	@$(MAKE) generate-sdk
	@echo "All done."


# Clean SDKs
clean-sdk:
	@echo "Cleaning SDK artifacts..."
	rm -rf sdk/
	@echo "✅ Cleaned SDK artifacts"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	cd provider && go mod tidy
	@echo "✅ Dependencies updated"

# Run tests
test:
	@echo "Running tests..."
	cd provider && go test ./...
	@echo "✅ Tests passed"

# Run linter
lint:
	@echo "Running linter..."
	cd provider && golangci-lint run
	@echo "✅ Linting passed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(SCHEMA_FILE)
	rm -rf sdk/
	@echo "✅ Cleaned"

# Development build (with debug info)
dev-build:
	@echo "Building $(PROVIDER_BINARY) for development..."
	@mkdir -p $(BUILD_DIR)
	cd provider && go build -gcflags="all=-N -l" -o ../$(BUILD_DIR)/$(PROVIDER_BINARY) ./cmd/pulumi-resource-$(PROVIDER)
	@echo "✅ Built $(BUILD_DIR)/$(PROVIDER_BINARY) with debug info"

# Cross-platform builds
build-linux:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-darwin:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-windows:
	@$(MAKE) build GOOS=windows GOARCH=amd64

build-all: build-linux build-darwin build-windows

# Release build (all platforms)
release: clean build-all schema
	@echo "✅ Release build complete"

# Help target
help:
	@echo "Available targets:"
	@echo "  build              - Build the provider binary"
	@echo "  install            - Install the provider binary to /usr/local/bin"
	@echo "  install-local      - Install the provider binary to ~/bin"
	@echo "  schema             - Generate the provider schema"
	@echo "  generate-sdk       - Generate TypeScript and Java SDKs"
	@echo "  generate-sdk-typescript - Generate TypeScript SDK only"
	@echo "  generate-sdk-java  - Generate Java SDK only"
	@echo "  generate-sdk-python - Generate Python SDK only"
	@echo "  generate-sdk-go    - Generate Go SDK only"
	@echo "  generate-sdk-dotnet - Generate .NET SDK only"
	@echo "  generate-all-sdks  - Generate all SDKs (TS, Java, Python, Go, .NET)"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  test               - Run tests"
	@echo "  lint               - Run linter"
	@echo "  clean              - Clean build artifacts and SDKs"
	@echo "  clean-sdk          - Clean SDK artifacts only"
	@echo "  dev-build          - Build with debug information"
	@echo "  build-all          - Build for all platforms"
	@echo "  release            - Create a release build"
	@echo "  help               - Show this help message"

# Check if required tools are installed
check-tools:
	@echo "Checking required tools..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed"; exit 1; }
	@echo "✅ Go is installed"
	@command -v golangci-lint >/dev/null 2>&1 || { echo "⚠️  golangci-lint is not installed (optional for development)"; }
	@echo "✅ Tool check complete"