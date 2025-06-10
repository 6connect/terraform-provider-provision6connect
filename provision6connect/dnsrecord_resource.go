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
	_ resource.Resource                = &dnsrecordResource{}
	_ resource.ResourceWithConfigure   = &dnsrecordResource{}
	_ resource.ResourceWithImportState = &dnsrecordResource{}
)

// NewDNSrecordResource is a helper function to simplify the provider implementation.
func NewDNSrecordResource() resource.Resource {
	return &dnsrecordResource{}
}

// resourcesModel maps resources schema data.
type dnsrecordModel struct {
	ID       types.String `tfsdk:"id"`
	ZoneID   types.String `tfsdk:"zone_id"`
	Name     types.String `tfsdk:"name"`
	Modified types.String `tfsdk:"modified"`

	Status      types.String `tfsdk:"status"`
	RecordType  types.String `tfsdk:"record_type"`
	RecordHost  types.String `tfsdk:"record_host"`
	RecordValue types.String `tfsdk:"record_value"`
	RecordTTL   types.Int64  `tfsdk:"record_ttl"`
}

// dnsrecordResource is the resource implementation.
type dnsrecordResource struct {
	client *provisionclient.Client
}

// Configure adds the provider configured client to the resource.
func (r *dnsrecordResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*provisionclient.Client)
}

// Metadata returns the resource type name.
func (r *dnsrecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnsrecord"
}

// Schema defines the schema for the resource.
func (r *dnsrecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "DNS Record Resource that represents a single DNS Record in ProVision",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the DNS Record.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Pretty name describing the DNS Record.",
				Required:    true,
			},
			"zone_id": schema.StringAttribute{
				Description: "Numeric identifier of the DNS Zone that contains the DNS Record.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				Description: "Date and Time of the last modification",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Current status set by ProVision of the DNS Record",
				Computed:    true,
			},
			"record_type": schema.StringAttribute{
				Description: "DNS Record Type Ex: A, AAAA, TXT, PTR, NS",
				Required:    true,
			},
			"record_host": schema.StringAttribute{
				Description: "FQDN of the DNS Record Ex: test.example.com.",
				Required:    true,
			},
			"record_value": schema.StringAttribute{
				Description: "DNS Record Value Ex: 192.168.0.1",
				Required:    true,
			},
			"record_ttl": schema.Int64Attribute{
				Description: "DNS Record TTL Ex: 900",
				Required:    true,
			},
		},
	}
}

// Create a new resource
func (r *dnsrecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dnsrecordModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newRecord := provisionclient.DNSRecord{
		Name:        plan.Name.ValueString(),
		ParentID:    provisionclient.PVID(plan.ZoneID.ValueString()),
		RecordType:  plan.RecordType.ValueString(),
		RecordValue: plan.RecordValue.ValueString(),
		RecordHost:  plan.RecordHost.ValueString(),
		RecordTTL:   int(plan.RecordTTL.ValueInt64()),
	}

	// Create new order
	dnsrecord, err := r.client.DNS.AddZoneRecord(newRecord)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating new ProVision DNS Record",
			"Could not create ProVision DNS Record, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(string(dnsrecord.ID))
	plan.Modified = types.StringValue(dnsrecord.Modified)
	plan.Status = types.StringValue(dnsrecord.Status)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *dnsrecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Read resource information
func (r *dnsrecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dnsrecordModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed zone records value
	records, err := r.client.DNS.GetZoneRecords(state.ZoneID.ValueString(), &map[string]string{
		"id":              state.ID.ValueString(),
		"load_attributes": "1",
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ProVision DNS Record",
			"Could not read ProVision DNS Record ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(records) == 0 {
		resp.Diagnostics.AddError(
			"Error Finding ProVision DNS Record",
			"ProVision DNS Record has not been found ",
		)
		return
	}

	dnsrecord := records[0]

	state.ZoneID = types.StringValue(string(dnsrecord.ParentID))
	state.Modified = types.StringValue(dnsrecord.Modified)
	state.Status = types.StringValue(dnsrecord.Status)
	state.RecordHost = types.StringValue(dnsrecord.RecordHost)
	state.RecordValue = types.StringValue(dnsrecord.RecordValue)
	state.RecordType = types.StringValue(dnsrecord.RecordType)
	state.RecordTTL = types.Int64Value(int64(dnsrecord.RecordTTL))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dnsrecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan dnsrecordModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Updating DNS Record ID "+plan.ID.ValueString())
	newRecord := provisionclient.DNSRecord{
		ID:          provisionclient.PVID(plan.ID.ValueString()),
		Name:        plan.Name.ValueString(),
		ParentID:    provisionclient.PVID(plan.ZoneID.ValueString()),
		RecordType:  plan.RecordType.ValueString(),
		RecordValue: plan.RecordValue.ValueString(),
		RecordHost:  plan.RecordHost.ValueString(),
		RecordTTL:   int(plan.RecordTTL.ValueInt64()),
	}

	// Update existing order
	dnsrecord, err := r.client.DNS.UpdateZoneRecord(newRecord)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ProVision DNS Record",
			"Could not update ProVision DNS Record, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(string(dnsrecord.ID))
	plan.ZoneID = types.StringValue(string(dnsrecord.ParentID))
	plan.Modified = types.StringValue(dnsrecord.Modified)
	plan.Status = types.StringValue(dnsrecord.Status)
	plan.RecordHost = types.StringValue(dnsrecord.RecordHost)
	plan.RecordValue = types.StringValue(dnsrecord.RecordValue)
	plan.RecordType = types.StringValue(dnsrecord.RecordType)
	plan.RecordTTL = types.Int64Value(int64(dnsrecord.RecordTTL))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsrecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dnsrecordModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DNS.DeleteZoneRecordByID(state.ZoneID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ProVision DNS Record",
			"Could not delete ProVision DNS Record, unexpected error: "+err.Error(),
		)
		return
	}
}
