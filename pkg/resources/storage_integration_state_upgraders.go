package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v2_7_0_StorageIntegrationUsePrivatelinkEndpointUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	storageIntegrationName, ok := rawState["name"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot read storage integration name from the state")
	}

	client := meta.(*provider.Context).Client
	storageIntegrationInSnowflake, err := client.StorageIntegrations.Describe(ctx, sdk.NewAccountObjectIdentifier(storageIntegrationName))
	if err != nil {
		return nil, err
	}
	oldDescribeOutput := schemas.DescribeStorageIntegrationToSchema(storageIntegrationInSnowflake)
	rawState[DescribeOutputAttributeName] = []any{oldDescribeOutput}
	rawState["use_privatelink_endpoint"] = BooleanDefault

	return rawState, nil
}
