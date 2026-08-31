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
	// OpenflowRuntimeActiveTimeout bounds the wait for a runtime to reach ACTIVE. Creation has been observed
	// to take four to five minutes.
	OpenflowRuntimeActiveTimeout = 30 * time.Minute
	// OpenflowRuntimeSuspendTimeout bounds the wait for SUSPEND to reach SUSPENDED.
	OpenflowRuntimeSuspendTimeout = 20 * time.Minute
	// OpenflowRuntimeTerminateTimeout bounds the wait for TERMINATE to reach DELETED.
	OpenflowRuntimeTerminateTimeout = 20 * time.Minute
)

type OpenflowRuntimeClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewOpenflowRuntimeClient(context *TestClientContext, idsGenerator *IdsGenerator) *OpenflowRuntimeClient {
	return &OpenflowRuntimeClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *OpenflowRuntimeClient) client() sdk.OpenflowRuntimes {
	return c.context.client.OpenflowRuntimes
}

// Create makes a minimal runtime in the given deployment. EXECUTE_AS_ROLE, NODE_TYPE, MIN_NODES and
// MAX_NODES are all required by CREATE OPENFLOW RUNTIME.
func (c *OpenflowRuntimeClient) Create(t *testing.T, deploymentId sdk.AccountObjectIdentifier, executeAsRoleId sdk.AccountObjectIdentifier) (*sdk.OpenflowRuntime, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateOpenflowRuntimeRequest(id, deploymentId, sdk.OpenflowRuntimeNodeTypeSmall, 1, 1, executeAsRoleId))
}

func (c *OpenflowRuntimeClient) CreateWithRequest(t *testing.T, req *sdk.CreateOpenflowRuntimeRequest) (*sdk.OpenflowRuntime, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	runtime, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return runtime, c.DropFunc(t, req.GetName())
}

// DropFunc tears a runtime down the way Snowflake requires: TERMINATE CASCADE, wait for DELETED, then DROP
// CASCADE. No SUSPEND step: TERMINATE is accepted from ACTIVE and from CREATE_FAILED, where SUSPEND is not.
// Waits for the runtime to settle, since TERMINATE is refused while it is CREATING.
func (c *OpenflowRuntimeClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		if _, err := c.client().ShowByIDSafely(ctx, id); err != nil {
			// Already gone; nothing to tear down.
			return
		}

		runtime := c.WaitUntilSettled(t, id, OpenflowRuntimeActiveTimeout)

		alreadyTerminating := slices.Contains([]sdk.OpenflowRuntimeStatus{
			sdk.OpenflowRuntimeStatusDeleting,
			sdk.OpenflowRuntimeStatusDeleted,
		}, runtime.Status)

		if !alreadyTerminating {
			err := c.client().Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithIfExists(true).WithTerminateCascade(true))
			require.NoError(t, err)
		}

		c.WaitForStatus(t, id, sdk.OpenflowRuntimeStatusDeleted, OpenflowRuntimeTerminateTimeout)

		err := c.client().Drop(ctx, sdk.NewDropOpenflowRuntimeRequest(id).WithIfExists(true).WithCascade(true))
		require.NoError(t, err)
	}
}

func (c *OpenflowRuntimeClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.OpenflowRuntime, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *OpenflowRuntimeClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.OpenflowRuntimeDetails, error) {
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

func (c *OpenflowRuntimeClient) Alter(t *testing.T, req *sdk.AlterOpenflowRuntimeRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

// WaitForStatus waits for the given status, bailing out on a terminal failure status. The ACTIVE timeout is
// thirty minutes, so without that a CREATE_FAILED runtime holds the suite for all of it.
func (c *OpenflowRuntimeClient) WaitForStatus(t *testing.T, id sdk.SchemaObjectIdentifier, status sdk.OpenflowRuntimeStatus, timeout time.Duration) {
	t.Helper()
	ctx := context.Background()

	var last sdk.OpenflowRuntimeStatus
	err := util.Retry(int(timeout/(10*time.Second)), 10*time.Second, func() (error, bool) {
		runtime, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			return err, true
		}
		last = runtime.Status
		if runtime.Status == status {
			return nil, true
		}
		if slices.Contains(sdk.OpenflowRuntimeFailureStatuses, runtime.Status) {
			return fmt.Errorf("runtime %s reached a terminal failure status %s, waiting for %s", id.FullyQualifiedName(), runtime.Status, status), true
		}
		return nil, false
	})
	require.NoErrorf(t, err, "runtime %s did not reach %s within %v (last status %s)", id.FullyQualifiedName(), status, timeout, last)
}

// WaitUntilSettled waits out an operation still in flight. Snowflake refuses SUSPEND, TERMINATE and DROP
// during those. Returns the runtime as observed on the last successful poll, so callers do not need to
// fetch it again.
func (c *OpenflowRuntimeClient) WaitUntilSettled(t *testing.T, id sdk.SchemaObjectIdentifier, timeout time.Duration) *sdk.OpenflowRuntime {
	t.Helper()
	ctx := context.Background()

	var last *sdk.OpenflowRuntime
	err := util.Retry(int(timeout/(10*time.Second)), 10*time.Second, func() (error, bool) {
		runtime, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			return err, true
		}
		last = runtime
		return nil, !slices.Contains(sdk.OpenflowRuntimeTransientStatuses, runtime.Status)
	})
	var lastStatus sdk.OpenflowRuntimeStatus
	if last != nil {
		lastStatus = last.Status
	}
	require.NoErrorf(t, err, "runtime %s did not settle within %v (last status %s)", id.FullyQualifiedName(), timeout, lastStatus)
	return last
}
