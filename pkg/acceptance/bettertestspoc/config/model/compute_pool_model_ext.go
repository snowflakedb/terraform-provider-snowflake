package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (c *ComputePoolModel) WithBackupInstanceFamilies(backupInstanceFamilies ...sdk.ComputePoolInstanceFamily) *ComputePoolModel {
	if len(backupInstanceFamilies) == 0 {
		return c.WithBackupInstanceFamiliesValue(config.EmptyListVariable())
	}
	familyVars := collections.Map(backupInstanceFamilies, func(f sdk.ComputePoolInstanceFamily) tfconfig.Variable {
		return tfconfig.StringVariable(string(f))
	})
	return c.WithBackupInstanceFamiliesValue(tfconfig.ListVariable(familyVars...))
}
