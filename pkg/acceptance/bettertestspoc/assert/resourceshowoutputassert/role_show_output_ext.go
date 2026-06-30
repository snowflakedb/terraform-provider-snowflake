package resourceshowoutputassert

func (r *RoleShowOutputAssert) HasCreatedOnNotEmpty() *RoleShowOutputAssert {
	r.ValuePresent("created_on")
	return r
}
