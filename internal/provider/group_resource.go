package provider

import (
	"context"
	"fmt"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &GroupResource{}

type GroupResource struct {
	client *client.Client
}

type GroupResourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Members      types.List   `tfsdk:"members"`
	IsAdmin      types.Bool   `tfsdk:"is_admin"`
	VPNLocations types.List   `tfsdk:"vpn_locations"`
}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Group resource represents a user group in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Group ID",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Group name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"members": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of member usernames",
			},
			"is_admin": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the group has admin privileges",
			},
			"vpn_locations": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "VPN locations associated with this group",
			},
		},
	}
}

func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert members to []string
	var members []string
	resp.Diagnostics.Append(plan.Members.ElementsAs(ctx, &members, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create group payload based on OpenAPI schema
	payload := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"members":  members,
		"is_admin": plan.IsAdmin.ValueBool(),
	}

	respObj, err := r.client.Post(ctx, "/api/v1/group", payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Group", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", err.Error())
		return
	}

	// Extract group info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupName := state.Name.ValueString()
	path := fmt.Sprintf("/api/v1/group/%s", groupName)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading Group", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse group", err.Error())
		return
	}

	// Update state from response - only set computed fields
	if val, ok := result["id"].(float64); ok {
		state.ID = types.Int64Value(int64(val))
	}
	if val, ok := result["is_admin"].(bool); ok {
		state.IsAdmin = types.BoolValue(val)
	}
	if val, ok := result["vpn_locations"].([]interface{}); ok {
		locations := make([]attr.Value, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				locations = append(locations, types.StringValue(s))
			}
		}
		state.VPNLocations = types.ListValueMust(types.StringType, locations)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert members to []string
	var members []string
	resp.Diagnostics.Append(plan.Members.ElementsAs(ctx, &members, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldGroupName := state.Name.ValueString()
	newGroupName := plan.Name.ValueString()

	var respObj *client.Response
	var err error

	if oldGroupName != newGroupName {
		// Rename group
		path := fmt.Sprintf("/api/v1/group/%s", oldGroupName)
		payload := map[string]interface{}{
			"name":     newGroupName,
			"members":  members,
			"is_admin": plan.IsAdmin.ValueBool(),
		}
		respObj, err = r.client.Put(ctx, path, payload)
	} else {
		// Update members and is_admin
		path := fmt.Sprintf("/api/v1/group/%s", newGroupName)
		payload := map[string]interface{}{
			"members":  members,
			"is_admin": plan.IsAdmin.ValueBool(),
		}
		respObj, err = r.client.Put(ctx, path, payload)
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Group", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupName := state.Name.ValueString()
	path := fmt.Sprintf("/api/v1/group/%s", groupName)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Group", err.Error())
		return
	}
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
