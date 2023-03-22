package resources

import (
	"context"
	"fmt"

	"terraform-provider-googleads/googleads/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shenzhencenter/google-ads-pb/enums"
	"github.com/shenzhencenter/google-ads-pb/resources"
	"github.com/shenzhencenter/google-ads-pb/services"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &budgetResource{}
	_ resource.ResourceWithConfigure = &budgetResource{}
)

// NewBudgetResource is a helper function to simplify the provider implementation.
func NewBudgetResource() resource.Resource {
	return &budgetResource{}
}

// budgetResource is the resource implementation.
type budgetResource struct {
	client *client.GoogleAdsClient
}

type budgetResourceModel struct {
	ResourceName     types.String `tfsdk:"resource_name"`
	Name             types.String `tfsdk:"name"`
	AmountMicros     types.Number `tfsdk:"amount_micros"`
	DeliveryMethod   types.String `tfsdk:"delivery_method"`
	ExplicitlyShared types.Bool   `tfsdk:"explicitly_shared"`
}

// Metadata returns the resource type name.
func (r *budgetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_budget"
}

// Schema defines the schema for the resource.
func (r *budgetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"amount_micros": schema.NumberAttribute{
				Required: true,
			},
			"delivery_method": schema.StringAttribute{
				Required: true,
			},
			"explicitly_shared": schema.BoolAttribute{
				Required: true,
			},
		},
	}
}

func (r *budgetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.GoogleAdsClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *budgetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Budget: Create")

	// Retrieve values from plan
	var plan budgetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request from plan
	client := services.NewCampaignBudgetServiceClient(&r.client.Connection)

	name := plan.Name.ValueString()
	micros, _ := plan.AmountMicros.ValueBigFloat().Int64()
	shared := plan.ExplicitlyShared.ValueBool()

	op := &services.CampaignBudgetOperation{
		Operation: &services.CampaignBudgetOperation_Create{Create: &resources.CampaignBudget{
			Name:             &name,
			AmountMicros:     &micros,
			DeliveryMethod:   enums.BudgetDeliveryMethodEnum_STANDARD,
			ExplicitlyShared: &shared,
		}},
	}

	mutateRequest := &services.MutateCampaignBudgetsRequest{
		CustomerId: r.client.CustomerId,
		Operations: []*services.CampaignBudgetOperation{op},
	}

	response, err := client.MutateCampaignBudgets(r.client.Context, mutateRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Budget",
			ParseClientError(err))
		return
	}

	// Map response body to schema and populate Computed attribute values
	resource_name := response.Results[0].ResourceName
	tflog.Info(ctx, "Created Budget", map[string]any{"resource_name": resource_name})

	plan.ResourceName = types.StringValue(resource_name)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

const GAQL_GetBudgetsByRN = `SELECT campaign_budget.resource_name, campaign_budget.amount_micros FROM campaign_budget WHERE campaign_budget.resource_name = '%s'`

// Read refreshes the Terraform state with the latest data.
func (r *budgetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Budget: Read")

	// Get current state
	var state budgetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Budget: Read", map[string]any{"resource_name": state.ResourceName.ValueString()})

	// Get refreshed order value from the API
	request := services.SearchGoogleAdsRequest{
		CustomerId: r.client.CustomerId,
		Query:      fmt.Sprintf(GAQL_GetBudgetsByRN, state.ResourceName.ValueString()),
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
		state.ResourceName = types.StringValue(resource.CampaignBudget.ResourceName)
		state.AmountMicros = types.NumberValue(ToBigFloat(*resource.CampaignBudget.AmountMicros))
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
func (r *budgetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Budget: Update")
	tflog.Info(ctx, "Assets are immutable and all fields force a new resource")

	resp.Diagnostics.AddError(
		"Google Ads Assets are immutable",
		"Any fields change should force a new resource. This is an error in the provider.")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *budgetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Budget: Delete")
	tflog.Info(ctx, "Assets are immutable, acting as if delete was successful")
}
