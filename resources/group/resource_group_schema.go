package group

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the group",
		},
		"is_admin": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether the group has admin privileges",
		},
		"members": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The members of the group",
		},
	}
}
