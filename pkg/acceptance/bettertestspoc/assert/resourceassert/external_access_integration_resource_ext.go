package resourceassert

func (e *ExternalAccessIntegrationResourceAssert) HasAllowedAuthenticationSecretsNone() *ExternalAccessIntegrationResourceAssert {
	e.ValueSet("allowed_authentication_secrets.#", "1")
	e.ValueSet("allowed_authentication_secrets.0.none", "true")
	e.ValueSet("allowed_authentication_secrets.0.all", "false")
	e.ValueSet("allowed_authentication_secrets.0.secrets.#", "0")
	return e
}

func (e *ExternalAccessIntegrationResourceAssert) HasAllowedAuthenticationSecretsAll() *ExternalAccessIntegrationResourceAssert {
	e.ValueSet("allowed_authentication_secrets.#", "1")
	e.ValueSet("allowed_authentication_secrets.0.all", "true")
	e.ValueSet("allowed_authentication_secrets.0.none", "false")
	e.ValueSet("allowed_authentication_secrets.0.secrets.#", "0")
	return e
}

func (e *ExternalAccessIntegrationResourceAssert) HasAllowedAuthenticationSecretsSecrets(expected ...string) *ExternalAccessIntegrationResourceAssert {
	e.ValueSet("allowed_authentication_secrets.#", "1")
	e.SetContainsExactlyStringValues("allowed_authentication_secrets.0.secrets", expected...)
	e.ValueSet("allowed_authentication_secrets.0.none", "false")
	e.ValueSet("allowed_authentication_secrets.0.all", "false")
	return e
}

func (e *ExternalAccessIntegrationResourceAssert) HasAllowedAuthenticationSecretsNotSet() *ExternalAccessIntegrationResourceAssert {
	e.ValueSet("allowed_authentication_secrets.#", "0")
	return e
}

func (e *ExternalAccessIntegrationResourceAssert) HasAllowedApiAuthenticationIntegrationsNone() *ExternalAccessIntegrationResourceAssert {
	e.ValueSet("allowed_api_authentication_integrations.#", "1")
	e.ValueSet("allowed_api_authentication_integrations.0.none", "true")
	e.ValueSet("allowed_api_authentication_integrations.0.integrations.#", "0")
	return e
}

func (e *ExternalAccessIntegrationResourceAssert) HasAllowedApiAuthenticationIntegrationsIntegrations(expected ...string) *ExternalAccessIntegrationResourceAssert {
	e.ValueSet("allowed_api_authentication_integrations.#", "1")
	e.SetContainsExactlyStringValues("allowed_api_authentication_integrations.0.integrations", expected...)
	e.ValueSet("allowed_api_authentication_integrations.0.none", "false")
	return e
}

func (e *ExternalAccessIntegrationResourceAssert) HasAllowedApiAuthenticationIntegrationsNotSet() *ExternalAccessIntegrationResourceAssert {
	e.ValueSet("allowed_api_authentication_integrations.#", "0")
	return e
}
