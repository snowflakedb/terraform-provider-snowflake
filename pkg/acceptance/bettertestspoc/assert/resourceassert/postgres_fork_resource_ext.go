package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (p *PostgresForkResourceAssert) HasAuthenticationAuthorityString(expected string) *PostgresForkResourceAssert {
	p.AddAssertion(assert.ValueSet("authentication_authority", expected))
	return p
}
