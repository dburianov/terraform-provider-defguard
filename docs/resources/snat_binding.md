## `defguard_snat_binding` Resource

Manages a SNAT (Source Network Address Translation) binding in Defguard.

### Example Usage

```hcl
resource "defguard_snat_binding" "example" {
  user_id      = 1
  location_id  = 1
  public_ip    = "203.0.113.1"
}
```

### Schema

- `id` (Int64, Computed) - SNAT binding ID
- `user_id` (Int64, Required) - User ID to bind to the public IP
- `location_id` (Int64, Required) - WireGuard location ID
- `public_ip` (String, Required) - Public IP address for SNAT
