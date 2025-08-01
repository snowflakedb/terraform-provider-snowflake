package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// FullTomlConfigForServiceUser is a temporary function used to test provider configuration
func FullTomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	return configDtoToTomlString(t, profile, sdk.NewConfigDTO().
		WithUser(userId.Name()).
		WithPrivateKey(privateKey).
		WithRole(roleId.Name()).
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithWarehouse(warehouseId.Name()).
		WithClientIp("1.2.3.4").
		WithProtocol("https").
		WithPort(443).
		WithOktaUrl(testvars.ExampleOktaUrlString).
		WithClientTimeout(10).
		WithJwtClientTimeout(20).
		WithLoginTimeout(30).
		WithRequestTimeout(40).
		WithJwtExpireTimeout(50).
		WithExternalBrowserTimeout(60).
		WithMaxRetryCount(1).
		WithAuthenticator(string(sdk.AuthenticationTypeJwt)).
		WithInsecureMode(true).
		WithOcspFailOpen(true).
		WithToken("token").
		WithKeepSessionAlive(true).
		WithDisableTelemetry(true).
		WithValidateDefaultParameters(true).
		WithClientRequestMfaToken(true).
		WithClientStoreTemporaryCredential(true).
		WithDriverTracing(string(sdk.DriverLogLevelWarning)).
		WithTmpDirPath(".").
		WithDisableQueryContextCache(true).
		WithIncludeRetryReason(true).
		WithDisableConsoleLogin(true).
		WithParams(map[string]*string{
			"foo": sdk.Pointer("bar"),
		}),
	)
}

// FullInvalidTomlConfigForServiceUser is a temporary function used to test provider configuration
func FullInvalidTomlConfigForServiceUser(t *testing.T, profile string) string {
	t.Helper()

	privateKey, _, _, _ := random.GenerateRSAKeyPair(t, "")
	dto := sdk.NewConfigDTO().
		WithUser("invalid").
		WithPrivateKey(privateKey).
		WithRole("invalid").
		WithAccountName("invalid").
		WithOrganizationName("invalid").
		WithWarehouse("invalid").
		WithClientIp("invalid").
		WithProtocol("invalid").
		WithPort(-1).
		WithOktaUrl("invalid").
		WithClientTimeout(-1).
		WithJwtClientTimeout(-1).
		WithLoginTimeout(-1).
		WithRequestTimeout(-1).
		WithJwtExpireTimeout(-1).
		WithExternalBrowserTimeout(-1).
		WithMaxRetryCount(-1).
		WithAuthenticator("snowflake").
		WithInsecureMode(true).
		WithOcspFailOpen(true).
		WithToken("token").
		WithKeepSessionAlive(true).
		WithDisableTelemetry(true).
		WithValidateDefaultParameters(false).
		WithClientRequestMfaToken(true).
		WithClientStoreTemporaryCredential(true).
		WithDriverTracing("invalid").
		WithTmpDirPath(".").
		WithDisableQueryContextCache(true).
		WithIncludeRetryReason(true).
		WithDisableConsoleLogin(true).
		WithParams(map[string]*string{
			"foo": sdk.Pointer("bar"),
		})
	return configDtoToTomlString(t, profile, dto)
}

// TomlConfigForServiceUser is a temporary function used to test provider configuration
func TomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	return configDtoToTomlString(t, profile, sdk.NewConfigDTO().
		WithUser(userId.Name()).
		WithPrivateKey(privateKey).
		WithRole(roleId.Name()).
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithWarehouse(warehouseId.Name()).
		WithAuthenticator(string(sdk.AuthenticationTypeJwt)),
	)
}

// TomlConfigForServiceUserWithPat is a temporary function used to test provider configuration
func TomlConfigForServiceUserWithPat(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, token string) string {
	t.Helper()

	return configDtoToTomlString(t, profile, sdk.NewConfigDTO().
		WithUser(userId.Name()).
		WithRole(roleId.Name()).
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithWarehouse(warehouseId.Name()).
		WithAuthenticator(string(sdk.AuthenticationTypeProgrammaticAccessToken)).
		WithToken(token),
	)
}

// TomlConfigForServiceUserWithEncryptedKey is a temporary function used to test provider configuration
func TomlConfigForServiceUserWithEncryptedKey(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string, pass string) string {
	t.Helper()

	return configDtoToTomlString(t, profile, sdk.NewConfigDTO().
		WithUser(userId.Name()).
		WithPrivateKey(privateKey).
		WithPrivateKeyPassphrase(pass).
		WithRole(roleId.Name()).
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithWarehouse(warehouseId.Name()).
		WithAuthenticator(string(sdk.AuthenticationTypeJwt)),
	)
}

// TomlIncorrectConfigForServiceUser is a temporary function used to test provider configuration
func TomlIncorrectConfigForServiceUser(t *testing.T, profile string, accountIdentifier sdk.AccountIdentifier) string {
	t.Helper()

	privateKey, _, _, _ := random.GenerateRSAKeyPair(t, "")
	return configDtoToTomlString(t, profile, sdk.NewConfigDTO().
		WithUser("non-existing-user").
		WithPrivateKey(privateKey).
		WithRole("non-existing-role").
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithAuthenticator(string(sdk.AuthenticationTypeJwt)),
	)
}

// TomlConfigForLegacyServiceUser is a temporary function used to test provider configuration
func TomlConfigForLegacyServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, pass string) string {
	t.Helper()

	return configDtoToTomlString(t, profile, sdk.NewConfigDTO().
		WithUser(userId.Name()).
		WithPassword(pass).
		WithRole(roleId.Name()).
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithWarehouse(warehouseId.Name()).
		WithAuthenticator(string(sdk.AuthenticationTypeSnowflake)),
	)
}

// TomlConfigForLegacyServiceUserWithoutAuthenticator is a temporary function used to test provider configuration
func TomlConfigForLegacyServiceUserWithoutAuthenticator(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, pass string) string {
	t.Helper()

	return configDtoToTomlString(t, profile, sdk.NewConfigDTO().
		WithUser(userId.Name()).
		WithPassword(pass).
		WithRole(roleId.Name()).
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithWarehouse(warehouseId.Name()),
	)
}

// TomlConfigForServiceUserWithModifiers is a temporary function used to test provider configuration allowing to modify the toml config
func TomlConfigForServiceUserWithModifiers(t *testing.T, profile string, serviceUser *TmpServiceUser, configDtoModifier func(cfg *sdk.ConfigDTO) *sdk.ConfigDTO) string {
	t.Helper()

	configDto := sdk.NewConfigDTO().
		WithOrganizationName(serviceUser.AccountId.OrganizationName()).
		WithAccountName(serviceUser.AccountId.AccountName()).
		WithUser(serviceUser.UserId.Name()).
		WithRole(serviceUser.RoleId.Name()).
		WithWarehouse(serviceUser.WarehouseId.Name()).
		WithPrivateKey(serviceUser.PrivateKey).
		WithAuthenticator(string(sdk.AuthenticationTypeJwt))

	return configDtoToTomlString(t, profile, configDtoModifier(configDto))
}

func configDtoToTomlString(t *testing.T, profile string, config *sdk.ConfigDTO) string {
	t.Helper()

	cfg := sdk.NewConfigFile().WithProfiles(map[string]sdk.ConfigDTO{profile: *config})
	bytes, err := cfg.MarshalToml()
	require.NoError(t, err)
	return string(bytes)
}
