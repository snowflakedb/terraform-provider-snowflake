package helpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeSnowflakeParameterID(t *testing.T) {
	testCases := map[string]struct {
		id                 string
		fullyQualifiedName string
	}{
		"decodes quoted account object identifier": {
			id:                 `"test.name"`,
			fullyQualifiedName: `"test.name"`,
		},
		"decodes quoted database object identifier": {
			id:                 `"db"."test.name"`,
			fullyQualifiedName: `"db"."test.name"`,
		},
		"decodes quoted schema object identifier": {
			id:                 `"db"."schema"."test.name"`,
			fullyQualifiedName: `"db"."schema"."test.name"`,
		},
		"decodes quoted table column identifier": {
			id:                 `"db"."schema"."table.name"."test.name"`,
			fullyQualifiedName: `"db"."schema"."table.name"."test.name"`,
		},
		"decodes unquoted account object identifier": {
			id:                 `name`,
			fullyQualifiedName: `"name"`,
		},
		"decodes unquoted database object identifier": {
			id:                 `db.name`,
			fullyQualifiedName: `"db"."name"`,
		},
		"decodes unquoted schema object identifier": {
			id:                 `db.schema.name`,
			fullyQualifiedName: `"db"."schema"."name"`,
		},
		"decodes unquoted table column identifier": {
			id:                 `db.schema.table.name`,
			fullyQualifiedName: `"db"."schema"."table"."name"`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			id, err := DecodeSnowflakeParameterID(tc.id)
			require.NoError(t, err)
			require.Equal(t, tc.fullyQualifiedName, id.FullyQualifiedName())
		})
	}

	t.Run("identifier with too many parts", func(t *testing.T) {
		id := `this.identifier.is.too.long.to.be.decoded`
		_, err := DecodeSnowflakeParameterID(id)
		require.ErrorContains(t, err, fmt.Sprintf("unable to classify identifier: %s", id))
	})

	t.Run("incompatible empty identifier", func(t *testing.T) {
		id := ""
		_, err := DecodeSnowflakeParameterID(id)
		require.ErrorContains(t, err, fmt.Sprintf("incompatible identifier: %s", id))
	})

	t.Run("incompatible multiline identifier", func(t *testing.T) {
		id := "db.\nname"
		_, err := DecodeSnowflakeParameterID(id)
		require.ErrorContains(t, err, fmt.Sprintf("unable to read identifier: %s", id))
	})
}

// TODO: add tests for non object identifiers
func TestEncodeSnowflakeID(t *testing.T) {
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
			identifier:        sdk.Pointer(sdk.NewAccountObjectIdentifier("database")),
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
			identifier:        sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database", "schema")),
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
			identifier:        sdk.Pointer(sdk.NewSchemaObjectIdentifier("database", "schema", "table")),
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
			identifier:        sdk.Pointer(sdk.NewTableColumnIdentifier("database", "schema", "table", "column")),
			expectedEncodedID: `database|schema|table|column`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			encodedID := EncodeSnowflakeID(tc.identifier)
			require.Equal(t, tc.expectedEncodedID, encodedID)
		})
	}

	t.Run("panics for unsupported object identifier", func(t *testing.T) {
		id := unsupportedObjectIdentifier{}
		require.PanicsWithValue(t, fmt.Sprintf("Unsupported object identifier: %v", id), func() {
			EncodeSnowflakeID(id)
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
				EncodeSnowflakeID(tt)
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

func Test_DecodeSnowflakeAccountIdentifier(t *testing.T) {
	t.Run("decodes account identifier", func(t *testing.T) {
		id, err := DecodeSnowflakeAccountIdentifier("abc.def")

		require.NoError(t, err)
		require.Equal(t, sdk.NewAccountIdentifier("abc", "def"), id)
	})

	t.Run("does not accept account locator", func(t *testing.T) {
		_, err := DecodeSnowflakeAccountIdentifier("ABC12345")

		require.ErrorContains(t, err, "identifier: ABC12345 seems to be account locator and these are not allowed - please use <organization_name>.<account_name>")
	})

	t.Run("identifier with too many parts", func(t *testing.T) {
		id := `this.identifier.is.too.long.to.be.decoded`
		_, err := DecodeSnowflakeAccountIdentifier(id)

		require.ErrorContains(t, err, fmt.Sprintf("unable to classify account identifier: %s", id))
	})

	t.Run("empty identifier", func(t *testing.T) {
		id := ""
		_, err := DecodeSnowflakeAccountIdentifier(id)

		require.ErrorContains(t, err, fmt.Sprintf("incompatible identifier: %s", id))
	})

	t.Run("multiline identifier", func(t *testing.T) {
		id := "db.\nname"
		_, err := DecodeSnowflakeAccountIdentifier(id)

		require.ErrorContains(t, err, fmt.Sprintf("unable to read identifier: %s", id))
	})
}

func TestParseRootLocation(t *testing.T) {
	tests := []struct {
		name         string
		location     string
		expectedId   string
		expectedPath string
		expectedErr  string
	}{
		{
			name:        "empty",
			location:    ``,
			expectedErr: "incompatible identifier",
		},
		{
			name:       "unquoted",
			location:   `@a.b.c`,
			expectedId: `"a"."b"."c"`,
		},
		{
			name:         "unquoted with path",
			location:     `@a.b.c/foo`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo`,
		},
		{
			name:       "partially quoted",
			location:   `@"a".b.c`,
			expectedId: `"a"."b"."c"`,
		},
		{
			name:         "partially quoted with path",
			location:     `@"a".b.c/foo`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo`,
		},
		{
			name:       "quoted",
			location:   `@"a"."b"."c"`,
			expectedId: `"a"."b"."c"`,
		},
		{
			name:         "quoted with path",
			location:     `@"a"."b"."c"/foo`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo`,
		},
		{
			name:         "unquoted with path with dots",
			location:     `@a.b.c/foo.d`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo.d`,
		},
		{
			name:         "quoted with path with dots",
			location:     `@"a"."b"."c"/foo.d`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo.d`,
		},
		{
			name:         "quoted with complex path",
			location:     `@"a"."b"."c"/foo.a/bar.b//hoge.c`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo.a/bar.b/hoge.c`,
		},
		{
			name:        "invalid location",
			location:    `@foo`,
			expectedErr: "expected 3 parts for location foo, got 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotPath, gotErr := ParseRootLocation(tt.location)
			if len(tt.expectedErr) > 0 {
				assert.ErrorContains(t, gotErr, tt.expectedErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.expectedId, gotId.FullyQualifiedName())
				assert.Equal(t, tt.expectedPath, gotPath)
			}
		})
	}
}

func Test_ContainsIdentifierIgnoreQuotes(t *testing.T) {
	testCases := []struct {
		Name          string
		Ids           []string
		Id            string
		ShouldContain bool
	}{
		{
			Name: "validation: nil Ids",
			Id:   "id",
		},
		{
			Name: "validation: empty Id",
			Ids:  []string{"id"},
			Id:   "",
		},
		{
			Name: "validation: Ids with too many parts",
			Ids:  []string{"this.id.has.too.many.parts"},
			Id:   "id",
		},
		{
			Name: "validation: Id with too many parts",
			Ids:  []string{"id"},
			Id:   "this.id.has.too.many.parts",
		},
		{
			Name: "validation: account object identifier in Ids ignore quotes with upper cased Id",
			Ids:  []string{"object", "db.schema", "db.schema.object"},
			Id:   "\"OBJECT\"",
		},
		{
			Name: "validation: account object identifier in Ids ignore quotes with upper cased id in Ids",
			Ids:  []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:   "\"object\"",
		},
		{
			Name:          "account object identifier in Ids",
			Ids:           []string{"object", "db.schema", "db.schema.object"},
			Id:            "\"object\"",
			ShouldContain: true,
		},
		{
			Name:          "database object identifier in Ids",
			Ids:           []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:            "\"db\".\"schema\"",
			ShouldContain: true,
		},
		{
			Name:          "schema object identifier in Ids",
			Ids:           []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:            "\"db\".\"schema\".\"object\"",
			ShouldContain: true,
		},
		{
			Name:          "account object identifier in Ids upper-cased",
			Ids:           []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:            "\"OBJECT\"",
			ShouldContain: true,
		},
		{
			Name:          "database object identifier in Ids upper-cased",
			Ids:           []string{"object", "DB.SCHEMA", "db.schema.object"},
			Id:            "\"DB\".\"SCHEMA\"",
			ShouldContain: true,
		},
		{
			Name:          "schema object identifier in Ids upper-cased",
			Ids:           []string{"object", "db.schema", "DB.SCHEMA.OBJECT"},
			Id:            "\"DB\".\"SCHEMA\".\"OBJECT\"",
			ShouldContain: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.ShouldContain, ContainsIdentifierIgnoringQuotes(tc.Ids, tc.Id))
		})
	}
}

// External volume helper tests

func Test_ContainsValue(t *testing.T) {
	haystack := sdk.AsStringList([]sdk.AuthenticationMethodsOption{sdk.AuthenticationMethodsAll, sdk.AuthenticationMethodsSaml, sdk.AuthenticationMethodsPassword})
	testCases := []struct {
		name     string
		haystack []string
		needle   string
		expected bool
	}{
		{
			name:     "needle is in haystack",
			haystack: haystack,
			needle:   string(sdk.AuthenticationMethodsSaml),
			expected: true,
		},
		{
			name:     "needle is not in haystack",
			haystack: haystack,
			needle:   "INVALID",
			expected: false,
		},
		{
			name:     "haystack is empty",
			haystack: []string{},
			needle:   string(sdk.AuthenticationMethodsSaml),
			expected: false,
		},
		{
			name:     "needle is empty",
			haystack: haystack,
			needle:   "",
			expected: false,
		},
		{
			name:     "needle is in haystack, normalized",
			haystack: haystack,
			needle:   strings.ToLower(string(sdk.AuthenticationMethodsSaml)),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ContainsEnum(tc.haystack, tc.needle, sdk.ToAuthenticationMethodsOption))
		})
	}
}
