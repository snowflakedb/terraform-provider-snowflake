package resourceshowoutputassert

func (n *NetworkPolicyShowOutputAssert) HasCreatedOnNotEmpty() *NetworkPolicyShowOutputAssert {
	n.ValuePresent("created_on")
	return n
}
