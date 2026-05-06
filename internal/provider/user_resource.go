package provider

import (
	"context"
	"fmt"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &UserResource{}

type UserResource struct {
	client *client.Client
}

type UserResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	Username               types.String `tfsdk:"username"`
	FirstName              types.String `tfsdk:"first_name"`
	LastName               types.String `tfsdk:"last_name"`
	Email                  types.String `tfsdk:"email"`
	Phone                  types.String `tfsdk:"phone"`
	IsAdmin                types.Bool   `tfsdk:"is_admin"`
	IsActive               types.Bool   `tfsdk:"is_active"`
	Enrolled               types.Bool   `tfsdk:"enrolled"`
	MFAEnabled             types.Bool   `tfsdk:"mfa_enabled"`
	TOTPEnabled            types.Bool   `tfsdk:"totp_enabled"`
	EmailMFAEnabled        types.Bool   `tfsdk:"email_mfa_enabled"`
	MFAMethod              types.String `tfsdk:"mfa_method"`
	AuthorizedApps         types.List   `tfsdk:"authorized_apps"`
	Groups                 types.List   `tfsdk:"groups"`
	LDAPPassRequiresChange types.Bool   `tfsdk:"ldap_pass_requires_change"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User resource represents a user in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "User ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Username (unique)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"first_name": schema.StringAttribute{
				Required:    true,
				Description: "User's first name",
			},
			"last_name": schema.StringAttribute{
				Required:    true,
				Description: "User's last name",
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "User's email address",
			},
			"phone": schema.StringAttribute{
				Optional:    true,
				Description: "User's phone number",
			},
			"is_admin": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the user has admin privileges",
			},
			"is_active": schema.BoolAttribute{
				Required:    true,
				Description: "Whether the user account is active",
			},
			"enrolled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the user has completed enrollment",
			},
			"mfa_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether MFA is enabled for the user",
			},
			"totp_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether TOTP is enabled for the user",
			},
			"email_mfa_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether email MFA is enabled for the user",
			},
			"mfa_method": schema.StringAttribute{
				Computed:    true,
				Description: "Current MFA method (None, OneTimePassword, Webauthn, Email)",
			},
			"authorized_apps": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of authorized OAuth2 apps",
			},
			"groups": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Groups the user belongs to",
			},
			"ldap_pass_requires_change": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether LDAP password requires change",
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert groups to []string
	var groups []string
	resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create user payload based on OpenAPI schema
	payload := map[string]interface{}{
		"username":   plan.Username.ValueString(),
		"first_name": plan.FirstName.ValueString(),
		"last_name":  plan.LastName.ValueString(),
		"email":      plan.Email.ValueString(),
		"is_admin":   plan.IsAdmin.ValueBool(),
		"is_active":  plan.IsActive.ValueBool(),
		"groups":     groups,
	}

	// Only include phone if it's set
	if !plan.Phone.IsUnknown() && !plan.Phone.IsNull() && plan.Phone.ValueString() != "" {
		payload["phone"] = plan.Phone.ValueString()
	}

	respObj, err := r.client.Post(ctx, "/api/v1/user", payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating User", err.Error())
		return
	}

	// Get username for reading user details
	username := plan.Username.ValueString()
	readPath := fmt.Sprintf("/api/v1/user/%s", username)

	respObj, err = r.client.Get(ctx, readPath)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading User After Create", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse user after create", err.Error())
		return
	}

	// Extract user info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}

	if val, ok := result["first_name"].(string); ok {
		plan.FirstName = types.StringValue(val)
	}

	if val, ok := result["last_name"].(string); ok {
		plan.LastName = types.StringValue(val)
	}

	if val, ok := result["email"].(string); ok {
		plan.Email = types.StringValue(val)
	}

	if val, ok := result["phone"].(string); ok {
		plan.Phone = types.StringValue(val)
	}

	if val, ok := result["is_admin"].(bool); ok {
		plan.IsAdmin = types.BoolValue(val)
	}

	if val, ok := result["is_active"].(bool); ok {
		plan.IsActive = types.BoolValue(val)
	}

	if val, ok := result["enrolled"].(bool); ok {
		plan.Enrolled = types.BoolValue(val)
	}

	if val, ok := result["mfa_enabled"].(bool); ok {
		plan.MFAEnabled = types.BoolValue(val)
	}

	if val, ok := result["totp_enabled"].(bool); ok {
		plan.TOTPEnabled = types.BoolValue(val)
	}

	if val, ok := result["email_mfa_enabled"].(bool); ok {
		plan.EmailMFAEnabled = types.BoolValue(val)
	}

	if val, ok := result["mfa_method"].(string); ok {
		plan.MFAMethod = types.StringValue(val)
	}

	if val, ok := result["ldap_pass_requires_change"].(bool); ok {
		plan.LDAPPassRequiresChange = types.BoolValue(val)
	}

	if val, ok := result["authorized_apps"].([]interface{}); ok {
		apps := make([]attr.Value, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				apps = append(apps, types.StringValue(s))
			}
		}
		plan.AuthorizedApps = types.ListValueMust(types.StringType, apps)
	}

	if val, ok := result["groups"].([]interface{}); ok {
		grps := make([]attr.Value, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				grps = append(grps, types.StringValue(s))
			}
		}
		plan.Groups = types.ListValueMust(types.StringType, grps)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := state.Username.ValueString()
	path := fmt.Sprintf("/api/v1/user/%s", username)

	respObj, err := r.client.Get(ctx, path)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading User", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse user", err.Error())
		return
	}

	// Update state from response
	if id, ok := result["id"].(float64); ok {
		state.ID = types.Int64Value(int64(id))
	}

	if val, ok := result["first_name"].(string); ok {
		state.FirstName = types.StringValue(val)
	}

	if val, ok := result["last_name"].(string); ok {
		state.LastName = types.StringValue(val)
	}

	if val, ok := result["email"].(string); ok {
		state.Email = types.StringValue(val)
	}

	if val, ok := result["phone"].(string); ok {
		state.Phone = types.StringValue(val)
	}

	if val, ok := result["is_admin"].(bool); ok {
		state.IsAdmin = types.BoolValue(val)
	}

	if val, ok := result["is_active"].(bool); ok {
		state.IsActive = types.BoolValue(val)
	}

	if val, ok := result["enrolled"].(bool); ok {
		state.Enrolled = types.BoolValue(val)
	}

	if val, ok := result["mfa_enabled"].(bool); ok {
		state.MFAEnabled = types.BoolValue(val)
	}

	if val, ok := result["totp_enabled"].(bool); ok {
		state.TOTPEnabled = types.BoolValue(val)
	}

	if val, ok := result["email_mfa_enabled"].(bool); ok {
		state.EmailMFAEnabled = types.BoolValue(val)
	}

	if val, ok := result["mfa_method"].(string); ok {
		state.MFAMethod = types.StringValue(val)
	}

	if val, ok := result["ldap_pass_requires_change"].(bool); ok {
		state.LDAPPassRequiresChange = types.BoolValue(val)
	}

	if val, ok := result["authorized_apps"].([]interface{}); ok {
		apps := make([]attr.Value, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				apps = append(apps, types.StringValue(s))
			}
		}
		state.AuthorizedApps = types.ListValueMust(types.StringType, apps)
	}

	if val, ok := result["groups"].([]interface{}); ok {
		grps := make([]attr.Value, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				grps = append(grps, types.StringValue(s))
			}
		}
		state.Groups = types.ListValueMust(types.StringType, grps)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert groups to []string
	var groups []string
	resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := plan.Username.ValueString()
	path := fmt.Sprintf("/api/v1/user/%s", username)

	// Create update payload
	payload := map[string]interface{}{
		"username":   plan.Username.ValueString(),
		"first_name": plan.FirstName.ValueString(),
		"last_name":  plan.LastName.ValueString(),
		"email":      plan.Email.ValueString(),
		"is_admin":   plan.IsAdmin.ValueBool(),
		"is_active":  plan.IsActive.ValueBool(),
		"groups":     groups,
	}

	// Only include phone if it's set
	if !plan.Phone.IsUnknown() && !plan.Phone.IsNull() && plan.Phone.ValueString() != "" {
		payload["phone"] = plan.Phone.ValueString()
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating User", err.Error())
		return
	}

	// Get user details to update state
	respObj, err = r.client.Get(ctx, path)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading User After Update", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse user after update", err.Error())
		return
	}

	// Extract user info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}

	if val, ok := result["first_name"].(string); ok {
		plan.FirstName = types.StringValue(val)
	}

	if val, ok := result["last_name"].(string); ok {
		plan.LastName = types.StringValue(val)
	}

	if val, ok := result["email"].(string); ok {
		plan.Email = types.StringValue(val)
	}

	if val, ok := result["phone"].(string); ok {
		plan.Phone = types.StringValue(val)
	}

	if val, ok := result["is_admin"].(bool); ok {
		plan.IsAdmin = types.BoolValue(val)
	}

	if val, ok := result["is_active"].(bool); ok {
		plan.IsActive = types.BoolValue(val)
	}

	if val, ok := result["enrolled"].(bool); ok {
		plan.Enrolled = types.BoolValue(val)
	}

	if val, ok := result["mfa_enabled"].(bool); ok {
		plan.MFAEnabled = types.BoolValue(val)
	}

	if val, ok := result["totp_enabled"].(bool); ok {
		plan.TOTPEnabled = types.BoolValue(val)
	}

	if val, ok := result["email_mfa_enabled"].(bool); ok {
		plan.EmailMFAEnabled = types.BoolValue(val)
	}

	if val, ok := result["mfa_method"].(string); ok {
		plan.MFAMethod = types.StringValue(val)
	}

	if val, ok := result["ldap_pass_requires_change"].(bool); ok {
		plan.LDAPPassRequiresChange = types.BoolValue(val)
	}

	if val, ok := result["authorized_apps"].([]interface{}); ok {
		apps := make([]attr.Value, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				apps = append(apps, types.StringValue(s))
			}
		}
		plan.AuthorizedApps = types.ListValueMust(types.StringType, apps)
	}

	if val, ok := result["groups"].([]interface{}); ok {
		grps := make([]attr.Value, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				grps = append(grps, types.StringValue(s))
			}
		}
		plan.Groups = types.ListValueMust(types.StringType, grps)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := state.Username.ValueString()
	path := fmt.Sprintf("/api/v1/user/%s", username)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting User", err.Error())
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
