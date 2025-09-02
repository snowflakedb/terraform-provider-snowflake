package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type AccountRoleGrantKind string

const (
	OnAccountAccountRoleGrantKind       AccountRoleGrantKind = "OnAccount"
	OnAccountObjectAccountRoleGrantKind AccountRoleGrantKind = "OnAccountObject"
	OnSchemaAccountRoleGrantKind        AccountRoleGrantKind = "OnSchema"
	OnSchemaObjectAccountRoleGrantKind  AccountRoleGrantKind = "OnSchemaObject"
)

type OnAccountGrantData struct{}

func (d *OnAccountGrantData) String() string {
	return ""
}

type OnAccountObjectGrantData struct {
	ObjectType sdk.ObjectType
	ObjectName sdk.AccountObjectIdentifier
}

func (d *OnAccountObjectGrantData) String() string {
	return helpers.EncodeResourceIdentifier(d.ObjectType.String(), d.ObjectName.FullyQualifiedName())
}

type GrantPrivilegesToAccountRoleId struct {
	RoleName        sdk.AccountObjectIdentifier
	WithGrantOption bool
	AlwaysApply     bool
	AllPrivileges   bool
	Privileges      []string
	Kind            AccountRoleGrantKind
	Data            fmt.Stringer
}

func NewGrantPrivilegesToAccountRoleIdOnAccount(roleName sdk.AccountObjectIdentifier, grantOption bool, alwaysApply bool, allPrivileges bool, privileges ...string) GrantPrivilegesToAccountRoleId {
	return GrantPrivilegesToAccountRoleId{
		RoleName:        roleName,
		WithGrantOption: grantOption,
		AlwaysApply:     alwaysApply,
		AllPrivileges:   allPrivileges,
		Privileges:      privileges,
		Kind:            OnAccountAccountRoleGrantKind,
		Data:            new(OnAccountGrantData),
	}
}

func NewGrantPrivilegesToAccountRoleIdOnAccountObject(roleName sdk.AccountObjectIdentifier, grantOption bool, alwaysApply bool, allPrivileges bool, grantedOn sdk.ObjectType, objectName sdk.AccountObjectIdentifier, privileges ...string) GrantPrivilegesToAccountRoleId {
	return GrantPrivilegesToAccountRoleId{
		RoleName:        roleName,
		WithGrantOption: grantOption,
		AlwaysApply:     alwaysApply,
		AllPrivileges:   allPrivileges,
		Privileges:      privileges,
		Kind:            OnAccountObjectAccountRoleGrantKind,
		Data: &OnAccountObjectGrantData{
			ObjectType: grantedOn,
			ObjectName: objectName,
		},
	}
}

func NewGrantPrivilegesToAccountRoleIdOnSchemaOnSchema(roleName sdk.AccountObjectIdentifier, grantOption bool, alwaysApply bool, allPrivileges bool, schemaName sdk.DatabaseObjectIdentifier, privileges ...string) GrantPrivilegesToAccountRoleId {
	return GrantPrivilegesToAccountRoleId{
		RoleName:        roleName,
		WithGrantOption: grantOption,
		AlwaysApply:     alwaysApply,
		AllPrivileges:   allPrivileges,
		Privileges:      privileges,
		Kind:            OnSchemaAccountRoleGrantKind,
		Data: &OnSchemaGrantData{
			Kind:       OnSchemaSchemaGrantKind,
			SchemaName: &schemaName,
		},
	}
}

func NewGrantPrivilegesToAccountRoleIdOnSchemaObjectOnObject(roleName sdk.AccountObjectIdentifier, grantOption bool, alwaysApply bool, allPrivileges bool, objectType sdk.ObjectType, objectName sdk.SchemaObjectIdentifier, privileges ...string) GrantPrivilegesToAccountRoleId {
	return GrantPrivilegesToAccountRoleId{
		RoleName:        roleName,
		WithGrantOption: grantOption,
		AlwaysApply:     alwaysApply,
		AllPrivileges:   allPrivileges,
		Privileges:      privileges,
		Kind:            OnSchemaObjectAccountRoleGrantKind,
		Data: &OnSchemaObjectGrantData{
			Kind: OnObjectSchemaObjectGrantKind,
			Object: &sdk.Object{
				ObjectType: objectType,
				Name:       objectName,
			},
		},
	}
}

func (g *GrantPrivilegesToAccountRoleId) String() string {
	var parts []string
	parts = append(parts, g.RoleName.FullyQualifiedName())
	parts = append(parts, strconv.FormatBool(g.WithGrantOption))
	parts = append(parts, strconv.FormatBool(g.AlwaysApply))
	if g.AllPrivileges {
		parts = append(parts, "ALL")
	} else {
		parts = append(parts, strings.Join(g.Privileges, ","))
	}
	parts = append(parts, string(g.Kind))
	data := g.Data.String()
	if len(data) > 0 {
		parts = append(parts, data)
	}
	return helpers.EncodeResourceIdentifier(parts...)
}

func ParseGrantPrivilegesToAccountRoleId(id string) (GrantPrivilegesToAccountRoleId, error) {
	var accountRoleId GrantPrivilegesToAccountRoleId

	parts := helpers.ParseResourceIdentifier(id)
	if len(parts) < 5 {
		return accountRoleId, sdk.NewError(`account role identifier should hold at least 5 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>"`)
	}

	roleId, err := sdk.ParseAccountObjectIdentifier(parts[0])
	if err != nil {
		return accountRoleId, err
	}
	accountRoleId.RoleName = roleId

	if parts[1] != "false" && parts[1] != "true" {
		return accountRoleId, sdk.NewError(fmt.Sprintf(`invalid WithGrantOption value: %s, should be either "true" or "false"`, parts[1]))
	}
	accountRoleId.WithGrantOption = parts[1] == "true"

	if parts[2] != "false" && parts[2] != "true" {
		return accountRoleId, sdk.NewError(fmt.Sprintf(`invalid AlwaysApply value: %s, should be either "true" or "false"`, parts[2]))
	}
	accountRoleId.AlwaysApply = parts[2] == "true"

	privileges := strings.Split(parts[3], ",")
	if len(privileges) == 0 || (len(privileges) == 1 && privileges[0] == "") {
		return accountRoleId, sdk.NewError(fmt.Sprintf(`invalid Privileges value: %s, should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all privileges`, parts[3]))
	}
	if len(privileges) == 1 && (privileges[0] == "ALL" || privileges[0] == "ALL PRIVILEGES") {
		accountRoleId.AllPrivileges = true
	} else {
		accountRoleId.Privileges = privileges
	}

	accountRoleId.Kind = AccountRoleGrantKind(parts[4])
	switch accountRoleId.Kind {
	case OnAccountAccountRoleGrantKind:
		accountRoleId.Data = new(OnAccountGrantData)
	case OnAccountObjectAccountRoleGrantKind:
		if len(parts) != 7 {
			return accountRoleId, sdk.NewError(`account role identifier should hold at least 7 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnAccountObject|<object_type>|<object_name>"`)
		}
		objectId, err := sdk.ParseAccountObjectIdentifier(parts[6])
		if err != nil {
			return accountRoleId, err
		}
		accountRoleId.Data = &OnAccountObjectGrantData{
			ObjectType: sdk.ObjectType(parts[5]),
			ObjectName: objectId,
		}
	case OnSchemaAccountRoleGrantKind:
		if len(parts) < 7 {
			return accountRoleId, sdk.NewError(`account role identifier should hold at least 7 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchema|<grant_on_schema_type>|<on_schema_grant_data>..."`)
		}
		onSchemaGrantData := OnSchemaGrantData{
			Kind: OnSchemaGrantKind(parts[5]),
		}
		switch onSchemaGrantData.Kind {
		case OnSchemaSchemaGrantKind:
			schemaId, err := sdk.ParseDatabaseObjectIdentifier(parts[6])
			if err != nil {
				return accountRoleId, err
			}
			onSchemaGrantData.SchemaName = sdk.Pointer(schemaId)
		case OnAllSchemasInDatabaseSchemaGrantKind, OnFutureSchemasInDatabaseSchemaGrantKind:
			databaseId, err := sdk.ParseAccountObjectIdentifier(parts[6])
			if err != nil {
				return accountRoleId, err
			}
			onSchemaGrantData.DatabaseName = sdk.Pointer(databaseId)
		default:
			return accountRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaGrantKind: %s", onSchemaGrantData.Kind))
		}
		accountRoleId.Data = &onSchemaGrantData
	case OnSchemaObjectAccountRoleGrantKind:
		if len(parts) < 7 {
			return accountRoleId, sdk.NewError(`account role identifier should hold at least 7 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|<grant_on_schema_object_type>|<on_schema_object_grant_data>..."`)
		}
		onSchemaObjectGrantData := OnSchemaObjectGrantData{
			Kind: OnSchemaObjectGrantKind(parts[5]),
		}
		switch onSchemaObjectGrantData.Kind {
		case OnObjectSchemaObjectGrantKind:
			if len(parts) != 8 {
				return accountRoleId, sdk.NewError(`account role identifier should hold 8 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|OnObject|<object_type>|<object_name>"`)
			}
			objectType := sdk.ObjectType(parts[6])
			var id sdk.ObjectIdentifier
			// TODO(SNOW-1569535): use a mapper from object type to parsing function
			if objectType.IsWithArguments() {
				id, err = sdk.ParseSchemaObjectIdentifierWithArguments(parts[7])
				if err != nil {
					return accountRoleId, err
				}
			} else {
				id, err = sdk.ParseSchemaObjectIdentifier(parts[7])
				if err != nil {
					return accountRoleId, err
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
					return accountRoleId, sdk.NewError(`account role identifier should hold 9 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnSchemaObject|On[All or Future]|<object_type_plural>|In[Database or Schema]|<identifier>"`)
				}
				bulkOperationGrantData.Kind = BulkOperationGrantKind(parts[7])
				switch bulkOperationGrantData.Kind {
				case InDatabaseBulkOperationGrantKind:
					databaseId, err := sdk.ParseAccountObjectIdentifier(parts[8])
					if err != nil {
						return accountRoleId, err
					}
					bulkOperationGrantData.Database = sdk.Pointer(databaseId)
				case InSchemaBulkOperationGrantKind:
					schemaId, err := sdk.ParseDatabaseObjectIdentifier(parts[8])
					if err != nil {
						return accountRoleId, err
					}
					bulkOperationGrantData.Schema = sdk.Pointer(schemaId)
				default:
					return accountRoleId, sdk.NewError(fmt.Sprintf("invalid BulkOperationGrantKind: %s", bulkOperationGrantData.Kind))
				}
			}
			onSchemaObjectGrantData.OnAllOrFuture = bulkOperationGrantData
		default:
			return accountRoleId, sdk.NewError(fmt.Sprintf("invalid OnSchemaObjectGrantKind: %s", onSchemaObjectGrantData.Kind))
		}
		accountRoleId.Data = &onSchemaObjectGrantData
	default:
		return accountRoleId, sdk.NewError(fmt.Sprintf("invalid AccountRoleGrantKind: %s", accountRoleId.Kind))
	}

	return accountRoleId, nil
}
