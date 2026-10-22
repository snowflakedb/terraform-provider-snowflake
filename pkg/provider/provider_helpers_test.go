package provider

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowflakedb/gosnowflake/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
