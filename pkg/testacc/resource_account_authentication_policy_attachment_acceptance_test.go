//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountAuthenticationPolicyAttachment_BasicUseCase(t *testing.T) {
	testClient().EnsureValidNonProdAccountIsUsed(t)

	policy, policyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(policyCleanup)
	policyName := policy.ID().FullyQualifiedName()

	policy2, policy2Cleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(policy2Cleanup)
	policyName2 := policy2.ID().FullyQualifiedName()

	policyForServiceUsers, policyForServiceUsersCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(policyForServiceUsersCleanup)
	policyForServiceUsersName := policyForServiceUsers.ID().FullyQualifiedName()

	basicModel := model.AccountAuthenticationPolicyAttachment("account", policyName)
	personModel := model.AccountAuthenticationPolicyAttachment("person", policyName).WithForAllPersonUsers(true)
	serviceModel := model.AccountAuthenticationPolicyAttachment("service", policyForServiceUsersName).WithForAllServiceUsers(true)

	basicModelUpdated := model.AccountAuthenticationPolicyAttachment("account", policyName2)

	personModelAsService := model.AccountAuthenticationPolicyAttachment("person", policyName).WithForAllServiceUsers(true)

	basicRef := basicModel.ResourceReference()
	personRef := personModel.ResourceReference()
	serviceRef := serviceModel.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.AccountAuthenticationPolicyAttachmentResource(t, basicRef).
			HasAuthenticationPolicy(policyName).
			HasForAllPersonUsers(false).
			HasForAllServiceUsers(false),
		resourceassert.AccountAuthenticationPolicyAttachmentResource(t, personRef).
			HasAuthenticationPolicy(policyName).
			HasForAllPersonUsers(true).
			HasForAllServiceUsers(false),
		resourceassert.AccountAuthenticationPolicyAttachmentResource(t, serviceRef).
			HasAuthenticationPolicy(policyForServiceUsersName).
			HasForAllPersonUsers(false).
			HasForAllServiceUsers(true),
	}

	updatedAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.AccountAuthenticationPolicyAttachmentResource(t, basicRef).
			HasAuthenticationPolicy(policyName2).
			HasForAllPersonUsers(false).
			HasForAllServiceUsers(false),
	}, basicAssertions[1:]...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountAuthenticationPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// Create
			{
				Config: config.FromModels(t, basicModel, personModel, serviceModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicRef, plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(personRef, plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(serviceRef, plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, basicAssertions...),
			},
			// Import
			{
				ResourceName:      basicRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      personRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      serviceRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change the account-wide attachment's policy
			{
				Config: config.FromModels(t, basicModelUpdated, personModel, serviceModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicRef, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(personRef, plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(serviceRef, plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, updatedAssertions...),
			},
			// Unset a single attachment (account-wide) outside of Terraform
			{
				PreConfig: func() {
					testClient().Account.Alter(t, sdk.NewAlterAccountRequest().
						WithUnset(*sdk.NewAccountUnsetRequest().WithAuthenticationPolicyUnset(*sdk.NewAccountAuthenticationPolicyUnsetRequest().WithAuthenticationPolicy(true))))
				},
				Config: config.FromModels(t, basicModelUpdated, personModel, serviceModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicRef, plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction(personRef, plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(serviceRef, plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t, updatedAssertions...),
			},
			// Remove the service-users attachment, freeing that scope for the switch below
			{
				Config: config.FromModels(t, basicModelUpdated, personModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicRef, plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(personRef, plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(serviceRef, plancheck.ResourceActionDestroy),
					},
				},
			},
			// Switch the person-users attachment to the service-users scope; changing the scope forces recreation
			{
				Config: config.FromModels(t, basicModelUpdated, personModelAsService),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicRef, plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(personRef, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.AccountAuthenticationPolicyAttachmentResource(t, personRef).
						HasAuthenticationPolicy(policyName).
						HasForAllPersonUsers(false).
						HasForAllServiceUsers(true),
				),
			},
		},
	})
}

func TestAcc_AccountAuthenticationPolicyAttachment_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	conflictingModel := model.AccountAuthenticationPolicyAttachment("t", id.FullyQualifiedName()).
		WithForAllPersonUsers(true).
		WithForAllServiceUsers(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountAuthenticationPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, conflictingModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
		},
	})
}

// TestAcc_AccountAuthenticationPolicyAttachment_migrateFromV2_18_0_ensureSmoothUpgradeWithNewResourceId verifies that
// the change of the resource id format does not cause a diff after upgrading from the last released version.
func TestAcc_AccountAuthenticationPolicyAttachment_migrateFromV2_18_0_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	testClient().EnsureValidNonProdAccountIsUsed(t)

	policy, policyCleanup := testClient().AuthenticationPolicy.Create(t)
	t.Cleanup(policyCleanup)
	policyName := policy.ID().FullyQualifiedName()

	basicModel := model.AccountAuthenticationPolicyAttachment("t", policyName)
	ref := basicModel.ResourceReference()

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountAuthenticationPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// Create with the last released version that still uses the old resource id format (helpers.EncodeSnowflakeID).
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.18.0"),
				Config:            config.FromModels(t, providerModel, basicModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", helpers.EncodeSnowflakeID(policy.ID()))),
				),
			},
			// Upgrade to the current version that additionally encodes the target scope in the resource id and ensure the
			// upgrade produces no plan. An account-wide attachment (no FOR ALL clause) is encoded with the ACCOUNT scope.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", helpers.EncodeResourceIdentifier(policy.ID().FullyQualifiedName(), string(sdk.AuthenticationPolicyTargetScopeAccount)))),
				),
			},
		},
	})
}
