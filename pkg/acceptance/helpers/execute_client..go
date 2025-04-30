package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExecuteClient struct {
	context *TestClientContext
}

func NewExecuteClient(context *TestClientContext) *ExecuteClient {
	return &ExecuteClient{
		context: context,
	}
}

func (c *ExecuteClient) client() *sdk.Client {
	return c.context.client
}

func (c *ExecuteClient) SQL(t *testing.T, sql string) {
	t.Helper()
	_, err := c.client().ExecForTests(context.Background(), sql)
	require.NoError(t, err)
}
