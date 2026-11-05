package testint

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random/integrationtests"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake/v2"
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
	defer cleanup()
	setup()
	exitVal := m.Run()
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
		err = c.Sessions.UseRole(context.Background(), sdk.NewUseRoleSessionRequest(snowflakeroles.Accountadmin))
		if err != nil {
			return err
		}
		defer func() {
			_ = c.Sessions.UseRole(context.Background(), sdk.NewUseRoleSessionRequest(snowflakeroles.Orgadmin))
		}()
	}

	// TODO [SNOW-1763603]: we can't use test client because of the testing.T parameter that is not present here; discuss
	itc.testClient = helpers.NewTestClient(
		c,
		TestDatabaseName,
		TestSchemaName,
		TestWarehouseName,
		integrationtests.ObjectsSuffix,
		testenvs.GetSnowflakeEnvironmentWithProdDefault(),
	)

	// Cleanups are registered only after a successful create. Otherwise a create that fails because
	// the object already exists (e.g. created by a concurrent run with the same suffix) would drop that object.
	db, dbCleanup, err := testClientHelper().CreateTestDatabase(itc.ctx, false)
	if err != nil {
		return err
	}
	itc.databaseCleanup = dbCleanup
	itc.database = db

	sc, scCleanup, err := testClientHelper().CreateTestSchema(itc.ctx, false)
	if err != nil {
		return err
	}
	itc.schemaCleanup = scCleanup
	itc.schema = sc

	wh, whCleanup, err := testClientHelper().CreateTestWarehouse(itc.ctx, false)
	if err != nil {
		return err
	}
	itc.warehouseCleanup = whCleanup
	itc.warehouse = wh

	// TODO [SNOW-1763603]: improve setup; this is a quick workaround for faster local testing
	if os.Getenv(string(testenvs.SimplifiedIntegrationTestsSetup)) == "" {
		if err := logCurrentSnowflakeVersion(itc.ctx, itc.testClient, testprofiles.Default); err != nil {
			return err
		}

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

		itc.secondaryTestClient = helpers.NewTestClient(
			secondaryClient,
			TestDatabaseName,
			TestSchemaName,
			TestWarehouseName,
			integrationtests.ObjectsSuffix,
			testenvs.GetSnowflakeEnvironmentWithProdDefault(),
		)
		if err := logCurrentSnowflakeVersion(itc.secondaryCtx, itc.secondaryTestClient, testprofiles.Secondary); err != nil {
			return err
		}

		secondaryDb, secondaryDbCleanup, err := secondaryTestClientHelper().CreateTestDatabase(itc.ctx, false)
		if err != nil {
			return err
		}
		itc.secondaryDatabaseCleanup = secondaryDbCleanup
		itc.secondaryDatabase = secondaryDb

		secondarySchema, secondarySchemaCleanup, err := secondaryTestClientHelper().CreateTestSchema(itc.ctx, false)
		if err != nil {
			return err
		}
		itc.secondarySchemaCleanup = secondarySchemaCleanup
		itc.secondarySchema = secondarySchema

		secondaryWarehouse, secondaryWarehouseCleanup, err := secondaryTestClientHelper().CreateTestWarehouse(itc.ctx, false)
		if err != nil {
			return err
		}
		itc.secondaryWarehouseCleanup = secondaryWarehouseCleanup
		itc.secondaryWarehouse = secondaryWarehouse

		err = errors.Join(
			testClientHelper().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(itc.ctx),
			testClientHelper().EnsureEnableIdentifierFirstLoginIsSetToTrue(itc.ctx),
			testClientHelper().EnsureTestUserProperties(itc.ctx),
			secondaryTestClientHelper().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(itc.secondaryCtx),
			secondaryTestClientHelper().EnsureEnableIdentifierFirstLoginIsSetToTrue(itc.secondaryCtx),
			secondaryTestClientHelper().EnsureTestUserProperties(itc.secondaryCtx),
		)
		if err != nil {
			return err
		}

		// TODO(SNOW-1842271): Adjust test setup to work properly with Accountadmin role for object tests and Orgadmin for account tests
		if os.Getenv(string(testenvs.TestAccountCreate)) == "" {
			err = testClientHelper().EnsureEssentialRolesExist(itc.ctx)
			if err != nil {
				return fmt.Errorf("ensuring essential roles exist in the test client: %w", err)
			}
			err = secondaryTestClientHelper().EnsureEssentialRolesExist(itc.secondaryCtx)
			if err != nil {
				return fmt.Errorf("ensuring essential roles exist in the secondary test client: %w", err)
			}
		}

		err = testClientHelper().EnsureImageRepositoryExist(itc.ctx)
		if err != nil {
			return err
		}
		err = secondaryTestClientHelper().EnsureImageRepositoryExist(itc.secondaryCtx)
		if err != nil {
			return err
		}
	}

	// Openflow builds on a deployment and a runtime provisioned out of band, so these assert the account is
	// set up rather than creating anything. No Openflow test uses the secondary account.
	if os.Getenv(string(testenvs.TestOpenflow)) != "" {
		err = errors.Join(
			testClientHelper().EnsureOpenflowPostgresCdcDefinitionExists(itc.ctx),
			testClientHelper().EnsureOpenflowDeploymentIsProvisioned(itc.ctx),
			testClientHelper().EnsureOpenflowRuntimeIsProvisioned(itc.ctx),
		)
		if err != nil {
			return err
		}
	}

	return nil
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

func logCurrentSnowflakeVersion(ctx context.Context, tc *helpers.TestClient, profile string) error {
	version, err := tc.GetCurrentSnowflakeVersion(ctx)
	if err != nil {
		return fmt.Errorf("getting current Snowflake version for profile %s: %w", profile, err)
	}
	log.Printf("[INFO] Current Snowflake version for profile %s: %s", profile, version)
	return nil
}

func testClientHelper() *helpers.TestClient {
	return itc.testClient
}

func secondaryTestClientHelper() *helpers.TestClient {
	return itc.secondaryTestClient
}
