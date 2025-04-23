package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
