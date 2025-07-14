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

var _ resource.ResourceWithConfigure = &ParameterHandlingResourcePlanModifierResource{}

func NewParameterHandlingResourcePlanModifierResource() resource.Resource {
	return &ParameterHandlingResourcePlanModifierResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[ParameterHandlingResourcePlanModifierOpts]("parameter_handling_resource_plan_modifier"),
	}
}

type ParameterHandlingResourcePlanModifierResource struct {
	common.HttpServerEmbeddable[ParameterHandlingResourcePlanModifierOpts]
}

type parameterHandlingResourcePlanModifierResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	StringValue types.String `tfsdk:"string_value"`
	Id          types.String `tfsdk:"id"`

	common.ActionsLogEmbeddable
}

type ParameterHandlingResourcePlanModifierOpts struct {
	StringValue *string
	Level       string
}

func (r *ParameterHandlingResourcePlanModifierResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_parameter_handling_resource_plan_modifier"
}

func (r *ParameterHandlingResourcePlanModifierResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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

// ModifyPlan inlines resources.ParameterValueComputedIf which is our previous SDK implementation.
func (r *ParameterHandlingResourcePlanModifierResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	// Do nothing if there is no state (resource is being created).
	if request.State.Raw.IsNull() {
		return
	}

	// Do nothing if there is no plan
	if request.Plan.Raw.IsNull() {
		return
	}

	// Do nothing if there is no config
	if request.Config.Raw.IsNull() {
		return
	}

	var config, plan, state *parameterHandlingResourcePlanModifierResourceModelV0
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	foundParameter, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
		return
	}

	// For cases where currently set value (in the config) is equal to the parameter, but not set on the right level.
	// The parameter is set somewhere higher in the hierarchy, and we need to "forcefully" set the value to
	// perform the actual set (and set the parameter on the correct level).
	if !config.StringValue.IsNull() && foundParameter.Level != "OBJECT" && foundParameter.StringValue != nil && *foundParameter.StringValue == state.StringValue.ValueString() {
		plan.StringValue = types.StringUnknown()
		response.Diagnostics.Append(response.Plan.Set(ctx, &plan)...)
		return
	}

	// For all other cases, if a parameter is set in the configuration, we can ignore parts needed for Computed fields.
	if !config.StringValue.IsNull() {
		return
	}

	// If the configuration is not set, perform SetNewComputed for cases like:
	// 1. Check if the parameter value differs from the one saved in state (if they differ, we'll update the computed value).
	// 2. Check if the parameter is set on the object level (if so, it means that it was set externally, and we have to unset it).
	//if foundParameter.StringValue !=  || parameter.Level == objectParameterLevel {
	//	plan.StringValue = types.StringUnknown()
	//	return
	//}
}

func (r *ParameterHandlingResourcePlanModifierResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("string_value"), *opts.StringValue)...)
		}
	}
}

func (r *ParameterHandlingResourcePlanModifierResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *parameterHandlingResourcePlanModifierResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &ParameterHandlingResourcePlanModifierOpts{}
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

func (r *ParameterHandlingResourcePlanModifierResource) setCreateActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *ParameterHandlingResourcePlanModifierOpts, data *parameterHandlingResourcePlanModifierResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "string_value", *opts.StringValue))
		}
		return actions
	})...)
}

func (r *ParameterHandlingResourcePlanModifierResource) create(opts *ParameterHandlingResourcePlanModifierOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingResourcePlanModifierResource) readAfterCreateOrUpdate(data *parameterHandlingResourcePlanModifierResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			data.StringValue = types.StringValue(*opts.StringValue)
		} else {
			data.StringValue = types.StringNull()
		}
	}
	return diags
}

func (r *ParameterHandlingResourcePlanModifierResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *parameterHandlingResourcePlanModifierResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *ParameterHandlingResourcePlanModifierResource) read(data *parameterHandlingResourcePlanModifierResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			if opts.StringValue != nil {
				data.StringValue = types.StringValue(*opts.StringValue)
			} else {
				data.StringValue = types.StringNull()
			}
		}
	}
	return diags
}

func (r *ParameterHandlingResourcePlanModifierResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *parameterHandlingResourcePlanModifierResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &ParameterHandlingResourcePlanModifierOpts{}
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

func (r *ParameterHandlingResourcePlanModifierResource) update(opts *ParameterHandlingResourcePlanModifierOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not update resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingResourcePlanModifierResource) setUpdateActionsOutput(ctx context.Context, response *resource.UpdateResponse, opts *ParameterHandlingResourcePlanModifierOpts, plan *parameterHandlingResourcePlanModifierResourceModelV0, state *parameterHandlingResourcePlanModifierResourceModelV0) {
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

func (r *ParameterHandlingResourcePlanModifierResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
