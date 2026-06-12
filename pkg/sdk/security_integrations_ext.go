package sdk

import (
	"fmt"
	"strings"
)

func (opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions) additionalValidations() error {
	if opts.OauthClient == OauthSecurityIntegrationClientOptionLooker && opts.OauthRedirectUri == nil {
		return NewError("OauthRedirectUri is required when OauthClient is LOOKER")
	}
	return nil
}

func (opts *CreateScimSecurityIntegrationOptions) additionalValidations() error {
	if opts.ScimClient == ScimSecurityIntegrationScimClientOptionAzure && opts.SyncPassword != nil {
		return NewError("SyncPassword is not supported for Azure scim client")
	}
	return nil
}

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
