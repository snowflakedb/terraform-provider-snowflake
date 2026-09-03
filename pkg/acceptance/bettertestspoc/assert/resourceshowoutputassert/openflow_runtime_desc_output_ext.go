package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowRuntimeDescribeOutputAssert) HasKeyNotEmpty() *OpenflowRuntimeDescribeOutputAssert {
	o.ValuePresent("key")
	return o
}

func (o *OpenflowRuntimeDescribeOutputAssert) HasServerUrlNotEmpty() *OpenflowRuntimeDescribeOutputAssert {
	o.ValuePresent("server_url")
	return o
}

func (o *OpenflowRuntimeDescribeOutputAssert) HasNodeTypeTierNotEmpty() *OpenflowRuntimeDescribeOutputAssert {
	o.ValuePresent("node_type_tier")
	return o
}

func (o *OpenflowRuntimeDescribeOutputAssert) HasExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *OpenflowRuntimeDescribeOutputAssert {
	o.ValueSet("external_access_integrations.#", strconv.Itoa(len(ids)))
	for _, id := range ids {
		o.SetContainsElem("external_access_integrations", id.Name())
	}
	return o
}
