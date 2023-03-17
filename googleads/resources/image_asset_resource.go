package resources

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

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
	Name         types.String `tfsdk:"name"`
	Path         types.String `tfsdk:"path"`
	ResourceName types.String `tfsdk:"resource_name"`
	LastUpdated  types.String `tfsdk:"last_updated"`
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
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
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
	tflog.Info(ctx, "Creating Image Asset")

	// Retrieve values from plan
	var plan imageAssetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePath := plan.Path.ValueString()
	tflog.Info(ctx, "Reading Image", map[string]any{"path": filePath})

	image, err := getImageFromFilePath(filePath)
	// TODO: Handle path doesn't exist
	if err != nil {
		// TODO: handle error
	}
	tflog.Info(ctx, "Image", map[string]any{"image": image})

	// Generate API request from plan
	assetService := services.NewAssetServiceClient(&r.client.Connection)

	url := "https://gaagl.page.link/Eit5"

	assetName := plan.Path.ValueString()
	assetOperation := &services.AssetOperation{
		Operation: &services.AssetOperation_Create{Create: &resources.Asset{
			Name: &assetName,
			Type: enums.AssetTypeEnum_IMAGE,
			AssetData: &resources.Asset_ImageAsset{ImageAsset: &common.ImageAsset{
				Data:     *image.Data,
				FileSize: &image.Size,
				MimeType: enums.MimeTypeEnum_IMAGE_JPEG,
				FullSize: &common.ImageDimension{
					WidthPixels:  &image.Width,
					HeightPixels: &image.Height,
					Url:          &url,
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
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *imageAssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *imageAssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *imageAssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

type ImageInfo struct {
	Data   *[]byte
	Type   string
	Size   int64
	Width  int64
	Height int64
}

func getImageFromFilePath(filePath string) (ImageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ImageInfo{}, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	filetype := http.DetectContentType(bytes)

	image := ImageInfo{
		Data:   &bytes,
		Type:   filetype,
		Size:   size,
		Height: 315,
		Width:  600,
	}

	return image, err
}
