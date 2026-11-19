//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// The parameter catalog supports SERVICE_CALLER_TOKEN_VALIDITY_SECS only on services, so on the account level the value is rendered by the kind.
func TestInt_Parameters_AccountParameterWithGenericFallback(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	parameter := sdk.AccountParameterServiceCallerTokenValiditySecs
	t.Cleanup(func() { testClientHelper().Parameter.UnsetAccountParameter(t, parameter) })

	require.NoError(t, client.Parameters.SetAccountParameterWithGenericFallback(ctx, string(parameter), sdk.ParameterValueKindInteger, "120"))
	p, err := client.Parameters.ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	require.Equal(t, "120", p.Value)
	require.Equal(t, sdk.ParameterTypeAccount, p.Level)

	require.NoError(t, client.Parameters.UnsetAccountParameterWithGenericFallback(ctx, string(parameter)))
	p, err = client.Parameters.ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	require.NotEqual(t, sdk.ParameterTypeAccount, p.Level)
}
