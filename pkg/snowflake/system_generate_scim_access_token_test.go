package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGenerateSCIMAccessToken(t *testing.T) {
	testCases := []struct {
		name            string
		integrationName string
		expected        string
	}{
		{
			name:            "basic",
			integrationName: "AAD_PROVISIONING",
			expected:        `SELECT SYSTEM$GENERATE_SCIM_ACCESS_TOKEN('AAD_PROVISIONING') AS "TOKEN"`,
		},
		{
			name:            "single quote is escaped",
			integrationName: "AAD'PROVISIONING",
			expected:        `SELECT SYSTEM$GENERATE_SCIM_ACCESS_TOKEN('AAD\'PROVISIONING') AS "TOKEN"`,
		},
		{
			name:            "backslash is escaped",
			integrationName: `AAD\PROVISIONING`,
			expected:        `SELECT SYSTEM$GENERATE_SCIM_ACCESS_TOKEN('AAD\\PROVISIONING') AS "TOKEN"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			sb := NewSystemGenerateSCIMAccessTokenBuilder(tc.integrationName)

			r.Equal(tc.expected, sb.Select())
		})
	}
}
