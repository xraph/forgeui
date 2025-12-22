# Makefile for ForgeUI
# Provides convenient commands for development and testing

.PHONY: help test lint coverage build clean install release-test release-snapshot

# Default target
help:
	@echo "ForgeUI Development Commands"
	@echo ""
	@echo "Testing & Quality:"
	@echo "  make test              Run all tests with race detector"
	@echo "  make test-coverage     Run tests with coverage report"
	@echo "  make lint              Run golangci-lint"
	@echo "  make coverage          Check coverage threshold (80%)"
	@echo ""
	@echo "Building:"
	@echo "  make build             Build CLI binary"
	@echo "  make install           Install CLI to GOPATH/bin"
	@echo "  make clean             Clean build artifacts"
	@echo ""
	@echo "Release:"
	@echo "  make release-test      Test GoReleaser configuration"
	@echo "  make release-snapshot  Create snapshot release (local testing)"
	@echo ""

# Run tests with race detector
test:
	@echo "Running tests with race detector..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m

# Check coverage threshold
coverage: test
	@echo "Checking coverage threshold..."
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	threshold=80.0; \
	echo "Total coverage: $${total}%"; \
	if [ $$(echo "$${total} < $${threshold}" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $${total}% is below threshold $${threshold}%"; \
		exit 1; \
	else \
		echo "✅ Coverage $${total}% meets threshold $${threshold}%"; \
	fi

# Build CLI binary
build:
	@echo "Building ForgeUI CLI..."
	go build -v -o forgeui ./cmd/forgeui
	@echo "Binary created: ./forgeui"
	@./forgeui --version

# Install CLI to GOPATH/bin
install:
	@echo "Installing ForgeUI CLI..."
	go install -v ./cmd/forgeui
	@echo "Installed to: $$(which forgeui)"
	@forgeui --version

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f forgeui
	rm -f coverage.out coverage.html coverage.txt
	rm -rf dist/
	rm -f *_coverage.out
	@echo "Clean complete"

# Test GoReleaser configuration
release-test:
	@echo "Testing GoReleaser configuration..."
	goreleaser check
	@echo "✅ GoReleaser configuration is valid"

# Create snapshot release for local testing
release-snapshot:
	@echo "Creating snapshot release..."
	@echo "This will build binaries for all platforms without publishing"
	goreleaser release --snapshot --clean
	@echo ""
	@echo "✅ Snapshot release created in dist/"
	@echo ""
	@echo "Test the binaries:"
	@ls -lh dist/forgeui_*/forgeui* | head -10

# Run all CI checks locally
ci: lint test coverage build
	@echo ""
	@echo "✅ All CI checks passed!"
	@echo "Ready to push or create a release"

# Quick check before committing
check: lint test
	@echo ""
	@echo "✅ Quick checks passed!"

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest
	@echo "✅ Dependencies installed"

# Run integration tests (if tagged)
test-integration:
	@echo "Running integration tests..."
	@if go test -v -tags=integration ./... 2>&1 | grep -q "no packages to test"; then \
		echo "No integration tests found"; \
	else \
		go test -v -tags=integration ./...; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .
	@echo "✅ Code formatted"

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	go mod verify
	go mod tidy
	@echo "✅ Dependencies verified"

# Run all checks (comprehensive)
all: verify fmt lint test-integration test coverage build
	@echo ""
	@echo "✅ All checks passed!"
	@echo "Project is ready for release"

