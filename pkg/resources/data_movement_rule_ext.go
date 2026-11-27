package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	dataMovementRuleCreateExt      = TrackingCreateWrapper(resources.DataMovementRule, CreateDataMovementRule)
	dataMovementRuleReadExt        = TrackingReadWrapper(resources.DataMovementRule, ReadDataMovementRuleFunc(true))
	dataMovementRuleUpdateExt      = TrackingUpdateWrapper(resources.DataMovementRule, UpdateDataMovementRule)
	dataMovementRuleDescriptionExt = "Resource used to manage data movement rules. For more information, check [data movement rule documentation](https://docs.snowflake.com/en/sql-reference/sql/create-data-movement-rule)."
	dataMovementRuleImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.DataMovementRule, ImportDataMovementRule),
	}
)
