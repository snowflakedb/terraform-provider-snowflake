package sdk

import (
	"fmt"
)

const (
	SecurityIntegrationCategory                                     = "SECURITY"
	ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode = "AUTHORIZATION_CODE"
	ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials = "CLIENT_CREDENTIALS" //nolint:gosec
	ApiAuthenticationSecurityIntegrationOauthGrantJwtBearer         = "JWT_BEARER"
)

type Saml2SecurityIntegrationSaml2RequestedNameidFormatOption string

const (
	Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified                Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress               Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName            Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:X509SubjectName"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:1.1:nameid-format:WindowsDomainQualifiedName"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos                   Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:2.0:nameid-format:kerberos"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent                 Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent"
	Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient                  Saml2SecurityIntegrationSaml2RequestedNameidFormatOption = "urn:oasis:names:tc:SAML:2.0:nameid-format:transient"
)

var AllSaml2SecurityIntegrationSaml2RequestedNameidFormats = []Saml2SecurityIntegrationSaml2RequestedNameidFormatOption{
	Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent,
	Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient,
}

func ToSaml2SecurityIntegrationSaml2RequestedNameidFormatOption(s string) (Saml2SecurityIntegrationSaml2RequestedNameidFormatOption, error) {
	switch s {
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatUnspecified, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatEmailAddress, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatX509SubjectName, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatWindowsDomainQualifiedName, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatKerberos, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatPersistent, nil
	case string(Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient):
		return Saml2SecurityIntegrationSaml2RequestedNameidFormatTransient, nil
	default:
		return "", fmt.Errorf("invalid Saml2SecurityIntegrationSaml2RequestedNameidFormatOption: %s", s)
	}
}
