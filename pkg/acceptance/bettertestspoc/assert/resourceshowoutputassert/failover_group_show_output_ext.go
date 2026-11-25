package resourceshowoutputassert

func (f *FailoverGroupShowOutputAssert) HasCreatedOnNotEmpty() *FailoverGroupShowOutputAssert {
	f.ValuePresent("created_on")
	return f
}

func (f *FailoverGroupShowOutputAssert) HasSnowflakeRegionNotEmpty() *FailoverGroupShowOutputAssert {
	f.ValuePresent("snowflake_region")
	return f
}

func (f *FailoverGroupShowOutputAssert) HasNextScheduledRefreshNotEmpty() *FailoverGroupShowOutputAssert {
	f.ValuePresent("next_scheduled_refresh")
	return f
}

func (f *FailoverGroupShowOutputAssert) HasSecondaryStateEmpty() *FailoverGroupShowOutputAssert {
	f.StringValueSet("secondary_state", "")
	return f
}

func (f *FailoverGroupShowOutputAssert) HasAllowedAccounts(expected ...string) *FailoverGroupShowOutputAssert {
	f.SetContainsExactlyStringValues("allowed_accounts", expected...)
	return f
}

func (f *FailoverGroupShowOutputAssert) HasObjectTypes(expected ...string) *FailoverGroupShowOutputAssert {
	f.SetContainsExactlyStringValues("object_types", expected...)
	return f
}

func (f *FailoverGroupShowOutputAssert) HasAllowedIntegrationTypes(expected ...string) *FailoverGroupShowOutputAssert {
	f.SetContainsExactlyStringValues("allowed_integration_types", expected...)
	return f
}
