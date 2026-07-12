## `defguard_network` Resource

Manages a WireGuard network in Defguard.

### Example Usage

```hcl
resource "defguard_network" "main" {
  name                    = "main-network"
  address                 = "10.0.0.0/24"
  port                    = 51820
  pubkey                  = "pubkey_value"
  endpoint                = "endpoint_value"
  allowed_ips             = "0.0.0.0/0"
  allowed_groups          = ["group1", "group2"]
  dns                     = "8.8.8.8"
  keepalive_interval      = 25
  peer_disconnect_threshold = 120
  acl_enabled             = false
  acl_default_allow       = true
  location_mfa_mode       = "disabled"
  service_location_mode   = "disabled"
}
```

### Schema

- `id` (Int64, Computed) - Network ID
- `name` (String, Required) - Network name
- `address` (String, Required) - Network address (CIDR format)
- `port` (Int64, Required) - Network port
- `pubkey` (String, Required) - Network public key (requires replace on update)
- `endpoint` (String, Required) - Network endpoint
- `allowed_ips` (String, Required) - Allowed IP ranges (CIDR format, comma-separated)
- `allowed_groups` (List of Strings, Required) - Groups allowed to connect
- `dns` (String, Optional) - DNS server
- `keepalive_interval` (Int64, Required) - Keepalive interval in seconds
- `peer_disconnect_threshold` (Int64, Required) - Peer disconnect threshold in seconds
- `acl_enabled` (Bool, Required) - Whether ACL is enabled
- `acl_default_allow` (Bool, Required) - Default ACL behavior (allow or deny)
- `location_mfa_mode` (String, Required) - MFA mode for locations (disabled, internal, external)
- `service_location_mode` (String, Required) - Service location mode (disabled, prelogon, alwayson)
- `connected` (Bool, Computed) - Whether the gateway is connected
