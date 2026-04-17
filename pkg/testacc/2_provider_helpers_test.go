package testacc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

// ConfigurationSameAsStepN should be used to obtain configuration for one of the previous steps to avoid duplication of configuration and var files.
// Based on config.TestStepDirectory.
func ConfigurationSameAsStepN(step int) func(config.TestStepConfigRequest) string {
	return func(req config.TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, strconv.Itoa(step))
	}
}

// ConfigurationDirectory should be used to obtain configuration if the same can be shared between multiple tests to avoid duplication of configuration and var files.
// Based on config.TestNameDirectory. Similar to config.StaticDirectory but prefixed provided directory with `testdata`.
func ConfigurationDirectory(directory string) func(config.TestStepConfigRequest) string {
	return func(req config.TestStepConfigRequest) string {
		return filepath.Join("testdata", directory)
	}
}

// ExternalProviderWithExactVersion returns a map of external providers with an exact version constraint
func ExternalProviderWithExactVersion(version string) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"snowflake": {
			VersionConstraint: fmt.Sprintf("=%s", version),
			Source:            "snowflakedb/snowflake",
		},
	}
}

func setConfigPathEnv(t *testing.T, configName string) {
	t.Helper()
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	configPath := filepath.Join(home, ".snowflake", configName)
	t.Setenv(snowflakeenvs.ConfigPath, configPath)
}

// SetV097CompatibleConfigWithServiceUserPathEnv sets a new config path in a relevant env variable for a file that is compatible with v0.97,
// and authenticates with a service user.
func SetV097CompatibleConfigWithServiceUserPathEnv(t *testing.T) {
	t.Helper()
	setConfigPathEnv(t, "config_v097_compatible_with_service_user")
}

// SetLegacyConfigPathEnv sets a new config path in a relevant env variable for a file that uses the legacy format.
func SetLegacyConfigPathEnv(t *testing.T) {
	t.Helper()
	setConfigPathEnv(t, "config_legacy")
}

// UnsetConfigPathEnv unsets a config path env
func UnsetConfigPathEnv(t *testing.T) {
	t.Helper()
	t.Setenv(snowflakeenvs.ConfigPath, "")
}
