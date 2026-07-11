package user

import (
	"context"
	"log"

	"terraform-provider-defguard/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the user to retrieve",
			},
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

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	username := d.Get("username").(string)
	log.Printf("[DEBUG] Reading user with username: %s", username)

	resp, err := client.Get(ctx, "/api/v1/user/"+username)
	if err != nil {
		return diag.Errorf("failed to read user: %v", err)
	}

	log.Printf("[DEBUG] Read user response: %s", string(resp.Body))

	d.SetId(username)

	return nil
}
