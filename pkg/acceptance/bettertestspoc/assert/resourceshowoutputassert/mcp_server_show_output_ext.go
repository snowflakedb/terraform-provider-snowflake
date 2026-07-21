package resourceshowoutputassert

func (c *McpServerShowOutputAssert) HasCreatedOnNotEmpty() *McpServerShowOutputAssert {
	c.ValuePresent("created_on")
	return c
}
