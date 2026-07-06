package resourceshowoutputassert

func (a *AuthenticationPolicyShowOutputAssert) HasCreatedOnNotEmpty() *AuthenticationPolicyShowOutputAssert {
	a.ValuePresent("created_on")
	return a
}
