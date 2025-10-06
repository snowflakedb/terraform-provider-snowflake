//go:build !account_level_tests

package testacc

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"regexp"
	"testing"
	"time"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testfiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcc_Provider_configHierarchy(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)
	incorrectConfig := testClient().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// make sure that we fail for incorrect profile
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// make sure that we succeed for the correct profile
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", TestDatabaseName),
				),
			},
			// incorrect user in provider config should not be rewritten by profile and cause error
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithAuthenticatorType(sdk.AuthenticationTypeJwt).WithProfile(tmpServiceUserConfig.Profile).WithUserId(NonExistingAccountObjectIdentifier), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// correct user and key in provider's config should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().
					WithAuthenticatorType(sdk.AuthenticationTypeJwt).
					WithProfile(incorrectConfig.Profile).
					WithUserId(tmpServiceUser.UserId).
					WithRoleId(tmpServiceUser.RoleId).
					WithPrivateKeyMultiline(tmpServiceUser.PrivateKey), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", TestDatabaseName),
				),
			},
			// incorrect user in env variable should not be rewritten by profile and cause error (profile authenticator is set to JWT and that's why the error is about incorrect token)
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, NonExistingAccountObjectIdentifier.Name())
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// correct user and private key in env should not be rewritten by a faulty config
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.User, tmpServiceUser.UserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, tmpServiceUser.PrivateKey)
					t.Setenv(snowflakeenvs.Role, tmpServiceUser.RoleId.Name())
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", TestDatabaseName),
				),
			},
			// user on provider level wins (it's incorrect - env and profile ones are)
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.User)
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().
					WithAuthenticatorType(sdk.AuthenticationTypeJwt).
					WithProfile(tmpServiceUserConfig.Profile).
					WithUserId(NonExistingAccountObjectIdentifier).
					WithRoleId(tmpServiceUser.RoleId).
					WithPrivateKeyMultiline(tmpServiceUser.PrivateKey), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// there is no config (by setting the dir to something different from .snowflake/config)
			{
				PreConfig: func() {
					dir, err := os.UserHomeDir()
					require.NoError(t, err)
					t.Setenv(snowflakeenvs.ConfigPath, dir)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().
					WithAuthenticatorType(sdk.AuthenticationTypeJwt).
					WithProfile(testprofiles.Default).
					WithUserId(tmpServiceUser.UserId).
					WithRoleId(tmpServiceUser.RoleId).
					WithPrivateKeyMultiline(tmpServiceUser.PrivateKey), datasourceModel()),
				ExpectError: regexp.MustCompile("account is empty"),
			},
			// provider's config should not be rewritten by env when there is no profile (incorrect user in config versus correct one in env) - proves #2242
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.ConfigPath)
					t.Setenv(snowflakeenvs.User, tmpServiceUser.UserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, tmpServiceUser.PrivateKey)
					t.Setenv(snowflakeenvs.AccountName, tmpServiceUser.AccountId.AccountName())
					t.Setenv(snowflakeenvs.OrganizationName, tmpServiceUser.AccountId.OrganizationName())
					t.Setenv(snowflakeenvs.Role, tmpServiceUser.RoleId.Name())
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithAuthenticatorType(sdk.AuthenticationTypeJwt).WithProfile(testprofiles.Default).WithUserId(NonExistingAccountObjectIdentifier), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// make sure the teardown is fine by using a correct env config at the end
			{
				PreConfig: func() {
					testenvs.AssertEnvSet(t, snowflakeenvs.ConfigPath)
					testenvs.AssertEnvSet(t, snowflakeenvs.User)
					testenvs.AssertEnvSet(t, snowflakeenvs.PrivateKey)
					testenvs.AssertEnvSet(t, snowflakeenvs.AccountName)
					testenvs.AssertEnvSet(t, snowflakeenvs.OrganizationName)
					testenvs.AssertEnvSet(t, snowflakeenvs.Role)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", TestDatabaseName),
				),
			},
		},
	})
}

func TestAcc_Provider_configureClientOnceSwitching(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)
	incorrectConfig := testClient().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// client setup is incorrect
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, incorrectConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// in this step we simulate the situation when we want to use client configured once, but it was faulty last time
			{
				PreConfig: func() {
					t.Setenv(string(testenvs.ConfigureClientOnce), "true")
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", TestDatabaseName),
				),
			},
		},
	})
}

func TestAcc_Provider_LegacyTomlConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullLegacyTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithUseLegacyTomlFile(true), datasourceModel()),
				Check: func(s *terraform.State) error {
					config := TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()
					assert.Equal(t, tmpServiceUser.OrgAndAccount(), config.Account)
					assert.Equal(t, tmpServiceUser.UserId.Name(), config.User)
					assert.Equal(t, tmpServiceUser.WarehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpServiceUser.RoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("1.2.3.4"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", tmpServiceUser.OrgAndAccount()), config.Host)
					assert.Equal(t, 443, config.Port)
					assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
					assert.False(t, config.PasscodeInPassword)
					assert.Equal(t, testvars.ExampleOktaUrl, config.OktaURL)
					assert.Equal(t, 30*time.Second, config.LoginTimeout)
					assert.Equal(t, 40*time.Second, config.RequestTimeout)
					assert.Equal(t, 50*time.Second, config.JWTExpireTimeout)
					assert.Equal(t, 10*time.Second, config.ClientTimeout)
					assert.Equal(t, 20*time.Second, config.JWTClientTimeout)
					assert.Equal(t, 60*time.Second, config.ExternalBrowserTimeout)
					assert.Equal(t, 1, config.MaxRetryCount)
					assert.Equal(t, "terraform-provider-snowflake", config.Application)
					assert.True(t, config.InsecureMode) //nolint:staticcheck
					assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
					assert.Equal(t, "token", config.Token)
					assert.True(t, config.KeepSessionAlive)
					assert.True(t, config.DisableTelemetry)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), config.Tracing)
					assert.Equal(t, ".", config.TmpDirPath)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
					assert.True(t, config.DisableQueryContextCache)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
					assert.Equal(t, map[string]*string{
						"foo": sdk.Pointer("bar"),
					}, config.Params)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), gosnowflake.GetLogger().GetLogLevel())

					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_TomlConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				Check: func(s *terraform.State) error {
					config := TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()
					assert.Equal(t, tmpServiceUser.OrgAndAccount(), config.Account)
					assert.Equal(t, tmpServiceUser.UserId.Name(), config.User)
					assert.Equal(t, tmpServiceUser.WarehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpServiceUser.RoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("1.2.3.4"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", tmpServiceUser.OrgAndAccount()), config.Host)
					assert.Equal(t, 443, config.Port)
					assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
					assert.False(t, config.PasscodeInPassword)
					assert.Equal(t, testvars.ExampleOktaUrl, config.OktaURL)
					assert.Equal(t, 30*time.Second, config.LoginTimeout)
					assert.Equal(t, 40*time.Second, config.RequestTimeout)
					assert.Equal(t, 50*time.Second, config.JWTExpireTimeout)
					assert.Equal(t, 10*time.Second, config.ClientTimeout)
					assert.Equal(t, 20*time.Second, config.JWTClientTimeout)
					assert.Equal(t, 60*time.Second, config.ExternalBrowserTimeout)
					assert.Equal(t, 1, config.MaxRetryCount)
					assert.Equal(t, "terraform-provider-snowflake", config.Application)
					assert.True(t, config.InsecureMode) //nolint:staticcheck
					assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
					assert.Equal(t, "token", config.Token)
					assert.True(t, config.KeepSessionAlive)
					assert.True(t, config.DisableTelemetry)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), config.Tracing)
					assert.Equal(t, ".", config.TmpDirPath)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
					assert.True(t, config.DisableQueryContextCache)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
					assert.Equal(t, map[string]*string{
						"foo": sdk.Pointer("bar"),
					}, config.Params)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), gosnowflake.GetLogger().GetLogLevel())
					assert.Equal(t, "oauth_client_id", config.OauthClientID)
					assert.Equal(t, "oauth_client_secret", config.OauthClientSecret)
					assert.Equal(t, "oauth_token_request_url", config.OauthTokenRequestURL)
					assert.Equal(t, "oauth_authorization_url", config.OauthAuthorizationURL)
					assert.Equal(t, "oauth_redirect_uri", config.OauthRedirectURI)
					assert.Equal(t, "oauth_scope", config.OauthScope)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.EnableSingleUseRefreshTokens)

					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_TomlConfigFailsIfFormatsMismatch(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
	})

	legacyTomlTmpServiceUserConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullLegacyTomlConfigForServiceUser(t, profile, tmpServiceUser.UserId, tmpServiceUser.RoleId, tmpServiceUser.WarehouseId, tmpServiceUser.AccountId, tmpServiceUser.PrivateKey)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Try reading the legacy format, but provide a file with the new format.
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithUseLegacyTomlFile(true), datasourceModel()),
				ExpectError: regexp.MustCompile("account is empty"),
			},
			// Try reading the new format, but provide a file with the legacy format.
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, legacyTomlTmpServiceUserConfig.Path)
					t.Setenv(string(snowflakeenvs.UseLegacyTomlFile), "")
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(legacyTomlTmpServiceUserConfig.Profile).WithUseLegacyTomlFile(false), datasourceModel()),
				ExpectError: regexp.MustCompile("account is empty"),
			},
		},
	})
}

func TestAcc_Provider_tomlConfigIsTooBig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	c := make([]byte, 11*1024*1024)
	tomlConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return string(c)
	})
	providerModel := providermodel.SnowflakeProvider().WithProfile(tomlConfig.Profile)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tomlConfig.Path)
				},
				Config:      config.FromModels(t, providerModel, datasourceModel()),
				ExpectError: regexp.MustCompile(fmt.Sprintf("could not load config file: config file %s is too big - maximum allowed size is 10MB", tomlConfig.Path)),
			},
		},
	})
}

func TestAcc_Provider_tomlConfigIsTooPermissive(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	permissions := fs.FileMode(0o755)

	configPath := testfiles.TestFileWithCustomPermissions(t, random.AlphaN(10), random.Bytes(), permissions)
	providerModel := providermodel.SnowflakeProvider().WithProfile(configPath)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, configPath)
				},
				Config:      config.FromModels(t, providerModel, datasourceModel()),
				ExpectError: regexp.MustCompile(fmt.Sprintf("could not load config file: config file %s has unsafe permissions - %#o", configPath, permissions)),
			},
		},
	})
}

func TestAcc_Provider_tomlConfigFilePermissionsCanBeSkipped(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigWithCustomPermissionsForServiceUser(t, tmpServiceUser, fs.FileMode(0o755))

	providerModelWithSkippedPermissionVerification := providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithSkipTomlFilePermissionVerification(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providerModelWithSkippedPermissionVerification, datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", TestDatabaseName),
				),
			},
		},
	})
}

func TestAcc_Provider_envConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullInvalidTomlConfigForServiceUser(t, profile)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.AccountName, tmpServiceUser.AccountId.AccountName())
					t.Setenv(snowflakeenvs.OrganizationName, tmpServiceUser.AccountId.OrganizationName())
					t.Setenv(snowflakeenvs.User, tmpServiceUser.UserId.Name())
					t.Setenv(snowflakeenvs.PrivateKey, tmpServiceUser.PrivateKey)
					t.Setenv(snowflakeenvs.Warehouse, tmpServiceUser.WarehouseId.Name())
					t.Setenv(snowflakeenvs.Protocol, "https")
					t.Setenv(snowflakeenvs.Port, "443")
					// do not set token - it should be propagated from TOML
					t.Setenv(snowflakeenvs.Role, tmpServiceUser.RoleId.Name())
					t.Setenv(snowflakeenvs.Authenticator, "SNOWFLAKE_JWT")
					t.Setenv(snowflakeenvs.ValidateDefaultParameters, "true")
					t.Setenv(snowflakeenvs.ClientIp, "2.2.2.2")
					t.Setenv(snowflakeenvs.Host, "")
					t.Setenv(snowflakeenvs.Passcode, "")
					t.Setenv(snowflakeenvs.PasscodeInPassword, "false")
					t.Setenv(snowflakeenvs.OktaUrl, testvars.ExampleOktaUrlFromEnvString)
					t.Setenv(snowflakeenvs.LoginTimeout, "100")
					t.Setenv(snowflakeenvs.RequestTimeout, "200")
					t.Setenv(snowflakeenvs.JwtExpireTimeout, "300")
					t.Setenv(snowflakeenvs.ClientTimeout, "400")
					t.Setenv(snowflakeenvs.JwtClientTimeout, "500")
					t.Setenv(snowflakeenvs.ExternalBrowserTimeout, "600")
					t.Setenv(snowflakeenvs.InsecureMode, "false")
					t.Setenv(snowflakeenvs.OcspFailOpen, "false")
					t.Setenv(snowflakeenvs.KeepSessionAlive, "false")
					t.Setenv(snowflakeenvs.DisableTelemetry, "false")
					t.Setenv(snowflakeenvs.ClientRequestMfaToken, "false")
					t.Setenv(snowflakeenvs.ClientStoreTemporaryCredential, "false")
					t.Setenv(snowflakeenvs.DisableQueryContextCache, "false")
					t.Setenv(snowflakeenvs.IncludeRetryReason, "false")
					t.Setenv(snowflakeenvs.MaxRetryCount, "2")
					t.Setenv(snowflakeenvs.DriverTracing, string(sdk.DriverLogLevelWarning))
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
					t.Setenv(snowflakeenvs.OauthClientId, "oauth_client_id")
					t.Setenv(snowflakeenvs.OauthClientSecret, "oauth_client_secret")
					t.Setenv(snowflakeenvs.OauthTokenRequestUrl, "oauth_token_request_url")
					t.Setenv(snowflakeenvs.OauthAuthorizationUrl, "oauth_authorization_url")
					t.Setenv(snowflakeenvs.OauthRedirectUri, "oauth_redirect_uri")
					t.Setenv(snowflakeenvs.OauthScope, "oauth_scope")
					t.Setenv(snowflakeenvs.EnableSingleUseRefreshTokens, "true")
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				Check: func(s *terraform.State) error {
					config := TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()

					assert.Equal(t, tmpServiceUser.OrgAndAccount(), config.Account)
					assert.Equal(t, tmpServiceUser.UserId.Name(), config.User)
					assert.Equal(t, tmpServiceUser.WarehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpServiceUser.RoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("2.2.2.2"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", tmpServiceUser.OrgAndAccount()), config.Host)
					assert.Equal(t, 443, config.Port)
					assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
					assert.False(t, config.PasscodeInPassword)
					assert.Equal(t, testvars.ExampleOktaUrlFromEnv, config.OktaURL)
					assert.Equal(t, 100*time.Second, config.LoginTimeout)
					assert.Equal(t, 200*time.Second, config.RequestTimeout)
					assert.Equal(t, 300*time.Second, config.JWTExpireTimeout)
					assert.Equal(t, 400*time.Second, config.ClientTimeout)
					assert.Equal(t, 500*time.Second, config.JWTClientTimeout)
					assert.Equal(t, 600*time.Second, config.ExternalBrowserTimeout)
					assert.Equal(t, 2, config.MaxRetryCount)
					assert.Equal(t, "terraform-provider-snowflake", config.Application)
					assert.True(t, config.InsecureMode) //nolint:staticcheck
					assert.Equal(t, gosnowflake.OCSPFailOpenFalse, config.OCSPFailOpen)
					assert.Equal(t, "token", config.Token)
					assert.True(t, config.KeepSessionAlive)
					assert.True(t, config.DisableTelemetry)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), config.Tracing)
					assert.Equal(t, "../", config.TmpDirPath)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.ClientRequestMfaToken)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.ClientStoreTemporaryCredential)
					assert.True(t, config.DisableQueryContextCache)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.IncludeRetryReason)
					assert.Equal(t, gosnowflake.ConfigBoolFalse, config.DisableConsoleLogin)
					assert.Equal(t, map[string]*string{
						"foo": sdk.Pointer("bar"),
					}, config.Params)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), gosnowflake.GetLogger().GetLogLevel())
					assert.Equal(t, "oauth_client_id", config.OauthClientID)
					assert.Equal(t, "oauth_client_secret", config.OauthClientSecret)
					assert.Equal(t, "oauth_token_request_url", config.OauthTokenRequestURL)
					assert.Equal(t, "oauth_authorization_url", config.OauthAuthorizationURL)
					assert.Equal(t, "oauth_redirect_uri", config.OauthRedirectURI)
					assert.Equal(t, "oauth_scope", config.OauthScope)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.EnableSingleUseRefreshTokens)
					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_tfConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.FullInvalidTomlConfigForServiceUser(t, profile)
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.OrganizationName, "invalid")
					t.Setenv(snowflakeenvs.AccountName, "invalid")
					t.Setenv(snowflakeenvs.User, "invalid")
					t.Setenv(snowflakeenvs.PrivateKey, "invalid")
					t.Setenv(snowflakeenvs.Warehouse, "invalid")
					t.Setenv(snowflakeenvs.Protocol, "invalid")
					t.Setenv(snowflakeenvs.Port, "-1")
					t.Setenv(snowflakeenvs.Token, "invalid")
					t.Setenv(snowflakeenvs.Role, "invalid")
					t.Setenv(snowflakeenvs.ValidateDefaultParameters, "false")
					t.Setenv(snowflakeenvs.ClientIp, "2.2.2.2")
					t.Setenv(snowflakeenvs.Host, "")
					t.Setenv(snowflakeenvs.Authenticator, "invalid")
					t.Setenv(snowflakeenvs.Passcode, "")
					t.Setenv(snowflakeenvs.PasscodeInPassword, "false")
					t.Setenv(snowflakeenvs.OktaUrl, testvars.ExampleOktaUrlFromEnvString)
					t.Setenv(snowflakeenvs.LoginTimeout, "100")
					t.Setenv(snowflakeenvs.RequestTimeout, "200")
					t.Setenv(snowflakeenvs.JwtExpireTimeout, "300")
					t.Setenv(snowflakeenvs.ClientTimeout, "400")
					t.Setenv(snowflakeenvs.JwtClientTimeout, "500")
					t.Setenv(snowflakeenvs.ExternalBrowserTimeout, "600")
					t.Setenv(snowflakeenvs.InsecureMode, "false")
					t.Setenv(snowflakeenvs.OcspFailOpen, "false")
					t.Setenv(snowflakeenvs.KeepSessionAlive, "false")
					t.Setenv(snowflakeenvs.DisableTelemetry, "false")
					t.Setenv(snowflakeenvs.ClientRequestMfaToken, "false")
					t.Setenv(snowflakeenvs.ClientStoreTemporaryCredential, "false")
					t.Setenv(snowflakeenvs.DisableQueryContextCache, "false")
					t.Setenv(snowflakeenvs.IncludeRetryReason, "false")
					t.Setenv(snowflakeenvs.MaxRetryCount, "2")
					t.Setenv(snowflakeenvs.DriverTracing, "invalid")
					t.Setenv(snowflakeenvs.TmpDirectoryPath, "../")
					t.Setenv(snowflakeenvs.DisableConsoleLogin, "false")
					t.Setenv(snowflakeenvs.OauthClientId, "oauth_client_id")
					t.Setenv(snowflakeenvs.OauthClientSecret, "oauth_client_secret")
					t.Setenv(snowflakeenvs.OauthTokenRequestUrl, "oauth_token_request_url")
					t.Setenv(snowflakeenvs.OauthAuthorizationUrl, "oauth_authorization_url")
					t.Setenv(snowflakeenvs.OauthRedirectUri, "oauth_redirect_uri")
					t.Setenv(snowflakeenvs.OauthScope, "oauth_scope")
					t.Setenv(snowflakeenvs.EnableSingleUseRefreshTokens, "true")
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().AllFields(tmpServiceUserConfig, tmpServiceUser), datasourceModel()),
				Check: func(s *terraform.State) error {
					config := TestAccProvider.Meta().(*internalprovider.Context).Client.GetConfig()

					assert.Equal(t, tmpServiceUser.OrgAndAccount(), config.Account)
					assert.Equal(t, tmpServiceUser.UserId.Name(), config.User)
					assert.Equal(t, tmpServiceUser.WarehouseId.Name(), config.Warehouse)
					assert.Equal(t, tmpServiceUser.RoleId.Name(), config.Role)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
					assert.Equal(t, net.ParseIP("3.3.3.3"), config.ClientIP)
					assert.Equal(t, "https", config.Protocol)
					assert.Equal(t, fmt.Sprintf("%s.snowflakecomputing.com", tmpServiceUser.OrgAndAccount()), config.Host)
					assert.Equal(t, 443, config.Port)
					assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
					assert.False(t, config.PasscodeInPassword)
					assert.Equal(t, testvars.ExampleOktaUrl, config.OktaURL)
					assert.Equal(t, 101*time.Second, config.LoginTimeout)
					assert.Equal(t, 201*time.Second, config.RequestTimeout)
					assert.Equal(t, 301*time.Second, config.JWTExpireTimeout)
					assert.Equal(t, 401*time.Second, config.ClientTimeout)
					assert.Equal(t, 501*time.Second, config.JWTClientTimeout)
					assert.Equal(t, 601*time.Second, config.ExternalBrowserTimeout)
					assert.Equal(t, 3, config.MaxRetryCount)
					assert.Equal(t, "terraform-provider-snowflake", config.Application)
					assert.True(t, config.InsecureMode) //nolint:staticcheck
					assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
					assert.Equal(t, "correct token", config.Token)
					assert.True(t, config.KeepSessionAlive)
					assert.True(t, config.DisableTelemetry)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), config.Tracing)
					assert.Equal(t, "../../", config.TmpDirPath)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
					assert.True(t, config.DisableQueryContextCache)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
					assert.Equal(t, map[string]*string{
						"foo": sdk.Pointer("piyo"),
					}, config.Params)
					assert.Equal(t, string(sdk.DriverLogLevelWarning), gosnowflake.GetLogger().GetLogLevel())
					assert.Equal(t, "oauth_client_id", config.OauthClientID)
					assert.Equal(t, "oauth_client_secret", config.OauthClientSecret)
					assert.Equal(t, "oauth_token_request_url", config.OauthTokenRequestURL)
					assert.Equal(t, "oauth_authorization_url", config.OauthAuthorizationURL)
					assert.Equal(t, "oauth_redirect_uri", config.OauthRedirectURI)
					assert.Equal(t, "oauth_scope", config.OauthScope)
					assert.Equal(t, gosnowflake.ConfigBoolTrue, config.EnableSingleUseRefreshTokens)
					return nil
				},
			},
		},
	})
}

func TestAcc_Provider_useNonExistentDefaultParams(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	nonExisting := "NON-EXISTENT"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.ConfigPath)

			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithRole(nonExisting).WithValidateDefaultParameters("true"), datasourceModel()),
				ExpectError: regexp.MustCompile("Role 'NON-EXISTENT' specified in the connect string does not exist or not authorized."),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithWarehouse(nonExisting).WithValidateDefaultParameters("true"), datasourceModel()),
				ExpectError: regexp.MustCompile("The requested warehouse does not exist or not authorized."),
			},
			// check that using a non-existing warehouse with disabled verification succeeds
			{
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithWarehouse(nonExisting).WithValidateDefaultParameters("false"), datasourceModel()),
			},
		},
	})
}

// prove we can use tri-value booleans, similarly to the ones in resources
func TestAcc_Provider_triValueBoolean(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.97.0"),
				Config:            config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Default).WithClientStoreTemporaryCredentialBool(true), datasourceModel()),
			},
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithClientStoreTemporaryCredentialBool(true), datasourceModel()),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithClientStoreTemporaryCredential("true"), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_triValueBooleanTransitions(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().StoreTempTomlConfig(t, func(profile string) string {
		return helpers.TomlConfigForServiceUserWithModifiers(t, profile, tmpServiceUser, func(cfg *sdk.ConfigDTO) *sdk.ConfigDTO {
			return cfg.
				WithClientRequestMfaToken(false).
				// WithDisableConsoleLogin(). // omitted - default
				// WithIncludeRetryReason(). // omitted - default
				WithClientStoreTemporaryCredential(true)
		})
	})
	providerConfigModel := providermodel.SnowflakeProvider().
		WithProfile(tmpServiceUserConfig.Profile).
		// WithClientRequestMfaTokenBool(). // omitted - default
		// WithDisableConsoleLoginBool(). // omitted - default
		WithIncludeRetryReasonBool(true).
		WithClientStoreTemporaryCredentialBool(false)

	printConfigBool := func(cb gosnowflake.ConfigBool) string {
		var s string
		switch cb {
		case gosnowflake.ConfigBoolTrue:
			s = "ConfigBoolTrue"
		case gosnowflake.ConfigBoolFalse:
			s = "ConfigBoolFalse"
		default:
			s = "ConfigBoolDefault"
		}
		return s
	}
	verifyDriverConfig := func(t *testing.T) func(*terraform.State) error {
		t.Helper()

		validateConfigBoolField := func(fieldName string, fieldValue gosnowflake.ConfigBool, expectedValue gosnowflake.ConfigBool) error {
			if fieldValue != expectedValue {
				return fmt.Errorf("expected %s: %s; got: %s", fieldName, printConfigBool(expectedValue), printConfigBool(fieldValue))
			}
			return nil
		}

		return func(state *terraform.State) error {
			cfg := lastConfiguredProviderContext.Client.GetConfig()
			return errors.Join(
				validateConfigBoolField("ClientRequestMfaToken", cfg.ClientRequestMfaToken, gosnowflake.ConfigBoolFalse),
				validateConfigBoolField("DisableConsoleLogin", cfg.DisableConsoleLogin, sdk.GosnowflakeBoolConfigDefault),
				validateConfigBoolField("IncludeRetryReason", cfg.IncludeRetryReason, gosnowflake.ConfigBoolTrue),
				validateConfigBoolField("ClientStoreTemporaryCredential", cfg.ClientStoreTemporaryCredential, gosnowflake.ConfigBoolFalse),
			)
		}
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, providerConfigModel, datasourceModel()),
				Check:                    verifyDriverConfig(t),
			},
		},
	})
}

func TestAcc_Provider_sessionParameters(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				// TODO [SNOW-1348325]: Use parameter data source with `IN SESSION` filtering.
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithParamsValue(
					tfconfig.ObjectVariable(
						map[string]tfconfig.Variable{
							"statement_timeout_in_seconds": tfconfig.IntegerVariable(31337),
						},
					),
				)) + executeShowSessionParameter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_execute.t", "query_results.#", "1"),
					resource.TestCheckResourceAttr("snowflake_execute.t", "query_results.0.value", "31337"),
				),
			},
		},
	})
}

func TestAcc_Provider_JwtAuth(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)
	tmpIncorrectServiceUserConfig := testClient().TempIncorrectTomlConfigForServiceUser(t, tmpServiceUser)
	tmpServiceUserWithEncryptedKeyConfig := testClient().TempTomlConfigForServiceUserWithEncryptedKey(t, tmpServiceUser)
	tmpIncorrectServiceUserWithEncryptedKeyConfig := testClient().TempIncorrectTomlConfigForServiceUserWithEncryptedKey(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// authenticate with incorrect private key
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpIncorrectServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpIncorrectServiceUserConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				ExpectError: regexp.MustCompile("JWT token is invalid"),
			},
			// authenticate with unencrypted private key
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
			},
			// check encrypted private key with incorrect password
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpIncorrectServiceUserWithEncryptedKeyConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpIncorrectServiceUserWithEncryptedKeyConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
				ExpectError: regexp.MustCompile("pkcs8: incorrect password"),
			},
			// authenticate with encrypted private key
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserWithEncryptedKeyConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserWithEncryptedKeyConfig.Profile).WithAuthenticatorType(sdk.AuthenticationTypeJwt), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_SnowflakeAuth(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpLegacyServiceUser := testClient().SetUpTemporaryLegacyServiceUser(t)
	tmpLegacyServiceUserConfig := testClient().TempTomlConfigForLegacyServiceUser(t, tmpLegacyServiceUser)
	tmpLegacyServiceUserConfigWithoutAuthenticator := testClient().TempTomlConfigForServiceUserWithoutAuthenticator(t, tmpLegacyServiceUser)
	incorrectLegacyServiceUserConfig := testClient().TempIncorrectTomlConfigForLegacyServiceUser(t, tmpLegacyServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, incorrectLegacyServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(incorrectLegacyServiceUserConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("Incorrect username or password was specified"),
			},
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpLegacyServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpLegacyServiceUserConfig.Profile), datasourceModel()),
			},
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpLegacyServiceUserConfigWithoutAuthenticator.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpLegacyServiceUserConfigWithoutAuthenticator.Profile), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_ProgrammaticAccessTokenAuth(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	userWithPat := testClient().SetUpTemporaryLegacyServiceUserWithPat(t)
	tmpIncorrectServiceUserConfig := testClient().TempIncorrectTomlConfigForTmpUser(t, &userWithPat.TmpUser)
	userWithPatConfig := testClient().TempTomlConfigForServiceUserWithPat(t, userWithPat)
	userWithPatTfConfig := providermodel.PatConfig(*userWithPat)

	userWithPatAsPasswordConfig := testClient().TempTomlConfigForServiceUserWithPatAsPassword(t, userWithPat)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Authenticate with PROGRAMMATIC_ACCESS_TOKEN authenticator and PAT passed in token field
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, userWithPatConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(userWithPatConfig.Profile), datasourceModel()),
			},
			// Authenticate with PROGRAMMATIC_ACCESS_TOKEN authenticator and PAT passed in token field (set in Terraform)
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpIncorrectServiceUserConfig.Path)
				},
				Config: config.FromModels(t, userWithPatTfConfig, datasourceModel()),
			},
			// Authenticate with the default authenticator and PAT passed in password field
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, userWithPatAsPasswordConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(userWithPatAsPasswordConfig.Profile), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_ProgrammaticAccessTokenAuth_tokenInTfConfig(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	userWithPat := testClient().SetUpTemporaryLegacyServiceUserWithPat(t)
	userWithPatConfig := testClient().TempTomlConfigForServiceUserWithPat(t, userWithPat)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Setting the proper PAT token in tf config results in error due to incorrect authenticator setting
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.4.0"),
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, userWithPatConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(userWithPatConfig.Profile).WithToken(userWithPat.Pat), datasourceModel()),
				ExpectError: regexp.MustCompile("Invalid OAuth access token"),
			},
			// Behavior is fixed in the new version
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, userWithPatConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(userWithPatConfig.Profile).WithToken(userWithPat.Pat), datasourceModel()),
			},
		},
	})
}

func TestAcc_Provider_invalidConfigurations(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() {
			t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithClientIp("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("expected client_ip to contain a valid IP"),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithProtocol("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("invalid protocol: invalid"),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithPort(123456789), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected "port" to be a valid port number or 0, got: 123456789`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithAuthenticator("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile("invalid authenticator type: invalid"),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithOktaUrl("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected "okta_url" to have a host, got invalid`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithLoginTimeout(-1), datasourceModel()),
				ExpectError: regexp.MustCompile(`expected login_timeout to be at least \(0\), got -1`),
			},
			{
				Config: config.FromModels(
					t,
					providermodel.SnowflakeProvider().
						WithProfile(tmpServiceUserConfig.Profile).
						WithTokenAccessorValue(
							tfconfig.ObjectVariable(
								map[string]tfconfig.Variable{
									"token_endpoint": tfconfig.StringVariable("invalid"),
									"refresh_token":  tfconfig.StringVariable("refresh_token"),
									"client_id":      tfconfig.StringVariable("client_id"),
									"client_secret":  tfconfig.StringVariable("client_secret"),
									"redirect_uri":   tfconfig.StringVariable("redirect_uri"),
								},
							),
						),
					datasourceModel(),
				),
				ExpectError: regexp.MustCompile(`expected "token_endpoint" to have a host, got invalid`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithDriverTracing("invalid"), datasourceModel()),
				ExpectError: regexp.MustCompile(`invalid driver log level: invalid`),
			},
			{
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile("non-existing"), datasourceModel()),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`profile "non-existing" not found in file %s`, tmpServiceUserConfig.Path)),
			},
			{
				Config:      providerConfigWithDatasourcePreviewFeatureEnabled(testprofiles.Default, "snowflake_invalid_feature"),
				ExpectError: regexp.MustCompile(`expected .* preview_features_enabled.* to be one of((.|\n)*), got snowflake_invalid_feature`),
			},
		},
	})
}

func TestAcc_Provider_PreviewFeaturesEnabled(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	t.Setenv(string(testenvs.EnableAllPreviewFeatures), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile).WithPreviewFeaturesEnabled(string(previewfeatures.DatabaseDatasource)), datasourceModel()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceModel().DatasourceReference(), "name"),
				),
			},
		},
	})
}

func TestAcc_Provider_PreviewFeaturesDisabled(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	t.Setenv(string(testenvs.EnableAllPreviewFeatures), "")

	tmpServiceUser := testClient().SetUpTemporaryServiceUser(t)
	tmpServiceUserConfig := testClient().TempTomlConfigForServiceUser(t, tmpServiceUser)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, tmpServiceUserConfig.Path)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(tmpServiceUserConfig.Profile), datasourceModel()),
				ExpectError: regexp.MustCompile("snowflake_database_datasource is currently a preview feature, and must be enabled by adding snowflake_database_datasource to `preview_features_enabled` in Terraform configuration"),
			},
		},
	})
}

func providerConfigWithDatasourcePreviewFeatureEnabled(profile, feature string) string {
	return fmt.Sprintf(`
provider "snowflake" {
	profile = "%[1]s"
	preview_features_enabled = ["%[2]s_datasource"]
}
data %[2]s t {}
`, profile, feature)
}

func datasourceModel() config.DatasourceModel {
	return datasourcemodel.Database("t", TestDatabaseName)
}

func executeShowSessionParameter() string {
	return `
resource snowflake_execute "t" {
    execute = "SELECT 1"
    query = "SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN SESSION"
    revert        = "SELECT 1"
}`
}
