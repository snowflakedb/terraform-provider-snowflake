package testfunctional

import (
	"context"
	"encoding/json"

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

var _ resource.ResourceWithConfigure = &ParameterHandlingPrivateResource{}

func NewParameterHandlingPrivateResource() resource.Resource {
	return &ParameterHandlingPrivateResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[ParameterHandlingPrivateOpts]("parameter_handling_private"),
	}
}

type ParameterHandlingPrivateResource struct {
	common.HttpServerEmbeddable[ParameterHandlingPrivateOpts]
}

type parameterHandlingPrivateResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	StringValue types.String `tfsdk:"string_value"`
	Id          types.String `tfsdk:"id"`

	common.ActionsLogEmbeddable
}

type ParameterHandlingPrivateJson struct {
	Value string `json:"value,omitempty"`
}

type ParameterHandlingPrivateOpts struct {
	StringValue *string
	Level       string
}

func (r *ParameterHandlingPrivateResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_parameter_handling_private"
}

func (r *ParameterHandlingPrivateResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
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

func (r *ParameterHandlingPrivateResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			if opts.Level == "OBJECT" {
				response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("string_value"), *opts.StringValue)...)
			}
		}
	}
}

func (r *ParameterHandlingPrivateResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *parameterHandlingPrivateResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &ParameterHandlingPrivateOpts{}
	stringAttributeCreate(data.StringValue, &opts.StringValue)

	r.setCreateActionsOutput(ctx, response, opts, data)

	response.Diagnostics.Append(r.create(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.readAfterCreate(ctx, &response)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *ParameterHandlingPrivateResource) setCreateActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *ParameterHandlingPrivateOpts, data *parameterHandlingPrivateResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "string_value", *opts.StringValue))
		}
		return actions
	})...)
}

func (r *ParameterHandlingPrivateResource) create(opts *ParameterHandlingPrivateOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingPrivateResource) readAfterCreate(ctx context.Context, response **resource.CreateResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			bytes, err := json.Marshal(ParameterHandlingPrivateJson{Value: *opts.StringValue})
			if err != nil {
				diags.AddError("Could not marshal json", err.Error())
				return diags
			}
			diags.Append((*response).Private.SetKey(ctx, "string_value_parameter", bytes)...)
		}
	}
	return diags
}

func (r *ParameterHandlingPrivateResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *parameterHandlingPrivateResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(ctx, data, request, &response)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *ParameterHandlingPrivateResource) read(ctx context.Context, data *parameterHandlingPrivateResourceModelV0, request resource.ReadRequest, response **resource.ReadResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			newValue := *opts.StringValue
			newLevel := opts.Level

			prevValueBytes, d := request.Private.GetKey(ctx, "string_value_parameter")
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}
			if prevValueBytes != nil {
				var prevValue ParameterHandlingPrivateJson
				err := json.Unmarshal(prevValueBytes, &prevValue)
				if err != nil {
					diags.AddError("Could not unmarshal json", err.Error())
					return diags
				}

				// if new value differs from the previous one
				if newValue != prevValue.Value {
					data.StringValue = types.StringValue(newValue)
					// if new level is not object we should set null
				} else if newLevel != "OBJECT" {
					data.StringValue = types.StringNull()
				} else {
					data.StringValue = types.StringValue(newValue)
				}
			}

			bytes, err := json.Marshal(ParameterHandlingPrivateJson{Value: newValue})
			if err != nil {
				diags.AddError("Could not marshal json", err.Error())
				return diags
			}
			diags.Append((*response).Private.SetKey(ctx, "string_value_parameter", bytes)...)
		}
	}
	return diags
}

func (r *ParameterHandlingPrivateResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *parameterHandlingPrivateResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &ParameterHandlingPrivateOpts{}
	stringAttributeUpdate(plan.StringValue, state.StringValue, &opts.StringValue, &opts.StringValue)

	r.setUpdateActionsOutput(ctx, response, opts, plan, state)

	response.Diagnostics.Append(r.update(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.readAfterUpdate(ctx, &response)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *ParameterHandlingPrivateResource) update(opts *ParameterHandlingPrivateOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not update resource", err.Error())
	}
	return diags
}

func (r *ParameterHandlingPrivateResource) readAfterUpdate(ctx context.Context, response **resource.UpdateResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		if opts.StringValue != nil {
			bytes, err := json.Marshal(ParameterHandlingPrivateJson{Value: *opts.StringValue})
			if err != nil {
				diags.AddError("Could not marshal json", err.Error())
				return diags
			}
			diags.Append((*response).Private.SetKey(ctx, "string_value_parameter", bytes)...)
		}
	}
	return diags
}

func (r *ParameterHandlingPrivateResource) setUpdateActionsOutput(ctx context.Context, response *resource.UpdateResponse, opts *ParameterHandlingPrivateOpts, plan *parameterHandlingPrivateResourceModelV0, state *parameterHandlingPrivateResourceModelV0) {
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

func (r *ParameterHandlingPrivateResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
