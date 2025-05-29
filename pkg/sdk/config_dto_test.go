package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigDTO_Marshal(t *testing.T) {
	file := NewConfigFile().WithProfiles(map[string]ConfigDTO{
		"default": *NewConfigDTO().WithAccountName("test_account").WithOrganizationName("test_org").WithUser("test_user").WithPassword("test_password").WithRole("test_role").WithWarehouse("test_warehouse"),
	})
	bytes, err := file.MarshalToml()
	require.NoError(t, err)
	require.Equal(t, `[default]
account_name = "test_account"
`, string(bytes))
}
