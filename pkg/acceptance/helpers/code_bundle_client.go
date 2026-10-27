package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type CodeBundleClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewCodeBundleClient(context *TestClientContext, idsGenerator *IdsGenerator) *CodeBundleClient {
	return &CodeBundleClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *CodeBundleClient) client() sdk.CodeBundles {
	return c.context.client.CodeBundles
}

func (c *CodeBundleClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().DropSafely(ctx, id)
		require.NoError(t, err)
	}
}

func (c *CodeBundleClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.CodeBundleDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeDetails(ctx, id)
}

func (c *CodeBundleClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.CodeBundle, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}
