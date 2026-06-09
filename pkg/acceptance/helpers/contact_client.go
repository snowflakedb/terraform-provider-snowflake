package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ContactClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewContactClient(context *TestClientContext, idsGenerator *IdsGenerator) *ContactClient {
	return &ContactClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ContactClient) client() *sdk.Client {
	return c.context.client
}

// TODO(SNOW-2175834): Replace raw SQL with SDK client once Contacts SDK is implemented.
func (c *ContactClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	stmt := fmt.Sprintf(`CREATE CONTACT %s`, id.FullyQualifiedName())
	_, err := c.client().ExecForTests(ctx, stmt)
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *ContactClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP CONTACT IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
