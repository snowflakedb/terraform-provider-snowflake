package testfunctional

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewZeroValuesResource() resource.Resource {
	return &ZeroValuesResource{}
}

type ZeroValuesResource struct{}

type zeroValuesResourceModelV0 struct {
	Name        types.String `tfsdk:"name"`
	BoolValue   types.Bool   `tfsdk:"bool_value"`
	IntValue    types.Int64  `tfsdk:"int_value"`
	StringValue types.String `tfsdk:"string_value"`
	ActionsLog  types.List   `tfsdk:"actions_log"`
	Id          types.String `tfsdk:"id"`
}

type zeroValuesOpts struct {
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
			"actions_log": schema.ListAttribute{
				ElementType: types.StringType,
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

func (r *ZeroValuesResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *zeroValuesResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &zeroValuesOpts{}
	booleanAttributeCreate(data.BoolValue, &opts.BoolValue)
	int64AttributeCreate(data.IntValue, &opts.IntValue)
	stringAttributeCreate(data.StringValue, &opts.StringValue)

	setActionsOutput(ctx, response, opts, data)

	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func setActionsOutput(ctx context.Context, response *resource.CreateResponse, opts *zeroValuesOpts, data *zeroValuesResourceModelV0) {
	var actions []string
	//diag := data.ActionsLog.ElementsAs(ctx, &actions, false)
	//if diag.HasError() {
	//	response.Diagnostics.Append(diag...)
	//	return
	//}
	//if opts.BoolValue != nil {
	//	actions = append(actions, actionlog.ActionEntry("CREATE", "BOOL", strconv.FormatBool(*opts.BoolValue)).ToString())
	//}
	//if opts.IntValue != nil {
	//	actions = append(actions, actionlog.ActionEntry("CREATE", "INT", strconv.Itoa(*opts.IntValue)).ToString())
	//}
	//if opts.StringValue != nil {
	//	actions = append(actions, actionlog.ActionEntry("CREATE", "STRING", *opts.StringValue).ToString())
	//}
	actionValues := make([]attr.Value, len(actions))
	for i, a := range actions {
		actionValues[i] = types.StringValue(a)
	}
	var diags diag.Diagnostics
	data.ActionsLog, diags = types.ListValue(types.StringType, actionValues)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

func (r *ZeroValuesResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *ZeroValuesResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *ZeroValuesResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
