package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ReplicationGroupClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewReplicationGroupClient(context *TestClientContext, idsGenerator *IdsGenerator) *ReplicationGroupClient {
	return &ReplicationGroupClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ReplicationGroupClient) client() sdk.ReplicationGroups {
	return c.context.client.ReplicationGroups
}

func (c *ReplicationGroupClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ReplicationGroup, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}
