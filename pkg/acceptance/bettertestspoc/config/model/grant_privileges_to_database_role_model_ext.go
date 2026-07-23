package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (g *GrantPrivilegesToDatabaseRoleModel) WithPrivileges(privileges []string) *GrantPrivilegesToDatabaseRoleModel {
	privilegeStringVariables := collections.Map(privileges, func(privilege string) config.Variable { return config.StringVariable(privilege) })
	g.WithPrivilegesValue(config.ListVariable(privilegeStringVariables...))
	return g
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnInheritedSchemasInDatabase(id sdk.AccountObjectIdentifier) *GrantPrivilegesToDatabaseRoleModel {
	g.WithOnSchemaValue(config.ObjectVariable(map[string]config.Variable{
		"inherited": config.StringVariable(id.FullyQualifiedName()),
	}))
	return g
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnInheritedSchemaObjectsInDatabase(pluralObjectType sdk.PluralObjectType, id sdk.AccountObjectIdentifier) *GrantPrivilegesToDatabaseRoleModel {
	g.WithOnSchemaObjectValue(config.ObjectVariable(map[string]config.Variable{
		"inherited": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"object_type_plural": config.StringVariable(string(pluralObjectType)),
				"in_database":        config.StringVariable(id.FullyQualifiedName()),
			}),
		),
	}))
	return g
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnInheritedSchemaObjectsInSchema(pluralObjectType sdk.PluralObjectType, schemaId sdk.DatabaseObjectIdentifier) *GrantPrivilegesToDatabaseRoleModel {
	g.WithOnSchemaObjectValue(config.ObjectVariable(map[string]config.Variable{
		"inherited": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"object_type_plural": config.StringVariable(string(pluralObjectType)),
				"in_schema":          config.StringVariable(schemaId.FullyQualifiedName()),
			}),
		),
	}))
	return g
}
