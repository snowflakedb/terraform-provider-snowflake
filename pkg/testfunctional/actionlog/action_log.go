package actionlog

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ActionLogEntry struct {
	Action types.String `tfsdk:"action"`
	Field  types.String `tfsdk:"field"`
	Value  types.String `tfsdk:"value"`
}

func GetActionLogEntrySchema() map[string]schema.Attribute {
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
