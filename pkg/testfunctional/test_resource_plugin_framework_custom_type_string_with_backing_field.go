package testfunctional

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customplanmodifiers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigure = &StringWithBackingFieldResource{}

func NewStringWithBackingFieldResource() resource.Resource {
	return &StringWithBackingFieldResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[StringWithBackingFieldOpts]("string_with_backing_field"),
	}
}

type StringWithBackingFieldResource struct {
	common.HttpServerEmbeddable[StringWithBackingFieldOpts]
}

type stringWithBackingFieldResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	StringValue types.String `tfsdk:"string_value"`
	Id          types.String `tfsdk:"id"`

	common.ActionsLogEmbeddable
}

type StringWithBackingFieldOpts struct {
	StringValue *string
}

func (r *StringWithBackingFieldResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_string_with_backing_field"
}

func (r *StringWithBackingFieldResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			"string_value": schema.StringAttribute{
				CustomType:  customtypes.StringWithBackingFieldType{},
				Description: "String value.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					customplanmodifiers.NewStringWithBackingFieldModifier(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			common.ActionsLogPropertyName: common.GetActionsLogSchema(),
		},
	}
}

func (r *StringWithBackingFieldResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *StringWithBackingFieldResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *stringWithBackingFieldResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &StringWithBackingFieldOpts{}
	stringAttributeCreate(data.StringValue, &opts.StringValue)

	response.Diagnostics.Append(r.create(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *StringWithBackingFieldResource) create(opts *StringWithBackingFieldOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *StringWithBackingFieldResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *stringWithBackingFieldResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.readStringWithBackingFieldResource(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *StringWithBackingFieldResource) readStringWithBackingFieldResource(data *stringWithBackingFieldResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			data.StringValue = types.StringValue(*opts.StringValue)
		}
	}
	return diags
}

func (r *StringWithBackingFieldResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *StringWithBackingFieldResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
