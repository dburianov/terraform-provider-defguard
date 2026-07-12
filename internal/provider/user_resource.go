package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/dburianov/terraform-provider-defguard/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &UserResource{}

// PasswordValidator validates password strength according to DefGuard requirements
type PasswordValidator struct{}

func (v PasswordValidator) Description(ctx context.Context) string {
	return "Password must be at least 10 characters and contain lowercase letters, uppercase letters, numbers, and special symbols"
}

func (v PasswordValidator) MarkdownDescription(ctx context.Context) string {
	return "Password must be at least 10 characters and contain:\n- Lowercase letters\n- Uppercase letters\n- Numbers\n- Special symbols"
}

func (v PasswordValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	password := req.ConfigValue.ValueString()

	// Check minimum length
	if len(password) < 10 {
		resp.Diagnostics.AddError(
			"Invalid Password",
			"Password must be at least 10 characters long",
		)
		return
	}

	// Check for lowercase letter
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLowercase {
		resp.Diagnostics.AddError(
			"Invalid Password",
			"Password must contain at least one lowercase letter",
		)
		return
	}

	// Check for uppercase letter
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		resp.Diagnostics.AddError(
			"Invalid Password",
			"Password must contain at least one uppercase letter",
		)
		return
	}

	// Check for number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		resp.Diagnostics.AddError(
			"Invalid Password",
			"Password must contain at least one number",
		)
		return
	}

	// Check for special symbol (non-alphanumeric)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	if !hasSpecial {
		resp.Diagnostics.AddError(
			"Invalid Password",
			"Password must contain at least one special symbol",
		)
		return
	}
}

type UserResource struct {
	client *client.Client
}

type UserResourceModel struct {
	ID                     types.Int64  `tfsdk:"id"`
	Username               types.String `tfsdk:"username"`
	FirstName              types.String `tfsdk:"first_name"`
	LastName               types.String `tfsdk:"last_name"`
	Name                   types.String `tfsdk:"name"`
	Email                  types.String `tfsdk:"email"`
	Phone                  types.String `tfsdk:"phone"`
	Password               types.String `tfsdk:"password"`
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
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "User's full name (first + last)",
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "User's email address",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "User's password (10+ chars, lowercase, uppercase, numbers, special symbols)",
				Validators: []validator.String{
					PasswordValidator{},
				},
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
				Optional:    true,
				Computed:    true,
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

	// Create user payload based on OpenAPI AddUserData schema
	payload := map[string]interface{}{
		"username":   plan.Username.ValueString(),
		"first_name": plan.FirstName.ValueString(),
		"last_name":  plan.LastName.ValueString(),
		"email":      plan.Email.ValueString(),
	}

	// Only include optional fields if provided (not null/unknown)
	if !plan.Phone.IsUnknown() && !plan.Phone.IsNull() {
		payload["phone"] = plan.Phone.ValueString()
	}
	payload["is_admin"] = plan.IsAdmin.ValueBool()
	payload["is_active"] = plan.IsActive.ValueBool()

	// Debug: log the payload being sent
	// Remove this debug logging before production use
	payloadJSON, _ := json.Marshal(payload)
	fmt.Fprintf(os.Stderr, "DEBUG User create payload JSON: %s\n", string(payloadJSON))

	// Only include groups in payload if provided (not null/unknown)
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		var groups []string
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		payload["groups"] = groups
	}

	respObj, err := r.client.Post(ctx, "/api/v1/user", payload)

	// Debug: log the response for inspection
	// Remove this debug logging before production use
	if respObj != nil {
		fmt.Printf("DEBUG Response status: %d, body len: %d, body: %s\n", respObj.StatusCode, len(respObj.Body), string(respObj.Body))
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating User", err.Error())
		return
	}

	// Check for HTTP error status codes (4xx, 5xx)
	if respObj.Err != nil {
		resp.Diagnostics.AddError("API Error Creating User", fmt.Sprintf("HTTP %d: %s", respObj.StatusCode, respObj.Err.Error()))
		return
	}

	var result map[string]interface{}
	fmt.Printf("DEBUG respObj.Err: %v\n", respObj.Err)
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", fmt.Sprintf("err=%v, body len=%d, body=%s", err, len(respObj.Body), string(respObj.Body)))
		return
	}

	// Extract user info from response
	if id, ok := result["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}
	if name, ok := result["name"].(string); ok {
		plan.Name = types.StringValue(name)
	}

	// Set default values for computed fields that may not be in the response
	plan.Enrolled = types.BoolValue(false)        // Default to false if not returned
	plan.MFAEnabled = types.BoolValue(false)      // Default to false if not returned
	plan.TOTPEnabled = types.BoolValue(false)     // Default to false if not returned
	plan.EmailMFAEnabled = types.BoolValue(false) // Default to false if not returned
	plan.LDAPPassRequiresChange = types.BoolValue(false)

	// Try to read back all computed fields from API for accurate values
	getPath := fmt.Sprintf("/api/v1/user/%s", plan.Username.ValueString())
	respObjGet, err := r.client.Get(ctx, getPath)
	if respObjGet != nil {
		fmt.Printf("DEBUG GET response status: %d\n", respObjGet.StatusCode)
	}
	if err == nil {
		var readResult map[string]interface{}
		if readErr := respObjGet.Unmarshal(&readResult); readErr == nil {
			if enrolled, ok := readResult["enrolled"].(bool); ok {
				plan.Enrolled = types.BoolValue(enrolled)
			}
			if mfaEnabled, ok := readResult["mfa_enabled"].(bool); ok {
				plan.MFAEnabled = types.BoolValue(mfaEnabled)
			}
			if totpEnabled, ok := readResult["totp_enabled"].(bool); ok {
				plan.TOTPEnabled = types.BoolValue(totpEnabled)
			}
			if emailMFAEnabled, ok := readResult["email_mfa_enabled"].(bool); ok {
				plan.EmailMFAEnabled = types.BoolValue(emailMFAEnabled)
			}
			if mfaMethod, ok := readResult["mfa_method"].(string); ok {
				plan.MFAMethod = types.StringValue(mfaMethod)
			} else {
				plan.MFAMethod = types.StringValue("None") // Default value
			}
			if ldapPassRequiresChange, ok := readResult["ldap_pass_requires_change"].(bool); ok {
				plan.LDAPPassRequiresChange = types.BoolValue(ldapPassRequiresChange)
			}
			if appsRaw, ok := readResult["authorized_apps"].([]interface{}); ok {
				var authorizedApps []attr.Value
				for _, app := range appsRaw {
					if appStr, ok := app.(string); ok {
						authorizedApps = append(authorizedApps, types.StringValue(appStr))
					}
				}
				if len(authorizedApps) > 0 {
					plan.AuthorizedApps, _ = types.ListValue(types.StringType, authorizedApps)
				} else {
					plan.AuthorizedApps, _ = types.ListValue(types.StringType, []attr.Value{})
				}
			} else {
				plan.AuthorizedApps, _ = types.ListValue(types.StringType, []attr.Value{})
			}

			// Handle groups - set default empty list if not present
			if groupsRaw, ok := readResult["groups"].([]interface{}); ok {
				var groups []attr.Value
				for _, g := range groupsRaw {
					if groupStr, ok := g.(string); ok {
						groups = append(groups, types.StringValue(groupStr))
					}
				}
				if len(groups) > 0 {
					plan.Groups, _ = types.ListValue(types.StringType, groups)
				} else {
					plan.Groups, _ = types.ListValue(types.StringType, []attr.Value{})
				}
			} else {
				plan.Groups, _ = types.ListValue(types.StringType, []attr.Value{})
			}
		}
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

	// Handle error from the request or response
	if err != nil || (respObj != nil && respObj.Err != nil) {
		// Check if this is a 404 - user was deleted outside Terraform
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		// Extract error message for better diagnostics
		errorMsg := err.Error()
		if respObj != nil && respObj.Err != nil {
			errorMsg = respObj.Err.Error()
		}
		resp.Diagnostics.AddError("API Error Reading User", errorMsg)
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
	if name, ok := result["name"].(string); ok {
		state.Name = types.StringValue(name)
	}

	// Read computed fields
	if enrolled, ok := result["enrolled"].(bool); ok {
		state.Enrolled = types.BoolValue(enrolled)
	}
	if mfaEnabled, ok := result["mfa_enabled"].(bool); ok {
		state.MFAEnabled = types.BoolValue(mfaEnabled)
	}
	if totpEnabled, ok := result["totp_enabled"].(bool); ok {
		state.TOTPEnabled = types.BoolValue(totpEnabled)
	}
	if emailMFAEnabled, ok := result["email_mfa_enabled"].(bool); ok {
		state.EmailMFAEnabled = types.BoolValue(emailMFAEnabled)
	}
	if mfaMethod, ok := result["mfa_method"].(string); ok {
		state.MFAMethod = types.StringValue(mfaMethod)
	}
	if ldapPassRequiresChange, ok := result["ldap_pass_requires_change"].(bool); ok {
		state.LDAPPassRequiresChange = types.BoolValue(ldapPassRequiresChange)
	}
	if appsRaw, ok := result["authorized_apps"].([]interface{}); ok {
		var authorizedApps []attr.Value
		for _, app := range appsRaw {
			if appStr, ok := app.(string); ok {
				authorizedApps = append(authorizedApps, types.StringValue(appStr))
			}
		}
		state.AuthorizedApps, _ = types.ListValue(types.StringType, authorizedApps)
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

	username := plan.Username.ValueString()
	path := fmt.Sprintf("/api/v1/user/%s", username)

	// Create update payload based on OpenAPI schema
	// Debug: log the payload being sent
	// Remove this debug logging before production use
	fmt.Fprintf(os.Stderr, "DEBUG plan.Username: '%s'\n", plan.Username.ValueString())
	fmt.Fprintf(os.Stderr, "DEBUG plan.FirstName: '%s'\n", plan.FirstName.ValueString())
	fmt.Fprintf(os.Stderr, "DEBUG plan.LastName: '%s'\n", plan.LastName.ValueString())
	fmt.Fprintf(os.Stderr, "DEBUG plan.Email: '%s'\n", plan.Email.ValueString())
	fmt.Fprintf(os.Stderr, "DEBUG plan.IsAdmin: %v\n", plan.IsAdmin.ValueBool())
	fmt.Fprintf(os.Stderr, "DEBUG plan.IsActive: %v\n", plan.IsActive.ValueBool())

	payload := map[string]interface{}{
		"username":   plan.Username.ValueString(),
		"first_name": plan.FirstName.ValueString(),
		"last_name":  plan.LastName.ValueString(),
		"email":      plan.Email.ValueString(),
	}

	// Only include optional fields if provided (not null/unknown)
	if !plan.Phone.IsUnknown() && !plan.Phone.IsNull() {
		payload["phone"] = plan.Phone.ValueString()
	}
	payload["is_admin"] = plan.IsAdmin.ValueBool()
	payload["is_active"] = plan.IsActive.ValueBool()

	fmt.Fprintf(os.Stderr, "DEBUG payload map: %+v\n", payload)

	// Only include groups in payload if provided (not null/unknown)
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		var groups []string
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		payload["groups"] = groups
	}

	respObj, err := r.client.Put(ctx, path, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating User", err.Error())
		return
	}

	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated user", err.Error())
		return
	}

	// Update state from response
	if name, ok := result["name"].(string); ok {
		plan.Name = types.StringValue(name)
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
