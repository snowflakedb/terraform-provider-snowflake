package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-3825230]: change raw sqls to proper client
type ExperimentClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExperimentClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExperimentClient {
	return &ExperimentClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExperimentClient) client() *sdk.Client {
	return c.context.client
}

func (c *ExperimentClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE EXPERIMENT %s`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *ExperimentClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP EXPERIMENT IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
