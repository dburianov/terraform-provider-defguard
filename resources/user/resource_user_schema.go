package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
				Description: "The password for the user (optional, can be set later)",
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
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The groups the user belongs to",
			},
			"mfa_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether multi-factor authentication is enabled",
			},
			"enrolled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the user has completed enrollment",
			},
			"email_mfa_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether email-based MFA is enabled",
			},
			"totp_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether TOTP-based MFA is enabled",
			},
			"ldap_pass_requires_change": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether LDAP password requires change",
			},
			"mfa_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "None",
				Description: "The MFA method used by the user",
			},
			"authorized_apps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The authorized applications for the user",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oauth2client_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The OAuth2 client ID",
						},
						"oauth2client_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The OAuth2 client name",
						},
						"user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The user ID",
						},
					},
				},
			},
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the user",
			},
		},
	}
}
