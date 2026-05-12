package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (n *NetworkRuleResourceAssert) HasTypeEnum(expected sdk.NetworkRuleType) *NetworkRuleResourceAssert {
	n.AddAssertion(assert.ValueSet("type", string(expected)))
	return n
}

func (n *NetworkRuleResourceAssert) HasModeEnum(expected sdk.NetworkRuleMode) *NetworkRuleResourceAssert {
	n.AddAssertion(assert.ValueSet("mode", string(expected)))
	return n
}
