package computednestedlisterror

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// There is an open issue for computed nested list: https://github.com/hashicorp/terraform-plugin-framework/issues/1104.
// TODO [mux-PR]: describe

func NewComputedNestedListResource() resource.Resource {
	return &ComputedNestedListResource{}
}

type ComputedNestedListResource struct{}

type computedNestedListResourceModelV0 struct {
	Name       types.String `tfsdk:"name"`
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

func (r *ComputedNestedListResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_computed_nested_list"
}

func (r *ComputedNestedListResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			"actions_log": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: getActionLogEntrySchema(),
				},
				Optional: true,
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

func (r *ComputedNestedListResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *computedNestedListResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	setActionsOutput(ctx, response, data)

	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func setActionsOutput(ctx context.Context, response *resource.CreateResponse, data *computedNestedListResourceModelV0) {
	var actions []ActionLogEntry
	diag := data.ActionsLog.ElementsAs(ctx, &actions, false)
	if diag.HasError() {
		response.Diagnostics.Append(diag...)
		return
	}
	actions = append(actions, actionEntry("DOES", "NOT", "MATTER"))
	data.ActionsLog, diag = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getActionLogEntryTypes()}, actions)
	if diag.HasError() {
		response.Diagnostics.Append(diag...)
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

func (r *ComputedNestedListResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *ComputedNestedListResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *ComputedNestedListResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
