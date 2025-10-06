package provider

import (
	"net"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_impl(t *testing.T) {
	_ = Provider()
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGetDriverConfigFromTerraform_EmptyConfiguration(t *testing.T) {
	d := schema.TestResourceDataRaw(t, GetProviderSchema(), map[string]interface{}{})

	config, err := getDriverConfigFromTerraform(d)

	require.NoError(t, err)
	assert.Equal(t, "terraform-provider-snowflake", config.Application)
	assert.Empty(t, config.User)
	assert.Empty(t, config.Password)
	assert.Empty(t, config.Account)
	assert.Empty(t, config.Warehouse)
	assert.Empty(t, config.Role)
	assert.Empty(t, config.Host)
	assert.Zero(t, config.Port)
	assert.Empty(t, config.Protocol)
	assert.Nil(t, config.ClientIP)
	assert.Equal(t, sdk.GosnowflakeAuthTypeEmpty, config.Authenticator)
	assert.Empty(t, config.ValidateDefaultParameters)
	assert.Empty(t, config.Passcode)
	assert.Empty(t, config.PasscodeInPassword)
	assert.Zero(t, config.LoginTimeout)
	assert.Zero(t, config.RequestTimeout)
	assert.Zero(t, config.JWTExpireTimeout)
	assert.Zero(t, config.ClientTimeout)
	assert.Zero(t, config.JWTClientTimeout)
	assert.Zero(t, config.ExternalBrowserTimeout)
	assert.Empty(t, config.InsecureMode) //nolint:staticcheck
	assert.Empty(t, config.OCSPFailOpen)
	assert.Empty(t, config.Token)
	assert.Empty(t, config.KeepSessionAlive)
	assert.Empty(t, config.DisableTelemetry)
	assert.Empty(t, config.ClientRequestMfaToken)
	assert.Empty(t, config.ClientStoreTemporaryCredential)
	assert.Empty(t, config.DisableQueryContextCache)
	assert.Empty(t, config.IncludeRetryReason)
	assert.Zero(t, config.MaxRetryCount)
	assert.Empty(t, config.Tracing)
	assert.Empty(t, config.TmpDirPath)
	assert.Empty(t, config.DisableConsoleLogin)
	assert.Empty(t, config.Params)
	assert.Empty(t, config.OauthClientID)
	assert.Empty(t, config.OauthClientSecret)
	assert.Empty(t, config.OauthTokenRequestURL)
	assert.Empty(t, config.OauthAuthorizationURL)
	assert.Empty(t, config.OauthRedirectURI)
	assert.Empty(t, config.OauthScope)
	assert.Empty(t, config.EnableSingleUseRefreshTokens)
}

func TestGetDriverConfigFromTerraform_AllFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, GetProviderSchema(), map[string]interface{}{
		"account_name":                      "test_account",
		"organization_name":                 "test_org",
		"user":                              "test_user",
		"password":                          "test_password",
		"warehouse":                         "test_warehouse",
		"role":                              "test_role",
		"host":                              "test_host",
		"port":                              443,
		"protocol":                          "https",
		"client_ip":                         "192.168.1.1",
		"authenticator":                     "SNOWFLAKE",
		"validate_default_parameters":       "true",
		"passcode":                          "123456",
		"passcode_in_password":              false,
		"login_timeout":                     60,
		"request_timeout":                   120,
		"jwt_expire_timeout":                300,
		"client_timeout":                    45,
		"jwt_client_timeout":                90,
		"external_browser_timeout":          180,
		"insecure_mode":                     false,
		"ocsp_fail_open":                    "true",
		"keep_session_alive":                true,
		"disable_telemetry":                 false,
		"client_request_mfa_token":          "true",
		"client_store_temporary_credential": "false",
		"disable_query_context_cache":       false,
		"include_retry_reason":              "true",
		"max_retry_count":                   5,
		"driver_tracing":                    "INFO",
		"tmp_directory_path":                "/tmp/snowflake",
		"disable_console_login":             "false",
		"params": map[string]interface{}{
			"QUERY_TAG": "test_tag",
			"TIMEZONE":  "UTC",
		},
		"oauth_client_id":                  "oauth_client_id",
		"oauth_client_secret":              "oauth_client_secret",
		"oauth_token_request_url":          "oauth_token_request_url",
		"oauth_authorization_url":          "oauth_authorization_url",
		"oauth_redirect_uri":               "oauth_redirect_uri",
		"oauth_scope":                      "oauth_scope",
		"enable_single_use_refresh_tokens": "true",
	})

	config, err := getDriverConfigFromTerraform(d)

	require.NoError(t, err)

	assert.Equal(t, "terraform-provider-snowflake", config.Application)
	assert.Equal(t, "test_org-test_account", config.Account)
	assert.Equal(t, "test_user", config.User)
	assert.Equal(t, "test_password", config.Password)
	assert.Equal(t, "test_warehouse", config.Warehouse)
	assert.Equal(t, "test_role", config.Role)
	assert.Equal(t, "test_host", config.Host)
	assert.Equal(t, 443, config.Port)
	assert.Equal(t, "https", config.Protocol)
	assert.Equal(t, net.ParseIP("192.168.1.1"), config.ClientIP)
	assert.Equal(t, gosnowflake.AuthTypeSnowflake, config.Authenticator)
	assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ValidateDefaultParameters)
	assert.Equal(t, "123456", config.Passcode)
	assert.False(t, config.PasscodeInPassword)
	assert.Equal(t, 60*time.Second, config.LoginTimeout)
	assert.Equal(t, 120*time.Second, config.RequestTimeout)
	assert.Equal(t, 300*time.Second, config.JWTExpireTimeout)
	assert.Equal(t, 45*time.Second, config.ClientTimeout)
	assert.Equal(t, 90*time.Second, config.JWTClientTimeout)
	assert.Equal(t, 180*time.Second, config.ExternalBrowserTimeout)
	assert.False(t, config.InsecureMode) //nolint:staticcheck
	assert.Equal(t, gosnowflake.OCSPFailOpenTrue, config.OCSPFailOpen)
	assert.Empty(t, config.Token)
	assert.True(t, config.KeepSessionAlive)
	assert.False(t, config.DisableTelemetry)
	assert.Equal(t, gosnowflake.ConfigBoolTrue, config.ClientRequestMfaToken)
	assert.Equal(t, gosnowflake.ConfigBoolFalse, config.ClientStoreTemporaryCredential)
	assert.False(t, config.DisableQueryContextCache)
	assert.Equal(t, gosnowflake.ConfigBoolTrue, config.IncludeRetryReason)
	assert.Equal(t, 5, config.MaxRetryCount)
	assert.Equal(t, "info", config.Tracing)
	assert.Equal(t, "/tmp/snowflake", config.TmpDirPath)
	assert.Equal(t, gosnowflake.ConfigBoolFalse, config.DisableConsoleLogin)
	assert.NotNil(t, config.Params)
	assert.Equal(t, "test_tag", *config.Params["QUERY_TAG"])
	assert.Equal(t, "UTC", *config.Params["TIMEZONE"])
	assert.Equal(t, "oauth_client_id", config.OauthClientID)
	assert.Equal(t, "oauth_client_secret", config.OauthClientSecret)
	assert.Equal(t, "oauth_token_request_url", config.OauthTokenRequestURL)
	assert.Equal(t, "oauth_authorization_url", config.OauthAuthorizationURL)
	assert.Equal(t, "oauth_redirect_uri", config.OauthRedirectURI)
	assert.Equal(t, "oauth_scope", config.OauthScope)
	assert.Equal(t, true, config.EnableSingleUseRefreshTokens)
}
