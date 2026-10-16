package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserSchema_WriteOnlyRsaPublicKeys(t *testing.T) {
	resourcesToCheck := map[string]map[string]*schema.Schema{
		"snowflake_user":                User().Schema,
		"snowflake_service_user":        ServiceUser().Schema,
		"snowflake_legacy_service_user": LegacyServiceUser().Schema,
	}

	keyPairs := []struct {
		field        string
		writeOnly    string
		writeOnlyVer string
	}{
		{"rsa_public_key", "rsa_public_key_wo", "rsa_public_key_wo_version"},
		{"rsa_public_key_2", "rsa_public_key_2_wo", "rsa_public_key_2_wo_version"},
	}

	for resourceName, resourceSchema := range resourcesToCheck {
		for _, keyPair := range keyPairs {
			t.Run(fmt.Sprintf("%s: %s", resourceName, keyPair.writeOnly), func(t *testing.T) {
				writeOnly, ok := resourceSchema[keyPair.writeOnly]
				require.True(t, ok, "%s should have the %s field", resourceName, keyPair.writeOnly)

				assert.Equal(t, schema.TypeString, writeOnly.Type)
				assert.True(t, writeOnly.WriteOnly, "%s should be write-only", keyPair.writeOnly)
				assert.True(t, writeOnly.Optional)

				assert.False(t, writeOnly.Computed)
				assert.False(t, writeOnly.ForceNew)
				assert.Nil(t, writeOnly.Default)
				assert.ElementsMatch(t, []string{keyPair.field}, writeOnly.ConflictsWith)
				assert.ElementsMatch(t, []string{keyPair.writeOnlyVer}, writeOnly.RequiredWith)
			})

			t.Run(fmt.Sprintf("%s: %s", resourceName, keyPair.writeOnlyVer), func(t *testing.T) {
				version, ok := resourceSchema[keyPair.writeOnlyVer]
				require.True(t, ok, "%s should have the %s field", resourceName, keyPair.writeOnlyVer)

				assert.Equal(t, schema.TypeInt, version.Type)
				assert.True(t, version.Optional)

				assert.False(t, version.WriteOnly, "%s must not be write-only", keyPair.writeOnlyVer)
				assert.ElementsMatch(t, []string{keyPair.writeOnly}, version.RequiredWith)
			})

			t.Run(fmt.Sprintf("%s: %s conflicts with its write-only variant", resourceName, keyPair.field), func(t *testing.T) {
				field, ok := resourceSchema[keyPair.field]
				require.True(t, ok, "%s should have the %s field", resourceName, keyPair.field)

				assert.False(t, field.WriteOnly)
				assert.ElementsMatch(t, []string{keyPair.writeOnly}, field.ConflictsWith)
			})
		}
	}
}

func Test_writeOnlyStringValue(t *testing.T) {
	testCases := []struct {
		name     string
		value    cty.Value
		expected *string
	}{
		{
			name:     "null value",
			value:    cty.NullVal(cty.String),
			expected: nil,
		},
		{
			name:     "unknown value",
			value:    cty.UnknownVal(cty.String),
			expected: nil,
		},
		{
			name:     "empty string is treated as unset",
			value:    cty.StringVal(""),
			expected: nil,
		},
		{
			name:     "regular value",
			value:    cty.StringVal("some-public-key"),
			expected: sdkString("some-public-key"),
		},
		{
			name:     "value with whitespace is not trimmed",
			value:    cty.StringVal(" some-public-key\n"),
			expected: sdkString(" some-public-key\n"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, writeOnlyStringValue(testCase.value))
		})
	}
}

func sdkString(s string) *string {
	return &s
}
