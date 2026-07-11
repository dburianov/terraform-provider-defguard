package device

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeviceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the device",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the device",
			},
			"wireguard_pubkey": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The WireGuard public key",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The user ID that owns the device",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username of the device owner",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Device description",
			},
			"configured": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the device is configured and ready to use",
			},
			"device_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of device (user or network)",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp",
			},
		},
	}
}
