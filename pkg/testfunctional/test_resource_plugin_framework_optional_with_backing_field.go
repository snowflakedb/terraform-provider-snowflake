package testfunctional

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigure = &OptionalWithBackingFieldResource{}

func NewOptionalWithBackingFieldResource() resource.Resource {
	return &OptionalWithBackingFieldResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[OptionalWithBackingFieldOpts]("optional_with_backing_field"),
	}
}

type OptionalWithBackingFieldResource struct {
	common.HttpServerEmbeddable[OptionalWithBackingFieldOpts]
}

type optionalWithBackingFieldResourceModelV0 struct {
	Name                    types.String `tfsdk:"name"`
	StringValue             types.String `tfsdk:"string_value"`
	StringValueBackingField types.String `tfsdk:"string_value_backing_field"`
	Id                      types.String `tfsdk:"id"`
}

type OptionalWithBackingFieldOpts struct {
	StringValue *string
}

func (r *OptionalWithBackingFieldResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_optional_with_backing_field"
}

func (r *OptionalWithBackingFieldResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			"string_value": schema.StringAttribute{
				Description: "String value.",
				Optional:    true,
			},
			"string_value_backing_field": schema.StringAttribute{
				Description: "String value backing field.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *OptionalWithBackingFieldResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *OptionalWithBackingFieldResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *optionalWithBackingFieldResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &OptionalWithBackingFieldOpts{}
	stringAttributeCreate(data.StringValue, &opts.StringValue)

	response.Diagnostics.Append(r.create(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.readAfterCreate(data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *OptionalWithBackingFieldResource) create(opts *OptionalWithBackingFieldOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *OptionalWithBackingFieldResource) readAfterCreate(data *optionalWithBackingFieldResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			data.StringValueBackingField = types.StringValue(*opts.StringValue)
		}
	}
	return diags
}

func (r *OptionalWithBackingFieldResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *optionalWithBackingFieldResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *OptionalWithBackingFieldResource) read(data *optionalWithBackingFieldResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			newValue := *opts.StringValue
			if newValue != data.StringValueBackingField.ValueString() {
				data.StringValue = types.StringValue(newValue)
			}
			data.StringValueBackingField = types.StringValue(newValue)
		}
	}
	return diags
}

func (r *OptionalWithBackingFieldResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *OptionalWithBackingFieldResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
