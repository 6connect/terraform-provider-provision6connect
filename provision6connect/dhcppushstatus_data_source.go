package provision6connect

import (
	"context"
	"time"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DHCPPushStatusMessage struct {
	MSGid       types.String `tfsdk:"msgid"`
	Message     types.String `tfsdk:"message"`
	State       types.String `tfsdk:"state"`
	DateCreated types.String `tfsdk:"date_created"`
}

// dhcppushstatusDataSourceModel maps the data source schema data.
type dhcppushstatusDataSourceModel struct {
	ServerID       types.String            `tfsdk:"server_id"`
	GroupID        types.String            `tfsdk:"group_id"`
	PoolID         types.String            `tfsdk:"pool_id"`
	PushPID        types.String            `tfsdk:"push_pid"`
	Delay          types.Int64             `tfsdk:"delay"`
	StatusMessages []DHCPPushStatusMessage `tfsdk:"status_messages"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dhcppushstatusDataSource{}
	_ datasource.DataSourceWithConfigure = &dhcppushstatusDataSource{}
)

// NewDHCPpushstatusDataSource is a helper function to simplify the provider implementation.
func NewDHCPpushstatusDataSource() datasource.DataSource {
	return &dhcppushstatusDataSource{}
}

// dhcppushstatusDataSource is the data source implementation.
type dhcppushstatusDataSource struct {
	client *provisionclient.Client
}

// Metadata returns the data source type name.
func (d *dhcppushstatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcppushstatus"
}

// Configure adds the provider configured client to the data source.
func (d *dhcppushstatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*provisionclient.Client)
}

// Schema defines the schema for the data source.
func (d *dhcppushstatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Description:         "Group Resource ID from the push request",
				MarkdownDescription: "Group Resource ID from the push request",
				Optional:            true,
			},
			"server_id": schema.StringAttribute{
				Description:         "Server Resource ID from the push request",
				MarkdownDescription: "Server Resource ID from the push request",
				Optional:            true,
			},
			"pool_id": schema.StringAttribute{
				Description:         "Pool Resource ID from the push request",
				MarkdownDescription: "Pool Resource ID from the push request",
				Optional:            true,
			},
			"delay": schema.Int64Attribute{
				Description:         "Time to wait before executing the status request",
				MarkdownDescription: "Time to wait before executing the status request",
				Optional:            true,
			},
			"push_pid": schema.StringAttribute{
				Description: "Push Request PID",
				Required:    true,
			},
			"status_messages": schema.ListNestedAttribute{
				Description: "Output containing a list of status messages",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"msgid": schema.StringAttribute{
							Description: "Status Message ID",
							Computed:    true,
						},
						"message": schema.StringAttribute{
							Description: "Message containing a description or the performed action",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "Current Execution State : running, finished, warning, error",
							Computed:    true,
						},
						"date_created": schema.StringAttribute{
							Description: "Date and Time of the message",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *dhcppushstatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dhcppushstatusDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Delay.IsNull() {
		time.Sleep(time.Duration(state.Delay.ValueInt64()) * time.Millisecond)
	}

	var messages []provisionclient.DHCPPushStatusMessage
	var err error

	if !state.GroupID.IsNull() {
		messages, err = d.client.DHCP.GetGroupPushStatus(state.GroupID.ValueString(), state.PushPID.ValueString())
	} else if !state.ServerID.IsNull() {
		messages, err = d.client.DHCP.GetServerPushStatus(state.ServerID.ValueString(), state.PushPID.ValueString())
	} else if !state.PoolID.IsNull() {
		messages, err = d.client.DHCP.GetPoolPushStatus(state.PoolID.ValueString(), state.PushPID.ValueString())
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

	for _, message := range messages {
		state.StatusMessages = append(state.StatusMessages, DHCPPushStatusMessage{
			MSGid:       types.StringValue(message.MSGid),
			Message:     types.StringValue(message.Message),
			State:       types.StringValue(message.State),
			DateCreated: types.StringValue(message.DateCreated),
		})
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
