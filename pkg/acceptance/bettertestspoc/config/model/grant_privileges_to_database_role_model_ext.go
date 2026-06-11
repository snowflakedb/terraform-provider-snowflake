package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (g *GrantPrivilegesToDatabaseRoleModel) WithPrivileges(privileges ...string) *GrantPrivilegesToDatabaseRoleModel {
	privilegeStringVariables := collections.Map(privileges, func(privilege string) config.Variable { return config.StringVariable(privilege) })
	g.WithPrivilegesValue(config.ListVariable(privilegeStringVariables...))
	return g
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithAccountObjectPrivileges(privileges ...sdk.AccountObjectPrivilege) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithPrivileges(collections.Map(privileges, func(p sdk.AccountObjectPrivilege) string { return string(p) })...)
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithSchemaPrivileges(privileges ...sdk.SchemaPrivilege) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithPrivileges(collections.Map(privileges, func(p sdk.SchemaPrivilege) string { return string(p) })...)
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithSchemaObjectPrivileges(privileges ...sdk.SchemaObjectPrivilege) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithPrivileges(collections.Map(privileges, func(p sdk.SchemaObjectPrivilege) string { return string(p) })...)
}
