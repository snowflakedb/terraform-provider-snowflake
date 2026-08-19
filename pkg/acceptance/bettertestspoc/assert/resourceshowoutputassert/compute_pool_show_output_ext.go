package resourceshowoutputassert

func (c *ComputePoolShowOutputAssert) HasCreatedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.ValuePresent("created_on")
	return c
}

func (c *ComputePoolShowOutputAssert) HasResumedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.ValuePresent("resumed_on")
	return c
}

func (c *ComputePoolShowOutputAssert) HasUpdatedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.ValuePresent("updated_on")
	return c
}

func (c *ComputePoolShowOutputAssert) HasApplicationEmpty() *ComputePoolShowOutputAssert {
	c.StringValueSet("application", "")
	return c
}

func (c *ComputePoolShowOutputAssert) HasBackupInstanceFamilies(expected ...string) *ComputePoolShowOutputAssert {
	c.ListContainsExactlyStringValuesInOrder("backup_instance_families", expected...)
	return c
}

func (c *ComputePoolDescribeOutputAssert) HasBackupInstanceFamilies(expected ...string) *ComputePoolDescribeOutputAssert {
	c.ListContainsExactlyStringValuesInOrder("backup_instance_families", expected...)
	return c
}
