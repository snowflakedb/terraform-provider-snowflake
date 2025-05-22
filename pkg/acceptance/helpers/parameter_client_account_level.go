//go:build account_level_tests

package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func (c *ParameterClient) UpdateAccountParameterTemporarily(t *testing.T, parameter sdk.AccountParameter, newValue string) func() {
	t.Helper()
	ctx := context.Background()

	param := c.ShowAccountParameter(t, parameter)
	oldValue := param.Value
	oldLevel := param.Level

	err := c.client().SetAccountParameter(ctx, parameter, newValue)
	require.NoError(t, err)

	return func() {
		if oldLevel == "" {
			c.UnsetAccountParameter(t, parameter)
		} else {
			err := c.client().SetAccountParameter(ctx, parameter, oldValue)
			require.NoError(t, err)
		}
	}
}
