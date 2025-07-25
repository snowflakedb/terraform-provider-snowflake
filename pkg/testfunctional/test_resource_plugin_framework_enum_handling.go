package testfunctional

import (
	"context"
	"fmt"
	"strings"

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
	Name                  types.String                        `tfsdk:"name"`
	EnumValue             customtypes.EnumValue[SomeEnumType] `tfsdk:"enum_value"`
	EnumValueBackingField types.String                        `tfsdk:"enum_value_backing_field"`
	Id                    types.String                        `tfsdk:"id"`
}

type EnumHandlingOpts struct {
	EnumValue *SomeEnumType
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
			"enum_value": schema.StringAttribute{
				CustomType:  customtypes.EnumType[SomeEnumType]{},
				Description: "String value - enum.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					customplanmodifiers.EnumSuppressor[SomeEnumType](),
				},
			},
			"enum_value_backing_field": schema.StringAttribute{
				Description: "String value backing field.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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

func (r *EnumHandlingResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if request.State.Raw.IsNull() || request.Plan.Raw.IsNull() {
		return
	}

	var plan, state *enumHandlingResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	// we run normal equal here, as customplanmodifiers.EnumSuppressor handles it
	if !plan.EnumValue.Equal(state.EnumValue) {
		plan.EnumValueBackingField = types.StringUnknown()
	}

	response.Diagnostics.Append(response.Plan.Set(ctx, &plan)...)
}

func (r *EnumHandlingResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		response.Diagnostics.AddError("Could not read resources state", err.Error())
	} else if opts.EnumValue != nil {
		response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("enum_value"), *opts.EnumValue)...)
	}
}

func (r *EnumHandlingResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *enumHandlingResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &EnumHandlingOpts{}
	err := StringEnumAttributeCreate(data.EnumValue, &opts.EnumValue)
	if err != nil {
		response.Diagnostics.AddError("Error creating some enum type", err.Error())
	}

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
	} else if opts.EnumValue != nil {
		data.EnumValueBackingField = types.StringValue(string(*opts.EnumValue))
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
	} else if opts.EnumValue != nil {
		newValue := *opts.EnumValue
		// We don't need conversion here as https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom#semantic-equality should handle it for us:
		// `When refreshing a resource, the response new state value from the Read method logic is compared to the request prior state value.`
		if string(newValue) != data.EnumValueBackingField.ValueString() {
			data.EnumValue = customtypes.NewEnumValue(newValue)
		}
		data.EnumValueBackingField = types.StringValue(string(newValue))
	}
	return diags
}

func (r *EnumHandlingResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *enumHandlingResourceModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	opts := &EnumHandlingOpts{}
	err := StringEnumAttributeUpdate(plan.EnumValue, state.EnumValue, &opts.EnumValue, &opts.EnumValue)
	if err != nil {
		response.Diagnostics.AddError("Error updating some enum type", err.Error())
	}

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

func (r *EnumHandlingResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
