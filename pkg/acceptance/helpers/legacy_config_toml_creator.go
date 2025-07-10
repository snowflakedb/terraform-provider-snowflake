package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// FullLegacyTomlConfigForServiceUser is a temporary function used to test provider configuration
func FullLegacyTomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	cfg := sdk.NewLegacyConfigFile().WithProfiles(map[string]sdk.LegacyConfigDTO{
		profile: *sdk.NewLegacyConfigDTO().
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
			WithParams(map[string]*string{"foo": sdk.Pointer("bar")}),
	})
	bytes, err := cfg.MarshalToml()
	require.NoError(t, err)
	return string(bytes)
}
