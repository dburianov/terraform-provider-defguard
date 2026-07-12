# Terraform Provider DefGuard - Automation Documentation

## Overview

This document describes the automation setup for developing, deploying, and testing the Terraform provider for DefGuard.

## Quick Start

### Prerequisites

- Go 1.23+
- Terraform 1.15.7+
- Access to a DefGuard instance

### Development

Build and install the provider locally:

```bash
# Using make
make dev-install

# Or using script
./scripts/dev.sh
```

### Deploy to Dev Environment

Deploy using the `prod/` folder configuration:

```bash
# Using make (requires DEFGUARD_ENDPOINT and DEFGUARD_COOKIE env vars)
DEFGUARD_ENDPOINT="https://dev.vpn.ddsc.ai" \
DEFGUARD_COOKIE="your-cookie-here" \
make deploy-dev

# Or using script directly (uses defaults from provider.tf)
./scripts/deploy-dev.sh
```

### Testing

Run all tests:

```bash
# Unit tests only
make test

# Acceptance tests
make testacc

# Full E2E suite
make test-e2e

# Or using script
./scripts/test.sh
```

### Validation

Validate Terraform configuration:

```bash
make validate
# or
terraform -chdir=prod validate
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `build` | Build the provider |
| `install` | Install to Go bin directory |
| `test` | Run unit tests |
| `testacc` | Run acceptance tests (TF_ACC=1) |
| `fmt` | Format Go code |
| `vet` | Run go vet |
| `dev-build` | Build with debug symbols |
| `dev-install` | Install for local development |
| `deploy-dev` | Deploy to dev environment |
| `deploy-prod` | Deploy to prod environment |
| `validate` | Validate terraform config in prod/ |
| `test-e2e` | Run full E2E test suite |
| `dev-setup` | Show environment variable setup |

## Scripts

| Script | Description |
|--------|-------------|
| `scripts/dev.sh` | Build and install provider with debug symbols |
| `scripts/deploy-dev.sh` | Init -> Plan -> Apply from prod/ folder |
| `scripts/test.sh` | Run unit, acceptance, and E2E tests |
| `scripts/validate.sh` | Validate terraform config in prod/ |

## Environment Variables

The provider configuration uses these variables:

- `DEFGUARD_ENDPOINT` - DefGuard instance URL (default: https://dev.vpn.ddsc.ai)
- `DEFGUARD_COOKIE` - Authentication cookie for API access

Set them before running terraform:

```bash
export DEFGUARD_ENDPOINT="https://dev.vpn.ddsc.ai"
export DEFGUARD_COOKIE="your-cookie-here"
```

## CI/CD Pipeline

The `.github/workflows/ci.yaml` defines the pipeline with these jobs:

1. **lint** - Run go fmt, vet, and build
2. **test** - Run unit and acceptance tests
3. **validate-terraform** - Validate terraform config
4. **deploy-dev** - Deploy to dev on main branch
5. **deploy-prod** - Deploy to prod on main branch

### Required GitHub Secrets

For deployment jobs, set these repository secrets:

- `DEFGUARD_COOKIE` - Authentication cookie for DefGuard API

## Directory Structure

```
.
├── scripts/            # Deployment and test automation scripts
│   ├── dev.sh         # Development build/install
│   ├── deploy-dev.sh  # Deploy to prod folder
│   ├── test.sh        # Run all tests
│   └── validate.sh    # Validate terraform config
├── prod/               # Production/test IAC configuration
│   ├── provider.tf    # Provider configuration with variables
│   ├── users.tf       # User resources
│   ├── group.tf       # Group resources
│   └── terraform.tfstate  # State file (local)
├── .github/workflows/
│   └── ci.yaml        # CI/CD pipeline definition
└── Makefile           # Build and deploy targets
```

## Workflow Examples

### Development Loop

```bash
# 1. Make changes to provider code
make dev-install

# 2. Test the changes
terraform -chdir=prod plan

# 3. Apply if satisfied
terraform -chdir=prod apply -auto-approve
```

### CI/CD Deployment

When pushing to `main` branch:
1. All tests run automatically
2. On success, deploy-dev job runs
3. deploy-prod job can be triggered manually for production deployment

## Provider Resources

The provider supports the following resources:

| Resource | Description |
|----------|-------------|
| `defguard_user` | Manages users in DefGuard |
| `defguard_group` | Manages user groups in DefGuard |
| `defguard_network` | Manages WireGuard networks in DefGuard |
| `defguard_device` | Manages WireGuard devices in DefGuard |
| `defguard_snat_binding` | Manages SNAT (Source NAT) bindings in DefGuard |

For detailed documentation on each resource, see [docs/resources/](docs/resources/).

## Notes

- State is stored locally in the `prod/` folder as `terraform.tfstate`
- The provider uses cookie-based authentication
- Always review terraform plan before applying in non-dev environments
