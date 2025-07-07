package testfunctional

import (
	"context"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/actionlog"
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
	Id          types.String `tfsdk:"id"`

	actionlog.ActionsLogEmbeddable
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
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			actionlog.ActionsLogPropertyName: actionlog.GetActionsLogSchema(),
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
	response.Diagnostics.Append(actionlog.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []actionlog.ActionLogEntry {
		actions := make([]actionlog.ActionLogEntry, 0)
		if opts.BoolValue != nil {
			actions = append(actions, actionlog.ActionEntry("CREATE", "bool_value", strconv.FormatBool(*opts.BoolValue)))
		}
		if opts.IntValue != nil {
			actions = append(actions, actionlog.ActionEntry("CREATE", "int_value", strconv.Itoa(*opts.IntValue)))
		}
		if opts.StringValue != nil {
			actions = append(actions, actionlog.ActionEntry("CREATE", "string_value", *opts.StringValue))
		}
		return actions
	})...)
}

func (r *ZeroValuesResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *ZeroValuesResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *ZeroValuesResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
