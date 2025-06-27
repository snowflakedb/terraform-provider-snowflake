package testacc

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewSomeResource() resource.Resource {
	return &SomeResource{}
}

type SomeResource struct{}

type someResourceModelV0 struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

func (r *SomeResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	// TODO [mux-PR]: add method for this logic
	response.TypeName = request.ProviderTypeName + "_some"
}

func (r *SomeResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this example resource.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this example resource.",
				PlanModifiers: []planmodifier.String{
					// TODO [mux-PR]: how it behaves with renames?
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SomeResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *someResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *SomeResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// TODO [mux-PR]: implement
}

func (r *SomeResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// TODO [mux-PR]: implement
}

func (r *SomeResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// TODO [mux-PR]: implement
}
