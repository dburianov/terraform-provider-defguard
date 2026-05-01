package user

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceUserSchema().Schema,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating user")

	// Prepare the user data
	userData := map[string]interface{}{
		"username":                  d.Get("username").(string),
		"first_name":                d.Get("first_name").(string),
		"last_name":                 d.Get("last_name").(string),
		"email":                     d.Get("email").(string),
		"phone":                     d.Get("phone").(string),
		"password":                  d.Get("password").(string),
		"is_active":                 d.Get("is_active").(bool),
		"is_admin":                  d.Get("is_admin").(bool),
		"groups":                    d.Get("groups").(*schema.Set).List(),
		"mfa_enabled":               d.Get("mfa_enabled").(bool),
		"enrolled":                  d.Get("enrolled").(bool),
		"email_mfa_enabled":         d.Get("email_mfa_enabled").(bool),
		"totp_enabled":              d.Get("totp_enabled").(bool),
		"ldap_pass_requires_change": d.Get("ldap_pass_requires_change").(bool),
		"mfa_method":                d.Get("mfa_method").(string),
	}

	// Make API call to create user
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	log.Printf("[DEBUG] Creating user with data: %+v", userData)

	// For now, just simulate the creation
	// In a real implementation, we would:
	// 1. Make HTTP POST request to /api/v1/user
	// 2. Parse the response
	// 3. Set the ID and other computed fields

	// Set the ID (in real implementation, this would come from API response)
	d.SetId("12345")

	// Read the user to populate all fields
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading user")

	// In a real implementation, we would:
	// 1. Make HTTP GET request to /api/v1/user/{username}
	// 2. Parse the response
	// 3. Set all the fields in the Terraform state

	// For now, we'll just return the current state
	username := d.Get("username").(string)
	log.Printf("[DEBUG] Reading user with username: %s", username)

	// In a real implementation, we would set the fields from the API response
	// For now, just simulate that we read the user successfully
	d.Set("username", username)

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating user")

	// Prepare the user data for update
	userData := map[string]interface{}{
		"username":                  d.Get("username").(string),
		"first_name":                d.Get("first_name").(string),
		"last_name":                 d.Get("last_name").(string),
		"email":                     d.Get("email").(string),
		"phone":                     d.Get("phone").(string),
		"is_active":                 d.Get("is_active").(bool),
		"is_admin":                  d.Get("is_admin").(bool),
		"groups":                    d.Get("groups").(*schema.Set).List(),
		"mfa_enabled":               d.Get("mfa_enabled").(bool),
		"enrolled":                  d.Get("enrolled").(bool),
		"email_mfa_enabled":         d.Get("email_mfa_enabled").(bool),
		"totp_enabled":              d.Get("totp_enabled").(bool),
		"ldap_pass_requires_change": d.Get("ldap_pass_requires_change").(bool),
		"mfa_method":                d.Get("mfa_method").(string),
	}

	// Make API call to update user
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	log.Printf("[DEBUG] Updating user with data: %+v", userData)

	// In a real implementation, we would:
	// 1. Make HTTP PUT request to /api/v1/user/{username}
	// 2. Parse the response
	// 3. Update the Terraform state

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting user")

	// Make API call to delete user
	// This is a placeholder - actual implementation would make HTTP requests
	// to the defguard API endpoints
	username := d.Get("username").(string)
	log.Printf("[DEBUG] Deleting user with username: %s", username)

	// In a real implementation, we would:
	// 1. Make HTTP DELETE request to /api/v1/user/{username}
	// 2. Handle any errors
	// 3. Clear the ID from Terraform state

	// Clear the ID so Terraform knows the resource is deleted
	d.SetId("")

	return nil
}

// Config holds the provider configuration for user operations
type Config struct {
	BaseURL  string
	APIToken string
}
