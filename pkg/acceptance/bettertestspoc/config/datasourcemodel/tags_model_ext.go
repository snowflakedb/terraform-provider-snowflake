package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (t *TagsModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *TagsModel {
	return t.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}
