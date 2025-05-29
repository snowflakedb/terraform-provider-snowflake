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
