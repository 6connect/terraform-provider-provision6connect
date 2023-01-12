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
	_ resource.Resource                = &pvresourceResource{}
	_ resource.ResourceWithConfigure   = &pvresourceResource{}
	_ resource.ResourceWithImportState = &pvresourceResource{}
)

// NewPVresourceResource is a helper function to simplify the provider implementation.
func NewPVresourceResource() resource.Resource {
	return &pvresourceResource{}
}

// resourcesModel maps resources schema data.
type pvresourceModel struct {
	ID       types.String      `tfsdk:"id"`
	ParentID types.String      `tfsdk:"parent_id"`
	Name     types.String      `tfsdk:"name"`
	Slug     types.String      `tfsdk:"slug"`
	Type     types.String      `tfsdk:"type"`
	Modified types.String      `tfsdk:"modified"`
	Attrs    map[string]string `tfsdk:"attrs"`
}

// pvresourceResource is the resource implementation.
type pvresourceResource struct {
	client *provisionclient.Client
}

// Configure adds the provider configured client to the resource.
func (r *pvresourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*provisionclient.Client)
}

// Metadata returns the resource type name.
func (r *pvresourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

// Schema defines the schema for the resource.
func (r *pvresourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "ProVision internal Resource preresentation into Terraform",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the Resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Resource Name",
				Required:    true,
			},
			"parent_id": schema.StringAttribute{
				Description: "Parent Resource identifier Number",
				Optional:    true,
			},
			"slug": schema.StringAttribute{
				Description: "Resource Slug",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Resource Type",
				Required:    true,
			},
			"modified": schema.StringAttribute{
				Description: "Date and Time of the last modification",
				Computed:    true,
			},
			"attrs": schema.MapAttribute{
				Description: "Resource Attributes List",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

// Create a new resource
func (r *pvresourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pvresourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newResource := provisionclient.Resource{
		Name:  plan.Name.ValueString(),
		Type:  plan.Type.ValueString(),
		Attrs: plan.Attrs,
	}

	if !plan.ParentID.IsNull() {
		newResource.ParentID = provisionclient.PVID(plan.ParentID.ValueString())
	}

	if !plan.Slug.IsNull() {
		newResource.Slug = plan.Slug.ValueString()
	}

	// Create new order
	pvresource, err := r.client.Resources.AddResource(newResource)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating new ProVision Resource",
			"Could not create ProVision Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(string(pvresource.ID))
	plan.ParentID = types.StringValue(string(pvresource.ParentID))
	plan.Modified = types.StringValue(pvresource.Modified)
	plan.Slug = types.StringValue(pvresource.Slug)
	plan.Type = types.StringValue(pvresource.Type)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *pvresourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Read resource information
func (r *pvresourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pvresourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	resources, err := r.client.Resources.GetResources(&map[string]string{
		"id":              state.ID.ValueString(),
		"load_attributes": "1",
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ProVision Resource",
			"Could not read ProVision Resource ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(resources) == 0 {
		resp.Diagnostics.AddError(
			"Error Finding ProVision Resource",
			"ProVision Resource has not been found ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	pvresource := resources[0]

	state.ParentID = types.StringValue(string(pvresource.ParentID))
	state.Modified = types.StringValue(pvresource.Modified)
	state.Slug = types.StringValue(pvresource.Slug)
	state.Name = types.StringValue(pvresource.Name)
	state.Type = types.StringValue(pvresource.Type)
	state.Attrs = pvresource.Attrs

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pvresourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan pvresourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Updating Resource ID "+plan.ID.ValueString())
	newResource := provisionclient.Resource{
		ID:    provisionclient.PVID(plan.ID.ValueString()),
		Name:  plan.Name.ValueString(),
		Type:  plan.Type.ValueString(),
		Attrs: plan.Attrs,
	}

	if !plan.ParentID.IsNull() {
		newResource.ParentID = provisionclient.PVID(plan.ParentID.ValueString())
	}

	if !plan.Slug.IsNull() {
		newResource.Slug = plan.Slug.ValueString()
	}

	// Update existing order
	pvresource, err := r.client.Resources.UpdateResource(newResource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating ProVision Resource",
			"Could not update ProVision Resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(string(pvresource.ID))
	plan.ParentID = types.StringValue(string(pvresource.ParentID))
	plan.Modified = types.StringValue(pvresource.Modified)
	plan.Slug = types.StringValue(pvresource.Slug)
	plan.Type = types.StringValue(pvresource.Type)
	plan.Name = types.StringValue(pvresource.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *pvresourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pvresourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.Resources.DeleteResourceByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ProVision Resource",
			"Could not delete ProVision Resource, unexpected error: "+err.Error(),
		)
		return
	}
}
