//go:build account_level_tests

package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func (c *ParameterClient) UpdateAccountParameterTemporarily(t *testing.T, parameter sdk.AccountParameter, newValue string) func() {
	t.Helper()

	param := c.ShowAccountParameter(t, parameter)
	oldValue := param.Value
	oldLevel := param.Level

	c.setAccountParameter(t, parameter, newValue)

	return func() {
		if oldLevel == "" {
			c.UnsetAccountParameter(t, parameter)
		} else {
			c.setAccountParameter(t, parameter, oldValue)
		}
	}
}

// setAccountParameter uses the catalog write path, so values are quoted by the SDK.
func (c *ParameterClient) setAccountParameter(t *testing.T, parameter sdk.AccountParameter, value string) {
	t.Helper()
	err := c.client().SetAccountParameter(context.Background(), parameter, value)
	require.NoError(t, err)
}
