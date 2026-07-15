# DefGuard Provider

The DefGuard provider allows you to manage infrastructure and users in DefGuard VPN platform using Terraform.

DefGuard is an open-source VPN solution that provides secure remote access with WireGuard, multi-factor authentication, and centralized user management.

## Example Usage

### Default Configuration

```hcl
terraform {
  required_providers {
    defguard = {
      source  = "dburianov/defguard"
      version = "~> 1.0.0"
    }
  }
}

provider "defguard" {}
```

### With Custom Endpoint

```hcl
provider "defguard" {
  endpoint = "https://vpn.example.com"
}
```

## Schema

### Optional

- `endpoint` (String) - DefGuard instance URL. If not set, uses environment variable `DEFGUARD_ENDPOINT`.
- `cookie` (String, Sensitive) - Authentication cookie for API access. If not set, uses environment variable `DEFGUARD_COOKIE`.

## Resources

- [defguard_user](resources/user.md) - Manages users in DefGuard
- [defguard_group](resources/group.md) - Manages user groups in DefGuard
- [defguard_network](resources/network.md) - Manages WireGuard networks in DefGuard
- [defguard_device](resources/device.md) - Manages WireGuard devices in DefGuard
- [defguard_snat_binding](resources/snat_binding.md) - Manages SNAT bindings in DefGuard