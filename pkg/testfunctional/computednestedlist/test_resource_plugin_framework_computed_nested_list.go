package computednestedlist

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

func NewComputedNestedListResource() resource.Resource {
	return &computedNestedListResource{}
}

type computedNestedListResource struct{}

type computedNestedListResourceModelV0 struct {
	Name       types.String `tfsdk:"name"`
	Option     types.String `tfsdk:"option"`
	ActionsLog types.List   `tfsdk:"actions_log"`
	Id         types.String `tfsdk:"id"`
}

type ActionLogEntry struct {
	Action types.String `tfsdk:"action"`
	Field  types.String `tfsdk:"field"`
	Value  types.String `tfsdk:"value"`
}

func getActionLogEntrySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Required: true,
		},
		"field": schema.StringAttribute{
			Required: true,
		},
		"value": schema.StringAttribute{
			Required: true,
		},
	}
}

func getActionLogEntryTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action": types.StringType,
		"field":  types.StringType,
		"value":  types.StringType,
	}
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
				Description: "Which implementation option should be tested. Available values: STRUCT, EXPLICIT",
				Required:    true,
			},
			"actions_log": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: getActionLogEntrySchema(),
				},
				Computed: true,
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
	default:
		response.Diagnostics.AddError("Use correct option", "Available options are: STRUCT, EXPLICIT")
		return
	}

	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func setActionsOutputThroughStruct(ctx context.Context, response *resource.CreateResponse, data *computedNestedListResourceModelV0) {
	var actions []ActionLogEntry
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

	actions := make([]ActionLogEntry, 0)
	actions = append(actions, actionEntry("SOME ACTION", "ON FIELD", "WITH VALUE"))
	actions = append(actions, actionEntry("SOME OTHER ACTION", "ON OTHER FIELD", "WITH OTHER VALUE"))

	for _, a := range actions {
		entry, diags := types.ObjectValue(getActionLogEntryTypes(), map[string]attr.Value{
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
	data.ActionsLog, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getActionLogEntryTypes()}, actions)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

func actionEntry(action string, field string, value string) ActionLogEntry {
	return ActionLogEntry{
		Action: types.StringValue(action),
		Field:  types.StringValue(field),
		Value:  types.StringValue(value),
	}
}

func (r *computedNestedListResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *computedNestedListResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *computedNestedListResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
