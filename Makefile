# Makefile for ForgeUI
# Provides convenient commands for development and testing

.PHONY: help test lint coverage build clean install release-test release-snapshot
.PHONY: fmt lint-fix f l t b c check ci all verify deps test-integration watch
.PHONY: test-verbose test-watch pre-commit update-deps info run-example bench
.PHONY: profile-cpu profile-mem check-deps audit vet build-release

# Default target
help:
	@echo "ForgeUI Development Commands"
	@echo ""
	@echo "Testing & Quality:"
	@echo "  make test (t)          Run all tests with race detector"
	@echo "  make test-coverage     Run tests with coverage report"
	@echo "  make test-integration  Run integration tests"
	@echo "  make test-verbose      Run tests with verbose output"
	@echo "  make test-watch        Watch and run tests on file changes"
	@echo "  make lint (l)          Run golangci-lint"
	@echo "  make lint-fix          Auto-fix linting issues"
	@echo "  make coverage          Check coverage threshold (80%)"
	@echo "  make fmt (f)           Format code with gofmt and goimports"
	@echo ""
	@echo "Building:"
	@echo "  make build (b)         Build CLI binary"
	@echo "  make install           Install CLI to GOPATH/bin"
	@echo "  make clean (c)         Clean build artifacts"
	@echo ""
	@echo "Development:"
	@echo "  make check             Quick pre-commit checks (lint + test)"
	@echo "  make ci                Run all CI checks locally"
	@echo "  make pre-commit        Comprehensive pre-commit checks"
	@echo "  make verify            Verify and tidy dependencies"
	@echo "  make deps              Install development dependencies"
	@echo "  make watch             Watch files and run tests"
	@echo ""
	@echo "Release:"
	@echo "  make release-test      Test GoReleaser configuration"
	@echo "  make release-snapshot  Create snapshot release (local testing)"
	@echo ""

# Run tests with race detector
test:
	@echo "Running tests with race detector..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Short alias for test
t: test

# Run tests with verbose output (no coverage)
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -race ./...

# Watch and run tests on file changes
test-watch:
	@echo "Watching for file changes..."
	@echo "Press Ctrl+C to stop"
	@while true; do \
		$(MAKE) test || true; \
		echo "\nWaiting for changes..."; \
		fswatch -1 -r --exclude=".*\.out$$" --exclude=".*\.html$$" --exclude="dist/" . || inotifywait -r -e modify,create,delete .; \
	done

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m

# Short alias for lint
l: lint

# Auto-fix linting issues
lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	golangci-lint run --fix --timeout=5m
	@echo "✅ Linting issues fixed where possible"

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

# Short alias for build
b: build

# Build with optimization flags
build-release:
	@echo "Building ForgeUI CLI (optimized)..."
	go build -v -ldflags="-s -w" -o forgeui ./cmd/forgeui
	@echo "Binary created: ./forgeui (stripped)"
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

# Short alias for clean
c: clean

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

# Quick check before committing (fast)
check: lint test
	@echo ""
	@echo "✅ Quick checks passed!"

# Comprehensive pre-commit checks
pre-commit: verify fmt lint test
	@echo ""
	@echo "✅ Pre-commit checks passed!"
	@echo "Ready to commit"

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing goreleaser..."
	@go install github.com/goreleaser/goreleaser@latest
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ All dependencies installed"

# Watch files and run commands on changes
watch:
	@echo "Watching files for changes..."
	@echo "Will run tests on Go file changes"
	@echo "Press Ctrl+C to stop"
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o --exclude=".*\.out$$" --exclude=".*\.html$$" --exclude="dist/" -e ".*" -i "\\.go$$" . | xargs -n1 -I{} sh -c 'clear && date && $(MAKE) test-verbose'; \
	elif command -v inotifywait >/dev/null 2>&1; then \
		while true; do \
			inotifywait -r -e modify,create,delete --include '.*\.go$$' .; \
			clear && date && $(MAKE) test-verbose; \
		done; \
	else \
		echo "❌ Neither fswatch nor inotifywait found"; \
		echo "Install fswatch: brew install fswatch (macOS) or apt-get install inotify-tools (Linux)"; \
		exit 1; \
	fi

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
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not found, using gofmt only"; \
		gofmt -s -w .; \
	fi
	@echo "✅ Code formatted"

# Short alias for fmt
f: fmt

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	go mod verify
	go mod tidy
	@echo "✅ Dependencies verified"

# Verify and update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated"

# Run all checks (comprehensive)
all: verify fmt lint test-integration test coverage build
	@echo ""
	@echo "✅ All checks passed!"
	@echo "Project is ready for release"

# Show project info
info:
	@echo "ForgeUI Project Information"
	@echo "============================"
	@echo "Go version: $$(go version)"
	@echo "Module: $$(go list -m)"
	@echo "Dependencies: $$(go list -m all | wc -l) modules"
	@echo "Packages: $$(go list ./... | wc -l) packages"
	@echo "Go files: $$(find . -name '*.go' -not -path './vendor/*' | wc -l) files"
	@echo "Test files: $$(find . -name '*_test.go' -not -path './vendor/*' | wc -l) files"
	@echo ""
	@if [ -f forgeui ]; then \
		echo "Binary size: $$(du -h forgeui | cut -f1)"; \
	fi

# Run example application
run-example:
	@echo "Running example application..."
	@cd example && go run .

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Profile CPU usage
profile-cpu:
	@echo "Running CPU profiling..."
	go test -cpuprofile=cpu.prof -bench=. ./...
	@echo "View profile with: go tool pprof cpu.prof"

# Profile memory usage
profile-mem:
	@echo "Running memory profiling..."
	go test -memprofile=mem.prof -bench=. ./...
	@echo "View profile with: go tool pprof mem.prof"

# Check for outdated dependencies
check-deps:
	@echo "Checking for outdated dependencies..."
	@go list -u -m all 2>/dev/null | grep '\[' || echo "All dependencies are up to date"

# Security audit
audit:
	@echo "Running security audit..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

# Check for common issues
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✅ No issues found"

