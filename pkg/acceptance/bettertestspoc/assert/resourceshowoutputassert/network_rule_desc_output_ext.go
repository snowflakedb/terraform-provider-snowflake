package resourceshowoutputassert

func (n *NetworkRuleDescribeOutputAssert) HasValueList(expected []string) *NetworkRuleDescribeOutputAssert {
	n.ListContainsExactlyStringValuesInOrder("value_list", expected...)
	return n
}
