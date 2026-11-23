package resources

import (
	"context"
)

// migrateTableIdAndNewBooleanStringFields rewrites the pipe-separated id to the fully qualified name (like
// migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName) and sets the boolean string fields added to
// snowflake_table in the same release to their unset value. Without it, the upgraded state has no value for them, so
// the next plan would replace the table (is_transient is ForceNew) or update the others.
func migrateTableIdAndNewBooleanStringFields(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	rawState, err := migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName(ctx, rawState, meta)
	if err != nil || rawState == nil {
		return rawState, err
	}
	for _, key := range []string{"is_transient", "enable_schema_evolution", "error_logging", "row_timestamp"} {
		rawState[key] = BooleanDefault
	}
	return rawState, nil
}
