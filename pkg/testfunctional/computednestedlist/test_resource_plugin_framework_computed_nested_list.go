package computednestedlist

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewComputedNestedListResource() resource.Resource {
	return &computedNestedListResource{}
}

type computedNestedListResource struct{}

type computedNestedListResourceModelV0 struct {
	Name   types.String `tfsdk:"name"`
	Option types.String `tfsdk:"option"`
	Id     types.String `tfsdk:"id"`

	common.ActionsLogEmbeddable
}

func (r *computedNestedListResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_computed_nested_list"
}

func (r *computedNestedListResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			"option": schema.StringAttribute{
				Description: "Which implementation option should be tested. Available values: STRUCT, EXPLICIT, DEDICATED",
				Required:    true,
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

func (r *computedNestedListResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *computedNestedListResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	switch data.Option.ValueString() {
	case "STRUCT":
		setActionsOutputThroughStruct(ctx, response, data)
	case "EXPLICIT":
		setActionsOutputExplicit(ctx, response, data)
	case "DEDICATED":
		setActionsOutputDedicated(ctx, response, data)
	default:
		response.Diagnostics.AddError("Use correct option", "Available options are: STRUCT, EXPLICIT, DEDICATED")
		return
	}

	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func setActionsOutputThroughStruct(ctx context.Context, response *resource.CreateResponse, data *computedNestedListResourceModelV0) {
	var actions []common.ActionLogEntry
	diags := data.ActionsLog.ElementsAs(ctx, &actions, false)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
	panic("we return above because of `Value Conversion Error` which happens only for `Computed` list")
	//actions = append(actions, actionEntry("DOES", "NOT", "MATTER"))
	//data.ActionsLog, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getActionLogEntryTypes()}, actions)
	//if diags.HasError() {
	//	response.Diagnostics.Append(diags...)
	//	return
	//}
}

func setActionsOutputExplicit(ctx context.Context, response *resource.CreateResponse, data *computedNestedListResourceModelV0) {
	existingEntries := data.ActionsLog.Elements()

	actions := make([]common.ActionLogEntry, 0)
	actions = append(actions, common.ActionEntry("SOME ACTION", "ON FIELD", "WITH VALUE"))
	actions = append(actions, common.ActionEntry("SOME OTHER ACTION", "ON OTHER FIELD", "WITH OTHER VALUE"))

	for _, a := range actions {
		entry, diags := types.ObjectValue(common.GetActionLogEntryTypes(), map[string]attr.Value{
			"action": a.Action,
			"field":  a.Field,
			"value":  a.Value,
		})
		if diags.HasError() {
			response.Diagnostics.Append(diags...)
			return
		}
		existingEntries = append(existingEntries, entry)
	}
	var diags diag.Diagnostics
	data.ActionsLog, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: common.GetActionLogEntryTypes()}, actions)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

func setActionsOutputDedicated(ctx context.Context, response *resource.CreateResponse, data *computedNestedListResourceModelV0) {
	response.Diagnostics.Append(common.AppendActions(ctx, &data.ActionsLogEmbeddable, func() []common.ActionLogEntry {
		actions := make([]common.ActionLogEntry, 0)
		actions = append(actions, common.ActionEntry("SOME ACTION", "ON FIELD", "WITH VALUE"))
		actions = append(actions, common.ActionEntry("SOME OTHER ACTION", "ON OTHER FIELD", "WITH OTHER VALUE"))
		return actions
	})...)
}

func (r *computedNestedListResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *computedNestedListResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *computedNestedListResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
