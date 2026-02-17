package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (n *NetworkRuleResourceAssert) HasValueList(expected []string) *NetworkRuleResourceAssert {
	n.AddAssertion(assert.ValueSet("value_list.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i, v := range expected {
		n.AddAssertion(assert.ValueSet(fmt.Sprintf("value_list.%d", i), v))
	}
	return n
}

func (n *NetworkRuleResourceAssert) HasTypeEnum(expected sdk.NetworkRuleType) *NetworkRuleResourceAssert {
	n.AddAssertion(assert.ValueSet("type", string(expected)))
	return n
}

func (n *NetworkRuleResourceAssert) HasModeEnum(expected sdk.NetworkRuleMode) *NetworkRuleResourceAssert {
	n.AddAssertion(assert.ValueSet("mode", string(expected)))
	return n
}
