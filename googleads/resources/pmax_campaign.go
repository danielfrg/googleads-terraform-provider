package resources

// import (
// 	"context"
// 	"fmt"

// 	"terraform-provider-googleads/googleads/client"

// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// 	"github.com/shenzhencenter/google-ads-pb/services"
// )

// // Ensure the implementation satisfies the expected interfaces.
// var (
// 	_ resource.Resource              = &pMaxCampaignResource{}
// 	_ resource.ResourceWithConfigure = &pMaxCampaignResource{}
// )

// // NewPMaxCampaignResource is a helper function to simplify the provider implementation.
// func NewPMaxCampaignResource() resource.Resource {
// 	return &pMaxCampaignResource{}
// }

// // pMaxCampaignResource is the resource implementation.
// type pMaxCampaignResource struct {
// 	client *client.GoogleAdsClient
// }

// type pMaxCampaignResourceModel struct {
// 	Headlines       types.List   `tfsdk:"headlines"`
// 	LongHeadLines   types.String `tfsdk:"long_headlines"`
// 	Descriptions    types.List   `tfsdk:"descriptions"`
// 	BusinessName    types.String `tfsdk:"business_name"`
// 	MarketingImages types.Object `tfsdk:"marketing_images"`
// 	LogoImages      types.Object `tfsdk:"logo_images"`
// }

// // Metadata returns the resource type name.
// func (r *pMaxCampaignResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_pmax_campaign"
// }

// // Schema defines the schema for the resource.
// func (r *pMaxCampaignResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		Attributes: map[string]schema.Attribute{
// 			"resource_name": schema.StringAttribute{
// 				Computed: true,
// 			},
// 			"headlines": schema.ListNestedAttribute{
// 				Required: true,
// 				NestedObject: schema.NestedAttributeObject{
// 					// Required: true,
// 					Attributes: map[string]schema.Attribute{
// 						"resource_name": schema.StringAttribute{
// 							Computed: true,
// 						},
// 						"text": schema.StringAttribute{
// 							Required: true,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// func (r *pMaxCampaignResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	r.client = req.ProviderData.(*client.GoogleAdsClient)
// }

// // Create creates the resource and sets the initial Terraform state.
// func (r *pMaxCampaignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	tflog.Info(ctx, "PMaxCampaign: Create")

// }

// const GAQL_GetPMaxCampaignByRN = `SELECT asset.image_asset.file_size, asset.image_asset.full_size.height_pixels, asset.image_asset.full_size.url, asset.image_asset.full_size.width_pixels, asset.image_asset.mime_type, asset.name FROM asset WHERE asset.resource_name = '%s'`

// // Read refreshes the Terraform state with the latest data.
// func (r *pMaxCampaignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	tflog.Info(ctx, "PMaxCampaign: Read")

// 	// Get current state
// 	var state pMaxCampaignResourceModel
// 	diags := req.State.Get(ctx, &state)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	tflog.Info(ctx, "PMaxCampaign: Read", map[string]any{"resource_name": state.ResourceName.ValueString()})

// 	// Get refreshed order value from the API
// 	request := services.SearchGoogleAdsRequest{
// 		CustomerId: r.client.CustomerId,
// 		Query:      fmt.Sprintf(GAQL_GetAssetsByRN, state.ResourceName.ValueString()),
// 	}
// 	response, err := services.NewGoogleAdsServiceClient(&r.client.Connection).Search(r.client.Context, &request)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Overwrite values with refreshed state
// 	if len(response.Results) == 0 {
// 		// Treat empty response as resource not found
// 		resp.State.RemoveResource(ctx)
// 		return
// 	}
// 	if len(response.Results) > 1 {
// 		// TODO: Handle multiple results
// 		panic("Multiple results returned for resource name: " + state.ResourceName.ValueString())
// 	}
// 	for _, resource := range response.Results {
// 		state.ResourceName = types.StringValue(resource.Asset.GetResourceName())
// 		state.Name = types.StringValue(resource.Asset.GetName())
// 		break
// 	}

// 	// Set refreshed state
// 	diags = resp.State.Set(ctx, &state)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// }

// // Update updates the resource and sets the updated Terraform state on success.
// func (r *pMaxCampaignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	tflog.Info(ctx, "PMaxCampaign: Update")
// }

// // Delete deletes the resource and removes the Terraform state on success.
// func (r *pMaxCampaignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// }
