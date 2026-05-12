package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasRelatedParametersNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters.#", "1"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters.0.oauth_add_privileged_roles_to_blocked_list.0.value", expected))
	return o
}
