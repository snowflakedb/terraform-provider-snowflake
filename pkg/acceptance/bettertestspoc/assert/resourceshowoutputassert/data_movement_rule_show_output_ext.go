package resourceshowoutputassert

func (d *DataMovementRuleShowOutputAssert) HasCreatedOnNotEmpty() *DataMovementRuleShowOutputAssert {
	d.ValuePresent("created_on")
	return d
}
