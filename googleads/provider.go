package googleads

import (
	"context"
	"os"

	"terraform-provider-googleads/googleads/client"
	"terraform-provider-googleads/googleads/resources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &googleadsProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &googleadsProvider{}
}

// googleadsProvider is the provider implementation.
type googleadsProvider struct{}

// googleadsProviderModel maps provider schema data to a Go type.
type googleadsProviderModel struct {
	DeveloperToken  types.String `tfsdk:"developer_token"`
	AccessToken     types.String `tfsdk:"access_token"`
	CustomerId      types.String `tfsdk:"customer_id"`
	LoginCustomerId types.String `tfsdk:"login_customer_id"`
}

// Metadata returns the provider type name.
func (p *googleadsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "googleads"
}

// Schema defines the provider-level schema for configuration data.
func (p *googleadsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"developer_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"access_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"customer_id": schema.StringAttribute{
				Required: true,
			},
			"login_customer_id": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Configure prepares a Google Ads API client for data sources and resources.
func (p *googleadsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring googleads client")

	// Retrieve provider data from configuration
	var config googleadsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	dev_token := os.Getenv("GOOGLEADS_DEVELOPER_TOKEN")
	access_token := os.Getenv("GOOGLEADS_ACCESS_TOKEN")
	customer_id := os.Getenv("GOOGLEADS_CUSTOMER_ID")
	login_customer_id := os.Getenv("GOOGLEADS_LOGIN_CUSTOMER_ID")

	if !config.DeveloperToken.IsNull() {
		dev_token = config.DeveloperToken.ValueString()
	}

	if !config.AccessToken.IsNull() {
		access_token = config.AccessToken.ValueString()
	}

	if !config.CustomerId.IsNull() {
		customer_id = config.CustomerId.ValueString()
	}

	if !config.LoginCustomerId.IsNull() {
		login_customer_id = config.LoginCustomerId.ValueString()
	}

	ctx = tflog.SetField(ctx, "googleads_dev_token", dev_token)
	// ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "googleads_dev_token")
	ctx = tflog.SetField(ctx, "googleads_access_token", access_token)
	// ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "googleads_access_token")
	ctx = tflog.SetField(ctx, "googleads_customer_id", customer_id)
	ctx = tflog.SetField(ctx, "googleads_login_customer_id", login_customer_id)

	tflog.Debug(ctx, "Creating Google Ads API client")

	// TODO: Error handling for missing configuration values

	// Create client
	client, err := client.NewGoogleAdsClient(dev_token, access_token, customer_id, login_customer_id)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Google Ads API Client",
			"An unexpected error occurred when creating the Google Ads API client:\n"+err.Error(),
		)
		return
	}

	// Make the client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Google Ads API client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *googleadsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCoffeesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *googleadsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewImageAssetResource,
		resources.NewTextAssetResource,
		resources.NewBudgetResource,
	}
}
