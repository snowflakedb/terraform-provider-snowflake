package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

type DatabaseRoleGrantKind string

const (
	OnDatabaseDatabaseRoleGrantKind     DatabaseRoleGrantKind = "OnDatabase"
	OnSchemaDatabaseRoleGrantKind       DatabaseRoleGrantKind = "OnSchema"
	OnSchemaObjectDatabaseRoleGrantKind DatabaseRoleGrantKind = "OnSchemaObject"
)

type GrantPrivilegesToDatabaseRoleId struct {
	DatabaseRoleName sdk.DatabaseObjectIdentifier
	WithGrantOption  bool
	AlwaysApply      bool
	AllPrivileges    bool
	Privileges       []string
	Kind             DatabaseRoleGrantKind
	Data             fmt.Stringer
}

func (g *GrantPrivilegesToDatabaseRoleId) String() string {
	var parts []string
	parts = append(parts, g.DatabaseRoleName.FullyQualifiedName())
	parts = append(parts, strconv.FormatBool(g.WithGrantOption))
	parts = append(parts, strconv.FormatBool(g.AlwaysApply))
	if g.AllPrivileges {
		parts = append(parts, "ALL")
	} else {
		parts = append(parts, strings.Join(g.Privileges, ","))
	}
	parts = append(parts, string(g.Kind))
	parts = append(parts, g.Data.String())
	return helpers.EncodeResourceIdentifier(parts...)
}

type OnDatabaseGrantData struct {
	DatabaseName sdk.AccountObjectIdentifier
}

func (d *OnDatabaseGrantData) String() string {
	return d.DatabaseName.FullyQualifiedName()
}

func ParseGrantPrivilegesToDatabaseRoleId(id string) (GrantPrivilegesToDatabaseRoleId, error) {
	var databaseRoleId GrantPrivilegesToDatabaseRoleId

	parts := helpers.ParseResourceIdentifier(id)
	if len(parts) < 6 {
		return databaseRoleId, sdk.NewError(`database role identifier should hold at least 6 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_data>"`)
	}

	roleId, err := sdk.ParseDatabaseObjectIdentifier(parts[0])
	if err != nil {
		return databaseRoleId, err
	}
	databaseRoleId.DatabaseRoleName = roleId

	if parts[1] != "false" && parts[1] != "true" {
		return databaseRoleId, sdk.NewError(fmt.Sprintf(`invalid WithGrantOption value: %s, should be either "true" or "false"`, parts[1]))
	}
	databaseRoleId.WithGrantOption = parts[1] == "true"

	if parts[2] != "false" && parts[2] != "true" {
		return databaseRoleId, sdk.NewError(fmt.Sprintf(`invalid AlwaysApply value: %s, should be either "true" or "false"`, parts[2]))
	}
	databaseRoleId.AlwaysApply = parts[2] == "true"

	privileges := strings.Split(parts[3], ",")
	if len(privileges) == 0 || (len(privileges) == 1 && privileges[0] == "") {
		return databaseRoleId, sdk.NewError(fmt.Sprintf(`invalid Privileges value: %s, should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all privileges`, parts[3]))
	}
	if len(privileges) == 1 && (privileges[0] == "ALL" || privileges[0] == "ALL PRIVILEGES") {
		databaseRoleId.AllPrivileges = true
	} else {
		databaseRoleId.Privileges = privileges
	}

	databaseRoleId.Kind = DatabaseRoleGrantKind(parts[4])

	switch databaseRoleId.Kind {
	case OnDatabaseDatabaseRoleGrantKind:
		databaseId, err := sdk.ParseAccountObjectIdentifier(parts[5])
		if err != nil {
			return databaseRoleId, err
		}
		databaseRoleId.Data = &OnDatabaseGrantData{
			DatabaseName: databaseId,
		}
	case OnSchemaDatabaseRoleGrantKind:
		if len(parts) < 7 {
			return databaseRoleId, sdk.NewError(`database role identifier should hold at least 7 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_on_schema_type>|<on_schema_grant_data>..."`)
		}
		onSchemaGrantData := OnSchemaGrantData{
			Kind: OnSchemaGrantKind(parts[5]),
		}
		switch onSchemaGrantData.Kind {
		case OnSchemaSchemaGrantKind:
			schemaId, err := sdk.ParseDatabaseObjectIdentifier(parts[6])
			if err != nil {
				return databaseRoleId, err
			}
			onSchemaGrantData.SchemaName = sdk.Pointer(schemaId)
		case OnAllSchemasInDatabaseSchemaGrantKind, OnFutureSchemasInDatabaseSchemaGrantKind:
			databaseId, err := sdk.ParseAccountObjectIdentifier(parts[6])
			if err != nil {
				return databaseRoleId, err
			}
			onSchemaGrantData.DatabaseName = sdk.Pointer(databaseId)
		default:
			return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaGrantKind: %s", onSchemaGrantData.Kind))
		}
		databaseRoleId.Data = &onSchemaGrantData
	case OnSchemaObjectDatabaseRoleGrantKind:
		if len(parts) < 7 {
			return databaseRoleId, sdk.NewError(`database role identifier should hold at least 7 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_on_schema_object_type>|<on_schema_object_grant_data>..."`)
		}
		onSchemaObjectGrantData := OnSchemaObjectGrantData{
			Kind: OnSchemaObjectGrantKind(parts[5]),
		}
		switch onSchemaObjectGrantData.Kind {
		case OnObjectSchemaObjectGrantKind:
			if len(parts) != 8 {
				return databaseRoleId, sdk.NewError(`database role identifier should hold 8 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnObject|<object_type>|<object_name>"`)
			}
			objectType := sdk.ObjectType(parts[6])
			var id sdk.ObjectIdentifier
			// TODO(SNOW-1569535): use a mapper from object type to parsing function
			if objectType.IsWithArguments() {
				id, err = sdk.ParseSchemaObjectIdentifierWithArguments(parts[7])
				if err != nil {
					return databaseRoleId, err
				}
			} else {
				id, err = sdk.ParseSchemaObjectIdentifier(parts[7])
				if err != nil {
					return databaseRoleId, err
				}
			}
			onSchemaObjectGrantData.Object = &sdk.Object{
				ObjectType: objectType,
				Name:       id,
			}
		case OnAllSchemaObjectGrantKind, OnFutureSchemaObjectGrantKind:
			bulkOperationGrantData := &BulkOperationGrantData{
				ObjectNamePlural: sdk.PluralObjectType(parts[6]),
			}
			if len(parts) > 7 {
				if len(parts) != 9 {
					return databaseRoleId, sdk.NewError(`database role identifier should hold 9 parts "<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|On[All or Future]|<object_type_plural>|In[Database or Schema]|<identifier>"`)
				}
				bulkOperationGrantData.Kind = BulkOperationGrantKind(parts[7])
				switch bulkOperationGrantData.Kind {
				case InDatabaseBulkOperationGrantKind:
					databaseId, err := sdk.ParseAccountObjectIdentifier(parts[8])
					if err != nil {
						return databaseRoleId, err
					}
					bulkOperationGrantData.Database = sdk.Pointer(databaseId)
				case InSchemaBulkOperationGrantKind:
					schemaId, err := sdk.ParseDatabaseObjectIdentifier(parts[8])
					if err != nil {
						return databaseRoleId, err
					}
					bulkOperationGrantData.Schema = sdk.Pointer(schemaId)
				default:
					return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid BulkOperationGrantKind: %s", bulkOperationGrantData.Kind))
				}
			}
			onSchemaObjectGrantData.OnAllOrFuture = bulkOperationGrantData
		default:
			return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaObjectGrantKind: %s", onSchemaObjectGrantData.Kind))
		}
		databaseRoleId.Data = &onSchemaObjectGrantData
	default:
		return databaseRoleId, sdk.NewError(fmt.Sprintf("invalid DatabaseRoleGrantKind: %s", databaseRoleId.Kind))
	}

	return databaseRoleId, nil
}
