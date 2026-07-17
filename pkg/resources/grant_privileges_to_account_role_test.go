package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestComputeInheritedAccountRolePrivileges(t *testing.T) {
	roleName := "test-role"
	roleId := sdk.NewAccountObjectIdentifier(roleName)

	onAccountWarehousesId := func(privileges ...string) GrantPrivilegesToAccountRoleId {
		return GrantPrivilegesToAccountRoleId{
			RoleName:   roleId,
			Privileges: privileges,
			Kind:       OnAccountObjectInheritedAccountRoleGrantKind,
			Data:       &OnAccountObjectInheritedGrantData{ObjectNamePlural: sdk.PluralObjectTypeWarehouses},
		}
	}
	onSchemaObjectTablesInDatabaseId := func(privileges ...string) GrantPrivilegesToAccountRoleId {
		return GrantPrivilegesToAccountRoleId{
			RoleName:   roleId,
			Privileges: privileges,
			Kind:       OnSchemaObjectInheritedAccountRoleGrantKind,
			Data: &OnSchemaObjectInheritedGrantData{
				ObjectNamePlural: sdk.PluralObjectTypeTables,
				Kind:             InDatabaseInheritedContainerKind,
				DatabaseName:     new(sdk.NewAccountObjectIdentifier("my_db")),
			},
		}
	}
	onSchemaObjectTablesInSchemaId := func(privileges ...string) GrantPrivilegesToAccountRoleId {
		return GrantPrivilegesToAccountRoleId{
			RoleName:   roleId,
			Privileges: privileges,
			Kind:       OnSchemaObjectInheritedAccountRoleGrantKind,
			Data: &OnSchemaObjectInheritedGrantData{
				ObjectNamePlural: sdk.PluralObjectTypeTables,
				Kind:             InSchemaInheritedContainerKind,
				SchemaName:       new(sdk.NewDatabaseObjectIdentifier("my_db", "my_schema")),
			},
		}
	}
	onSchemaInAccountId := func(privileges ...string) GrantPrivilegesToAccountRoleId {
		return GrantPrivilegesToAccountRoleId{
			RoleName:   roleId,
			Privileges: privileges,
			Kind:       OnSchemaInheritedAccountRoleGrantKind,
			Data:       &OnSchemaInheritedGrantData{Kind: InAccountInheritedContainerKind},
		}
	}

	testCases := []struct {
		Name     string
		Id       GrantPrivilegesToAccountRoleId
		Grants   []sdk.Grant
		Strict   bool
		Expected []string
	}{
		{
			Name: "on account object - matching inherited grant is returned",
			Id:   onAccountWarehousesId("USAGE"),
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, roleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE"},
		},
		{
			Name: "non-inherited grants are ignored - is_inherited is nil",
			Id:   onAccountWarehousesId("USAGE"),
			Grants: []sdk.Grant{
				{Privilege: "USAGE", GrantedOn: sdk.ObjectTypeWarehouse, GrantedTo: sdk.ObjectTypeRole, GranteeName: roleId, IsInherited: nil},
			},
			Expected: nil,
		},
		{
			Name: "non-inherited grants are ignored - is_inherited is false",
			Id:   onAccountWarehousesId("USAGE"),
			Grants: []sdk.Grant{
				{Privilege: "USAGE", GrantedOn: sdk.ObjectTypeWarehouse, GrantedTo: sdk.ObjectTypeRole, GranteeName: roleId, IsInherited: new(false)},
			},
			Expected: nil,
		},
		{
			Name: "grants to a different role are ignored",
			Id:   onAccountWarehousesId("USAGE"),
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, "other-role", sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: nil,
		},
		{
			Name: "grants on a different object type are ignored",
			Id:   onAccountWarehousesId("USAGE"),
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeDatabase, roleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: nil,
		},
		{
			Name: "grants from a different container are ignored",
			Id:   onAccountWarehousesId("USAGE"),
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, roleName, sdk.GrantInheritedFromDatabase, "my_db", ""),
			},
			Expected: nil,
		},
		{
			Name: "inherited grant without inherited_from is ignored",
			Id:   onAccountWarehousesId("USAGE"),
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, roleName, "", "", ""),
			},
			Expected: nil,
		},
		{
			Name:   "non-strict - external privileges not in config are ignored",
			Id:     onAccountWarehousesId("USAGE"),
			Strict: false,
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, roleName, sdk.GrantInheritedFromAccount, "", ""),
				inheritedGrant("MONITOR", sdk.ObjectTypeWarehouse, roleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE"},
		},
		{
			Name:   "strict - external privileges are detected",
			Id:     onAccountWarehousesId("USAGE"),
			Strict: true,
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, roleName, sdk.GrantInheritedFromAccount, "", ""),
				inheritedGrant("MONITOR", sdk.ObjectTypeWarehouse, roleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE", "MONITOR"},
		},
		{
			Name: "on schema object in database - different database is ignored",
			Id:   onSchemaObjectTablesInDatabaseId("SELECT"),
			Grants: []sdk.Grant{
				inheritedGrant("SELECT", sdk.ObjectTypeTable, roleName, sdk.GrantInheritedFromDatabase, "other_db", ""),
			},
			Expected: nil,
		},
		{
			Name: "on schema object in schema - different schema is ignored",
			Id:   onSchemaObjectTablesInSchemaId("SELECT"),
			Grants: []sdk.Grant{
				inheritedGrant("SELECT", sdk.ObjectTypeTable, roleName, sdk.GrantInheritedFromSchema, "my_db", "other_schema"),
			},
			Expected: nil,
		},
		{
			Name: "on schema in account - matching schemas inherited from account",
			Id:   onSchemaInAccountId("USAGE"),
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeSchema, roleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := computeInheritedAccountRolePrivileges(tc.Id, tc.Grants, tc.Strict)
			assert.ElementsMatch(t, tc.Expected, actual)
		})
	}
}

func inheritedGrant(privilege string, grantedOn sdk.ObjectType, grantee string, inheritedFrom sdk.GrantInheritedFrom, inheritedFromDatabase string, inheritedFromSchema string) sdk.Grant {
	grant := sdk.Grant{
		Privilege:   privilege,
		GrantedOn:   grantedOn,
		GrantedTo:   sdk.ObjectTypeRole,
		GranteeName: sdk.NewAccountObjectIdentifier(grantee),
		IsInherited: new(true),
	}
	if inheritedFrom != "" {
		grant.InheritedFrom = &inheritedFrom
	}
	if inheritedFromDatabase != "" {
		grant.InheritedFromDatabase = &inheritedFromDatabase
	}
	if inheritedFromSchema != "" {
		grant.InheritedFromSchema = &inheritedFromSchema
	}
	return grant
}
