package device

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeviceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the device",
			},
			"wireguard_pubkey": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The WireGuard public key of the device",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the user that owns the device",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username of the user that owns the device",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the device",
			},
			"configured": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the device is configured",
			},
			"device_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "user",
				Description: "The type of device (user or network)",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation timestamp of the device",
			},
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the device",
			},
		},
	}
}
