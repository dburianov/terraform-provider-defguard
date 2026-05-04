terraform {
  required_providers {
    defguard = {
      source = "registry.terraform.io/dburianov/defguard"
    }
  }
}

provider "defguard" {
  url     = "http://localhost:8000"
  api_key = var.defguard_api_key
  session = var.defguard_session
}

variable "defguard_api_key" {
  type        = string
  description = "Defguard API key"
}

variable "defguard_session" {
  type        = string
  description = "Defguard session cookie"
}

# Create a new group
resource "defguard_group" "admin" {
  name     = "admin"
  members  = ["admin_user"]
  is_admin = true
}

# Create a new user
resource "defguard_user" "example" {
  username   = "johndoe"
  first_name = "John"
  last_name  = "Doe"
  email      = "john.doe@example.com"
  phone      = "+1234567890"
  is_admin   = false
  is_active  = true
  groups     = ["admin"]
}

# Create a new network
resource "defguard_network" "main" {
  name                      = "main-network"
  address                   = "10.0.0.0/24"
  port                      = 51820
  pubkey                    = "pubkey_value"
  endpoint                  = "endpoint_value"
  allowed_ips               = "0.0.0.0/0"
  allowed_groups            = ["admin"]
  dns                       = "8.8.8.8"
  keepalive_interval        = 25
  peer_disconnect_threshold = 120
  acl_enabled               = false
  acl_default_allow         = true
  location_mfa_mode         = "disabled"
  service_location_mode     = "disabled"
}

# Create a new device for the user
resource "defguard_device" "example" {
  name             = "laptop"
  user_id          = defguard_user.example.id
  wireguard_pubkey = "device_pubkey_value"
}

# Create a SNAT binding
resource "defguard_snat_binding" "example" {
  user_id     = defguard_user.example.id
  location_id = 1
  public_ip   = "203.0.113.1"
}
