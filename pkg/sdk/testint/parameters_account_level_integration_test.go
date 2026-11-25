//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// The parameter catalog does not support DEFAULT_MAX_PER_USER_CONCURRENT_STATEMENT_LIMIT, so the value is rendered by the kind.
func TestInt_Parameters_AccountParameterWithGenericFallback(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	parameter := sdk.AccountParameter("DEFAULT_MAX_PER_USER_CONCURRENT_STATEMENT_LIMIT")
	t.Cleanup(func() { testClientHelper().Parameter.UnsetAccountParameter(t, parameter) })

	require.NoError(t, client.Parameters.SetAccountParameterWithGenericFallback(ctx, string(parameter), sdk.ParameterValueKindInteger, "2400"))
	p, err := client.Parameters.ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	require.Equal(t, "2400", p.Value)
	require.Equal(t, sdk.ParameterTypeAccount, p.Level)

	require.NoError(t, client.Parameters.UnsetAccountParameterWithGenericFallback(ctx, string(parameter)))
	p, err = client.Parameters.ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	require.NotEqual(t, sdk.ParameterTypeAccount, p.Level)
}
