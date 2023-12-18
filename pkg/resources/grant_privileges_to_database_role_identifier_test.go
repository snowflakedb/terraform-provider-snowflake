package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseGrantPrivilegesToDatabaseRoleId(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier string
		Expected   GrantPrivilegesToDatabaseRoleId
		Error      string
	}{
		{
			Name:       "grant database role on database",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnDatabase|"on-database-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnDatabaseDatabaseRoleGrantKind,
				Data: OnDatabaseGrantData{
					DatabaseName: sdk.NewAccountObjectIdentifier("on-database-name"),
				},
			},
		},
		{
			Name:       "grant database role on schema with schema name",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnSchema|"database-name"."schema-name"`, // TODO: OnSchema OnSchema x2
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: OnSchemaGrantData{
					Kind:       OnSchemaSchemaGrantKind,
					SchemaName: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
				},
			},
		},
		{
			Name:       "grant database role on all schemas in database",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: OnSchemaGrantData{
					Kind:         OnAllSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
				},
			},
		},
		{
			Name:       "grant database role on future schemas in database",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnFutureSchemasInDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaDatabaseRoleGrantKind,
				Data: OnSchemaGrantData{
					Kind:         OnFutureSchemasInDatabaseSchemaGrantKind,
					DatabaseName: sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
				},
			},
		},
		{
			Name:       "grant database role on schema object with on object option",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE|"database-name"."schema-name"."table-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: OnSchemaObjectGrantData{
					Kind: OnObjectSchemaObjectGrantKind,
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeTable,
						Name:       sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in database",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             sdk.Pointer(InDatabaseBulkOperationGrantKind),
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in schema",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: OnSchemaObjectGrantData{
					Kind: OnAllSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             sdk.Pointer(InSchemaBulkOperationGrantKind),
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on future option",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in database",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES|InDatabase|"database-name-123"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             sdk.Pointer(InDatabaseBulkOperationGrantKind),
						Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name-123")),
					},
				},
			},
		},
		{
			Name:       "grant database role on schema object with on all option in schema",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnFuture|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantPrivilegesToDatabaseRoleId{
				DatabaseRoleName: sdk.NewAccountObjectIdentifier("database-name"),
				WithGrantOption:  false,
				Privileges:       []string{"CREATE SCHEMA", "USAGE", "MONITOR"},
				Kind:             OnSchemaObjectDatabaseRoleGrantKind,
				Data: OnSchemaObjectGrantData{
					Kind: OnFutureSchemaObjectGrantKind,
					OnAllOrFuture: &BulkOperationGrantData{
						ObjectNamePlural: "TABLES",
						Kind:             sdk.Pointer(InSchemaBulkOperationGrantKind),
						Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
					},
				},
			},
		},
		{
			Name:       "validation: grant database role not enough parts",
			Identifier: `"database-name"|false`,
			Error:      "database role identifier should hold at least 4 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnDatabase kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnDatabase`,
			Error:      "database role identifier should hold at least 4 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchema kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|OnAllSchemasInDatabase`,
			Error:      "database role identifier should hold at least 6 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchemaObject kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject`,
			Error:      "database role identifier should hold at least 6 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchemaObject kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnObject|TABLE`,
			Error:      "database role identifier should hold 7 parts",
		},
		{
			Name:       "validation: grant database role not enough parts for OnSchemaObject.InDatabase kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|OnAll|TABLES|InDatabase`,
			Error:      "database role identifier should hold 8 parts",
		},
		{
			Name:       "validation: grant database role invalid DatabaseRoleGrantKind kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|some-kind|some-data`,
			Error:      "invalid DatabaseRoleGrantKind: some-kind",
		},
		{
			Name:       "validation: grant database role invalid OnSchemaGrantKind kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchema|some-kind|some-data`,
			Error:      "invalid OnSchemaGrantKind: some-kind",
		},
		{
			Name:       "validation: grant database role invalid OnSchemaObjectGrantKind kind",
			Identifier: `"database-name"|false|CREATE SCHEMA,USAGE,MONITOR|OnSchemaObject|some-kind|some-data`,
			Error:      "invalid OnSchemaObjectGrantKind: some-kind",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			id, err := ParseGrantPrivilegesToDatabaseRoleId(tt.Identifier)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}
