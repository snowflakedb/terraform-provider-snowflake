package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (g *GrantPrivilegesToAccountRoleModel) WithPrivileges(privileges ...string) *GrantPrivilegesToAccountRoleModel {
	privilegeStringVariables := collections.Map(privileges, func(privilege string) config.Variable { return config.StringVariable(privilege) })
	g.WithPrivilegesValue(config.ListVariable(privilegeStringVariables...))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnAccountObject(objectType sdk.ObjectType, id sdk.AccountObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
	g.WithOnAccountObjectValue(config.ObjectVariable(map[string]config.Variable{
		"object_type": config.StringVariable(objectType.String()),
		"object_name": config.StringVariable(id.Name()),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnAllSchemasInDatabase(id sdk.AccountObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaValue(config.ObjectVariable(map[string]config.Variable{
		"all_schemas_in_database": config.StringVariable(id.FullyQualifiedName()),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnFutureSchemasInDatabase(id sdk.AccountObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaValue(config.ObjectVariable(map[string]config.Variable{
		"future_schemas_in_database": config.StringVariable(id.FullyQualifiedName()),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnSchemaName(id sdk.DatabaseObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaValue(config.ObjectVariable(map[string]config.Variable{
		"schema_name": config.StringVariable(id.FullyQualifiedName()),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnAllSchemaObjectsInSchema(pluralObjectType sdk.PluralObjectType, schemaId sdk.DatabaseObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaObjectValue(config.ObjectVariable(map[string]config.Variable{
		"all": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"object_type_plural": config.StringVariable(string(pluralObjectType)),
				"in_schema":          config.StringVariable(schemaId.FullyQualifiedName()),
			}),
		),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnFutureSchemaObjectsInSchema(pluralObjectType sdk.PluralObjectType, schemaId sdk.DatabaseObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaObjectValue(config.ObjectVariable(map[string]config.Variable{
		"future": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"object_type_plural": config.StringVariable(string(pluralObjectType)),
				"in_schema":          config.StringVariable(schemaId.FullyQualifiedName()),
			}),
		),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnInheritedAccountObjects(pluralObjectType sdk.PluralObjectType) *GrantPrivilegesToAccountRoleModel {
	g.WithOnAccountObjectValue(config.ObjectVariable(map[string]config.Variable{
		"inherited": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"object_type_plural": config.StringVariable(string(pluralObjectType)),
			}),
		),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnInheritedSchemasInAccount() *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaValue(config.ObjectVariable(map[string]config.Variable{
		"inherited": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"in_account": config.BoolVariable(true),
			}),
		),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnInheritedSchemasInDatabase(id sdk.AccountObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaValue(config.ObjectVariable(map[string]config.Variable{
		"inherited": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"in_database": config.StringVariable(id.FullyQualifiedName()),
			}),
		),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnInheritedSchemaObjectsInAccount(pluralObjectType sdk.PluralObjectType) *GrantPrivilegesToAccountRoleModel {
	g.WithOnSchemaObjectValue(config.ObjectVariable(map[string]config.Variable{
		"inherited": config.ListVariable(
			config.ObjectVariable(map[string]config.Variable{
				"object_type_plural": config.StringVariable(string(pluralObjectType)),
				"in_account":         config.BoolVariable(true),
			}),
		),
	}))
	return g
}

func (g *GrantPrivilegesToAccountRoleModel) WithOnInheritedSchemaObjectsInDatabase(pluralObjectType sdk.PluralObjectType, id sdk.AccountObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
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

func (g *GrantPrivilegesToAccountRoleModel) WithOnInheritedSchemaObjectsInSchema(pluralObjectType sdk.PluralObjectType, schemaId sdk.DatabaseObjectIdentifier) *GrantPrivilegesToAccountRoleModel {
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
