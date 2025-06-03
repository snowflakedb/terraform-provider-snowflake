package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (c *GitRepositoryShowOutputAssert) HasCreatedOnNotEmpty() *GitRepositoryShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}

func (c *GitRepositoryShowOutputAssert) HasGitCredentialsEmpty() *GitRepositoryShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("git_credentials", ""))
	return c
}
