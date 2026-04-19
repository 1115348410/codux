.PHONY: build test clean run fmt vet lint

# Variables
BINARY_NAME=codux
VERSION?=0.1.0
BUILD_DIR=dist

# Default target
all: build

# Build
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -v .

# Build with debug info
build-debug:
	go build -tags debug -o $(BUILD_DIR)/$(BINARY_NAME) -v .

# Run
run:
	go run .

# Test
test:
	go test -v ./...

# Test with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	go clean

# Format
fmt:
	go fmt ./...

# Vet
vet:
	go vet ./...

# Install dependencies
deps:
	go mod download

# Update dependencies
deps-update:
	go get -u ./...
	go mod tidy

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

# Build for macOS
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Build all platforms
build-all: build-linux build-macos build-windows

# Install
install:
	go install .

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-debug    - Build with debug info"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  deps           - Download dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  install        - Install the application"
	@echo "  build-all      - Build for all platforms"

