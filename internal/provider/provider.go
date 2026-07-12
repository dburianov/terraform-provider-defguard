package provider

import (
	"context"
	"fmt"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &defguardProvider{}

// nilIfUnknown returns nil if the string is unknown or null, otherwise returns the string value
func nilIfUnknown(s types.String) interface{} {
	if s.IsUnknown() || s.IsNull() {
		return nil
	}
	return s.ValueString()
}

type defguardProvider struct {
	version string
	client  *client.Client
}

type defguardProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIToken types.String `tfsdk:"api_token"` // API token for authentication (when cookie is not used)
	Cookie   types.String `tfsdk:"cookie"`    // Session cookie value for authentication
	Username types.String `tfsdk:"username"`  // Username for authentication
	Password types.String `tfsdk:"password"`  // Password for authentication
	Insecure types.Bool   `tfsdk:"insecure"`  // Skip TLS verification
}

func (p *defguardProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "defguard"
	resp.Version = p.version
}

func (p *defguardProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required:    true,
				Description: "Defguard API endpoint URL (e.g., https://defguard.example.com/api/v1)",
			},
			"api_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "API token for authentication. Cannot be used together with 'cookie'.",
			},
			"cookie": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Session cookie value for authentication. Cannot be used together with 'api_token'.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username for authentication. If provided with password, api_token/cookie is not required.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for authentication. Required when username is provided.",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Skip TLS verification (not recommended for production)",
			},
		},
	}
}

func (p *defguardProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data defguardProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := data.Endpoint.ValueString()

	apiToken := ""
	if !data.APIToken.IsUnknown() && !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
	}

	cookie := ""
	if !data.Cookie.IsUnknown() && !data.Cookie.IsNull() {
		cookie = data.Cookie.ValueString()
	}

	username := ""
	password := ""
	if !data.Username.IsUnknown() && !data.Username.IsNull() {
		username = data.Username.ValueString()
	}
	if !data.Password.IsUnknown() && !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	// Validate: cannot use both api_token and cookie
	if apiToken != "" && cookie != "" {
		resp.Diagnostics.AddError(
			"Authentication Configuration Error",
			"Cannot use 'api_token' and 'cookie' together. Please choose one authentication method.",
		)
		return
	}

	p.client = client.NewClient(endpoint, apiToken)

	// Set session cookie name (default: defguard_session)
	p.client.SetSessionCookie("defguard_session")

	// If username and password are provided, authenticate and store cookie
	if username != "" && password != "" {
		ctx := context.Background()
		_, err := p.client.LoginWithCredentials(ctx, username, password)
		if err != nil {
			resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Failed to login: %v", err))
			return
		}
	}

	// Set session cookie value if provided via 'cookie' field
	if cookie != "" {
		err := p.client.SetSessionValue(cookie)
		if err != nil {
			resp.Diagnostics.AddError("Cookie Error", fmt.Sprintf("Failed to set session cookie: %v", err))
			return
		}
	}

	resp.ResourceData = p.client
}

func (p *defguardProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
		NewGroupResource,
		NewNetworkResource,
		NewSNATBindingResource,
		NewUserResource,
		NewACLAliasResource,
		NewACLDestinationResource,
		NewACLRuleResource,
		NewOpenIDProviderResource,
	}
}

func (p *defguardProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &defguardProvider{
			version: version,
		}
	}
}
