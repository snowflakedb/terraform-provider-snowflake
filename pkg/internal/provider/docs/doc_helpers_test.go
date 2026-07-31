package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PossibleValuesListedStrings(t *testing.T) {
	values := []string{"abc", "DEF"}

	result := PossibleValuesListed(values)

	assert.Equal(t, "`abc` | `DEF`", result)
}

func Test_PossibleValuesListedInts(t *testing.T) {
	values := []int{42, 21}

	result := PossibleValuesListed(values)

	assert.Equal(t, "`42` | `21`", result)
}

func Test_PossibleValuesListed_empty(t *testing.T) {
	var values []string

	result := PossibleValuesListed(values)

	assert.Empty(t, result)
}

func Test_GetDeprecatedObjectReplacements(t *testing.T) {
	testCases := []struct {
		Name               string
		DeprecationMessage string
		Expected           []Replacement
	}{
		{
			Name:               "one replacement listed",
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use one of the new resources instead: `snowflake_storage_integration_aws`.",
			Expected:           []Replacement{{Name: "snowflake_storage_integration_aws"}},
		},
		{
			Name:               "multiple replacements listed",
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use one of the new resources instead: `snowflake_storage_integration_aws` | `snowflake_storage_integration_azure` | `snowflake_storage_integration_gcs`.",
			Expected: []Replacement{
				{Name: "snowflake_storage_integration_aws"},
				{Name: "snowflake_storage_integration_azure"},
				{Name: "snowflake_storage_integration_gcs"},
			},
		},
		{
			Name:               "replacement being a prefix of another replacement",
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use one of the new resources instead: `snowflake_stage_external_s3` | `snowflake_stage_external_s3_compatible`.",
			Expected: []Replacement{
				{Name: "snowflake_stage_external_s3"},
				{Name: "snowflake_stage_external_s3_compatible"},
			},
		},
		{
			Name:               "multiple data source replacements listed",
			DeprecationMessage: "This data source is deprecated and will be removed in a future major version release. Please use one of the new data sources instead: `snowflake_users` | `snowflake_roles`.",
			Expected: []Replacement{
				{Name: "snowflake_users"},
				{Name: "snowflake_roles"},
			},
		},
		{
			Name:               "no replacement",
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release.",
			Expected:           nil,
		},
		{
			Name:               "empty message",
			DeprecationMessage: "",
			Expected:           nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result := GetDeprecatedObjectReplacements(testCase.DeprecationMessage)

			assert.Equal(t, testCase.Expected, result)
		})
	}
}
