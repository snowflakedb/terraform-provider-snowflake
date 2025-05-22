package testint

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random/integrationtests"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
)

const IntegrationTestPrefix = "int_test_"

var (
	TestWarehouseName = fmt.Sprintf("%swh_%s", IntegrationTestPrefix, integrationtests.ObjectsSuffix)
	TestDatabaseName  = fmt.Sprintf("%sdb_%s", IntegrationTestPrefix, integrationtests.ObjectsSuffix)
	TestSchemaName    = fmt.Sprintf("%ssc_%s", IntegrationTestPrefix, integrationtests.ObjectsSuffix)

	NonExistingAccountObjectIdentifier                                             = sdk.NewAccountObjectIdentifier("does_not_exist")
	NonExistingDatabaseObjectIdentifier                                            = sdk.NewDatabaseObjectIdentifier(TestDatabaseName, "does_not_exist")
	NonExistingDatabaseObjectIdentifierWithNonExistingDatabase                     = sdk.NewDatabaseObjectIdentifier("does_not_exist", "does_not_exist")
	NonExistingSchemaObjectIdentifier                                              = sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, "does_not_exist")
	NonExistingSchemaObjectIdentifierWithNonExistingSchema                         = sdk.NewSchemaObjectIdentifier(TestDatabaseName, "does_not_exist", "does_not_exist")
	NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema              = sdk.NewSchemaObjectIdentifier("does_not_exist", "does_not_exist", "does_not_exist")
	NonExistingSchemaObjectIdentifierWithArguments                                 = sdk.NewSchemaObjectIdentifierWithArguments(TestDatabaseName, TestSchemaName, "does_not_exist", sdk.DataTypeInt)
	NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingSchema            = sdk.NewSchemaObjectIdentifierWithArguments(TestDatabaseName, "does_not_exist", "does_not_exist", sdk.DataTypeInt)
	NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingDatabaseAndSchema = sdk.NewSchemaObjectIdentifierWithArguments("does_not_exist", "does_not_exist", "does_not_exist", sdk.DataTypeInt)
)

var itc integrationTestContext

func TestMain(m *testing.M) {
	exitVal := execute(m)
	os.Exit(exitVal)
}

func execute(m *testing.M) int {
	defer timer("tests")()
	setup()
	exitVal := m.Run()
	cleanup()
	return exitVal
}

func setup() {
	log.Println("[DEBUG] Running integration tests setup")

	err := itc.initialize()
	if err != nil {
		log.Printf("[DEBUG] Integration test context initialisation failed with: `%s`", err)
		cleanup()
		os.Exit(1)
	}
}

func cleanup() {
	log.Println("[DEBUG] Running integration tests cleanup")
	if itc.databaseCleanup != nil {
		defer itc.databaseCleanup()
	}
	if itc.schemaCleanup != nil {
		defer itc.schemaCleanup()
	}
	if itc.warehouseCleanup != nil {
		defer itc.warehouseCleanup()
	}
	if itc.secondaryDatabaseCleanup != nil {
		defer itc.secondaryDatabaseCleanup()
	}
	if itc.secondarySchemaCleanup != nil {
		defer itc.secondarySchemaCleanup()
	}
	if itc.secondaryWarehouseCleanup != nil {
		defer itc.secondaryWarehouseCleanup()
	}
}

type integrationTestContext struct {
	config *gosnowflake.Config
	client *sdk.Client
	ctx    context.Context

	database         *sdk.Database
	databaseCleanup  func()
	schema           *sdk.Schema
	schemaCleanup    func()
	warehouse        *sdk.Warehouse
	warehouseCleanup func()

	secondaryClient *sdk.Client
	secondaryCtx    context.Context

	secondaryDatabase         *sdk.Database
	secondaryDatabaseCleanup  func()
	secondarySchema           *sdk.Schema
	secondarySchemaCleanup    func()
	secondaryWarehouse        *sdk.Warehouse
	secondaryWarehouseCleanup func()

	testClient          *helpers.TestClient
	secondaryTestClient *helpers.TestClient
}

func (itc *integrationTestContext) initialize() error {
	log.Println("[DEBUG] Initializing integration test context")

	testObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.TestObjectsSuffix))
	requireTestObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.RequireTestObjectsSuffix))
	if requireTestObjectSuffix != "" && testObjectSuffix == "" {
		return fmt.Errorf("Test object suffix is required for this test run. Set %s env.", testenvs.TestObjectsSuffix)
	}

	defaultConfig, err := sdk.ProfileConfig(testprofiles.Default)
	if err != nil {
		return err
	}
	if defaultConfig == nil {
		return errors.New("config is required to run integration tests")
	}
	itc.config = defaultConfig

	c, err := sdk.NewClient(defaultConfig)
	if err != nil {
		return err
	}
	itc.client = c
	itc.ctx = context.Background()

	// TODO(SNOW-1842271): Adjust test setup to work properly with Accountadmin role for object tests and Orgadmin for account tests
	if os.Getenv(string(testenvs.TestAccountCreate)) != "" {
		err = c.Sessions.UseRole(context.Background(), snowflakeroles.Accountadmin)
		if err != nil {
			return err
		}
		defer func() { _ = c.Sessions.UseRole(context.Background(), snowflakeroles.Orgadmin) }()
	}

	// TODO [SNOW-1763603]: we can't use test client because of the testing.T parameter that is not present here; discuss
	itc.testClient = helpers.NewTestClient(c, TestDatabaseName, TestSchemaName, TestWarehouseName, integrationtests.ObjectsSuffix)

	db, dbCleanup, err := createDb(itc.client, itc.ctx, false, itc.testClient)
	itc.databaseCleanup = dbCleanup
	if err != nil {
		return err
	}
	itc.database = db

	sc, scCleanup, err := createSc(itc.client, itc.ctx, itc.database, false)
	itc.schemaCleanup = scCleanup
	if err != nil {
		return err
	}
	itc.schema = sc

	wh, whCleanup, err := createWh(itc.client, itc.ctx, false)
	itc.warehouseCleanup = whCleanup
	if err != nil {
		return err
	}
	itc.warehouse = wh

	// TODO [SNOW-1763603]: improve setup; this is a quick workaround for faster local testing
	if os.Getenv(string(testenvs.SimplifiedIntegrationTestsSetup)) == "" {
		config, err := sdk.ProfileConfig(testprofiles.Secondary)
		if err != nil {
			return err
		}

		if config.Account == defaultConfig.Account {
			log.Println("[WARN] default and secondary configs are set to the same account; it may cause problems in tests requiring multiple accounts")
		}

		secondaryClient, err := sdk.NewClient(config)
		if err != nil {
			return err
		}
		itc.secondaryClient = secondaryClient
		itc.secondaryCtx = context.Background()

		itc.secondaryTestClient = helpers.NewTestClient(secondaryClient, TestDatabaseName, TestSchemaName, TestWarehouseName, integrationtests.ObjectsSuffix)

		secondaryDb, secondaryDbCleanup, err := createDb(itc.secondaryClient, itc.secondaryCtx, true, itc.secondaryTestClient)
		itc.secondaryDatabaseCleanup = secondaryDbCleanup
		if err != nil {
			return err
		}
		itc.secondaryDatabase = secondaryDb

		secondarySchema, secondarySchemaCleanup, err := createSc(itc.secondaryClient, itc.secondaryCtx, itc.database, true)
		itc.secondarySchemaCleanup = secondarySchemaCleanup
		if err != nil {
			return err
		}
		itc.secondarySchema = secondarySchema

		secondaryWarehouse, secondaryWarehouseCleanup, err := createWh(itc.secondaryClient, itc.secondaryCtx, true)
		itc.secondaryWarehouseCleanup = secondaryWarehouseCleanup
		if err != nil {
			return err
		}
		itc.secondaryWarehouse = secondaryWarehouse

		err = helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(itc.client, itc.ctx)
		if err != nil {
			return err
		}
		err = helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(itc.secondaryClient, itc.secondaryCtx)
		if err != nil {
			return err
		}

		// TODO(SNOW-1842271): Adjust test setup to work properly with Accountadmin role for object tests and Orgadmin for account tests
		if os.Getenv(string(testenvs.TestAccountCreate)) == "" {
			err = helpers.EnsureScimProvisionerRolesExist(itc.client, itc.ctx)
			if err != nil {
				return err
			}
			err = helpers.EnsureScimProvisionerRolesExist(itc.secondaryClient, itc.secondaryCtx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createDb(client *sdk.Client, ctx context.Context, ifNotExists bool, testClient *helpers.TestClient) (*sdk.Database, func(), error) {
	id := sdk.NewAccountObjectIdentifier(TestDatabaseName)
	cleanup := func() {
		_ = client.Databases.Drop(ctx, id, &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)})
	}
	opts := testClient.Database.TestParametersSet()
	opts.IfNotExists = sdk.Bool(ifNotExists)
	err := client.Databases.Create(ctx, id, opts)
	if err != nil {
		return nil, cleanup, err
	}
	database, err := client.Databases.ShowByID(ctx, id)
	return database, cleanup, err
}

func createSc(client *sdk.Client, ctx context.Context, db *sdk.Database, ifNotExists bool) (*sdk.Schema, func(), error) {
	id := sdk.NewDatabaseObjectIdentifier(db.Name, TestSchemaName)
	cleanup := func() {
		_ = client.Schemas.Drop(ctx, id, &sdk.DropSchemaOptions{IfExists: sdk.Bool(true)})
	}
	err := client.Schemas.Create(ctx, id, &sdk.CreateSchemaOptions{IfNotExists: sdk.Bool(ifNotExists)})
	if err != nil {
		return nil, cleanup, err
	}
	schema, err := client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(db.Name, TestSchemaName))
	return schema, cleanup, err
}

func createWh(client *sdk.Client, ctx context.Context, ifNotExists bool) (*sdk.Warehouse, func(), error) {
	id := sdk.NewAccountObjectIdentifier(TestWarehouseName)
	cleanup := func() {
		_ = client.Warehouses.Drop(ctx, id, &sdk.DropWarehouseOptions{IfExists: sdk.Bool(true)})
	}
	err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(ifNotExists)})
	if err != nil {
		return nil, cleanup, err
	}
	warehouse, err := client.Warehouses.ShowByID(ctx, id)
	return warehouse, cleanup, err
}

// timer measures time from invocation point to the end of method.
// It's supposed to be used like:
//
//	defer timer("something to measure name")()
func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("[DEBUG] %s took %v", name, time.Since(start))
	}
}

// TODO [SNOW-1653619]: distinguish between testClient and testClientHelper in tests (one is tested, the second helps with the tests, they should be the different ones)
func testClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.client
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	return itc.ctx
}

func testSecondaryClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.secondaryClient
}

func testSecondaryContext(t *testing.T) context.Context {
	t.Helper()
	return itc.secondaryCtx
}

func testConfig(t *testing.T) *gosnowflake.Config {
	t.Helper()
	return itc.config
}

func testClientHelper() *helpers.TestClient {
	return itc.testClient
}

func secondaryTestClientHelper() *helpers.TestClient {
	return itc.secondaryTestClient
}
