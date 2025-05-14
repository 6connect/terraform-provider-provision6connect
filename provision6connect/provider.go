package provision6connect

import (
	"context"
	"os"

	provisionclient "github.com/6connect/golangclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &provision6connectProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &provision6connectProvider{}
}

// provision6connectProviderModel maps provider schema data to a Go type.
type provision6connectProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// provision6connectProvider is the provider implementation.
type provision6connectProvider struct{}

// Metadata returns the provider type name.
func (p *provision6connectProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "provision6connect"
}

// Schema defines the provider-level schema for configuration data.
func (p *provision6connectProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with 6connect ProVision.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for 6connect ProVision. May also be provided via PROVISION_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for 6connect ProVision. May also be provided via PROVISION_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for 6connect ProVision. May also be provided via PROVISION_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *provision6connectProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config provision6connectProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown 6connect ProVision Host",
			"The provider cannot create the 6connect ProVision client as there is an unknown configuration value for the 6connect ProVision host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PROVISION_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown 6connect ProVision Username",
			"The provider cannot create the 6connect ProVision client as there is an unknown configuration value for the 6connect ProVision username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PROVISION_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown 6connect ProVision Password",
			"The provider cannot create the 6connect ProVision client as there is an unknown configuration value for the 6connect ProVision password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PROVISION_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("PROVISION_HOST")
	username := os.Getenv("PROVISION_USERNAME")
	password := os.Getenv("PROVISION_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing 6connect ProVision Host",
			"The provider cannot create the 6connect ProVision client as there is a missing or empty value for the 6connect ProVision host. "+
				"Set the host value in the configuration or use the PROVISION_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing 6connect ProVision Username",
			"The provider cannot create the 6connect ProVision client as there is a missing or empty value for the 6connect ProVision username. "+
				"Set the username value in the configuration or use the PROVISION_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing 6connect ProVision Password",
			"The provider cannot create the 6connect ProVision client as there is a missing or empty value for the 6connect ProVision password. "+
				"Set the password value in the configuration or use the PROVISION_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new client using the configuration values
	client, err := provisionclient.NewClient(host, username, password, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create 6connect ProVision Client",
			"An unexpected error occurred when creating the 6connect ProVision client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"6connect ProVision Client Error: "+err.Error(),
		)
		return
	}

	// Make the 6connect ProVision client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *provision6connectProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewResourcesDataSource,
		NewNetblocksDataSource,
		NewFirstavailableipDataSource,
		NewDNSpushDataSource,
		NewDNSpushstatusDataSource,
		NewDHCPpushDataSource,
		NewDHCPpushstatusDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *provision6connectProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPVresourceResource,
		NewIPAMsmartassignResource,
		NewIPAMdirectassignResource,
		NewIPAMnetblockResource,
		NewDNSrecordResource,
		NewDNSzoneResource,
	}
}
