# Terraform Provider for Defguard

This Terraform provider allows you to manage Defguard resources including users, groups, networks, devices, and SNAT bindings.

## Requirements

- Terraform >= 1.0
- Go >= 1.21

## Building the Provider

1. Clone the repository
2. Run:

```bash
go build -o terraform-provider-defguard
```

## Using the Provider

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    defguard = {
      source = "registry.terraform.io/dburianov/defguard"
    }
  }
}

provider "defguard" {
  url      = "http://localhost:8000"
  api_key  = var.defguard_api_key
  session  = var.defguard_session
}
```

## Resources

- `defguard_user` - Manage users
- `defguard_group` - Manage groups
- `defguard_network` - Manage WireGuard networks
- `defguard_device` - Manage devices
- `defguard_snat_binding` - Manage SNAT bindings

## Development

Run tests:

```bash
make test
```

## License

Apache 2.0
