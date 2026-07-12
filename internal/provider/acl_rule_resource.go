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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ACLRuleResource{}

type ACLRuleResource struct {
	client *client.Client
}

type ACLRuleResourceModel struct {
	ID                           types.Int64  `tfsdk:"id"`
	Name                         types.String `tfsdk:"name"`
	Enabled                      types.Bool   `tfsdk:"enabled"`
	AllLocations                 types.Bool   `tfsdk:"all_locations"`
	Locations                    types.List   `tfsdk:"locations"`
	AllowAllUsers                types.Bool   `tfsdk:"allow_all_users"`
	DenyAllUsers                 types.Bool   `tfsdk:"deny_all_users"`
	AllowAllGroups               types.Bool   `tfsdk:"allow_all_groups"`
	DenyAllGroups                types.Bool   `tfsdk:"deny_all_groups"`
	AllowAllNetworkDevices       types.Bool   `tfsdk:"allow_all_network_devices"`
	DenyAllNetworkDevices        types.Bool   `tfsdk:"deny_all_network_devices"`
	AllowedUsers                 types.List   `tfsdk:"allowed_users"`
	DeniedUsers                  types.List   `tfsdk:"denied_users"`
	AllowedGroups                types.List   `tfsdk:"allowed_groups"`
	DeniedGroups                 types.List   `tfsdk:"denied_groups"`
	AllowedNetworkDevices        types.List   `tfsdk:"allowed_network_devices"`
	DeniedNetworkDevices         types.List   `tfsdk:"denied_network_devices"`
	UseManualDestinationSettings types.Bool   `tfsdk:"use_manual_destination_settings"`
	Addresses                    types.String `tfsdk:"addresses"`
	Ports                        types.String `tfsdk:"ports"`
	Protocols                    types.List   `tfsdk:"protocols"`
	AnyAddress                   types.Bool   `tfsdk:"any_address"`
	AnyPort                      types.Bool   `tfsdk:"any_port"`
	AnyProtocol                  types.Bool   `tfsdk:"any_protocol"`
	Aliases                      types.List   `tfsdk:"aliases"`
	Destinations                 types.List   `tfsdk:"destinations"`
	Expires                      types.String `tfsdk:"expires"`
}

func NewACLRuleResource() resource.Resource {
	return &ACLRuleResource{}
}

func (r *ACLRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl_rule"
}

func (r *ACLRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "ACL rule resource represents an ACL rule in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ACL rule ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Rule name",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the rule is enabled",
			},
			"all_locations": schema.BoolAttribute{
				Optional:    true,
				Description: "Apply to all locations",
			},
			"locations": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Location IDs to apply this rule (if not all_locations)",
			},
			"allow_all_users": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow all users",
			},
			"deny_all_users": schema.BoolAttribute{
				Optional:    true,
				Description: "Deny all users",
			},
			"allow_all_groups": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow all groups",
			},
			"deny_all_groups": schema.BoolAttribute{
				Optional:    true,
				Description: "Deny all groups",
			},
			"allow_all_network_devices": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow all network devices",
			},
			"deny_all_network_devices": schema.BoolAttribute{
				Optional:    true,
				Description: "Deny all network devices",
			},
			"allowed_users": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Allowed user IDs",
			},
			"denied_users": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Denied user IDs",
			},
			"allowed_groups": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Allowed group IDs",
			},
			"denied_groups": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Denied group IDs",
			},
			"allowed_network_devices": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Allowed network device IDs",
			},
			"denied_network_devices": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Denied network device IDs",
			},
			"use_manual_destination_settings": schema.BoolAttribute{
				Optional:    true,
				Description: "Use manual destination settings",
			},
			"addresses": schema.StringAttribute{
				Optional:    true,
				Description: "Manual addresses (comma-separated)",
			},
			"ports": schema.StringAttribute{
				Optional:    true,
				Description: "Manual ports (comma-separated)",
			},
			"protocols": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Protocol IDs",
			},
			"any_address": schema.BoolAttribute{
				Optional:    true,
				Description: "Match any address",
			},
			"any_port": schema.BoolAttribute{
				Optional:    true,
				Description: "Match any port",
			},
			"any_protocol": schema.BoolAttribute{
				Optional:    true,
				Description: "Match any protocol",
			},
			"aliases": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Alias IDs to use",
			},
			"destinations": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Destination IDs to use",
			},
			"expires": schema.StringAttribute{
				Optional:    true,
				Description: "Rule expiration date (RFC3339 format)",
			},
		},
	}
}

func (r *ACLRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ACLRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ACLRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert lists to []int64
	var locations, allowedUsers, deniedUsers, allowedGroups, deniedGroups,
		allowedNetworkDevices, deniedNetworkDevices, protocols, aliases, destinations []int64

	resp.Diagnostics.Append(plan.Locations.ElementsAs(ctx, &locations, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.AllowedUsers.ElementsAs(ctx, &allowedUsers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.DeniedUsers.ElementsAs(ctx, &deniedUsers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.AllowedGroups.ElementsAs(ctx, &allowedGroups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.DeniedGroups.ElementsAs(ctx, &deniedGroups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.AllowedNetworkDevices.ElementsAs(ctx, &allowedNetworkDevices, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.DeniedNetworkDevices.ElementsAs(ctx, &deniedNetworkDevices, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.Protocols.ElementsAs(ctx, &protocols, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.Aliases.ElementsAs(ctx, &aliases, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create rule payload based on OpenAPI EditAclRule schema
	payload := map[string]interface{}{
		"name":                            plan.Name.ValueString(),
		"enabled":                         plan.Enabled.ValueBool(),
		"all_locations":                   plan.AllLocations.ValueBool(),
		"locations":                       locations,
		"allow_all_users":                 plan.AllowAllUsers.ValueBool(),
		"deny_all_users":                  plan.DenyAllUsers.ValueBool(),
		"allow_all_groups":                plan.AllowAllGroups.ValueBool(),
		"deny_all_groups":                 plan.DenyAllGroups.ValueBool(),
		"allow_all_network_devices":       plan.AllowAllNetworkDevices.ValueBool(),
		"deny_all_network_devices":        plan.DenyAllNetworkDevices.ValueBool(),
		"allowed_users":                   allowedUsers,
		"denied_users":                    deniedUsers,
		"allowed_groups":                  allowedGroups,
		"denied_groups":                   deniedGroups,
		"allowed_network_devices":         allowedNetworkDevices,
		"denied_network_devices":          deniedNetworkDevices,
		"use_manual_destination_settings": plan.UseManualDestinationSettings.ValueBool(),
		"addresses":                       plan.Addresses.ValueString(),
		"ports":                           plan.Ports.ValueString(),
		"protocols":                       protocols,
		"any_address":                     plan.AnyAddress.ValueBool(),
		"any_port":                        plan.AnyPort.ValueBool(),
		"any_protocol":                    plan.AnyProtocol.ValueBool(),
		"aliases":                         aliases,
		"destinations":                    destinations,
	}

	if !plan.Expires.IsUnknown() && !plan.Expires.IsNull() {
		payload["expires"] = plan.Expires.ValueString()
	}

	respObj, err := r.client.Post(ctx, "/api/v1/acl/rule", payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating ACL Rule", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", err.Error())
		return
	}

	// Extract rule info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACLRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ACLRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/rule/%d", ruleID)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading ACL Rule", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse rule", err.Error())
		return
	}

	// Update state from response
	if id, ok := result["id"].(float64); ok {
		state.ID = types.Int64Value(int64(id))
	}
	if name, ok := result["name"].(string); ok {
		state.Name = types.StringValue(name)
	}
	if enabled, ok := result["enabled"].(bool); ok {
		state.Enabled = types.BoolValue(enabled)
	}
	if allLocations, ok := result["all_locations"].(bool); ok {
		state.AllLocations = types.BoolValue(allLocations)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACLRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ACLRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ACLRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert lists to []int64
	var locations, allowedUsers, deniedUsers, allowedGroups, deniedGroups,
		allowedNetworkDevices, deniedNetworkDevices, protocols, aliases, destinations []int64

	resp.Diagnostics.Append(plan.Locations.ElementsAs(ctx, &locations, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.AllowedUsers.ElementsAs(ctx, &allowedUsers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.DeniedUsers.ElementsAs(ctx, &deniedUsers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.AllowedGroups.ElementsAs(ctx, &allowedGroups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.DeniedGroups.ElementsAs(ctx, &deniedGroups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.AllowedNetworkDevices.ElementsAs(ctx, &allowedNetworkDevices, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.DeniedNetworkDevices.ElementsAs(ctx, &deniedNetworkDevices, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.Protocols.ElementsAs(ctx, &protocols, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.Aliases.ElementsAs(ctx, &aliases, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := plan.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/rule/%d", ruleID)

	payload := map[string]interface{}{
		"name":                            plan.Name.ValueString(),
		"enabled":                         plan.Enabled.ValueBool(),
		"all_locations":                   plan.AllLocations.ValueBool(),
		"locations":                       locations,
		"allow_all_users":                 plan.AllowAllUsers.ValueBool(),
		"deny_all_users":                  plan.DenyAllUsers.ValueBool(),
		"allow_all_groups":                plan.AllowAllGroups.ValueBool(),
		"deny_all_groups":                 plan.DenyAllGroups.ValueBool(),
		"allow_all_network_devices":       plan.AllowAllNetworkDevices.ValueBool(),
		"deny_all_network_devices":        plan.DenyAllNetworkDevices.ValueBool(),
		"allowed_users":                   allowedUsers,
		"denied_users":                    deniedUsers,
		"allowed_groups":                  allowedGroups,
		"denied_groups":                   deniedGroups,
		"allowed_network_devices":         allowedNetworkDevices,
		"denied_network_devices":          deniedNetworkDevices,
		"use_manual_destination_settings": plan.UseManualDestinationSettings.ValueBool(),
		"addresses":                       plan.Addresses.ValueString(),
		"ports":                           plan.Ports.ValueString(),
		"protocols":                       protocols,
		"any_address":                     plan.AnyAddress.ValueBool(),
		"any_port":                        plan.AnyPort.ValueBool(),
		"any_protocol":                    plan.AnyProtocol.ValueBool(),
		"aliases":                         aliases,
		"destinations":                    destinations,
	}

	if !plan.Expires.IsUnknown() && !plan.Expires.IsNull() {
		payload["expires"] = plan.Expires.ValueString()
	} else if !state.Expires.IsUnknown() && !state.Expires.IsNull() {
		payload["expires"] = state.Expires.ValueString()
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating ACL Rule", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated rule", err.Error())
		return
	}

	// Update state from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}
	if name, ok := result["name"].(string); ok {
		plan.Name = types.StringValue(name)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACLRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ACLRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/rule/%d", ruleID)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting ACL Rule", err.Error())
		return
	}
}

func (r *ACLRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
