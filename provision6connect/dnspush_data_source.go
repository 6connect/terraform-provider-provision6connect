package provision6connect

import (
	"context"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dnspushDataSourceModel maps the data source schema data.
type dnspushDataSourceModel struct {
	ServerID types.String `tfsdk:"server_id"`
	GroupID  types.String `tfsdk:"group_id"`
	ZoneID   types.String `tfsdk:"zone_id"`
	PushPID  types.String `tfsdk:"push_pid"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnspushDataSource{}
	_ datasource.DataSourceWithConfigure = &dnspushDataSource{}
)

// NewDNSpushDataSource is a helper function to simplify the provider implementation.
func NewDNSpushDataSource() datasource.DataSource {
	return &dnspushDataSource{}
}

// dnspushDataSource is the data source implementation.
type dnspushDataSource struct {
	client *provisionclient.Client
}

// Metadata returns the data source type name.
func (d *dnspushDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnspush"
}

// Configure adds the provider configured client to the data source.
func (d *dnspushDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*provisionclient.Client)
}

// Schema defines the schema for the data source.
func (d *dnspushDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Executes Push Request to the DNS Module. Either group_id, server_id or zone_id must be specified. The push status id is retured in push_pid.",
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
			"zone_id": schema.StringAttribute{
				Description:         "Zone Resource ID to push",
				MarkdownDescription: "Zone Resource ID to push",
				Optional:            true,
			},
			"push_pid": schema.StringAttribute{
				Description: "Push PID is returned here that be used for a pushstatus request",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dnspushDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnspushDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pushpid *string
	var err error

	if !state.GroupID.IsNull() {
		pushpid, err = d.client.DNS.PushGroupByID(state.GroupID.ValueString())
	} else if !state.ServerID.IsNull() {
		pushpid, err = d.client.DNS.PushServerByID(state.ServerID.ValueString())
	} else if !state.ZoneID.IsNull() {
		pushpid, err = d.client.DNS.PushZoneByID(state.ZoneID.ValueString())
	} else {
		resp.Diagnostics.AddError(
			"Either group_id or zone_id or server_id are required",
			"Either group_id or zone_id or server_id are required",
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
