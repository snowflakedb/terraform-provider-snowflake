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

func TestLoadConfigFileLegacy(t *testing.T) {
	cfg := NewLegacyConfigFile().WithProfiles(map[string]LegacyConfigDTO{
		"default": *NewLegacyConfigDTO().
			WithAccountName("TEST_ACCOUNT").
			WithOrganizationName("TEST_ORG").
			WithUser("TEST_USER").
			WithPassword("abcd1234").
			WithRole("ACCOUNTADMIN"),
		"securityadmin": *NewLegacyConfigDTO().
			WithAccountName("TEST_ACCOUNT").
			WithOrganizationName("TEST_ORG").
			WithUser("TEST_USER").
			WithPassword("abcd1234").
			WithRole("SECURITYADMIN"),
	})
	bytes, err := cfg.MarshalToml()
	require.NoError(t, err)
	configPath := testhelpers.TestFile(t, "config", bytes)

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

	cfg := NewLegacyConfigFile().WithProfiles(map[string]LegacyConfigDTO{
		"securityadmin": *NewLegacyConfigDTO().
			WithAccountName("accountname").
			WithOrganizationName("organizationname").
			WithUser("user").
			WithPassword("password").
			WithHost("host").
			WithWarehouse("warehouse").
			WithRole("role").
			WithClientIp("1.1.1.1").
			WithProtocol("http").
			WithPasscode("passcode").
			WithPort(1).
			WithPasscodeInPassword(true).
			WithOktaUrl(testvars.ExampleOktaUrlString).
			WithClientTimeout(10).
			WithJwtClientTimeout(20).
			WithLoginTimeout(30).
			WithRequestTimeout(40).
			WithJwtExpireTimeout(50).
			WithExternalBrowserTimeout(60).
			WithMaxRetryCount(1).
			WithAuthenticator("SNOWFLAKE_JWT").
			WithInsecureMode(true).
			WithOcspFailOpen(true).
			WithToken("token").
			WithKeepSessionAlive(true).
			WithPrivateKey(encryptedKey).
			WithPrivateKeyPassphrase("password").
			WithDisableTelemetry(true).
			WithValidateDefaultParameters(true).
			WithClientRequestMfaToken(true).
			WithClientStoreTemporaryCredential(true).
			WithDriverTracing("tracing").
			WithTmpDirPath(".").
			WithDisableQueryContextCache(true).
			WithIncludeRetryReason(true).
			WithDisableConsoleLogin(true).
			WithParams(map[string]*string{"foo": Pointer("bar")}),
	})
	bytes, err := cfg.MarshalToml()
	require.NoError(t, err)
	c := string(bytes)
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
		assert.Equal(t, string(DriverLogLevelTrace), config.Tracing)
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

	t.Run("with old account field", func(t *testing.T) {
		c := `
		[default]
		account='ACCOUNT'
		accountname='TEST_ACCOUNT'
		organizationname='TEST_ORG'
		`
		configPath := testhelpers.TestFile(t, "config", []byte(c))

		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("default", WithUseLegacyTomlFormat(true))
		require.NoError(t, err)
		require.NotNil(t, config)
		assert.Equal(t, "TEST_ORG-TEST_ACCOUNT", config.Account)
	})
}

func TestLegacyConfigDTODriverConfig(t *testing.T) {
	privateKey, _ := random.GenerateRSAPrivateKeyEncrypted(t, "pass")
	tests := []struct {
		name     string
		input    *LegacyConfigDTO
		expected func(t *testing.T, got gosnowflake.Config, err error)
	}{
		{
			name: "minimal config with account and org",
			input: NewLegacyConfigDTO().
				WithAccountName("acc").
				WithOrganizationName("org").
				WithUser("user").
				WithPassword("pass"),
			expected: func(t *testing.T, got gosnowflake.Config, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, "org-acc", got.Account)
				assert.Equal(t, "user", got.User)
				assert.Equal(t, "pass", got.Password)
			},
		},
		{
			name: "all fields set",
			input: NewLegacyConfigDTO().
				WithAccountName("acc").
				WithOrganizationName("org").
				WithUser("user").
				WithPassword("pass").
				WithHost("host").
				WithWarehouse("wh").
				WithRole("role").
				WithParams(map[string]*string{"foo": Pointer("bar")}).
				WithClientIp("1.2.3.4").
				WithProtocol("https").
				WithPasscode("code").
				WithPort(1234).
				WithPasscodeInPassword(true).
				WithOktaUrl("https://okta.example.com").
				WithClientTimeout(10).
				WithJwtClientTimeout(20).
				WithLoginTimeout(30).
				WithRequestTimeout(40).
				WithJwtExpireTimeout(50).
				WithExternalBrowserTimeout(60).
				WithMaxRetryCount(2).
				WithAuthenticator(string(AuthenticationTypeJwt)).
				WithInsecureMode(true).
				WithOcspFailOpen(true).
				WithToken("token").
				WithKeepSessionAlive(true).
				WithPrivateKey(privateKey).
				WithPrivateKeyPassphrase("passphrase").
				WithDisableTelemetry(true).
				WithValidateDefaultParameters(true).
				WithClientRequestMfaToken(true).
				WithClientStoreTemporaryCredential(true).
				WithDriverTracing(string(DriverLogLevelDebug)).
				WithTmpDirPath("/tmp").
				WithDisableQueryContextCache(true).
				WithIncludeRetryReason(true).
				WithDisableConsoleLogin(true),
			expected: func(t *testing.T, got gosnowflake.Config, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Equal(t, "org-acc", got.Account)
				assert.Equal(t, "user", got.User) // LegacyConfigDTO does not have Username override
				assert.Equal(t, "pass", got.Password)
				assert.Equal(t, "host", got.Host)
				assert.Equal(t, "wh", got.Warehouse)
				assert.Equal(t, "role", got.Role)
				assert.Equal(t, map[string]*string{"foo": Pointer("bar")}, got.Params)
				assert.Equal(t, "1.2.3.4", got.ClientIP.String())
				assert.Equal(t, "https", got.Protocol)
				assert.Equal(t, "code", got.Passcode)
				assert.Equal(t, 1234, got.Port)
				assert.True(t, got.PasscodeInPassword)
				assert.Equal(t, "https://okta.example.com", got.OktaURL.String())
				assert.Equal(t, 10*time.Second, got.ClientTimeout)
				assert.Equal(t, 20*time.Second, got.JWTClientTimeout)
				assert.Equal(t, 30*time.Second, got.LoginTimeout)
				assert.Equal(t, 40*time.Second, got.RequestTimeout)
				assert.Equal(t, 50*time.Second, got.JWTExpireTimeout)
				assert.Equal(t, 60*time.Second, got.ExternalBrowserTimeout)
				assert.Equal(t, 2, got.MaxRetryCount)
				assert.Equal(t, gosnowflake.AuthTypeJwt, got.Authenticator)
				assert.True(t, got.InsecureMode)
				assert.Equal(t, gosnowflake.OCSPFailOpenTrue, got.OCSPFailOpen)
				assert.Equal(t, "token", got.Token)
				assert.True(t, got.KeepSessionAlive)
				assert.True(t, got.DisableTelemetry)
				assert.Equal(t, gosnowflake.ConfigBoolTrue, got.ValidateDefaultParameters)
				assert.Equal(t, gosnowflake.ConfigBoolTrue, got.ClientRequestMfaToken)
				assert.Equal(t, gosnowflake.ConfigBoolTrue, got.ClientStoreTemporaryCredential)
				assert.Equal(t, string(DriverLogLevelDebug), got.Tracing)
				assert.Equal(t, "/tmp", got.TmpDirPath)
				assert.True(t, got.DisableQueryContextCache)
				assert.Equal(t, gosnowflake.ConfigBoolTrue, got.IncludeRetryReason)
				assert.Equal(t, gosnowflake.ConfigBoolTrue, got.DisableConsoleLogin)

				gotKey, err := x509.MarshalPKCS8PrivateKey(got.PrivateKey)
				require.NoError(t, err)
				gotUnencryptedKey := pem.EncodeToMemory(
					&pem.Block{
						Type:  "PRIVATE KEY",
						Bytes: gotKey,
					},
				)
				assert.Equal(t, privateKey, string(gotUnencryptedKey))
			},
		},
	}

	invalid := []struct {
		name  string
		input *LegacyConfigDTO
		err   error
	}{
		{
			name: "invalid okta url",
			input: NewLegacyConfigDTO().
				WithOktaUrl(":invalid:"),
			err: fmt.Errorf("parse \":invalid:\": missing protocol scheme"),
		},
		{
			name: "invalid authenticator",
			input: NewLegacyConfigDTO().
				WithAuthenticator("invalid"),
			err: fmt.Errorf("invalid authenticator type: invalid"),
		},
		{
			name: "invalid privatekey",
			input: NewLegacyConfigDTO().
				WithPrivateKey("not_a_valid_pem"),
			err: fmt.Errorf("could not parse private key, key is not in PEM format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.DriverConfig()
			tt.expected(t, got, err)
		})
	}

	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input.DriverConfig()
			require.ErrorContains(t, err, tt.err.Error())
		})
	}
}
