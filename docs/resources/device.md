## `defguard_device` Resource

Manages a WireGuard device in Defguard.

### Example Usage

```hcl
resource "defguard_device" "example" {
  name             = "example-device"
  user_id          = 1
  wireguard_pubkey = "pubkey_value_here"
}
```

### Schema

- `id` (Int64, Computed) - Device ID
- `name` (String, Required) - Device name
- `user_id` (Int64, Required) - User ID who owns this device
- `wireguard_pubkey` (String, Required) - WireGuard public key
- `created` (String, Computed) - Creation timestamp
