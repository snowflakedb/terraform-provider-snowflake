package provider

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowflakedb/gosnowflake/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_toMinTLSVersion(t *testing.T) {
	valid := []struct {
		name  string
		input string
		want  minTLSVersion
	}{
		{name: "unset", input: "", want: ""},
		{name: "empty", input: "   ", want: ""},
		{name: "1.2", input: "1.2", want: minTLSVersion12},
		{name: "1.3", input: "1.3", want: minTLSVersion13},
		{name: "TLSv1.2", input: "  TLSv1.2  ", want: minTLSVersion12},
		{name: "TLSv1.3", input: "tlsv1.3", want: minTLSVersion13},
	}

	for _, tc := range valid {
		t.Run(tc.name, func(t *testing.T) {
			got, err := toMinTLSVersion(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("invalid", func(t *testing.T) {
		got, err := toMinTLSVersion("1.0")
		require.ErrorContains(t, err, "invalid min tls version: 1.0")
		assert.Empty(t, got)
	})
}

func Test_oauthRefreshHTTPClient(t *testing.T) {
	versionOf := func(t *testing.T, client *http.Client) uint16 {
		t.Helper()
		require.NotNil(t, client.Transport)
		transport, ok := client.Transport.(*http.Transport)
		require.True(t, ok)
		require.NotNil(t, transport.TLSClientConfig)
		return transport.TLSClientConfig.MinVersion
	}

	t.Run("1.2", func(t *testing.T) {
		assert.Equal(t, uint16(tls.VersionTLS12), versionOf(t, oauthRefreshHTTPClient(minTLSVersion12)))
	})

	t.Run("1.3", func(t *testing.T) {
		assert.Equal(t, uint16(tls.VersionTLS13), versionOf(t, oauthRefreshHTTPClient(minTLSVersion13)))
	})
}

func Test_Provider_toProtocol(t *testing.T) {
	type test struct {
		input string
		want  protocol
	}

	valid := []test{
		// Case insensitive.
		{input: "http", want: protocolHttp},

		// Supported Values.
		{input: "HTTP", want: protocolHttp},
		{input: "HTTPS", want: protocolHttps},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := toProtocol(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := toProtocol(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Provider_applyTfcWorkloadIdentityToken(t *testing.T) {
	oidcConfig := func() *gosnowflake.Config {
		return &gosnowflake.Config{
			Authenticator:            gosnowflake.AuthTypeWorkloadIdentityFederation,
			WorkloadIdentityProvider: tfcWorkloadIdentityProviderOidc,
		}
	}

	t.Run("empty tag leaves the token untouched", func(t *testing.T) {
		config := oidcConfig()
		config.Token = "existing token"

		require.NoError(t, applyTfcWorkloadIdentityToken("", config))
		assert.Equal(t, "existing token", config.Token)
	})

	t.Run("empty tag does not validate the authenticator", func(t *testing.T) {
		config := &gosnowflake.Config{Authenticator: gosnowflake.AuthTypeJwt}

		require.NoError(t, applyTfcWorkloadIdentityToken("", config))
		assert.Empty(t, config.Token)
	})

	t.Run("token is read from the tagged environment variable", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "tfc token")
		config := oidcConfig()

		require.NoError(t, applyTfcWorkloadIdentityToken("SNOWFLAKE", config))
		assert.Equal(t, "tfc token", config.Token)
	})

	t.Run("token takes precedence over a token set through other means", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "tfc token")
		config := oidcConfig()
		config.Token = "token from the token field"

		require.NoError(t, applyTfcWorkloadIdentityToken("SNOWFLAKE", config))
		assert.Equal(t, "tfc token", config.Token)
	})

	t.Run("tag is upper-cased when building the environment variable name", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "tfc token")
		config := oidcConfig()

		require.NoError(t, applyTfcWorkloadIdentityToken("snowflake", config))
		assert.Equal(t, "tfc token", config.Token)
	})

	t.Run("workload identity provider is matched case-insensitively", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "tfc token")
		config := oidcConfig()
		config.WorkloadIdentityProvider = "oidc"

		require.NoError(t, applyTfcWorkloadIdentityToken("SNOWFLAKE", config))
		assert.Equal(t, "tfc token", config.Token)
	})

	t.Run("unset environment variable is an error", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "")
		config := oidcConfig()

		err := applyTfcWorkloadIdentityToken("SNOWFLAKE", config)
		require.ErrorContains(t, err, "TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE is not set or is empty")
		assert.Empty(t, config.Token)
	})

	t.Run("authenticator other than WORKLOAD_IDENTITY is an error", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "tfc token")
		config := oidcConfig()
		config.Authenticator = gosnowflake.AuthTypeJwt

		err := applyTfcWorkloadIdentityToken("SNOWFLAKE", config)
		require.ErrorContains(t, err, `requires "authenticator" to be "WORKLOAD_IDENTITY"`)
		assert.Empty(t, config.Token)
	})

	t.Run("workload identity provider other than OIDC is an error", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "tfc token")
		config := oidcConfig()
		config.WorkloadIdentityProvider = "AWS"

		err := applyTfcWorkloadIdentityToken("SNOWFLAKE", config)
		require.ErrorContains(t, err, `"workload_identity_provider" to be "OIDC"`)
		assert.Empty(t, config.Token)
	})

	t.Run("unset workload identity provider is an error", func(t *testing.T) {
		t.Setenv("TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE", "tfc token")
		config := oidcConfig()
		config.WorkloadIdentityProvider = ""

		err := applyTfcWorkloadIdentityToken("SNOWFLAKE", config)
		require.ErrorContains(t, err, `"workload_identity_provider" to be "OIDC"`)
		assert.Empty(t, config.Token)
	})
}

func Test_Provider_tfcWorkloadIdentityTokenTagEnvDefault(t *testing.T) {
	t.Setenv(snowflakeenvs.TfcWorkloadIdentityTokenTag, "SNOWFLAKE")

	d := schema.TestResourceDataRaw(t, GetProviderSchema(), map[string]any{})

	assert.Equal(t, "SNOWFLAKE", d.Get("tfc_workload_identity_token_tag"))
}
