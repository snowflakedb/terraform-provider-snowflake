package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type FailoverGroupClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewFailoverGroupClient(context *TestClientContext, idsGenerator *IdsGenerator) *FailoverGroupClient {
	return &FailoverGroupClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *FailoverGroupClient) client() sdk.FailoverGroups {
	return c.context.client.FailoverGroups
}

func (c *FailoverGroupClient) Create(t *testing.T) (*sdk.FailoverGroup, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateFailoverGroupRequest(
		c.ids.RandomAccountObjectIdentifier(),
		[]sdk.FailoverGroupObjectType{sdk.FailoverGroupObjectTypeRoles},
		[]sdk.AccountIdentifier{c.ids.AccountIdentifierWithLocator()},
	))
}

func (c *FailoverGroupClient) CreateWithRequest(t *testing.T, req *sdk.CreateFailoverGroupRequest) (*sdk.FailoverGroup, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	failoverGroup, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return failoverGroup, c.DropFunc(t, req.GetName())
}

func (c *FailoverGroupClient) RemoveAllowedAccounts(t *testing.T, id sdk.AccountObjectIdentifier, accounts ...sdk.AccountIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterFailoverGroupRequest(id).WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedAccounts(accounts)))
	require.NoError(t, err)
}

func (c *FailoverGroupClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropFailoverGroupRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *FailoverGroupClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.FailoverGroup, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *FailoverGroupClient) CreateSecondary(t *testing.T, id sdk.AccountObjectIdentifier, primary sdk.ExternalObjectIdentifier) {
	t.Helper()
	require.NoError(t, c.client().CreateSecondary(context.Background(), sdk.NewCreateSecondaryFailoverGroupRequest(id, primary)))
}

func (c *FailoverGroupClient) Alter(t *testing.T, request *sdk.AlterFailoverGroupRequest) {
	t.Helper()
	require.NoError(t, c.client().Alter(context.Background(), request))
}

func (c *FailoverGroupClient) AlterSecondary(t *testing.T, request *sdk.AlterSecondaryFailoverGroupRequest) {
	t.Helper()
	require.NoError(t, c.client().AlterSecondary(context.Background(), request))
}

// WaitForPrimaryToBeVisible polls SHOW FAILOVER GROUPS on this account until the primary appears.
// A new primary is not always visible to the target account immediately.
func (c *FailoverGroupClient) WaitForPrimaryToBeVisible(t *testing.T, primary sdk.ExternalObjectIdentifier) {
	t.Helper()
	ctx := context.Background()
	t.Logf("Waiting for primary failover group %s to be visible", primary.FullyQualifiedName())
	require.Eventually(t, func() bool {
		groups, err := c.client().Show(ctx, sdk.NewShowFailoverGroupRequest())
		if err != nil {
			t.Logf("Error showing failover groups, retrying: %v", err)
			return false
		}
		for _, group := range groups {
			if group.IsPrimary && group.Name == primary.Name() && sameAccount(group.Primary.AccountIdentifier(), primary.AccountIdentifier()) {
				t.Logf("Primary failover group %s is visible", primary.FullyQualifiedName())
				return true
			}
		}
		t.Logf("Primary failover group %s is not visible yet, retrying", primary.FullyQualifiedName())
		return false
	}, time.Minute, 5*time.Second)
}

func (c *FailoverGroupClient) Drop(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	require.Eventually(t, func() bool {
		return c.client().DropSafely(context.Background(), id) == nil
	}, 10*time.Second, time.Second)
}

// SHOW prints the primary without the account locator, while FailoverGroup.ExternalID carries it.
func sameAccount(a, b sdk.AccountIdentifier) bool {
	return a.OrganizationName() == b.OrganizationName() && a.AccountName() == b.AccountName()
}
