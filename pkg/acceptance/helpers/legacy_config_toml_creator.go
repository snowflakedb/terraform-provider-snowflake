package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// FullLegacyTomlConfigForServiceUser is a temporary function used to test provider configuration
// TODO [SNOW-1827309]: use toml marshaling from "github.com/pelletier/go-toml/v2"
// TODO [SNOW-1827309]: add builders for our toml config struct
func FullLegacyTomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
privatekey = '''%[7]s'''
role = '%[3]s'
organizationname = '%[5]s'
accountname = '%[6]s'
warehouse = '%[4]s'
clientip = '1.2.3.4'
protocol = 'https'
port = 443
oktaurl = '%[8]s'
clienttimeout = 10
jwtclienttimeout = 20
logintimeout = 30
requesttimeout = 40
jwtexpiretimeout = 50
externalbrowsertimeout = 60
maxretrycount = 1
authenticator = 'SNOWFLAKE_JWT'
insecuremode = true
ocspfailopen = true
token = 'token'
keepsessionalive = true
disabletelemetry = true
validatedefaultparameters = true
clientrequestmfatoken = true
clientstoretemporarycredential = true
tracing = 'warning'
tmpdirpath = '.'
disablequerycontextcache = true
includeretryreason = true
disableconsolelogin = true

[%[1]s.params]
foo = 'bar'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey, testvars.ExampleOktaUrlString)
}

// FullInvalidLegacyTomlConfigForServiceUser is a temporary function used to test provider configuration
func FullInvalidLegacyTomlConfigForServiceUser(t *testing.T, profile string) string {
	t.Helper()

	privateKey, _, _, _ := random.GenerateRSAKeyPair(t, "")
	return fmt.Sprintf(`
[%[1]s]
user = 'invalid'
privatekey = '''%[2]s'''
role = 'invalid'
accountname = 'invalid'
organizationname = 'invalid'
warehouse = 'invalid'
clientip = 'invalid'
protocol = 'invalid'
port = -1
oktaurl = 'invalid'
clienttimeout = -1
jwtclienttimeout = -1
logintimeout = -1
requesttimeout = -1
jwtexpiretimeout = -1
externalbrowsertimeout = -1
maxretrycount = -1
authenticator = 'snowflake'
insecuremode = true
ocspfailopen = true
token = 'token'
keepsessionalive = true
disabletelemetry = true
validatedefaultparameters = false
clientrequestmfatoken = true
clientstoretemporarycredential = true
tracing = 'invalid'
tmpdirpath = '.'
disablequerycontextcache = true
includeretryreason = true
disableconsolelogin = true

[%[1]s.params]
foo = 'bar'`, profile, privateKey)
}

// LegacyTomlConfigForServiceUser is a temporary function used to test provider configuration
func LegacyTomlConfigForServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
privatekey = '''%[7]s'''
role = '%[3]s'
organizationname = '%[5]s'
accountname = '%[6]s'
warehouse = '%[4]s'
authenticator = 'SNOWFLAKE_JWT'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey)
}

// LegacyTomlConfigForServiceUserWithEncryptedKey is a temporary function used to test provider configuration
func LegacyTomlConfigForServiceUserWithEncryptedKey(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, privateKey string, pass string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
privatekey = '''%[7]s'''
privatekeypassphrase = '%[8]s'
role = '%[3]s'
organizationname = '%[5]s'
accountname = '%[6]s'
warehouse = '%[4]s'
authenticator = 'SNOWFLAKE_JWT'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey, pass)
}

// LegacyTomlIncorrectConfigForServiceUser is a temporary function used to test provider configuration
func LegacyTomlIncorrectConfigForServiceUser(t *testing.T, profile string, accountIdentifier sdk.AccountIdentifier) string {
	t.Helper()

	privateKey, _, _, _ := random.GenerateRSAKeyPair(t, "")
	return fmt.Sprintf(`
[%[1]s]
user = 'non-existing-user'
privatekey = '''%[4]s'''
role = 'non-existing-role'
organizationname = '%[2]s'
accountname = '%[3]s'
authenticator = 'SNOWFLAKE_JWT'
`, profile, accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), privateKey)
}

// LegacyTomlConfigForLegacyServiceUser is a temporary function used to test provider configuration
func LegacyTomlConfigForLegacyServiceUser(t *testing.T, profile string, userId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, accountIdentifier sdk.AccountIdentifier, pass string) string {
	t.Helper()

	return fmt.Sprintf(`
[%[1]s]
user = '%[2]s'
password = '%[7]s'
role = '%[3]s'
organizationname = '%[5]s'
accountname = '%[6]s'
warehouse = '%[4]s'
authenticator = 'SNOWFLAKE'
`, profile, userId.Name(), roleId.Name(), warehouseId.Name(), accountIdentifier.OrganizationName(), accountIdentifier.AccountName(), pass)
}
