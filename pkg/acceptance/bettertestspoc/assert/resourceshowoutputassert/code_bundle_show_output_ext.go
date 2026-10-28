package resourceshowoutputassert

func (c *CodeBundleShowOutputAssert) HasCreatedOnNotEmpty() *CodeBundleShowOutputAssert {
	c.ValuePresent("created_on")
	return c
}
