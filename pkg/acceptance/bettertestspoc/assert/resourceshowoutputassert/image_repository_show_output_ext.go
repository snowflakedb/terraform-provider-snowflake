package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (p *ImageRepositoryShowOutputAssert) HasCreatedOnNotEmpty() *ImageRepositoryShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return p
}

func (p *ImageRepositoryShowOutputAssert) HasRepositoryUrlNotEmpty() *ImageRepositoryShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("repository_url"))
	return p
}
