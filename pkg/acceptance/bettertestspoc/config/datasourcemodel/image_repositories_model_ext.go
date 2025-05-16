package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (i *ImageRepositoriesModel) WithEmptyIn() *ImageRepositoriesModel {
	return i.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (i *ImageRepositoriesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *ImageRepositoriesModel {
	return i.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}
