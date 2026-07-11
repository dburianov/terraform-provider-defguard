package user

import (
	"context"
	"log"

	"terraform-provider-defguard/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceUserSchema(),
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating user")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	userData := map[string]interface{}{
		"username":   d.Get("username").(string),
		"first_name": d.Get("first_name").(string),
		"last_name":  d.Get("last_name").(string),
		"email":      d.Get("email").(string),
		"is_active":  d.Get("is_active").(bool),
		"is_admin":   d.Get("is_admin").(bool),
		"groups":     d.Get("groups").([]interface{}),
	}

	if phone, ok := d.GetOk("phone"); ok {
		userData["phone"] = phone.(string)
	}

	if password, ok := d.GetOk("password"); ok {
		userData["password"] = password.(string)
	}

	_, err := client.Post(ctx, "/api/v1/user", userData)
	if err != nil {
		return diag.Errorf("failed to create user: %v", err)
	}

	log.Printf("[DEBUG] Creating user with data: %+v", userData)

	d.SetId(d.Get("username").(string))

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading user")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	username := d.Id()
	log.Printf("[DEBUG] Reading user with username: %s", username)

	resp, err := client.Get(ctx, "/api/v1/user/"+username)
	if err != nil {
		return diag.Errorf("failed to read user: %v", err)
	}

	log.Printf("[DEBUG] Read user response: %s", string(resp.Body))

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating user")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	username := d.Id()

	userData := map[string]interface{}{
		"username":   d.Get("username").(string),
		"first_name": d.Get("first_name").(string),
		"last_name":  d.Get("last_name").(string),
		"email":      d.Get("email").(string),
		"is_active":  d.Get("is_active").(bool),
		"is_admin":   d.Get("is_admin").(bool),
		"groups":     d.Get("groups").([]interface{}),
	}

	if phone, ok := d.GetOk("phone"); ok {
		userData["phone"] = phone.(string)
	}

	_, err := client.Put(ctx, "/api/v1/user/"+username, userData)
	if err != nil {
		return diag.Errorf("failed to update user: %v", err)
	}

	log.Printf("[DEBUG] Updating user with data: %+v", userData)

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting user")

	client, ok := meta.(*client.Client)
	if !ok {
		return diag.Errorf("invalid client type")
	}

	username := d.Id()
	log.Printf("[DEBUG] Deleting user with username: %s", username)

	_, err := client.Delete(ctx, "/api/v1/user/"+username, nil)
	if err != nil {
		return diag.Errorf("failed to delete user: %v", err)
	}

	d.SetId("")

	return nil
}
