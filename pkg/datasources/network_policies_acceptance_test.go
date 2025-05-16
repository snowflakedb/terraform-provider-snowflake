//go:build !account_level_tests

package datasources_test

import (
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicies_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	id2 := acc.TestClient().Ids.RandomAccountObjectIdentifier()
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

	networkPolicyModel1 := model.NetworkPolicy("test", id.Name()).
		WithComment(comment).
		WithAllowedNetworkRules(allowedNetworkRuleId1, allowedNetworkRuleId2).
		WithBlockedNetworkRules(blockedNetworkRuleId1, blockedNetworkRuleId2).
		WithAllowedIps("1.1.1.1", "2.2.2.2").
		WithBlockedIps("3.3.3.3", "4.4.4.4")
	networkPolicyModel2 := model.NetworkPolicy("test2", id2.Name())
	networkPoliciesModel := datasourcemodel.NetworkPolicies("test").
		WithLike(id.Name()).
		WithDependsOn(networkPolicyModel1.ResourceReference(), networkPolicyModel2.ResourceReference())
	networkPoliciesModel2WithDescribe := datasourcemodel.NetworkPolicies("test").
		WithWithDescribe(true).
		WithLike(id2.Name()).
		WithDependsOn(networkPolicyModel1.ResourceReference(), networkPolicyModel2.ResourceReference())
	networkPoliciesModel1WithDescribe := datasourcemodel.NetworkPolicies("test").
		WithWithDescribe(true).
		WithLike(id.Name()).
		WithDependsOn(networkPolicyModel1.ResourceReference(), networkPolicyModel2.ResourceReference())
	networkPoliciesModel1WithoutDescribe := datasourcemodel.NetworkPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithDependsOn(networkPolicyModel1.ResourceReference(), networkPolicyModel2.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, networkPolicyModel1, networkPolicyModel2, networkPoliciesModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.comment", comment),

					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.blocked_ip_list"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.blocked_network_rule_list"),
				),
			},
			{
				Config: accconfig.FromModels(t, networkPolicyModel1, networkPolicyModel2, networkPoliciesModel1WithDescribe),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.comment", comment),

					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.allowed_ip_list"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.blocked_ip_list"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.allowed_network_rule_list"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.blocked_network_rule_list"),
				),
			},
			{
				Config: accconfig.FromModels(t, networkPolicyModel1, networkPolicyModel2, networkPoliciesModel2WithDescribe),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.name", id2.Name()),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_ip_list", "0"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_ip_list", "0"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_network_rules", "0"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_network_rules", "0"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.comment", ""),

					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.#", "1"),
					resource.TestCheckNoResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.allowed_ip_list"),
					resource.TestCheckNoResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.blocked_ip_list"),
					resource.TestCheckNoResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.allowed_network_rule_list"),
					resource.TestCheckNoResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.0.blocked_network_rule_list"),
				),
			},
			{
				Config: accconfig.FromModels(t, networkPolicyModel1, networkPolicyModel2, networkPoliciesModel1WithoutDescribe),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_ip_list", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_allowed_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.entries_in_blocked_network_rules", "2"),
					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.show_output.0.comment", comment),

					resource.TestCheckResourceAttr(networkPoliciesModel.DatasourceReference(), "network_policies.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicies_NetworkPolicyNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      networkPolicyConfigBasicWithPostConditions(),
				ExpectError: regexp.MustCompile("there should be at least one network policy"),
			},
		},
	})
}

func networkPolicyConfigBasicWithPostConditions() string {
	return `
	data "snowflake_network_policies" "test" {
		like = "non_existing_network_policy"
	  	lifecycle {
			postcondition {
		  		condition     = length(self.network_policies) > 0
		  		error_message = "there should be at least one network policy"
			}
	  	}
	}
	`
}
