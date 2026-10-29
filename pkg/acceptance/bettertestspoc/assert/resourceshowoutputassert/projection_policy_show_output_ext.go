package resourceshowoutputassert

func (b *ProjectionPolicyShowOutputAssert) HasCreatedOnNotEmpty() *ProjectionPolicyShowOutputAssert {
	b.ValuePresent("created_on")
	return b
}
