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

	detailsSchema, err := schemas.ExternalVolumeDetailsToSchema(details)
	if err != nil {
		return nil, err
	}

	rawState[DescribeOutputAttributeName] = []any{detailsSchema}

	return rawState, nil
}
