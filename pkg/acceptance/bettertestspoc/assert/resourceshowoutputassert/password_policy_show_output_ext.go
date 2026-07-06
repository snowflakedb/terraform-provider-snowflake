package resourceshowoutputassert

func (p *PasswordPolicyShowOutputAssert) HasCreatedOnNotEmpty() *PasswordPolicyShowOutputAssert {
	p.ValuePresent("created_on")
	return p
}
