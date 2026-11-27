package resourceshowoutputassert

func (d *DataMovementRuleDescribeOutputAssert) HasCreatedOnNotEmpty() *DataMovementRuleDescribeOutputAssert {
	d.ValuePresent("created_on")
	return d
}
