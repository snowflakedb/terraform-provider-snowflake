package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// ImageRepositoriesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func ImageRepositoriesDatasourceShowOutput(t *testing.T, name string) *ImageRepositoryShowOutputAssert {
	t.Helper()

	s := ImageRepositoryShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "image_repositories.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (p *ImageRepositoryShowOutputAssert) HasCreatedOnNotEmpty() *ImageRepositoryShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return p
}

func (p *ImageRepositoryShowOutputAssert) HasRepositoryUrlNotEmpty() *ImageRepositoryShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("repository_url"))
	return p
}
