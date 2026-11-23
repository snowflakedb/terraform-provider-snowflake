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
	}
	return rawState, nil
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
