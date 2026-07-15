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
	"github.com/snowflakedb/gosnowflake/v2"
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

type snowflakeTestEnvironmentContext struct {
	config     *gosnowflake.Config
	client     *sdk.Client
	testClient *helpers.TestClient
	database   *sdk.Database
	schema     *sdk.Schema
	warehouse  *sdk.Warehouse
}

type acceptanceTestContext struct {
	defaultTestEnv           snowflakeTestEnvironmentContext
	secondaryTestEnv         snowflakeTestEnvironmentContext
	azureTestEnv             snowflakeTestEnvironmentContext
	snowflakeDefaultsTestEnv snowflakeTestEnvironmentContext

	cleanups []func()
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

	if err := atc.initializeSnowflakeEnvironment(ctx, testprofiles.Default, &atc.defaultTestEnv); err != nil {
		return err
	}

	// TODO [SNOW-1763603]: what do we do with SimplifiedIntegrationTestsSetup
	if os.Getenv(string(testenvs.SimplifiedIntegrationTestsSetup)) == "" {
		if err := atc.initializeSnowflakeEnvironment(ctx, testprofiles.Secondary, &atc.secondaryTestEnv); err != nil {
			return err
		}

		if atc.secondaryTestEnv.config.Account == atc.defaultTestEnv.config.Account {
			accTestLog.Printf("[WARN] Default and secondary configs are set to the same account; it may cause problems in tests requiring multiple accounts")
		}

		if errs := errors.Join(
			testClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
			testClient().EnsureEnableIdentifierFirstLoginIsSetToTrue(ctx),
			testClient().EnsureEssentialRolesExist(ctx),
			secondaryTestClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
			secondaryTestClient().EnsureEnableIdentifierFirstLoginIsSetToTrue(ctx),
			secondaryTestClient().EnsureEssentialRolesExist(ctx),
		); errs != nil {
			return errs
		}

		// TODO(SNOW-3198924): For now, tests requiring multiple Snowflake instances on other clouds are done only on non-prod environment
		if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeNonProdEnvironment {
			if err := atc.initializeSnowflakeEnvironment(ctx, testprofiles.Azure, &atc.azureTestEnv); err != nil {
				return err
			}

			if errs := errors.Join(
				azureTestClient().EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx),
				azureTestClient().EnsureEnableIdentifierFirstLoginIsSetToTrue(ctx),
				azureTestClient().EnsureEssentialRolesExist(ctx),
			); errs != nil {
				return errs
			}

			if err := atc.initializeSnowflakeEnvironment(ctx, testprofiles.SnowflakeDefaults, &atc.snowflakeDefaultsTestEnv); err != nil {
				return err
			}

			// no setup assertions as the snowflake defaults account is expected to have no predefined objects
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
	envCtx *snowflakeTestEnvironmentContext,
) error {
	config, client, err := setUpSdkClient(profile, "acceptance")
	if err != nil {
		return err
	}
	tc := helpers.NewTestClient(
		client,
		TestDatabaseName,
		TestSchemaName,
		TestWarehouseName,
		acceptancetests.ObjectsSuffix,
		testenvs.GetSnowflakeEnvironmentWithProdDefault(),
	)

	database, databaseCleanup, err := tc.CreateTestDatabase(ctx, true)
	if err != nil {
		return err
	}
	atc.cleanups = append(atc.cleanups, databaseCleanup)

	schema, schemaCleanup, err := tc.CreateTestSchema(ctx, true)
	if err != nil {
		return err
	}
	atc.cleanups = append(atc.cleanups, schemaCleanup)

	warehouse, warehouseCleanup, err := tc.CreateTestWarehouse(ctx, true)
	if err != nil {
		return err
	}
	atc.cleanups = append(atc.cleanups, warehouseCleanup)

	*envCtx = snowflakeTestEnvironmentContext{
		config:     config,
		client:     client,
		testClient: tc,
		database:   database,
		schema:     schema,
		warehouse:  warehouse,
	}

	return nil
}

func cleanup() {
	accTestLog.Printf("[INFO] Running acceptance tests cleanup")
	for _, cleanupFunc := range atc.cleanups {
		if cleanupFunc != nil {
			defer cleanupFunc()
		}
	}
}

func testClient() *helpers.TestClient {
	return atc.defaultTestEnv.testClient
}

func secondaryTestClient() *helpers.TestClient {
	return atc.secondaryTestEnv.testClient
}

func azureTestClient() *helpers.TestClient {
	return atc.azureTestEnv.testClient
}

func snowflakeDefaultsTestClient() *helpers.TestClient {
	return atc.snowflakeDefaultsTestEnv.testClient
}
