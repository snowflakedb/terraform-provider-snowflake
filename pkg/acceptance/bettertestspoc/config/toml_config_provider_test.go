package config_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TomlConfigProvider(t *testing.T) {
	t.Run("test provider json config", func(t *testing.T) {
		// model := providermodel.SnowflakeProvider().WithAccountName("asdf")
		model := providermodel.SnowflakeProviderToml("default")
		model.AccountName = sdk.Pointer("asdf")
		expectedResult := `{
    "provider": {
        "snowflake": {}
    }
}`

		result, err := config.DefaultTomlConfigProvider.ProviderTomlFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result))
	})
}
