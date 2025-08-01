package testfunctional

import (
	"context"
	"strconv"

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

var _ resource.ResourceWithConfigure = &ZeroValuesResource{}

func NewZeroValuesResource() resource.Resource {
	return &ZeroValuesResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[ZeroValuesOpts]("zero_values_handling"),
	}
}

type ZeroValuesResource struct {
	common.HttpServerEmbeddable[ZeroValuesOpts]
}

type zeroValuesResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	BoolValue   types.Bool   `tfsdk:"bool_value"`
	IntValue    types.Int64  `tfsdk:"int_value"`
	StringValue types.String `tfsdk:"string_value"`
	Id          types.String `tfsdk:"id"`

	common.ActionsLogEmbeddable
}

type ZeroValuesOpts struct {
	BoolValue   *bool
	IntValue    *int
	StringValue *string
}

func (r *ZeroValuesResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zero_values"
}

func (r *ZeroValuesResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			// TODO [mux-PRs]: another setup for Optional+Computed combo
			"bool_value": schema.BoolAttribute{
				Description: "Boolean value.",
				Optional:    true,
			},
			"int_value": schema.Int64Attribute{
				Description: "Int value.",
				Optional:    true,
			},
			"string_value": schema.StringAttribute{
				Description: "String value.",
				Optional:    true,
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

func (r *ZeroValuesResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *ZeroValuesResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *zeroValuesResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &ZeroValuesOpts{}
	_ = BooleanAttributeCreate(data.BoolValue, &opts.BoolValue)
	_ = Int64AttributeCreate(data.IntValue, &opts.IntValue)
	_ = StringAttributeCreate(data.StringValue, &opts.StringValue)

	setCreateActionsOutput(ctx, response, opts, data)

	response.Diagnostics.Append(r.create(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func setCreateActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *ZeroValuesOpts, data *zeroValuesResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.BoolValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "bool_value", strconv.FormatBool(*opts.BoolValue)))
		}
		if opts.IntValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "int_value", strconv.Itoa(*opts.IntValue)))
		}
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "string_value", *opts.StringValue))
		}
		return actions
	})...)
}

func (r *ZeroValuesResource) create(opts *ZeroValuesOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *ZeroValuesResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *zeroValuesResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *ZeroValuesResource) read(data *zeroValuesResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.BoolValue != nil {
			data.BoolValue = types.BoolValue(*opts.BoolValue)
		}
		if opts.IntValue != nil {
			data.IntValue = types.Int64Value(int64(*opts.IntValue))
		}
		if opts.StringValue != nil {
			data.StringValue = types.StringValue(*opts.StringValue)
		}
	}
	return diags
}

func (r *ZeroValuesResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *zeroValuesResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &ZeroValuesOpts{}
	booleanAttributeUpdate(plan.BoolValue, state.BoolValue, &opts.BoolValue, &opts.BoolValue)
	int64AttributeUpdate(plan.IntValue, state.IntValue, &opts.IntValue, &opts.IntValue)
	stringAttributeUpdate(plan.StringValue, state.StringValue, &opts.StringValue, &opts.StringValue)

	setUpdateActionsOutput(ctx, response, opts, plan, state)

	response.Diagnostics.Append(r.update(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func setUpdateActionsOutput(ctx context.Context, response *resource.UpdateResponse, opts *ZeroValuesOpts, plan *zeroValuesResourceModelV0, state *zeroValuesResourceModelV0) {
	plan.ActionsLogEmbeddable = state.ActionsLogEmbeddable
	response.Diagnostics.Append(common.AppendActions(ctx, &plan.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.BoolValue != nil {
			actions = append(actions, common.ActionEntry("UPDATE - SET", "bool_value", strconv.FormatBool(*opts.BoolValue)))
		} else {
			actions = append(actions, common.ActionEntry("UPDATE - UNSET", "bool_value", "nil"))
		}
		if opts.IntValue != nil {
			actions = append(actions, common.ActionEntry("UPDATE - SET", "int_value", strconv.Itoa(*opts.IntValue)))
		} else {
			actions = append(actions, common.ActionEntry("UPDATE - UNSET", "int_value", "nil"))
		}
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("UPDATE - SET", "string_value", *opts.StringValue))
		} else {
			actions = append(actions, common.ActionEntry("UPDATE - UNSET", "string_value", "nil"))
		}
		return actions
	})...)
}

func (r *ZeroValuesResource) update(opts *ZeroValuesOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *ZeroValuesResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
