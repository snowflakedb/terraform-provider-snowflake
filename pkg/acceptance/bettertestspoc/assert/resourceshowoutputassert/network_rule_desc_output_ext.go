package resourceshowoutputassert

func (n *NetworkRuleDescribeOutputAssert) HasValueList(expected []string) *NetworkRuleDescribeOutputAssert {
	n.ListContainsExactlyStringValuesInOrder("value_list", expected...)
	return n
}

func (n *NetworkRuleDescribeOutputAssert) HasCommentEmpty() *NetworkRuleDescribeOutputAssert {
	n.ValueSet("comment", "")
	return n
}

func (n *NetworkRuleDescribeOutputAssert) HasCreatedOnNotEmpty() *NetworkRuleDescribeOutputAssert {
	n.ValuePresent("created_on")
	return n
}
