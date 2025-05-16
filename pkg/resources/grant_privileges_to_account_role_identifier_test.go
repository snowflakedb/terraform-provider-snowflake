package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestParseGrantPrivilegesToAccountRoleId(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier string
		Expected   GrantPrivilegesToAccountRoleId
		Error      string
	}{
		{
			Name:       "grant account role on account",
			Identifier: `"account-role"|false|false|CREATE DATABASE,CREATE USER|OnAccount`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE DATABASE", "CREATE USER"},
				Kind:            OnAccountAccountRoleGrantKind,
				Data:            new(OnAccountGrantData),
			},
		},
		{
			Name:       "grant account role on account - empty role name",
			Identifier: `|false|false|CREATE DATABASE,CREATE USER|OnAccount`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier(""),
				WithGrantOption: false,
				Privileges:      []string{"CREATE DATABASE", "CREATE USER"},
				Kind:            OnAccountAccountRoleGrantKind,
				Data:            new(OnAccountGrantData),
			},
		},
		{
			Name:       "grant account role on account - always apply with grant option",
			Identifier: `"account-role"|true|true|CREATE DATABASE,CREATE USER|OnAccount`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: true,
				AlwaysApply:     true,
				Privileges:      []string{"CREATE DATABASE", "CREATE USER"},
				Kind:            OnAccountAccountRoleGrantKind,
				Data:            new(OnAccountGrantData),
			},
		},
		{
			Name:       "grant account role on account - all privileges",
			Identifier: `"account-role"|false|false|ALL|OnAccount`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				AllPrivileges:   true,
				Privileges:      nil,
				Kind:            OnAccountAccountRoleGrantKind,
				Data:            new(OnAccountGrantData),
			},
		},
		{
			Name:       "grant account role on account object",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnAccountObject|DATABASE|"database-name"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnAccountObjectAccountRoleGrantKind,
				Data: &OnAccountObjectGrantData{
					ObjectType: sdk.ObjectTypeDatabase,
					ObjectName: sdk.NewAccountObjectIdentifier("database-name"),
				},
			},
		},
		{
			Name:       "grant account role on schema with schema name",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaAccountRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:       OnSchemaSchemaGrantKind,
					SchemaName: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
				},
			},
		},
		{
			Name:       "grant account role on all schemas in database",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaAccountRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnAllSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
				},
			},
		},
		{
			Name:       "grant account role on future schemas in database",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnFutureSchemasInDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaAccountRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnFutureSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
				},
			},
		},
		{
			Name:       "grant account role on schema object with on object option",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE|"database-name"."schema-name"."table-name"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnObjectSchemaObjectGrantKind,
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeTable,
						Name:       sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
					},
				},
			},
		},
		{
			Name:       "grant account role on function",
			Identifier: `"account-role"|false|false|USAGE|OnSchemaObject|OnObject|FUNCTION|"database-name"."schema-name"."function-name"(FLOAT)`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"USAGE"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnObjectSchemaObjectGrantKind,
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeFunction,
						Name:       sdk.NewSchemaObjectIdentifierWithArguments("database-name", "schema-name", "function-name", sdk.DataTypeFloat),
					},
				},
			},
		},
		{
			Name:       "grant account role on function without arguments",
			Identifier: `"account-role"|false|false|USAGE|OnSchemaObject|OnObject|FUNCTION|"database-name"."schema-name"."function-name"()`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"USAGE"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnObjectSchemaObjectGrantKind,
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeFunction,
						Name:       sdk.NewSchemaObjectIdentifierWithArguments("database-name", "schema-name", "function-name", []sdk.DataType{}...),
					},
				},
			},
		},
		{
			Name:       "grant account role on schema object with on all option",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
					},
				},
			},
		},
		{
			Name:       "grant account role on schema object with on all option in database",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InDatabaseBulkOperationGrantKind,
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
					},
				},
			},
		},
		{
			Name:       "grant account role on schema object with on all option in schema",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InSchemaBulkOperationGrantKind,
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
		},
		{
			Name:       "grant account role on schema object with on future option",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
					},
				},
			},
		},
		{
			Name:       "grant account role on schema object with on all option in database",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES|InDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InDatabaseBulkOperationGrantKind,
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
					},
				},
			},
		},
		{
			Name:       "grant account role on schema object with on all option in schema",
			Identifier: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             InSchemaBulkOperationGrantKind,
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
		},
		{
			Name:       "validation: grant account role not enough parts",
			Identifier: `"database-name"."role-name"|false|false`,
			Error:      "account role identifier should hold at least 5 parts",
		},
		{
			Name:       "validation: grant account role not enough parts for OnAccountObject kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnAccountObject`,
			Error:      `account role identifier should hold at least 7 parts "<role_name>|<with_grant_option>|<always_apply>|<privileges>|OnAccountObject|<object_type>|<object_name>"`,
		},
		{
			Name:       "validation: grant account role not enough parts for OnSchema kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase`,
			Error:      "account role identifier should hold at least 7 parts",
		},
		{
			Name:       "validation: grant account role not enough parts for OnSchemaObject kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject`,
			Error:      "account role identifier should hold at least 7 parts",
		},
		{
			Name:       "validation: grant account role not enough parts for OnSchemaObject kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE`,
			Error:      "account role identifier should hold 8 parts",
		},
		{
			Name:       "validation: grant account role not enough parts for OnSchemaObject.InDatabase kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase`,
			Error:      "account role identifier should hold 9 parts",
		},
		{
			Name:       "validation: grant account role invalid AccountRoleGrantKind kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|some-kind|some-data`,
			Error:      "invalid AccountRoleGrantKind: some-kind",
		},
		{
			Name:       "validation: grant account role invalid OnSchemaGrantKind kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|some-kind|some-data`,
			Error:      "invalid OnSchemaGrantKind: some-kind",
		},
		{
			Name:       "validation: grant account role invalid OnSchemaObjectGrantKind kind",
			Identifier: `"role-name"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|some-kind|some-data`,
			Error:      "invalid OnSchemaObjectGrantKind: some-kind",
		},
		{
			Name:       "validation: grant account role empty privileges",
			Identifier: `"account-role"|false|false||OnAccount`,
			Error:      `invalid Privileges value: , should be either a comma separated list of privileges or "ALL" / "ALL PRIVILEGES" for all privileges`,
		},
		{
			Name:       "validation: grant account role empty with grant option",
			Identifier: `"account-role"||false|ALL PRIVILEGES|OnAccount`,
			Error:      `invalid WithGrantOption value: , should be either "true" or "false"`,
		},
		{
			Name:       "validation: grant account role empty always apply",
			Identifier: `"account-role"|false||ALL PRIVILEGES|OnAccount`,
			Error:      `invalid AlwaysApply value: , should be either "true" or "false"`,
		},
		{
			Name:       "validation: account role empty type",
			Identifier: `"account-role"|false|false|ALL PRIVILEGES||"on-database-name"`,
			Error:      "invalid AccountRoleGrantKind: ",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			id, err := ParseGrantPrivilegesToAccountRoleId(tt.Identifier)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}

func TestGrantPrivilegesToAccountRoleIdString(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier GrantPrivilegesToAccountRoleId
		Expected   string
		Error      string
	}{
		{
			Name: "grant account role on account",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: true,
				AllPrivileges:   true,
				Kind:            OnAccountAccountRoleGrantKind,
				AlwaysApply:     true,
				Data:            new(OnAccountGrantData),
			},
			Expected: `"account-role"|true|true|ALL|OnAccount`,
		},
		{
			Name: "grant account role on account object",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: true,
				AllPrivileges:   true,
				Kind:            OnAccountObjectAccountRoleGrantKind,
				AlwaysApply:     true,
				Data: &OnAccountObjectGrantData{
					ObjectType: sdk.ObjectTypeDatabase,
					ObjectName: sdk.NewAccountObjectIdentifier("database-name"),
				},
			},
			Expected: `"account-role"|true|true|ALL|OnAccountObject|DATABASE|"database-name"`,
		},
		{
			Name: "grant account role on schema on schema",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaAccountRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:       OnSchemaSchemaGrantKind,
					SchemaName: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
				},
			},
			Expected: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnSchema|"database-name"."schema-name"`,
		},
		{
			Name: "grant account role on all schemas in database",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaAccountRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnAllSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
				},
			},
			Expected: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase|"database-name"`,
		},
		{
			Name: "grant account role on future schemas in database",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaAccountRoleGrantKind,
				Data: &OnSchemaGrantData{
					Kind:         OnFutureSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
				},
			},
			Expected: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnFutureSchemasInDatabase|"database-name"`,
		},
		{
			Name: "grant account role on schema object on object",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnObjectSchemaObjectGrantKind,
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeTable,
						Name:       sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
					},
				},
			},
			Expected: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE|"database-name"."schema-name"."table-name"`,
		},
		{
			Name: "grant account role on schema object on all tables in database",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: sdk.PluralObjectTypeTables,
						Kind:             InDatabaseBulkOperationGrantKind,
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
					},
				},
			},
			Expected: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase|"database-name"`,
		},
		{
			Name: "grant account role on schema object on all tables in schema",
			Identifier: GrantPrivilegesToAccountRoleId{
				RoleName:        sdk.NewAccountObjectIdentifier("account-role"),
				WithGrantOption: false,
				Privileges:      []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:            OnSchemaObjectAccountRoleGrantKind,
				Data: &OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: sdk.PluralObjectTypeTables,
						Kind:             InSchemaBulkOperationGrantKind,
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
			Expected: `"account-role"|false|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InSchema|"database-name"."schema-name"`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Identifier.String())
		})
	}
}
