# Defguard Terraform Provider

The Defguard Terraform provider is used to manage resources in a Defguard infrastructure through Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go tools:

```bash
go build -o terraform-provider-defguard
```

## Authentication

The provider supports two authentication methods:

### API Token

Use an API token for authentication:

```hcl
provider "defguard" {
  base_url  = "https://your-defguard-instance.com"
  api_token = "your-api-token"
}
```

Or set the environment variable `DEFGUARD_API_TOKEN`.

### Session Cookie

Use a session cookie for authentication:

```hcl
provider "defguard" {
  base_url  = "https://your-defguard-instance.com"
  session   = "your-session-cookie"
}
```

Or set the environment variable `DEFGUARD_SESSION_COOKIE`.

Both methods can be used independently. At least one is required.

## Using the provider

```hcl
terraform {
  required_providers {
    defguard = {
      source = "your-organization/defguard"
      version = "0.1.0"
    }
  }
}

provider "defguard" {
  base_url  = "https://your-defguard-instance.com"
  # Use either api_token or session
}
```

## Example Usage

```hcl
# Create a user
resource "defguard_user" "example_user" {
  username   = "exampleuser"
  first_name = "Example"
  last_name  = "User"
  email      = "example@example.com"
  is_active  = true
}

# Create a group
resource "defguard_group" "example_group" {
  name      = "example-group"
  is_admin  = false
  members   = [defguard_user.example_user.username]
}

# Create a device for the user
resource "defguard_device" "example_device" {
  name               = "example-device"
  wireguard_pubkey   = "examplepubkey123"
  username           = defguard_user.example_user.username
}
```

## Argument Reference

### Provider

- `base_url` - (Required) The base URL of the Defguard API
- `api_token` - (Optional) The API token for authentication. Can also be set via `DEFGUARD_API_TOKEN` environment variable.
- `session` - (Optional) The session cookie for authentication. Can also be set via `DEFGUARD_SESSION_COOKIE` environment variable. At least one of `api_token` or `session` is required.

### Resources

#### defguard_user

Manages a user in Defguard.

- `username` - (Required) The username of the user
- `first_name` - (Required) The first name of the user
- `last_name` - (Required) The last name of the user
- `email` - (Required) The email address of the user
- `password` - (Optional) The password for the user
- `phone` - (Optional) The phone number of the user
- `is_active` - (Optional) Whether the user is active (default: true)
- `is_admin` - (Optional) Whether the user is an administrator (default: false)
- `groups` - (Optional) The groups the user belongs to

#### defguard_group

Manages a group in Defguard.

- `name` - (Required) The name of the group
- `is_admin` - (Optional) Whether the group has admin privileges (default: false)
- `members` - (Optional) The members of the group

#### defguard_device

Manages a device in Defguard.

- `name` - (Required) The name of the device
- `wireguard_pubkey` - (Required) The WireGuard public key of the device
- `username` - (Required) The username of the user that owns the device
- `description` - (Optional) Description of the device

### Data Sources

#### defguard_user

Retrieves information about a user by username.

- `username` - (Required) The username of the user to retrieve

Outputs:

- `id` - The ID of the user

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a new Pull Request