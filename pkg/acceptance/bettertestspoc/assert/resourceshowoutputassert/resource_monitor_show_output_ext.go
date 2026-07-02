package resourceshowoutputassert

func (r *ResourceMonitorShowOutputAssert) HasStartTimeNotEmpty() *ResourceMonitorShowOutputAssert {
	r.ValuePresent("start_time")
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasEndTimeNotEmpty() *ResourceMonitorShowOutputAssert {
	r.ValuePresent("end_time")
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasCreatedOnNotEmpty() *ResourceMonitorShowOutputAssert {
	r.ValuePresent("created_on")
	return r
}

func (r *ResourceMonitorShowOutputAssert) HasOwnerNotEmpty() *ResourceMonitorShowOutputAssert {
	r.ValuePresent("owner")
	return r
}

// TODO [next PRs]: the show output has improper name for the field and the logic is also incorrect for output mapping; solving as a separate change
func (r *ResourceMonitorShowOutputAssert) HasSuspendImmediateAt(expected int) *ResourceMonitorShowOutputAssert {
	r.IntValueSet("suspend_immediate_at", expected)
	return r
}
