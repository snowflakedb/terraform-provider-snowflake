package testfunctional

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/tmpplanmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigure = &OptionalComputedResource{}

func NewOptionalComputedResource() resource.Resource {
	return &OptionalComputedResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[OptionalComputedOpts]("optional_computed_handling"),
	}
}

type OptionalComputedResource struct {
	common.HttpServerEmbeddable[OptionalComputedOpts]
}

type optionalComputedResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	StringValue types.String `tfsdk:"string_value"`
	Id          types.String `tfsdk:"id"`
}

type OptionalComputedOpts struct {
	StringValue *string
}

func (r *OptionalComputedResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_optional_computed"
}

func (r *OptionalComputedResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					tmpplanmodifiers.OptionalComputedString(),
				},
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

func (r *OptionalComputedResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *OptionalComputedResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *optionalComputedResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &OptionalComputedOpts{}
	stringAttributeCreate(data.StringValue, &opts.StringValue)

	response.Diagnostics.Append(r.create(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.readCreateUpdate(data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *OptionalComputedResource) create(opts *OptionalComputedOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *OptionalComputedResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *optionalComputedResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *OptionalComputedResource) read(data *optionalComputedResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			if data.StringValue.ValueString() != *opts.StringValue {
				data.StringValue = types.StringNull()
			}
		} else {
			data.StringValue = types.StringNull()
		}
	}
	return diags
}

func (r *OptionalComputedResource) readCreateUpdate(data *optionalComputedResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else if opts.StringValue != nil {
		data.StringValue = types.StringValue(*opts.StringValue)
	}
	return diags
}

func (r *OptionalComputedResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *optionalComputedResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &OptionalComputedOpts{}
	stringAttributeUpdate(plan.StringValue, state.StringValue, &opts.StringValue, &opts.StringValue)

	response.Diagnostics.Append(r.update(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.readCreateUpdate(plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *OptionalComputedResource) update(opts *OptionalComputedOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *OptionalComputedResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
