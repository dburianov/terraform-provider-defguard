package user

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the user to retrieve",
			},
			// Add other computed fields as needed
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the user",
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading user data source")

	// In a real implementation, we would:
	// 1. Make HTTP GET request to /api/v1/user/{username}
	// 2. Parse the response
	// 3. Set all the fields in the Terraform state

	// For now, just simulate that we read the user successfully
	username := d.Get("username").(string)
	d.SetId(username)

	return nil
}
