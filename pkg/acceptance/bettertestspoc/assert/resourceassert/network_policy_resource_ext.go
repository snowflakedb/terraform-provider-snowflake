package resourceassert

import (
	"fmt"
)

// TODO [SNOW-3127764]: Use set asserts in acceptance tests instead
func (n *NetworkPolicyResourceAssert) HasAllowedIpListLength(expected int) *NetworkPolicyResourceAssert {
	n.ValueSet("allowed_ip_list.#", fmt.Sprintf("%d", expected))
	return n
}

func (n *NetworkPolicyResourceAssert) HasBlockedIpListLength(expected int) *NetworkPolicyResourceAssert {
	n.ValueSet("blocked_ip_list.#", fmt.Sprintf("%d", expected))
	return n
}

func (n *NetworkPolicyResourceAssert) HasAllowedNetworkRuleListLength(expected int) *NetworkPolicyResourceAssert {
	n.ValueSet("allowed_network_rule_list.#", fmt.Sprintf("%d", expected))
	return n
}

func (n *NetworkPolicyResourceAssert) HasBlockedNetworkRuleListLength(expected int) *NetworkPolicyResourceAssert {
	n.ValueSet("blocked_network_rule_list.#", fmt.Sprintf("%d", expected))
	return n
}
