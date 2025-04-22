package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

// TODO [this PR]: check usages
func (r *RowAccessPolicyModel) WithArgument(argument []sdk.TableColumnSignature) *RowAccessPolicyModel {
	maps := make([]config.Variable, len(argument))
	for i, v := range argument {
		maps[i] = config.MapVariable(map[string]config.Variable{
			"name": config.StringVariable(v.Name),
			"type": config.StringVariable(v.Type.ToSql()),
		})
	}
	r.Argument = tfconfig.SetVariable(maps...)
	return r
}
