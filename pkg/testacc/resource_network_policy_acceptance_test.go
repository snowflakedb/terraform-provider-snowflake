//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	allowedNetworkRule1, allowedNetworkRule1Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule1Cleanup)

	allowedNetworkRule2, allowedNetworkRule2Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule2Cleanup)

	blockedNetworkRule1, blockedNetworkRule1Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule1Cleanup)

	blockedNetworkRule2, blockedNetworkRule2Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule2Cleanup)

	allowedNetworkRuleId1 := allowedNetworkRule1.ID()
	allowedNetworkRuleId2 := allowedNetworkRule2.ID()
	blockedNetworkRuleId1 := blockedNetworkRule1.ID()
	blockedNetworkRuleId2 := blockedNetworkRule2.ID()

	basic := model.NetworkPolicy("test", id.Name())

	complete := model.NetworkPolicy("test", newId.Name()).
		WithComment(comment).
		WithAllowedNetworkRules(allowedNetworkRuleId1, allowedNetworkRuleId2).
		WithBlockedNetworkRules(blockedNetworkRuleId1, blockedNetworkRuleId2).
		WithAllowedIps("1.1.1.1", "2.2.2.2").
		WithBlockedIps("3.3.3.3", "4.4.4.4")

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.NetworkPolicy(t, id).
			HasName(id.Name()).
			HasComment("").
			HasEntriesInAllowedIpList(0).
			HasEntriesInBlockedIpList(0).
			HasEntriesInAllowedNetworkRules(0).
			HasEntriesInBlockedNetworkRules(0),

		resourceassert.NetworkPolicyResource(t, basic.ResourceReference()).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasCommentString("").
			HasAllowedIpListLength(0).
			HasBlockedIpListLength(0).
			HasAllowedNetworkRuleListLength(0).
			HasBlockedNetworkRuleListLength(0),

		resourceshowoutputassert.NetworkPolicyShowOutput(t, basic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasComment("").
			HasEntriesInAllowedIpList(0).
			HasEntriesInBlockedIpList(0).
			HasEntriesInAllowedNetworkRules(0).
			HasEntriesInBlockedNetworkRules(0),

		assert.Check(resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(basic.ResourceReference(), "describe_output.#", "1"),
			resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
			resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
			resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.allowed_ip_list"),
			resource.TestCheckNoResourceAttr(basic.ResourceReference(), "describe_output.0.blocked_ip_list"),
		)),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		objectassert.NetworkPolicy(t, newId).
			HasName(newId.Name()).
			HasComment(comment).
			HasEntriesInAllowedIpList(2).
			HasEntriesInBlockedIpList(2).
			HasEntriesInAllowedNetworkRules(2).
			HasEntriesInBlockedNetworkRules(2),

		resourceassert.NetworkPolicyResource(t, complete.ResourceReference()).
			HasNameString(newId.Name()).
			HasFullyQualifiedNameString(newId.FullyQualifiedName()).
			HasCommentString(comment).
			HasAllowedIpListLength(2).
			HasBlockedIpListLength(2).
			HasAllowedNetworkRuleListLength(2).
			HasBlockedNetworkRuleListLength(2),

		resourceshowoutputassert.NetworkPolicyShowOutput(t, complete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasComment(comment).
			HasEntriesInAllowedIpList(2).
			HasEntriesInBlockedIpList(2).
			HasEntriesInAllowedNetworkRules(2).
			HasEntriesInBlockedNetworkRules(2),

		assert.Check(resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(complete.ResourceReference(), "describe_output.#", "1"),
			resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
			resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
			resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.allowed_ip_list"),
			resource.TestCheckResourceAttrSet(complete.ResourceReference(), "describe_output.0.blocked_ip_list"),
		)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      complete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					// TODO: Change
					testClient().NetworkPolicy.Update(t, sdk.NewAlterNetworkPolicyRequest(id).WithUnset(
						*sdk.NewNetworkPolicyUnsetRequest().
							WithAllowedIpList(true).
							WithBlockedIpList(true).
							WithAllowedNetworkRuleList(true).
							WithBlockedNetworkRuleList(true).
							WithComment(true),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Create - with optionals
			{
				Taint: []string{complete.ResourceReference()},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_NetworkPolicy_Rename(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	networkPolicyModelBasic := model.NetworkPolicy("test", id.Name())
	networkPolicyModelBasicNewId := model.NetworkPolicy("test", newId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, networkPolicyModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "id", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.name", id.Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, networkPolicyModelBasicNewId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(networkPolicyModelBasicNewId.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasicNewId.ResourceReference(), "id", newId.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasicNewId.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasicNewId.ResourceReference(), "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(networkPolicyModelBasicNewId.ResourceReference(), "show_output.0.name", newId.Name()),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_InvalidBlockedIpListValue(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config:      networkPolicyConfigInvalidBlockedIpListValue(id),
				ExpectError: regexp.MustCompile(`invalid value \(0.0.0.0/0\) set for a field \[{{} blocked_ip_list} {{} {{{{}`),
			},
		},
	})
}

func networkPolicyConfigInvalidBlockedIpListValue(networkPolicyId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%[1]s"
		blocked_ip_list = ["1.1.1.1", "0.0.0.0/0"]
	}`, networkPolicyId.Name())
}

func TestAcc_NetworkPolicy_InvalidNetworkRuleIds(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config:      networkPolicyConfigInvalidAllowedNetworkRules(id),
				ExpectError: regexp.MustCompile(`sdk\.TableColumnIdentifier\. The correct form of the fully qualified name for`),
			},
			{
				Config:      networkPolicyConfigInvalidAllowedNetworkRules(id),
				ExpectError: regexp.MustCompile(`sdk\.DatabaseObjectIdentifier\. The correct form of the fully qualified name`),
			},
			{
				Config:      networkPolicyConfigInvalidBlockedNetworkRules(id),
				ExpectError: regexp.MustCompile(`sdk\.TableColumnIdentifier\. The correct form of the fully qualified name for`),
			},
			{
				Config:      networkPolicyConfigInvalidBlockedNetworkRules(id),
				ExpectError: regexp.MustCompile(`sdk\.DatabaseObjectIdentifier\. The correct form of the fully qualified name`),
			},
		},
	})
}

func networkPolicyConfigInvalidAllowedNetworkRules(networkPolicyId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%[1]s"
		allowed_network_rule_list = ["a.b", "a.b.c.d"]
	}`, networkPolicyId.Name())
}

func networkPolicyConfigInvalidBlockedNetworkRules(networkPolicyId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%[1]s"
		blocked_network_rule_list = ["a.b", "a.b.c.d"]
	}`, networkPolicyId.Name())
}

func TestAcc_NetworkPolicy_Issue2236(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	allowedNetworkRule1, allowedNetworkRule1Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule1Cleanup)

	allowedNetworkRule2, allowedNetworkRule2Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule2Cleanup)

	blockedNetworkRule1, blockedNetworkRule1Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule1Cleanup)

	blockedNetworkRule2, blockedNetworkRule2Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule2Cleanup)

	allowedNetworkRuleId1 := allowedNetworkRule1.ID()
	allowedNetworkRuleId2 := allowedNetworkRule2.ID()
	blockedNetworkRuleId1 := blockedNetworkRule1.ID()
	blockedNetworkRuleId2 := blockedNetworkRule2.ID()

	networkPolicyWithNetworkRules := model.NetworkPolicy("test", id.Name()).
		WithAllowedNetworkRulesUnquotedNamePart(allowedNetworkRuleId1, allowedNetworkRuleId2).
		WithBlockedNetworkRulesUnquotedNamePart(blockedNetworkRuleId1, blockedNetworkRuleId2)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("0.93.0"),
				// Identifier quoting mismatch (no diff suppression)
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					func() { SetV097CompatibleConfigPathEnv(t) }()
				},
				Config: accconfig.FromModels(t, networkPolicyWithNetworkRules),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "blocked_network_rule_list.#", "2"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, networkPolicyWithNetworkRules),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "blocked_network_rule_list.#", "2"),

					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyWithNetworkRules.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPolicyWithNetworkRules.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	networkPolicyModelBasic := model.NetworkPolicy("test", id.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, networkPolicyModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, networkPolicyModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_WithQuotedName(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             networkPolicyConfigBasicWithQuotedName(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   networkPolicyConfigBasicWithQuotedName(id),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_network_policy.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_network_policy.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "id", id.Name()),
				),
			},
		},
	})
}

func networkPolicyConfigBasicWithQuotedName(networkPolicyId sdk.AccountObjectIdentifier) string {
	quotedId := fmt.Sprintf(`\"%s\"`, networkPolicyId.Name())
	return fmt.Sprintf(`resource "snowflake_network_policy" "test" {
		name = "%v"
	}`, quotedId)
}
