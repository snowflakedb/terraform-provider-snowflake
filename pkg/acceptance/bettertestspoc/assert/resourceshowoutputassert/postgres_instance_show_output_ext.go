package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (p *PostgresInstanceShowOutputAssert) HasCreatedOnNotEmpty() *PostgresInstanceShowOutputAssert {
	p.AddAssertion(assert.ValuePresent("created_on"))
	return p
}

func (p *PostgresInstanceShowOutputAssert) HasIsHa(expected bool) *PostgresInstanceShowOutputAssert {
	p.BoolValueSet("is_ha", expected)
	return p
}
