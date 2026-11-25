package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// v0941ResourceIdentifierWithArguments migrates functions, procedures, and external functions to use the new identifier type.
// They're already using old identifier with arguments, but the only case where parentheses weren't specified
// (which are essential in the new identifier) is for empty argument list.
func v0941ResourceIdentifierWithArguments(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldId := rawState["id"].(string)
	if !strings.HasSuffix(oldId, ")") {
		oldId += "()"
	}
	id, err := sdk.ParseSchemaObjectIdentifierWithArguments(oldId)
	if err != nil {
		return nil, err
	}
	rawState["id"] = id.FullyQualifiedName()

	return rawState, nil
}

func migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	oldId := helpers.DecodeSnowflakeIDLegacy(rawState["id"].(string))
	rawState["id"] = oldId.FullyQualifiedName()
	return rawState, nil
}
