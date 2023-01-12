package provision6connect

import (
	"context"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dnszoneResource{}
	_ resource.ResourceWithConfigure   = &dnszoneResource{}
	_ resource.ResourceWithImportState = &dnszoneResource{}
)

// NewDNSzoneResource is a helper function to simplify the provider implementation.
func NewDNSzoneResource() resource.Resource {
	return &dnszoneResource{}
}

// resourcesModel maps resources schema data.
type dnszoneModel struct {
	ID          types.String `tfsdk:"id"`
	ParentID    types.String `tfsdk:"parent_id"`
	GroupID     types.String `tfsdk:"group_id"`
	Name        types.String `tfsdk:"name"`
	Modified    types.String `tfsdk:"modified"`
	Status      types.String `tfsdk:"status"`
	ZoneType    types.String `tfsdk:"zone_type"`
	ZoneExpire  types.Int64  `tfsdk:"zone_expire"`
	ZoneHost    types.String `tfsdk:"zone_host"`
	ZoneMail    types.String `tfsdk:"zone_mail"`
	ZoneMinimum types.Int64  `tfsdk:"zone_minimum"`
	ZoneRefresh types.Int64  `tfsdk:"zone_refresh"`
	ZoneRetry   types.Int64  `tfsdk:"zone_retry"`
	ZoneSerial  types.Int64  `tfsdk:"zone_serial"`
	ZoneTTL     types.Int64  `tfsdk:"zone_ttl"`
}

// dnszoneResource is the resource implementation.
type dnszoneResource struct {
	client *provisionclient.Client
}

// Configure adds the provider configured client to the resource.
func (r *dnszoneResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*provisionclient.Client)
}

// Metadata returns the resource type name.
func (r *dnszoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnszone"
}

// Schema defines the schema for the resource.
func (r *dnszoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the DNS Zone.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "DNS Zone Name, it must be FQDN, for Example: example.com.",
				Required:    true,
			},
			"parent_id": schema.StringAttribute{
				Description: "Parent ID for the Zone mainly because of permissions, if it is not set ProVision will set TLR by default.",
				Optional:    true,
				Computed:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "Group Identifier for the Zone",
				Optional:    true,
				Computed:    true,
			},
			"modified": schema.StringAttribute{
				Description: "Date and Time of the last modification",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Current status set by ProVision of the DNS Zone",
				Computed:    true,
			},
			"zone_type": schema.StringAttribute{
				Description: "Type for the current zone, it can be forward or reverse",
				Computed:    true,
				Optional:    true,
			},
			"zone_expire": schema.Int64Attribute{
				Description: "DNS Zone Expire Time",
				Computed:    true,
				Optional:    true,
			},
			"zone_host": schema.StringAttribute{
				Description: "DNS Zone Host in FQDN format",
				Computed:    true,
				Optional:    true,
			},
			"zone_mail": schema.StringAttribute{
				Description: "DNS Zone Mail in FQDN format",
				Computed:    true,
				Optional:    true,
			},
			"zone_minimum": schema.Int64Attribute{
				Description: "DNS Zone Minimum Time",
				Computed:    true,
				Optional:    true,
			},
			"zone_refresh": schema.Int64Attribute{
				Description: "DNS Zone Refresh Time",
				Computed:    true,
				Optional:    true,
			},
			"zone_retry": schema.Int64Attribute{
				Description: "DNS Zone Retry Time",
				Computed:    true,
				Optional:    true,
			},
			"zone_serial": schema.Int64Attribute{
				Description: "DNS Zone Serial",
				Computed:    true,
				Optional:    true,
			},
			"zone_ttl": schema.Int64Attribute{
				Description: "DNS Zone TTL",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Create a new resource
func (r *dnszoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dnszoneModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newZone := provisionclient.DNSZone{
		Name: plan.Name.ValueString(),
	}

	if !plan.GroupID.IsNull() {
		newZone.GroupID = provisionclient.PVID(plan.GroupID.ValueString())
	}

	if !plan.ZoneHost.IsNull() {
		newZone.ZoneHost = plan.ZoneHost.ValueString()
	}
	if !plan.ZoneMail.IsNull() {
		newZone.ZoneMail = plan.ZoneMail.ValueString()
	}
	if !plan.ZoneType.IsNull() {
		newZone.ZoneType = plan.ZoneType.ValueString()
	}
	if !plan.ZoneExpire.IsNull() {
		newZone.ZoneExpire = int(plan.ZoneExpire.ValueInt64())
	}
	if !plan.ZoneMinimum.IsNull() {
		newZone.ZoneMinimum = int(plan.ZoneMinimum.ValueInt64())
	}
	if !plan.ZoneRefresh.IsNull() {
		newZone.ZoneRefresh = int(plan.ZoneRefresh.ValueInt64())
	}
	if !plan.ZoneRetry.IsNull() {
		newZone.ZoneRetry = int(plan.ZoneRetry.ValueInt64())
	}
	if !plan.ZoneSerial.IsNull() {
		newZone.ZoneSerial = int(plan.ZoneSerial.ValueInt64())
	}
	if !plan.ZoneTTL.IsNull() {
		newZone.ZoneTTL = int(plan.ZoneTTL.ValueInt64())
	}

	// Create new order
	dnszone, err := r.client.DNS.AddZone(newZone)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating new ProVision DNS Zone",
			"Could not create ProVision DNS Zone, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(string(dnszone.ID))
	plan.Modified = types.StringValue(plan.Modified.ValueString())
	plan.ZoneHost = types.StringValue(dnszone.ZoneHost)
	plan.ZoneMail = types.StringValue(dnszone.ZoneMail)
	plan.ZoneType = types.StringValue(dnszone.ZoneType)
	plan.ParentID = types.StringValue(string(dnszone.ParentID))
	plan.Status = types.StringValue(dnszone.Status)
	plan.ZoneExpire = types.Int64Value(int64(dnszone.ZoneExpire))
	plan.ZoneMinimum = types.Int64Value(int64(dnszone.ZoneMinimum))
	plan.ZoneRefresh = types.Int64Value(int64(dnszone.ZoneRefresh))
	plan.ZoneRetry = types.Int64Value(int64(dnszone.ZoneRetry))
	plan.ZoneSerial = types.Int64Value(int64(dnszone.ZoneSerial))
	plan.ZoneTTL = types.Int64Value(int64(dnszone.ZoneTTL))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *dnszoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Read resource information
func (r *dnszoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dnszoneModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	zones, err := r.client.DNS.GetZoneByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ProVision DNS Zone",
			"Could not read ProVision DNS Zone ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(zones) == 0 {
		resp.Diagnostics.AddError(
			"Error Finding ProVision DNS Zone",
			"ProVision DNS Zone has not been found ",
		)
		return
	}

	dnszone := zones[0]

	state.Modified = types.StringValue(state.Modified.ValueString())
	state.ZoneHost = types.StringValue(dnszone.ZoneHost)
	state.ZoneMail = types.StringValue(dnszone.ZoneMail)
	state.ZoneType = types.StringValue(dnszone.ZoneType)
	state.ParentID = types.StringValue(string(dnszone.ParentID))
	state.Status = types.StringValue(dnszone.Status)
	state.ZoneExpire = types.Int64Value(int64(dnszone.ZoneExpire))
	state.ZoneMinimum = types.Int64Value(int64(dnszone.ZoneMinimum))
	state.ZoneRefresh = types.Int64Value(int64(dnszone.ZoneRefresh))
	state.ZoneRetry = types.Int64Value(int64(dnszone.ZoneRetry))
	state.ZoneSerial = types.Int64Value(int64(dnszone.ZoneSerial))
	state.ZoneTTL = types.Int64Value(int64(dnszone.ZoneTTL))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dnszoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan dnszoneModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Updating DNS Zone ID "+plan.ID.ValueString())

	newZone := provisionclient.DNSZone{
		ID:          provisionclient.PVID(plan.ID.ValueString()),
		Name:        plan.Name.ValueString(),
		ParentID:    provisionclient.PVID(plan.ParentID.ValueString()),
		ZoneType:    plan.ZoneType.ValueString(),
		ZoneHost:    plan.ZoneHost.ValueString(),
		ZoneMail:    plan.ZoneMail.ValueString(),
		ZoneExpire:  int(plan.ZoneExpire.ValueInt64()),
		ZoneMinimum: int(plan.ZoneMinimum.ValueInt64()),
		ZoneRefresh: int(plan.ZoneRefresh.ValueInt64()),
		ZoneRetry:   int(plan.ZoneRetry.ValueInt64()),
		ZoneSerial:  int(plan.ZoneSerial.ValueInt64()),
		ZoneTTL:     int(plan.ZoneTTL.ValueInt64()),
	}

	// Update existing order
	dnszone, err := r.client.DNS.UpdateZone(newZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ProVision DNS Zone",
			"Could not update ProVision DNS Zone, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ZoneHost = types.StringValue(dnszone.ZoneHost)
	plan.ZoneMail = types.StringValue(dnszone.ZoneMail)
	plan.ZoneType = types.StringValue(dnszone.ZoneType)
	plan.ParentID = types.StringValue(string(dnszone.ParentID))
	plan.Status = types.StringValue(dnszone.Status)
	plan.ZoneExpire = types.Int64Value(int64(dnszone.ZoneExpire))
	plan.ZoneMinimum = types.Int64Value(int64(dnszone.ZoneMinimum))
	plan.ZoneRefresh = types.Int64Value(int64(dnszone.ZoneRefresh))
	plan.ZoneRetry = types.Int64Value(int64(dnszone.ZoneRetry))
	plan.ZoneSerial = types.Int64Value(int64(dnszone.ZoneSerial))
	plan.ZoneTTL = types.Int64Value(int64(dnszone.ZoneTTL))
	plan.Modified = types.StringValue(plan.Modified.ValueString())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnszoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dnszoneModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DNS.DeleteZoneByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ProVision DNS Zone",
			"Could not delete ProVision DNS Zone, unexpected error: "+err.Error(),
		)
		return
	}
}
