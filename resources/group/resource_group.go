package group

import (
	"context"
	"log"

	"terraform-provider-defguard/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceGroupSchema(),
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating group")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	// Prepare the group data
	groupData := map[string]interface{}{
		"name":     d.Get("name").(string),
		"is_admin": d.Get("is_admin").(bool),
		"members":  d.Get("members").([]interface{}),
	}

	_, err := client.Post(ctx, "/api/v1/group", groupData)
	if err != nil {
		return diag.Errorf("failed to create group: %v", err)
	}

	log.Printf("[DEBUG] Creating group with data: %+v", groupData)

	// Set the ID from name (since we don't get it back from API in this flow)
	d.SetId(d.Get("name").(string))

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading group")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	groupID := d.Id()
	log.Printf("[DEBUG] Reading group with ID: %s", groupID)

	resp, err := client.Get(ctx, "/api/v1/group-info")
	if err != nil {
		return diag.Errorf("failed to read group: %v", err)
	}

	log.Printf("[DEBUG] Read group response: %s", string(resp.Body))

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating group")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	groupID := d.Id()

	// Prepare the group data for update
	groupData := map[string]interface{}{
		"name":     d.Get("name").(string),
		"is_admin": d.Get("is_admin").(bool),
		"members":  d.Get("members").([]interface{}),
	}

	_, err := client.Put(ctx, "/api/v1/group/"+groupID, groupData)
	if err != nil {
		return diag.Errorf("failed to update group: %v", err)
	}

	log.Printf("[DEBUG] Updating group with data: %+v", groupData)

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting group")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	groupID := d.Id()
	log.Printf("[DEBUG] Deleting group with ID: %s", groupID)

	_, err := client.Delete(ctx, "/api/v1/group/"+groupID, nil)
	if err != nil {
		return diag.Errorf("failed to delete group: %v", err)
	}

	d.SetId("")

	return nil
}
