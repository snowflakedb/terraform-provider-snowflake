package helpers

import (
	"context"
	"fmt"
	"slices"
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

var refreshTerminalPhases = []string{"COMPLETED", "FAILED", "CANCELED"}

// WaitForNoRefreshInProgress polls the refresh progress of the group on this account until its last refresh has no phase in progress.
// The progress lists the phases of the last refresh: the phase in progress has no END_TIME, and neither has the closing COMPLETED, FAILED, or CANCELED row.
// The first phase appears about a second after a refresh starts, so a refresh that has just started is not visible yet.
func (c *FailoverGroupClient) WaitForNoRefreshInProgress(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()
	query := fmt.Sprintf(`SELECT PHASE_NAME, END_TIME FROM TABLE(SNOWFLAKE.INFORMATION_SCHEMA.REPLICATION_GROUP_REFRESH_PROGRESS('%s'))`, id.Name())
	require.Eventually(t, func() bool {
		phases, err := c.context.client.QueryUnsafe(ctx, query)
		if err != nil {
			t.Logf("Error reading the refresh progress of %s, retrying: %v", id.FullyQualifiedName(), err)
			return false
		}
		for _, phase := range phases {
			name, endTime := phase["PHASE_NAME"], phase["END_TIME"]
			if name != nil && !slices.Contains(refreshTerminalPhases, fmt.Sprint(*name)) && (endTime == nil || *endTime == nil) {
				t.Logf("Refresh of %s is in progress (%v), waiting", id.FullyQualifiedName(), *name)
				return false
			}
		}
		return true
	}, 2*time.Minute, time.Second)
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
