package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestParseGrantOwnershipId(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier string
		Expected   GrantOwnershipId
		Error      string
	}{
		{
			Name:       "grant ownership on database to account role",
			Identifier: `ToAccountRole|"account-role"|COPY|OnObject|DATABASE|"database-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeDatabase,
					ObjectName: sdk.NewAccountObjectIdentifier("database-name"),
				},
			},
		},
		{
			Name:       "grant ownership on schema to account role",
			Identifier: `ToAccountRole|"account-role"|COPY|OnObject|SCHEMA|"database-name"."schema-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeSchema,
					ObjectName: sdk.NewDatabaseObjectIdentifier("database-name", "schema-name"),
				},
			},
		},
		{
			Name:       "grant ownership on table to account role",
			Identifier: `ToAccountRole|"account-role"|COPY|OnObject|TABLE|"database-name"."schema-name"."table-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeTable,
					ObjectName: sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
				},
			},
		},
		{
			Name:       "grant ownership on schema to database role",
			Identifier: `ToDatabaseRole|"database-name"."database-role"|REVOKE|OnObject|SCHEMA|"database-name"."schema-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToDatabaseGrantOwnershipTargetRoleKind,
				DatabaseRoleName:             sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(RevokeOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeSchema,
					ObjectName: sdk.NewDatabaseObjectIdentifier("database-name", "schema-name"),
				},
			},
		},
		{
			Name:       "grant ownership on all tables in database to account role",
			Identifier: `ToAccountRole|"account-role"||OnAll|TABLES|InDatabase|"database-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				Kind:                         OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
				},
			},
		},
		{
			Name:       "grant ownership on all tables in database to account role - empty database id",
			Identifier: `ToAccountRole|"account-role"||OnAll|TABLES|InDatabase|`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				Kind:                         OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("")),
				},
			},
		},
		{
			Name:       "grant ownership on all tables in schema to account role",
			Identifier: `ToAccountRole|"account-role"||OnAll|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				Kind:                         OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
				},
			},
		},
		{
			Name:       "grant ownership on future tables in database to account role",
			Identifier: `ToAccountRole|"account-role"|COPY|OnFuture|TABLES|InDatabase|"database-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database-name")),
				},
			},
		},
		{
			Name:       "grant ownership on future tables in database to account role - empty database",
			Identifier: `ToAccountRole|"account-role"|COPY|OnFuture|TABLES|InDatabase|`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("")),
				},
			},
		},
		{
			Name:       "grant ownership on future tables in schema to account role",
			Identifier: `ToAccountRole|"account-role"|COPY|OnFuture|TABLES|InSchema|"database-name"."schema-name"`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database-name", "schema-name")),
				},
			},
		},
		{
			Name:       "grant ownership on function to account role",
			Identifier: `ToAccountRole|"account-role"|COPY|OnObject|FUNCTION|"database-name"."schema-name"."function-name"(FLOAT)`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeFunction,
					ObjectName: sdk.NewSchemaObjectIdentifierWithArguments("database-name", "schema-name", "function-name", sdk.DataTypeFloat),
				},
			},
		},
		{
			Name:       "grant ownership on function without arguments to database role",
			Identifier: `ToDatabaseRole|"database-name"."database-role"|REVOKE|OnObject|FUNCTION|"database-name"."schema-name"."function-name"()`,
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToDatabaseGrantOwnershipTargetRoleKind,
				DatabaseRoleName:             sdk.NewDatabaseObjectIdentifier("database-name", "database-role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(RevokeOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeFunction,
					ObjectName: sdk.NewSchemaObjectIdentifierWithArguments("database-name", "schema-name", "function-name", []sdk.DataType{}...),
				},
			},
		},
		{
			Name:       "validation: not enough parts",
			Identifier: `ToDatabaseRole|"database-name"."role-name"|`,
			Error:      "ownership identifier should hold at least 5 parts",
		},
		{
			Name:       "validation: invalid to role enum",
			Identifier: `SomeInvalidEnum|"database-name"."role-name"|OnObject|DATABASE|"some-database"`,
			Error:      "unknown GrantOwnershipTargetRoleKind: SomeInvalidEnum, valid options are ToAccountRole | ToDatabaseRole",
		},
		{
			Name:       "invalid outbound privilege option resulting in no outbound privileges option set",
			Identifier: `ToAccountRole|"account-role"|InvalidOption|OnFuture|TABLES|InSchema|"database-name"."schema-name"`,
			Error:      `unknown OutboundPrivilegesBehavior: InvalidOption, valid options are COPY | REVOKE`,
		},
		{
			Name:       "validation: not enough parts for OnObject kind",
			Identifier: `ToAccountRole|"account-role"|COPY|OnObject|DATABASE`,
			Error:      `grant ownership identifier should consist of 6 parts`,
		},
		{
			Name:       "validation: not enough parts for OnAll kind",
			Identifier: `ToAccountRole|"account-role"|COPY|OnAll|TABLES|InDatabase`,
			Error:      `grant ownership identifier should consist of 7 parts`,
		},
		{
			Name:       "validation: OnAll in InvalidOption",
			Identifier: `ToAccountRole|"account-role"|COPY|OnAll|TABLES|InvalidOption|"some-identifier"`,
			Error:      "invalid BulkOperationGrantKind: InvalidOption, valid options are InDatabase | InSchema",
		},
		{
			Name:       "validation: OnAll in database - missing schema identifier",
			Identifier: `ToAccountRole|"account-role"|COPY|OnAll|InvalidTarget|InSchema|`,
			Error:      "incompatible identifier: ",
		},
		{
			Name:       "validation: not enough parts for OnFuture kind",
			Identifier: `ToAccountRole|"account-role"|COPY|OnFuture|TABLES`,
			Error:      `grant ownership identifier should consist of 7 parts`,
		},
		{
			Name:       "validation: OnFuture in InvalidOption",
			Identifier: `ToAccountRole|"account-role"|COPY|OnFuture|TABLES|InvalidOption|"some-identifier"`,
			Error:      "invalid BulkOperationGrantKind: InvalidOption, valid options are InDatabase | InSchema",
		},
		{
			Name:       "validation: OnFuture in database - missing schema identifier",
			Identifier: `ToAccountRole|"account-role"|COPY|OnFuture|InvalidTarget|InSchema|`,
			Error:      "incompatible identifier: ",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			id, err := ParseGrantOwnershipId(tt.Identifier)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.NotNil(t, id)
				assert.Equal(t, tt.Expected, *id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}

func TestGrantOwnershipIdString(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier GrantOwnershipId
		Expected   string
		Error      string
	}{
		{
			Name: "grant ownership on database to account role",
			Identifier: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account_role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeDatabase,
					ObjectName: sdk.NewAccountObjectIdentifier("database_name"),
				},
			},
			Expected: `ToAccountRole|"account_role"|COPY|OnObject|DATABASE|"database_name"`,
		},
		{
			Name: "grant ownership on schema to account role",
			Identifier: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account_role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(RevokeOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeSchema,
					ObjectName: sdk.NewDatabaseObjectIdentifier("database_name", "schema_name"),
				},
			},
			Expected: `ToAccountRole|"account_role"|REVOKE|OnObject|SCHEMA|"database_name"."schema_name"`,
		},
		{
			Name: "grant ownership on table to account role",
			Identifier: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("account_role"),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeTable,
					ObjectName: sdk.NewSchemaObjectIdentifier("database_name", "schema_name", "table_name"),
				},
			},
			Expected: `ToAccountRole|"account_role"||OnObject|TABLE|"database_name"."schema_name"."table_name"`,
		},
		{
			Name: "grant ownership on all tables in database to database role",
			Identifier: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToDatabaseGrantOwnershipTargetRoleKind,
				DatabaseRoleName:             sdk.NewDatabaseObjectIdentifier("database_name", "database_role_name"),
				Kind:                         OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database_name")),
				},
			},
			Expected: `ToDatabaseRole|"database_name"."database_role_name"||OnAll|TABLES|InDatabase|"database_name"`,
		},
		{
			Name: "grant ownership on all tables in schema to database role",
			Identifier: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToDatabaseGrantOwnershipTargetRoleKind,
				DatabaseRoleName:             sdk.NewDatabaseObjectIdentifier("database_name", "database_role_name"),
				Kind:                         OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database_name", "schema_name")),
				},
			},
			Expected: `ToDatabaseRole|"database_name"."database_role_name"||OnAll|TABLES|InSchema|"database_name"."schema_name"`,
		},
		{
			Name: "grant ownership on future tables in database to database role",
			Identifier: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToDatabaseGrantOwnershipTargetRoleKind,
				DatabaseRoleName:             sdk.NewDatabaseObjectIdentifier("database_name", "database_role_name"),
				Kind:                         OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("database_name")),
				},
			},
			Expected: `ToDatabaseRole|"database_name"."database_role_name"||OnFuture|TABLES|InDatabase|"database_name"`,
		},
		{
			Name: "grant ownership on future tables in schema to database role",
			Identifier: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToDatabaseGrantOwnershipTargetRoleKind,
				DatabaseRoleName:             sdk.NewDatabaseObjectIdentifier("database_name", "database_role_name"),
				Kind:                         OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database_name", "schema_name")),
				},
			},
			Expected: `ToDatabaseRole|"database_name"."database_role_name"||OnFuture|TABLES|InSchema|"database_name"."schema_name"`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Identifier.String())
		})
	}
}

func TestCreateGrantOwnershipIdFromSchema(t *testing.T) {
	testCases := []struct {
		Name     string
		Config   map[string]any
		Expected GrantOwnershipId
	}{
		{
			Name: "grant ownership on schema to account role with copied outbound privileges",
			Config: map[string]any{
				"account_role_name":   "test_acc_role_name",
				"outbound_privileges": "COPY",
				"on": []any{
					map[string]any{
						"object_type": "SCHEMA",
						"object_name": "\"test_database\".\"test_schema\"",
					},
				},
			},
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("test_acc_role_name"),
				OutboundPrivilegesBehavior:   sdk.Pointer(CopyOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeSchema,
					ObjectName: sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
				},
			},
		},
		{
			Name: "grant ownership on schema to database role with revoked outbound privileges",
			Config: map[string]any{
				"database_role_name":  "\"test_database\".\"test_database_role\"",
				"outbound_privileges": "REVOKE",
				"on": []any{
					map[string]any{
						"object_type": "SCHEMA",
						"object_name": "\"test_database\".\"test_schema\"",
					},
				},
			},
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToDatabaseGrantOwnershipTargetRoleKind,
				DatabaseRoleName:             sdk.NewDatabaseObjectIdentifier("test_database", "test_database_role"),
				OutboundPrivilegesBehavior:   sdk.Pointer(RevokeOutboundPrivilegesBehavior),
				Kind:                         OnObjectGrantOwnershipKind,
				Data: &OnObjectGrantOwnershipData{
					ObjectType: sdk.ObjectTypeSchema,
					ObjectName: sdk.NewDatabaseObjectIdentifier("test_database", "test_schema"),
				},
			},
		},
		{
			Name: "grant ownership on all tables in database to account role",
			Config: map[string]any{
				"account_role_name": "test_acc_role",
				"on": []any{
					map[string]any{
						"all": []any{
							map[string]any{
								"object_type_plural": "tables",
								"in_database":        "test_database",
							},
						},
					},
				},
			},
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("test_acc_role"),
				Kind:                         OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("test_database")),
				},
			},
		},
		{
			Name: "grant ownership on all tables in schema to account role",
			Config: map[string]any{
				"account_role_name": "test_acc_role",
				"on": []any{
					map[string]any{
						"all": []any{
							map[string]any{
								"object_type_plural": "tables",
								"in_schema":          "\"test_database\".\"test_schema\"",
							},
						},
					},
				},
			},
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("test_acc_role"),
				Kind:                         OnAllGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
		},
		{
			Name: "grant ownership on future tables in database to account role",
			Config: map[string]any{
				"account_role_name": "test_acc_role",
				"on": []any{
					map[string]any{
						"future": []any{
							map[string]any{
								"object_type_plural": "tables",
								"in_database":        "test_database",
							},
						},
					},
				},
			},
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("test_acc_role"),
				Kind:                         OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InDatabaseBulkOperationGrantKind,
					Database:         sdk.Pointer(sdk.NewAccountObjectIdentifier("test_database")),
				},
			},
		},
		{
			Name: "grant ownership on future tables in schema to account role",
			Config: map[string]any{
				"account_role_name": "test_acc_role",
				"on": []any{
					map[string]any{
						"future": []any{
							map[string]any{
								"object_type_plural": "tables",
								"in_schema":          "\"test_database\".\"test_schema\"",
							},
						},
					},
				},
			},
			Expected: GrantOwnershipId{
				GrantOwnershipTargetRoleKind: ToAccountGrantOwnershipTargetRoleKind,
				AccountRoleName:              sdk.NewAccountObjectIdentifier("test_acc_role"),
				Kind:                         OnFutureGrantOwnershipKind,
				Data: &BulkOperationGrantData{
					ObjectNamePlural: sdk.PluralObjectTypeTables,
					Kind:             InSchemaBulkOperationGrantKind,
					Schema:           sdk.Pointer(sdk.NewDatabaseObjectIdentifier("test_database", "test_schema")),
				},
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, grantOwnershipSchema, tt.Config)
			id, err := createGrantOwnershipIdFromSchema(d)
			assert.NoError(t, err)
			assert.NotNil(t, id)
			assert.Equal(t, tt.Expected.GrantOwnershipTargetRoleKind, id.GrantOwnershipTargetRoleKind)
			assert.Equal(t, tt.Expected.AccountRoleName, id.AccountRoleName)
			assert.Equal(t, tt.Expected.DatabaseRoleName, id.DatabaseRoleName)
			assert.Equal(t, tt.Expected.OutboundPrivilegesBehavior, id.OutboundPrivilegesBehavior)
			assert.Equal(t, tt.Expected.Kind, id.Kind)
			assert.Equal(t, tt.Expected.Data.String(), id.Data.String())
		})
	}
}
