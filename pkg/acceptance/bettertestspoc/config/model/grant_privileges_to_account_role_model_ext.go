package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (g *GrantPrivilegesToAccountRoleModel) WithPrivileges(privileges []string) *GrantPrivilegesToAccountRoleModel {
	privilegeStringVariables := collections.Map(privileges, func(privilege string) config.Variable { return config.StringVariable(privilege) })
	g.WithPrivilegesValue(config.ListVariable(privilegeStringVariables...))
	return g
}
