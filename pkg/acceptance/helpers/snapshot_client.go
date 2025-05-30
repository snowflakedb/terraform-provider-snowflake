package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-2129584): change raw sqls to proper client
type SnapshotClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSnapshotClient(context *TestClientContext, idsGenerator *IdsGenerator) *SnapshotClient {
	return &SnapshotClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SnapshotClient) client() *sdk.Client {
	return c.context.client
}

func (c *SnapshotClient) Create(t *testing.T, serviceId sdk.SchemaObjectIdentifier, volume string) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	query := fmt.Sprintf(`CREATE SNAPSHOT %s FROM SERVICE %s VOLUME "%s" INSTANCE 0`, id.FullyQualifiedName(), serviceId.FullyQualifiedName(), volume)
	_, err := c.client().ExecForTests(ctx, query)
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *SnapshotClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP SNAPSHOT IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
