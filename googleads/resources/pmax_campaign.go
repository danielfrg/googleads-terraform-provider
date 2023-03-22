package resources

import (
	"context"
	"fmt"
	"strings"

	"terraform-provider-googleads/googleads/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shenzhencenter/google-ads-pb/common"
	"github.com/shenzhencenter/google-ads-pb/enums"
	"github.com/shenzhencenter/google-ads-pb/resources"
	"github.com/shenzhencenter/google-ads-pb/services"
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
	ResourceName           types.String `tfsdk:"resource_name"`
	AssetGroupResourceName types.String `tfsdk:"asset_group_resource_name"`
	Name                   types.String `tfsdk:"name"`
	Status                 types.String `tfsdk:"status"`
	Budget                 types.String `tfsdk:"budget"`
	TargetRoas             types.Number `tfsdk:"target_roas"`
	Headlines              types.List   `tfsdk:"headlines"`
	LongHeadLines          types.List   `tfsdk:"long_headlines"`
	Descriptions           types.List   `tfsdk:"descriptions"`
	BusinessName           types.String `tfsdk:"business_name"`
	MarketingImages        types.List   `tfsdk:"marketing_images"`
	LogoImages             types.List   `tfsdk:"logo_images"`
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
			"name": schema.StringAttribute{
				Required: true,
			},
			"status": schema.StringAttribute{
				Required: true,
			},
			"budget": schema.StringAttribute{
				Required: true,
			},
			"target_roas": schema.NumberAttribute{
				Required: true,
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
	var plan pMaxCampaignResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan

	name := plan.Name.ValueString()
	budgetRN := plan.Budget.ValueString()
	targetRoas, _ := plan.TargetRoas.ValueBigFloat().Float64()

	assetGroupRN := fmt.Sprintf("customers/%s/assetGroups/%d", r.client.CustomerId, -1)
	campaignRN := fmt.Sprintf("customers/%s/campaigns/%d", r.client.CustomerId, -2)

	tflog.Info(ctx, "TEMP", map[string]any{"resource_name": assetGroupRN})
	tflog.Info(ctx, "TEMP", map[string]any{"resource_name": campaignRN})

	operations := []*services.MutateOperation{}

	op_AssetGroup := &services.MutateOperation{
		Operation: &services.MutateOperation_AssetGroupOperation{
			AssetGroupOperation: &services.AssetGroupOperation{
				Operation: &services.AssetGroupOperation_Create{Create: &resources.AssetGroup{
					ResourceName: assetGroupRN,
					Name:         name,
					Campaign:     campaignRN,
				},
				},
			},
		},
	}
	operations = append(operations, op_AssetGroup)

	for _, headlineValue := range plan.Headlines.Elements() {
		headlineRN := headlineValue.String()
		op_LinkAsset := &services.MutateOperation{
			Operation: &services.MutateOperation_AssetGroupAssetOperation{
				AssetGroupAssetOperation: &services.AssetGroupAssetOperation{
					Operation: &services.AssetGroupAssetOperation_Create{Create: &resources.AssetGroupAsset{
						AssetGroup: assetGroupRN,
						Asset:      strings.Trim(headlineRN, "\""),
						FieldType:  enums.AssetFieldTypeEnum_HEADLINE,
					}},
				},
			},
		}

		operations = append(operations, op_LinkAsset)
	}

	op_Campaign := &services.MutateOperation{
		Operation: &services.MutateOperation_CampaignOperation{
			CampaignOperation: &services.CampaignOperation{
				Operation: &services.CampaignOperation_Create{Create: &resources.Campaign{
					ResourceName:           campaignRN,
					Name:                   &name,
					CampaignBudget:         &budgetRN,
					Status:                 enums.CampaignStatusEnum_PAUSED,
					AdvertisingChannelType: enums.AdvertisingChannelTypeEnum_PERFORMANCE_MAX,
					BiddingStrategyType:    enums.BiddingStrategyTypeEnum_MAXIMIZE_CONVERSION_VALUE,
					CampaignBiddingStrategy: &resources.Campaign_MaximizeConversionValue{
						MaximizeConversionValue: &common.MaximizeConversionValue{
							TargetRoas: targetRoas,
						},
					},
				},
				},
			},
		},
	}
	operations = append(operations, op_Campaign)

	mutateRequest := &services.MutateGoogleAdsRequest{
		CustomerId:       r.client.CustomerId,
		MutateOperations: operations,
	}

	client := services.NewGoogleAdsServiceClient(&r.client.Connection)
	response, err := client.Mutate(r.client.Context, mutateRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating PMax Campaign",
			ParseClientError(err))
		return
	}

	// Map response body to schema and populate Computed attribute values
	resource_name := response.MutateOperationResponses[0].Response.(*services.MutateOperationResponse_AssetGroupResult).AssetGroupResult.ResourceName

	resource_name_2 := response.MutateOperationResponses[0].Response.(*services.MutateOperationResponse_CampaignResult).CampaignResult.ResourceName

	tflog.Info(ctx, "Created PMax Campaign", map[string]any{"resource_name": resource_name})

	plan.ResourceName = types.StringValue(resource_name_2)
	plan.AssetGroupResourceName = types.StringValue(resource_name)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

const GAQL_GetPMaxCampaignByRN = `SELECT campaign.resource_name FROM campaign WHERE campaign.resource_name = '%s'`

// Read refreshes the Terraform state with the latest data.
func (r *pMaxCampaignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pMaxCampaignResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "PMaxCampaign: Read", map[string]any{"resource_name": state.ResourceName.ValueString()})

	// Get refreshed order value from the API
	request := services.SearchGoogleAdsRequest{
		CustomerId: r.client.CustomerId,
		Query:      fmt.Sprintf(GAQL_GetPMaxCampaignByRN, state.ResourceName.ValueString()),
	}
	response, err := services.NewGoogleAdsServiceClient(&r.client.Connection).Search(r.client.Context, &request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Budget",
			ParseClientError(err))
		return
	}

	// Overwrite values with refreshed state
	if len(response.Results) == 0 {
		// Treat empty response as resource not found
		resp.State.RemoveResource(ctx)
		return
	}
	if len(response.Results) > 1 {
		// TODO: Handle multiple results
		panic("Multiple results returned for resource name: " + state.ResourceName.ValueString())
	}
	for _, resource := range response.Results {
		state.ResourceName = types.StringValue(resource.Campaign.ResourceName)
		break
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *pMaxCampaignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "PMaxCampaign: Update")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *pMaxCampaignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "PMaxCampaign: Delete")

	// Get current state
	var state pMaxCampaignResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan
	client := services.NewCampaignServiceClient(&r.client.Connection)

	resourceName := state.ResourceName.ValueString()

	op := &services.CampaignOperation{
		Operation: &services.CampaignOperation_Remove{Remove: resourceName},
	}

	tflog.Info(ctx, "CustomerId", map[string]any{"CustomerId": r.client.CustomerId})
	mutateRequest := &services.MutateCampaignsRequest{
		CustomerId: r.client.CustomerId,
		Operations: []*services.CampaignOperation{op},
	}

	response, err := client.MutateCampaigns(r.client.Context, mutateRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing Campaign",
			ParseClientError(err))
		return
	}

	resource_name := response.Results[0].ResourceName
	tflog.Info(ctx, "Removed PMax Campaign", map[string]any{"resource_name": resource_name})

}
