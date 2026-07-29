package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestParseGrantPrivilegesToShareId(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier string
		Expected   GrantPrivilegesToShareId
		Error      string
	}{
		{
			Name:       "grant privileges on database to share",
			Identifier: `"share-name"|REFERENCE_USAGE|OnDatabase|"on-database-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"REFERENCE_USAGE"},
				Kind:       OnDatabaseShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("on-database-name"),
			},
		},
		{
			Name:       "grant privileges on database to share - database identifier with dots",
			Identifier: `"share-name"|SELECT|OnDatabase|"one.two.three.four.five.six.seven.eight.nine.ten"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"SELECT"},
				Kind:       OnDatabaseShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("one.two.three.four.five.six.seven.eight.nine.ten"),
			},
		},
		{
			Name:       "grant privileges on schema to share",
			Identifier: `"share-name"|USAGE|OnSchema|"on-database-name"."on-schema-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"USAGE"},
				Kind:       OnSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("on-database-name", "on-schema-name"),
			},
		},
		{
			Name:       "grant privileges on function to share",
			Identifier: `"share-name"|USAGE|OnFunction|"on-database-name"."on-schema-name".on-function-name(INT, VARCHAR)`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"USAGE"},
				Kind:       OnFunctionShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifierWithArguments("on-database-name", "on-schema-name", "on-function-name", sdk.DataTypeInt, sdk.DataTypeVARCHAR),
			},
		},
		{
			Name:       "grant privileges on function without arguments to share",
			Identifier: `"share-name"|READ|OnFunction|"on-database-name"."on-schema-name"."on-view-name"()`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"READ"},
				Kind:       OnFunctionShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifierWithArguments("on-database-name", "on-schema-name", "on-view-name", []sdk.DataType{}...),
			},
		},
		{
			Name:       "grant privileges on table to share",
			Identifier: `"share-name"|EVOLVE SCHEMA|OnTable|"on-database-name"."on-schema-name"."on-table-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"EVOLVE SCHEMA"},
				Kind:       OnTableShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("on-database-name", "on-schema-name", "on-table-name"),
			},
		},
		{
			Name:       "grant privileges on all tables in schema to share",
			Identifier: `"share-name"|EVOLVE SCHEMA,SELECT|OnAllTablesInSchema|"on-database-name"."on-schema-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"EVOLVE SCHEMA", "SELECT"},
				Kind:       OnAllTablesInSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("on-database-name", "on-schema-name"),
			},
		},
		{
			Name:       "grant privileges on tag to share",
			Identifier: `"share-name"|READ|OnTag|"database-name"."schema-name"."on-tag-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"READ"},
				Kind:       OnTagShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "on-tag-name"),
			},
		},
		{
			Name:       "grant privileges on view to share",
			Identifier: `"share-name"|READ|OnView|"on-database-name"."on-schema-name"."on-view-name"`,
			Expected: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"READ"},
				Kind:       OnViewShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("on-database-name", "on-schema-name", "on-view-name"),
			},
		},
		{
			Name:       "validation: not enough parts",
			Identifier: `"share-name"|SELECT|OnDatabase`,
			Error:      `snowflake_grant_privileges_to_share id is composed out of 4 parts "<share_name>|<privileges>|<grant_on_type>|<grant_on_identifier>", but got 3 parts: ["share-name" SELECT OnDatabase]`,
		},
		{
			Name:       "validation: empty privileges",
			Identifier: `"share-name"||OnDatabase|"database-name"`,
			Error:      `invalid Privileges value: [], should be comma separated list of privileges`,
		},
		{
			Name:       "validation: unsupported kind",
			Identifier: `"share-name"|SELECT|OnSomething|"object-name"`,
			Error:      `unexpected share grant kind: OnSomething`,
		},
		{
			Name:       "validation: invalid account object identifier",
			Identifier: `"share-name"|SELECT|OnTable|"one"."two"`,
			Error:      `unexpected number of parts 2 in identifier "one"."two", expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`,
		},
		{
			Name:       "validation: invalid database object identifier",
			Identifier: `"share-name"|SELECT|OnSchema|one.two.three`,
			Error:      `unexpected number of parts 3 in identifier one.two.three, expected 2 in a form of "<database_name>.<database_object_name>`,
		},
		{
			Name:       "validation: invalid schema object identifier",
			Identifier: `"share-name"|SELECT|OnTable|one`,
			Error:      `unexpected number of parts 1 in identifier one, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			id, err := ParseGrantPrivilegesToShareId(tt.Identifier)
			if tt.Error == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.Expected, id)
			} else {
				assert.ErrorContains(t, err, tt.Error)
			}
		})
	}
}

func TestGrantPrivilegesToShareIdString(t *testing.T) {
	testCases := []struct {
		Name       string
		Identifier GrantPrivilegesToShareId
		Expected   string
		Error      string
	}{
		{
			Name: "grant privileges on database to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"REFERENCE_USAGE"},
				Kind:       OnDatabaseShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("database-name"),
			},
			Expected: `"share-name"|REFERENCE_USAGE|OnDatabase|"database-name"`,
		},
		{
			Name: "grant privileges on schema to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"USAGE"},
				Kind:       OnSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("database-name", "schema-name"),
			},
			Expected: `"share-name"|USAGE|OnSchema|"database-name"."schema-name"`,
		},
		{
			Name: "grant privileges on table to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"EVOLVE SCHEMA", "SELECT"},
				Kind:       OnTableShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "table-name"),
			},
			Expected: `"share-name"|EVOLVE SCHEMA,SELECT|OnTable|"database-name"."schema-name"."table-name"`,
		},
		{
			Name: "grant privileges on all tables in schema to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"EVOLVE SCHEMA", "SELECT"},
				Kind:       OnAllTablesInSchemaShareGrantKind,
				Identifier: sdk.NewDatabaseObjectIdentifier("database-name", "schema-name"),
			},
			Expected: `"share-name"|EVOLVE SCHEMA,SELECT|OnAllTablesInSchema|"database-name"."schema-name"`,
		},
		{
			Name: "grant privileges on tag to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"READ"},
				Kind:       OnTagShareGrantKind,
				Identifier: sdk.NewAccountObjectIdentifier("tag-name"),
			},
			Expected: `"share-name"|READ|OnTag|"tag-name"`,
		},
		{
			Name: "grant privileges on view to share",
			Identifier: GrantPrivilegesToShareId{
				ShareName:  sdk.NewAccountObjectIdentifier("share-name"),
				Privileges: []string{"SELECT"},
				Kind:       OnViewShareGrantKind,
				Identifier: sdk.NewSchemaObjectIdentifier("database-name", "schema-name", "view-name"),
			},
			Expected: `"share-name"|SELECT|OnView|"database-name"."schema-name"."view-name"`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Identifier.String())
		})
	}
}

func TestIsGrantedOnEquivalent(t *testing.T) {
	testCases := []struct {
		Name     string
		Expected sdk.ObjectType
		Actual   sdk.ObjectType
		Result   bool
	}{
		{
			Name:     "identical types",
			Expected: sdk.ObjectTypeTable,
			Actual:   sdk.ObjectTypeTable,
			Result:   true,
		},
		{
			Name:     "table vs external table",
			Expected: sdk.ObjectTypeTable,
			Actual:   sdk.ObjectTypeExternalTable,
			Result:   true,
		},
		{
			Name:     "table vs dynamic table",
			Expected: sdk.ObjectTypeTable,
			Actual:   sdk.ObjectTypeDynamicTable,
			Result:   true,
		},
		{
			Name:     "table vs iceberg table",
			Expected: sdk.ObjectTypeTable,
			Actual:   sdk.ObjectTypeIcebergTable,
			Result:   true,
		},
		{
			Name:     "external table vs table - not equivalent (expected is not TABLE)",
			Expected: sdk.ObjectTypeExternalTable,
			Actual:   sdk.ObjectTypeTable,
			Result:   false,
		},
		{
			Name:     "different base types",
			Expected: sdk.ObjectTypeTable,
			Actual:   sdk.ObjectTypeView,
			Result:   false,
		},
		{
			Name:     "view vs table",
			Expected: sdk.ObjectTypeView,
			Actual:   sdk.ObjectTypeTable,
			Result:   false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			result := isGrantedOnEquivalent(tt.Expected, tt.Actual)
			assert.Equal(t, tt.Result, result)
		})
	}
}
