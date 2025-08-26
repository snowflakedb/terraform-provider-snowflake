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

// TODO [SNOW-2298287]: function acceptance.TestAccPreCheck was needed for the database, schema, and warehouse creation; it's not needed when TestMain is used; it was left for now, as we may want to do other validations here (and not to do a big revolution in all the acceptance tests that will be copied)
func TestAccPreCheck(t *testing.T) {
	t.Helper()
}

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

// SetV097CompatibleConfigPathEnv sets a new config path in a relevant env variable for a file that is compatible with v0.97.
func SetV097CompatibleConfigPathEnv(t *testing.T) {
	t.Helper()
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	configPath := filepath.Join(home, ".snowflake", "config_v097_compatible")
	t.Setenv(snowflakeenvs.ConfigPath, configPath)
}

// SetLegacyConfigPathEnv sets a new config path in a relevant env variable for a file that uses the legacy format.
func SetLegacyConfigPathEnv(t *testing.T) {
	t.Helper()
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	configPath := filepath.Join(home, ".snowflake", "config_legacy")
	t.Setenv(snowflakeenvs.ConfigPath, configPath)
}

// UnsetConfigPathEnv unsets a config path env
func UnsetConfigPathEnv(t *testing.T) {
	t.Helper()
	t.Setenv(snowflakeenvs.ConfigPath, "")
}
