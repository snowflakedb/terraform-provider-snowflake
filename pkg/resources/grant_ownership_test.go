package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOnObjectIdentifier(t *testing.T) {
	testCases := []struct {
		Name       string
		ObjectType sdk.ObjectType
		ObjectName string
		Expected   sdk.ObjectIdentifier
		Error      string
	}{
		{
			Name:       "database - account object identifier",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "test_database",
			Expected:   sdk.NewAccountObjectIdentifier("test_database"),
		},
		{
			Name:       "database - account object identifier - quoted",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "\"test_database\"",
			Expected:   sdk.NewAccountObjectIdentifier("test_database"),
		},
		{
			Name:       "schema - database object identifier",
			ObjectType: sdk.ObjectTypeSchema,
			ObjectName: "test_database.test_schema",
			Expected:   sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
		},
		{
			Name:       "schema - database object identifier - quoted",
			ObjectType: sdk.ObjectTypeSchema,
			ObjectName: "\"test_database\".\"test_schema\"",
			Expected:   sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
		},
		{
			Name:       "table - schema object identifier",
			ObjectType: sdk.ObjectTypeTable,
			ObjectName: "test_database.test_schema.test_table",
			Expected:   sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_table"),
		},
		{
			Name:       "table - schema object identifier - quoted",
			ObjectType: sdk.ObjectTypeTable,
			ObjectName: "\"test_database\".\"test_schema\".\"test_table\"",
			Expected:   sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_table"),
		},
		{
			Name:       "dbt project - schema object identifier",
			ObjectType: sdk.ObjectTypeDbtProject,
			ObjectName: "test_database.test_schema.test_dbt_project",
			Expected:   sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_dbt_project"),
		},
		{
			Name:       "semantic view - schema object identifier",
			ObjectType: sdk.ObjectTypeSemanticView,
			ObjectName: "test_database.test_schema.test_semantic_view",
			Expected:   sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_semantic_view"),
		},
		{
			Name:       "account - organization and account name",
			ObjectType: sdk.ObjectTypeAccount,
			ObjectName: "myorg.myaccount",
			Expected:   sdk.NewAccountIdentifier("myorg", "myaccount"),
		},
		{
			Name:       "validation - invalid account identifier (single part)",
			ObjectType: sdk.ObjectTypeAccount,
			ObjectName: "only_one_part",
			Error:      `expected 2 in a form of "<organization_name>.<account_name>"`,
		},
		{
			Name:       "account object identifier with dots",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "\"database.name.with.dots\"",
			Expected:   sdk.NewAccountObjectIdentifier("database.name.with.dots"),
		},
		{
			Name:       "validation - valid identifier",
			ObjectType: sdk.ObjectTypeDatabase,
			ObjectName: "to.many.parts.in.this.identifier",
			Expected:   sdk.NewAccountObjectIdentifier("to.many.parts.in.this.identifier"),
		},
		{
			Name:       "unknown type - one part is account object identifier",
			ObjectType: sdk.ObjectTypePostgresInstance,
			ObjectName: "pg_instance",
			Expected:   sdk.NewAccountObjectIdentifier("pg_instance"),
		},
		{
			Name:       "unknown type - two parts is database object identifier",
			ObjectType: sdk.ObjectType("SOME NEW TYPE"),
			ObjectName: "test_database.test_object",
			Expected:   sdk.NewDatabaseObjectIdentifier("test_database", "test_object"),
		},
		{
			Name:       "unknown type - three parts is schema object identifier",
			ObjectType: sdk.ObjectType("SOME NEW TYPE"),
			ObjectName: "test_database.test_schema.test_object",
			Expected:   sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_object"),
		},
		{
			Name:       "unknown type - four parts is table column identifier",
			ObjectType: sdk.ObjectType("SOME NEW TYPE"),
			ObjectName: "test_database.test_schema.test_table.column_name",
			Expected:   sdk.NewTableColumnIdentifier("test_database", "test_schema", "test_table", "column_name"),
		},
		{
			Name:       "unknown type - quoted one part with dots is account object identifier",
			ObjectType: sdk.ObjectTypePostgresInstance,
			ObjectName: `"name.with.dots"`,
			Expected:   sdk.NewAccountObjectIdentifier("name.with.dots"),
		},
		{
			Name:       "unknown type - five parts is rejected",
			ObjectType: sdk.ObjectType("SOME NEW TYPE"),
			ObjectName: "a.b.c.d.e",
			Error:      "unsupported identifier: a.b.c.d.e (number of parts: 5)",
		},
		{
			Name:       "validation - invalid database object identifier",
			ObjectType: sdk.ObjectTypeSchema,
			ObjectName: "test_database.test_schema.test_table",
			Error:      "unexpected number of parts 3 in identifier test_database.test_schema.test_table, expected 2 in a form of \"<database_name>.<database_object_name>\"",
		},
		{
			Name:       "validation - invalid schema object identifier",
			ObjectType: sdk.ObjectTypeTable,
			ObjectName: "test_database.test_schema.test_table.column_name",
			Error:      "unexpected number of parts 4 in identifier test_database.test_schema.test_table.column_name, expected 3 in a form of \"<database_name>.<schema_name>.<schema_object_name>\"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			id, err := GetOnObjectIdentifier(tt.ObjectType, tt.ObjectName)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}

func TestGetOwnershipGrantOn(t *testing.T) {
	testCases := []struct {
		Name     string
		On       map[string]any
		Expected sdk.OwnershipGrantOnRequest
		Error    string
	}{
		{
			Name: "database object type",
			On: map[string]any{
				"object_type": "DATABASE",
				"object_name": "test_database",
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeDatabase,
				Name:       sdk.NewAccountObjectIdentifier("test_database"),
			}),
		},
		{
			Name: "schema object type",
			On: map[string]any{
				"object_type": "SCHEMA",
				"object_name": "test_database.test_schema",
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeSchema,
				Name:       sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
			}),
		},
		{
			Name: "table object type",
			On: map[string]any{
				"object_type": "TABLE",
				"object_name": "test_database.test_schema.test_table",
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeTable,
				Name:       sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_table"),
			}),
		},
		{
			Name: "dbt project object type",
			On: map[string]any{
				"object_type": "DBT PROJECT",
				"object_name": "test_database.test_schema.test_dbt_project",
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeDbtProject,
				Name:       sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_dbt_project"),
			}),
		},
		{
			Name: "on all tables in database",
			On: map[string]any{
				"all": []any{
					map[string]any{
						"object_type_plural": "TABLES",
						"in_database":        "test_database",
					},
				},
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInDatabase(sdk.NewAccountObjectIdentifier("test_database"))),
		},
		{
			Name: "on all tables in schema",
			On: map[string]any{
				"all": []any{
					map[string]any{
						"object_type_plural": "TABLES",
						"in_schema":          "test_database.test_schema",
					},
				},
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInSchema(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"))),
		},
		{
			Name: "on future tables in database",
			On: map[string]any{
				"future": []any{
					map[string]any{
						"object_type_plural": "TABLES",
						"in_database":        "test_database",
					},
				},
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInDatabase(sdk.NewAccountObjectIdentifier("test_database"))),
		},
		{
			Name: "on future tables in schema",
			On: map[string]any{
				"future": []any{
					map[string]any{
						"object_type_plural": "TABLES",
						"in_schema":          "test_database.test_schema",
					},
				},
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInSchema(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"))),
		},
		{
			Name: "unknown account object type - postgres instance",
			On: map[string]any{
				"object_type": "POSTGRES INSTANCE",
				"object_name": "pg1",
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypePostgresInstance,
				Name:       sdk.NewAccountObjectIdentifier("pg1"),
			}),
		},
		{
			Name: "database object type in lowercase",
			On: map[string]any{
				"object_type": "database",
				"object_name": "test_database",
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeDatabase,
				Name:       sdk.NewAccountObjectIdentifier("test_database"),
			}),
		},
		{
			Name: "grant all in database plural object type in lowercase",
			On: map[string]any{
				"future": []any{
					map[string]any{
						"object_type_plural": "tables",
						"in_schema":          "test_database.test_schema",
					},
				},
			},
			Expected: *sdk.NewOwnershipGrantOnRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInSchema(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"))),
		},
		{
			Name: "validation - invalid schema object type",
			On: map[string]any{
				"object_type": "SCHEMA",
				"object_name": "test_database.test_schema.test_table",
			},
			Error: "unexpected number of parts 3 in identifier test_database.test_schema.test_table, expected 2 in a form of \"<database_name>.<database_object_name>\"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, grantOwnershipSchema, map[string]any{
				"on": []any{tt.On},
			})
			grantOn, err := getOwnershipGrantOn(d)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.NotNil(t, grantOn)
				assert.Equal(t, tt.Expected, *grantOn)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}

func TestPrepareShowGrantsRequestForGrantOwnership(t *testing.T) {
	testCases := []struct {
		Name              string
		Identifier        GrantOwnershipId
		Expected          *sdk.ShowGrantsRequest
		ExpectedGrantedOn []sdk.ObjectType
	}{
		{
			Name: "show for object - database",
			Identifier: GrantOwnershipId{
				Kind: OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeDatabase,
					ObjectName: sdk.NewAccountObjectIdentifier("test_database"),
				},
			},
			Expected: sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeDatabase,
				Name:       sdk.NewAccountObjectIdentifier("test_database"),
			})),
			ExpectedGrantedOn: []sdk.ObjectType{sdk.ObjectTypeDatabase},
		},
		{
			Name: "show for object - schema",
			Identifier: GrantOwnershipId{
				Kind: OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeSchema,
					ObjectName: sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
				},
			},
			Expected: sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeSchema,
				Name:       sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
			})),
			ExpectedGrantedOn: []sdk.ObjectType{sdk.ObjectTypeSchema},
		},
		{
			Name: "show for object - database role",
			Identifier: GrantOwnershipId{
				Kind: OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeDatabaseRole,
					ObjectName: sdk.NewDatabaseObjectIdentifier("test_database", "test_database_role"),
				},
			},
			Expected: sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
				ObjectType: sdk.ObjectTypeDatabaseRole,
				Name:       sdk.NewDatabaseObjectIdentifier("test_database", "test_database_role"),
			})),
			// Before BCR-2371 (2026_06 bundle) SHOW GRANTS ON DATABASE ROLE reports granted_on = ROLE, after it reports DATABASE_ROLE; both must match.
			ExpectedGrantedOn: []sdk.ObjectType{sdk.ObjectTypeRole, sdk.ObjectTypeDatabaseRole},
		},
		{
			Name: "show for all in database",
			Identifier: GrantOwnershipId{
				Kind: OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("test_database")),
				},
			},
			Expected:          nil,
			ExpectedGrantedOn: nil,
		},
		{
			Name: "show for all in schema",
			Identifier: GrantOwnershipId{
				Kind: OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
			Expected:          nil,
			ExpectedGrantedOn: nil,
		},
		{
			Name: "show for future in database",
			Identifier: GrantOwnershipId{
				Kind: OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("test_database")),
				},
			},
			Expected:          sdk.NewShowGrantsRequest().WithFuture(true).WithIn(*sdk.NewShowGrantsInRequest().WithDatabase(sdk.NewAccountObjectIdentifier("test_database"))),
			ExpectedGrantedOn: []sdk.ObjectType{sdk.ObjectTypeTable},
		},
		{
			Name: "show for future in schema",
			Identifier: GrantOwnershipId{
				Kind: OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
			Expected:          sdk.NewShowGrantsRequest().WithFuture(true).WithIn(*sdk.NewShowGrantsInRequest().WithSchema(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"))),
			ExpectedGrantedOn: []sdk.ObjectType{sdk.ObjectTypeTable},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			opts, grantedOn := prepareShowGrantsRequestForGrantOwnership(&tt.Identifier)
			if tt.Expected == nil {
				assert.Nil(t, opts)
			} else {
				assert.NotNil(t, opts)
				assert.Equal(t, *tt.Expected, *opts)
			}
			assert.Equal(t, tt.ExpectedGrantedOn, grantedOn)
		})
	}
}

func TestValidAccountRoleNameGetOwnershipGrantTo(t *testing.T) {
	testCases := []struct {
		Name        string
		AccountRole *string
		Expected    sdk.OwnershipGrantToRequest
		Error       string
	}{
		{
			Name:        "account role name",
			AccountRole: sdk.String("account_role_name"),
			Expected:    *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(sdk.NewAccountObjectIdentifier("account_role_name")),
		},
		{
			Name:        "account role name - quoted",
			AccountRole: sdk.String("\"account_role_name\""),
			Expected:    *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(sdk.NewAccountObjectIdentifier("account_role_name")),
		},
		{
			Name:        "account role name - with dots",
			AccountRole: sdk.String("account.role.with.dots"),
			Expected:    *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(sdk.NewAccountObjectIdentifier("account.role.with.dots")),
		},
		{
			Name:        "account role name - with dots quoted",
			AccountRole: sdk.String("\"account.role.with.dots\""),
			Expected:    *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(sdk.NewAccountObjectIdentifier("account.role.with.dots")),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			grantTo, err := getOwnershipGrantTo(schema.TestResourceDataRaw(t, grantOwnershipSchema, map[string]any{
				"account_role_name": *tt.AccountRole,
			}))
			if tt.Error != "" {
				require.ErrorContains(t, err, tt.Error)
			} else {
				require.NoError(t, err)
				assert.Equal(t, *tt.Expected.AccountRoleName, *grantTo.AccountRoleName)
			}
		})
	}
}

func TestValidDatabaseRoleNameGetOwnershipGrantTo(t *testing.T) {
	testCases := []struct {
		Name         string
		DatabaseRole *string
		Expected     sdk.OwnershipGrantToRequest
	}{
		{
			Name:         "database role name",
			DatabaseRole: sdk.String("test_database.database_role_name"),
			Expected:     *sdk.NewOwnershipGrantToRequest().WithDatabaseRoleName(sdk.NewDatabaseObjectIdentifier("test_database", "database_role_name")),
		},
		{
			Name:         "database role name - quoted",
			DatabaseRole: sdk.String("\"test_database\".\"database_role_name\""),
			Expected:     *sdk.NewOwnershipGrantToRequest().WithDatabaseRoleName(sdk.NewDatabaseObjectIdentifier("test_database", "database_role_name")),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			grantTo, err := getOwnershipGrantTo(schema.TestResourceDataRaw(t, grantOwnershipSchema, map[string]any{
				"database_role_name": *tt.DatabaseRole,
			}))
			require.NoError(t, err)

			assert.Equal(t, *tt.Expected.DatabaseRoleName, *grantTo.DatabaseRoleName)
		})
	}
}

func TestInvalidDatabaseRoleGetOwnershipGrantTo(t *testing.T) {
	d := schema.TestResourceDataRaw(t, grantOwnershipSchema, map[string]any{
		"database_role_name": "account_role_name",
	})

	_, err := getOwnershipGrantTo(d)
	require.Error(t, err)
}

func TestGetOwnershipGrantOpts(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier GrantOwnershipId
		Expected   *sdk.OwnershipCurrentGrantsRequest
	}{
		{
			Name: "outbound privileges copy",
			Identifier: GrantOwnershipId{
				OutboundPrivilegesBehavior: sdk.Pointer(CopyOutboundPrivilegesBehavior),
			},
			Expected: sdk.NewOwnershipCurrentGrantsRequest(sdk.Copy),
		},
		{
			Name: "outbound privileges revoke",
			Identifier: GrantOwnershipId{
				OutboundPrivilegesBehavior: sdk.Pointer(RevokeOutboundPrivilegesBehavior),
			},
			Expected: sdk.NewOwnershipCurrentGrantsRequest(sdk.Revoke),
		},
		{
			Name: "no outbound privileges option",
			Identifier: GrantOwnershipId{
				OutboundPrivilegesBehavior: nil,
			},
			Expected: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			currentGrants := getOwnershipCurrentGrants(&tt.Identifier)
			assert.Equal(t, tt.Expected, currentGrants)
		})
	}
}
