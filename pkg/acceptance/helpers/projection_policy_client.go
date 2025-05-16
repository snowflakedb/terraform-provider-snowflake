package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1564959): change raw sqls to proper client
type ProjectionPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewProjectionPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *ProjectionPolicyClient {
	return &ProjectionPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ProjectionPolicyClient) client() *sdk.Client {
	return c.context.client
}

func (c *ProjectionPolicyClient) CreateProjectionPolicy(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE PROJECTION POLICY %s AS () RETURNS PROJECTION_CONSTRAINT -> PROJECTION_CONSTRAINT(ALLOW => false)`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropProjectionPolicyFunc(t, id)
}

func (c *ProjectionPolicyClient) DropProjectionPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP PROJECTION POLICY IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
