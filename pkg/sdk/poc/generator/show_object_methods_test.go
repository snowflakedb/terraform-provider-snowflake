package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdentifierStringToObjectIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected objectIdentifierKind
	}{
		{"AccountObjectIdentifier", AccountObjectIdentifier},
		{"DatabaseObjectIdentifier", DatabaseObjectIdentifier},
		{"SchemaObjectIdentifier", SchemaObjectIdentifier},
		{"SchemaObjectIdentifierWithArguments", SchemaObjectIdentifierWithArguments},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ToObjectIdentifierKind(test.input)
			require.NoError(t, err)
			require.Equal(t, test.expected, result)
		})
	}
}

func TestIdentifierStringToObjectIdentifier_Invalid(t *testing.T) {
	tests := []struct {
		input string
		err   string
	}{
		{"accountobjectidentifier", "invalid string identifier type: accountobjectidentifier"},
		{"Account", "invalid string identifier type: Account"},
		{"databaseobjectidentifier", "invalid string identifier type: databaseobjectidentifier"},
		{"Database", "invalid string identifier type: Database"},
		{"schemaobjectidentifier", "invalid string identifier type: schemaobjectidentifier"},
		{"Schema", "invalid string identifier type: Schema"},
		{"schemaobjectidentifierwitharguments", "invalid string identifier type: schemaobjectidentifierwitharguments"},
		{"schemawitharguemnts", "invalid string identifier type: schemawitharguemnts"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToObjectIdentifierKind(tc.input)
			require.ErrorContains(t, err, tc.err)
		})
	}
}
