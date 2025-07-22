package testenvs

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type env string

const (
	BusinessCriticalAccount env = "SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"

	// TestNonProdModifiableAccountLocator represents an account locator that can be used for account modification tests.
	TestNonProdModifiableAccountLocator env = "TEST_SF_TF_NON_PROD_MODIFIABLE_ACCOUNT_LOCATOR"
	TestAccountCreate                   env = "TEST_SF_TF_TEST_ACCOUNT_CREATE"
	TestFailoverGroups                  env = "TEST_SF_TF_TEST_FAILOVER_GROUPS"

	AwsExternalBucketUrl   env = "TEST_SF_TF_AWS_EXTERNAL_BUCKET_URL"
	AwsExternalKeyId       env = "TEST_SF_TF_AWS_EXTERNAL_KEY_ID"
	AwsExternalSecretKey   env = "TEST_SF_TF_AWS_EXTERNAL_SECRET_KEY" // #nosec G101
	AwsExternalRoleArn     env = "TEST_SF_TF_AWS_EXTERNAL_ROLE_ARN"
	AzureExternalBucketUrl env = "TEST_SF_TF_AZURE_EXTERNAL_BUCKET_URL"
	AzureExternalTenantId  env = "TEST_SF_TF_AZURE_EXTERNAL_TENANT_ID"
	AzureExternalSasToken  env = "TEST_SF_TF_AZURE_EXTERNAL_SAS_TOKEN" // #nosec G101
	GcsExternalBucketUrl   env = "TEST_SF_TF_GCS_EXTERNAL_BUCKET_URL"

	EnableObjectRenamingTest env = "TEST_SF_TF_ENABLE_OBJECT_RENAMING"
	SkipManagedAccountTest   env = "TEST_SF_TF_SKIP_MANAGED_ACCOUNT_TEST"
	SkipSamlIntegrationTest  env = "TEST_SF_TF_SKIP_SAML_INTEGRATION_TEST"

	EnableAcceptance            env = resource.EnvTfAcc
	EnableSweep                 env = "TEST_SF_TF_ENABLE_SWEEP"
	EnableManual                env = "TEST_SF_TF_ENABLE_MANUAL_TESTS"
	ConfigureClientOnce         env = "SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE"
	EnableAllPreviewFeatures    env = "SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES"
	TestObjectsSuffix           env = "TEST_SF_TF_TEST_OBJECT_SUFFIX"
	RequireTestObjectsSuffix    env = "TEST_SF_TF_REQUIRE_TEST_OBJECT_SUFFIX"
	RequireGeneratedRandomValue env = "TEST_SF_TF_REQUIRE_GENERATED_RANDOM_VALUE"
	GeneratedRandomValue        env = "TEST_SF_TF_GENERATED_RANDOM_VALUE"

	SimplifiedIntegrationTestsSetup env = "TEST_SF_TF_SIMPLIFIED_INTEGRATION_TESTS_SETUP"

	TestResourceDataTypeDiffHandlingEnv env = "TEST_SF_TF_TEST_RESOURCE_DATA_DIFF_HANDLING_ENV"
)

func GetOrSkipTest(t *testing.T, envName Env) string {
	t.Helper()
	env := os.Getenv(fmt.Sprintf("%v", envName))
	if env == "" {
		t.Skipf("Skipping %s, env %v missing", t.Name(), envName)
	}
	return env
}

func SkipTestIfSet(t *testing.T, envName Env, reason string) {
	t.Helper()
	env := os.Getenv(fmt.Sprintf("%v", envName))
	if env != "" {
		t.Skipf("Skipping %s, because env %v is set. Reason: \"%s\"", t.Name(), envName, reason)
	}
}

type Env interface {
	xxxProtected()
}

func (e env) xxxProtected() {}
