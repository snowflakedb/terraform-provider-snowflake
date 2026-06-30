package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (n *NetworkRuleResourceAssert) HasTypeEnum(expected sdk.NetworkRuleType) *NetworkRuleResourceAssert {
	n.ValueSet("type", string(expected))
	return n
}

func (n *NetworkRuleResourceAssert) HasModeEnum(expected sdk.NetworkRuleMode) *NetworkRuleResourceAssert {
	n.ValueSet("mode", string(expected))
	return n
}
