package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

func v0_95_0_RowAccessPolicyStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	rawState["body"] = rawState["row_access_expression"]
	delete(rawState, "row_access_expression")

	signature := rawState["signature"].(map[string]any)
	args := make([]map[string]any, 0)
	for k, v := range signature {
		args = append(args, map[string]any{
			"name": strings.ToUpper(k),
			"type": v,
		})
	}
	rawState["argument"] = args
	delete(rawState, "signature")

	return migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx, rawState, meta)
}

func v200RowAccessPolicyStateUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	arguments := rawState["argument"].([]any)
	if len(arguments) == 0 {
		return rawState, nil
	}
	args := make([]map[string]any, 0)
	for _, v := range arguments {
		argument := v.(map[string]any)
		columnDataType, err := datatypes.ParseDataType(argument["type"].(string))
		if err != nil {
			return nil, fmt.Errorf("updating the snowflake_row_access_policy resource state for the v2.0.0 provider version, error: %w", err)
		}
		args = append(args, map[string]any{
			"name": argument["name"].(string),
			"type": columnDataType.ToSql(),
		})
	}
	rawState["argument"] = args

	return rawState, nil
}
