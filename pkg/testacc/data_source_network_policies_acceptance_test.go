//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_NetworkPolicies_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	np1 := model.NetworkPolicy("np1", id1.Name())

	np2 := model.NetworkPolicy("np2", id2.Name())

	likePrefix := datasourcemodel.NetworkPolicies("test").
		WithLike(prefix+"%").
		WithDependsOn(np1.ResourceReference(), np2.ResourceReference())

	likeExact := datasourcemodel.NetworkPolicies("test").
		WithLike(id1.Name()).
		WithDependsOn(np1.ResourceReference(), np2.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			// like (prefix)
			{
				Config: config.FromModels(t, np1, np2, likePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "network_policies.#", "2"),
				),
			},
			// like (exact)
			{
				Config: config.FromModels(t, np1, np2, likeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(likeExact.DatasourceReference(), "network_policies.#", "1"),
					resource.TestCheckResourceAttr(likeExact.DatasourceReference(), "network_policies.0.show_output.0.name", id1.Name()),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicies_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	allowedNetworkRule1, allowedNetworkRule1Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule1Cleanup)

	allowedNetworkRule2, allowedNetworkRule2Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(allowedNetworkRule2Cleanup)

	blockedNetworkRule1, blockedNetworkRule1Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule1Cleanup)

	blockedNetworkRule2, blockedNetworkRule2Cleanup := testClient().NetworkRule.CreateIngress(t)
	t.Cleanup(blockedNetworkRule2Cleanup)

	networkPolicyModel := model.NetworkPolicy("test", id.Name()).
		WithComment(comment).
		WithAllowedNetworkRules(allowedNetworkRule1.ID(), allowedNetworkRule2.ID()).
		WithBlockedNetworkRules(blockedNetworkRule1.ID(), blockedNetworkRule2.ID()).
		WithAllowedIps("1.1.1.1", "2.2.2.2").
		WithBlockedIps("3.3.3.3", "4.4.4.4")

	withoutDescribe := datasourcemodel.NetworkPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithDependsOn(networkPolicyModel.ResourceReference())

	withDescribe := datasourcemodel.NetworkPolicies("test").
		WithWithDescribe(true).
		WithLike(id.Name()).
		WithDependsOn(networkPolicyModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, networkPolicyModel, withoutDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.NetworkPoliciesDatasourceShowOutput(t, withoutDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasEntriesInAllowedIpList(2).
						HasEntriesInBlockedIpList(2).
						HasEntriesInAllowedNetworkRules(2).
						HasEntriesInBlockedNetworkRules(2).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr(withoutDescribe.DatasourceReference(), "network_policies.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(withoutDescribe.DatasourceReference(), "network_policies.0.describe_output.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, networkPolicyModel, withDescribe),
				Check: assertThat(t,
					resourceshowoutputassert.NetworkPoliciesDatasourceShowOutput(t, withDescribe.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasEntriesInAllowedIpList(2).
						HasEntriesInBlockedIpList(2).
						HasEntriesInAllowedNetworkRules(2).
						HasEntriesInBlockedNetworkRules(2).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr(withDescribe.DatasourceReference(), "network_policies.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(withDescribe.DatasourceReference(), "network_policies.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttrSet(withDescribe.DatasourceReference(), "network_policies.0.describe_output.0.allowed_ip_list")),
					assert.Check(resource.TestCheckResourceAttrSet(withDescribe.DatasourceReference(), "network_policies.0.describe_output.0.blocked_ip_list")),
					assert.Check(resource.TestCheckResourceAttrSet(withDescribe.DatasourceReference(), "network_policies.0.describe_output.0.allowed_network_rule_list")),
					assert.Check(resource.TestCheckResourceAttrSet(withDescribe.DatasourceReference(), "network_policies.0.describe_output.0.blocked_network_rule_list")),
				),
			},
		},
	})
}

func TestAcc_NetworkPolicies_NetworkPolicyNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
