package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (o *OauthIntegrationForCustomClientsResourceAssert) HasPreAuthorizedRolesList(values ...string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("pre_authorized_roles_list.#", fmt.Sprintf("%d", len(values))))
	for i, value := range values {
		o.AddAssertion(assert.ValueSet(fmt.Sprintf("pre_authorized_roles_list.%d", i), value))
	}
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasRelatedParametersNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters.#", "1"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters.0.oauth_add_privileged_roles_to_blocked_list.0.value", expected))
	return o
}
