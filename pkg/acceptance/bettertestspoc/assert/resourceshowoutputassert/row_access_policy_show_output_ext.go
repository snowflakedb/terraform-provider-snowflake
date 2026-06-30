package resourceshowoutputassert

func (r *RowAccessPolicyShowOutputAssert) HasCreatedOnNotEmpty() *RowAccessPolicyShowOutputAssert {
	r.ValuePresent("created_on")
	return r
}
