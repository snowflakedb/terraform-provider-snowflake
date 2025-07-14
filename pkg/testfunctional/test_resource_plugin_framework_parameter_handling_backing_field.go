package testfunctional

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.ResourceWithConfigure = &ParameterHandlingBackingFieldResource{}

func NewParameterHandlingBackingFieldResource() resource.Resource {
	return &ParameterHandlingBackingFieldResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[ParameterHandlingBackingFieldOpts]("parameter_handling_backing_field"),
	}
}

type ParameterHandlingBackingFieldResource struct {
	common.HttpServerEmbeddable[ParameterHandlingBackingFieldOpts]
}

type parameterHandlingBackingFieldResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	StringValue types.String `tfsdk:"string_value"`
	Id          types.String `tfsdk:"id"`

	StringValueBackingField types.Object `tfsdk:"string_value_backing_field"`

	common.ActionsLogEmbeddable
}

type ParameterBackingField struct {
	Value types.String `tfsdk:"value"`
	Level types.String `tfsdk:"level"`
}

type ParameterHandlingBackingFieldOpts struct {
	StringValue *string
	Level       string
}

func (r *ParameterHandlingBackingFieldResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_parameter_handling_backing_field"
}

func (r *ParameterHandlingBackingFieldResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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
			"string_value_backing_field": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"value": types.StringType,
					"level": types.StringType,
				},
				Description: "Parameter backing field.",
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

func (r *ParameterHandlingBackingFieldResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			// TODO: set conditionally based on level
			response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("string_value"), *opts.StringValue)...)
		}
	}
}

func (r *ParameterHandlingBackingFieldResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *parameterHandlingBackingFieldResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &ParameterHandlingBackingFieldOpts{}
	stringAttributeCreate(data.StringValue, &opts.StringValue)

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

func (r *ParameterHandlingBackingFieldResource) setCreateActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *ParameterHandlingBackingFieldOpts, data *parameterHandlingBackingFieldResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "string_value", *opts.StringValue))
		}
		return actions
	})...)
}

func (r *ParameterHandlingBackingFieldResource) create(opts *ParameterHandlingBackingFieldOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingBackingFieldResource) readAfterCreateOrUpdate(data *parameterHandlingBackingFieldResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			elementTypes := map[string]attr.Type{
				"value": types.StringType,
				"level": types.StringType,
			}
			elements := map[string]attr.Value{
				"value": types.StringValue(*opts.StringValue),
				"level": types.StringValue(opts.Level),
			}
			objectValue, d := types.ObjectValue(elementTypes, elements)
			if d.HasError() {
				diags.Append(d...)
				return diags
			}
			data.StringValueBackingField = objectValue
		}
	}
	return diags
}

func (r *ParameterHandlingBackingFieldResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *parameterHandlingBackingFieldResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(ctx, data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *ParameterHandlingBackingFieldResource) read(ctx context.Context, data *parameterHandlingBackingFieldResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			newValue := *opts.StringValue
			newLevel := opts.Level

			if !data.StringValueBackingField.IsNull() {
				var param *ParameterBackingField
				diags.Append(data.StringValueBackingField.As(ctx, &param, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return diags
				}
				// if new value differs from the previous one
				if newValue != param.Value.ValueString() {
					data.StringValue = types.StringValue(newValue)
					// if new level is not object we should set null
				} else if newLevel != "OBJECT" {
					data.StringValue = types.StringNull()
				} else {
					data.StringValue = types.StringValue(newValue)
				}
			}

			elementTypes := map[string]attr.Type{
				"value": types.StringType,
				"level": types.StringType,
			}
			elements := map[string]attr.Value{
				"value": types.StringValue(newValue),
				"level": types.StringValue(newLevel),
			}
			objectValue, d := types.ObjectValue(elementTypes, elements)
			if d.HasError() {
				diags.Append(d...)
				return diags
			}
			data.StringValueBackingField = objectValue
		}
	}
	return diags
}

func (r *ParameterHandlingBackingFieldResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *parameterHandlingBackingFieldResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &ParameterHandlingBackingFieldOpts{}
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

func (r *ParameterHandlingBackingFieldResource) update(opts *ParameterHandlingBackingFieldOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not update resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingBackingFieldResource) setUpdateActionsOutput(ctx context.Context, response *resource.UpdateResponse, opts *ParameterHandlingBackingFieldOpts, plan *parameterHandlingBackingFieldResourceModelV0, state *parameterHandlingBackingFieldResourceModelV0) {
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

func (r *ParameterHandlingBackingFieldResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
