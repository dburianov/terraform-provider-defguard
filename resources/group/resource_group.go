package group

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceGroupSchema().Schema,
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating group")

	// Prepare the group data
	groupData := map[string]interface{}{
		"name":          d.Get("name").(string),
		"is_admin":      d.Get("is_admin").(bool),
		"members":       d.Get("members").(*schema.Set).List(),
		"vpn_locations": d.Get("vpn_locations").(*schema.Set).List(),
	}

	// Make API call to create group
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	log.Printf("[DEBUG] Creating group with data: %+v", groupData)

	// For now, just simulate the creation
	// In a real implementation, we would:
	// 1. Make HTTP POST request to /api/v1/group
	// 2. Parse the response
	// 3. Set the ID and other computed fields

	// Set the ID (in real implementation, this would come from API response)
	d.SetId("12345")

	// Read the group to populate all fields
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading group")

	// In a real implementation, we would:
	// 1. Make HTTP GET request to /api/v1/group/{name}
	// 2. Parse the response
	// 3. Set all the fields in the Terraform state

	// For now, we'll just return the current state
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Reading group with name: %s", name)

	// In a real implementation, we would set the fields from the API response
	// For now, just simulate that we read the group successfully
	d.Set("name", name)

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating group")

	// Prepare the group data for update
	groupData := map[string]interface{}{
		"name":          d.Get("name").(string),
		"is_admin":      d.Get("is_admin").(bool),
		"members":       d.Get("members").(*schema.Set).List(),
		"vpn_locations": d.Get("vpn_locations").(*schema.Set).List(),
	}

	// Make API call to update group
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	log.Printf("[DEBUG] Updating group with data: %+v", groupData)

	// In a real implementation, we would:
	// 1. Make HTTP PUT request to /api/v1/group/{name}
	// 2. Parse the response
	// 3. Update the Terraform state

	return nil
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting group")

	// Make API call to delete group
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting group with name: %s", name)

	// In a real implementation, we would:
	// 1. Make HTTP DELETE request to /api/v1/group/{name}
	// 2. Handle any errors
	// 3. Clear the ID from Terraform state

	// Clear the ID so Terraform knows the resource is deleted
	d.SetId("")

	return nil
}
