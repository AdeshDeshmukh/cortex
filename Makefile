.PHONY: help build test clean install run lint fmt coverage

BINARY_NAME=cortex
VERSION?=dev
BUILD_DIR=dist
GO_FILES=$(shell find . -name '*.go' -type f)

help:
	@echo "🧠 Cortex - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make build       Build binary"
	@echo "  make test        Run tests"
	@echo "  make run         Run without building"
	@echo "  make fmt         Format code"
	@echo "  make lint        Run linters"
	@echo ""
	@echo "Deployment:"
	@echo "  make install     Install to system"
	@echo "  make release     Build all platforms"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean       Remove build artifacts"
	@echo ""

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY_NAME) cmd/cortex/main.go
	@echo "✅ Binary created: ./$(BINARY_NAME)"

test:
	@echo "🧪 Running tests..."
	@go test -v -race -coverprofile=coverage.txt ./...
	@echo "✅ Tests complete"

coverage: test
	@echo "📊 Generating coverage report..."
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "✅ Open coverage.html in browser"

run:
	@go run cmd/cortex/main.go

install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"
	@echo "   Run 'cortex --help' from anywhere"

clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.txt coverage.html
	@rm -rf $(BUILD_DIR)
	@go clean -testcache
	@echo "✅ Clean complete"

fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint:
	@echo "🔍 Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed"; \
		echo "   Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

release:
	@echo "📦 Building releases..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 cmd/cortex/main.go
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 cmd/cortex/main.go
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 cmd/cortex/main.go
	@GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe cmd/cortex/main.go
	@echo "✅ Releases built in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

.DEFAULT_GOAL := help