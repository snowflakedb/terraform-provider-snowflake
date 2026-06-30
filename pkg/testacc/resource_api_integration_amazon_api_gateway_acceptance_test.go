//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ApiIntegrationAmazonApiGateway_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{awsAllowedPrefix}, awsRoleArn, apiProvider, true)
	withOptionals := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{awsAllowedPrefix}, awsRoleArn, apiProvider, true).
		WithApiBlockedPrefixes([]string{awsBlockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAmazonApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiProviderString(apiProvider).
			HasApiAwsRoleArnString(awsRoleArn).
			HasApiAllowedPrefixes(awsAllowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty().
			HasNoApiKey(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationAwsDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasApiAwsRoleArn(awsRoleArn).
			HasApiKey("").
			HasComment(""),
		objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(awsRoleArn).
			HasAllowedPrefixes(awsAllowedPrefix).
			HasNoBlockedPrefixes().
			HasApiKey("").
			HasComment("").
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty(),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAmazonApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiProviderString(apiProvider).
			HasApiAwsRoleArnString(awsRoleArn).
			HasApiAllowedPrefixes(awsAllowedPrefix).
			HasApiBlockedPrefixes(awsBlockedPrefix).
			HasCommentString(comment).
			HasNoApiKey(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationAwsDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasApiAwsRoleArn(awsRoleArn).
			HasApiKey("").
			HasComment(comment),
		objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(awsRoleArn).
			HasAllowedPrefixes(awsAllowedPrefix).
			HasBlockedPrefixes(awsBlockedPrefix).
			HasApiKey("").
			HasComment(comment).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Import - with optionals
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - external changes
			{
				PreConfig: func() {
					testClient().ApiIntegration.Alter(t, sdk.NewAlterApiIntegrationRequest(id).WithSet(
						*sdk.NewApiIntegrationSetRequest().WithComment(externalComment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Destroy
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
			},
			// Create - with optionals
			{
				PreConfig: func() {
					_, err := testClient().ApiIntegration.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
		},
	})
}

func TestAcc_ApiIntegrationAmazonApiGateway_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)

	comment := random.Comment()
	apiKey := random.AlphanumericN(10)

	allAttributes := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{awsAllowedPrefix}, awsRoleArn, apiProvider, true).
		WithApiBlockedPrefixes([]string{awsBlockedPrefix}).
		WithApiKey(apiKey).
		WithComment(comment)

	ref := allAttributes.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAmazonApiGatewayResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasApiProviderString(apiProvider).
			HasApiAwsRoleArnString(awsRoleArn).
			HasApiAllowedPrefixes(awsAllowedPrefix).
			HasApiBlockedPrefixes(awsBlockedPrefix).
			HasCommentString(comment).
			HasApiKey(apiKey),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationAwsDescribeOutput(t, ref).
			HasApiProvider(apiProvider).
			HasApiAwsRoleArn(awsRoleArn).
			HasApiKeyNotEmpty().
			HasComment(comment),
		objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(awsRoleArn).
			HasAllowedPrefixes(awsAllowedPrefix).
			HasBlockedPrefixes(awsBlockedPrefix).
			HasApiKeyNotEmpty().
			HasComment(comment).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import with all attributes (api_key is not readable after import, set to "")
			{
				Config:                  config.FromModels(t, allAttributes),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func TestAcc_ApiIntegrationAmazonApiGateway_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	invalidProvider := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{awsAllowedPrefix}, awsRoleArn, "INVALID_PROVIDER", true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidProvider),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid api integration aws api provider type: invalid_provider"),
			},
		},
	})
}

func TestAcc_ApiIntegrationAmazonApiGateway_Import_WrongProviderType(t *testing.T) {
	// Create an Azure API integration outside of Terraform to use as the import target.
	azureIntegration, azureCleanup := testClient().ApiIntegration.CreateAzure(t)
	t.Cleanup(azureCleanup)

	awsId := testClient().Ids.RandomAccountObjectIdentifier()
	awsModel := model.ApiIntegrationAmazonApiGateway(
		"t", awsId.Name(),
		[]string{awsAllowedPrefix},
		awsRoleArn,
		string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway),
		true,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			// Attempt to import an Azure API integration via the Amazon API Gateway resource — expects a provider type mismatch error.
			{
				Config:        config.FromModels(t, awsModel),
				ResourceName:  awsModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: azureIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_amazon_api_gateway"),
			},
		},
	})
}

// TestAcc_ApiIntegrationAmazonApiGateway_Import verifies that importing a resource created outside Terraform
// populates the ForceNew field api_provider correctly so that no destroy-before-create plan is produced.
func TestAcc_ApiIntegrationAmazonApiGateway_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)
	comment := random.Comment()

	testModel := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{awsAllowedPrefix}, awsRoleArn, apiProvider, true).
		WithApiBlockedPrefixes([]string{awsBlockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(
						t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: awsBlockedPrefix}}).
							WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(
								sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway,
								awsRoleArn,
							)),
					)
					t.Cleanup(cleanup)
				},
				Config:             config.FromModels(t, testModel),
				ResourceName:       testModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config: config.FromModels(t, testModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// TestAcc_ApiIntegrationAmazonApiGateway_Import_WithApiKey verifies that importing a resource with api_key
// does not trigger a destroy-before-create plan. Because Snowflake does not return api_key, the plan will show
// an in-place update to sync the value into state; subsequent plans should be empty.
func TestAcc_ApiIntegrationAmazonApiGateway_Import_WithApiKey(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)

	apiKey := random.AlphanumericN(10)

	testModel := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{awsAllowedPrefix}, awsRoleArn, apiProvider, true).
		WithApiKey(apiKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(
						t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}, true).
							WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(
								sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway,
								awsRoleArn,
							).WithApiKey(apiKey)),
					)
					t.Cleanup(cleanup)
				},
				Config:             config.FromModels(t, testModel),
				ResourceName:       testModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config: config.FromModels(t, testModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(testModel.ResourceReference(), "api_key", tfjson.ActionUpdate, nil, sdk.String(apiKey)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_ApiIntegrationAmazonApiGateway_ExternalProviderTypeMismatch(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway)

	awsModel := model.ApiIntegrationAmazonApiGateway("t", id.Name(), []string{awsAllowedPrefix}, awsRoleArn, apiProvider, true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAmazonApiGateway),
		Steps: []resource.TestStep{
			// Create AWS resource.
			{
				Config: config.FromModels(t, awsModel),
			},
			// External change: drop the AWS integration and recreate with the same name as Azure API Management.
			// The next read should detect the provider type mismatch and return an error.
			{
				PreConfig: func() {
					testClient().ApiIntegration.DropApiIntegrationFunc(t, id)()
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(
						t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}}, true).
							WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId)),
					)
					t.Cleanup(cleanup)
				},
				Config:      config.FromModels(t, awsModel),
				ExpectError: regexp.MustCompile("could not normalize api_provider value"),
			},
		},
	})
}
