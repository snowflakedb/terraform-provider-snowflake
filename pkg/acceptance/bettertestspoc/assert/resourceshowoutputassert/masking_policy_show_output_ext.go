package resourceshowoutputassert

func (p *MaskingPolicyShowOutputAssert) HasCreatedOnNotEmpty() *MaskingPolicyShowOutputAssert {
	p.ValuePresent("created_on")
	return p
}

func (p *MaskingPolicyShowOutputAssert) HasOwnerNotEmpty() *MaskingPolicyShowOutputAssert {
	p.ValuePresent("owner")
	return p
}

func (p *MaskingPolicyShowOutputAssert) HasOwnerRoleTypeNotEmpty() *MaskingPolicyShowOutputAssert {
	p.ValuePresent("owner_role_type")
	return p
}
