package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (n *NetworkPolicyResourceAssert) HasAllowedIpListLength(expected int) *NetworkPolicyResourceAssert {
	n.AddAssertion(assert.ValueSet("allowed_ip_list.#", fmt.Sprintf("%d", expected)))
	return n
}

func (n *NetworkPolicyResourceAssert) HasBlockedIpListLength(expected int) *NetworkPolicyResourceAssert {
	n.AddAssertion(assert.ValueSet("blocked_ip_list.#", fmt.Sprintf("%d", expected)))
	return n
}

func (n *NetworkPolicyResourceAssert) HasAllowedNetworkRuleListLength(expected int) *NetworkPolicyResourceAssert {
	n.AddAssertion(assert.ValueSet("allowed_network_rule_list.#", fmt.Sprintf("%d", expected)))
	return n
}

func (n *NetworkPolicyResourceAssert) HasBlockedNetworkRuleListLength(expected int) *NetworkPolicyResourceAssert {
	n.AddAssertion(assert.ValueSet("blocked_network_rule_list.#", fmt.Sprintf("%d", expected)))
	return n
}
