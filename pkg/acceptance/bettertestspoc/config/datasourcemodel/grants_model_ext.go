package datasourcemodel

import (
	"fmt"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func GrantsOnAccount(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"account": tfconfig.BoolVariable(true),
			}),
		)
}

func GrantsOnAccountObject(
	datasourceName string,
	id sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.Name()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnDatabaseObject(
	datasourceName string,
	id sdk.DatabaseObjectIdentifier,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.FullyQualifiedName()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnSchemaObject(
	datasourceName string,
	id sdk.SchemaObjectIdentifier,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.FullyQualifiedName()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnSchemaObjectWithArguments(
	datasourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	objectType sdk.ObjectType,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable(id.FullyQualifiedName()),
				"object_type": tfconfig.StringVariable(fmt.Sprintf("%s", objectType)),
			}),
		)
}

func GrantsOnMissingObjectType(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_name": tfconfig.StringVariable("DATABASE"),
			}),
		)
}

func GrantsOnEmpty(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			}),
		)
}

func GrantsToDatabaseRole(
	datasourceName string,
	id sdk.DatabaseObjectIdentifier,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database_role": tfconfig.StringVariable(id.FullyQualifiedName()),
			}),
		)
}

func GrantsToUser(
	datasourceName string,
	id sdk.AccountObjectIdentifier,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"user": tfconfig.StringVariable(id.Name()),
			}),
		)
}

func GrantsToAccountRole(
	datasourceName string,
	roleName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"account_role": tfconfig.StringVariable(roleName),
			}),
		)
}

func InheritedGrantsInAccount(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithInheritedGrantsInValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"account": tfconfig.BoolVariable(true),
			}),
		)
}

func GrantsToShare(
	datasourceName string,
	shareName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"share": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"share_name": tfconfig.StringVariable(shareName),
				})),
			}),
		)
}

func InheritedGrantsInDatabase(
	datasourceName string,
	id sdk.AccountObjectIdentifier,
) *GrantsModel {
	return Grants(datasourceName).
		WithInheritedGrantsInValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database": tfconfig.StringVariable(id.Name()),
			}),
		)
}

func GrantsToInvalidEmpty(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			}),
		)
}

func GrantsToInvalidShareNameMissing(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"share": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			}),
		)
}

func GrantsToInvalidDatabaseRoleIdInvalid(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database_role": tfconfig.StringVariable("role"),
			}),
		)
}

func GrantsToInvalidApplicationRoleIdInvalid(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"application_role": tfconfig.StringVariable("role"),
			}),
		)
}

func GrantsOfAccountRole(
	datasourceName string,
	roleName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOfValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"account_role": tfconfig.StringVariable(roleName),
			}),
		)
}

func InheritedGrantsInSchema(
	datasourceName string,
	id sdk.DatabaseObjectIdentifier,
) *GrantsModel {
	return Grants(datasourceName).
		WithInheritedGrantsInValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema": tfconfig.StringVariable(id.FullyQualifiedName()),
			}),
		)
}

func GrantsOfDatabaseRole(
	datasourceName string,
	id sdk.DatabaseObjectIdentifier,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOfValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database_role": tfconfig.StringVariable(id.FullyQualifiedName()),
			}),
		)
}

func GrantsOfShare(
	datasourceName string,
	shareName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOfValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"share": tfconfig.StringVariable(shareName),
			}),
		)
}

func GrantsOfInvalidEmpty(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOfValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			}),
		)
}

func GrantsOfInvalidDatabaseRoleIdInvalid(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOfValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database_role": tfconfig.StringVariable("role"),
			}),
		)
}

func GrantsOfInvalidApplicationRoleIdInvalid(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOfValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"application_role": tfconfig.StringVariable("role"),
			}),
		)
}

func GrantsFutureInDatabase(
	datasourceName string,
	database string,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsInValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database": tfconfig.StringVariable(database),
			}),
		)
}

func GrantsFutureInSchema(
	datasourceName string,
	schemaFQN string,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsInValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema": tfconfig.StringVariable(schemaFQN),
			}),
		)
}

func GrantsFutureInInvalidEmpty(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsInValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			}),
		)
}

func GrantsFutureInInvalidSchemaNotFullyQualified(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsInValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema": tfconfig.StringVariable("schema"),
			}),
		)
}

func GrantsFutureToAccountRole(
	datasourceName string,
	roleName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"account_role": tfconfig.StringVariable(roleName),
			}),
		)
}

func GrantsFutureToDatabaseRole(
	datasourceName string,
	id sdk.DatabaseObjectIdentifier,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database_role": tfconfig.StringVariable(id.FullyQualifiedName()),
			}),
		)
}

func GrantsFutureToInvalidEmpty(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			}),
		)
}

func GrantsFutureToInvalidDatabaseRoleIdInvalid(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithFutureGrantsToValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"database_role": tfconfig.StringVariable("role"),
			}),
		)
}
