package testacc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random/acceptancetests"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
)

const AcceptanceTestPrefix = "acc_test_"

var (
	TestDatabaseName  = fmt.Sprintf("%sdb_%s", AcceptanceTestPrefix, acceptancetests.ObjectsSuffix)
	TestSchemaName    = fmt.Sprintf("%ssc_%s", AcceptanceTestPrefix, acceptancetests.ObjectsSuffix)
	TestWarehouseName = fmt.Sprintf("%swh_%s", AcceptanceTestPrefix, acceptancetests.ObjectsSuffix)

	NonExistingAccountObjectIdentifier  = sdk.NewAccountObjectIdentifier("does_not_exist")
	NonExistingDatabaseObjectIdentifier = sdk.NewDatabaseObjectIdentifier(TestDatabaseName, "does_not_exist")
	NonExistingSchemaObjectIdentifier   = sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, "does_not_exist")
)

// TODO [SNOW-1325306]: make logging level configurable
// TODO [SNOW-1325306]: adjust during logger rework (e.g. use in model builders); maybe use log/slog
var accTestLog = log.New(os.Stdout, "", log.LstdFlags)

type acceptanceTestContext struct {
	config              *gosnowflake.Config
	secondaryConfig     *gosnowflake.Config
	testClient          *helpers.TestClient
	secondaryTestClient *helpers.TestClient
	client              *sdk.Client
	secondaryClient     *sdk.Client

	database         *sdk.Database
	databaseCleanup  func()
	schema           *sdk.Schema
	schemaCleanup    func()
	warehouse        *sdk.Warehouse
	warehouseCleanup func()

	secondaryDatabase         *sdk.Database
	secondaryDatabaseCleanup  func()
	secondarySchema           *sdk.Schema
	secondarySchemaCleanup    func()
	secondaryWarehouse        *sdk.Warehouse
	secondaryWarehouseCleanup func()
}

var atc acceptanceTestContext

func TestMain(m *testing.M) {
	exitVal := execute(m)
	os.Exit(exitVal)
}

func execute(m *testing.M) int {
	defer timer("acceptance tests", accTestLog)()
	defer cleanup()
	setup()
	exitVal := m.Run()
	return exitVal
}

func setup() {
	accTestLog.Printf("[INFO] Running acceptance tests setup")

	err := atc.initialize()
	if err != nil {
		accTestLog.Printf("[ERROR] Acceptance test context initialization failed with: `%s`", err)
		cleanup()
		os.Exit(1)
	}
}

// TODO [next PRs]: extract more convenience methods for reuse
// TODO [next PRs]: potentially extract test context logic into separate directory
func (atc *acceptanceTestContext) initialize() error {
	accTestLog.Printf("[INFO] Initializing acceptance test context")

	enableAcceptance := os.Getenv(fmt.Sprintf("%v", testenvs.EnableAcceptance))
	if enableAcceptance == "" {
		return fmt.Errorf("acceptance tests cannot be run; set %s env to run them", testenvs.EnableAcceptance)
	}

	testObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.TestObjectsSuffix))
	requireTestObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.RequireTestObjectsSuffix))
	if requireTestObjectSuffix != "" && testObjectSuffix == "" {
		return fmt.Errorf("test object suffix is required for this test run; set %s env", testenvs.TestObjectsSuffix)
	}

	defaultConfig, client, err := setUpSdkClient(testprofiles.Default, "acceptance")
	if err != nil {
		return err
	}
	atc.config = defaultConfig
	atc.client = client
	atc.testClient = helpers.NewTestClient(client, TestDatabaseName, TestSchemaName, TestWarehouseName, acceptancetests.ObjectsSuffix)

	ctx := context.Background()
	db, dbCleanup, err := testClient().CreateTestDatabase(ctx, false)
	atc.databaseCleanup = dbCleanup
	if err != nil {
		return err
	}
	atc.database = db

	sc, scCleanup, err := testClient().CreateTestSchema(ctx, false)
	atc.schemaCleanup = scCleanup
	if err != nil {
		return err
	}
	atc.schema = sc

	wh, whCleanup, err := testClient().CreateTestWarehouse(ctx, false)
	atc.warehouseCleanup = whCleanup
	if err != nil {
		return err
	}
	atc.warehouse = wh

	// TODO [SNOW-1763603]: what do we do with SimplifiedIntegrationTestsSetup
	if os.Getenv(string(testenvs.SimplifiedIntegrationTestsSetup)) == "" {
		secondaryConfig, secondaryClient, err := setUpSdkClient(testprofiles.Secondary, "acceptance")
		if err != nil {
			return err
		}
		atc.secondaryConfig = secondaryConfig
		atc.secondaryClient = secondaryClient
		atc.secondaryTestClient = helpers.NewTestClient(secondaryClient, TestDatabaseName, TestSchemaName, TestWarehouseName, acceptancetests.ObjectsSuffix)

		if secondaryConfig.Account == defaultConfig.Account {
			accTestLog.Printf("[WARN] Default and secondary configs are set to the same account; it may cause problems in tests requiring multiple accounts")
		}

		secondaryDb, secondaryDbCleanup, err := secondaryTestClient().CreateTestDatabase(ctx, true)
		atc.secondaryDatabaseCleanup = secondaryDbCleanup
		if err != nil {
			return err
		}
		atc.secondaryDatabase = secondaryDb

		secondarySchema, secondarySchemaCleanup, err := secondaryTestClient().CreateTestSchema(ctx, true)
		atc.secondarySchemaCleanup = secondarySchemaCleanup
		if err != nil {
			return err
		}
		atc.secondarySchema = secondarySchema

		secondaryWarehouse, secondaryWarehouseCleanup, err := secondaryTestClient().CreateTestWarehouse(ctx, true)
		atc.secondaryWarehouseCleanup = secondaryWarehouseCleanup
		if err != nil {
			return err
		}
		atc.secondaryWarehouse = secondaryWarehouse

		errs := errors.Join(
			testClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
			secondaryTestClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
			testClient().EnsureScimProvisionerRolesExist(ctx),
			secondaryTestClient().EnsureScimProvisionerRolesExist(ctx),
		)
		if errs != nil {
			return errs
		}
	}

	if err := setUpProvider(); err != nil {
		return fmt.Errorf("cannot set up the provider for the acceptance tests, err: %w", err)
	}

	return nil
}

func cleanup() {
	accTestLog.Printf("[INFO] Running acceptance tests cleanup")
	if atc.databaseCleanup != nil {
		defer atc.databaseCleanup()
	}
	if atc.schemaCleanup != nil {
		defer atc.schemaCleanup()
	}
	if atc.warehouseCleanup != nil {
		defer atc.warehouseCleanup()
	}
	if atc.secondaryDatabaseCleanup != nil {
		defer atc.secondaryDatabaseCleanup()
	}
	if atc.secondarySchemaCleanup != nil {
		defer atc.secondarySchemaCleanup()
	}
	if atc.secondaryWarehouseCleanup != nil {
		defer atc.secondaryWarehouseCleanup()
	}
}

func testClient() *helpers.TestClient {
	return atc.testClient
}

func secondaryTestClient() *helpers.TestClient {
	return atc.secondaryTestClient
}
