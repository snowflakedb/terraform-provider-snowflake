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
	return c.CreateDatabaseWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateDatabaseOptions{})
}

// CreateDatabaseWithParametersSet should be used to create database which sets the parameters that can be altered on the account level in other tests; this way, the test is not affected by the changes.
func (c *DatabaseClient) CreateDatabaseWithParametersSet(t *testing.T) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithParametersSetWithId(t, c.ids.RandomAccountObjectIdentifier())
}

// CreateDatabaseWithParametersSetWithId should be used to create database which sets the parameters that can be altered on the account level in other tests; this way, the test is not affected by the changes.
func (c *DatabaseClient) CreateDatabaseWithParametersSetWithId(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithOptions(t, id, c.TestParametersSet())
}

// CreateTestDatabaseIfNotExists should be used to create the main database used throughout the acceptance tests.
// It's created only if it does not exist already.
func (c *DatabaseClient) CreateTestDatabaseIfNotExists(t *testing.T) (*sdk.Database, func()) {
	t.Helper()

	opts := c.TestParametersSet()
	opts.IfNotExists = sdk.Bool(true)

	return c.CreateDatabaseWithOptions(t, c.ids.DatabaseId(), opts)
}

func (c *DatabaseClient) TestParametersSet() *sdk.CreateDatabaseOptions {
	return &sdk.CreateDatabaseOptions{
		DataRetentionTimeInDays:    sdk.Int(testDatabaseDataRetentionTimeInDays),
		MaxDataExtensionTimeInDays: sdk.Int(testDatabaseMaxDataExtensionTimeInDays),
		// according to the docs SNOWFLAKE is a valid value (https://docs.snowflake.com/en/sql-reference/parameters#catalog)
		Catalog: sdk.Pointer(TestDatabaseCatalog),
	}
}

func (c *DatabaseClient) CreateDatabaseWithIdentifier(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	return c.CreateDatabaseWithOptions(t, id, &sdk.CreateDatabaseOptions{})
}

func (c *DatabaseClient) CreateDatabaseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, id, opts)
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

	if err := c.client().Drop(ctx, id, &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)}); err != nil {
		return err
	}
	if err := c.context.client.Sessions.UseSchema(ctx, c.ids.SchemaId()); err != nil {
		return err
	}
	return nil
}

func (c *DatabaseClient) CreateSecondaryDatabaseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, externalId sdk.ExternalObjectIdentifier, opts *sdk.CreateSecondaryDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes creating secondary db right after primary creation resulted in error
	time.Sleep(1 * time.Second)

	err := c.client().CreateSecondary(ctx, id, externalId, opts)
	require.NoError(t, err)

	// TODO [926148]: make this wait better with tests stabilization
	// waiting because sometimes secondary database is not shown as SHOW REPLICATION DATABASES results right after creation
	time.Sleep(1 * time.Second)

	database, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := c.client().Drop(ctx, id, nil)
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

	primaryDatabase, primaryDatabaseCleanup := c.CreateDatabase(t)

	err := c.client().AlterReplication(ctx, primaryDatabase.ID(), &sdk.AlterDatabaseReplicationOptions{
		EnableReplication: &sdk.EnableReplication{
			ToAccounts:         enableReplicationTo,
			IgnoreEditionCheck: sdk.Bool(true),
		},
	})
	require.NoError(t, err)

	sessionDetails, err := c.context.client.ContextFunctions.CurrentSessionDetails(ctx)
	require.NoError(t, err)

	externalPrimaryId := sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName), primaryDatabase.ID())
	return primaryDatabase, externalPrimaryId, primaryDatabaseCleanup
}

func (c *DatabaseClient) UpdateDataRetentionTime(t *testing.T, id sdk.AccountObjectIdentifier, days int) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterDatabaseOptions{
		Set: &sdk.DatabaseSet{
			DataRetentionTimeInDays: sdk.Int(days),
		},
	})
	require.NoError(t, err)
}

func (c *DatabaseClient) UnsetCatalog(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterDatabaseOptions{
		Unset: &sdk.DatabaseUnset{
			Catalog: sdk.Bool(true),
		},
	})
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
	err := c.client().CreateShared(context.Background(), databaseId, externalShareId, c.testParametersSetSharedDatabase())
	require.NoError(t, err)

	database, err := c.Show(t, databaseId)
	require.NoError(t, err)

	return database, c.DropDatabaseFunc(t, databaseId)
}

func (c *DatabaseClient) testParametersSetSharedDatabase() *sdk.CreateSharedDatabaseOptions {
	return &sdk.CreateSharedDatabaseOptions{
		// according to the docs SNOWFLAKE is a valid value (https://docs.snowflake.com/en/sql-reference/parameters#catalog)
		Catalog: sdk.Pointer(TestDatabaseCatalog),
	}
}

func (c *DatabaseClient) ShowAllReplicationDatabases(t *testing.T) ([]sdk.ReplicationDatabase, error) {
	t.Helper()
	ctx := context.Background()

	return c.context.client.ReplicationFunctions.ShowReplicationDatabases(ctx, nil)
}

func (c *DatabaseClient) Alter(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.AlterDatabaseOptions) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, opts)
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
