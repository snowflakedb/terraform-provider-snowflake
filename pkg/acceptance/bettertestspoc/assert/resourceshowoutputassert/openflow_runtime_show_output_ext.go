package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowRuntimeShowOutputAssert) HasCreatedOnNotEmpty() *OpenflowRuntimeShowOutputAssert {
	o.ValuePresent("created_on")
	return o
}

func (o *OpenflowRuntimeShowOutputAssert) HasUpdatedOnNotEmpty() *OpenflowRuntimeShowOutputAssert {
	o.ValuePresent("updated_on")
	return o
}

func (o *OpenflowRuntimeShowOutputAssert) HasKeyNotEmpty() *OpenflowRuntimeShowOutputAssert {
	o.ValuePresent("key")
	return o
}

func (o *OpenflowRuntimeShowOutputAssert) HasExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *OpenflowRuntimeShowOutputAssert {
	o.ValueSet("external_access_integrations.#", strconv.Itoa(len(ids)))
	for _, id := range ids {
		o.SetContainsElem("external_access_integrations", id.Name())
	}
	return o
}
