package testfunctional

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigure = &staleErrorReproResource{}

func NewStaleErrorReproResource() resource.Resource {
	return &staleErrorReproResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[StaleReproRead]("stale_error_repro"),
	}
}

type staleErrorReproResource struct {
	common.HttpServerEmbeddable[StaleReproRead]
}

// StaleReproRead is the response type returned by the HTTP server handler.
// The handler generates a new random value on every GET, so successive Read()
// calls always return different values — the condition that makes the stale
// plan error 100% reproducible in terraform-plugin-testing v1.14.0.
type StaleReproRead struct {
	RandomInt int64 `json:"random_int"`
}

type staleErrorReproResourceModel struct {
	Name      types.String `tfsdk:"name"`
	Id        types.String `tfsdk:"id"`
	RandomInt types.Int64  `tfsdk:"random_int"`
}

func (r *staleErrorReproResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stale_error_repro"
}

func (r *staleErrorReproResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// random_int is re-fetched from the HTTP server on every Read call.
			// The server generates a new value each time, so any two consecutive reads
			// produce different state — guaranteeing the Refresh() inserted between
			// CreatePlan and Apply by terraform-plugin-testing v1.14.0 (triggered when
			// a destroy step has a non-nil Check) always writes a new state serial,
			// making the plan stale on every single run.
			"random_int": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (r *staleErrorReproResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data staleErrorReproResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(sdk.NewAccountObjectIdentifier(data.Name.ValueString()).FullyQualifiedName())

	resp.Diagnostics.Append(r.readInto(&data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staleErrorReproResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data staleErrorReproResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readInto(&data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staleErrorReproResource) readInto(data *staleErrorReproResourceModel) diag.Diagnostics {
	diags := diag.Diagnostics{}
	result, err := r.Get()
	if err != nil {
		diags.AddError("Could not read resource state", err.Error())
		return diags
	}
	data.RandomInt = types.Int64Value(result.RandomInt)
	return diags
}

func (r *staleErrorReproResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *staleErrorReproResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
