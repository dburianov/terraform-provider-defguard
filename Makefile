# Default target
default: build

# Build the provider
build:
	go build -v ./...

# Install the provider locally
install: build
	go install -v ./...

# Run unit tests
test:
	go test -v -cover -timeout=120s -parallel=4 ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

# Format Go code
fmt:
	go fmt ./...

# Run Go vet
vet:
	go vet ./...

# Clean dependencies
clean:
	go mod tidy

# Generate documentation
docs:
	go generate ./...

# Run all checks
check: fmt vet test

# Development build with debug info
dev-build:
	go build -gcflags="-N -l" -v ./...

# Install for local development
dev-install: dev-build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/dburianov/defguard/0.0.1/linux_amd64
	cp terraform-provider-defguard ~/.terraform.d/plugins/registry.terraform.io/dburianov/defguard/0.0.1/linux_amd64/

# Deploy to dev (prod folder)
deploy-dev: build
	./scripts/deploy-dev.sh

# Deploy to prod environment
deploy-prod:
	@echo "Deploying to prod..."
	TF_LOG=INFO terraform -chdir=prod apply -auto-approve

# Validate terraform configuration in prod
validate:
	terraform -chdir=prod validate

# Run E2E tests from prod folder
test-e2e: build
	./scripts/test.sh

# Dev setup with environment variables
dev-setup:
	@echo "Set up the following environment variables for development:"
	@echo "  export DEFGUARD_ENDPOINT=https://dev.vpn.ddsc.ai"
	@echo "  export DEFGUARD_COOKIE=your-cookie-here"

.PHONY: build install test testacc fmt vet clean docs check dev-build dev-install deploy-dev deploy-prod validate test-e2e dev-setup