package provision6connect

import (
	"context"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// netblocksDataSourceModel maps the data source schema data.
type netblocksDataSourceModel struct {
	Search    map[string]string `tfsdk:"search"`
	Netblocks []netblocksModel  `tfsdk:"netblocks"`
}

// netblocksModel maps netblocks schema data.
type netblocksModel struct {
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
	Range               []string     `tfsdk:"range"`
	Tags                []string     `tfsdk:"tags"`
	UtilizationStatus   types.String `tfsdk:"utilization_status"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &netblocksDataSource{}
	_ datasource.DataSourceWithConfigure = &netblocksDataSource{}
)

// NewNetblocksDataSource is a helper function to simplify the provider implementation.
func NewNetblocksDataSource() datasource.DataSource {
	return &netblocksDataSource{}
}

// netblocksDataSource is the data source implementation.
type netblocksDataSource struct {
	client *provisionclient.Client
}

// Metadata returns the data source type name.
func (d *netblocksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_netblocks"
}

// Configure adds the provider configured client to the data source.
func (d *netblocksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*provisionclient.Client)
}

// Schema defines the schema for the data source.
func (d *netblocksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Netblock Data Source for query netblock information on ProVision",
		Attributes: map[string]schema.Attribute{
			"search": schema.MapAttribute{
				Description:         "In the Search List you can provide parameters that are accepted by the ProVision IPAM Netblocks GET API",
				ElementType:         types.StringType,
				MarkdownDescription: "The map will be used into the API request to retrieve netblock data",
				Optional:            true,
			},
			"netblocks": schema.ListNestedAttribute{
				Description: "Contains a list of the NetBlocks found by the search query",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Numeric identifier of the NetBlock.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "IP Type can be either ipv4 or ipv6",
							Computed:    true,
						},
						"top_aggregate": schema.StringAttribute{
							Description: "Top Aggregate Netblock ID",
							Computed:    true,
						},
						"cidr": schema.StringAttribute{
							Description: "CIDR of the netblock",
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
							Computed:    true,
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
							Computed: true,
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
							Computed: true,
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
							Computed: true,
						},
						"org_id": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"region_id": schema.StringAttribute{
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
							Computed: true,
						},
						"meta2": schema.StringAttribute{
							Computed: true,
						},
						"meta3": schema.StringAttribute{
							Computed: true,
						},
						"meta4": schema.StringAttribute{
							Computed: true,
						},
						"meta5": schema.StringAttribute{
							Computed: true,
						},
						"meta6": schema.StringAttribute{
							Computed: true,
						},
						"meta7": schema.StringAttribute{
							Computed: true,
						},
						"meta8": schema.StringAttribute{
							Computed: true,
						},
						"meta9": schema.StringAttribute{
							Computed: true,
						},
						"meta10": schema.StringAttribute{
							Computed: true,
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
							ElementType: types.StringType,
							Computed:    true,
						},
						"utilization_status": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *netblocksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state netblocksDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	netblocks, err := d.client.IPAM.GetNetblocks(&state.Search)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ProVision Netblocks",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, netblock := range netblocks {
		/*
			rangeList := make([]attr.Value, len(netblock.Range))
			for _, val := range netblock.Range {
				rangeList = append(rangeList, types.StringValue(val))
			}
			rangeListValue, _ := types.ListValue(types.StringType, rangeList)

			tagsList := make([]attr.Value, len(netblock.Tags))
			for _, val := range netblock.Tags {
				tagsList = append(tagsList, types.StringValue(val))
			}
			tagsListValue, _ := types.ListValue(types.StringType, tagsList)
		*/
		netblockState := netblocksModel{
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
			Range:               netblock.Range,
			Tags:                netblock.Tags,
			UtilizationStatus:   types.StringValue(netblock.UtilizationStatus),
		}

		state.Netblocks = append(state.Netblocks, netblockState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
