package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

// v0941ResourceIdentifierWithArguments migrates functions, procedures, and external functions to use the new identifier type.
// They're already using old identifier with arguments, but the only case where parentheses weren't specified
// (which are essential in the new identifier) is for empty argument list.
func v0941ResourceIdentifierWithArguments(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	id := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rawState["id"].(string))
	rawState["id"] = sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), id.Name(), id.Arguments()...).FullyQualifiedName()

	return rawState, nil
}

func migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	oldId := helpers.DecodeSnowflakeID(rawState["id"].(string))
	rawState["id"] = oldId.FullyQualifiedName()
	return rawState, nil
}
