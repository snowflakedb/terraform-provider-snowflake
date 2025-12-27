package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (g *GrantPrivilegesToAccountRoleResourceAssert) HasPrivileges(privileges ...string) *GrantPrivilegesToAccountRoleResourceAssert {
	g.AddAssertion(assert.ValueSet("privileges.#", strconv.FormatInt(int64(len(privileges)), 10)))
	for i, v := range privileges {
		g.AddAssertion(assert.ValueSet(fmt.Sprintf("privileges.%d", i), v))
	}
	return g
}
