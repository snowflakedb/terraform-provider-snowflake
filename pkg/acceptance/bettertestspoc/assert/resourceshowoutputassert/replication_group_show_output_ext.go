package resourceshowoutputassert

func (r *ReplicationGroupShowOutputAssert) HasCreatedOnNotEmpty() *ReplicationGroupShowOutputAssert {
	r.ValuePresent("created_on")
	return r
}

func (r *ReplicationGroupShowOutputAssert) HasSnowflakeRegionNotEmpty() *ReplicationGroupShowOutputAssert {
	r.ValuePresent("snowflake_region")
	return r
}

func (r *ReplicationGroupShowOutputAssert) HasNextScheduledRefreshNotEmpty() *ReplicationGroupShowOutputAssert {
	r.ValuePresent("next_scheduled_refresh")
	return r
}

func (r *ReplicationGroupShowOutputAssert) HasSecondaryStateEmpty() *ReplicationGroupShowOutputAssert {
	r.StringValueSet("secondary_state", "")
	return r
}

func (r *ReplicationGroupShowOutputAssert) HasAllowedAccounts(expected ...string) *ReplicationGroupShowOutputAssert {
	r.ListContainsExactlyStringValuesInOrder("allowed_accounts", expected...)
	return r
}
