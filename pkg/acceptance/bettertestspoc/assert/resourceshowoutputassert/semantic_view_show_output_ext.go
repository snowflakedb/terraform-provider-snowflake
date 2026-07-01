package resourceshowoutputassert

func (c *SemanticViewShowOutputAssert) HasCreatedOnNotEmpty() *SemanticViewShowOutputAssert {
	c.ValuePresent("created_on")
	return c
}
