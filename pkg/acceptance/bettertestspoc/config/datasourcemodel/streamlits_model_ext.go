package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (s *StreamlitsModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *StreamlitsModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}
