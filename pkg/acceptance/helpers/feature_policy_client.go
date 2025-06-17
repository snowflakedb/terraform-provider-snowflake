package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type FeaturePolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewFeaturePolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *FeaturePolicyClient {
	return &FeaturePolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *FeaturePolicyClient) client() *sdk.Client {
	return c.context.client
}

func (c *FeaturePolicyClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	// TODO(SNOW-2158888): Replace with client method when available
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE FEATURE POLICY %s BLOCKED_OBJECT_TYPES_FOR_CREATION = ("TASKS")`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *FeaturePolicyClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// TODO(SNOW-2158888): Replace with client method when available
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP FEATURE POLICY IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
