package resourceshowoutputassert

func (c *McpServerDescribeOutputAssert) HasCreatedOnNotEmpty() *McpServerDescribeOutputAssert {
	c.ValuePresent("created_on")
	return c
}
