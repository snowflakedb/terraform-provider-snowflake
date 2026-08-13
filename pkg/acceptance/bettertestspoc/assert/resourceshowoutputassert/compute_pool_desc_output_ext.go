package resourceshowoutputassert

func (c *ComputePoolDescribeOutputAssert) HasCreatedOnNotEmpty() *ComputePoolDescribeOutputAssert {
	c.ValuePresent("created_on")
	return c
}

func (c *ComputePoolDescribeOutputAssert) HasResumedOnNotEmpty() *ComputePoolDescribeOutputAssert {
	c.ValuePresent("resumed_on")
	return c
}

func (c *ComputePoolDescribeOutputAssert) HasUpdatedOnNotEmpty() *ComputePoolDescribeOutputAssert {
	c.ValuePresent("updated_on")
	return c
}

func (c *ComputePoolDescribeOutputAssert) HasApplicationEmpty() *ComputePoolDescribeOutputAssert {
	c.StringValueSet("application", "")
	return c
}
