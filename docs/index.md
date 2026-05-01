# Defguard Provider

The Defguard provider is used to interact with the Defguard infrastructure management system. It allows you to manage users, groups, and devices through Terraform.

## Example Usage

```hcl
provider "defguard" {
  base_url  = "https://your-defguard-instance.com"
  api_token = "your-api-token"
}

resource "defguard_user" "example_user" {
  username   = "exampleuser"
  first_name = "Example"
  last_name  = "User"
  email      = "example@example.com"
  is_active  = true
}

resource "defguard_group" "example_group" {
  name      = "example-group"
  is_admin  = false
  members   = [defguard_user.example_user.username]
}

resource "defguard_device" "example_device" {
  name               = "example-device"
  wireguard_pubkey   = "examplepubkey123"
  username           = defguard_user.example_user.username
}
```

## Authentication

The Defguard provider requires an API token to authenticate with the Defguard API. You can provide this token in several ways:

1. Using the `api_token` argument in the provider block
2. Using the `DEFGUARD_API_TOKEN` environment variable
3. Using a credentials file (if supported by your Defguard instance)

## Argument Reference

The provider supports the following arguments:

- `base_url` - (Required) The base URL of the Defguard API
- `api_token` - (Required) The API token for authentication

## Resources

### defguard_user

Manages a Defguard user.

- `username` - (Required) The username of the user
- `first_name` - (Required) The first name of the user
- `last_name` - (Required) The last name of the user
- `email` - (Required) The email address of the user
- `password` - (Optional) The password for the user
- `phone` - (Optional) The phone number of the user
- `is_active` - (Optional) Whether the user is active (default: true)
- `is_admin` - (Optional) Whether the user is an administrator (default: false)
- `groups` - (Optional) The groups the user belongs to
- `mfa_enabled` - (Optional) Whether multi-factor authentication is enabled (default: false)
- `enrolled` - (Optional) Whether the user has completed enrollment (default: false)
- `email_mfa_enabled` - (Optional) Whether email-based MFA is enabled (default: false)
- `totp_enabled` - (Optional) Whether TOTP-based MFA is enabled (default: false)
- `ldap_pass_requires_change` - (Optional) Whether LDAP password requires change (default: false)
- `mfa_method` - (Optional) The MFA method used by the user (default: "None")

### defguard_group

Manages a Defguard group.

- `name` - (Required) The name of the group
- `is_admin` - (Optional) Whether the group has admin privileges (default: false)
- `members` - (Optional) The members of the group
- `vpn_locations` - (Optional) The VPN locations associated with the group

### defguard_device

Manages a Defguard device.

- `name` - (Required) The name of the device
- `wireguard_pubkey` - (Required) The WireGuard public key of the device
- `user_id` - (Optional) The ID of the user that owns the device
- `username` - (Optional) The username of the user that owns the device
- `description` - (Optional) Description of the device
- `configured` - (Optional) Whether the device is configured (default: true)
- `device_type` - (Optional) The type of device (user or network) (default: "user")

## Data Sources

### defguard_user

Retrieves information about a Defguard user.

- `username` - (Required) The username of the user to retrieve