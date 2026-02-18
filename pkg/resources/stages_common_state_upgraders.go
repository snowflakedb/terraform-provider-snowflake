package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func v2_14_0_StageStateUpgrader(describeToSchema func(sdk.StageDetails) (map[string]any, error)) schema.StateUpgradeFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		if rawState == nil {
			return rawState, nil
		}

		idRaw, ok := rawState["id"].(string)
		if !ok {
			return nil, fmt.Errorf("cannot read stage id from the state")
		}
		id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
		if err != nil {
			return nil, err
		}

		client := meta.(*provider.Context).Client
		properties, err := client.Stages.Describe(ctx, id)
		if err != nil {
			return nil, err
		}

		details, err := sdk.ParseStageDetails(properties)
		if err != nil {
			return nil, err
		}

		detailsSchema, err := describeToSchema(*details)
		if err != nil {
			return nil, err
		}

		rawState[DescribeOutputAttributeName] = []any{detailsSchema}

		return rawState, nil
	}
}

var (
	v2_14_0_ExternalAzureStageStateUpgrader        = v2_14_0_StageStateUpgrader(schemas.StageDescribeToSchema)
	v2_14_0_ExternalGcsStageStateUpgrader          = v2_14_0_StageStateUpgrader(schemas.StageDescribeToSchema)
	v2_14_0_ExternalS3CompatibleStageStateUpgrader = v2_14_0_StageStateUpgrader(schemas.AwsCompatibleStageDescribeToSchema)
	v2_14_0_InternalStageStateUpgrader             = v2_14_0_StageStateUpgrader(schemas.StageDescribeToSchema)
	v2_14_0_ExternalS3StageStateUpgrader           = v2_14_0_StageStateUpgrader(schemas.AwsStageDescribeToSchema)
)
