package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func v0_95_0_MaskingPolicyStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["body"] = rawState["masking_expression"]

	signature := rawState["signature"].([]any)
	if len(signature) != 1 {
		return nil, fmt.Errorf("corrupted signature in state: expected list of length 1, got %d", len(signature))
	}
	columns := signature[0].(map[string]any)["column"].([]any)
	args := make([]map[string]any, 0)
	for _, v := range columns {
		column := v.(map[string]any)
		args = append(args, map[string]any{
			"name": strings.ToUpper(column["name"].(string)),
			"type": column["type"].(string),
		})
	}
	rawState["argument"] = args

	return migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx, rawState, meta)
}

func v200MaskingPolicyStateUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	arguments := rawState["argument"].([]any)
	if len(arguments) > 0 {
		return nil, fmt.Errorf("updating the snowflake_masking_policy resource state for the v2.0.0 provider version: expected at least 1 argument")
	}
	args := make([]map[string]any, 0)
	for _, v := range arguments {
		argument := v.(map[string]any)
		columnDataType, err := datatypes.ParseDataType(argument["type"].(string))
		if err != nil {
			return nil, fmt.Errorf("updating the snowflake_masking_policy resource state for the v2.0.0 provider version, error: %w", err)
		}
		args = append(args, map[string]any{
			"name": argument["name"].(string),
			"type": columnDataType.ToSql(),
		})
	}
	rawState["argument"] = args

	returnDataTypeRaw := rawState["return_data_type"].(string)
	returnDataType, err := datatypes.ParseDataType(returnDataTypeRaw)
	if err != nil {
		return nil, fmt.Errorf("updating the snowflake_masking_policy resource state for the v2.0.0 provider version, error: %w", err)
	}
	rawState["return_data_type"] = returnDataType.ToSql()

	return rawState, nil
}
