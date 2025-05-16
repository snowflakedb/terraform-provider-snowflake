//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicy_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	allowedNetworkRule1, allowedNetworkRule1Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule1Cleanup)

	allowedNetworkRule2, allowedNetworkRule2Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule2Cleanup)

	blockedNetworkRule1, blockedNetworkRule1Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule1Cleanup)

	blockedNetworkRule2, blockedNetworkRule2Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule2Cleanup)

	allowedNetworkRuleId1 := allowedNetworkRule1.ID()
	allowedNetworkRuleId2 := allowedNetworkRule2.ID()
	blockedNetworkRuleId1 := blockedNetworkRule1.ID()
	blockedNetworkRuleId2 := blockedNetworkRule2.ID()

	networkPolicyModelBasic := model.NetworkPolicy("test", id.Name())
	networkPolicyModelComplete := model.NetworkPolicy("test", id.Name()).
		WithComment(comment).
		WithAllowedNetworkRules(allowedNetworkRuleId1, allowedNetworkRuleId2).
		WithBlockedNetworkRules(blockedNetworkRuleId1, blockedNetworkRuleId2).
		WithAllowedIps("1.1.1.1", "2.2.2.2").
		WithBlockedIps("3.3.3.3", "4.4.4.4")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, networkPolicyModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "allowed_ip_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "blocked_ip_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "allowed_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "blocked_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelBasic.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_blocked_network_rules", "0"),

					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.allowed_ip_list"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.blocked_ip_list"),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, networkPolicyModelBasic),
				ResourceName: "snowflake_network_policy.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_ip_list"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_ip_list"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_network_rule_list"),
					importchecks.TestCheckResourceAttrNotInInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_network_rule_list"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, networkPolicyModelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(networkPolicyModelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "blocked_ip_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "blocked_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.blocked_ip_list"),
				),
			},
			// import - complete
			{
				Config:       accconfig.FromModels(t, networkPolicyModelComplete),
				ResourceName: "snowflake_network_policy.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// change externally
			{
				PreConfig: func() {
					acc.TestClient().NetworkPolicy.Update(t, sdk.NewAlterNetworkPolicyRequest(id).WithUnset(
						*sdk.NewNetworkPolicyUnsetRequest().
							WithAllowedIpList(true).
							WithBlockedIpList(true).
							WithAllowedNetworkRuleList(true).
							WithBlockedNetworkRuleList(true).
							WithComment(true),
					))
				},
				Config: accconfig.FromModels(t, networkPolicyModelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(networkPolicyModelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "blocked_ip_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "blocked_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.blocked_ip_list"),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, networkPolicyModelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(networkPolicyModelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "allowed_ip_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "blocked_ip_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "allowed_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "blocked_network_rule_list.#", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelBasic.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.comment", ""),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "show_output.0.entries_in_blocked_network_rules", "0"),

					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.allowed_ip_list"),
					resource.TestCheckNoResourceAttr(networkPolicyModelBasic.ResourceReference(), "describe_output.0.blocked_ip_list"),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	allowedNetworkRule1, allowedNetworkRule1Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule1Cleanup)

	allowedNetworkRule2, allowedNetworkRule2Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule2Cleanup)

	blockedNetworkRule1, blockedNetworkRule1Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule1Cleanup)

	blockedNetworkRule2, blockedNetworkRule2Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule2Cleanup)

	allowedNetworkRuleId1 := allowedNetworkRule1.ID()
	allowedNetworkRuleId2 := allowedNetworkRule2.ID()
	blockedNetworkRuleId1 := blockedNetworkRule1.ID()
	blockedNetworkRuleId2 := blockedNetworkRule2.ID()

	networkPolicyModelComplete := model.NetworkPolicy("test", id.Name()).
		WithComment(comment).
		WithAllowedNetworkRules(allowedNetworkRuleId1, allowedNetworkRuleId2).
		WithBlockedNetworkRules(blockedNetworkRuleId1, blockedNetworkRuleId2).
		WithAllowedIps("1.1.1.1", "2.2.2.2").
		WithBlockedIps("3.3.3.3", "4.4.4.4")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, networkPolicyModelComplete),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "allowed_ip_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "blocked_ip_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "blocked_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.comment", comment),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "show_output.0.entries_in_blocked_network_rules", "2"),

					resource.TestCheckResourceAttr(networkPolicyModelComplete.ResourceReference(), "describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.blocked_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet(networkPolicyModelComplete.ResourceReference(), "describe_output.0.blocked_ip_list"),
				),
			},
			{
				Config:       accconfig.FromModels(t, networkPolicyModelComplete),
				ResourceName: "snowflake_network_policy.test",
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_ip_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "allowed_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "blocked_network_rule_list.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_Rename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	networkPolicyModelBasic := model.NetworkPolicy("test", id.Name())
	networkPolicyModelBasicNewId := model.NetworkPolicy("test", newId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	allowedNetworkRule1, allowedNetworkRule1Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule1Cleanup)

	allowedNetworkRule2, allowedNetworkRule2Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule2Cleanup)

	blockedNetworkRule1, blockedNetworkRule1Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule1Cleanup)

	blockedNetworkRule2, blockedNetworkRule2Cleanup := acc.TestClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule2Cleanup)

	allowedNetworkRuleId1 := allowedNetworkRule1.ID()
	allowedNetworkRuleId2 := allowedNetworkRule2.ID()
	blockedNetworkRuleId1 := blockedNetworkRule1.ID()
	blockedNetworkRuleId2 := blockedNetworkRule2.ID()

	networkPolicyWithNetworkRules := model.NetworkPolicy("test", id.Name()).
		WithAllowedNetworkRulesUnquotedNamePart(allowedNetworkRuleId1, allowedNetworkRuleId2).
		WithBlockedNetworkRulesUnquotedNamePart(blockedNetworkRuleId1, blockedNetworkRuleId2)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.93.0"),
				// Identifier quoting mismatch (no diff suppression)
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					func() { acc.SetV097CompatibleConfigPathEnv(t) }()
				},
				Config: accconfig.FromModels(t, networkPolicyWithNetworkRules),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "allowed_network_rule_list.#", "2"),
					resource.TestCheckResourceAttr(networkPolicyWithNetworkRules.ResourceReference(), "blocked_network_rule_list.#", "2"),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	networkPolicyModelBasic := model.NetworkPolicy("test", id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.94.1"),
				Config:            accconfig.FromModels(t, networkPolicyModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, networkPolicyModelBasic),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkPolicyModelBasic.ResourceReference(), "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicy_WithQuotedName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				PreConfig:          func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders:  acc.ExternalProviderWithExactVersion("0.94.1"),
				ExpectNonEmptyPlan: true,
				Config:             networkPolicyConfigBasicWithQuotedName(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_network_policy.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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
