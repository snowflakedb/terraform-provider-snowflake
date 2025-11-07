package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type StreamlitClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewStreamlitClient(context *TestClientContext, idsGenerator *IdsGenerator) *StreamlitClient {
	return &StreamlitClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *StreamlitClient) client() sdk.Streamlits {
	return c.context.client.Streamlits
}

func (c *StreamlitClient) Update(t *testing.T, request *sdk.AlterStreamlitRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *StreamlitClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Streamlit, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *StreamlitClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropStreamlitRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
