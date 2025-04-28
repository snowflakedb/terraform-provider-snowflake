package sdk

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/snowflakedb/gosnowflake"
	"github.com/snowflakedb/terraform-provider-snowflake/v2/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1827309]: use toml config builder instead of hardcoding
func TestLoadConfigFile(t *testing.T) {
	c := `
	[default]
	account_name='TEST_ACCOUNT'
	organization_name='TEST_ORG'
	user='TEST_USER'
	password='abcd1234'
	role='ACCOUNTADMIN'

	[securityadmin]
	account_name='TEST_ACCOUNT'
	organization_name='TEST_ORG'
	user='TEST_USER'
	password='abcd1234'
	role='SECURITYADMIN'
	`
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	m, err := LoadConfigFile[*ConfigDTO](configPath, true)
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

func TestLoadConfigFileWithUnknownFields(t *testing.T) {
	c := `
	[default]
	unknown='TEST_ACCOUNT'
	account_name='TEST_ACCOUNT'
	`
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	m, err := LoadConfigFile[*ConfigDTO](configPath, true)
	require.NoError(t, err)
	assert.Equal(t, map[string]*ConfigDTO{
		"default": {
			AccountName: Pointer("TEST_ACCOUNT"),
		},
	}, m)
}

func TestLoadConfigFileWithInvalidFieldTypeFails(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		wantType  string
	}{
		{name: "AccountName", fieldName: "account_name", wantType: "*string"},
		{name: "OrganizationName", fieldName: "organization_name", wantType: "*string"},
		{name: "User", fieldName: "user", wantType: "*string"},
		{name: "Username", fieldName: "username", wantType: "*string"},
		{name: "Password", fieldName: "password", wantType: "*string"},
		{name: "Host", fieldName: "host", wantType: "*string"},
		{name: "Warehouse", fieldName: "warehouse", wantType: "*string"},
		{name: "Role", fieldName: "role", wantType: "*string"},
		{name: "Params", fieldName: "params", wantType: "*map[string]*string"},
		{name: "ClientIp", fieldName: "client_ip", wantType: "*string"},
		{name: "Protocol", fieldName: "protocol", wantType: "*string"},
		{name: "Passcode", fieldName: "passcode", wantType: "*string"},
		{name: "PasscodeInPassword", fieldName: "passcode_in_password", wantType: "*bool"},
		{name: "OktaUrl", fieldName: "okta_url", wantType: "*string"},
		{name: "Authenticator", fieldName: "authenticator", wantType: "*string"},
		{name: "InsecureMode", fieldName: "insecure_mode", wantType: "*bool"},
		{name: "OcspFailOpen", fieldName: "ocsp_fail_open", wantType: "*bool"},
		{name: "Token", fieldName: "token", wantType: "*string"},
		{name: "KeepSessionAlive", fieldName: "keep_session_alive", wantType: "*bool"},
		{name: "PrivateKey", fieldName: "private_key", wantType: "*string"},
		{name: "PrivateKeyPassphrase", fieldName: "private_key_passphrase", wantType: "*string"},
		{name: "DisableTelemetry", fieldName: "disable_telemetry", wantType: "*bool"},
		{name: "ValidateDefaultParameters", fieldName: "validate_default_parameters", wantType: "*bool"},
		{name: "ClientRequestMfaToken", fieldName: "client_request_mfa_token", wantType: "*bool"},
		{name: "ClientStoreTemporaryCredential", fieldName: "client_store_temporary_credential", wantType: "*bool"},
		{name: "DriverTracing", fieldName: "driver_tracing", wantType: "*string"},
		{name: "TmpDirPath", fieldName: "tmp_dir_path", wantType: "*string"},
		{name: "DisableQueryContextCache", fieldName: "disable_query_context_cache", wantType: "*bool"},
		{name: "IncludeRetryReason", fieldName: "include_retry_reason", wantType: "*bool"},
		{name: "DisableConsoleLogin", fieldName: "disable_console_login", wantType: "*bool"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s has to have a correct type", tt.name), func(t *testing.T) {
			config := fmt.Sprintf(`
		[default]
		%s=42
		`, tt.fieldName)
			configPath := testhelpers.TestFile(t, "config", []byte(config))

			_, err := LoadConfigFile[*ConfigDTO](configPath, true)
			require.ErrorContains(t, err, fmt.Sprintf("toml: cannot decode TOML integer into struct field sdk.ConfigDTO.%s of type %s", tt.name, tt.wantType))
		})
	}
}

func TestLoadConfigFileWithInvalidFieldTypeIntFails(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
	}{
		{name: "Port", fieldName: "port"},
		{name: "ClientTimeout", fieldName: "client_timeout"},
		{name: "JwtClientTimeout", fieldName: "jwt_client_timeout"},
		{name: "LoginTimeout", fieldName: "login_timeout"},
		{name: "RequestTimeout", fieldName: "request_timeout"},
		{name: "JwtExpireTimeout", fieldName: "jwt_expire_timeout"},
		{name: "ExternalBrowserTimeout", fieldName: "external_browser_timeout"},
		{name: "MaxRetryCount", fieldName: "max_retry_count"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s has to have a correct type", tt.name), func(t *testing.T) {
			config := fmt.Sprintf(`
		[default]
		%s=value
		`, tt.fieldName)
			configPath := testhelpers.TestFile(t, "config", []byte(config))

			_, err := LoadConfigFile[*ConfigDTO](configPath, true)
			require.ErrorContains(t, err, "toml: incomplete number")
		})
	}
}

func TestLoadConfigFileWithInvalidTOMLFails(t *testing.T) {
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
			account_name=
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
			account_name="value"
			[default]
			organization_name="value"
			`,
			err: "toml: table default already exists",
		},
		{
			name: "multiple keys with the same name",
			config: `
			[default]
			password="sensitive"
			account_name="foo"
			account_name="bar"
			`,
			err: "toml: key account_name is already defined",
		},
		{
			name: "more than one key in a line",
			config: `
			[default]
			password="sensitive"
			account_name="account" organizationname="organizationname"
			`,
			err: "toml: expected newline but got U+006F 'o'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := testhelpers.TestFile(t, "config", []byte(tt.config))

			_, err := LoadConfigFile[*ConfigDTO](configPath, true)
			require.ErrorContains(t, err, tt.err)
			require.NotContains(t, err.Error(), "sensitive")
		})
	}
}

func TestProfileConfig(t *testing.T) {
	unencryptedKey, encryptedKey := random.GenerateRSAPrivateKeyEncrypted(t, "password")

	c := fmt.Sprintf(`
	[securityadmin]
	account='account'
	account_name='accountname'
	organization_name='organizationname'
	user='user'
	password='password'
	host='host'
	warehouse='warehouse'
	role='role'
	client_ip='1.1.1.1'
	protocol='http'
	passcode='passcode'
	port=1
	passcode_in_password=true
	okta_url='%[3]s'
	client_timeout=10
	jwt_client_timeout=20
	login_timeout=30
	request_timeout=40
	jwt_expire_timeout=50
	external_browser_timeout=60
	max_retry_count=1
	authenticator='SNOWFLAKE_JWT'
	insecure_mode=true
	ocsp_fail_open=true
	token='token'
	keep_session_alive=true
	private_key="""%[1]s"""
	private_key_passphrase='%[2]s'
	disable_telemetry=true
	validate_default_parameters=true
	client_request_mfa_token=true
	client_store_temporary_credential=true
	driver_tracing='tracing'
	tmp_dir_path='.'
	disable_query_context_cache=true
	include_retry_reason=true
	disable_console_login=true

	[securityadmin.params]
	foo = 'bar'
	`, encryptedKey, "password", testvars.ExampleOktaUrlString)
	configPath := testhelpers.TestFile(t, "config", []byte(c))

	t.Run("with found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("securityadmin")
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
		assert.Equal(t, true, config.PasscodeInPassword)
		assert.Equal(t, testvars.ExampleOktaUrlString, config.OktaURL.String())
		assert.Equal(t, 10*time.Second, config.ClientTimeout)
		assert.Equal(t, 20*time.Second, config.JWTClientTimeout)
		assert.Equal(t, 30*time.Second, config.LoginTimeout)
		assert.Equal(t, 40*time.Second, config.RequestTimeout)
		assert.Equal(t, 50*time.Second, config.JWTExpireTimeout)
		assert.Equal(t, 60*time.Second, config.ExternalBrowserTimeout)
		assert.Equal(t, 1, config.MaxRetryCount)
		assert.Equal(t, true, config.InsecureMode) //nolint:staticcheck
		assert.Equal(t, "token", config.Token)
		assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
		assert.Equal(t, true, config.KeepSessionAlive)
		assert.Equal(t, unencryptedKey, string(gotUnencryptedKey))
		assert.Equal(t, true, config.DisableTelemetry)
		assert.Equal(t, "tracing", config.Tracing)
		assert.Equal(t, ".", config.TmpDirPath)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientStoreTemporaryCredential)
		assert.Equal(t, true, config.DisableQueryContextCache)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
		assert.Equal(t, gosnowflake.ConfigBoolTrue, config.DisableConsoleLogin)
	})

	t.Run("with not found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("orgadmin")
		require.NoError(t, err)
		require.Nil(t, config)
	})

	t.Run("with not found config", func(t *testing.T) {
		filename := random.AlphaN(8)
		t.Setenv(snowflakeenvs.ConfigPath, filename)

		config, err := ProfileConfig("orgadmin")
		require.ErrorContains(t, err, fmt.Sprintf("could not load config file: reading information about the config file: stat %s: no such file or directory", filename))
		require.Nil(t, config)
	})
}

func TestParsingPrivateKeyDoesNotReturnSensitiveValues(t *testing.T) {
	unencryptedKey, encryptedKey := random.GenerateRSAPrivateKeyEncrypted(t, "password")

	// Make the key invalid.
	sensitive := "sensitive"
	unencryptedKey = unencryptedKey[:50] + sensitive + unencryptedKey[50:]
	_, err := ParsePrivateKey([]byte(unencryptedKey), []byte{})
	require.Error(t, err)
	require.NotContains(t, err.Error(), "PRIVATE KEY")
	require.NotContains(t, err.Error(), sensitive)

	// Use an invalid password.
	badPassword := "bad_password"
	_, err = ParsePrivateKey([]byte(encryptedKey), []byte(badPassword))
	require.Error(t, err)
	require.NotContains(t, err.Error(), "PRIVATE KEY")
	require.NotContains(t, err.Error(), badPassword)
}

func Test_MergeConfig(t *testing.T) {
	config1 := &gosnowflake.Config{
		Account:                   "account1",
		User:                      "user1",
		Password:                  "password1",
		Warehouse:                 "warehouse1",
		Role:                      "role1",
		ValidateDefaultParameters: 1,
		Params: map[string]*string{
			"foo": Pointer("1"),
		},
		ClientIP:                       net.ParseIP("1.1.1.1"),
		Protocol:                       "protocol1",
		Host:                           "host1",
		Port:                           1,
		Authenticator:                  gosnowflake.AuthTypeSnowflake,
		Passcode:                       "passcode1",
		PasscodeInPassword:             false,
		OktaURL:                        testvars.ExampleOktaUrl,
		LoginTimeout:                   1,
		RequestTimeout:                 1,
		JWTExpireTimeout:               1,
		ClientTimeout:                  1,
		JWTClientTimeout:               1,
		ExternalBrowserTimeout:         1,
		MaxRetryCount:                  1,
		InsecureMode:                   false,
		OCSPFailOpen:                   1,
		Token:                          "token1",
		KeepSessionAlive:               false,
		PrivateKey:                     random.GenerateRSAPrivateKey(t),
		DisableTelemetry:               false,
		Tracing:                        "tracing1",
		TmpDirPath:                     "tmpdirpath1",
		ClientRequestMfaToken:          gosnowflake.ConfigBoolFalse,
		ClientStoreTemporaryCredential: gosnowflake.ConfigBoolFalse,
		DisableQueryContextCache:       false,
		IncludeRetryReason:             1,
		DisableConsoleLogin:            gosnowflake.ConfigBoolFalse,
	}

	config2 := &gosnowflake.Config{
		Account:                   "account2",
		User:                      "user2",
		Password:                  "password2",
		Warehouse:                 "warehouse2",
		Role:                      "role2",
		ValidateDefaultParameters: 1,
		Params: map[string]*string{
			"foo": Pointer("2"),
		},
		ClientIP:                       net.ParseIP("2.2.2.2"),
		Protocol:                       "protocol2",
		Host:                           "host2",
		Port:                           2,
		Authenticator:                  gosnowflake.AuthTypeOAuth,
		Passcode:                       "passcode2",
		PasscodeInPassword:             true,
		OktaURL:                        testvars.ExampleOktaUrlFromEnv,
		LoginTimeout:                   2,
		RequestTimeout:                 2,
		JWTExpireTimeout:               2,
		ClientTimeout:                  2,
		JWTClientTimeout:               2,
		ExternalBrowserTimeout:         2,
		MaxRetryCount:                  2,
		InsecureMode:                   true,
		OCSPFailOpen:                   2,
		Token:                          "token2",
		KeepSessionAlive:               true,
		PrivateKey:                     random.GenerateRSAPrivateKey(t),
		DisableTelemetry:               true,
		Tracing:                        "tracing2",
		TmpDirPath:                     "tmpdirpath2",
		ClientRequestMfaToken:          gosnowflake.ConfigBoolTrue,
		ClientStoreTemporaryCredential: gosnowflake.ConfigBoolTrue,
		DisableQueryContextCache:       true,
		IncludeRetryReason:             gosnowflake.ConfigBoolTrue,
		DisableConsoleLogin:            gosnowflake.ConfigBoolTrue,
	}

	t.Run("base config empty", func(t *testing.T) {
		config := MergeConfig(&gosnowflake.Config{}, config1)

		require.Equal(t, config1, config)
	})

	t.Run("merge config empty", func(t *testing.T) {
		config := MergeConfig(config1, &gosnowflake.Config{})

		require.Equal(t, config1, config)
	})

	t.Run("both configs filled - base config takes precedence", func(t *testing.T) {
		config := MergeConfig(config1, config2)
		require.Equal(t, config1, config)
	})

	t.Run("special authenticator value", func(t *testing.T) {
		config := MergeConfig(&gosnowflake.Config{
			Authenticator: gosnowflakeAuthTypeEmpty,
		}, config1)

		require.Equal(t, config1, config)
	})
}

func Test_ToAuthenticationType(t *testing.T) {
	type test struct {
		input string
		want  gosnowflake.AuthType
	}

	valid := []test{
		// Case insensitive.
		{input: "snowflake", want: gosnowflake.AuthTypeSnowflake},

		// Supported Values.
		{input: "SNOWFLAKE", want: gosnowflake.AuthTypeSnowflake},
		{input: "OAUTH", want: gosnowflake.AuthTypeOAuth},
		{input: "EXTERNALBROWSER", want: gosnowflake.AuthTypeExternalBrowser},
		{input: "OKTA", want: gosnowflake.AuthTypeOkta},
		{input: "SNOWFLAKE_JWT", want: gosnowflake.AuthTypeJwt},
		{input: "TOKENACCESSOR", want: gosnowflake.AuthTypeTokenAccessor},
		{input: "USERNAMEPASSWORDMFA", want: gosnowflake.AuthTypeUsernamePasswordMFA},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToAuthenticatorType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToAuthenticatorType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ToExtendedAuthenticatorType(t *testing.T) {
	type test struct {
		input string
		want  gosnowflake.AuthType
	}

	valid := []test{
		// Case insensitive.
		{input: "snowflake", want: gosnowflake.AuthTypeSnowflake},

		// Supported Values.
		{input: "SNOWFLAKE", want: gosnowflake.AuthTypeSnowflake},
		{input: "OAUTH", want: gosnowflake.AuthTypeOAuth},
		{input: "EXTERNALBROWSER", want: gosnowflake.AuthTypeExternalBrowser},
		{input: "OKTA", want: gosnowflake.AuthTypeOkta},
		{input: "SNOWFLAKE_JWT", want: gosnowflake.AuthTypeJwt},
		{input: "TOKENACCESSOR", want: gosnowflake.AuthTypeTokenAccessor},
		{input: "USERNAMEPASSWORDMFA", want: gosnowflake.AuthTypeUsernamePasswordMFA},
		{input: "", want: gosnowflakeAuthTypeEmpty},
	}

	invalid := []test{
		{input: "   "},
		{input: "foo"},
		{input: "JWT"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToExtendedAuthenticatorType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToExtendedAuthenticatorType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Provider_toDriverLogLevel(t *testing.T) {
	type test struct {
		input string
		want  DriverLogLevel
	}

	valid := []test{
		// Case insensitive.
		{input: "WARNING", want: DriverLogLevelWarning},

		// Supported Values.
		{input: "trace", want: DriverLogLevelTrace},
		{input: "debug", want: DriverLogLevelDebug},
		{input: "info", want: DriverLogLevelInfo},
		{input: "print", want: DriverLogLevelPrint},
		{input: "warning", want: DriverLogLevelWarning},
		{input: "error", want: DriverLogLevelError},
		{input: "fatal", want: DriverLogLevelFatal},
		{input: "panic", want: DriverLogLevelPanic},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToDriverLogLevel(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToDriverLogLevel(tc.input)
			require.Error(t, err)
		})
	}
}
