package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func (d *DatabaseShowOutputAssert) HasCreatedOnNotEmpty() *DatabaseShowOutputAssert {
	d.ValuePresent("created_on")
	return d
}

func (d *DatabaseShowOutputAssert) HasIsCurrentNotEmpty() *DatabaseShowOutputAssert {
	d.ValuePresent("is_current")
	return d
}

func (d *DatabaseShowOutputAssert) HasOwnerNotEmpty() *DatabaseShowOutputAssert {
	d.ValuePresent("owner")
	return d
}

func (d *DatabaseShowOutputAssert) HasRetentionTimeNotEmpty() *DatabaseShowOutputAssert {
	d.ValuePresent("retention_time")
	return d
}

func (d *DatabaseShowOutputAssert) HasOwnerRoleTypeNotEmpty() *DatabaseShowOutputAssert {
	d.ValuePresent("owner_role_type")
	return d
}

func (d *DatabaseShowOutputAssert) HasOriginEmpty() *DatabaseShowOutputAssert {
	d.ValueSet("origin", "")
	return d
}

func (d *DatabaseShowOutputAssert) HasOrigin(expected sdk.ObjectIdentifier) *DatabaseShowOutputAssert {
	d.StringValueSet("origin", expected.FullyQualifiedName())
	return d
}
