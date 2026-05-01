package device

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceCreate,
		ReadContext:   resourceDeviceRead,
		UpdateContext: resourceDeviceUpdate,
		DeleteContext: resourceDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceDeviceSchema().Schema,
	}
}

func resourceDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating device")

	// Prepare the device data
	deviceData := map[string]interface{}{
		"name":             d.Get("name").(string),
		"wireguard_pubkey": d.Get("wireguard_pubkey").(string),
		"user_id":          d.Get("user_id").(int),
		"username":         d.Get("username").(string),
		"description":      d.Get("description").(string),
		"configured":       d.Get("configured").(bool),
		"device_type":      d.Get("device_type").(string),
	}

	// Make API call to create device
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	log.Printf("[DEBUG] Creating device with data: %+v", deviceData)

	// For now, just simulate the creation
	// In a real implementation, we would:
	// 1. Make HTTP POST request to /api/v1/device/{username}
	// 2. Parse the response
	// 3. Set the ID and other computed fields

	// Set the ID (in real implementation, this would come from API response)
	d.SetId("12345")

	// Read the device to populate all fields
	return resourceDeviceRead(ctx, d, meta)
}

func resourceDeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading device")

	// In a real implementation, we would:
	// 1. Make HTTP GET request to /api/v1/device/{device_id}
	// 2. Parse the response
	// 3. Set all the fields in the Terraform state

	// For now, we'll just return the current state
	id := d.Get("id").(int)
	log.Printf("[DEBUG] Reading device with ID: %d", id)

	// In a real implementation, we would set the fields from the API response
	// For now, just simulate that we read the device successfully
	d.Set("id", id)

	return nil
}

func resourceDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating device")

	// Prepare the device data for update
	deviceData := map[string]interface{}{
		"name":             d.Get("name").(string),
		"wireguard_pubkey": d.Get("wireguard_pubkey").(string),
		"user_id":          d.Get("user_id").(int),
		"username":         d.Get("username").(string),
		"description":      d.Get("description").(string),
		"configured":       d.Get("configured").(bool),
		"device_type":      d.Get("device_type").(string),
	}

	// Make API call to update device
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	log.Printf("[DEBUG] Updating device with data: %+v", deviceData)

	// In a real implementation, we would:
	// 1. Make HTTP PUT request to /api/v1/device/{device_id}
	// 2. Parse the response
	// 3. Update the Terraform state

	return nil
}

func resourceDeviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting device")

	// Make API call to delete device
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	id := d.Get("id").(int)
	log.Printf("[DEBUG] Deleting device with ID: %d", id)

	// In a real implementation, we would:
	// 1. Make HTTP DELETE request to /api/v1/device/{device_id}
	// 2. Handle any errors
	// 3. Clear the ID from Terraform state

	// Clear the ID so Terraform knows the resource is deleted
	d.SetId("")

	return nil
}
