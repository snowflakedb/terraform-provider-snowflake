package resources

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// v2_15_0_ImageRepositoryEncryptionUpgrader is needed because:
// - we are adding encryption to the show output and to the config
// - this means that the old output in handleExternalChangesToObject would be empty after migration
// - that caused the plan to show a diff for encryption (which is force new and would cause a recreation)
// - this upgrade fixes that by setting the default value for encryption
func v2_15_0_ImageRepositoryEncryptionUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldShowOutputRaw, ok := rawState[ShowOutputAttributeName].([]any)
	if !ok || len(oldShowOutputRaw) != 1 {
		return rawState, nil
	}
	oldShowOutput, ok := oldShowOutputRaw[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cannot read image repository show output from the state")
	}

	idRaw, ok := rawState["id"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot read image repository id from the state")
	}
	id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
	if err != nil {
		return nil, err
	}

	client := meta.(*provider.Context).Client
	imageRepositoryInSnowflake, err := client.ImageRepositories.ShowByIDSafely(ctx, id)
	if err != nil {
		return nil, err
	}

	oldShowOutput["encryption"] = string(imageRepositoryInSnowflake.Encryption)
	rawState[ShowOutputAttributeName] = []any{oldShowOutput}

	return rawState, nil
}
