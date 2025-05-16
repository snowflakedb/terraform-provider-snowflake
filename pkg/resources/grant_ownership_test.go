package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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
			Name:       "validation - unsupported type",
			ObjectType: sdk.ObjectTypeShare,
			ObjectName: "some_share",
			Error:      "object_type SHARE is not supported",
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
		tt := tt
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
		Expected sdk.OwnershipGrantOn
		Error    string
	}{
		{
			Name: "database object type",
			On: map[string]any{
				"object_type": "DATABASE",
				"object_name": "test_database",
			},
			Expected: sdk.OwnershipGrantOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeDatabase,
					Name:       sdk.NewAccountObjectIdentifier("test_database"),
				},
			},
		},
		{
			Name: "schema object type",
			On: map[string]any{
				"object_type": "SCHEMA",
				"object_name": "test_database.test_schema",
			},
			Expected: sdk.OwnershipGrantOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeSchema,
					Name:       sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
				},
			},
		},
		{
			Name: "table object type",
			On: map[string]any{
				"object_type": "TABLE",
				"object_name": "test_database.test_schema.test_table",
			},
			Expected: sdk.OwnershipGrantOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeTable,
					Name:       sdk.NewSchemaObjectIdentifier("test_database", "test_schema", "test_table"),
				},
			},
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
			Expected: sdk.OwnershipGrantOn{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InDatabase:       sdk.Pointer(sdk.NewAccountObjectIdentifier("test_database")),
				},
			},
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
			Expected: sdk.OwnershipGrantOn{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InSchema:         sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
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
			Expected: sdk.OwnershipGrantOn{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InDatabase:       sdk.Pointer(sdk.NewAccountObjectIdentifier("test_database")),
				},
			},
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
			Expected: sdk.OwnershipGrantOn{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InSchema:         sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
		},
		{
			Name: "database object type in lowercase",
			On: map[string]any{
				"object_type": "database",
				"object_name": "test_database",
			},
			Expected: sdk.OwnershipGrantOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeDatabase,
					Name:       sdk.NewAccountObjectIdentifier("test_database"),
				},
			},
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
			Expected: sdk.OwnershipGrantOn{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeTables,
					InSchema:         sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
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
		tt := tt
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
		Name                   string
		Identifier             GrantOwnershipId
		ExpectedShowGrantsOpts *sdk.ShowGrantOptions
		ExpectedGrantedOn      sdk.ObjectType
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
			ExpectedShowGrantsOpts: &sdk.ShowGrantOptions{
				On: &sdk.ShowGrantsOn{
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeDatabase,
						Name:       sdk.NewAccountObjectIdentifier("test_database"),
					},
				},
			},
			ExpectedGrantedOn: sdk.ObjectTypeDatabase,
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
			ExpectedShowGrantsOpts: &sdk.ShowGrantOptions{
				On: &sdk.ShowGrantsOn{
					Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeSchema,
						Name:       sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
					},
				},
			},
			ExpectedGrantedOn: sdk.ObjectTypeSchema,
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
			ExpectedShowGrantsOpts: nil,
			ExpectedGrantedOn:      "",
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
			ExpectedShowGrantsOpts: nil,
			ExpectedGrantedOn:      "",
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
			ExpectedShowGrantsOpts: &sdk.ShowGrantOptions{
				Future: sdk.Bool(true),
				In: &sdk.ShowGrantsIn{
					Database: sdk.Pointer(sdk.NewAccountObjectIdentifier("test_database")),
				},
			},
			ExpectedGrantedOn: sdk.ObjectTypeTable,
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
			ExpectedShowGrantsOpts: &sdk.ShowGrantOptions{
				Future: sdk.Bool(true),
				In: &sdk.ShowGrantsIn{
					Schema: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
			ExpectedGrantedOn: sdk.ObjectTypeTable,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			opts, grantedOn := prepareShowGrantsRequestForGrantOwnership(&tt.Identifier)
			if tt.ExpectedShowGrantsOpts == nil {
				assert.Nil(t, opts)
			} else {
				assert.NotNil(t, opts)
				assert.Equal(t, *tt.ExpectedShowGrantsOpts, *opts)
			}
			assert.Equal(t, tt.ExpectedGrantedOn, grantedOn)
		})
	}
}

func TestValidAccountRoleNameGetOwnershipGrantTo(t *testing.T) {
	testCases := []struct {
		Name        string
		AccountRole *string
		Expected    sdk.OwnershipGrantTo
		Error       string
	}{
		{
			Name:        "account role name",
			AccountRole: sdk.String("account_role_name"),
			Expected: sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(sdk.NewAccountObjectIdentifier("account_role_name")),
			},
		},
		{
			Name:        "account role name - quoted",
			AccountRole: sdk.String("\"account_role_name\""),
			Expected: sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(sdk.NewAccountObjectIdentifier("account_role_name")),
			},
		},
		{
			Name:        "account role name - with dots",
			AccountRole: sdk.String("account.role.with.dots"),
			Expected: sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(sdk.NewAccountObjectIdentifier("account.role.with.dots")),
			},
		},
		{
			Name:        "account role name - with dots quoted",
			AccountRole: sdk.String("\"account.role.with.dots\""),
			Expected: sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(sdk.NewAccountObjectIdentifier("account.role.with.dots")),
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
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
		Expected     sdk.OwnershipGrantTo
	}{
		{
			Name:         "database role name",
			DatabaseRole: sdk.String("test_database.database_role_name"),
			Expected: sdk.OwnershipGrantTo{
				DatabaseRoleName: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "database_role_name")),
			},
		},
		{
			Name:         "database role name - quoted",
			DatabaseRole: sdk.String("\"test_database\".\"database_role_name\""),
			Expected: sdk.OwnershipGrantTo{
				DatabaseRoleName: sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "database_role_name")),
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
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
		Expected   *sdk.GrantOwnershipOptions
	}{
		{
			Name: "outbound privileges copy",
			Identifier: GrantOwnershipId{
				OutboundPrivilegesBehavior: sdk.Pointer(CopyOutboundPrivilegesBehavior),
			},
			Expected: &sdk.GrantOwnershipOptions{
				CurrentGrants: &sdk.OwnershipCurrentGrants{
					OutboundPrivileges: sdk.Copy,
				},
			},
		},
		{
			Name: "outbound privileges revoke",
			Identifier: GrantOwnershipId{
				OutboundPrivilegesBehavior: sdk.Pointer(RevokeOutboundPrivilegesBehavior),
			},
			Expected: &sdk.GrantOwnershipOptions{
				CurrentGrants: &sdk.OwnershipCurrentGrants{
					OutboundPrivileges: sdk.Revoke,
				},
			},
		},
		{
			Name: "no outbound privileges option",
			Identifier: GrantOwnershipId{
				OutboundPrivilegesBehavior: nil,
			},
			Expected: &sdk.GrantOwnershipOptions{},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			opts := getOwnershipGrantOpts(&tt.Identifier)
			assert.NotNil(t, opts)
			assert.Equal(t, *tt.Expected, *opts)
		})
	}
}
