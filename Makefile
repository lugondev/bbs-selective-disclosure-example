.PHONY: help build test test-integration run-demo clean fmt vet

# Default target
help:
	@echo "Available targets:"
	@echo "  build            - Build the demo application"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  run-demo         - Run the demo application"
	@echo "  clean            - Clean build artifacts"
	@echo "  fmt              - Format Go code"
	@echo "  vet              - Run go vet"
	@echo "  lint             - Run golangci-lint (if available)"

# Build the demo application
build:
	@echo "Building demo application..."
	go build -o bin/demo ./cmd/demo

# Run all tests
test: fmt vet test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v ./pkg/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v ./test/integration/...

# Run the demo application
run-demo: build
	@echo "Running BBS+ Selective Disclosure Demo..."
	./bin/demo

# Run demo without building (go run)
demo:
	@echo "Running BBS+ Selective Disclosure Demo..."
	go run ./cmd/demo

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint (if available)
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

# Initialize project dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Generate test coverage report
test-coverage:
	@echo "Generating test coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Development setup complete!"

# Verify project structure
verify:
	@echo "Verifying project structure..."
	@test -d pkg/bbs && echo "✓ BBS+ package exists" || echo "✗ BBS+ package missing"
	@test -d pkg/did && echo "✓ DID package exists" || echo "✗ DID package missing"
	@test -d pkg/vc && echo "✓ VC package exists" || echo "✗ VC package missing"
	@test -d internal/issuer && echo "✓ Issuer use case exists" || echo "✗ Issuer use case missing"
	@test -d internal/holder && echo "✓ Holder use case exists" || echo "✗ Holder use case missing"
	@test -d internal/verifier && echo "✓ Verifier use case exists" || echo "✗ Verifier use case missing"
	@test -f cmd/demo/main.go && echo "✓ Demo application exists" || echo "✗ Demo application missing"
	@test -f test/integration/full_lifecycle_test.go && echo "✓ Integration tests exist" || echo "✗ Integration tests missing"

# Create bin directory
bin:
	mkdir -p bin
