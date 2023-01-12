package provision6connect

import (
	"context"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dhcppushDataSourceModel maps the data source schema data.
type dhcppushDataSourceModel struct {
	ServerID types.String `tfsdk:"server_id"`
	GroupID  types.String `tfsdk:"group_id"`
	PoolID   types.String `tfsdk:"pool_id"`
	PushPID  types.String `tfsdk:"push_pid"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dhcppushDataSource{}
	_ datasource.DataSourceWithConfigure = &dhcppushDataSource{}
)

// NewDHCPpushDataSource is a helper function to simplify the provider implementation.
func NewDHCPpushDataSource() datasource.DataSource {
	return &dhcppushDataSource{}
}

// dhcppushDataSource is the data source implementation.
type dhcppushDataSource struct {
	client *provisionclient.Client
}

// Metadata returns the data source type name.
func (d *dhcppushDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcppush"
}

// Configure adds the provider configured client to the data source.
func (d *dhcppushDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*provisionclient.Client)
}

// Schema defines the schema for the data source.
func (d *dhcppushDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Executes Push Request to the DHCP Module. Either group_id, server_id or pool_id must be specified. The push status id is retured in push_pid.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Description:         "Group Resource ID to push",
				MarkdownDescription: "Group Resource ID to push",
				Optional:            true,
			},
			"server_id": schema.StringAttribute{
				Description:         "Server Resource ID to push",
				MarkdownDescription: "Server Resource ID to push",
				Optional:            true,
			},
			"pool_id": schema.StringAttribute{
				Description:         "Pool Resource ID to push",
				MarkdownDescription: "Pool Resource ID to push",
				Optional:            true,
			},
			"push_pid": schema.StringAttribute{
				Description: "Push PID is returned here that can be used for a pushstatus request",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dhcppushDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dhcppushDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pushpid *string
	var err error

	if !state.GroupID.IsNull() {
		pushpid, err = d.client.DHCP.PushGroupByID(state.GroupID.ValueString())
	} else if !state.ServerID.IsNull() {
		pushpid, err = d.client.DHCP.PushServerByID(state.ServerID.ValueString())
	} else if !state.PoolID.IsNull() {
		pushpid, err = d.client.DHCP.PushPoolByID(state.PoolID.ValueString())
	} else {
		resp.Diagnostics.AddError(
			"Either group_id or pool_id or server_id are required",
			"Either group_id or pool_id or server_id are required",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"The Push Request has returned an error",
			err.Error(),
		)
		return
	}

	state.PushPID = types.StringValue(*pushpid)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
