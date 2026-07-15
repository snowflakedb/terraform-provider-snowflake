package resourceassert

func (c *McpServerResourceAssert) HasCommentEmpty() *McpServerResourceAssert {
	c.StringValueSet("comment", "")
	return c
}
