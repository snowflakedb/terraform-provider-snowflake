package resourceshowoutputassert

func (b *JoinPolicyShowOutputAssert) HasCreatedOnNotEmpty() *JoinPolicyShowOutputAssert {
	b.ValuePresent("created_on")
	return b
}
