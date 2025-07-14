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

var _ resource.ResourceWithConfigure = &ParameterHandlingReadLogicResource{}

func NewParameterHandlingReadLogicResource() resource.Resource {
	return &ParameterHandlingReadLogicResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[ParameterHandlingReadLogicOpts]("parameter_handling_read_logic"),
	}
}

type ParameterHandlingReadLogicResource struct {
	common.HttpServerEmbeddable[ParameterHandlingReadLogicOpts]
}

type parameterHandlingReadLogicResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	StringValue types.String `tfsdk:"string_value"`
	Id          types.String `tfsdk:"id"`

	common.ActionsLogEmbeddable
}

type ParameterHandlingReadLogicOpts struct {
	StringValue *string
	Level       string
}

func (r *ParameterHandlingReadLogicResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_parameter_handling_read_logic"
}

func (r *ParameterHandlingReadLogicResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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

func (r *ParameterHandlingReadLogicResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			// TODO
			response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("string_value"), *opts.StringValue)...)
		}
	}
}

func (r *ParameterHandlingReadLogicResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *parameterHandlingReadLogicResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &ParameterHandlingReadLogicOpts{}
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

func (r *ParameterHandlingReadLogicResource) setCreateActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *ParameterHandlingReadLogicOpts, data *parameterHandlingReadLogicResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "string_value", *opts.StringValue))
		}
		return actions
	})...)
}

func (r *ParameterHandlingReadLogicResource) create(opts *ParameterHandlingReadLogicOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingReadLogicResource) readAfterCreateOrUpdate(data *parameterHandlingReadLogicResourceModelV0) diag.Diagnostics {
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

func (r *ParameterHandlingReadLogicResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *parameterHandlingReadLogicResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *ParameterHandlingReadLogicResource) read(data *parameterHandlingReadLogicResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// TODO: denormalization
	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			// If the level differs we set the state to null, to trigger setting.
			// It's not ideal as the plan will output null -> value plan.
			// Can't set to unknown because then "The returned state contains unknown values." error is returned.
			if opts.StringValue != nil && opts.Level == "OBJECT" {
				data.StringValue = types.StringValue(*opts.StringValue)
			} else {
				data.StringValue = types.StringNull()
			}
		}
	}
	return diags
}

func (r *ParameterHandlingReadLogicResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *parameterHandlingReadLogicResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &ParameterHandlingReadLogicOpts{}
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

func (r *ParameterHandlingReadLogicResource) update(opts *ParameterHandlingReadLogicOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not update resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingReadLogicResource) setUpdateActionsOutput(ctx context.Context, response *resource.UpdateResponse, opts *ParameterHandlingReadLogicOpts, plan *parameterHandlingReadLogicResourceModelV0, state *parameterHandlingReadLogicResourceModelV0) {
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

func (r *ParameterHandlingReadLogicResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
