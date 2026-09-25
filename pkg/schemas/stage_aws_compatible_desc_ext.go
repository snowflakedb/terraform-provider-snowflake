package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (stageAwsCompatibleToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return collections.MergeMaps(stageCommonDescribeSchema(), map[string]*schema.Schema{
		"location": stageCompatibleLocationDescribeSchema(),
	})
}

func (stageAwsCompatibleToSchemaMapper) additionalToSchema(src *sdk.StageAwsCompatible, dst map[string]any) {
	mapStageCommonDescribe(stageCommonFromAwsCompatible(src), dst)
	mapStageCompatibleLocation(src.Location, dst)
}
