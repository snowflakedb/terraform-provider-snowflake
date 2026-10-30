package helpers

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
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

func (c *ReplicationGroupClient) CreateWithRequest(t *testing.T, request *sdk.CreateReplicationGroupRequest) {
	t.Helper()
	require.NoError(t, c.client().Create(context.Background(), request))
}

func (c *ReplicationGroupClient) Alter(t *testing.T, request *sdk.AlterReplicationGroupRequest) {
	t.Helper()
	require.NoError(t, c.client().Alter(context.Background(), request))
}

// WaitForPrimaryToBeVisible polls SHOW REPLICATION GROUPS on this account until the primary appears.
// A new primary is not always visible to the target account immediately.
func (c *ReplicationGroupClient) WaitForPrimaryToBeVisible(t *testing.T, primary sdk.ExternalObjectIdentifier) {
	t.Helper()
	ctx := context.Background()
	t.Logf("Waiting for primary replication group %s to be visible", primary.FullyQualifiedName())
	require.Eventually(t, func() bool {
		groups, err := c.client().Show(ctx, sdk.NewShowReplicationGroupRequest())
		if err != nil {
			t.Logf("Error showing replication groups, retrying: %v", err)
			return false
		}
		for _, group := range groups {
			if group.IsPrimary && group.Name == primary.Name() && reflect.DeepEqual(group.Primary, primary) {
				t.Logf("Primary replication group %s is visible", primary.FullyQualifiedName())
				return true
			}
		}
		t.Logf("Primary replication group %s is not visible yet, retrying", primary.FullyQualifiedName())
		return false
	}, time.Minute, 5*time.Second)
}

func (c *ReplicationGroupClient) Drop(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	require.Eventually(t, func() bool {
		return c.client().DropSafely(context.Background(), id) == nil
	}, 10*time.Second, time.Second)
}
