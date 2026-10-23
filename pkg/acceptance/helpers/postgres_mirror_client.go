package helpers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// The source objects the fixture instance must hold.
const (
	PostgresMirrorFixtureDatabase       = "postgres"
	PostgresMirrorFixtureTable1         = "tf_mirror.t1"
	PostgresMirrorFixtureTable2         = "tf_mirror.t2"
	PostgresMirrorFixtureTable3         = "tf_mirror.t3"
	PostgresMirrorFixtureMixedCaseTable = `tf_mirror."MixedCase"`
	PostgresMirrorFixtureSchemaA        = "tf_mirror_schema_a"
	PostgresMirrorFixtureSchemaB        = "tf_mirror_schema_b"
)

func PostgresMirrorNames(names ...string) []sdk.StringListItemWrapper {
	return collections.Map(names, func(name string) sdk.StringListItemWrapper {
		return sdk.StringListItemWrapper{Value: name}
	})
}

type PostgresMirrorClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewPostgresMirrorClient(context *TestClientContext, idsGenerator *IdsGenerator) *PostgresMirrorClient {
	return &PostgresMirrorClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *PostgresMirrorClient) client() sdk.PostgresMirrors {
	return c.context.client.PostgresMirrors
}

// FixtureInstance returns the pre-provisioned Postgres instance and skips the test when none is configured.
// Picking an instance automatically could create mirrors on someone else's instance.
func (c *PostgresMirrorClient) FixtureInstance(t *testing.T) string {
	t.Helper()
	return testenvs.GetOrSkipTest(t, testenvs.PostgresMirrorInstance)
}

func (c *PostgresMirrorClient) CreateWithTables(t *testing.T, instance string, tables ...string) *sdk.PostgresMirrorDetails {
	t.Helper()
	request := sdk.NewCreatePostgresMirrorRequest(*sdk.NewCreatePostgresMirrorArgumentsRequest(
		c.ids.RandomPostgresMirrorName(),
		instance,
		PostgresMirrorFixtureDatabase,
		c.ids.RandomAccountObjectIdentifier().Name(),
	).WithPostgresTables(PostgresMirrorNames(tables...)))
	return c.CreateWithRequest(t, request)
}

func (c *PostgresMirrorClient) CreateWithRequest(t *testing.T, request *sdk.CreatePostgresMirrorRequest) *sdk.PostgresMirrorDetails {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)
	t.Cleanup(c.DropFunc(t, request.GetName()))

	details, err := c.Describe(t, request.GetName())
	require.NoError(t, err)
	return details
}

// DropFunc drops the mirror together with its target database, which the mirror protects from a plain DROP DATABASE.
func (c *PostgresMirrorClient) DropFunc(t *testing.T, name string) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().DropSafely(ctx, sdk.NewDropPostgresMirrorRequest(*sdk.NewDropPostgresMirrorArgumentsRequest(name).WithDropTargetDatabase(true)))
		require.NoError(t, err)
	}
}

// DropKeptTargetDatabaseFunc drops a target database that DROP_MIRROR kept. The SNOWFLAKE application still owns it,
// so the current role takes ownership first, which needs MANAGE GRANTS.
func (c *PostgresMirrorClient) DropKeptTargetDatabaseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.context.client.Databases.ShowByIDSafely(ctx, id)
		if errors.Is(err, sdk.ErrObjectNotFound) {
			return
		}
		require.NoError(t, err)

		currentRole, err := c.context.client.ContextFunctions.CurrentRole(ctx)
		require.NoError(t, err)
		err = c.context.client.Grants.GrantOwnership(ctx, sdk.NewGrantOwnershipRequest(
			*sdk.NewOwnershipGrantOnRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: id}),
			*sdk.NewOwnershipGrantToRequest().WithAccountRoleName(currentRole.Value),
		))
		require.NoError(t, err)
		require.NoError(t, c.context.client.Databases.DropSafely(ctx, id))
	}
}

func (c *PostgresMirrorClient) Describe(t *testing.T, name string) (*sdk.PostgresMirrorDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, sdk.NewDescribePostgresMirrorRequest(*sdk.NewDescribePostgresMirrorArgumentsRequest(name)))
}

func (c *PostgresMirrorClient) Alter(t *testing.T, request *sdk.AlterPostgresMirrorRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *PostgresMirrorClient) WaitForInitialSnapshot(t *testing.T, instance string, name string, minimumTables int, timeout time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := c.client().WaitForInitialSnapshot(ctx, instance, name, minimumTables)
	require.NoError(t, err)
}
