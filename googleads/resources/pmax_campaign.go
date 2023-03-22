package resources

import (
	"context"

	"terraform-provider-googleads/googleads/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &pMaxCampaignResource{}
	_ resource.ResourceWithConfigure = &pMaxCampaignResource{}
)

// NewPMaxCampaignResource is a helper function to simplify the provider implementation.
func NewPMaxCampaignResource() resource.Resource {
	return &pMaxCampaignResource{}
}

// pMaxCampaignResource is the resource implementation.
type pMaxCampaignResource struct {
	client *client.GoogleAdsClient
}

type pMaxCampaignResourceModel struct {
	Headlines       types.List   `tfsdk:"headlines"`
	LongHeadLines   types.String `tfsdk:"long_headlines"`
	Descriptions    types.List   `tfsdk:"descriptions"`
	BusinessName    types.String `tfsdk:"business_name"`
	MarketingImages types.Object `tfsdk:"marketing_images"`
	LogoImages      types.Object `tfsdk:"logo_images"`
}

// Metadata returns the resource type name.
func (r *pMaxCampaignResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pmax_campaign"
}

// Schema defines the schema for the resource.
func (r *pMaxCampaignResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Computed: true,
			},
			"asset_group_resource_name": schema.StringAttribute{
				Computed: true,
			},
			"headlines": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"long_headlines": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"descriptions": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"business_name": schema.StringAttribute{
				Required: true,
			},
			"marketing_images": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"logo_images": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *pMaxCampaignResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.GoogleAdsClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *pMaxCampaignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "PMaxCampaign: Create")

	// Retrieve values from plan
	var plan textAssetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

const GAQL_GetPMaxCampaignByRN = `SELECT campaign.resource_name FROM campaign WHERE asset.resource_name = '%s'`

// Read refreshes the Terraform state with the latest data.
func (r *pMaxCampaignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "PMaxCampaign: Read")

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *pMaxCampaignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "PMaxCampaign: Update")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *pMaxCampaignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "PMaxCampaign: Delete")
}
