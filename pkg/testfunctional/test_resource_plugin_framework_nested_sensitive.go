package testfunctional

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewNestedSensitiveResource() resource.Resource {
	return &NestedSensitiveResource{}
}

type NestedSensitiveResource struct{}

type nestedSensitiveResourceModelV0 struct {
	Name   types.String `tfsdk:"name"`
	Id     types.String `tfsdk:"id"`
	Output types.List   `tfsdk:"output"`
}

func (r *NestedSensitiveResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nested_sensitive"
}

func (r *NestedSensitiveResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"output": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_sensitive": schema.StringAttribute{
							Computed:  true,
							Sensitive: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (r *NestedSensitiveResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *nestedSensitiveResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())
	var diags diag.Diagnostics
	data.Output, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{"string_sensitive": types.StringType}}, []attr.Value{
		types.ObjectValueMust(map[string]attr.Type{"string_sensitive": types.StringType}, map[string]attr.Value{
			"string_sensitive": types.StringValue("SECRET"),
		}),
	})
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *NestedSensitiveResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *NestedSensitiveResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *NestedSensitiveResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
