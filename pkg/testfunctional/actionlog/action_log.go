package actionlog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ActionsLogPropertyName = "actions_log"

type ActionsLogEmbeddable struct {
	ActionsLog types.List `tfsdk:"actions_log"`
}

type ActionLogEntry struct {
	Action types.String `tfsdk:"action"`
	Field  types.String `tfsdk:"field"`
	Value  types.String `tfsdk:"value"`
}

func GetActionsLogSchema() schema.Attribute {
	return schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: getActionLogEntrySchema(),
		},
		Computed: true,
	}
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

func GetActionLogEntryTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action": types.StringType,
		"field":  types.StringType,
		"value":  types.StringType,
	}
}

func ActionEntry(action string, field string, value string) ActionLogEntry {
	return ActionLogEntry{
		Action: types.StringValue(action),
		Field:  types.StringValue(field),
		Value:  types.StringValue(value),
	}
}

func AppendActions(ctx context.Context, actionsLog *ActionsLogEmbeddable, actionsProvider func() []ActionLogEntry) diag.Diagnostics {
	existingEntries := actionsLog.ActionsLog.Elements()

	actions := actionsProvider()

	for _, a := range actions {
		entry, diags := types.ObjectValue(GetActionLogEntryTypes(), map[string]attr.Value{
			"action": a.Action,
			"field":  a.Field,
			"value":  a.Value,
		})
		if diags.HasError() {
			return diags
		}
		existingEntries = append(existingEntries, entry)
	}
	var diags diag.Diagnostics
	actionsLog.ActionsLog, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: GetActionLogEntryTypes()}, existingEntries)
	return diags
}
