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

var _ resource.Resource = &DeviceResource{}

type DeviceResource struct {
	client *client.Client
}

type DeviceResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	UserID          types.Int64  `tfsdk:"user_id"`
	WireguardPubkey types.String `tfsdk:"wireguard_pubkey"`
	Created         types.String `tfsdk:"created"`
	UserIDValue     types.Int64  `tfsdk:"-"`
}

func NewDeviceResource() resource.Resource {
	return &DeviceResource{}
}

func (r *DeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *DeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Device resource represents a WireGuard device in Defguard.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Device ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Device name",
			},
			"user_id": schema.Int64Attribute{
				Required:    true,
				Description: "User ID who owns this device",
			},
			"wireguard_pubkey": schema.StringAttribute{
				Required:    true,
				Description: "WireGuard public key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
			},
		},
	}
}

func (r *DeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create device payload based on OpenAPI schema
	payload := map[string]interface{}{
		"name":             plan.Name.ValueString(),
		"wireguard_pubkey": plan.WireguardPubkey.ValueString(),
	}

	// Add user_id to the path
	username := fmt.Sprintf("user_%d", plan.UserID.ValueInt64())
	devicePath := fmt.Sprintf("/api/v1/device/%s", username)

	respObj, err := r.client.Post(ctx, devicePath, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Device", err.Error())
		return
	}

	// Parse response - the API returns AddDeviceResult with device inside
	// For now, we'll just use the device object
	var result map[string]interface{}
	if err := respObj.Unmarshal(&result); err != nil {
		resp.Diagnostics.AddError("Failed to parse response", err.Error())
		return
	}

	// Extract device from result
	deviceMap, ok := result["device"].(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError("Failed to extract device from response", "device field not found or invalid type")
		return
	}

	// Set values from response
	if id, ok := deviceMap["id"].(float64); ok {
		plan.ID = types.Int64Value(int64(id))
	}
	if created, ok := deviceMap["created"].(string); ok {
		plan.Created = types.StringValue(created)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := state.ID.ValueInt64()
	devicePath := fmt.Sprintf("/api/v1/device/%d", deviceID)

	respObj, err := r.client.Get(ctx, devicePath)
	if err != nil {
		if respObj != nil && respObj.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error Reading Device", err.Error())
		return
	}

	var device DeviceResourceModel
	if err := respObj.Unmarshal(&device); err != nil {
		resp.Diagnostics.AddError("Failed to parse device", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &device)...)
}

func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := plan.ID.ValueInt64()
	devicePath := fmt.Sprintf("/api/v1/device/%d", deviceID)

	// Create update payload
	payload := map[string]interface{}{
		"name":             plan.Name.ValueString(),
		"wireguard_pubkey": plan.WireguardPubkey.ValueString(),
	}

	respObj, err := r.client.Put(ctx, devicePath, payload)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Device", err.Error())
		return
	}

	var device DeviceResourceModel
	if err := respObj.Unmarshal(&device); err != nil {
		resp.Diagnostics.AddError("Failed to parse updated device", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &device)...)
}

func (r *DeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceID := state.ID.ValueInt64()
	path := fmt.Sprintf("/api/v1/device/%d", deviceID)

	_, err := r.client.Delete(ctx, path, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Device", err.Error())
		return
	}
}

func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
