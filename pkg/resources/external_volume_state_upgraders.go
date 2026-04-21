package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v2_14_0_ExternalVolumeStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	idRaw, ok := rawState["id"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot read external volume id from the state")
	}
	id, err := sdk.ParseAccountObjectIdentifier(idRaw)
	if err != nil {
		return nil, err
	}

	client := meta.(*provider.Context).Client
	properties, err := client.ExternalVolumes.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	details, err := sdk.ParseExternalVolumeDescribed(properties)
	if err != nil {
		return nil, err
	}

	detailsSchema := schemas.ExternalVolumeDetailsToSchema(details)

	rawState[DescribeOutputAttributeName] = []any{detailsSchema}

	// Clear storage_aws_external_id from storage_location entries.
	// Previously this was a Computed field populated by SF; now it's Optional (user-configured only).
	if storageLocations, ok := rawState["storage_location"].([]any); ok {
		for _, loc := range storageLocations {
			if locMap, ok := loc.(map[string]any); ok {
				locMap["storage_aws_external_id"] = ""
			}
		}
	}

	return rawState, nil
}
