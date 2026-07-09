package resourceassert

func (o *OauthIntegrationForCustomClientsResourceAssert) HasRelatedParametersNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.ValueSet("related_parameters.#", "1")
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasRelatedParametersOauthAddPrivilegedRolesToBlockedList(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.ValueSet("related_parameters.0.oauth_add_privileged_roles_to_blocked_list.0.value", expected)
	return o
}
