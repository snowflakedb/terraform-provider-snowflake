package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (g *GitRepositoriesModel) WithRowsAndFrom(rows int, from string) *GitRepositoriesModel {
	return g.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}

func (g *GitRepositoriesModel) WithEmptyIn() *GitRepositoriesModel {
	return g.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{}),
	)
}

func (g *GitRepositoriesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *GitRepositoriesModel {
	return g.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}
