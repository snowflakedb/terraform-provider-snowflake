package helpers

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const (
	// OpenflowDeploymentTerminateTimeout bounds the wait for TERMINATE to reach DELETED. Terminating a
	// deployment tears down infrastructure asynchronously and can take several minutes.
	OpenflowDeploymentTerminateTimeout = 20 * time.Minute
	// OpenflowDeploymentActiveTimeout bounds the wait for a Snowflake-managed deployment to reach ACTIVE,
	// which takes several minutes.
	OpenflowDeploymentActiveTimeout = 30 * time.Minute
)

type OpenflowDeploymentClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewOpenflowDeploymentClient(context *TestClientContext, idsGenerator *IdsGenerator) *OpenflowDeploymentClient {
	return &OpenflowDeploymentClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *OpenflowDeploymentClient) client() sdk.OpenflowDeployments {
	return c.context.client.OpenflowDeployments
}

// CreateByoc creates a BYOC deployment, which settles at INACTIVE without waiting on infrastructure.
// Creating a Snowflake-managed deployment provisions it for real and takes far longer.
func (c *OpenflowDeploymentClient) CreateByoc(t *testing.T) (*sdk.OpenflowDeployment, func()) {
	t.Helper()
	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateOpenflowDeploymentRequest(id, sdk.OpenflowDeploymentTypeByoc).
		WithVpcType(sdk.OpenflowVpcTypeManaged))
}

// CreateWithRequest creates a deployment and waits for it to settle, since Snowflake refuses ALTER SET,
// TERMINATE and DROP while one is CREATING. The returned deployment reflects its settled status: ACTIVE for
// Snowflake-managed, INACTIVE for BYOC.
func (c *OpenflowDeploymentClient) CreateWithRequest(t *testing.T, req *sdk.CreateOpenflowDeploymentRequest) (*sdk.OpenflowDeployment, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	cleanup := c.DropFunc(t, req.GetName())

	c.WaitUntilSettled(t, req.GetName(), OpenflowDeploymentActiveTimeout)

	deployment, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return deployment, cleanup
}

// DropFunc tears a deployment down the way Snowflake requires: TERMINATE, wait for DELETED, then DROP.
// Every step is guarded, so a test that already removed the deployment does not fail its own cleanup.
func (c *OpenflowDeploymentClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		deployment, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			// Already gone; nothing to tear down.
			return
		}

		if !slices.Contains([]sdk.OpenflowDeploymentStatus{
			sdk.OpenflowDeploymentStatusDeleting,
			sdk.OpenflowDeploymentStatusDeleted,
		}, deployment.Status) {
			// TERMINATE is rejected while the deployment is CREATING, which is the state cleanup runs in when
			// a test fails its first assertion.
			c.WaitUntilSettled(t, id, OpenflowDeploymentActiveTimeout)

			err = c.client().Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithIfExists(true).WithTerminate(true))
			require.NoError(t, err)
		}

		c.WaitForStatus(t, id, sdk.OpenflowDeploymentStatusDeleted, OpenflowDeploymentTerminateTimeout)

		err = c.client().DropSafely(ctx, id)
		require.NoError(t, err)
	}
}

func (c *OpenflowDeploymentClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.OpenflowDeployment, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *OpenflowDeploymentClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.OpenflowDeploymentDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}

func (c *OpenflowDeploymentClient) Alter(t *testing.T, req *sdk.AlterOpenflowDeploymentRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

// ShowParameters returns the deployment-level parameters. This is the only way to read EVENT_TABLE back,
// as it appears in neither SHOW nor DESCRIBE output.
func (c *OpenflowDeploymentClient) ShowParameters(t *testing.T, id sdk.AccountObjectIdentifier) []*sdk.Parameter {
	t.Helper()
	ctx := context.Background()

	parameters, err := c.client().ShowParameters(ctx, id)
	require.NoError(t, err)
	return parameters
}

// ActiveDeploymentForRuntimes returns the deployment named by TEST_SF_TF_OPENFLOW_DEPLOYMENT. That it
// exists and is ACTIVE is a suite prerequisite, checked once by
// TestClient.EnsureOpenflowDeploymentIsProvisioned rather than here.
func (c *OpenflowDeploymentClient) ActiveDeploymentForRuntimes(t *testing.T) sdk.AccountObjectIdentifier {
	t.Helper()

	name := os.Getenv(string(testenvs.OpenflowDeployment))
	require.NotEmptyf(t, name, "%s must name an ACTIVE Openflow deployment", testenvs.OpenflowDeployment)
	return sdk.NewAccountObjectIdentifier(name)
}

// WaitForStatus waits for the given status, bailing out on a terminal failure status rather than holding
// the suite for the full timeout and then reporting only that a condition was never satisfied.
func (c *OpenflowDeploymentClient) WaitForStatus(t *testing.T, id sdk.AccountObjectIdentifier, status sdk.OpenflowDeploymentStatus, timeout time.Duration) {
	t.Helper()
	ctx := context.Background()

	var last sdk.OpenflowDeploymentStatus
	err := util.Retry(int(timeout/(10*time.Second)), 10*time.Second, func() (error, bool) {
		deployment, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			return err, true
		}
		last = deployment.Status
		if deployment.Status == status {
			return nil, true
		}
		if slices.Contains(sdk.OpenflowDeploymentFailureStatuses, deployment.Status) {
			return fmt.Errorf("deployment %s reached a terminal failure status %s, waiting for %s", id.FullyQualifiedName(), deployment.Status, status), true
		}
		return nil, false
	})
	require.NoErrorf(t, err, "deployment %s did not reach %s within %v (last status %s)", id.FullyQualifiedName(), status, timeout, last)
}

// WaitUntilSettled waits for a deployment to stop coming up, whether it lands on ACTIVE, INACTIVE or
// CREATE_FAILED. Only settled deployments accept ALTER SET, TERMINATE or DROP. A Snowflake-managed deployment
// goes CREATING, PROVISIONING, ACTIVE, so PROVISIONING is transient too; BYOC goes straight to INACTIVE.
func (c *OpenflowDeploymentClient) WaitUntilSettled(t *testing.T, id sdk.AccountObjectIdentifier, timeout time.Duration) {
	t.Helper()
	ctx := context.Background()

	var last sdk.OpenflowDeploymentStatus
	err := util.Retry(int(timeout/(10*time.Second)), 10*time.Second, func() (error, bool) {
		deployment, err := c.client().ShowByIDSafely(ctx, id)
		if err != nil {
			return err, true
		}
		last = deployment.Status
		return nil, !slices.Contains(sdk.OpenflowDeploymentTransientStatuses, deployment.Status)
	})
	require.NoErrorf(t, err, "deployment %s did not settle within %v (last status %s)", id.FullyQualifiedName(), timeout, last)
}
