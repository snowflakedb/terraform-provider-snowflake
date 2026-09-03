package testenvs

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type env string

const (
	envPrefix                   = "TEST_SF_TF_"
	BusinessCriticalAccount env = "SNOWFLAKE_BUSINESS_CRITICAL_ACCOUNT"

	// TestNonProdModifiableAccountLocator represents an account locator that can be used for account modification tests.
	TestNonProdModifiableAccountLocator env = "TEST_SF_TF_NON_PROD_MODIFIABLE_ACCOUNT_LOCATOR"
	TestAccountCreate                   env = "TEST_SF_TF_TEST_ACCOUNT_CREATE"
	TestFailoverGroups                  env = "TEST_SF_TF_TEST_FAILOVER_GROUPS"
	// TestOpenflow gates tests needing an account with the Openflow object model enabled and configured.
	TestOpenflow env = envPrefix + "TEST_OPENFLOW"

	AwsExternalBucketUrl   env = "TEST_SF_TF_AWS_EXTERNAL_BUCKET_URL"
	AwsExternalKeyId       env = "TEST_SF_TF_AWS_EXTERNAL_KEY_ID"
	AwsExternalSecretKey   env = "TEST_SF_TF_AWS_EXTERNAL_SECRET_KEY" // #nosec G101
	AwsExternalRoleArn     env = "TEST_SF_TF_AWS_EXTERNAL_ROLE_ARN"
	AwsExternalSnsTopicArn env = "TEST_SF_TF_AWS_EXTERNAL_SNS_TOPIC_ARN"
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
	EnableAllPreviewFeatures    env = "SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES"
	TestObjectsSuffix           env = "TEST_SF_TF_TEST_OBJECT_SUFFIX"
	RequireTestObjectsSuffix    env = "TEST_SF_TF_REQUIRE_TEST_OBJECT_SUFFIX"
	RequireGeneratedRandomValue env = "TEST_SF_TF_REQUIRE_GENERATED_RANDOM_VALUE"
	GeneratedRandomValue        env = "TEST_SF_TF_GENERATED_RANDOM_VALUE"
	SnowflakeTestingEnvironment env = "TEST_SF_TF_SNOWFLAKE_TESTING_ENVIRONMENT"

	SimplifiedIntegrationTestsSetup env = "TEST_SF_TF_SIMPLIFIED_INTEGRATION_TESTS_SETUP"

	TestResourceNullListHandlingEnv     env = "TEST_SF_TF_TEST_RESOURCE_DATA_NULL_LIST_HANDLING_ENV"
	TestResourceDataTypeDiffHandlingEnv env = "TEST_SF_TF_TEST_RESOURCE_DATA_DIFF_HANDLING_ENV"

	// Oauth-related
	OauthWithClientCredentialsClientId     env = envPrefix + "OAUTH_WITH_CLIENT_CREDENTIALS_CLIENT_ID"
	OauthWithClientCredentialsClientSecret env = envPrefix + "OAUTH_WITH_CLIENT_CREDENTIALS_CLIENT_SECRET"
	OauthWithClientCredentialsIssuer       env = envPrefix + "OAUTH_WITH_CLIENT_CREDENTIALS_ISSUER"

	// OpenflowDeployment names a pre-provisioned, ACTIVE SNOWFLAKE Openflow deployment for the
	// runtime and connector tests to build on. Those objects can only be created inside a deployment that
	// has reached ACTIVE, and provisioning one takes minutes - so the tests reuse an existing deployment
	// rather than creating and tearing one down per run. The name is required and the run fails without
	// it: an account can hold several ACTIVE deployments belonging to other people, and picking one
	// automatically would create and destroy objects inside it.
	OpenflowDeployment env = envPrefix + "OPENFLOW_DEPLOYMENT"

	// OpenflowRuntime names a pre-provisioned, ACTIVE Openflow runtime by fully qualified name, for the
	// subtests that only need a runtime to exist. Creating one takes four to five minutes, so it is
	// provisioned out of band like the deployment above rather than per run.
	//
	// It carries a display name, comment and one integration so that the subtests reading it see every
	// optional column populated. Provision it with:
	//
	//	CREATE DATABASE IF NOT EXISTS TF_PROVIDER_TESTS;
	//	CREATE SCHEMA IF NOT EXISTS TF_PROVIDER_TESTS.OPENFLOW;
	//
	//	-- The runtime needs an integration attached so external_access_integrations is non-empty. An egress
	//	-- rule with no values reaches nothing, which is all the tests require of it.
	//	CREATE NETWORK RULE TF_PROVIDER_TESTS.OPENFLOW.TF_PROVIDER_TESTS_OPENFLOW_RULE
	//	  TYPE = HOST_PORT MODE = EGRESS VALUE_LIST = ();
	//	CREATE EXTERNAL ACCESS INTEGRATION TF_PROVIDER_TESTS_OPENFLOW_EAI
	//	  ALLOWED_NETWORK_RULES = (TF_PROVIDER_TESTS.OPENFLOW.TF_PROVIDER_TESTS_OPENFLOW_RULE) ENABLED = TRUE;
	//
	//	CREATE OPENFLOW RUNTIME TF_PROVIDER_TESTS.OPENFLOW.TF_PROVIDER_TESTS_RUNTIME
	//	  IN DEPLOYMENT <the deployment named by TEST_SF_TF_OPENFLOW_DEPLOYMENT>
	//	  EXECUTE_AS_ROLE = <the role the tests connect as>
	//	  NODE_TYPE = 'SMALL' MIN_NODES = 1 MAX_NODES = 1
	//	  EXTERNAL_ACCESS_INTEGRATIONS = (TF_PROVIDER_TESTS_OPENFLOW_EAI)
	//	  DISPLAY_NAME = 'TF Provider Tests Runtime'
	//	  COMMENT = 'Shared fixture for the Terraform provider Openflow integration tests';
	//
	// The execute-as role must hold usage on the integration, which it does when the same role creates both.
	OpenflowRuntime env = envPrefix + "OPENFLOW_RUNTIME"

	OpenCatalogAccountLocator             env = envPrefix + "OPEN_CATALOG_ACCOUNT_LOCATOR"
	OpenCatalogPrimaryOAuthClientId       env = envPrefix + "OPEN_CATALOG_PRIMARY_OAUTH_CLIENT_ID"
	OpenCatalogPrimaryOAuthClientSecret   env = envPrefix + "OPEN_CATALOG_PRIMARY_OAUTH_CLIENT_SECRET"
	OpenCatalogSecondaryOAuthClientId     env = envPrefix + "OPEN_CATALOG_SECONDARY_OAUTH_CLIENT_ID"
	OpenCatalogSecondaryOAuthClientSecret env = envPrefix + "OPEN_CATALOG_SECONDARY_OAUTH_CLIENT_SECRET"
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

func SkipTestIfSetTo(t *testing.T, envName Env, value string, reason string) {
	t.Helper()
	env := os.Getenv(fmt.Sprintf("%v", envName))
	if env == value {
		t.Skipf("Skipping %s, because env %v is set to %s. Reason: \"%s\"", t.Name(), envName, value, reason)
	}
}

func SkipTestIfValueIn(t *testing.T, envName Env, values []string, reason string) {
	t.Helper()
	env := os.Getenv(fmt.Sprintf("%v", envName))
	if slices.Contains(values, env) {
		t.Skipf("Skipping %s, because env %v is set to %s. Reason: \"%s\"", t.Name(), envName, env, reason)
	}
}

type Env interface {
	xxxProtected()
}

func (e env) xxxProtected() {}
