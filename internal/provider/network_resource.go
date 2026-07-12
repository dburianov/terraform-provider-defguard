package provider

import (
	"context"
	"fmt"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &NetworkResource{}

type NetworkResource struct {
	client *client.Client
}

type NetworkResourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Address                 types.String `tfsdk:"address"`
	Port                    types.Int64  `tfsdk:"port"`
	Pubkey                  types.String `tfsdk:"pubkey"`
	Endpoint                types.String `tfsdk:"endpoint"`
	AllowedIPs              types.String `tfsdk:"allowed_ips"`
	AllowedGroups           types.List   `tfsdk:"allowed_groups"`
	DNS                     types.String `tfsdk:"dns"`
	KeepaliveInterval       types.Int64  `tfsdk:"keepalive_interval"`
	PeerDisconnectThreshold types.Int64  `tfsdk:"peer_disconnect_threshold"`
	ACLEnabled              types.Bool   `tfsdk:"acl_enabled"`
	ACLDefaultAllow         types.Bool   `tfsdk:"acl_default_allow"`
	LocationMFAMode         types.String `tfsdk:"location_mfa_mode"`
	ServiceLocationMode     types.String `tfsdk:"service_location_mode"`
	Connected               types.Bool   `tfsdk:"connected"`
}

func NewNetworkResource() resource.Resource {
	return &NetworkResource{}
}

func (r *NetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *NetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network resource represents a WireGuard network in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Network ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Network name",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "Network address (CIDR format)",
			},
			"port": schema.Int64Attribute{
				Required:    true,
				Description: "Network port",
			},
			"pubkey": schema.StringAttribute{
				Required:    true,
				Description: "Network public key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"endpoint": schema.StringAttribute{
				Required:    true,
				Description: "Network endpoint",
			},
			"allowed_ips": schema.StringAttribute{
				Required:    true,
				Description: "Allowed IP ranges (CIDR format, comma-separated)",
			},
			"allowed_groups": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Groups allowed to connect",
			},
			"dns": schema.StringAttribute{
				Optional:    true,
				Description: "DNS server",
			},
			"keepalive_interval": schema.Int64Attribute{
				Required:    true,
				Description: "Keepalive interval in seconds",
			},
			"peer_disconnect_threshold": schema.Int64Attribute{
				Required:    true,
				Description: "Peer disconnect threshold in seconds",
			},
			"acl_enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether ACL is enabled",
			},
			"acl_default_allow": schema.BoolAttribute{
				Required:    true,
				Description: "Default ACL behavior (allow or deny)",
			},
			"location_mfa_mode": schema.StringAttribute{
				Required:    true,
				Description: "MFA mode for locations (disabled, internal, external)",
			},
			"service_location_mode": schema.StringAttribute{
				Required:    true,
				Description: "Service location mode (disabled, prelogon, alwayson)",
			},
			"connected": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the gateway is connected",
			},
		},
	}
}

func (r *NetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert allowed_groups to []string
	var allowedGroups []string
	resp.Diagnostics.Append(plan.AllowedGroups.ElementsAs(ctx, &allowedGroups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create network payload based on OpenAPI WireguardNetworkData schema
	payload := map[string]interface{}{
		"name":                      plan.Name.ValueString(),
		"address":                   plan.Address.ValueString(),
		"port":                      plan.Port.ValueInt64(),
		"pubkey":                    plan.Pubkey.ValueString(),
		"endpoint":                  plan.Endpoint.ValueString(),
		"allowed_ips":               plan.AllowedIPs.ValueString(),
		"allowed_groups":            allowedGroups,
		"dns":                       nilIfUnknown(plan.DNS),
		"keepalive_interval":        plan.KeepaliveInterval.ValueInt64(),
		"peer_disconnect_threshold": plan.PeerDisconnectThreshold.ValueInt64(),
		"acl_enabled":               plan.ACLEnabled.ValueBool(),
		"acl_default_allow":         plan.ACLDefaultAllow.ValueBool(),
		"location_mfa_mode":         plan.LocationMFAMode.ValueString(),
		"service_location_mode":     plan.ServiceLocationMode.ValueString(),
	}

	respObj, err := r.client.Post(ctx, "/api/v1/network", payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Network", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", err.Error())
		return
	}

	// Extract network info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/network/%d", networkID)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading Network", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse network", err.Error())
		return
	}

	// Update state from response
	if id, ok := result["id"].(float64); ok {
		state.ID = types.Int64Value(int64(id))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state NetworkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert allowed_groups to []string
	var allowedGroups []string
	resp.Diagnostics.Append(plan.AllowedGroups.ElementsAs(ctx, &allowedGroups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkID := plan.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/network/%d", networkID)

	// Create update payload based on OpenAPI WireguardNetworkData schema
	payload := map[string]interface{}{
		"name":                      plan.Name.ValueString(),
		"address":                   plan.Address.ValueString(),
		"port":                      plan.Port.ValueInt64(),
		"pubkey":                    plan.Pubkey.ValueString(),
		"endpoint":                  plan.Endpoint.ValueString(),
		"allowed_ips":               plan.AllowedIPs.ValueString(),
		"allowed_groups":            allowedGroups,
		"dns":                       nilIfUnknown(plan.DNS),
		"keepalive_interval":        plan.KeepaliveInterval.ValueInt64(),
		"peer_disconnect_threshold": plan.PeerDisconnectThreshold.ValueInt64(),
		"acl_enabled":               plan.ACLEnabled.ValueBool(),
		"acl_default_allow":         plan.ACLDefaultAllow.ValueBool(),
		"location_mfa_mode":         plan.LocationMFAMode.ValueString(),
		"service_location_mode":     plan.ServiceLocationMode.ValueString(),
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Network", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated network", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/network/%d", networkID)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Network", err.Error())
		return
	}
}

func (r *NetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
