package testfunctional

import (
	"context"
	"math/rand/v2"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &staleErrorReproResource{}

func NewStaleErrorReproResource() resource.Resource {
	return &staleErrorReproResource{}
}

type staleErrorReproResource struct{}

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
			// random_int changes on every Read call, so any two consecutive reads produce
			// different state — guaranteeing the Refresh() inserted between CreatePlan and
			// Apply by terraform-plugin-testing v1.14.0 always writes a new state serial,
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
	data.RandomInt = types.Int64Value(rand.Int64()) // #nosec G404
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staleErrorReproResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data staleErrorReproResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.RandomInt = types.Int64Value(rand.Int64()) // #nosec G404
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staleErrorReproResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *staleErrorReproResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
