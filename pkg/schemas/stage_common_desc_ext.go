package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (stageCommonToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return stageCommonDescribeSchema()
}

func (stageCommonToSchemaMapper) additionalToSchema(src *sdk.StageCommon, dst map[string]any) {
	mapStageCommonDescribe(src, dst)
}
