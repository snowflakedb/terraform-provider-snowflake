package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (g *GrantPrivilegesToAccountRoleResourceAssert) HasPrivileges(privileges ...string) *GrantPrivilegesToAccountRoleResourceAssert {
	g.AddAssertion(assert.ValueSet("privileges.#", strconv.FormatInt(int64(len(privileges)), 10)))
	for _, v := range privileges {
		g.AddAssertion(assert.SetElem("privileges", v))
	}
	return g
}
