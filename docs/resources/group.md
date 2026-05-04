## `defguard_group` Resource

Manages a user group in Defguard.

### Example Usage

```hcl
resource "defguard_group" "admin" {
  name         = "admin"
  members      = ["user1", "user2"]
  is_admin     = true
  vpn_locations = ["location1", "location2"]
}
```

### Schema

- `id` (Int64, Computed) - Group ID
- `name` (String, Required) - Group name (requires replace on update)
- `members` (List of Strings, Required) - List of member usernames
- `is_admin` (Bool, Required) - Whether the group has admin privileges
- `vpn_locations` (List of Strings, Computed) - VPN locations associated with this group
