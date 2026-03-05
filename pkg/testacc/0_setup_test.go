package testacc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
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
	config     *gosnowflake.Config
	testClient *helpers.TestClient
	client     *sdk.Client

	secondaryConfig     *gosnowflake.Config
	secondaryTestClient *helpers.TestClient
	secondaryClient     *sdk.Client

	azureConfig     *gosnowflake.Config
	azureTestClient *helpers.TestClient
	azureClient     *sdk.Client

	cleanups []func()

	database  *sdk.Database
	schema    *sdk.Schema
	warehouse *sdk.Warehouse

	secondaryDatabase  *sdk.Database
	secondarySchema    *sdk.Schema
	secondaryWarehouse *sdk.Warehouse

	azureDatabase  *sdk.Database
	azureSchema    *sdk.Schema
	azureWarehouse *sdk.Warehouse
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

// TODO [SNOW-2298294]: extract more convenience methods for reuse
// TODO [SNOW-2298294]: potentially extract test context logic into separate directory
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

	ctx := context.Background()

	if err := atc.initializeSnowflakeEnvironment(
		ctx,
		testprofiles.Default,
		&atc.config,
		&atc.client,
		&atc.testClient,
		&atc.database,
		&atc.schema,
		&atc.warehouse,
	); err != nil {
		return err
	}

	// TODO [SNOW-1763603]: what do we do with SimplifiedIntegrationTestsSetup
	if os.Getenv(string(testenvs.SimplifiedIntegrationTestsSetup)) == "" {
		if err := atc.initializeSnowflakeEnvironment(
			ctx,
			testprofiles.Secondary,
			&atc.secondaryConfig,
			&atc.secondaryClient,
			&atc.secondaryTestClient,
			&atc.secondaryDatabase,
			&atc.secondarySchema,
			&atc.secondaryWarehouse,
		); err != nil {
			return err
		}

		if atc.secondaryConfig.Account == atc.config.Account {
			accTestLog.Printf("[WARN] Default and secondary configs are set to the same account; it may cause problems in tests requiring multiple accounts")
		}

		if errs := errors.Join(
			testClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
			testClient().EnsureEssentialRolesExist(ctx),

			secondaryTestClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
			secondaryTestClient().EnsureEssentialRolesExist(ctx),
		); errs != nil {
			return errs
		}

		// TODO(SNOW-3198924): For now, tests on requiring multiple Snowflake instances on other clouds are done only on non-prod environments
		if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeNonProdEnvironment {
			if errs := errors.Join(
				atc.initializeSnowflakeEnvironment(
					ctx,
					testprofiles.Azure,
					&atc.azureConfig,
					&atc.azureClient,
					&atc.azureTestClient,
					&atc.azureDatabase,
					&atc.azureSchema,
					&atc.azureWarehouse,
				),
				azureTestClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
				azureTestClient().EnsureEssentialRolesExist(ctx),
			); errs != nil {
				return errs
			}
		}
	}

	if err := setUpProvider(); err != nil {
		return fmt.Errorf("cannot set up the provider for the acceptance tests, err: %w", err)
	}

	return nil
}

func (atc *acceptanceTestContext) initializeSnowflakeEnvironment(
	ctx context.Context,
	profile string,
	configField **gosnowflake.Config,
	clientField **sdk.Client,
	testClientField **helpers.TestClient,
	databaseField **sdk.Database,
	schemaField **sdk.Schema,
	warehouseField **sdk.Warehouse,
) error {
	config, client, err := setUpSdkClient(profile, "acceptance")
	if err != nil {
		return err
	}
	*configField = config
	*clientField = client
	*testClientField = helpers.NewTestClient(
		client,
		TestDatabaseName,
		TestSchemaName,
		TestWarehouseName,
		acceptancetests.ObjectsSuffix,
		testenvs.GetSnowflakeEnvironmentWithProdDefault(),
	)

	database, databaseCleanup, err := (*testClientField).CreateTestDatabase(ctx, true)
	if err != nil {
		return err
	}
	atc.cleanups = append(atc.cleanups, databaseCleanup)
	*databaseField = database

	schema, schemaCleanup, err := (*testClientField).CreateTestSchema(ctx, true)
	if err != nil {
		return err
	}
	atc.cleanups = append(atc.cleanups, schemaCleanup)
	*schemaField = schema

	warehouse, warehouseCleanup, err := (*testClientField).CreateTestWarehouse(ctx, true)
	if err != nil {
		return err
	}
	atc.cleanups = append(atc.cleanups, warehouseCleanup)
	*warehouseField = warehouse

	return nil
}

func cleanup() {
	accTestLog.Printf("[INFO] Running acceptance tests cleanup")
	for _, cleanupFunc := range slices.Backward(atc.cleanups) {
		if cleanupFunc != nil {
			cleanupFunc()
		}
	}
}

func testClient() *helpers.TestClient {
	return atc.testClient
}

func secondaryTestClient() *helpers.TestClient {
	return atc.secondaryTestClient
}

func azureTestClient() *helpers.TestClient {
	return atc.azureTestClient
}
