package helpers

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const (
	// OpenflowConnectorRunningTimeout bounds the wait for START to reach RUNNING.
	OpenflowConnectorRunningTimeout = 20 * time.Minute
	// OpenflowConnectorStopTimeout bounds the wait for STOP to reach STOPPED.
	OpenflowConnectorStopTimeout = 20 * time.Minute
	// OpenflowConnectorTerminateTimeout bounds the wait for TERMINATE to reach DELETED.
	OpenflowConnectorTerminateTimeout = 20 * time.Minute
)

type OpenflowConnectorClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewOpenflowConnectorClient(context *TestClientContext, idsGenerator *IdsGenerator) *OpenflowConnectorClient {
	return &OpenflowConnectorClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *OpenflowConnectorClient) client() sdk.OpenflowConnectors {
	return c.context.client.OpenflowConnectors
}

// CreateFromDefinition instantiates a connector from a Snowflake-managed definition.
func (c *OpenflowConnectorClient) CreateFromDefinition(t *testing.T, runtimeId sdk.SchemaObjectIdentifier, definition string) (*sdk.OpenflowConnector, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
	return c.CreateWithRequest(t, sdk.NewCreateOpenflowConnectorRequest(id, runtimeId).WithFromDefinition(definition))
}

// CreateFromDefinitionAndCommit creates a connector and commits its live version, which several operations
// require: START is refused until a default version exists, and DEFAULT_VERSION can only name a committed
// one. The returned connector reflects the settled state after the commit.
func (c *OpenflowConnectorClient) CreateFromDefinitionAndCommit(t *testing.T, runtimeId sdk.SchemaObjectIdentifier, definition string) (*sdk.OpenflowConnector, func()) {
	t.Helper()
	ctx := context.Background()

	connector, cleanup := c.CreateFromDefinition(t, runtimeId, definition)
	require.NoError(t, c.client().Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(connector.ID()).
		WithCommit(*sdk.NewOpenflowConnectorCommitRequest())))

	return c.WaitUntilSettled(t, connector.ID(), OpenflowConnectorStopTimeout), cleanup
}

// CreateWithRequest creates a connector and waits for it to settle, since Snowflake refuses COMMIT,
// TERMINATE and DROP while one is CREATING. A new connector settles on STOPPED, as a draft.
func (c *OpenflowConnectorClient) CreateWithRequest(t *testing.T, req *sdk.CreateOpenflowConnectorRequest) (*sdk.OpenflowConnector, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	cleanup := c.DropFunc(t, req.GetName())

	connector := c.WaitUntilSettled(t, req.GetName(), OpenflowConnectorStopTimeout)

	return connector, cleanup
}

// DropFunc tears a connector down the way Snowflake requires: STOP if running, TERMINATE, wait for
// DELETED, then DROP. It waits for the connector to settle first, since TERMINATE is refused while
// CREATING.
func (c *OpenflowConnectorClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		if _, err := c.client().ShowByIDSafely(ctx, id); err != nil {
			// Already gone; nothing to tear down.
			return
		}

		connector := c.WaitUntilSettled(t, id, OpenflowConnectorStopTimeout)

		alreadyTerminating := slices.Contains([]sdk.OpenflowConnectorStatus{
			sdk.OpenflowConnectorStatusDeleting,
			sdk.OpenflowConnectorStatusDeleted,
		}, connector.Status)

		if !alreadyTerminating {
			if slices.Contains([]sdk.OpenflowConnectorStatus{
				sdk.OpenflowConnectorStatusRunning,
				sdk.OpenflowConnectorStatusStarting,
				sdk.OpenflowConnectorStatusStartFailed,
			}, connector.Status) {
				err := c.client().Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithIfExists(true).WithStop(true))
				require.NoError(t, err)
				c.WaitForStatus(t, id, sdk.OpenflowConnectorStatusStopped, OpenflowConnectorStopTimeout)
			}

			err := c.client().Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithIfExists(true).WithTerminate(true))
			require.NoError(t, err)
		}

		c.WaitForStatus(t, id, sdk.OpenflowConnectorStatusDeleted, OpenflowConnectorTerminateTimeout)

		err := c.client().DropSafely(ctx, id)
		require.NoError(t, err)
	}
}

func (c *OpenflowConnectorClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.OpenflowConnector, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *OpenflowConnectorClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.OpenflowConnectorDetails, error) {
	t.Helper()
	ctx := context.Background()

	details, err := c.client().Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	// DESCRIBE does not return database_name or schema_name, so thread the identifier through.
	details.Id = id
	return details, nil
}

func (c *OpenflowConnectorClient) Alter(t *testing.T, req *sdk.AlterOpenflowConnectorRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

// WaitForStatus waits for the given status, bailing out on a terminal failure status. require.Eventually
// cannot bail early, so a START_FAILED connector would sit there for the full twenty minutes.
func (c *OpenflowConnectorClient) WaitForStatus(t *testing.T, id sdk.SchemaObjectIdentifier, status sdk.OpenflowConnectorStatus, timeout time.Duration) {
	t.Helper()
	ctx := context.Background()

	var last sdk.OpenflowConnectorStatus
	err := util.Retry(int(timeout/(10*time.Second)), 10*time.Second, func() (error, bool) {
		connector, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			return err, true
		}
		last = connector.Status
		if connector.Status == status {
			return nil, true
		}
		if slices.Contains(sdk.OpenflowConnectorFailureStatuses, connector.Status) {
			return fmt.Errorf("connector %s reached a terminal failure status %s, waiting for %s", id.FullyQualifiedName(), connector.Status, status), true
		}
		return nil, false
	})
	require.NoErrorf(t, err, "connector %s did not reach %s within %v (last status %s)", id.FullyQualifiedName(), status, timeout, last)
}

// WaitForAnyStatus waits for any one of the given statuses, for cases with more than one legitimate
// outcome: starting an unconfigured connector lands on either RUNNING or START_FAILED.
func (c *OpenflowConnectorClient) WaitForAnyStatus(t *testing.T, id sdk.SchemaObjectIdentifier, statuses []sdk.OpenflowConnectorStatus, timeout time.Duration) {
	t.Helper()
	ctx := context.Background()

	var last sdk.OpenflowConnectorStatus
	err := util.Retry(int(timeout/(10*time.Second)), 10*time.Second, func() (error, bool) {
		connector, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			return err, true
		}
		last = connector.Status
		if slices.Contains(statuses, connector.Status) {
			return nil, true
		}
		if slices.Contains(sdk.OpenflowConnectorFailureStatuses, connector.Status) {
			return fmt.Errorf("connector %s reached an unexpected terminal failure status %s, waiting for one of %v", id.FullyQualifiedName(), connector.Status, statuses), true
		}
		return nil, false
	})
	require.NoErrorf(t, err, "connector %s did not reach one of %v within %v (last status %s)", id.FullyQualifiedName(), statuses, timeout, last)
}

// WaitUntilSettled waits out an operation still in flight. Snowflake refuses COMMIT, TERMINATE and DROP
// during those. Returns the connector as observed on the last successful poll, so callers do not need to
// fetch it again.
func (c *OpenflowConnectorClient) WaitUntilSettled(t *testing.T, id sdk.SchemaObjectIdentifier, timeout time.Duration) *sdk.OpenflowConnector {
	t.Helper()
	ctx := context.Background()

	var last *sdk.OpenflowConnector
	err := util.Retry(int(timeout/(10*time.Second)), 10*time.Second, func() (error, bool) {
		connector, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			return err, true
		}
		last = connector
		return nil, !slices.Contains(sdk.OpenflowConnectorTransientStatuses, connector.Status)
	})
	var lastStatus sdk.OpenflowConnectorStatus
	if last != nil {
		lastStatus = last.Status
	}
	require.NoErrorf(t, err, "connector %s did not settle within %v (last status %s)", id.FullyQualifiedName(), timeout, lastStatus)
	return last
}
