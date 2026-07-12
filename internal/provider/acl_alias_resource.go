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

var _ resource.Resource = &ACLAliasResource{}

type ACLAliasResource struct {
	client *client.Client
}

type ACLAliasResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Kind      types.String `tfsdk:"kind"`
	State     types.String `tfsdk:"state"`
	Addresses types.String `tfsdk:"addresses"`
	Ports     types.String `tfsdk:"ports"`
	Protocols types.List   `tfsdk:"protocols"`
	Rules     types.List   `tfsdk:"rules"`
	ParentID  types.Int64  `tfsdk:"parent_id"`
}

func NewACLAliasResource() resource.Resource {
	return &ACLAliasResource{}
}

func (r *ACLAliasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl_alias"
}

func (r *ACLAliasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "ACL alias resource represents an ACL alias in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ACL alias ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Alias name",
			},
			"kind": schema.StringAttribute{
				Required:    true,
				Description: "Alias kind (Destination or Component)",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Alias state (Applied, Modified)",
			},
			"addresses": schema.StringAttribute{
				Optional:    true,
				Description: "Addresses (comma-separated list of IP addresses or CIDR ranges)",
			},
			"ports": schema.StringAttribute{
				Optional:    true,
				Description: "Ports (comma-separated list of port numbers or ranges)",
			},
			"protocols": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Protocol IDs",
			},
			"rules": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "ACL rule IDs that use this alias",
			},
			"parent_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Parent alias ID (for nested aliases)",
			},
		},
	}
}

func (r *ACLAliasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ACLAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ACLAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert protocols to []int64
	var protocols []int64
	resp.Diagnostics.Append(plan.Protocols.ElementsAs(ctx, &protocols, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert rules to []int64
	var rules []int64
	resp.Diagnostics.Append(plan.Rules.ElementsAs(ctx, &rules, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create alias payload based on OpenAPI EditAclAlias schema
	payload := map[string]interface{}{
		"name":      plan.Name.ValueString(),
		"addresses": plan.Addresses.ValueString(),
		"ports":     plan.Ports.ValueString(),
		"protocols": protocols,
	}

	respObj, err := r.client.Post(ctx, "/api/v1/acl/alias", payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating ACL Alias", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", err.Error())
		return
	}

	// Extract alias info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}
	if state, ok := result["state"].(string); ok {
		plan.State = types.StringValue(state)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACLAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ACLAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aliasID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/alias/%d", aliasID)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading ACL Alias", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse alias", err.Error())
		return
	}

	// Update state from response
	if id, ok := result["id"].(float64); ok {
		state.ID = types.Int64Value(int64(id))
	}
	if name, ok := result["name"].(string); ok {
		state.Name = types.StringValue(name)
	}
	if kind, ok := result["kind"].(string); ok {
		state.Kind = types.StringValue(kind)
	}
	if addresses, ok := result["addresses"].(string); ok {
		state.Addresses = types.StringValue(addresses)
	}
	if ports, ok := result["ports"].(string); ok {
		state.Ports = types.StringValue(ports)
	}
	if stateStr, ok := result["state"].(string); ok {
		state.State = types.StringValue(stateStr)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACLAliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ACLAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ACLAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert protocols to []int64
	var protocols []int64
	resp.Diagnostics.Append(plan.Protocols.ElementsAs(ctx, &protocols, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aliasID := plan.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/alias/%d", aliasID)

	payload := map[string]interface{}{
		"name":      plan.Name.ValueString(),
		"addresses": plan.Addresses.ValueString(),
		"ports":     plan.Ports.ValueString(),
		"protocols": protocols,
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating ACL Alias", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated alias", err.Error())
		return
	}

	// Update state from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}
	if name, ok := result["name"].(string); ok {
		plan.Name = types.StringValue(name)
	}
	if kind, ok := result["kind"].(string); ok {
		plan.Kind = types.StringValue(kind)
	}
	if addresses, ok := result["addresses"].(string); ok {
		plan.Addresses = types.StringValue(addresses)
	}
	if ports, ok := result["ports"].(string); ok {
		plan.Ports = types.StringValue(ports)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACLAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ACLAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aliasID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/alias/%d", aliasID)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting ACL Alias", err.Error())
		return
	}
}

func (r *ACLAliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
