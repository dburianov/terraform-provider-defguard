## `defguard_user` Resource

Manages a user in Defguard.

### Example Usage

```hcl
resource "defguard_user" "example" {
  username               = "johndoe"
  first_name             = "John"
  last_name              = "Doe"
  email                  = "john.doe@example.com"
  phone                  = "+1234567890"
  is_admin               = false
  is_active              = true
  groups                 = ["users", "developers"]
}
```

### Schema

- `id` (Int64, Computed) - User ID
- `username` (String, Required) - Username (unique, requires replace on update)
- `first_name` (String, Required) - User's first name
- `last_name` (String, Required) - User's last name
- `email` (String, Required) - User's email address
- `phone` (String, Optional) - User's phone number
- `is_admin` (Bool, Required) - Whether the user has admin privileges
- `is_active` (Bool, Required) - Whether the user account is active
- `enrolled` (Bool, Computed) - Whether the user has completed enrollment
- `mfa_enabled` (Bool, Computed) - Whether MFA is enabled for the user
- `totp_enabled` (Bool, Computed) - Whether TOTP is enabled for the user
- `email_mfa_enabled` (Bool, Computed) - Whether email MFA is enabled for the user
- `mfa_method` (String, Computed) - Current MFA method
- `authorized_apps` (List of Strings, Computed) - List of authorized OAuth2 apps
- `groups` (List of Strings, Required) - Groups the user belongs to
- `ldap_pass_requires_change` (Bool, Computed) - Whether LDAP password requires change
