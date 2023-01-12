package provision6connect

import (
	"context"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// firstavailableipDataSourceModel maps the data source schema data.
type firstavailableipDataSourceModel struct {
	NetblockID       types.String `tfsdk:"netblock_id"`
	NetblockCIDR     types.String `tfsdk:"netblock_cidr"`
	Firstavailableip types.String `tfsdk:"firstavailableip"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &firstavailableipDataSource{}
	_ datasource.DataSourceWithConfigure = &firstavailableipDataSource{}
)

// NewFirstavailableipDataSource is a helper function to simplify the provider implementation.
func NewFirstavailableipDataSource() datasource.DataSource {
	return &firstavailableipDataSource{}
}

// firstavailableipDataSource is the data source implementation.
type firstavailableipDataSource struct {
	client *provisionclient.Client
}

// Metadata returns the data source type name.
func (d *firstavailableipDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firstavailableip"
}

// Configure adds the provider configured client to the data source.
func (d *firstavailableipDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*provisionclient.Client)
}

// Schema defines the schema for the data source.
func (d *firstavailableipDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"netblock_id": schema.StringAttribute{
				Description:         "Get Next Available IP by Netblock ID",
				MarkdownDescription: "Get Next Available IP by Netblock ID",
				Optional:            true,
			},
			"netblock_cidr": schema.StringAttribute{
				Description:         "Get Next Available IP by CIDR",
				MarkdownDescription: "Get Next Available IP by CIDR",
				Optional:            true,
			},
			"firstavailableip": schema.StringAttribute{
				Description: "Output the first available IP",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *firstavailableipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firstavailableipDataSourceModel
	search := ""
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.NetblockID.IsNull() {
		search = state.NetblockID.ValueString()
	} else if !state.NetblockCIDR.IsNull() {
		search = state.NetblockCIDR.ValueString()
	} else {
		resp.Diagnostics.AddError(
			"Either netblock_cidr or netblock_id are required",
			"Either netblock_cidr or netblock_id are required",
		)
		return
	}

	firstavailableip, err := d.client.IPAM.GetFirstAvailable(search)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ProVision FirstAvailableIP",
			err.Error(),
		)
		return
	}

	state.Firstavailableip = types.StringValue(*firstavailableip)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
