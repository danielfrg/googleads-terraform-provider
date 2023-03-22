package resources

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"

	"terraform-provider-googleads/googleads/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shenzhencenter/google-ads-pb/common"
	"github.com/shenzhencenter/google-ads-pb/enums"
	"github.com/shenzhencenter/google-ads-pb/resources"
	"github.com/shenzhencenter/google-ads-pb/services"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &imageAssetResource{}
	_ resource.ResourceWithConfigure = &imageAssetResource{}
)

// NewImageAssetResource is a helper function to simplify the provider implementation.
func NewImageAssetResource() resource.Resource {
	return &imageAssetResource{}
}

// imageAssetResource is the resource implementation.
type imageAssetResource struct {
	client *client.GoogleAdsClient
}

type imageAssetResourceModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Name         types.String `tfsdk:"name"`
	Path         types.String `tfsdk:"path"`
	Url          types.String `tfsdk:"url"`
	Hash         types.String `tfsdk:"hash"`
}

// Metadata returns the resource type name.
func (r *imageAssetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image_asset"
}

// Schema defines the schema for the resource.
func (r *imageAssetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Optional: true,
			},
			"hash": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *imageAssetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.GoogleAdsClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *imageAssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "ImageAsset: Create")

	// Retrieve values from plan
	var plan imageAssetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePath := plan.Path.ValueString()
	// url := plan.Url.ValueString()

	image, err := GetImageFromFilePath(filePath)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading image",
			err.Error())
		return
	}

	// Generate API request from plan
	assetService := services.NewAssetServiceClient(&r.client.Connection)

	assetName := plan.Name.ValueString()
	mimeType := enums.MimeTypeEnum_IMAGE_JPEG

	assetOperation := &services.AssetOperation{
		Operation: &services.AssetOperation_Create{Create: &resources.Asset{
			Name: &assetName,
			Type: enums.AssetTypeEnum_IMAGE,
			AssetData: &resources.Asset_ImageAsset{ImageAsset: &common.ImageAsset{
				Data:     *image.Data,
				FileSize: &image.Size,
				MimeType: mimeType,
				FullSize: &common.ImageDimension{
					WidthPixels:  &image.Width,
					HeightPixels: &image.Height,
					// Url:          &url,
				},
			}}},
		},
	}

	mutateRequest := &services.MutateAssetsRequest{
		CustomerId: r.client.CustomerId,
		Operations: []*services.AssetOperation{assetOperation},
	}

	response, err := assetService.MutateAssets(r.client.Context, mutateRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ImageAsset",
			ParseClientError(err))
		return
	}

	// Map response body to schema and populate Computed attribute values
	resource_name := response.Results[0].ResourceName
	tflog.Info(ctx, "Created ImageAsset", map[string]any{"resource_name": resource_name})

	plan.ResourceName = types.StringValue(resource_name)
	plan.Hash = types.StringValue(image.Hash)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

const GAQL_GetImageAssetsByRN = `SELECT asset.resource_name, asset.image_asset.file_size, asset.image_asset.full_size.height_pixels, asset.image_asset.full_size.url, asset.image_asset.full_size.width_pixels, asset.image_asset.mime_type, asset.name FROM asset WHERE asset.resource_name = '%s'`

// Read refreshes the Terraform state with the latest data.
func (r *imageAssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "ImageAsset: Read")

	// Get current state
	var state imageAssetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ImageAsset: Read", map[string]any{"resource_name": state.ResourceName.ValueString()})

	// Get refreshed order value from the API
	request := services.SearchGoogleAdsRequest{
		CustomerId: r.client.CustomerId,
		Query:      fmt.Sprintf(GAQL_GetImageAssetsByRN, state.ResourceName.ValueString()),
	}
	response, err := services.NewGoogleAdsServiceClient(&r.client.Connection).Search(r.client.Context, &request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ImageAsset",
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
		state.ResourceName = types.StringValue(resource.Asset.GetResourceName())
		state.Name = types.StringValue(resource.Asset.GetName())
		// state.Url = types.StringValue(*resource.Asset.GetImageAsset().FullSize.Url)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *imageAssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "ImageAsset: Update")
	tflog.Info(ctx, "Assets are immutable and all fields force a new resource")

	resp.Diagnostics.AddError(
		"Google Ads Assets are immutable",
		"Any fields change should force a new resource. This is an error in the provider.")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *imageAssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "ImageAsset: Delete")
	tflog.Info(ctx, "Assets are immutable, acting as if delete was successful")
}

type ImageInfo struct {
	Data     *[]byte
	Hash     string
	Size     int64
	Height   int64
	Width    int64
	MimeType string
}

func GetImageFromFilePath(filePath string) (ImageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ImageInfo{}, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	i, err := file.Read(buffer)
	_ = i
	if err != nil {
		return ImageInfo{}, err
	}

	h := sha256.New()
	h.Write(buffer)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	mimeType := http.DetectContentType(buffer)

	image := ImageInfo{
		Data:     &buffer,
		Hash:     hash,
		MimeType: mimeType,
		Size:     size,
		Height:   315,
		Width:    600,
	}

	return image, err
}
