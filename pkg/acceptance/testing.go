package acceptance

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random/acceptancetests"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
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

var (
	TestAccProvider *schema.Provider
	v5Server        tfprotov5.ProviderServer
	v6Server        tfprotov6.ProviderServer
	atc             acceptanceTestContext

	configureClientErrorDiag diag.Diagnostics
	configureProviderCtx     *internalprovider.Context
)

func init() {
	testObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.TestObjectsSuffix))
	requireTestObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.RequireTestObjectsSuffix))
	if requireTestObjectSuffix != "" && testObjectSuffix == "" {
		log.Printf("Test object suffix is required for this test run. Set %s env.", testenvs.TestObjectsSuffix)
		os.Exit(1)
	}

	TestAccProvider = provider.Provider()
	// TODO [SNOW-2054208]: improve during the package extraction
	TestAccProvider.ResourcesMap["snowflake_object_renaming"] = resources.ObjectRenamingListsAndSets()
	TestAccProvider.ResourcesMap["snowflake_test_resource_data_type_diff_handling"] = resources.TestResourceDataTypeDiffHandling()
	TestAccProvider.ResourcesMap["snowflake_test_resource_data_type_diff_handling_list"] = resources.TestResourceDataTypeDiffHandlingList()
	TestAccProvider.ConfigureContextFunc = ConfigureProviderWithConfigCache

	v5Server = TestAccProvider.GRPCProvider()
	var err error
	v6Server, err = tf5to6server.UpgradeServer(
		context.Background(),
		func() tfprotov5.ProviderServer {
			return v5Server
		},
	)
	if err != nil {
		log.Panicf("Cannot upgrade server from proto v5 to proto v6, failing, err: %v", err)
	}
	_ = testAccProtoV6ProviderFactoriesNew

	// TODO [SNOW-2054383]: Use the new TOML format.
	defaultConfig, err := sdk.ProfileConfig(testprofiles.Default, sdk.WithUseLegacyTomlFormat(true))
	if err != nil {
		log.Panicf("Could not read configuration from profile: %v", err)
	}
	if defaultConfig == nil {
		log.Panic("Config is required to run acceptance tests")
	}
	atc.config = defaultConfig

	client, err := sdk.NewClient(defaultConfig)
	if err != nil {
		log.Panicf("Cannot instantiate new client, err: %v", err)
	}
	atc.client = client

	// TODO [SNOW-2054383]: Use the new TOML format.
	cfg, err := sdk.ProfileConfig(testprofiles.Secondary, sdk.WithUseLegacyTomlFormat(true))
	if err != nil {
		log.Panicf("Config for the secondary client is needed to run acceptance tests, err: %v", err)
	}
	secondaryClient, err := sdk.NewClient(cfg)
	if err != nil {
		log.Panicf("Cannot instantiate new secondary client, err: %v", err)
	}
	atc.secondaryClient = secondaryClient

	atc.testClient = helpers.NewTestClient(client, TestDatabaseName, TestSchemaName, TestWarehouseName, acceptancetests.ObjectsSuffix)
	atc.secondaryTestClient = helpers.NewTestClient(secondaryClient, TestDatabaseName, TestSchemaName, TestWarehouseName, acceptancetests.ObjectsSuffix)
}

type acceptanceTestContext struct {
	config              *gosnowflake.Config
	client              *sdk.Client
	secondaryClient     *sdk.Client
	testClient          *helpers.TestClient
	secondaryTestClient *helpers.TestClient
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return v6Server, nil
	},
}

// if we do not reuse the created objects there is no `Previously configured provider being re-configured.` warning
// currently left for possible usage after other improvements
var testAccProtoV6ProviderFactoriesNew = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return tf5to6server.UpgradeServer(
			context.Background(),
			provider.Provider().GRPCProvider,
		)
	},
}

func ConfigureProviderWithConfigCache(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	accTestEnabled, err := oswrapper.GetenvBool("TF_ACC")
	if err != nil {
		accTestEnabled = false
		log.Printf("TF_ACC environmental variable has incorrect format: %v, using %v as a default value", err, accTestEnabled)
	}
	configureClientOnceEnabled, err := oswrapper.GetenvBool("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE")
	if err != nil {
		configureClientOnceEnabled = false
		log.Printf("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE environmental variable has incorrect format: %v, using %v as a default value", err, configureClientOnceEnabled)
	}
	// hacky way to speed up our acceptance tests
	if accTestEnabled && configureClientOnceEnabled {
		log.Printf("[DEBUG] Returning cached provider configuration result")
		if configureProviderCtx != nil {
			log.Printf("[DEBUG] Returning cached provider configuration context")
			return configureProviderCtx, nil
		}
		if configureClientErrorDiag.HasError() {
			log.Printf("[DEBUG] Returning cached provider configuration error")
			return nil, configureClientErrorDiag
		}
	}
	log.Printf("[DEBUG] No cached provider configuration found or caching is not enabled; configuring a new provider")

	providerCtx, clientErrorDiag := provider.ConfigureProvider(ctx, d)

	if providerCtx != nil && accTestEnabled && oswrapper.Getenv("SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES") == "true" {
		providerCtx.(*internalprovider.Context).EnabledFeatures = previewfeatures.AllPreviewFeatures
	}

	// needed for tests verifying different provider setups
	if accTestEnabled && configureClientOnceEnabled {
		configureProviderCtx = providerCtx.(*internalprovider.Context)
		configureClientErrorDiag = clientErrorDiag
	} else {
		configureProviderCtx = nil
		configureClientErrorDiag = make(diag.Diagnostics, 0)
	}

	if clientErrorDiag.HasError() {
		return nil, clientErrorDiag
	}

	return providerCtx, nil
}

var once sync.Once

func TestAccPreCheck(t *testing.T) {
	// use singleton design pattern to ensure we only create these resources once
	// there is no cleanup currently, sweepers take care of it
	once.Do(func() {
		ctx := context.Background()

		_, _ = TestClient().Database.CreateTestDatabaseIfNotExists(t)
		_, _ = TestClient().Schema.CreateTestSchemaIfNotExists(t)
		_, _ = TestClient().Warehouse.CreateTestWarehouseIfNotExists(t)

		_, _ = SecondaryTestClient().Database.CreateTestDatabaseIfNotExists(t)
		_, _ = SecondaryTestClient().Schema.CreateTestSchemaIfNotExists(t)
		_, _ = SecondaryTestClient().Warehouse.CreateTestWarehouseIfNotExists(t)

		if err := helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(atc.client, ctx); err != nil {
			t.Fatal(err)
		}

		if err := helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(atc.secondaryClient, ctx); err != nil {
			t.Fatal(err)
		}

		if err := helpers.EnsureScimProvisionerRolesExist(atc.client, ctx); err != nil {
			t.Fatal(err)
		}

		if err := helpers.EnsureScimProvisionerRolesExist(atc.secondaryClient, ctx); err != nil {
			t.Fatal(err)
		}
	})
}

// ConfigurationSameAsStepN should be used to obtain configuration for one of the previous steps to avoid duplication of configuration and var files.
// Based on config.TestStepDirectory.
func ConfigurationSameAsStepN(step int) func(config.TestStepConfigRequest) string {
	return func(req config.TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, strconv.Itoa(step))
	}
}

// ConfigurationDirectory should be used to obtain configuration if the same can be shared between multiple tests to avoid duplication of configuration and var files.
// Based on config.TestNameDirectory. Similar to config.StaticDirectory but prefixed provided directory with `testdata`.
func ConfigurationDirectory(directory string) func(config.TestStepConfigRequest) string {
	return func(req config.TestStepConfigRequest) string {
		return filepath.Join("testdata", directory)
	}
}

func DefaultConfig(t *testing.T) *gosnowflake.Config {
	t.Helper()
	return atc.config
}

func TestClient() *helpers.TestClient {
	return atc.testClient
}

func SecondaryTestClient() *helpers.TestClient {
	return atc.secondaryTestClient
}

// ExternalProviderWithExactVersion returns a map of external providers with an exact version constraint
func ExternalProviderWithExactVersion(version string) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"snowflake": {
			VersionConstraint: fmt.Sprintf("=%s", version),
			Source:            "snowflakedb/snowflake",
		},
	}
}

// SetV097CompatibleConfigPathEnv sets a new config path in a relevant env variable for a file that is compatible with v0.97.
func SetV097CompatibleConfigPathEnv(t *testing.T) {
	t.Helper()
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	configPath := filepath.Join(home, ".snowflake", "config_v097_compatible")
	t.Setenv(snowflakeenvs.ConfigPath, configPath)
}

// UnsetConfigPathEnv unsets a config path env
func UnsetConfigPathEnv(t *testing.T) {
	t.Helper()
	t.Setenv(snowflakeenvs.ConfigPath, "")
}
