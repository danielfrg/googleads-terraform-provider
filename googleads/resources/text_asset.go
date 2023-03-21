package resources

import (
	"context"
	"fmt"

	"terraform-provider-googleads/googleads/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shenzhencenter/google-ads-pb/common"
	"github.com/shenzhencenter/google-ads-pb/enums"
	"github.com/shenzhencenter/google-ads-pb/resources"
	"github.com/shenzhencenter/google-ads-pb/services"
	"google.golang.org/grpc/status"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &textAssetResource{}
	_ resource.ResourceWithConfigure = &textAssetResource{}
)

// NewTextAssetResource is a helper function to simplify the provider implementation.
func NewTextAssetResource() resource.Resource {
	return &textAssetResource{}
}

// textAssetResource is the resource implementation.
type textAssetResource struct {
	client *client.GoogleAdsClient
}

type textAssetResourceModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Text         types.String `tfsdk:"text"`
}

// Metadata returns the resource type name.
func (r *textAssetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_text_asset"
}

// Schema defines the schema for the resource.
func (r *textAssetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Computed: true,
			},
			"text": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *textAssetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.GoogleAdsClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *textAssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "TextAsset: Create")

	// Retrieve values from plan
	var plan textAssetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	text := plan.Text.ValueString()

	// Generate API request from plan
	assetService := services.NewAssetServiceClient(&r.client.Connection)

	assetOperation := &services.AssetOperation{
		Operation: &services.AssetOperation_Create{Create: &resources.Asset{
			// Name:      &assetName,
			Type: enums.AssetTypeEnum_TEXT,
			AssetData: &resources.Asset_TextAsset{TextAsset: &common.TextAsset{
				Text: &text,
			}}},
		},
	}

	mutateRequest := &services.MutateAssetsRequest{
		CustomerId: r.client.CustomerId,
		Operations: []*services.AssetOperation{assetOperation},
	}

	response, err := assetService.MutateAssets(r.client.Context, mutateRequest)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			tflog.Info(ctx, fmt.Sprintf("%s %s %s %s", e.Code(), e.Message(), e.Details(), e.Err()))

		} else {
			fmt.Printf("not able to parse error returned %v", err)
		}
	}

	// Map response body to schema and populate Computed attribute values
	resource_name := response.Results[0].ResourceName
	tflog.Info(ctx, "Created Image Asset", map[string]any{"resource_name": resource_name})

	plan.ResourceName = types.StringValue(resource_name)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

const GAQL_GetTextAssetsByRN = `SELECT asset.resource_name, asset.text_asset.text, asset.name FROM asset WHERE asset.resource_name = '%s'`

// Read refreshes the Terraform state with the latest data.
func (r *textAssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "TextAsset: Read")

	// Get current state
	var state textAssetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "TextAsset: Read", map[string]any{"resource_name": state.ResourceName.ValueString()})

	// Get refreshed order value from the API
	request := services.SearchGoogleAdsRequest{
		CustomerId: r.client.CustomerId,
		Query:      fmt.Sprintf(GAQL_GetTextAssetsByRN, state.ResourceName.ValueString()),
	}
	response, err := services.NewGoogleAdsServiceClient(&r.client.Connection).Search(r.client.Context, &request)
	if err != nil {
		panic(err)
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
		state.ResourceName = types.StringValue(resource.Asset.GetResourceName())
		state.Text = types.StringValue(*resource.Asset.GetTextAsset().Text)
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
func (r *textAssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "TextAsset: Update")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *textAssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
