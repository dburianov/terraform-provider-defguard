package provider

import (
	"context"

	userDataSource "terraform-provider-defguard/data-sources/user"
	device "terraform-provider-defguard/resources/device"
	group "terraform-provider-defguard/resources/group"
	userResource "terraform-provider-defguard/resources/user"

	"terraform-provider-defguard/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API token for authentication",
				DefaultFunc: schema.EnvDefaultFunc("DEFGUARD_API_TOKEN", nil),
			},
			"session": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Session cookie for authentication",
				DefaultFunc: schema.EnvDefaultFunc("DEFGUARD_SESSION_COOKIE", nil),
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base URL for defguard API",
				DefaultFunc: schema.EnvDefaultFunc("DEFGUARD_BASE_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"defguard_user":   userResource.ResourceUser(),
			"defguard_group":  group.ResourceGroup(),
			"defguard_device": device.ResourceDevice(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"defguard_user": userDataSource.DataSourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	apiToken := d.Get("api_token").(string)
	sessionCookie := d.Get("session").(string)

	if baseURL == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Base URL is required",
			Detail:   "Base URL must be provided for defguard API",
		})
		return nil, diags
	}

	if apiToken == "" && sessionCookie == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Authentication is required",
			Detail:   "API token or session cookie must be provided for defguard API",
		})
		return nil, diags
	}

	// Create API client with the provided configuration
	client := client.NewClient(baseURL, apiToken, sessionCookie)

	return client, diags
}
