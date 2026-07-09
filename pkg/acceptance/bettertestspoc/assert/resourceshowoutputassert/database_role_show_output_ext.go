package resourceshowoutputassert

func (d *DatabaseRoleShowOutputAssert) HasCreatedOnNotEmpty() *DatabaseRoleShowOutputAssert {
	d.ValuePresent("created_on")
	return d
}

func (d *DatabaseRoleShowOutputAssert) HasOwnerNotEmpty() *DatabaseRoleShowOutputAssert {
	d.ValuePresent("owner")
	return d
}

func (d *DatabaseRoleShowOutputAssert) HasOwnerRoleTypeNotEmpty() *DatabaseRoleShowOutputAssert {
	d.ValuePresent("owner_role_type")
	return d
}
