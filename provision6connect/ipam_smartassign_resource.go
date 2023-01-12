package provision6connect

import (
	"context"
	"strings"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &ipamsmartassignResource{}
	_ resource.ResourceWithConfigure = &ipamsmartassignResource{}
)

// NewIPAMsmartassignResource is a helper function to simplify the provider implementation.
func NewIPAMsmartassignResource() resource.Resource {
	return &ipamsmartassignResource{}
}

// resourcesModel maps resources schema data.
type ipamsmartassignModel struct {
	ID                  types.String `tfsdk:"id"`
	Type                types.String `tfsdk:"type"`
	TopAggregate        types.String `tfsdk:"top_aggregate"`
	CIDR                types.String `tfsdk:"cidr"`
	Address             types.String `tfsdk:"address"`
	EndAddress          types.String `tfsdk:"end_address"`
	IsAggregate         types.Bool   `tfsdk:"is_aggregate"`
	Assigned            types.Bool   `tfsdk:"assigned"`
	SparseAllocationId  types.String `tfsdk:"sparse_allocation_id"`
	IsImportant         types.Bool   `tfsdk:"is_important"`
	Swipped             types.Bool   `tfsdk:"swipped"`
	LastUpdateTime      types.String `tfsdk:"last_update_time"`
	LIRID               types.String `tfsdk:"lir_id"`
	Mask                types.Int64  `tfsdk:"mask"`
	NetMask             types.String `tfsdk:"netmask"`
	ASN                 types.String `tfsdk:"asn"`
	AllowSubAssignments types.Bool   `tfsdk:"allow_sub_assignments"`
	Child1              types.String `tfsdk:"child1"`
	Child2              types.String `tfsdk:"child2"`
	ResourceID          types.String `tfsdk:"resource_id"`
	ResourceName        types.String `tfsdk:"resource_name"`
	Description         types.String `tfsdk:"description"`
	Parent              types.String `tfsdk:"parent"`
	RIR                 types.String `tfsdk:"rir"`
	Notes               types.String `tfsdk:"notes"`
	GenericCode         types.String `tfsdk:"generic_code"`
	AssignTime          types.String `tfsdk:"assign_time"`
	SWIPTime            types.String `tfsdk:"swip_time"`
	NetHandle           types.String `tfsdk:"net_handle"`
	CustomerHandle      types.String `tfsdk:"customer_handle"`
	VLANID              types.String `tfsdk:"vlan_id"`
	ORGID               types.String `tfsdk:"org_id"`
	Region              types.String `tfsdk:"region"`
	RegionID            types.String `tfsdk:"region_id"`
	RuleID              types.String `tfsdk:"rule_id"`
	ReservedTime        types.String `tfsdk:"reserved_time"`
	ReservedBy          types.String `tfsdk:"reserved_by"`
	DHCPResourceID      types.String `tfsdk:"dhcp_resource_id"`
	CMNETBLOCKID        types.String `tfsdk:"cmnetblock_resource_id"`
	UMBRELLAID          types.String `tfsdk:"umbrella_resource_id"`
	Meta1               types.String `tfsdk:"meta1"`
	Meta2               types.String `tfsdk:"meta2"`
	Meta3               types.String `tfsdk:"meta3"`
	Meta4               types.String `tfsdk:"meta4"`
	Meta5               types.String `tfsdk:"meta5"`
	Meta6               types.String `tfsdk:"meta6"`
	Meta7               types.String `tfsdk:"meta7"`
	Meta8               types.String `tfsdk:"meta8"`
	Meta9               types.String `tfsdk:"meta9"`
	Meta10              types.String `tfsdk:"meta10"`
	NAT                 types.String `tfsdk:"nat"`
	HostCount           types.String `tfsdk:"host_count"`
	RegionName          types.String `tfsdk:"region_name"`
	Range               types.List   `tfsdk:"range"`
	Tags                types.List   `tfsdk:"tags"`
	UtilizationStatus   types.String `tfsdk:"utilization_status"`
	//not in the netblock, but into the Sheme
	AssignedResourceID types.String `tfsdk:"assigned_resource_id"`
}

func netblockToState(ctx context.Context, netblock *provisionclient.Netblock) ipamsmartassignModel {
	range_list, _ := types.ListValueFrom(ctx, types.StringType, netblock.Range)
	tags_list, _ := types.ListValueFrom(ctx, types.StringType, netblock.Tags)
	return ipamsmartassignModel{
		ID:                  types.StringValue(string(netblock.ID)),
		Type:                types.StringValue(netblock.Type),
		TopAggregate:        types.StringValue(string(netblock.TopAggregate)),
		CIDR:                types.StringValue(netblock.CIDR),
		Address:             types.StringValue(netblock.Address),
		EndAddress:          types.StringValue(netblock.EndAddress),
		IsAggregate:         types.BoolValue(netblock.IsAggregate),
		Assigned:            types.BoolValue(netblock.Assigned),
		SparseAllocationId:  types.StringValue(string(netblock.SparseAllocationId)),
		IsImportant:         types.BoolValue(netblock.IsImportant),
		Swipped:             types.BoolValue(netblock.Swipped),
		LastUpdateTime:      types.StringValue(netblock.LastUpdateTime),
		LIRID:               types.StringValue(string(netblock.LIRID)),
		Mask:                types.Int64Value(int64(netblock.Mask)),
		NetMask:             types.StringValue(netblock.NetMask),
		ASN:                 types.StringValue(string(netblock.ASN)),
		AllowSubAssignments: types.BoolValue(netblock.AllowSubAssignments),
		Child1:              types.StringValue(string(netblock.Child1)),
		Child2:              types.StringValue(string(netblock.Child2)),
		ResourceID:          types.StringValue(string(netblock.ResourceID)),
		ResourceName:        types.StringValue(netblock.ResourceName),
		Description:         types.StringValue(netblock.Description),
		Parent:              types.StringValue(string(netblock.Parent)),
		RIR:                 types.StringValue(netblock.RIR),
		Notes:               types.StringValue(netblock.Notes),
		GenericCode:         types.StringValue(netblock.GenericCode),
		AssignTime:          types.StringValue(netblock.AssignTime),
		SWIPTime:            types.StringValue(netblock.SWIPTime),
		NetHandle:           types.StringValue(netblock.NetHandle),
		CustomerHandle:      types.StringValue(netblock.CustomerHandle),
		VLANID:              types.StringValue(string(netblock.VLANID)),
		ORGID:               types.StringValue(string(netblock.ORGID)),
		Region:              types.StringValue(netblock.Region),
		RegionID:            types.StringValue(string(netblock.RegionID)),
		RuleID:              types.StringValue(string(netblock.RuleID)),
		ReservedTime:        types.StringValue(netblock.ReservedTime),
		ReservedBy:          types.StringValue(string(netblock.ReservedBy)),
		DHCPResourceID:      types.StringValue(string(netblock.DHCPResourceID)),
		CMNETBLOCKID:        types.StringValue(string(netblock.CMNETBLOCKID)),
		UMBRELLAID:          types.StringValue(string(netblock.UMBRELLAID)),
		Meta1:               types.StringValue(netblock.Meta1),
		Meta2:               types.StringValue(netblock.Meta2),
		Meta3:               types.StringValue(netblock.Meta3),
		Meta4:               types.StringValue(netblock.Meta4),
		Meta5:               types.StringValue(netblock.Meta5),
		Meta6:               types.StringValue(netblock.Meta6),
		Meta7:               types.StringValue(netblock.Meta7),
		Meta8:               types.StringValue(netblock.Meta8),
		Meta9:               types.StringValue(netblock.Meta9),
		Meta10:              types.StringValue(netblock.Meta10),
		NAT:                 types.StringValue(netblock.NAT),
		HostCount:           types.StringValue(netblock.HostCount),
		RegionName:          types.StringValue(netblock.RegionName),
		Range:               range_list,
		Tags:                tags_list,
		UtilizationStatus:   types.StringValue(netblock.UtilizationStatus),
	}
}

// ipamsmartassignResource is the resource implementation.
type ipamsmartassignResource struct {
	client *provisionclient.Client
}

// Configure adds the provider configured client to the resource.
func (r *ipamsmartassignResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*provisionclient.Client)
}

// Metadata returns the resource type name.
func (r *ipamsmartassignResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smartassign"
}

// Schema defines the schema for the resource.
func (r *ipamsmartassignResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Smart Assign IPAM Netblock by given RIR,Mask,Type,Resource ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the NetBlock.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "IP Type can be either ipv4 or ipv6",
				Required:    true,
			},
			"top_aggregate": schema.StringAttribute{
				Description: "Top Aggregate Netblock ID",
				Optional:    true,
				Computed:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "CIDR of the netblock that will have the assignment",
				Computed:    true,
			},
			"address": schema.StringAttribute{
				Description: "Numeric Start IP Address",
				Computed:    true,
			},
			"end_address": schema.StringAttribute{
				Description: "Numeric End IP Address",
				Computed:    true,
			},
			"is_aggregate": schema.BoolAttribute{
				Description: "Set to True if the netblock is an Aggregate",
				Computed:    true,
			},
			"assigned": schema.BoolAttribute{
				Computed: true,
			},
			"sparse_allocation_id": schema.StringAttribute{
				Computed: true,
			},
			"is_important": schema.BoolAttribute{
				Computed: true,
			},
			"swipped": schema.BoolAttribute{
				Computed: true,
			},
			"last_update_time": schema.StringAttribute{
				Computed: true,
			},
			"lir_id": schema.StringAttribute{
				Computed: true,
			},
			"mask": schema.Int64Attribute{
				Description: "Numeric representation of the mask",
				Required:    true,
			},
			"netmask": schema.StringAttribute{
				Computed: true,
			},
			"asn": schema.StringAttribute{
				Computed: true,
			},
			"allow_sub_assignments": schema.BoolAttribute{
				Computed: true,
			},
			"child1": schema.StringAttribute{
				Computed: true,
			},
			"child2": schema.StringAttribute{
				Computed: true,
			},
			"resource_id": schema.StringAttribute{
				Description: "Assigned Resource ID",
				Required:    true,
			},
			"resource_name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"parent": schema.StringAttribute{
				Computed: true,
			},
			"rir": schema.StringAttribute{
				Description: "RIR of the Netblock",
				Required:    true,
			},
			"notes": schema.StringAttribute{
				Computed: true,
			},
			"generic_code": schema.StringAttribute{
				Computed: true,
			},
			"assign_time": schema.StringAttribute{
				Computed: true,
			},
			"swip_time": schema.StringAttribute{
				Computed: true,
			},
			"net_handle": schema.StringAttribute{
				Computed: true,
			},
			"customer_handle": schema.StringAttribute{
				Computed: true,
			},
			"vlan_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"org_id": schema.StringAttribute{
				Computed: true,
			},
			"region": schema.StringAttribute{
				Computed: true,
			},
			"region_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"rule_id": schema.StringAttribute{
				Computed: true,
			},
			"reserved_time": schema.StringAttribute{
				Computed: true,
			},
			"reserved_by": schema.StringAttribute{
				Computed: true,
			},
			"dhcp_resource_id": schema.StringAttribute{
				Computed: true,
			},
			"cmnetblock_resource_id": schema.StringAttribute{
				Computed: true,
			},
			"umbrella_resource_id": schema.StringAttribute{
				Computed: true,
			},
			"meta1": schema.StringAttribute{
				Description: "Meta1 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta2": schema.StringAttribute{
				Description: "Meta2 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta3": schema.StringAttribute{
				Description: "Meta3 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta4": schema.StringAttribute{
				Description: "Meta4 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta5": schema.StringAttribute{
				Description: "Meta5 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta6": schema.StringAttribute{
				Description: "Meta6 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta7": schema.StringAttribute{
				Description: "Meta7 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta8": schema.StringAttribute{
				Description: "Meta8 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta9": schema.StringAttribute{
				Description: "Meta9 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"meta10": schema.StringAttribute{
				Description: "Meta10 IPAM attribute",
				Optional:    true,
				Computed:    true,
			},
			"nat": schema.StringAttribute{
				Computed: true,
			},
			"host_count": schema.StringAttribute{
				Computed: true,
			},
			"region_name": schema.StringAttribute{
				Computed: true,
			},

			"range": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},

			"tags": schema.ListAttribute{
				Description: "Netblock Tags list",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},

			"utilization_status": schema.StringAttribute{
				Computed: true,
			},
			"assigned_resource_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Create a new resource
func (r *ipamsmartassignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ipamsmartassignModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	/*
		if !plan.ParentID.IsNull() {
			newResource.ParentID = provisionclient.PVID(plan.ParentID.ValueString())
		}

		if !plan.Slug.IsNull() {
			newResource.Slug = plan.Slug.ValueString()
		}
	*/
	params := make(map[string]interface{})

	if !plan.Tags.IsNull() {
		tags := plan.Tags.Elements()
		tags_list := make([]string, len(tags))
		for key, val := range tags {
			tags_list[key] = val.String()
		}

		params["tags"] = strings.Join(tags_list, ",")
	}

	if !plan.TopAggregate.IsNull() {
		params["top_aggregate"] = plan.TopAggregate.ValueString()
	}
	if !plan.AssignedResourceID.IsNull() {
		params["assigned_resource_id"] = plan.AssignedResourceID.ValueString()
	}
	if !plan.VLANID.IsNull() {
		params["vlan"] = plan.VLANID.ValueString()
	}
	if !plan.RegionID.IsNull() {
		params["region_id"] = plan.RegionID.ValueString()
	}
	if !plan.Meta1.IsNull() {
		params["meta1"] = plan.Meta1.ValueString()
	}
	if !plan.Meta2.IsNull() {
		params["meta2"] = plan.Meta2.ValueString()
	}
	if !plan.Meta3.IsNull() {
		params["meta3"] = plan.Meta3.ValueString()
	}
	if !plan.Meta4.IsNull() {
		params["meta4"] = plan.Meta4.ValueString()
	}
	if !plan.Meta5.IsNull() {
		params["meta5"] = plan.Meta5.ValueString()
	}
	if !plan.Meta6.IsNull() {
		params["meta6"] = plan.Meta6.ValueString()
	}
	if !plan.Meta7.IsNull() {
		params["meta7"] = plan.Meta7.ValueString()
	}
	if !plan.Meta8.IsNull() {
		params["meta8"] = plan.Meta8.ValueString()
	}
	if !plan.Meta9.IsNull() {
		params["meta9"] = plan.Meta9.ValueString()
	}
	if !plan.Meta10.IsNull() {
		params["meta10"] = plan.Meta10.ValueString()
	}

	// Do Smart Assign
	tflog.Info(ctx, "Executing SmartAssign Request...")
	netblock, err := r.client.IPAM.SmartAssign(
		plan.ResourceID.ValueString(),
		plan.Type.ValueString(),
		plan.RIR.ValueString(),
		int(plan.Mask.ValueInt64()),
		params,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error using SmartAssign",
			"Could not SmartAssign a netblock, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	netblockState := netblockToState(ctx, netblock)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, netblockState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Read resource information
func (r *ipamsmartassignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ipamsmartassignModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	netblock, err := r.client.IPAM.GetNetblockByID(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ProVision Netblock",
			"Could not read ProVision Netblock ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	netblockState := netblockToState(ctx, netblock)

	// Set refreshed state
	diags = resp.State.Set(ctx, &netblockState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ipamsmartassignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ipamsmartassignModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Updating Netblock ID "+plan.ID.ValueString())
	newNetblock := provisionclient.Netblock{
		ID:                  provisionclient.PVID(plan.ID.ValueString()),
		AllowSubAssignments: plan.AllowSubAssignments.ValueBool(),
		//Tags:                plan.Tags,
		RIR:      plan.RIR.ValueString(),
		VLANID:   provisionclient.PVID(plan.VLANID.ValueString()),
		RuleID:   provisionclient.PVID(plan.RuleID.ValueString()),
		ASN:      provisionclient.PVID(plan.ASN.ValueString()),
		RegionID: provisionclient.PVID(plan.RegionID.ValueString()),
		LIRID:    provisionclient.PVID(plan.LIRID.ValueString()),
		Meta1:    plan.Meta1.ValueString(),
		Meta2:    plan.Meta2.ValueString(),
		Meta3:    plan.Meta3.ValueString(),
		Meta4:    plan.Meta4.ValueString(),
		Meta5:    plan.Meta5.ValueString(),
		Meta6:    plan.Meta6.ValueString(),
		Meta7:    plan.Meta7.ValueString(),
		Meta8:    plan.Meta8.ValueString(),
		Meta9:    plan.Meta9.ValueString(),
		Meta10:   plan.Meta10.ValueString(),
	}

	// Update existing order
	netblock, err := r.client.IPAM.UpdateNetblock(newNetblock)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ProVision NetBlock",
			"Could not update ProVision NetBlock, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	netblockState := netblockToState(ctx, netblock)

	diags = resp.State.Set(ctx, netblockState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ipamsmartassignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ipamsmartassignModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	_, err := r.client.IPAM.UnassignNetblockByID(state.ID.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ProVision NetBlock",
			"Could not delete ProVision NetBlock, unexpected error: "+err.Error(),
		)
		return
	}
}
