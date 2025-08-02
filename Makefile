# GoViz Makefile
# Professional build system for GoViz CLI tool

# Metadata
BINARY_NAME := goviz
VERSION ?= v0.1.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go settings
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -w -s -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)
BUILD_FLAGS := -ldflags="$(LDFLAGS)"

# Directories
DIST_DIR := dist
SCRIPTS_DIR := scripts

# Default target
.PHONY: all
all: clean test build

# Help target
.PHONY: help
help: ## Show this help message
	@echo "GoViz Build System"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
.PHONY: build
build: ## Build the binary
	@echo "🔨 Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) .
	@echo "✅ Build completed: ./$(BINARY_NAME)"

.PHONY: build-debug
build-debug: ## Build with debug information
	@echo "🔨 Building $(BINARY_NAME) with debug info..."
	$(GOBUILD) -gcflags="all=-N -l" -o $(BINARY_NAME) .
	@echo "✅ Debug build completed: ./$(BINARY_NAME)"

.PHONY: install
install: build ## Install binary to system
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installation completed"

.PHONY: uninstall
uninstall: ## Remove binary from system
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstallation completed"

# Testing targets
.PHONY: test
test: ## Run tests
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	$(GOTEST) -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "🧪 Running tests with race detection..."
	$(GOTEST) -race ./...

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Code quality targets
.PHONY: lint
lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "💄 Formatting code..."
	$(GOCMD) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "🔍 Running go vet..."
	$(GOCMD) vet ./...

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "🧹 Tidying go modules..."
	$(GOMOD) tidy

.PHONY: verify
verify: fmt vet lint test ## Run all verification steps
	@echo "✅ All verification steps completed"

# Release targets
.PHONY: build-releases
build-releases: ## Build releases for all platforms
	@echo "🏗️  Building releases for all platforms..."
	@chmod +x $(SCRIPTS_DIR)/build-releases.sh
	VERSION=$(VERSION) $(SCRIPTS_DIR)/build-releases.sh

.PHONY: release
release: verify build-releases ## Create a full release (verify + build all platforms)
	@echo "🚀 Release $(VERSION) ready in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

# Utility targets
.PHONY: run
run: build ## Build and run the application
	@echo "🚀 Running $(BINARY_NAME)..."
	./$(BINARY_NAME) --help

.PHONY: demo
demo: build ## Run demo commands
	@echo "🎮 Running demo..."
	./$(BINARY_NAME) generate --format tree
	@echo ""
	./$(BINARY_NAME) analyze
	@echo ""
	./$(BINARY_NAME) doctor

.PHONY: clean
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

.PHONY: deps
deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: upgrade-deps
upgrade-deps: ## Upgrade all dependencies
	@echo "⬆️  Upgrading dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Development helpers
.PHONY: watch
watch: ## Watch for changes and rebuild (requires entr)
	@echo "👀 Watching for changes..."
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -r make build; \
	else \
		echo "⚠️  entr not found. Install with your package manager"; \
	fi

.PHONY: serve-docs
serve-docs: ## Serve documentation locally (requires Python)
	@echo "📚 Serving documentation on http://localhost:8000"
	@if command -v python3 >/dev/null 2>&1; then \
		python3 -m http.server 8000; \
	elif command -v python >/dev/null 2>&1; then \
		python -m SimpleHTTPServer 8000; \
	else \
		echo "⚠️  Python not found"; \
	fi

# Security targets
.PHONY: vuln-check
vuln-check: ## Check for vulnerabilities
	@echo "🔒 Checking for vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "⚠️  govulncheck not found. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

.PHONY: audit
audit: vuln-check ## Run security audit
	@echo "🔍 Running security audit..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Information targets
.PHONY: info
info: ## Show build information
	@echo "📋 Build Information:"
	@echo "  Binary Name: $(BINARY_NAME)"
	@echo "  Version:     $(VERSION)"
	@echo "  Commit:      $(COMMIT)"
	@echo "  Build Time:  $(BUILD_TIME)"
	@echo "  Go Version:  $(shell $(GOCMD) version)"

.PHONY: size
size: build ## Show binary size
	@echo "📏 Binary size:"
	@ls -lh $(BINARY_NAME) | awk '{print "  Size: " $$5}'
	@file $(BINARY_NAME)

# Setup targets for new developers
.PHONY: setup
setup: deps ## Setup development environment
	@echo "🛠️  Setting up development environment..."
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "✅ Development environment ready"

.PHONY: pre-commit
pre-commit: verify vuln-check ## Run pre-commit checks
	@echo "✅ Pre-commit checks passed"

# CI/CD targets
.PHONY: ci
ci: verify build-releases ## CI pipeline
	@echo "🤖 CI pipeline completed"

.PHONY: cd
cd: release ## CD pipeline
	@echo "🚀 CD pipeline completed"