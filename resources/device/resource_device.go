package device

import (
	"context"
	"log"

	"terraform-provider-defguard/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDevice() *schema.Resource {
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

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	username := d.Get("username").(string)

	deviceData := map[string]interface{}{
		"name":             d.Get("name").(string),
		"wireguard_pubkey": d.Get("wireguard_pubkey").(string),
	}

	_, err := client.Post(ctx, "/api/v1/device/user/"+username, deviceData)
	if err != nil {
		return diag.Errorf("failed to create device: %v", err)
	}

	log.Printf("[DEBUG] Creating device with data: %+v", deviceData)

	d.SetId(d.Get("name").(string))

	return resourceDeviceRead(ctx, d, meta)
}

func resourceDeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading device")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	deviceID := d.Id()
	log.Printf("[DEBUG] Reading device with ID: %s", deviceID)

	resp, err := client.Get(ctx, "/api/v1/device/"+deviceID)
	if err != nil {
		return diag.Errorf("failed to read device: %v", err)
	}

	log.Printf("[DEBUG] Read device response: %s", string(resp.Body))

	return nil
}

func resourceDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating device")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	deviceID := d.Id()

	deviceData := map[string]interface{}{
		"name":             d.Get("name").(string),
		"wireguard_pubkey": d.Get("wireguard_pubkey").(string),
	}

	if description, ok := d.GetOk("description"); ok {
		deviceData["description"] = description.(string)
	}

	_, err := client.Put(ctx, "/api/v1/device/"+deviceID, deviceData)
	if err != nil {
		return diag.Errorf("failed to update device: %v", err)
	}

	log.Printf("[DEBUG] Updating device with data: %+v", deviceData)

	return resourceDeviceRead(ctx, d, meta)
}

func resourceDeviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting device")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	deviceID := d.Id()
	log.Printf("[DEBUG] Deleting device with ID: %s", deviceID)

	_, err := client.Delete(ctx, "/api/v1/device/"+deviceID, nil)
	if err != nil {
		return diag.Errorf("failed to delete device: %v", err)
	}

	d.SetId("")

	return nil
}
