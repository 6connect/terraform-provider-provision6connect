package provision6connect

import (
	"context"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resourcesDataSourceModel maps the data source schema data.
type resourcesDataSourceModel struct {
	Search    map[string]string `tfsdk:"search"`
	Resources []resourcesModel  `tfsdk:"resources"`
}

// resourcesModel maps resources schema data.
type resourcesModel struct {
	ID       types.String      `tfsdk:"id"`
	ParentID types.String      `tfsdk:"parent_id"`
	Name     types.String      `tfsdk:"name"`
	Slug     types.String      `tfsdk:"slug"`
	Type     types.String      `tfsdk:"type"`
	Modified types.String      `tfsdk:"modified"`
	Attrs    map[string]string `tfsdk:"attrs"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &resourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &resourcesDataSource{}
)

// NewResourcesDataSource is a helper function to simplify the provider implementation.
func NewResourcesDataSource() datasource.DataSource {
	return &resourcesDataSource{}
}

// resourcesDataSource is the data source implementation.
type resourcesDataSource struct {
	client *provisionclient.Client
}

// Metadata returns the data source type name.
func (d *resourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resources"
}

// Configure adds the provider configured client to the data source.
func (d *resourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*provisionclient.Client)
}

// Schema defines the schema for the data source.
func (d *resourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resources Data Source for query resource information on ProVision",
		Attributes: map[string]schema.Attribute{
			"search": schema.MapAttribute{
				ElementType:         types.StringType,
				Description:         "The map will be used into the API request to retrieve resource data",
				MarkdownDescription: "The map will be used into the API request to retrieve resource data",
				Optional:            true,
			},
			"resources": schema.ListNestedAttribute{
				Description: "Resource List of the resource found by the search query",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Numeric identifier of the Resource",
							Computed:    true,
						},
						"parent_id": schema.StringAttribute{
							Description: "Parent Resource identifier Number",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Resource Name",
							Computed:    true,
						},
						"slug": schema.StringAttribute{
							Description: "Resource Slug",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Resource Type",
							Computed:    true,
						},
						"modified": schema.StringAttribute{
							Description: "Date and Time of the last modification",
							Computed:    true,
						},
						"attrs": schema.MapAttribute{
							Description: "Resource Attributes List",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *resourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state resourcesDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	_, ok := state.Search["load_attributes"]
	if !ok {
		state.Search["load_attributes"] = "1"
	}
	resources, err := d.client.Resources.GetResources(&state.Search)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ProVision Resources",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, resource := range resources {
		/*
			attrMap := make(map[string]attr.Value, len(resource.Attrs))
			for key, val := range resource.Attrs {
				attrMap[key] = types.StringValue(val)
			}

			mapval, _ := types.MapValue(types.StringType, attrMap)
		*/
		resourceState := resourcesModel{
			ID:       types.StringValue(string(resource.ID)),
			ParentID: types.StringValue(string(resource.ParentID)),
			Name:     types.StringValue(resource.Name),
			Slug:     types.StringValue(resource.Slug),
			Type:     types.StringValue(resource.Type),
			Modified: types.StringValue(resource.Modified),
			Attrs:    resource.Attrs,
		}
		state.Resources = append(state.Resources, resourceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
