package provider

import (
	"context"
	"fmt"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &OpenIDProviderResource{}

type OpenIDProviderResource struct {
	client *client.Client
}

type OpenIDProviderResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	BaseURL                    types.String `tfsdk:"base_url"`
	Kind                       types.String `tfsdk:"kind"`
	ClientID                   types.String `tfsdk:"client_id"`
	ClientSecret               types.String `tfsdk:"client_secret"`
	DirectorySyncEnabled       types.Bool   `tfsdk:"directory_sync_enabled"`
	DirectorySyncInterval      types.Int64  `tfsdk:"directory_sync_interval"`
	DirectorySyncUserBehavior  types.String `tfsdk:"directory_sync_user_behavior"`
	DirectorySyncAdminBehavior types.String `tfsdk:"directory_sync_admin_behavior"`
	DirectorySyncTarget        types.String `tfsdk:"directory_sync_target"`
	PrefetchUsers              types.Bool   `tfsdk:"prefetch_users"`
	CreateAccount              types.Bool   `tfsdk:"create_account"`
	UsernameHandling           types.String `tfsdk:"username_handling"`
	AdminEmail                 types.String `tfsdk:"admin_email"`
	DisplayName                types.String `tfsdk:"display_name"`
	GoogleServiceAccountEmail  types.String `tfsdk:"google_service_account_email"`
	GoogleServiceAccountKey    types.String `tfsdk:"google_service_account_key"`
	JumpcloudAPIKey            types.String `tfsdk:"jumpcloud_api_key"`
	OktaDirsyncClientID        types.String `tfsdk:"okta_dirsync_client_id"`
	OktaPrivateJWK             types.String `tfsdk:"okta_private_jwk"`
	DirectorySyncGroupMatch    types.String `tfsdk:"directory_sync_group_match"`
}

func NewOpenIDProviderResource() resource.Resource {
	return &OpenIDProviderResource{}
}

func (r *OpenIDProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openid_provider"
}

func (r *OpenIDProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "OpenID Provider resource represents an OpenID Connect provider in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Provider name (used as ID)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Provider name",
			},
			"base_url": schema.StringAttribute{
				Required:    true,
				Description: "Base URL of the OpenID provider",
			},
			"kind": schema.StringAttribute{
				Required:    true,
				Description: "Provider kind (Custom, Google, Microsoft, Okta, JumpCloud, Zitadel)",
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "Client ID",
			},
			"client_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Client secret",
			},
			"directory_sync_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable directory synchronization",
			},
			"directory_sync_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Directory sync interval in seconds",
			},
			"directory_sync_user_behavior": schema.StringAttribute{
				Optional:    true,
				Description: "User behavior for directory sync",
			},
			"directory_sync_admin_behavior": schema.StringAttribute{
				Optional:    true,
				Description: "Admin behavior for directory sync",
			},
			"directory_sync_target": schema.StringAttribute{
				Optional:    true,
				Description: "Directory sync target",
			},
			"prefetch_users": schema.BoolAttribute{
				Optional:    true,
				Description: "Prefetch users on startup",
			},
			"create_account": schema.BoolAttribute{
				Optional:    true,
				Description: "Create accounts for new users",
			},
			"username_handling": schema.StringAttribute{
				Optional:    true,
				Description: "Username handling strategy (RemoveForbidden, ReplaceForbidden, PruneEmailDomain)",
			},
			"admin_email": schema.StringAttribute{
				Optional:    true,
				Description: "Admin email address",
			},
			"display_name": schema.StringAttribute{
				Optional:    true,
				Description: "Display name for the provider",
			},
			"google_service_account_email": schema.StringAttribute{
				Optional:    true,
				Description: "Google service account email",
			},
			"google_service_account_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Google service account key",
			},
			"jumpcloud_api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "JumpCloud API key",
			},
			"okta_dirsync_client_id": schema.StringAttribute{
				Optional:    true,
				Description: "Okta directory sync client ID",
			},
			"okta_private_jwk": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Okta private JSON Web Key",
			},
			"directory_sync_group_match": schema.StringAttribute{
				Optional:    true,
				Description: "Group match pattern for directory sync",
			},
		},
	}
}

func (r *OpenIDProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *OpenIDProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OpenIDProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert all optional fields to appropriate types
	payload := map[string]interface{}{
		"name":                          plan.Name.ValueString(),
		"base_url":                      plan.BaseURL.ValueString(),
		"kind":                          plan.Kind.ValueString(),
		"client_id":                     plan.ClientID.ValueString(),
		"client_secret":                 plan.ClientSecret.ValueString(),
		"directory_sync_enabled":        plan.DirectorySyncEnabled.ValueBool(),
		"directory_sync_interval":       plan.DirectorySyncInterval.ValueInt64(),
		"directory_sync_user_behavior":  plan.DirectorySyncUserBehavior.ValueString(),
		"directory_sync_admin_behavior": plan.DirectorySyncAdminBehavior.ValueString(),
		"directory_sync_target":         plan.DirectorySyncTarget.ValueString(),
		"prefetch_users":                plan.PrefetchUsers.ValueBool(),
		"create_account":                plan.CreateAccount.ValueBool(),
		"username_handling":             plan.UsernameHandling.ValueString(),
	}

	if !plan.AdminEmail.IsUnknown() && !plan.AdminEmail.IsNull() {
		payload["admin_email"] = plan.AdminEmail.ValueString()
	}
	if !plan.DisplayName.IsUnknown() && !plan.DisplayName.IsNull() {
		payload["display_name"] = plan.DisplayName.ValueString()
	}
	if !plan.GoogleServiceAccountEmail.IsUnknown() && !plan.GoogleServiceAccountEmail.IsNull() {
		payload["google_service_account_email"] = plan.GoogleServiceAccountEmail.ValueString()
	}
	if !plan.GoogleServiceAccountKey.IsUnknown() && !plan.GoogleServiceAccountKey.IsNull() {
		payload["google_service_account_key"] = plan.GoogleServiceAccountKey.ValueString()
	}
	if !plan.JumpcloudAPIKey.IsUnknown() && !plan.JumpcloudAPIKey.IsNull() {
		payload["jumpcloud_api_key"] = plan.JumpcloudAPIKey.ValueString()
	}
	if !plan.OktaDirsyncClientID.IsUnknown() && !plan.OktaDirsyncClientID.IsNull() {
		payload["okta_dirsync_client_id"] = plan.OktaDirsyncClientID.ValueString()
	}
	if !plan.OktaPrivateJWK.IsUnknown() && !plan.OktaPrivateJWK.IsNull() {
		payload["okta_private_jwk"] = plan.OktaPrivateJWK.ValueString()
	}
	if !plan.DirectorySyncGroupMatch.IsUnknown() && !plan.DirectorySyncGroupMatch.IsNull() {
		payload["directory_sync_group_match"] = plan.DirectorySyncGroupMatch.ValueString()
	}

	_, err := r.client.Post(ctx, "/api/v1/openid/provider", payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating OpenID Provider", err.Error())
		return
	}

	// Extract provider info from response (name is used as ID)
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OpenIDProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OpenIDProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerName := state.Name.ValueString()
	path := fmt.Sprintf("/api/v1/openid/provider/%s", providerName)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading OpenID Provider", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse provider", err.Error())
		return
	}

	// Update state from response
	if name, ok := result["name"].(string); ok {
		state.Name = types.StringValue(name)
		state.ID = types.StringValue(name)
	}
	if baseURL, ok := result["base_url"].(string); ok {
		state.BaseURL = types.StringValue(baseURL)
	}
	if kind, ok := result["kind"].(string); ok {
		state.Kind = types.StringValue(kind)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OpenIDProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OpenIDProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OpenIDProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"name":                          plan.Name.ValueString(),
		"base_url":                      plan.BaseURL.ValueString(),
		"kind":                          plan.Kind.ValueString(),
		"client_id":                     plan.ClientID.ValueString(),
		"client_secret":                 plan.ClientSecret.ValueString(),
		"directory_sync_enabled":        plan.DirectorySyncEnabled.ValueBool(),
		"directory_sync_interval":       plan.DirectorySyncInterval.ValueInt64(),
		"directory_sync_user_behavior":  plan.DirectorySyncUserBehavior.ValueString(),
		"directory_sync_admin_behavior": plan.DirectorySyncAdminBehavior.ValueString(),
		"directory_sync_target":         plan.DirectorySyncTarget.ValueString(),
		"prefetch_users":                plan.PrefetchUsers.ValueBool(),
		"create_account":                plan.CreateAccount.ValueBool(),
		"username_handling":             plan.UsernameHandling.ValueString(),
	}

	if !plan.AdminEmail.IsUnknown() && !plan.AdminEmail.IsNull() {
		payload["admin_email"] = plan.AdminEmail.ValueString()
	}
	if !plan.DisplayName.IsUnknown() && !plan.DisplayName.IsNull() {
		payload["display_name"] = plan.DisplayName.ValueString()
	}
	if !plan.GoogleServiceAccountEmail.IsUnknown() && !plan.GoogleServiceAccountEmail.IsNull() {
		payload["google_service_account_email"] = plan.GoogleServiceAccountEmail.ValueString()
	}
	if !plan.GoogleServiceAccountKey.IsUnknown() && !plan.GoogleServiceAccountKey.IsNull() {
		payload["google_service_account_key"] = plan.GoogleServiceAccountKey.ValueString()
	}
	if !plan.JumpcloudAPIKey.IsUnknown() && !plan.JumpcloudAPIKey.IsNull() {
		payload["jumpcloud_api_key"] = plan.JumpcloudAPIKey.ValueString()
	}
	if !plan.OktaDirsyncClientID.IsUnknown() && !plan.OktaDirsyncClientID.IsNull() {
		payload["okta_dirsync_client_id"] = plan.OktaDirsyncClientID.ValueString()
	}
	if !plan.OktaPrivateJWK.IsUnknown() && !plan.OktaPrivateJWK.IsNull() {
		payload["okta_private_jwk"] = plan.OktaPrivateJWK.ValueString()
	}
	if !plan.DirectorySyncGroupMatch.IsUnknown() && !plan.DirectorySyncGroupMatch.IsNull() {
		payload["directory_sync_group_match"] = plan.DirectorySyncGroupMatch.ValueString()
	}

	respObj, err := r.client.Put(ctx, "/api/v1/openid/provider/"+plan.Name.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating OpenID Provider", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated provider", err.Error())
		return
	}

	// Update state from response
	if name, ok := result["name"].(string); ok {
		state.Name = types.StringValue(name)
		state.ID = types.StringValue(name)
	}
	if baseURL, ok := result["base_url"].(string); ok {
		state.BaseURL = types.StringValue(baseURL)
	}
	if kind, ok := result["kind"].(string); ok {
		state.Kind = types.StringValue(kind)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OpenIDProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OpenIDProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerName := state.Name.ValueString()
	path := fmt.Sprintf("/api/v1/openid/provider/%s", providerName)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting OpenID Provider", err.Error())
		return
	}
}

func (r *OpenIDProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
