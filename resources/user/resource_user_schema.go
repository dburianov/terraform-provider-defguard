package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The username of the user",
		},
		"first_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The first name of the user",
		},
		"last_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The last name of the user",
		},
		"email": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The email address of the user",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The password for the user",
		},
		"phone": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The phone number of the user",
		},
		"is_active": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Whether the user is active",
		},
		"is_admin": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether the user is an administrator",
		},
		"groups": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The groups the user belongs to",
		},
	}
}
