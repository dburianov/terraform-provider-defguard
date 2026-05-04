package provider

import (
	"context"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &defguardProvider{}

type defguardProvider struct {
	version string
	client  *client.Client
}

type defguardProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIToken types.String `tfsdk:"api_token"`
	Insecure types.Bool   `tfsdk:"insecure"`
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
				Description: "API token for authentication. If not set, cookie authentication will be used.",
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

	p.client = client.NewClient(endpoint, apiToken)

	resp.ResourceData = p.client
}

func (p *defguardProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
		NewGroupResource,
		NewNetworkResource,
		NewSNATBindingResource,
		NewUserResource,
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
