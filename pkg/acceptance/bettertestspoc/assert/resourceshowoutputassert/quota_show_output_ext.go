package resourceshowoutputassert

func (q *QuotaShowOutputAssert) HasCreatedOnNotEmpty() *QuotaShowOutputAssert {
	q.ValuePresent("created_on")
	return q
}

func (q *QuotaShowOutputAssert) HasOwnerNotEmpty() *QuotaShowOutputAssert {
	q.ValuePresent("owner")
	return q
}
