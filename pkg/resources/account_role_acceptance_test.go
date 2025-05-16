//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountRole_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	currentRole := acc.TestClient().Context.CurrentRole(t)

	accountRoleModel := model.AccountRole("role", id.Name())
	accountRoleModelWithComment := model.AccountRole("role", id.Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, accountRoleModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(accountRoleModel.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.comment", ""),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, accountRoleModel),
				ResourceName: accountRoleModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, accountRoleModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(accountRoleModelWithComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(accountRoleModelWithComment.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.comment", comment),
				),
			},
			// import - complete
			{
				Config:       accconfig.FromModels(t, accountRoleModel),
				ResourceName: accountRoleModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, accountRoleModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(accountRoleModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(accountRoleModel.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr(accountRoleModel.ResourceReference(), "show_output.0.comment", ""),
				),
			},
		},
	})
}

func TestAcc_AccountRole_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	currentRole := acc.TestClient().Context.CurrentRole(t)

	accountRoleModel := model.AccountRole("role", id.Name())
	accountRoleModelWithComment := model.AccountRole("role", id.Name()).
		WithComment(comment)
	accountRoleModelNewIdAndNewComment := model.AccountRole("role", newId.Name()).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(accountRoleModelWithComment.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr(accountRoleModelWithComment.ResourceReference(), "show_output.0.comment", comment),
				),
			},
			{
				Config:       accconfig.FromModels(t, accountRoleModel),
				ResourceName: accountRoleModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// rename + comment change
			{
				Config: accconfig.FromModels(t, accountRoleModelNewIdAndNewComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "comment", newComment),

					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.name", newId.Name()),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.is_default", "false"),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.is_current", "false"),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.is_inherited", "false"),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.assigned_to_users", "0"),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.granted_to_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.granted_roles", "0"),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.owner", currentRole.Name()),
					resource.TestCheckResourceAttr(accountRoleModelNewIdAndNewComment.ResourceReference(), "show_output.0.comment", newComment),
				),
			},
		},
	})
}

func TestAcc_AccountRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	accountRoleModelWithComment := model.AccountRole("role", id.Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, accountRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, accountRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_AccountRole_WithQuotedName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`"%s"`, id.Name())
	comment := random.Comment()

	accountRoleModelWithComment := model.AccountRole("role", quotedId).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.AccountRole),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             accconfig.FromModels(t, accountRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, accountRoleModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_account_role.role", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_account_role.role", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_role.role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account_role.role", "id", id.Name()),
				),
			},
		},
	})
}
