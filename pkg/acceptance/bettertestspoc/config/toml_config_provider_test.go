package config_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/tomlconfigmodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TomlConfigProvider(t *testing.T) {
	testCases := []struct {
		name     string
		model    *tomlconfigmodel.SnowflakeConfigModel
		expected string
	}{
		{
			name:  "simple config",
			model: tomlconfigmodel.SnowflakeTomlConfig("default").WithAccountName(sdk.Pointer("account")),
			expected: `[default]
account_name = 'account'
`,
		},
		{
			name: "more fields",
			model: tomlconfigmodel.SnowflakeTomlConfig("default").
				WithAccountName(sdk.Pointer("account")).
				WithAuthenticator(sdk.Pointer("externalbrowser")).
				WithOrganizationName(sdk.Pointer("org")).
				WithUser(sdk.Pointer("user")).
				WithUsername(sdk.Pointer("username")).
				WithPassword(sdk.Pointer("password")).
				WithHost(sdk.Pointer("host")).
				WithWarehouse(sdk.Pointer("warehouse")).
				WithMaxRetryCount(sdk.Pointer(5)),
			expected: `[default]
account_name = 'account'
organization_name = 'org'
user = 'user'
username = 'username'
password = 'password'
host = 'host'
warehouse = 'warehouse'
max_retry_count = 5
authenticator = 'externalbrowser'
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := config.DefaultTomlConfigProvider.ProviderTomlFromModel(tc.model)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
