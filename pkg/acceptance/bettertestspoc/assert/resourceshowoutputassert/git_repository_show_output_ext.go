package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// GitRepositoriesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func GitRepositoriesDatasourceShowOutput(t *testing.T, name string) *GitRepositoryShowOutputAssert {
	t.Helper()

	s := GitRepositoryShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "git_repositories.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (c *GitRepositoryShowOutputAssert) HasCreatedOnNotEmpty() *GitRepositoryShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}

func (c *GitRepositoryShowOutputAssert) HasGitCredentialsEmpty() *GitRepositoryShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("git_credentials", ""))
	return c
}
