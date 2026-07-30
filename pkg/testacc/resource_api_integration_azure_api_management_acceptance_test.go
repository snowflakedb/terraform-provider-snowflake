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

func TestAcc_ApiIntegrationAzureApiManagement_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.ApiIntegrationAzureApiManagement("t", id.Name(), []string{azureAllowedPrefix}, azureAdApplicationId, azureTenantId, true)
	withOptionals := model.ApiIntegrationAzureApiManagement("t", id.Name(), []string{azureAllowedPrefix}, azureAdApplicationId, azureTenantId, true).
		WithApiBlockedPrefixes([]string{azureBlockedPrefix}).
		WithComment(comment)

	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAzureApiManagementResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasAzureTenantIdString(azureTenantId).
			HasAzureAdApplicationIdString(azureAdApplicationId).
			HasApiAllowedPrefixes(azureAllowedPrefix).
			HasApiBlockedPrefixesEmpty().
			HasCommentEmpty().
			HasNoApiKey(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(""),
		resourceshowoutputassert.ApiIntegrationAzureDescribeOutput(t, ref).
			HasApiProvider(string(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement)).
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasApiKey("").
			HasComment(""),
		objectassert.ApiIntegrationAzureDetails(t, id).
			HasEnabled(true).
			HasApiProviderType(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement).
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasAllowedPrefixes(azureAllowedPrefix).
			HasNoBlockedPrefixes().
			HasApiKey("").
			HasComment("").
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty(),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAzureApiManagementResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasAzureTenantIdString(azureTenantId).
			HasAzureAdApplicationIdString(azureAdApplicationId).
			HasApiAllowedPrefixes(azureAllowedPrefix).
			HasApiBlockedPrefixes(azureBlockedPrefix).
			HasCommentString(comment).
			HasNoApiKey(),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationAzureDescribeOutput(t, ref).
			HasApiProvider(string(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement)).
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasApiKey("").
			HasComment(comment),
		objectassert.ApiIntegrationAzureDetails(t, id).
			HasEnabled(true).
			HasApiProviderType(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement).
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasAllowedPrefixes(azureAllowedPrefix).
			HasBlockedPrefixes(azureBlockedPrefix).
			HasApiKey("").
			HasComment(comment).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAzureApiManagement),
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

func TestAcc_ApiIntegrationAzureApiManagement_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	apiKey := random.AlphanumericN(10)

	allAttributes := model.ApiIntegrationAzureApiManagement("t", id.Name(), []string{azureAllowedPrefix}, azureAdApplicationId, azureTenantId, true).
		WithApiBlockedPrefixes([]string{azureBlockedPrefix}).
		WithApiKey(apiKey).
		WithComment(comment)

	ref := allAttributes.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationAzureApiManagementResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasAzureTenantIdString(azureTenantId).
			HasAzureAdApplicationIdString(azureAdApplicationId).
			HasApiAllowedPrefixes(azureAllowedPrefix).
			HasApiBlockedPrefixes(azureBlockedPrefix).
			HasCommentString(comment).
			HasApiKey(apiKey),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationAzureDescribeOutput(t, ref).
			HasApiProvider(string(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement)).
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasApiKeyNotEmpty().
			HasComment(comment),
		objectassert.ApiIntegrationAzureDetails(t, id).
			HasEnabled(true).
			HasApiProviderType(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement).
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasAllowedPrefixes(azureAllowedPrefix).
			HasBlockedPrefixes(azureBlockedPrefix).
			HasApiKeyNotEmpty().
			HasComment(comment).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAzureApiManagement),
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
			// Import with all attributes (api_key is write-only: Snowflake returns a masked value and the provider
			// cannot detect external changes to it; api_key is excluded from import state verification)
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

func TestAcc_ApiIntegrationAzureApiManagement_Import_WrongApiProvider(t *testing.T) {
	azureId := testClient().Ids.RandomAccountObjectIdentifier()

	// Create an AWS API integration outside Terraform to use as the import target.
	awsIntegration, awsCleanup := testClient().ApiIntegration.CreateAws(t)
	t.Cleanup(awsCleanup)

	azureModel := model.ApiIntegrationAzureApiManagement(
		"t", azureId.Name(),
		[]string{azureAllowedPrefix},
		azureAdApplicationId,
		azureTenantId,
		true,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAzureApiManagement),
		Steps: []resource.TestStep{
			// Create a valid Azure integration to have a resource in state.
			{
				Config: config.FromModels(t, azureModel),
			},
			// Attempt to import an AWS integration via the Azure API Management resource — expects a type mismatch error.
			{
				ResourceName:  azureModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: awsIntegration.Name,
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_azure_api_management"),
			},
		},
	})
}

// TestAcc_ApiIntegrationAzureApiManagement_Import verifies that importing a resource created outside Terraform
// populates state correctly so that no destroy-before-create plan is produced.
func TestAcc_ApiIntegrationAzureApiManagement_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()

	testModel := model.ApiIntegrationAzureApiManagement("t", id.Name(), []string{azureAllowedPrefix}, azureAdApplicationId, azureTenantId, true).
		WithApiBlockedPrefixes([]string{azureBlockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAzureApiManagement),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(
						t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: azureBlockedPrefix}}).
							WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(
								azureTenantId,
								azureAdApplicationId,
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

// TestAcc_ApiIntegrationAzureApiManagement_Import_WithApiKey verifies that importing a resource with api_key
// does not trigger a destroy-before-create plan. Because Snowflake does not return api_key, the plan will show
// an in-place update to sync the value into state; subsequent plans should be empty.
func TestAcc_ApiIntegrationAzureApiManagement_Import_WithApiKey(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiKey := random.AlphanumericN(10)

	testModel := model.ApiIntegrationAzureApiManagement("t", id.Name(), []string{azureAllowedPrefix}, azureAdApplicationId, azureTenantId, true).
		WithApiKey(apiKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationAzureApiManagement),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(
						t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}}, true).
							WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(
								azureTenantId,
								azureAdApplicationId,
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
