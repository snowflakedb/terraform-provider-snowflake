package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (g *GrantPrivilegesToAccountRoleResourceAssert) HasOnAccountObject(objectType sdk.ObjectType, id sdk.AccountObjectIdentifier) *GrantPrivilegesToAccountRoleResourceAssert {
	g.AddAssertion(assert.ValueSet("on_account_object.#", "1"))
	g.AddAssertion(assert.ValueSet("on_account_object.0.object_type", objectType.String()))
	g.AddAssertion(assert.ValueSet("on_account_object.0.object_name", id.Name()))
	return g
}

func (g *GrantPrivilegesToAccountRoleResourceAssert) HasResourceId(expected string) *GrantPrivilegesToAccountRoleResourceAssert {
	g.AddAssertion(assert.ValueSet("id", expected))
	return g
}
