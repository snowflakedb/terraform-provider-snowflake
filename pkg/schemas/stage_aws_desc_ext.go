package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (stageAwsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return stageAwsDescribeSchema()
}

func (stageAwsToSchemaMapper) additionalToSchema(src *sdk.StageAws, dst map[string]any) {
	mapStageCommonDescribe(stageCommonFromAws(src), dst)
	mapStagePrivateLink(src.PrivateLink, dst)
	mapStageAwsLocation(src.Location, dst)
}
