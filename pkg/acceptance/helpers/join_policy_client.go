package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-3648593): change raw sqls to proper client
type JoinPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewJoinPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *JoinPolicyClient {
	return &JoinPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *JoinPolicyClient) client() *sdk.Client {
	return c.context.client
}

func (c *JoinPolicyClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE JOIN POLICY %s AS () RETURNS JOIN_CONSTRAINT -> JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE)`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *JoinPolicyClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP JOIN POLICY IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
