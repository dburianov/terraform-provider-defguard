package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API token for authentication",
				DefaultFunc: schema.EnvDefaultFunc("DEFGUARD_API_TOKEN", nil),
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base URL for defguard API",
				DefaultFunc: schema.EnvDefaultFunc("DEFGUARD_BASE_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"defguard_user":   &schema.Resource{},
			"defguard_group":  &schema.Resource{},
			"defguard_device": &schema.Resource{},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"defguard_user": &schema.Resource{},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	apiToken := d.Get("api_token").(string)

	if baseURL == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Base URL is required",
			Detail:   "Base URL must be provided for defguard API",
		})
		return nil, diags
	}

	if apiToken == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "API token is required",
			Detail:   "API token must be provided for defguard API",
		})
		return nil, diags
	}

	// Create API client with the provided configuration
	config := &Config{
		BaseURL:  baseURL,
		APIToken: apiToken,
	}

	return config, diags
}

// Config holds the provider configuration
type Config struct {
	BaseURL  string
	APIToken string
}
