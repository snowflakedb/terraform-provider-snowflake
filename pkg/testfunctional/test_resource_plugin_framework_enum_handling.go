package testfunctional

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (e SomeEnumType) FromString(s string) (SomeEnumType, error) {
	return ToSomeEnumType(s)
}

type SomeEnumType string

const (
	SomeEnumTypeVersion1 SomeEnumType = "VERSION_1"
	SomeEnumTypeVersion2 SomeEnumType = "VERSION_2"
	SomeEnumTypeVersion3 SomeEnumType = "VERSION_3"
)

func ToSomeEnumType(s string) (SomeEnumType, error) {
	switch strings.ToUpper(s) {
	case string(SomeEnumTypeVersion1):
		return SomeEnumTypeVersion1, nil
	case string(SomeEnumTypeVersion2):
		return SomeEnumTypeVersion2, nil
	case string(SomeEnumTypeVersion3):
		return SomeEnumTypeVersion3, nil
	default:
		return "", fmt.Errorf("invalid some enum type: %s", s)
	}
}

var _ resource.ResourceWithConfigure = &EnumHandlingResource{}

func NewEnumHandlingResource() resource.Resource {
	return &EnumHandlingResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[EnumHandlingOpts]("enum_handling"),
	}
}

type EnumHandlingResource struct {
	common.HttpServerEmbeddable[EnumHandlingOpts]
}

type enumHandlingResourceModelV0 struct {
	Name        types.String                        `tfsdk:"name"`
	StringValue customtypes.EnumValue[SomeEnumType] `tfsdk:"string_value"`
	Id          types.String                        `tfsdk:"id"`

	common.ActionsLogEmbeddable
}

type EnumHandlingOpts struct {
	StringValue *SomeEnumType
}

func (r *EnumHandlingResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_enum_handling"
}

func (r *EnumHandlingResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			"string_value": schema.StringAttribute{
				CustomType:  customtypes.EnumType[SomeEnumType]{},
				Description: "String value - enum.",
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

func (r *EnumHandlingResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
	} else if opts.StringValue != nil {
		// TODO
		response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("string_value"), *opts.StringValue)...)
	}
}

func (r *EnumHandlingResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *enumHandlingResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &EnumHandlingOpts{}
	err := stringEnumAttributeCreate(data.StringValue, &opts.StringValue, ToSomeEnumType)
	if err != nil {
		response.Diagnostics.AddError("Error creating some enum type", err.Error())
	}

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

func (r *EnumHandlingResource) setCreateActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *EnumHandlingOpts, data *enumHandlingResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("CREATE", "string_value", string(*opts.StringValue)))
		}
		return actions
	})...)
}

func (r *EnumHandlingResource) create(opts *EnumHandlingOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *EnumHandlingResource) readAfterCreateOrUpdate(data *enumHandlingResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else if opts.StringValue != nil {
		// TODO
		// data.StringValueBackingField = types.StringValue(*opts.StringValue)
	}
	return diags
}

func (r *EnumHandlingResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *enumHandlingResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *EnumHandlingResource) read(data *enumHandlingResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else if opts.StringValue != nil {
		if data.StringValue.IsNull() {
			data.StringValue = customtypes.NewEnumValue(*opts.StringValue)
		} else {
			areTheSame, err := sameAfterNormalization(data.StringValue.ValueString(), string(*opts.StringValue), ToSomeEnumType)
			if err != nil {
				diags.AddError("Could not read resources state", err.Error())
				return diags
			}
			if !areTheSame {
				data.StringValue = customtypes.NewEnumValue(*opts.StringValue)
			}
		}
	}
	return diags
}

func (r *EnumHandlingResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *enumHandlingResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &EnumHandlingOpts{}
	err := stringEnumAttributeUpdate(plan.StringValue, state.StringValue, &opts.StringValue, &opts.StringValue, ToSomeEnumType)
	if err != nil {
		response.Diagnostics.AddError("Error updating some enum type", err.Error())
	}

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

func (r *EnumHandlingResource) update(opts *EnumHandlingOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not update resource", err.Error())
	}
	return diags
}

func (r *EnumHandlingResource) setUpdateActionsOutput(ctx context.Context, response *resource.UpdateResponse, opts *EnumHandlingOpts, plan *enumHandlingResourceModelV0, state *enumHandlingResourceModelV0) {
	plan.ActionsLogEmbeddable = state.ActionsLogEmbeddable
	response.Diagnostics.Append(common.AppendActions(ctx, &plan.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		if opts.StringValue != nil {
			actions = append(actions, common.ActionEntry("UPDATE - SET", "string_value", string(*opts.StringValue)))
		} else {
			actions = append(actions, common.ActionEntry("UPDATE - UNSET", "string_value", "nil"))
		}
		return actions
	})...)
}

func (r *EnumHandlingResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
