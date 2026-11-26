package resourceshowoutputassert

func (b *AggregationPolicyShowOutputAssert) HasCreatedOnNotEmpty() *AggregationPolicyShowOutputAssert {
	b.ValuePresent("created_on")
	return b
}
