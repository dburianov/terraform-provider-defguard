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

var _ resource.Resource = &SNATBindingResource{}

type SNATBindingResource struct {
	client *client.Client
}

type SNATBindingResourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	UserID     types.Int64  `tfsdk:"user_id"`
	LocationID types.Int64  `tfsdk:"location_id"`
	PublicIP   types.String `tfsdk:"public_ip"`
}

func NewSNATBindingResource() resource.Resource {
	return &SNATBindingResource{}
}

func (r *SNATBindingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snat_binding"
}

func (r *SNATBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SNAT binding resource represents a SNAT (Source Network Address Translation) binding in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "SNAT binding ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Required:    true,
				Description: "User ID to bind to the public IP",
			},
			"location_id": schema.Int64Attribute{
				Required:    true,
				Description: "WireGuard location ID",
			},
			"public_ip": schema.StringAttribute{
				Required:    true,
				Description: "Public IP address for SNAT",
			},
		},
	}
}

func (r *SNATBindingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SNATBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SNATBindingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locationID := plan.LocationID.ValueInt64()
	path := fmt.Sprintf("/api/v1/network/%d/snat", locationID)

	payload := map[string]interface{}{
		"user_id":   plan.UserID.ValueInt64(),
		"public_ip": plan.PublicIP.ValueString(),
	}

	respObj, err := r.client.Post(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating SNAT Binding", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", err.Error())
		return
	}

	// Extract SNAT binding info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SNATBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SNATBindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locationID := state.LocationID.ValueInt64()
	userID := state.UserID.ValueInt64()
	path := fmt.Sprintf("/api/v1/network/%d/snat/%d", locationID, userID)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading SNAT Binding", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse SNAT binding", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SNATBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SNATBindingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SNATBindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locationID := plan.LocationID.ValueInt64()
	userID := plan.UserID.ValueInt64()
	path := fmt.Sprintf("/api/v1/network/%d/snat/%d", locationID, userID)

	payload := map[string]interface{}{
		"public_ip": plan.PublicIP.ValueString(),
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating SNAT Binding", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated SNAT binding", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SNATBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SNATBindingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locationID := state.LocationID.ValueInt64()
	userID := state.UserID.ValueInt64()
	path := fmt.Sprintf("/api/v1/network/%d/snat/%d", locationID, userID)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting SNAT Binding", err.Error())
		return
	}
}

func (r *SNATBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
