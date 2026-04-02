package testenvs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSnowflakeEnvironment(t *testing.T) {
	t.Run("valid values", func(t *testing.T) {
		testCases := map[string]struct {
			input    string
			expected SnowflakeEnvironment
		}{
			"exact PROD": {
				input:    string(SnowflakeProdEnvironment),
				expected: SnowflakeProdEnvironment,
			},
			"lowercase prod": {
				input:    "prod",
				expected: SnowflakeProdEnvironment,
			},
			"exact NON_PROD": {
				input:    string(SnowflakeNonProdEnvironment),
				expected: SnowflakeNonProdEnvironment,
			},
			"lowercase non_prod": {
				input:    "non_prod",
				expected: SnowflakeNonProdEnvironment,
			},
			"mixed case prod": {
				input:    "PrOd",
				expected: SnowflakeProdEnvironment,
			},
			"exact PRE_PROD_GOV": {
				input:    string(SnowflakePreProdGovEnvironment),
				expected: SnowflakePreProdGovEnvironment,
			},
			"lowercase pre_prod_gov": {
				input:    "pre_prod_gov",
				expected: SnowflakePreProdGovEnvironment,
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				actual, err := parseSnowflakeEnvironment(tc.input)
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			})
		}
	})

	t.Run("invalid values", func(t *testing.T) {
		testCases := []string{
			"",
			"dev",
			"prod-env",
			"non prod",
		}

		for _, input := range testCases {
			input := input
			t.Run(input, func(t *testing.T) {
				_, err := parseSnowflakeEnvironment(input)
				require.Error(t, err)
				require.ErrorContains(t, err, "invalid Snowflake environment")
			})
		}
	})
}
