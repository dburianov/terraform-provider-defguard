#!/bin/bash
set -e

echo "=========================================="
echo "Running CI checks locally (same as GitHub Actions)"
echo "=========================================="

echo ""
echo "Step 1: Download dependencies"
go mod download

echo ""
echo "Step 2: Run tests"
go test -v ./...

echo ""
echo "Step 3: Run vet"
go vet ./...

echo ""
echo "Step 4: Run golangci-lint"
export PATH="$PATH:$(go env GOPATH)/bin"
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo "WARNING: golangci-lint not installed. Install with:"
    echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
fi

echo ""
echo "Step 5: Build"
go build -v ./...

echo ""
echo "Step 6: Run acceptance tests (optional)"
read -p "Run acceptance tests? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    TF_ACC=1 go test -v ./internal/provider
fi

echo ""
echo "=========================================="
echo "✓ All CI checks completed successfully!"
echo "=========================================="
