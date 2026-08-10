//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountSessionPolicyAttachment_BasicUseCase(t *testing.T) {
	sessionPolicy, sessionPolicyCleanup := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup)
	sessionPolicyName := sessionPolicy.ID().FullyQualifiedName()

	sessionPolicy2, sessionPolicyCleanup2 := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyCleanup2)
	sessionPolicyName2 := sessionPolicy2.ID().FullyQualifiedName()

	sessionPolicyForServiceUsers, sessionPolicyForServiceUsersCleanup := testClient().SessionPolicy.CreateSessionPolicy(t)
	t.Cleanup(sessionPolicyForServiceUsersCleanup)
	sessionPolicyForServiceUsersName := sessionPolicyForServiceUsers.ID().FullyQualifiedName()

	basicModel := model.AccountSessionPolicyAttachment("account", sessionPolicyName)
	personModel := model.AccountSessionPolicyAttachment("person", sessionPolicyName).WithForAllPersonUsers(true)
	serviceModel := model.AccountSessionPolicyAttachment("service", sessionPolicyForServiceUsersName).WithForAllServiceUsers(true)

	basicModelUpdated := model.AccountSessionPolicyAttachment("account", sessionPolicyName2)

	personModelAsService := model.AccountSessionPolicyAttachment("person", sessionPolicyName).WithForAllServiceUsers(true)

	basicRef := basicModel.ResourceReference()
	personRef := personModel.ResourceReference()
	serviceRef := serviceModel.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.AccountSessionPolicyAttachmentResource(t, basicRef).
			HasSessionPolicyName(sessionPolicyName).
			HasForAllPersonUsers(false).
			HasForAllServiceUsers(false),
		resourceassert.AccountSessionPolicyAttachmentResource(t, personRef).
			HasSessionPolicyName(sessionPolicyName).
			HasForAllPersonUsers(true).
			HasForAllServiceUsers(false),
		resourceassert.AccountSessionPolicyAttachmentResource(t, serviceRef).
			HasSessionPolicyName(sessionPolicyForServiceUsersName).
			HasForAllPersonUsers(false).
			HasForAllServiceUsers(true),
	}

	updatedAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.AccountSessionPolicyAttachmentResource(t, basicRef).
			HasSessionPolicyName(sessionPolicyName2).
			HasForAllPersonUsers(false).
			HasForAllServiceUsers(false),
	}, basicAssertions[1:]...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountSessionPolicyAttachmentDestroy(t),
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
						WithUnset(*sdk.NewAccountUnsetRequest().WithSessionPolicyUnset(*sdk.NewAccountSessionPolicyUnsetRequest().WithSessionPolicy(true))))
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
					resourceassert.AccountSessionPolicyAttachmentResource(t, personRef).
						HasSessionPolicyName(sessionPolicyName).
						HasForAllPersonUsers(false).
						HasForAllServiceUsers(true),
				),
			},
		},
	})
}

func TestAcc_AccountSessionPolicyAttachment_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	conflictingModel := model.AccountSessionPolicyAttachment("t", id.FullyQualifiedName()).
		WithForAllPersonUsers(true).
		WithForAllServiceUsers(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckAccountSessionPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, conflictingModel),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
		},
	})
}
