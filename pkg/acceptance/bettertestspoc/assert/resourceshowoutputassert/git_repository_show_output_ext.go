package resourceshowoutputassert

func (c *GitRepositoryShowOutputAssert) HasCreatedOnNotEmpty() *GitRepositoryShowOutputAssert {
	c.ValuePresent("created_on")
	return c
}

func (c *GitRepositoryShowOutputAssert) HasGitCredentialsEmpty() *GitRepositoryShowOutputAssert {
	c.StringValueSet("git_credentials", "")
	return c
}
