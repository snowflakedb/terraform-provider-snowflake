package resourceshowoutputassert

func (t *TableShowOutputAssert) HasDroppedOnNotEmpty() *TableShowOutputAssert {
	t.ValuePresent("dropped_on")
	return t
}
