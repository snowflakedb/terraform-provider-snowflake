//go:build non_account_level_tests

package testint

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1827310]: use generated config for these tests
// TODO [SNOW-2054366]: Use dedicated users for these tests.
func TestInt_Client_NewClient(t *testing.T) {
	restoreDriverLogLevelAfterTest(t)

	t.Run("with default config (legacy)", func(t *testing.T) {
		tmpServiceUser := testClientHelper().SetUpTemporaryServiceUser(t)
		tmpServiceUserConfig := testClientHelper().StoreTempTomlConfigWithProfile(t, testprofiles.Default, func(profile string) string {
			return helpers.FullLegacyTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
		})
		t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		config := sdk.DefaultConfig(sdk.WithUseLegacyTomlFormat(true))
		_, err := sdk.NewClient(config)
		require.NoError(t, err)
	})

	t.Run("with config", func(t *testing.T) {
		tmpServiceUser := testClientHelper().SetUpTemporaryServiceUser(t)
		tmpServiceUserConfig := testClientHelper().StoreTempTomlConfig(t, func(profile string) string {
			return helpers.FullTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
		})
		t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)

		config, err := sdk.ProfileConfig(tmpServiceUserConfig.Profile)
		require.NoError(t, err)
		_, err = sdk.NewClient(config)
		require.NoError(t, err)
	})

	t.Run("with config (legacy)", func(t *testing.T) {
		tmpServiceUser := testClientHelper().SetUpTemporaryServiceUser(t)
		tmpServiceUserConfig := testClientHelper().StoreTempTomlConfig(t, func(profile string) string {
			return helpers.FullLegacyTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
		})
		t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)

		config, err := sdk.ProfileConfig(tmpServiceUserConfig.Profile, sdk.WithUseLegacyTomlFormat(true))
		require.NoError(t, err)
		_, err = sdk.NewClient(config)
		require.NoError(t, err)
	})

	t.Run("with missing config", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config := sdk.DefaultConfig()
		_, err = sdk.NewClient(config)
		require.ErrorContains(t, err, "260000: account is empty")
	})

	t.Run("with incorrect config", func(t *testing.T) {
		tmpServiceUser := testClientHelper().SetUpTemporaryServiceUser(t)
		tmpServiceUserConfig := testClientHelper().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)

		t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)

		config, err := sdk.ProfileConfig(tmpServiceUserConfig.Profile)
		require.NoError(t, err)
		require.NotNil(t, config)

		_, err = sdk.NewClient(config)
		require.ErrorContains(t, err, "JWT token is invalid")
	})

	t.Run("with too big file", func(t *testing.T) {
		c := make([]byte, 11*1024*1024)
		tomlConfig := testClientHelper().StoreTempTomlConfig(t, func(profile string) string {
			return string(c)
		})

		t.Setenv(snowflakeenvs.ConfigPath, tomlConfig.Path)

		_, err := sdk.ProfileConfig(tomlConfig.Profile)
		require.ErrorContains(t, err, fmt.Sprintf("could not load config file: config file %s is too big - maximum allowed size is 10MB", tomlConfig.Path))
	})

	t.Run("with incorrect privileges and enabled check", func(t *testing.T) {
		if oswrapper.IsRunningOnWindows() {
			t.Skip("checking file permissions on Windows is currently done in manual tests package")
		}
		permissions := fs.FileMode(0o755)
		tmpServiceUser := testClientHelper().SetUpTemporaryServiceUser(t)
		tmpServiceUserConfig := testClientHelper().TempTomlConfigWithCustomPermissionsForServiceUser(t, tmpServiceUser, permissions)

		t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)

		_, err := sdk.ProfileConfig(tmpServiceUserConfig.Profile)
		require.ErrorContains(t, err, fmt.Sprintf("could not load config file: config file %s has unsafe permissions - %#o", tmpServiceUserConfig.Path, permissions))
	})

	t.Run("with incorrect privileges and disabled check", func(t *testing.T) {
		if oswrapper.IsRunningOnWindows() {
			t.Skip("checking file permissions on Windows is currently done in manual tests package")
		}
		permissions := fs.FileMode(0o755)
		tmpServiceUser := testClientHelper().SetUpTemporaryServiceUser(t)
		tmpServiceUserConfig := testClientHelper().TempTomlConfigWithCustomPermissionsForServiceUser(t, tmpServiceUser, permissions)

		t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)

		config, err := sdk.ProfileConfig(tmpServiceUserConfig.Profile, sdk.WithVerifyPermissions(false))
		require.NoError(t, err)
		require.NotNil(t, config)

		_, err = sdk.NewClient(config)
		require.NoError(t, err)
	})

	t.Run("with missing config - should not care about correct env variables", func(t *testing.T) {
		config, err := sdk.ProfileConfig(testprofiles.Default)
		require.NoError(t, err)
		require.NotNil(t, config)

		account := config.Account
		parts := strings.Split(account, "-")
		t.Setenv(snowflakeenvs.OrganizationName, parts[0])
		t.Setenv(snowflakeenvs.AccountName, parts[1])

		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config = sdk.DefaultConfig()
		_, err = sdk.NewClient(config)
		require.ErrorContains(t, err, "260000: account is empty")
	})

	t.Run("registers snowflake driver", func(t *testing.T) {
		config := sdk.DefaultConfig()
		_, err := sdk.NewClient(config)
		require.NoError(t, err)

		assert.ElementsMatch(t, sql.Drivers(), []string{"snowflake"})
	})
}

func TestInt_Client_AdditionalMetadata(t *testing.T) {
	client := testClient(t)
	metadata := tracking.Metadata{SchemaVersion: "1", Version: "v1.13.1002-rc-test", Resource: resources.Database.String(), Operation: tracking.CreateOperation}

	assertQueryMetadata := func(t *testing.T, queryId string) {
		t.Helper()
		queryText := testClientHelper().InformationSchema.GetQueryHistoryByQueryId(t, 20, queryId).QueryText
		parsedMetadata, err := tracking.ParseMetadata(queryText)
		require.NoError(t, err)
		require.Equal(t, metadata, parsedMetadata)
	}

	t.Run("query one", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		row := struct {
			One int `db:"ONE"`
		}{}
		err := client.QueryOneForTests(ctx, &row, "SELECT 1 AS ONE")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})

	t.Run("query", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		var rows []struct {
			One int `db:"ONE"`
		}
		err := client.QueryForTests(ctx, &rows, "SELECT 1 AS ONE")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})

	t.Run("exec", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		_, err := client.ExecForTests(ctx, "SELECT 1")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})
}

func TestInt_Client(t *testing.T) {
	restoreDriverLogLevelAfterTest(t)

	t.Run("ping", func(t *testing.T) {
		client := defaultTestClient(t)
		err := client.Ping()
		require.NoError(t, err)
	})

	t.Run("close", func(t *testing.T) {
		client := defaultTestClient(t)
		err := client.Close()
		require.NoError(t, err)
	})

	t.Run("exec", func(t *testing.T) {
		client := defaultTestClient(t)
		ctx := context.Background()
		_, err := client.ExecForTests(ctx, "SELECT 1")
		require.NoError(t, err)
	})

	t.Run("query", func(t *testing.T) {
		client := defaultTestClient(t)
		ctx := context.Background()
		rows := []struct {
			One int `db:"ONE"`
		}{}
		err := client.QueryForTests(ctx, &rows, "SELECT 1 AS ONE")
		require.NoError(t, err)
		require.NotNil(t, rows)
		require.Len(t, rows, 1)
		require.Equal(t, 1, rows[0].One)
	})

	t.Run("queryOne", func(t *testing.T) {
		client := defaultTestClient(t)
		ctx := context.Background()
		row := struct {
			One int `db:"ONE"`
		}{}
		err := client.QueryOneForTests(ctx, &row, "SELECT 1 AS ONE")
		require.NoError(t, err)
		require.Equal(t, 1, row.One)
	})

	// TODO [SNOW-2054366]: Use dedicated users for these tests.
	t.Run("newCLientDriverLoggingLevel", func(t *testing.T) {
		t.Run("get default gosnowflake driver logging level", func(t *testing.T) {
			config := sdk.DefaultConfig()
			_, err := sdk.NewClient(config)
			require.NoError(t, err)

			var expected string
			if os.Getenv("GITHUB_ACTIONS") != "" {
				expected = "FATAL"
			} else {
				expected = "ERROR"
			}
			assert.Equal(t, expected, gosnowflake.GetLogger().GetLogLevel())
		})

		t.Run("set gosnowflake driver logging level with config", func(t *testing.T) {
			config := sdk.DefaultConfig()
			config.Tracing = "trace"
			_, err := sdk.NewClient(config)
			require.NoError(t, err)

			assert.Equal(t, "TRACE", gosnowflake.GetLogger().GetLogLevel())
		})
	})
}

// restoreDriverLogLevelAfterTest snapshots the global gosnowflake logger level and restores it once the test finishes.
// gosnowflake.GetLogger() is a single global logger for the whole test binary. The driver sets the global logger level
// in its init, and every time a connection is opened with a non-empty Tracing the driver overwrites that global level.
// Without restoring it, the level leaks into subsequent tests in the binary.
func restoreDriverLogLevelAfterTest(t *testing.T) {
	t.Helper()
	defaultLogLevel := gosnowflake.GetLogger().GetLogLevel()
	t.Cleanup(func() {
		err := gosnowflake.GetLogger().SetLogLevel(defaultLogLevel)
		require.NoError(t, err)
	})
}

func defaultTestClient(t *testing.T) *sdk.Client {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}
