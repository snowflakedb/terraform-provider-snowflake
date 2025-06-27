package sdk

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1827309]: use toml config builder instead of hardcoding
func TestLoadConfigFileLegacy(t *testing.T) {
	c := `
	[default]
	accountname='TEST_ACCOUNT'
	organizationname='TEST_ORG'
	user='TEST_USER'
	password='abcd1234'
	role='ACCOUNTADMIN'

	[securityadmin]
	accountname='TEST_ACCOUNT'
	organizationname='TEST_ORG'
	user='TEST_USER'
	password='abcd1234'
	role='SECURITYADMIN'
	`
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	m, err := LoadConfigFile[*LegacyConfigDTO](configPath, true)
	require.NoError(t, err)
	assert.Equal(t, "TEST_ACCOUNT", *m["default"].AccountName)
	assert.Equal(t, "TEST_ORG", *m["default"].OrganizationName)
	assert.Equal(t, "TEST_USER", *m["default"].User)
	assert.Equal(t, "abcd1234", *m["default"].Password)
	assert.Equal(t, "ACCOUNTADMIN", *m["default"].Role)
	assert.Equal(t, "TEST_ACCOUNT", *m["securityadmin"].AccountName)
	assert.Equal(t, "TEST_ORG", *m["securityadmin"].OrganizationName)
	assert.Equal(t, "TEST_USER", *m["securityadmin"].User)
	assert.Equal(t, "abcd1234", *m["securityadmin"].Password)
	assert.Equal(t, "SECURITYADMIN", *m["securityadmin"].Role)
}

func TestLoadConfigFileWithUnknownFieldsLegacy(t *testing.T) {
	c := `
	[default]
	unknown='TEST_ACCOUNT'
	accountname='TEST_ACCOUNT'
	`
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	m, err := LoadConfigFile[*LegacyConfigDTO](configPath, true)
	require.NoError(t, err)
	assert.Equal(t, map[string]*LegacyConfigDTO{
		"default": {
			AccountName: Pointer("TEST_ACCOUNT"),
		},
	}, m)
}

func TestLoadConfigFileWithInvalidFieldTypeFailsLegacy(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		wantType  string
	}{
		{name: "AccountName", fieldName: "accountname", wantType: "*string"},
		{name: "OrganizationName", fieldName: "organizationname", wantType: "*string"},
		{name: "User", fieldName: "user", wantType: "*string"},
		{name: "Username", fieldName: "username", wantType: "*string"},
		{name: "Password", fieldName: "password", wantType: "*string"},
		{name: "Host", fieldName: "host", wantType: "*string"},
		{name: "Warehouse", fieldName: "warehouse", wantType: "*string"},
		{name: "Role", fieldName: "role", wantType: "*string"},
		{name: "Params", fieldName: "params", wantType: "*map[string]*string"},
		{name: "ClientIp", fieldName: "clientip", wantType: "*string"},
		{name: "Protocol", fieldName: "protocol", wantType: "*string"},
		{name: "Passcode", fieldName: "passcode", wantType: "*string"},
		{name: "PasscodeInPassword", fieldName: "passcodeinpassword", wantType: "*bool"},
		{name: "OktaUrl", fieldName: "oktaurl", wantType: "*string"},
		{name: "Authenticator", fieldName: "authenticator", wantType: "*string"},
		{name: "InsecureMode", fieldName: "insecuremode", wantType: "*bool"},
		{name: "OcspFailOpen", fieldName: "ocspfailopen", wantType: "*bool"},
		{name: "Token", fieldName: "token", wantType: "*string"},
		{name: "KeepSessionAlive", fieldName: "keepsessionalive", wantType: "*bool"},
		{name: "PrivateKey", fieldName: "privatekey", wantType: "*string"},
		{name: "PrivateKeyPassphrase", fieldName: "privatekeypassphrase", wantType: "*string"},
		{name: "DisableTelemetry", fieldName: "disabletelemetry", wantType: "*bool"},
		{name: "ValidateDefaultParameters", fieldName: "validatedefaultparameters", wantType: "*bool"},
		{name: "ClientRequestMfaToken", fieldName: "clientrequestmfatoken", wantType: "*bool"},
		{name: "ClientStoreTemporaryCredential", fieldName: "clientstoretemporarycredential", wantType: "*bool"},
		{name: "DriverTracing", fieldName: "tracing", wantType: "*string"},
		{name: "TmpDirPath", fieldName: "tmpdirpath", wantType: "*string"},
		{name: "DisableQueryContextCache", fieldName: "disablequerycontextcache", wantType: "*bool"},
		{name: "IncludeRetryReason", fieldName: "includeretryreason", wantType: "*bool"},
		{name: "DisableConsoleLogin", fieldName: "disableconsolelogin", wantType: "*bool"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s has to have a correct type", tt.name), func(t *testing.T) {
			config := fmt.Sprintf(`
		[default]
		%s=42
		`, tt.fieldName)
			configPath := testhelpers.TestFile(t, "config", []byte(config))

			_, err := LoadConfigFile[*LegacyConfigDTO](configPath, true)
			require.ErrorContains(t, err, fmt.Sprintf("toml: cannot decode TOML integer into struct field sdk.LegacyConfigDTO.%s of type %s", tt.name, tt.wantType))
		})
	}
}

func TestLoadConfigFileWithInvalidFieldTypeIntFailsLegacy(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
	}{
		{name: "Port", fieldName: "port"},
		{name: "ClientTimeout", fieldName: "clienttimeout"},
		{name: "JwtClientTimeout", fieldName: "jwtclienttimeout"},
		{name: "LoginTimeout", fieldName: "logintimeout"},
		{name: "RequestTimeout", fieldName: "requesttimeout"},
		{name: "JwtExpireTimeout", fieldName: "jwtexpiretimeout"},
		{name: "ExternalBrowserTimeout", fieldName: "externalbrowsertimeout"},
		{name: "MaxRetryCount", fieldName: "maxretrycount"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s has to have a correct type", tt.name), func(t *testing.T) {
			config := fmt.Sprintf(`
		[default]
		%s=value
		`, tt.fieldName)
			configPath := testhelpers.TestFile(t, "config", []byte(config))

			_, err := LoadConfigFile[*LegacyConfigDTO](configPath, true)
			require.ErrorContains(t, err, "toml: incomplete number")
		})
	}
}

func TestLoadConfigFileWithInvalidTOMLFailsLegacy(t *testing.T) {
	tests := []struct {
		name   string
		config string
		err    string
	}{
		{
			name: "key without a value",
			config: `
			[default]
			password="sensitive"
			accountname=
			`,
			err: "toml: incomplete number",
		},
		{
			name: "value without a key",
			config: `
			[default]
			password="sensitive"
			="value"
			`,
			err: "toml: invalid character at start of key: =",
		},
		{
			name: "multiple profiles with the same name",
			config: `
			[default]
			password="sensitive"
			accountname="value"
			[default]
			organizationname="value"
			`,
			err: "toml: table default already exists",
		},
		{
			name: "multiple keys with the same name",
			config: `
			[default]
			password="sensitive"
			accountname="foo"
			accountname="bar"
			`,
			err: "toml: key accountname is already defined",
		},
		{
			name: "more than one key in a line",
			config: `
			[default]
			password="sensitive"
			accountname="account" organizationname="organizationname"
			`,
			err: "toml: expected newline but got U+006F 'o'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := testhelpers.TestFile(t, "config", []byte(tt.config))

			_, err := LoadConfigFile[*LegacyConfigDTO](configPath, true)
			require.ErrorContains(t, err, tt.err)
			require.NotContains(t, err.Error(), "sensitive")
		})
	}
}

func TestProfileConfigLegacy(t *testing.T) {
	unencryptedKey, encryptedKey := random.GenerateRSAPrivateKeyEncrypted(t, "password")

	c := fmt.Sprintf(`
	[securityadmin]
	account='account'
	accountname='accountname'
	organizationname='organizationname'
	user='user'
	password='password'
	host='host'
	warehouse='warehouse'
	role='role'
	clientip='1.1.1.1'
	protocol='http'
	passcode='passcode'
	port=1
	passcodeinpassword=true
	oktaurl='%[3]s'
	clienttimeout=10
	jwtclienttimeout=20
	logintimeout=30
	requesttimeout=40
	jwtexpiretimeout=50
	externalbrowsertimeout=60
	maxretrycount=1
	authenticator='SNOWFLAKE_JWT'
	insecuremode=true
	ocspfailopen=true
	token='token'
	keepsessionalive=true
	privatekey="""%[1]s"""
	privatekeypassphrase='%[2]s'
	disabletelemetry=true
	validatedefaultparameters=true
	clientrequestmfatoken=true
	clientstoretemporarycredential=true
	tracing='tracing'
	tmpdirpath='.'
	disablequerycontextcache=true
	includeretryreason=true
	disableconsolelogin=true

	[securityadmin.params]
	foo = 'bar'
	`, encryptedKey, "password", testvars.ExampleOktaUrlString)
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	t.Run("with found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("securityadmin", WithUseLegacyTomlFormat(true))
		require.NoError(t, err)
		require.NotNil(t, config.PrivateKey)

		gotKey, err := x509.MarshalPKCS8PrivateKey(config.PrivateKey)
		require.NoError(t, err)
		gotUnencryptedKey := pem.EncodeToMemory(
			&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: gotKey,
			},
		)

		assert.Equal(t, "organizationname-accountname", config.Account)
		assert.Equal(t, "user", config.User)
		assert.Equal(t, "password", config.Password)
		assert.Equal(t, "warehouse", config.Warehouse)
		assert.Equal(t, "role", config.Role)
		assert.Equal(t, map[string]*string{"foo": Pointer("bar")}, config.Params)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
		assert.Equal(t, "1.1.1.1", config.ClientIP.String())
		assert.Equal(t, "http", config.Protocol)
		assert.Equal(t, "host", config.Host)
		assert.Equal(t, 1, config.Port)
		assert.Equal(t, gosnowflake.AuthTypeJwt, config.Authenticator)
		assert.Equal(t, "passcode", config.Passcode)
		assert.True(t, config.PasscodeInPassword)
		assert.Equal(t, testvars.ExampleOktaUrlString, config.OktaURL.String())
		assert.Equal(t, 10*time.Second, config.ClientTimeout)
		assert.Equal(t, 20*time.Second, config.JWTClientTimeout)
		assert.Equal(t, 30*time.Second, config.LoginTimeout)
		assert.Equal(t, 40*time.Second, config.RequestTimeout)
		assert.Equal(t, 50*time.Second, config.JWTExpireTimeout)
		assert.Equal(t, 60*time.Second, config.ExternalBrowserTimeout)
		assert.Equal(t, 1, config.MaxRetryCount)
		assert.True(t, config.InsecureMode) //nolint:staticcheck
		assert.Equal(t, "token", config.Token)
		assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
		assert.True(t, config.KeepSessionAlive)
		assert.Equal(t, unencryptedKey, string(gotUnencryptedKey))
		assert.True(t, config.DisableTelemetry)
		assert.Equal(t, "tracing", config.Tracing)
		assert.Equal(t, ".", config.TmpDirPath)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
		assert.True(t, config.DisableQueryContextCache)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
	})

	t.Run("with not found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("orgadmin", WithUseLegacyTomlFormat(true))
		require.NoError(t, err)
		require.Nil(t, config)
	})

	t.Run("with not found config", func(t *testing.T) {
		filename := random.AlphaN(8)
		t.Setenv(snowflakeenvs.ConfigPath, filename)

		config, err := ProfileConfig("orgadmin", WithUseLegacyTomlFormat(true))
		require.ErrorContains(t, err, fmt.Sprintf("could not load config file: reading information about the config file: stat %s: no such file or directory", filename))
		require.Nil(t, config)
	})
}
