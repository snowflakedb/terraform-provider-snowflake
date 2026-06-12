package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
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

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnSchemaName(schemaFQN string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"schema_name": tfconfig.StringVariable(schemaFQN),
	}))
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnAllSchemasInDatabase(databaseFQN string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all_schemas_in_database": tfconfig.StringVariable(databaseFQN),
	}))
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnFutureSchemasInDatabase(databaseFQN string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"future_schemas_in_database": tfconfig.StringVariable(databaseFQN),
	}))
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnSchemaObjectObject(objectType, objectName string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"object_type": tfconfig.StringVariable(objectType),
		"object_name": tfconfig.StringVariable(objectName),
	}))
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnSchemaObjectAllInDatabase(objectTypePlural, databaseFQN string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(objectTypePlural),
			"in_database":        tfconfig.StringVariable(databaseFQN),
		})),
	}))
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnSchemaObjectAllInSchema(objectTypePlural, schemaFQN string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(objectTypePlural),
			"in_schema":          tfconfig.StringVariable(schemaFQN),
		})),
	}))
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnSchemaObjectFutureInDatabase(objectTypePlural, databaseFQN string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"future": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(objectTypePlural),
			"in_database":        tfconfig.StringVariable(databaseFQN),
		})),
	}))
}

func (g *GrantPrivilegesToDatabaseRoleModel) WithOnSchemaObjectFutureInSchema(objectTypePlural, schemaFQN string) *GrantPrivilegesToDatabaseRoleModel {
	return g.WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"future": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(objectTypePlural),
			"in_schema":          tfconfig.StringVariable(schemaFQN),
		})),
	}))
}
