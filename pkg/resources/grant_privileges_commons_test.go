package resources

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestComputeInheritedPrivileges(t *testing.T) {
	accountRoleName := "test-role"
	accountRoleId := sdk.NewAccountObjectIdentifier(accountRoleName)
	databaseRoleName := "my_role"

	onAccountWarehouses := &OnAccountObjectInheritedGrantData{ObjectNamePlural: sdk.PluralObjectTypeWarehouses}
	onSchemaInAccount := &OnSchemaInheritedGrantData{Kind: InAccountInheritedContainerKind}
	onSchemaObjectTablesInDatabase := &OnSchemaObjectInheritedGrantData{ObjectNamePlural: sdk.PluralObjectTypeTables, Kind: InDatabaseInheritedContainerKind, DatabaseName: new(sdk.NewAccountObjectIdentifier("my_db"))}
	onSchemaObjectTablesInSchema := &OnSchemaObjectInheritedGrantData{ObjectNamePlural: sdk.PluralObjectTypeTables, Kind: InSchemaInheritedContainerKind, SchemaName: new(sdk.NewDatabaseObjectIdentifier("my_db", "my_schema"))}

	testCases := []struct {
		Name        string
		Data        fmt.Stringer
		RoleName    string
		GranteeType sdk.ObjectType
		Privileges  []string
		Grants      []sdk.Grant
		Strict      bool
		Expected    []string
	}{
		{
			Name:        "on account object - matching inherited grant is returned",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, accountRoleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE"},
		},
		{
			Name:        "non-inherited grants are ignored - is_inherited is nil",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				{Privilege: "USAGE", GrantedOn: sdk.ObjectTypeWarehouse, GrantedTo: sdk.ObjectTypeRole, GranteeName: accountRoleId, IsInherited: nil},
			},
			Expected: nil,
		},
		{
			Name:        "non-inherited grants are ignored - is_inherited is false",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				{Privilege: "USAGE", GrantedOn: sdk.ObjectTypeWarehouse, GrantedTo: sdk.ObjectTypeRole, GranteeName: accountRoleId, IsInherited: new(false)},
			},
			Expected: nil,
		},
		{
			Name:        "grants to a different role are ignored",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, "other-role", sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: nil,
		},
		{
			Name:        "grants on a different object type are ignored",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeDatabase, accountRoleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: nil,
		},
		{
			Name:        "grants from a different container are ignored",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, accountRoleName, sdk.GrantInheritedFromDatabase, "my_db", ""),
			},
			Expected: nil,
		},
		{
			Name:        "inherited grant without inherited_from is ignored",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, accountRoleName, "", "", ""),
			},
			Expected: nil,
		},
		{
			Name:        "non-strict - external privileges not in config are ignored",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Strict:      false,
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, accountRoleName, sdk.GrantInheritedFromAccount, "", ""),
				inheritedGrant("MONITOR", sdk.ObjectTypeWarehouse, accountRoleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE"},
		},
		{
			Name:        "strict - external privileges are detected",
			Data:        onAccountWarehouses,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Strict:      true,
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeWarehouse, accountRoleName, sdk.GrantInheritedFromAccount, "", ""),
				inheritedGrant("MONITOR", sdk.ObjectTypeWarehouse, accountRoleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE", "MONITOR"},
		},
		{
			Name:        "on schema object in database - different database is ignored",
			Data:        onSchemaObjectTablesInDatabase,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"SELECT"},
			Grants: []sdk.Grant{
				inheritedGrant("SELECT", sdk.ObjectTypeTable, accountRoleName, sdk.GrantInheritedFromDatabase, "other_db", ""),
			},
			Expected: nil,
		},
		{
			Name:        "on schema object in schema - different schema is ignored",
			Data:        onSchemaObjectTablesInSchema,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"SELECT"},
			Grants: []sdk.Grant{
				inheritedGrant("SELECT", sdk.ObjectTypeTable, accountRoleName, sdk.GrantInheritedFromSchema, "my_db", "other_schema"),
			},
			Expected: nil,
		},
		{
			Name:        "on schema in account - matching schemas inherited from account",
			Data:        onSchemaInAccount,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"USAGE"},
			Grants: []sdk.Grant{
				inheritedGrant("USAGE", sdk.ObjectTypeSchema, accountRoleName, sdk.GrantInheritedFromAccount, "", ""),
			},
			Expected: []string{"USAGE"},
		},
		{
			Name:        "account role scope ignores grants to a database role",
			Data:        onSchemaObjectTablesInDatabase,
			RoleName:    accountRoleName,
			GranteeType: sdk.ObjectTypeRole,
			Privileges:  []string{"SELECT"},
			Grants: []sdk.Grant{
				inheritedDatabaseRoleGrant("SELECT", sdk.ObjectTypeTable, sdk.NewDatabaseObjectIdentifier("my_db", accountRoleName), sdk.GrantInheritedFromDatabase, "my_db", ""),
			},
			Expected: nil,
		},
		{
			Name:        "database role scope ignores grants to an account role",
			Data:        onSchemaObjectTablesInDatabase,
			RoleName:    databaseRoleName,
			GranteeType: sdk.ObjectTypeDatabaseRole,
			Privileges:  []string{"SELECT"},
			Grants: []sdk.Grant{
				inheritedGrant("SELECT", sdk.ObjectTypeTable, databaseRoleName, sdk.GrantInheritedFromDatabase, "my_db", ""),
			},
			Expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := computeInheritedPrivileges(tc.Data, tc.RoleName, tc.GranteeType, tc.Privileges, tc.Grants, tc.Strict)
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

func inheritedDatabaseRoleGrant(privilege string, grantedOn sdk.ObjectType, grantee sdk.DatabaseObjectIdentifier, inheritedFrom sdk.GrantInheritedFrom, inheritedFromDatabase string, inheritedFromSchema string) sdk.Grant {
	grant := sdk.Grant{
		Privilege:   privilege,
		GrantedOn:   grantedOn,
		GrantedTo:   sdk.ObjectTypeDatabaseRole,
		GranteeName: grantee,
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
