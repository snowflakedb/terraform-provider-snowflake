package sdk

import (
	"fmt"
	"strings"
)

func (r *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateExternalOauthSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateOauthForCustomClientsSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateSaml2SecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (r *CreateScimSecurityIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (s SecurityIntegrationProperty) GetName() string {
	return s.Name
}

func (s SecurityIntegrationProperty) GetDefault() string {
	return s.Default
}

func (s *SecurityIntegration) SubType() (string, error) {
	typeParts := strings.Split(s.IntegrationType, "-")
	if len(typeParts) < 2 {
		return "", fmt.Errorf("expected \"<type> - <subtype>\", got: %s", s.IntegrationType)
	}
	return strings.TrimSpace(typeParts[1]), nil
}

func securityIntegrationNetworkPolicyQuoted(id *AccountObjectIdentifier) *string {
	if id == nil {
		return nil
	}
	// TODO(SNOW-2236323): Use a proper generation option instead.
	// We need to use a custom parsing here, see SNOW-1833593 for more details.
	return Pointer(fmt.Sprintf("'%s'", id.FullyQualifiedName()))
}
