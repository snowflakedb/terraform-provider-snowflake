package helpers

import (
	"context"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const (
	testDatabaseDataRetentionTimeInDays    = 1
	testDatabaseMaxDataExtensionTimeInDays = 1
)

var TestDatabaseCatalog = sdk.NewAccountObjectIdentifier("SNOWFLAKE")

type DatabaseClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewDatabaseClient(context *TestClientContext, idsGenerator *IdsGenerator) *DatabaseClient {
	return &DatabaseClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *DatabaseClient) client() sdk.Databases {
	return c.context.client.Databases
}

func (c *DatabaseClient) CreateDatabase(t *testing.T) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(c.ids.RandomAccountObjectIdentifier()))
}

// CreateDatabaseWithParametersSet should be used to create database which sets the parameters that can be altered on the account level in other tests; this way, the test is not affected by the changes.
func (c *DatabaseClient) CreateDatabaseWithParametersSet(t *testing.T) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithParametersSetWithId(t, c.ids.RandomAccountObjectIdentifier())
}

// CreateDatabaseWithParametersSetWithId should be used to create database which sets the parameters that can be altered on the account level in other tests; this way, the test is not affected by the changes.
func (c *DatabaseClient) CreateDatabaseWithParametersSetWithId(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithRequest(t, c.TestParametersSet(id))
}

// CreateTestDatabaseIfNotExists should be used to create the main database used throughout the acceptance tests.
// It's created only if it does not exist already.
func (c *DatabaseClient) CreateTestDatabaseIfNotExists(t *testing.T) (*sdk.Database, func()) {
	t.Helper()

	id := c.ids.DatabaseId()
	req := c.TestParametersSet(id).WithIfNotExists(true)

	return c.CreateDatabaseWithRequest(t, req)
}

func (c *DatabaseClient) TestParametersSet(id sdk.AccountObjectIdentifier) *sdk.CreateDatabaseRequest {
	return sdk.NewCreateDatabaseRequest(id).
		WithDataRetentionTimeInDays(testDatabaseDataRetentionTimeInDays).
		WithMaxDataExtensionTimeInDays(testDatabaseMaxDataExtensionTimeInDays).
		// according to the docs SNOWFLAKE is a valid value (https://docs.snowflake.com/en/sql-reference/parameters#catalog)
		WithCatalog(TestDatabaseCatalog)
}

func (c *DatabaseClient) CreateDatabaseWithIdentifier(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(id))
}

func (c *DatabaseClient) CreateDatabaseWithRequest(t *testing.T, request *sdk.CreateDatabaseRequest) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()

	id := request.ID()
	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	database, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return database, c.DropDatabaseFunc(t, id)
}

func (c *DatabaseClient) DropDatabaseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	return func() { require.NoError(t, c.DropDatabase(t, id)) }
}

func (c *DatabaseClient) DropDatabase(t *testing.T, id sdk.AccountObjectIdentifier) error {
	t.Helper()
	ctx := context.Background()

	if err := c.client().Drop(ctx, sdk.NewDropDatabaseRequest(id).WithIfExists(true)); err != nil {
		return err
	}
	if err := c.context.client.Sessions.UseSchema(ctx, c.ids.SchemaId()); err != nil {
		return err
	}
	return nil
}

func (c *DatabaseClient) CreateSecondaryDatabaseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, externalId sdk.ExternalObjectIdentifier, request *sdk.CreateSecondaryDatabaseRequest) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes creating secondary db right after primary creation resulted in error
	time.Sleep(1 * time.Second)

	err := c.client().CreateSecondary(ctx, request)
	require.NoError(t, err)

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes secondary database is not shown as SHOW REPLICATION DATABASES results right after creation
	time.Sleep(1 * time.Second)

	database, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := c.client().Drop(ctx, sdk.NewDropDatabaseRequest(id))
		require.NoError(t, err)

		// TODO [926148]: make this wait better with tests stabilization
		// waiting because sometimes dropping primary db right after dropping the secondary resulted in error
		time.Sleep(1 * time.Second)
		err = c.context.client.Sessions.UseSchema(ctx, c.ids.SchemaId())
		require.NoError(t, err)
	}
}

func (c *DatabaseClient) CreatePrimaryDatabase(t *testing.T, enableReplicationTo []sdk.AccountIdentifier) (*sdk.Database, sdk.ExternalObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	primaryDatabase, _ := c.CreateDatabase(t)

	err := c.client().AlterReplication(ctx, sdk.NewAlterReplicationDatabaseRequest(primaryDatabase.ID()).WithEnableReplication(
		*sdk.NewEnableReplicationRequest().WithToAccounts(enableReplicationTo).WithIgnoreEditionCheck(true),
	))
	require.NoError(t, err)

	sessionDetails, err := c.context.client.ContextFunctions.CurrentSessionDetails(ctx)
	require.NoError(t, err)

	externalPrimaryId := sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName), primaryDatabase.ID())

	// Dropping the primary database can fail until the removal of its secondary database has propagated across
	// regions, so retry the cleanup until it succeeds.
	// TODO(SNOW-1562172): Replace this retry-based workaround once there is a deterministic way to detect that the cross-region replication change has propagated.
	cleanup := func() {
		require.Eventually(t, func() bool {
			return c.DropDatabase(t, primaryDatabase.ID()) == nil
		}, 2*time.Minute, 5*time.Second)
	}

	return primaryDatabase, externalPrimaryId, cleanup
}

// WaitForReplicationToTakeEffect waits until replication of the given primary database has taken effect on
// this account, polling SHOW REPLICATION DATABASES until the primary appears. Enabling replication triggers
// cross-region operations that take a few seconds to take effect, so the primary database may not be
// immediately usable for creating a secondary database.
// TODO(SNOW-1562172): Replace this polling-based workaround once there is a deterministic way to detect that the cross-region replication change has propagated.
func (c *DatabaseClient) WaitForReplicationToTakeEffect(t *testing.T, primaryDatabaseId sdk.ExternalObjectIdentifier) {
	t.Helper()
	ctx := context.Background()
	t.Logf("Waiting for replication of primary database %s to take effect", primaryDatabaseId.FullyQualifiedName())
	require.Eventually(t, func() bool {
		replicationDatabases, err := c.context.client.ReplicationFunctions.ShowReplicationDatabases(ctx, &sdk.ShowReplicationDatabasesOptions{
			WithPrimary: &primaryDatabaseId,
		})
		if err != nil {
			t.Logf("Error showing replication databases, retrying: %v", err)
			return false
		}
		for _, replicationDatabase := range replicationDatabases {
			if replicationDatabase.IsPrimary && replicationDatabase.Name == primaryDatabaseId.Name() {
				t.Logf("Replication of primary database %s has taken effect", primaryDatabaseId.FullyQualifiedName())
				return true
			}
		}
		t.Logf("Replication of primary database %s has not taken effect yet, retrying", primaryDatabaseId.FullyQualifiedName())
		return false
	}, 2*time.Minute, 5*time.Second)
}

func (c *DatabaseClient) UpdateDataRetentionTime(t *testing.T, id sdk.AccountObjectIdentifier, days int) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterDatabaseRequest(id).WithSet(
		*sdk.NewDatabaseSetRequest().WithDataRetentionTimeInDays(days),
	))
	require.NoError(t, err)
}

func (c *DatabaseClient) UpdateLogLevel(t *testing.T, id sdk.AccountObjectIdentifier, level sdk.LogLevel) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterDatabaseRequest(id).WithSet(
		*sdk.NewDatabaseSetRequest().WithLogLevel(level),
	))
	require.NoError(t, err)
}

func (c *DatabaseClient) UnsetCatalog(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterDatabaseRequest(id).WithUnset(
		*sdk.NewDatabaseUnsetRequest().WithCatalog(true),
	))
	require.NoError(t, err)
}

func (c *DatabaseClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Database, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *DatabaseClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.DatabaseDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}

// TODO [SNOW-1562172]: Create a better solution for this type of situations
// We have to create test database from share before the actual test to check if the newly created share is ready
// after previous test (there's some kind of issue or delay between cleaning up a share and creating a new one right after).
func (c *DatabaseClient) CreateDatabaseFromShareTemporarily(t *testing.T, externalShareId sdk.ExternalObjectIdentifier) {
	t.Helper()

	db, _ := c.CreateDatabaseFromShare(t, externalShareId)

	err := c.DropDatabase(t, db.ID())
	require.NoError(t, err)
}

func (c *DatabaseClient) CreateDatabaseFromShare(t *testing.T, externalShareId sdk.ExternalObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()

	databaseId := c.ids.RandomAccountObjectIdentifier()
	err := c.client().CreateShared(context.Background(), sdk.NewCreateSharedDatabaseRequest(databaseId, externalShareId).
		// according to the docs SNOWFLAKE is a valid value (https://docs.snowflake.com/en/sql-reference/parameters#catalog)
		WithCatalog(TestDatabaseCatalog))
	require.NoError(t, err)

	database, err := c.Show(t, databaseId)
	require.NoError(t, err)

	return database, c.DropDatabaseFunc(t, databaseId)
}

func (c *DatabaseClient) ShowAllReplicationDatabases(t *testing.T) ([]sdk.ReplicationDatabase, error) {
	t.Helper()
	ctx := context.Background()

	return c.context.client.ReplicationFunctions.ShowReplicationDatabases(ctx, nil)
}

func (c *DatabaseClient) Alter(t *testing.T, request *sdk.AlterDatabaseRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *DatabaseClient) AlterReplication(t *testing.T, request *sdk.AlterReplicationDatabaseRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().AlterReplication(ctx, request)
	require.NoError(t, err)
}

func (c *DatabaseClient) AlterFailover(t *testing.T, request *sdk.AlterFailoverDatabaseRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().AlterFailover(ctx, request)
	require.NoError(t, err)
}

func (c *DatabaseClient) TestDatabaseDataRetentionTimeInDays() int {
	return testDatabaseDataRetentionTimeInDays
}

func (c *DatabaseClient) TestDatabaseMaxDataExtensionTimeInDays() int {
	return testDatabaseMaxDataExtensionTimeInDays
}

func (c *DatabaseClient) TestDatabaseCatalog() sdk.AccountObjectIdentifier {
	return TestDatabaseCatalog
}
