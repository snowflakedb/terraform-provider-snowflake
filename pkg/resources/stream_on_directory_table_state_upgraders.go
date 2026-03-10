package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v2_14_0_StreamOnDirectoryTableStateUpgrader(ctx context.Context, state map[string]any, meta any) (map[string]any, error) {
	if state == nil {
		return state, nil
	}

	stageId, err := sdk.ParseSchemaObjectIdentifier(state["stage"].(string))
	if err != nil {
		return nil, err
	}

	state["stage"] = stageId.FullyQualifiedName()

	return state, nil
}
