package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
)

// migrateTableToV2_24_0 upgrades the state written before v2.24.0. Each change of that release is one step. It is
// best-effort: values it does not recognize are left unchanged, because the next refresh rewrites the columns from
// DESCRIBE TABLE anyway.
func migrateTableToV2_24_0(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	upgradeTableBooleanStringDefaults(rawState)
	columns, ok := rawState["column"].([]any)
	if !ok {
		return rawState, nil
	}
	for _, c := range columns {
		column, ok := c.(map[string]any)
		if !ok {
			continue
		}
		upgradeColumnDefault(column)
		upgradeColumnMaskingPolicy(column)
		upgradeColumnNotNull(column)
	}
	return rawState, nil
}

// upgradeColumnNotNull replaces nullable with its inverse, not_null. A column without a boolean nullable, or one that
// already has not_null, keeps not_null as it is, and a missing not_null becomes false.
func upgradeColumnNotNull(column map[string]any) {
	nullable, hasNullable := column["nullable"]
	delete(column, "nullable")
	if _, ok := column["not_null"].(bool); ok {
		return
	}
	if nullable, ok := nullable.(bool); ok && hasNullable {
		column["not_null"] = !nullable
		return
	}
	column["not_null"] = false
}

// upgradeColumnMaskingPolicy converts the masking policy name to the masking_policy block. Without it, the old string
// value cannot be decoded into the list.
func upgradeColumnMaskingPolicy(column map[string]any) {
	switch maskingPolicy := column["masking_policy"].(type) {
	case []any:
		return
	case string:
		if maskingPolicy != "" {
			column["masking_policy"] = []any{map[string]any{"policy_name": maskingPolicy, "using": []any{}}}
			return
		}
	}
	column["masking_policy"] = []any{}
}

// upgradeTableBooleanStringDefaults sets the boolean string fields added in v2.24.0 to their schema default. Without
// it, the missing is_transient plans "default" and recreates the table, because is_transient forces a new resource.
func upgradeTableBooleanStringDefaults(rawState map[string]any) {
	for _, key := range []string{"is_transient", "enable_schema_evolution", "error_logging", "row_timestamp"} {
		if value, ok := rawState[key].(string); !ok || value == "" {
			rawState[key] = BooleanDefault
		}
	}
}

// upgradeColumnDefault converts the constant, expression, and sequence default kinds to the single expression.
// A text constant becomes a quoted string literal, and a sequence becomes <sequence>.NEXTVAL.
func upgradeColumnDefault(column map[string]any) {
	defaults, ok := column["default"].([]any)
	if !ok || len(defaults) == 0 {
		return
	}
	def, ok := defaults[0].(map[string]any)
	if !ok {
		column["default"] = []any{}
		return
	}
	nonEmpty := func(key string) (string, bool) {
		value, ok := def[key].(string)
		return value, ok && value != ""
	}
	var expression string
	if value, ok := nonEmpty("expression"); ok {
		expression = value
	} else if value, ok := nonEmpty("constant"); ok {
		expression = value
		if columnType, ok := column["type"].(string); ok {
			if dataType, err := datatypes.ParseDataType(columnType); err == nil && datatypes.IsTextDataType(dataType) {
				expression = snowflake.EscapeSnowflakeString(value)
			}
		}
	} else if value, ok := nonEmpty("sequence"); ok {
		expression = fmt.Sprintf("%s.NEXTVAL", value)
	}
	if expression == "" {
		column["default"] = []any{}
		return
	}
	column["default"] = []any{map[string]any{"expression": expression}}
}
