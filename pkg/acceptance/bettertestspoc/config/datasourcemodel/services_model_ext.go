package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (i *ServicesModel) WithEmptyIn() *ServicesModel {
	return i.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (i *ServicesModel) WithInComputePool(computePoolId sdk.AccountObjectIdentifier) *ServicesModel {
	return i.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"compute_pool": tfconfig.StringVariable(computePoolId.Name()),
		}),
	)
}

func (i *ServicesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *ServicesModel {
	return i.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}
