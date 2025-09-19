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

	common.ActionsLogEmbeddable
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
			common.ActionsLogPropertyName: common.GetActionsLogSchema(),
		},
	}
}

func (r *OptionalWithBackingFieldResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
	} else if opts.StringValue != nil {
		response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("string_value"), *opts.StringValue)...)
	}
}

func (r *OptionalWithBackingFieldResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *optionalWithBackingFieldResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &OptionalWithBackingFieldOpts{}
	_ = StringAttributeCreate(data.StringValue, &opts.StringValue)

	r.setCreateActionsOutput(ctx, response, opts, data)

	response.Diagnostics.Append(r.create(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.readAfterCreateOrUpdate(data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *OptionalWithBackingFieldResource) setCreateActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *OptionalWithBackingFieldOpts, data *optionalWithBackingFieldResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "string_value", *opts.StringValue))
		}
		return actions
	})...)
}

func (r *OptionalWithBackingFieldResource) create(opts *OptionalWithBackingFieldOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *OptionalWithBackingFieldResource) readAfterCreateOrUpdate(data *optionalWithBackingFieldResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else if opts.StringValue != nil {
		data.StringValueBackingField = types.StringValue(*opts.StringValue)
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
	} else if opts.StringValue != nil {
		newValue := *opts.StringValue
		if newValue != data.StringValueBackingField.ValueString() {
			data.StringValue = types.StringValue(newValue)
		}
		data.StringValueBackingField = types.StringValue(newValue)
	}
	return diags
}

func (r *OptionalWithBackingFieldResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *optionalWithBackingFieldResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &OptionalWithBackingFieldOpts{}
	stringAttributeUpdate(plan.StringValue, state.StringValue, &opts.StringValue, &opts.StringValue)

	r.setUpdateActionsOutput(ctx, response, opts, plan, state)

	response.Diagnostics.Append(r.update(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.readAfterCreateOrUpdate(plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *OptionalWithBackingFieldResource) update(opts *OptionalWithBackingFieldOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not update resource", err.Error())
	}
	return diags
}

func (r *OptionalWithBackingFieldResource) setUpdateActionsOutput(ctx context.Context, response *resource.UpdateResponse, opts *OptionalWithBackingFieldOpts, plan *optionalWithBackingFieldResourceModelV0, state *optionalWithBackingFieldResourceModelV0) {
	plan.ActionsLogEmbeddable = state.ActionsLogEmbeddable
	response.Diagnostics.Append(common.AppendActions(ctx, &plan.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("UPDATE - SET", "string_value", *opts.StringValue))
		} else {
			actions = append(actions, common.ActionEntry("UPDATE - UNSET", "string_value", "nil"))
		}
		return actions
	})...)
}

func (r *OptionalWithBackingFieldResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
