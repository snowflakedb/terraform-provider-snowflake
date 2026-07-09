package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccountIdentifierFromFullyQualifiedName(t *testing.T) {
	type test struct {
		input string
		want  AccountIdentifier
	}

	tests := []test{
		{input: "BSB98216", want: AccountIdentifier{accountLocator: "BSB98216"}},
		{input: "SNOW.MY_TEST_ACCOUNT", want: AccountIdentifier{organizationName: "SNOW", accountName: "MY_TEST_ACCOUNT"}},
		{input: "\"SNOW\".\"MY_TEST_ACCOUNT\"", want: AccountIdentifier{organizationName: "SNOW", accountName: "MY_TEST_ACCOUNT"}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			id := NewAccountIdentifierFromFullyQualifiedName(tc.input)
			require.Equal(t, tc.want, id)
		})
	}
}

func TestNewSchemaObjectIdentifierFromFullyQualifiedName(t *testing.T) {
	type test struct {
		input string
		want  SchemaObjectIdentifier
	}

	tests := []test{
		{input: "\"MY_DB\".\"MY_SCHEMA\".\"multiply\"(number, number)", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "multiply", arguments: []DataType{DataTypeNumber, DataTypeNumber}}},
		{input: "MY_DB.MY_SCHEMA.add(number, number)", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "add", arguments: []DataType{DataTypeNumber, DataTypeNumber}}},
		{input: "\"MY_DB\".\"MY_SCHEMA\".\"MY_UDF\"()", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "MY_UDF", arguments: []DataType{}}},
		{input: "\"MY_DB\".\"MY_SCHEMA\".\"MY_PIPE\"", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "MY_PIPE", arguments: nil}},
		{input: "MY_DB.MY_SCHEMA.MY_STAGE", want: SchemaObjectIdentifier{databaseName: "MY_DB", schemaName: "MY_SCHEMA", name: "MY_STAGE", arguments: nil}},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			id := NewSchemaObjectIdentifierFromFullyQualifiedName(tc.input)
			require.Equal(t, tc.want, id)
		})
	}
}

func TestDatabaseObjectIdentifier(t *testing.T) {
	t.Run("create from strings", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifier("aaa", "bbb")

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("create from quoted strings", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifier(`"aaa"`, `"bbb"`)

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("create from fully qualified name", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifierFromFullyQualifiedName("aaa.bbb")

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("create from quoted fully qualified name", func(t *testing.T) {
		identifier := NewDatabaseObjectIdentifierFromFullyQualifiedName(`"aaa"."bbb"`)

		assert.Equal(t, "aaa", identifier.DatabaseName())
		assert.Equal(t, "bbb", identifier.Name())
	})

	t.Run("get fully qualified name", func(t *testing.T) {
		identifier := DatabaseObjectIdentifier{"aaa", "bbb"}

		assert.Equal(t, `"aaa"."bbb"`, identifier.FullyQualifiedName())
	})
}

// The tests below verify how object identifiers emit SQL when the underlying
// names contain "weird" characters such as single and double quotes.
func TestAccountObjectIdentifier_Injection(t *testing.T) {
	tests := []struct {
		name        string
		id          AccountObjectIdentifier
		wantFQN     string
		wantEscaped string
	}{
		{
			name:        "regular name",
			id:          AccountObjectIdentifier{name: "MY_DB"},
			wantFQN:     `"MY_DB"`,
			wantEscaped: `"MY_DB"`,
		},
		{
			name:        "embedded double quote",
			id:          AccountObjectIdentifier{name: `my"db`},
			wantFQN:     `"my"db"`,  // naive wrapping breaks out of quoting
			wantEscaped: `"my""db"`, // escaped: the double quote is doubled
		},
		{
			name:        "embedded single quote",
			id:          AccountObjectIdentifier{name: `my'db`},
			wantFQN:     `"my'db"`,
			wantEscaped: `"my'db"`, // single quotes are irrelevant inside double-quoted identifiers
		},
		{
			name:        "embedded backslash",
			id:          AccountObjectIdentifier{name: `my\db`},
			wantFQN:     `"my\db"`,
			wantEscaped: `"my\db"`,
		},
		{
			name:        "injection attempt",
			id:          AccountObjectIdentifier{name: `a"; DROP TABLE t; --`},
			wantFQN:     `"a"; DROP TABLE t; --"`,
			wantEscaped: `"a""; DROP TABLE t; --"`,
		},
		{
			name:        "empty name",
			id:          AccountObjectIdentifier{name: ""},
			wantFQN:     ``,
			wantEscaped: ``,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantFQN, tc.id.FullyQualifiedName())
			assert.Equal(t, tc.wantEscaped, tc.id.FullyQualifiedNameEscaped())
		})
	}
}

func TestDatabaseObjectIdentifier_Injection(t *testing.T) {
	tests := []struct {
		name        string
		id          DatabaseObjectIdentifier
		wantFQN     string
		wantEscaped string
	}{
		{
			name:        "regular names",
			id:          DatabaseObjectIdentifier{databaseName: "DB", name: "OBJ"},
			wantFQN:     `"DB"."OBJ"`,
			wantEscaped: `"DB"."OBJ"`,
		},
		{
			name:        "double quote in name",
			id:          DatabaseObjectIdentifier{databaseName: "DB", name: `a"b`},
			wantFQN:     `"DB"."a"b"`,
			wantEscaped: `"DB"."a""b"`,
		},
		{
			name:        "double quote in database name",
			id:          DatabaseObjectIdentifier{databaseName: `d"b`, name: "OBJ"},
			wantFQN:     `"d"b"."OBJ"`,
			wantEscaped: `"d""b"."OBJ"`,
		},
		{
			name:        "single quote in name",
			id:          DatabaseObjectIdentifier{databaseName: "DB", name: `a'b`},
			wantFQN:     `"DB"."a'b"`,
			wantEscaped: `"DB"."a'b"`,
		},
		{
			name:        "injection attempt in name",
			id:          DatabaseObjectIdentifier{databaseName: "DB", name: `a"; DROP TABLE t; --`},
			wantFQN:     `"DB"."a"; DROP TABLE t; --"`,
			wantEscaped: `"DB"."a""; DROP TABLE t; --"`,
		},
		{
			name:        "empty",
			id:          DatabaseObjectIdentifier{},
			wantFQN:     ``,
			wantEscaped: ``,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantFQN, tc.id.FullyQualifiedName())
			assert.Equal(t, tc.wantEscaped, tc.id.FullyQualifiedNameEscaped())
		})
	}
}

func TestSchemaObjectIdentifier_Injection(t *testing.T) {
	tests := []struct {
		name        string
		id          SchemaObjectIdentifier
		wantFQN     string
		wantEscaped string
	}{
		{
			name:        "regular names",
			id:          SchemaObjectIdentifier{databaseName: "DB", schemaName: "SC", name: "OBJ"},
			wantFQN:     `"DB"."SC"."OBJ"`,
			wantEscaped: `"DB"."SC"."OBJ"`,
		},
		{
			name:        "double quote in name",
			id:          SchemaObjectIdentifier{databaseName: "DB", schemaName: "SC", name: `a"b`},
			wantFQN:     `"DB"."SC"."a"b"`,
			wantEscaped: `"DB"."SC"."a""b"`,
		},
		{
			name:        "double quote in schema name",
			id:          SchemaObjectIdentifier{databaseName: "DB", schemaName: `s"c`, name: "OBJ"},
			wantFQN:     `"DB"."s"c"."OBJ"`,
			wantEscaped: `"DB"."s""c"."OBJ"`,
		},
		{
			name:        "single quote in name",
			id:          SchemaObjectIdentifier{databaseName: "DB", schemaName: "SC", name: `a'b`},
			wantFQN:     `"DB"."SC"."a'b"`,
			wantEscaped: `"DB"."SC"."a'b"`,
		},
		{
			name:        "injection attempt in name",
			id:          SchemaObjectIdentifier{databaseName: "DB", schemaName: "SC", name: `a"; DROP TABLE t; --`},
			wantFQN:     `"DB"."SC"."a"; DROP TABLE t; --"`,
			wantEscaped: `"DB"."SC"."a""; DROP TABLE t; --"`,
		},
		{
			name:        "empty",
			id:          SchemaObjectIdentifier{},
			wantFQN:     ``,
			wantEscaped: ``,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantFQN, tc.id.FullyQualifiedName())
			assert.Equal(t, tc.wantEscaped, tc.id.FullyQualifiedNameEscaped())
		})
	}
}

func TestExternalObjectIdentifier_Injection(t *testing.T) {
	tests := []struct {
		name        string
		id          ExternalObjectIdentifier
		wantFQN     string
		wantEscaped string
	}{
		{
			name:        "regular names with account locator",
			id:          NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("LOC"), NewAccountObjectIdentifier("OBJ")),
			wantFQN:     `"LOC"."OBJ"`,
			wantEscaped: `"LOC"."OBJ"`,
		},
		{
			name:        "regular names with organization and account name",
			id:          NewExternalObjectIdentifier(NewAccountIdentifier("ORG", "ACC"), NewAccountObjectIdentifier("OBJ")),
			wantFQN:     `"ORG"."ACC"."OBJ"`,
			wantEscaped: `"ORG"."ACC"."OBJ"`,
		},
		{
			// NOTE: FullyQualifiedNameEscaped() currently delegates to FullyQualifiedName()
			// and does NOT escape embedded double quotes - this is an injection gap (SNOW-3696846).
			name:        "double quote in object name is NOT escaped",
			id:          NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("LOC"), AccountObjectIdentifier{name: `o"b`}),
			wantFQN:     `"LOC"."o"b"`,
			wantEscaped: `"LOC"."o"b"`,
		},
		{
			// NOTE: the account locator is likewise emitted unescaped (SNOW-3696846).
			name:        "double quote in account locator is NOT escaped",
			id:          NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator(`l"c`), NewAccountObjectIdentifier("OBJ")),
			wantFQN:     `"l"c"."OBJ"`,
			wantEscaped: `"l"c"."OBJ"`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantFQN, tc.id.FullyQualifiedName())
			assert.Equal(t, tc.wantEscaped, tc.id.FullyQualifiedNameEscaped())
		})
	}
}
