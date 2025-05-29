package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// FullTomlConfigForServiceUser is a temporary function used to test provider configuration
func FullTomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
private_key = '''%[7]s'''
role = '%[3]s'
organization_name = '%[5]s'
account_name = '%[6]s'
warehouse = '%[4]s'
client_ip = '1.2.3.4'
protocol = 'https'
port = 443
okta_url = '%[8]s'
client_timeout = 10
jwt_client_timeout = 20
login_timeout = 30
request_timeout = 40
jwt_expire_timeout = 50
external_browser_timeout = 60
max_retry_count = 1
authenticator = 'SNOWFLAKE_JWT'
insecure_mode = true
ocsp_fail_open = true
token = 'token'
keep_session_alive = true
disable_telemetry = true
validate_default_parameters = true
client_request_mfa_token = true
client_store_temporary_credential = true
driver_tracing = 'warning'
tmp_dir_path = '.'
disable_query_context_cache = true
include_retry_reason = true
disable_console_login = true

[%[1]s.params]
foo = 'bar'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey, testvars.ExampleOktaUrlString)
}

// FullInvalidTomlConfigForServiceUser is a temporary function used to test provider configuration
func FullInvalidTomlConfigForServiceUser(t *testing.T, profile string) string {
	t.Helper()

	privateKey, _, _, _ := random.GenerateRSAKeyPair(t, "")
	return fmt.Sprintf(`
[%[1]s]
user = 'invalid'
private_key = '''%[2]s'''
role = 'invalid'
account_name = 'invalid'
organization_name = 'invalid'
warehouse = 'invalid'
client_ip = 'invalid'
protocol = 'invalid'
port = -1
okta_url = 'invalid'
client_timeout = -1
jwt_client_timeout = -1
login_timeout = -1
request_timeout = -1
jwt_expire_timeout = -1
external_browser_timeout = -1
max_retry_count = -1
authenticator = 'snowflake'
insecure_mode = true
ocsp_fail_open = true
token = 'token'
keep_session_alive = true
disable_telemetry = true
validate_default_parameters = false
client_request_mfa_token = true
client_store_temporary_credential = true
tracing = 'invalid'
tmp_dir_path = '.'
disable_query_context_cache = true
include_retry_reason = true
disable_console_login = true

[%[1]s.params]
foo = 'bar'`, profile, privateKey)
}

// TomlConfigForServiceUser is a temporary function used to test provider configuration
func TomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
private_key = '''%[7]s'''
role = '%[3]s'
organization_name = '%[5]s'
account_name = '%[6]s'
warehouse = '%[4]s'
authenticator = 'SNOWFLAKE_JWT'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey)
}

// TomlConfigForServiceUserWithEncryptedKey is a temporary function used to test provider configuration
func TomlConfigForServiceUserWithEncryptedKey(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string, pass string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
private_key = '''%[7]s'''
private_key_passphrase = '%[8]s'
role = '%[3]s'
organization_name = '%[5]s'
account_name = '%[6]s'
warehouse = '%[4]s'
authenticator = 'SNOWFLAKE_JWT'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey, pass)
}

// TomlIncorrectConfigForServiceUser is a temporary function used to test provider configuration
func TomlIncorrectConfigForServiceUser(t *testing.T, profile string, accountIdentifier sdk.AccountIdentifier) string {
	t.Helper()

	privateKey, _, _, _ := random.GenerateRSAKeyPair(t, "")
	return fmt.Sprintf(`
[%[1]s]
user = 'non-existing-user'
private_key = '''%[4]s'''
role = 'non-existing-role'
organization_name = '%[2]s'
account_name = '%[3]s'
authenticator = 'SNOWFLAKE_JWT'
`, profile, accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey)
}

// TomlConfigForLegacyServiceUser is a temporary function used to test provider configuration
func TomlConfigForLegacyServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, pass string) string {
	t.Helper()

	config := sdk.NewConfigDTO().
		WithUser(userId.Name()).
		WithPassword(pass).
		WithRole(roleId.Name()).
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithWarehouse(warehouseId.Name()).
		WithAuthenticator("SNOWFLAKE")
	cfg := sdk.NewConfigFile().WithProfiles(map[string]sdk.ConfigDTO{profile: *config})
	bytes, err := cfg.Marshal()
	require.NoError(t, err)
	return string(bytes)
}
