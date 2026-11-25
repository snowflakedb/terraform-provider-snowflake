package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO: add tests for non object identifiers
func TestEncodeSnowflakeIDLegacy(t *testing.T) {
	testCases := map[string]struct {
		identifier        sdk.ObjectIdentifier
		expectedEncodedID string
	}{
		"encodes account object identifier": {
			identifier:        sdk.NewAccountObjectIdentifier("database"),
			expectedEncodedID: `database`,
		},
		"encodes quoted account object identifier": {
			identifier:        sdk.NewAccountObjectIdentifier("\"database\""),
			expectedEncodedID: `database`,
		},
		"encodes account object identifier with a dot": {
			identifier:        sdk.NewAccountObjectIdentifier("data.base"),
			expectedEncodedID: `data.base`,
		},
		"encodes pointer to account object identifier": {
			identifier:        new(sdk.NewAccountObjectIdentifier("database")),
			expectedEncodedID: `database`,
		},
		"encodes database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("database", "schema"),
			expectedEncodedID: `database|schema`,
		},
		"encodes quoted database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("\"database\"", "\"schema\""),
			expectedEncodedID: `database|schema`,
		},
		"encodes database object identifier with dots": {
			identifier:        sdk.NewDatabaseObjectIdentifier("data.base", "sche.ma"),
			expectedEncodedID: `data.base|sche.ma`,
		},
		"encodes pointer to database object identifier": {
			identifier:        new(sdk.NewDatabaseObjectIdentifier("database", "schema")),
			expectedEncodedID: `database|schema`,
		},
		"encodes schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("database", "schema", "table"),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes quoted schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("\"database\"", "\"schema\"", "\"table\""),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes schema object identifier with dots": {
			identifier:        sdk.NewSchemaObjectIdentifier("data.base", "sche.ma", "tab.le"),
			expectedEncodedID: `data.base|sche.ma|tab.le`,
		},
		"encodes pointer to schema object identifier": {
			identifier:        new(sdk.NewSchemaObjectIdentifier("database", "schema", "table")),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes table column identifier": {
			identifier:        sdk.NewTableColumnIdentifier("database", "schema", "table", "column"),
			expectedEncodedID: `database|schema|table|column`,
		},
		"encodes quoted table column identifier": {
			identifier:        sdk.NewTableColumnIdentifier("\"database\"", "\"schema\"", "\"table\"", "\"column\""),
			expectedEncodedID: `database|schema|table|column`,
		},
		"encodes table column identifier with dots": {
			identifier:        sdk.NewTableColumnIdentifier("data.base", "sche.ma", "tab.le", "col.umn"),
			expectedEncodedID: `data.base|sche.ma|tab.le|col.umn`,
		},
		"encodes pointer to table column identifier": {
			identifier:        new(sdk.NewTableColumnIdentifier("database", "schema", "table", "column")),
			expectedEncodedID: `database|schema|table|column`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			encodedID := EncodeSnowflakeIDLegacy(tc.identifier)
			require.Equal(t, tc.expectedEncodedID, encodedID)
		})
	}

	t.Run("panics for unsupported object identifier", func(t *testing.T) {
		id := unsupportedObjectIdentifier{}
		require.PanicsWithValue(t, fmt.Sprintf("Unsupported object identifier: %v", id), func() {
			EncodeSnowflakeIDLegacy(id)
		})
	})

	nilTestCases := []any{
		(*sdk.AccountObjectIdentifier)(nil),
		(*sdk.DatabaseObjectIdentifier)(nil),
		(*sdk.SchemaObjectIdentifier)(nil),
		(*sdk.TableColumnIdentifier)(nil),
	}

	for i, tt := range nilTestCases {
		t.Run(fmt.Sprintf("handle nil pointer to object identifier %d", i), func(t *testing.T) {
			require.PanicsWithValue(t, "Nil object identifier received", func() {
				EncodeSnowflakeIDLegacy(tt)
			})
		})
	}
}

type unsupportedObjectIdentifier struct{}

func (i unsupportedObjectIdentifier) Name() string {
	return "name"
}

func (i unsupportedObjectIdentifier) FullyQualifiedName() string {
	return "fully qualified name"
}

func (i unsupportedObjectIdentifier) FullyQualifiedNameEscaped() string {
	return "fully qualified name escaped"
}
