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

var _ resource.Resource = &ACLDestinationResource{}

type ACLDestinationResource struct {
	client *client.Client
}

type ACLDestinationResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Kind        types.String `tfsdk:"kind"`
	State       types.String `tfsdk:"state"`
	Addresses   types.String `tfsdk:"addresses"`
	Ports       types.String `tfsdk:"ports"`
	Protocols   types.List   `tfsdk:"protocols"`
	Rules       types.List   `tfsdk:"rules"`
	ParentID    types.Int64  `tfsdk:"parent_id"`
	AnyAddress  types.Bool   `tfsdk:"any_address"`
	AnyPort     types.Bool   `tfsdk:"any_port"`
	AnyProtocol types.Bool   `tfsdk:"any_protocol"`
}

func NewACLDestinationResource() resource.Resource {
	return &ACLDestinationResource{}
}

func (r *ACLDestinationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl_destination"
}

func (r *ACLDestinationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "ACL destination resource represents an ACL destination in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ACL destination ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Destination name",
			},
			"kind": schema.StringAttribute{
				Computed:    true,
				Description: "Destination kind (always Destination)",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Destination state (Applied, Modified)",
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
				Description: "ACL rule IDs that use this destination",
			},
			"parent_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Parent alias ID (for nested aliases)",
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
		},
	}
}

func (r *ACLDestinationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ACLDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ACLDestinationResourceModel
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

	// Create destination payload based on OpenAPI EditAclDestination schema
	payload := map[string]interface{}{
		"name":         plan.Name.ValueString(),
		"addresses":    plan.Addresses.ValueString(),
		"ports":        plan.Ports.ValueString(),
		"protocols":    protocols,
		"any_address":  plan.AnyAddress.ValueBool(),
		"any_port":     plan.AnyPort.ValueBool(),
		"any_protocol": plan.AnyProtocol.ValueBool(),
	}

	respObj, err := r.client.Post(ctx, "/api/v1/acl/destination", payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating ACL Destination", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", err.Error())
		return
	}

	// Extract destination info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}
	if state, ok := result["state"].(string); ok {
		plan.State = types.StringValue(state)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACLDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ACLDestinationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/destination/%d", destinationID)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading ACL Destination", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse destination", err.Error())
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
	if anyAddress, ok := result["any_address"].(bool); ok {
		state.AnyAddress = types.BoolValue(anyAddress)
	}
	if anyPort, ok := result["any_port"].(bool); ok {
		state.AnyPort = types.BoolValue(anyPort)
	}
	if anyProtocol, ok := result["any_protocol"].(bool); ok {
		state.AnyProtocol = types.BoolValue(anyProtocol)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACLDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ACLDestinationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ACLDestinationResourceModel
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

	destinationID := plan.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/destination/%d", destinationID)

	payload := map[string]interface{}{
		"name":         plan.Name.ValueString(),
		"addresses":    plan.Addresses.ValueString(),
		"ports":        plan.Ports.ValueString(),
		"protocols":    protocols,
		"any_address":  plan.AnyAddress.ValueBool(),
		"any_port":     plan.AnyPort.ValueBool(),
		"any_protocol": plan.AnyProtocol.ValueBool(),
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating ACL Destination", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated destination", err.Error())
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

func (r *ACLDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ACLDestinationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	destinationID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/acl/destination/%d", destinationID)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting ACL Destination", err.Error())
		return
	}
}

func (r *ACLDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
