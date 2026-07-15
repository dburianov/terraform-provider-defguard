package provider

import (
	"context"
	"fmt"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Group name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"members": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of member usernames",
			},
			"is_admin": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the group has admin privileges",
			},
			"vpn_locations": schema.ListAttribute{
				Computed:      true,
				ElementType:   types.StringType,
				Description:   "VPN locations associated with this group",
				PlanModifiers: []planmodifier.List{},
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

	groupName := plan.Name.ValueString()

	// Check if group with same name already exists
	existingGroup := r.findGroupByName(ctx, groupName)
	if existingGroup != nil {
		resp.Diagnostics.AddError(
			"Group Already Exists",
			fmt.Sprintf("A group with name '%s' already exists (ID: %d). Use import or update the existing group instead of creating a new one.", groupName, existingGroup["id"]),
		)
		return
	}

	// Build group payload based on OpenAPI EditGroupInfo schema
	payload := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"is_admin": plan.IsAdmin.ValueBool(),
	}

	// Only include members in payload if provided (not null/unknown)
	if !plan.Members.IsNull() && !plan.Members.IsUnknown() {
		var members []string
		resp.Diagnostics.Append(plan.Members.ElementsAs(ctx, &members, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		payload["members"] = members
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

	// The POST /api/v1/group endpoint returns EditGroupInfo without ID.
	// We need to fetch all groups and filter by name to get the full info including ID.
	getResp, err := r.client.Get(ctx, "/api/v1/group-info")
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Groups After Creation", err.Error())
		return
	}

	var groups []map[string]interface{}
	if err := getResp.Unmarshal(&groups); err != nil {
		resp.Diagnostics.AddError("Failed to parse group info after creation", err.Error())
		return
	}

	// Find the created group by name (should always succeed now)
	var groupInfo map[string]interface{}
	for _, g := range groups {
		if name, ok := g["name"].(string); ok && name == groupName {
			groupInfo = g
			break
		}
	}

	if groupInfo == nil {
		resp.Diagnostics.AddError("API Error Reading Group After Creation", fmt.Sprintf("Group '%s' not found after creation - this should never happen", groupName))
		return
	}

	// Extract group info from filtered response which includes the ID
	if id, ok := groupInfo["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	} else if idStr, ok := groupInfo["id"].(string); ok {
		var intID int64
		fmt.Sscanf(idStr, "%d", &intID)
		plan.ID = types.Int64Value(intID)
	}
	if name, ok := groupInfo["name"].(string); ok {
		plan.Name = types.StringValue(name)
	}
	if isAdmin, ok := groupInfo["is_admin"].(bool); ok {
		plan.IsAdmin = types.BoolValue(isAdmin)
	}
	if locationsRaw, ok := groupInfo["vpn_locations"].([]interface{}); ok {
		var vpnLocations []attr.Value
		for _, loc := range locationsRaw {
			if locStr, ok := loc.(string); ok {
				vpnLocations = append(vpnLocations, types.StringValue(locStr))
			}
		}
		plan.VPNLocations, _ = types.ListValue(types.StringType, vpnLocations)
	} else if plan.VPNLocations.IsNull() || plan.VPNLocations.IsUnknown() {
		// If vpn_locations not in response, set to empty list
		plan.VPNLocations, _ = types.ListValue(types.StringType, []attr.Value{})
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

	respObj, err := r.client.Get(ctx, "/api/v1/group-info")
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Group", err.Error())
		return
	}

	var groups []map[string]interface{}
	if err := respObj.Unmarshal(&groups); err != nil {
		resp.Diagnostics.AddError("Failed to parse group info", err.Error())
		return
	}

	// Find the group by name in the list of all groups
	var groupInfo map[string]interface{}
	for _, g := range groups {
		if name, ok := g["name"].(string); ok && name == groupName {
			groupInfo = g
			break
		}
	}

	if groupInfo == nil {
		resp.Diagnostics.AddError("Group Not Found", fmt.Sprintf("Group '%s' not found", groupName))
		return
	}

	// Update state from response
	if id, ok := groupInfo["id"].(float64); ok {
		state.ID = types.Int64Value(int64(id))
	}
	if name, ok := groupInfo["name"].(string); ok {
		state.Name = types.StringValue(name)
	}
	if isAdmin, ok := groupInfo["is_admin"].(bool); ok {
		state.IsAdmin = types.BoolValue(isAdmin)
	}
	if locationsRaw, ok := groupInfo["vpn_locations"].([]interface{}); ok {
		var vpnLocations []attr.Value
		for _, loc := range locationsRaw {
			if locStr, ok := loc.(string); ok {
				vpnLocations = append(vpnLocations, types.StringValue(locStr))
			}
		}
		state.VPNLocations, _ = types.ListValue(types.StringType, vpnLocations)
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

	groupName := plan.Name.ValueString()
	path := fmt.Sprintf("/api/v1/group/%s", groupName)

	payload := map[string]interface{}{
		"name":     plan.Name.ValueString(),
		"is_admin": plan.IsAdmin.ValueBool(),
	}

	// Only include members in payload if provided (not null/unknown)
	if !plan.Members.IsNull() && !plan.Members.IsUnknown() {
		var members []string
		resp.Diagnostics.Append(plan.Members.ElementsAs(ctx, &members, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		payload["members"] = members
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Group", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated group", err.Error())
		return
	}

	// Update state from response
	if name, ok := result["name"].(string); ok {
		plan.Name = types.StringValue(name)
	}
	if isAdmin, ok := result["is_admin"].(bool); ok {
		plan.IsAdmin = types.BoolValue(isAdmin)
	}
	if locationsRaw, ok := result["vpn_locations"].([]interface{}); ok {
		var vpnLocations []attr.Value
		for _, loc := range locationsRaw {
			if locStr, ok := loc.(string); ok {
				vpnLocations = append(vpnLocations, types.StringValue(locStr))
			}
		}
		plan.VPNLocations, _ = types.ListValue(types.StringType, vpnLocations)
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
	// Import can be done using either ID or name
	importID := req.ID

	// First try to find if it's a numeric ID by looking up all groups
	var groups []map[string]interface{}
	getResp, err := r.client.Get(ctx, "/api/v1/group-info")
	if err == nil {
		if getResp.Unmarshal(&groups) == nil {
			for _, g := range groups {
				// Check if this is an ID match
				if idVal, ok := g["id"].(float64); ok && fmt.Sprintf("%d", int64(idVal)) == importID {
					// Found by ID - set the name
					if name, ok := g["name"].(string); ok {
						resp.Diagnostics.Append(resp.State.Set(ctx, GroupResourceModel{
							ID:           types.Int64Value(int64(idVal)),
							Name:         types.StringValue(name),
							Members:      types.ListNull(types.StringType),
							IsAdmin:      types.BoolUnknown(),
							VPNLocations: types.ListNull(types.StringType),
						})...)
					}
					return
				}
			}
		}
	}

	// If not found by ID, try to find by name
	for _, g := range groups {
		if name, ok := g["name"].(string); ok && name == importID {
			if idVal, ok := g["id"].(float64); ok {
				resp.Diagnostics.Append(resp.State.Set(ctx, GroupResourceModel{
					ID:           types.Int64Value(int64(idVal)),
					Name:         types.StringValue(name),
					Members:      types.ListNull(types.StringType),
					IsAdmin:      types.BoolUnknown(),
					VPNLocations: types.ListNull(types.StringType),
				})...)
			}
			return
		}
	}

	resp.Diagnostics.AddError(
		"Cannot find group for import",
		fmt.Sprintf("Group with ID/name '%s' not found", importID),
	)
}

// findGroupByName searches for a group by name in the list of all groups
// Returns nil if not found, or the group info map if found
func (r *GroupResource) findGroupByName(ctx context.Context, name string) map[string]interface{} {
	respObj, err := r.client.Get(ctx, "/api/v1/group-info")
	if err != nil {
		return nil
	}

	var groups []map[string]interface{}
	if err := respObj.Unmarshal(&groups); err != nil {
		return nil
	}

	for _, g := range groups {
		if groupName, ok := g["name"].(string); ok && groupName == name {
			return g
		}
	}

	return nil
}

// RemoveUserFromGroupRequest represents the request body for removing a user from a group
type RemoveUserFromGroupRequest struct {
	Username string `json:"username"`
}

// removeUserFromGroup removes a user from a group by group ID and username
func (r *GroupResource) removeUserFromGroup(ctx context.Context, groupID int64, username string) error {
	path := fmt.Sprintf("/api/v1/group/%d/user/%s", groupID, username)
	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to remove user %s from group %d: %w", username, groupID, err)
	}
	return nil
}

// AddGroupMemberRequest represents the request body for adding a user to a group
type AddGroupMemberRequest struct {
	Username string `json:"username"`
}

// addGroupMember adds a user to a group by group ID and username
func (r *GroupResource) addGroupMember(ctx context.Context, groupID int64, username string) error {
	path := fmt.Sprintf("/api/v1/group/%d/user/%s", groupID, username)
	payload := map[string]interface{}{
		"username": username,
	}
	respObj, err := r.client.Post(ctx, path, payload)
	if err != nil {
		return fmt.Errorf("failed to add user %s to group %d: %w", username, groupID, err)
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
