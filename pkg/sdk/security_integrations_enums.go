package sdk

import (
	"fmt"
	"strings"
)

const (
	SecurityIntegrationCategory                                     = "SECURITY"
	ApiAuthenticationSecurityIntegrationOauthGrantAuthorizationCode = "AUTHORIZATION_CODE"
	ApiAuthenticationSecurityIntegrationOauthGrantClientCredentials = "CLIENT_CREDENTIALS" //nolint:gosec
	ApiAuthenticationSecurityIntegrationOauthGrantJwtBearer         = "JWT_BEARER"
)

type ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption string

const (
	ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption = "CLIENT_SECRET_POST"
)

var AllApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption = []ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption{
	ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost,
}

func ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption(s string) (ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost):
		return ApiAuthenticationSecurityIntegrationOauthClientAuthMethodClientSecretPost, nil
	default:
		return "", fmt.Errorf("invalid ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption: %s", s)
	}
}

type ExternalOauthSecurityIntegrationTypeOption string

const (
	ExternalOauthSecurityIntegrationTypeOkta         ExternalOauthSecurityIntegrationTypeOption = "OKTA"
	ExternalOauthSecurityIntegrationTypeAzure        ExternalOauthSecurityIntegrationTypeOption = "AZURE"
	ExternalOauthSecurityIntegrationTypePingFederate ExternalOauthSecurityIntegrationTypeOption = "PING_FEDERATE"
	ExternalOauthSecurityIntegrationTypeCustom       ExternalOauthSecurityIntegrationTypeOption = "CUSTOM"
)

var AllExternalOauthSecurityIntegrationTypes = []ExternalOauthSecurityIntegrationTypeOption{
	ExternalOauthSecurityIntegrationTypeOkta,
	ExternalOauthSecurityIntegrationTypeAzure,
	ExternalOauthSecurityIntegrationTypePingFederate,
	ExternalOauthSecurityIntegrationTypeCustom,
}

func ToExternalOauthSecurityIntegrationTypeOption(s string) (ExternalOauthSecurityIntegrationTypeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ExternalOauthSecurityIntegrationTypeOkta):
		return ExternalOauthSecurityIntegrationTypeOkta, nil
	case string(ExternalOauthSecurityIntegrationTypeAzure):
		return ExternalOauthSecurityIntegrationTypeAzure, nil
	case string(ExternalOauthSecurityIntegrationTypePingFederate):
		return ExternalOauthSecurityIntegrationTypePingFederate, nil
	case string(ExternalOauthSecurityIntegrationTypeCustom):
		return ExternalOauthSecurityIntegrationTypeCustom, nil
	default:
		return "", fmt.Errorf("invalid ExternalOauthSecurityIntegrationTypeOption: %s", s)
	}
}

type ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption string

const (
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName    ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption = "LOGIN_NAME"
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption = "EMAIL_ADDRESS"
)

var AllExternalOauthSecurityIntegrationSnowflakeUserMappingAttributes = []ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption{
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName,
	ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress,
}

func ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption(s string) (ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName):
		return ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName, nil
	case string(ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress):
		return ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeEmailAddress, nil
	default:
		return "", fmt.Errorf("invalid ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption: %s", s)
	}
}

type ExternalOauthSecurityIntegrationAnyRoleModeOption string

const (
	ExternalOauthSecurityIntegrationAnyRoleModeDisable            ExternalOauthSecurityIntegrationAnyRoleModeOption = "DISABLE"
	ExternalOauthSecurityIntegrationAnyRoleModeEnable             ExternalOauthSecurityIntegrationAnyRoleModeOption = "ENABLE"
	ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege ExternalOauthSecurityIntegrationAnyRoleModeOption = "ENABLE_FOR_PRIVILEGE"
)

var AllExternalOauthSecurityIntegrationAnyRoleModes = []ExternalOauthSecurityIntegrationAnyRoleModeOption{
	ExternalOauthSecurityIntegrationAnyRoleModeDisable,
	ExternalOauthSecurityIntegrationAnyRoleModeEnable,
	ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege,
}

func ToExternalOauthSecurityIntegrationAnyRoleModeOption(s string) (ExternalOauthSecurityIntegrationAnyRoleModeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ExternalOauthSecurityIntegrationAnyRoleModeDisable):
		return ExternalOauthSecurityIntegrationAnyRoleModeDisable, nil
	case string(ExternalOauthSecurityIntegrationAnyRoleModeEnable):
		return ExternalOauthSecurityIntegrationAnyRoleModeEnable, nil
	case string(ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege):
		return ExternalOauthSecurityIntegrationAnyRoleModeEnableForPrivilege, nil
	default:
		return "", fmt.Errorf("invalid ExternalOauthSecurityIntegrationAnyRoleModeOption: %s", s)
	}
}

type OauthSecurityIntegrationUseSecondaryRolesOption string

const (
	OauthSecurityIntegrationUseSecondaryRolesImplicit OauthSecurityIntegrationUseSecondaryRolesOption = "IMPLICIT"
	OauthSecurityIntegrationUseSecondaryRolesNone     OauthSecurityIntegrationUseSecondaryRolesOption = "NONE"
)

var AllOauthSecurityIntegrationUseSecondaryRoles = []OauthSecurityIntegrationUseSecondaryRolesOption{
	OauthSecurityIntegrationUseSecondaryRolesImplicit,
	OauthSecurityIntegrationUseSecondaryRolesNone,
}

func ToOauthSecurityIntegrationUseSecondaryRolesOption(s string) (OauthSecurityIntegrationUseSecondaryRolesOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(OauthSecurityIntegrationUseSecondaryRolesImplicit):
		return OauthSecurityIntegrationUseSecondaryRolesImplicit, nil
	case string(OauthSecurityIntegrationUseSecondaryRolesNone):
		return OauthSecurityIntegrationUseSecondaryRolesNone, nil
	default:
		return "", fmt.Errorf("invalid OauthSecurityIntegrationUseSecondaryRolesOption: %s", s)
	}
}

type OauthSecurityIntegrationClientTypeOption string

const (
	OauthSecurityIntegrationClientTypePublic       OauthSecurityIntegrationClientTypeOption = "PUBLIC"
	OauthSecurityIntegrationClientTypeConfidential OauthSecurityIntegrationClientTypeOption = "CONFIDENTIAL"
)

var AllOauthSecurityIntegrationClientTypes = []OauthSecurityIntegrationClientTypeOption{
	OauthSecurityIntegrationClientTypePublic,
	OauthSecurityIntegrationClientTypeConfidential,
}

func ToOauthSecurityIntegrationClientTypeOption(s string) (OauthSecurityIntegrationClientTypeOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(OauthSecurityIntegrationClientTypePublic):
		return OauthSecurityIntegrationClientTypePublic, nil
	case string(OauthSecurityIntegrationClientTypeConfidential):
		return OauthSecurityIntegrationClientTypeConfidential, nil
	default:
		return "", fmt.Errorf("invalid OauthSecurityIntegrationClientTypeOption: %s", s)
	}
}

type OauthSecurityIntegrationClientOption string

const (
	OauthSecurityIntegrationClientLooker         OauthSecurityIntegrationClientOption = "LOOKER"
	OauthSecurityIntegrationClientTableauDesktop OauthSecurityIntegrationClientOption = "TABLEAU_DESKTOP"
	OauthSecurityIntegrationClientTableauServer  OauthSecurityIntegrationClientOption = "TABLEAU_SERVER"
)

var AllOauthSecurityIntegrationClients = []OauthSecurityIntegrationClientOption{
	OauthSecurityIntegrationClientLooker,
	OauthSecurityIntegrationClientTableauDesktop,
	OauthSecurityIntegrationClientTableauServer,
}

func ToOauthSecurityIntegrationClientOption(s string) (OauthSecurityIntegrationClientOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(OauthSecurityIntegrationClientLooker):
		return OauthSecurityIntegrationClientLooker, nil
	case string(OauthSecurityIntegrationClientTableauDesktop):
		return OauthSecurityIntegrationClientTableauDesktop, nil
	case string(OauthSecurityIntegrationClientTableauServer):
		return OauthSecurityIntegrationClientTableauServer, nil
	default:
		return "", fmt.Errorf("invalid OauthSecurityIntegrationClientOption: %s", s)
	}
}

type Saml2SecurityIntegrationSaml2ProviderOption string

const (
	Saml2SecurityIntegrationSaml2ProviderOkta   Saml2SecurityIntegrationSaml2ProviderOption = "OKTA"
	Saml2SecurityIntegrationSaml2ProviderAdfs   Saml2SecurityIntegrationSaml2ProviderOption = "ADFS"
	Saml2SecurityIntegrationSaml2ProviderCustom Saml2SecurityIntegrationSaml2ProviderOption = "CUSTOM"
)

var AllSaml2SecurityIntegrationSaml2Providers = []Saml2SecurityIntegrationSaml2ProviderOption{
	Saml2SecurityIntegrationSaml2ProviderOkta,
	Saml2SecurityIntegrationSaml2ProviderAdfs,
	Saml2SecurityIntegrationSaml2ProviderCustom,
}

func ToSaml2SecurityIntegrationSaml2ProviderOption(s string) (Saml2SecurityIntegrationSaml2ProviderOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(Saml2SecurityIntegrationSaml2ProviderOkta):
		return Saml2SecurityIntegrationSaml2ProviderOkta, nil
	case string(Saml2SecurityIntegrationSaml2ProviderAdfs):
		return Saml2SecurityIntegrationSaml2ProviderAdfs, nil
	case string(Saml2SecurityIntegrationSaml2ProviderCustom):
		return Saml2SecurityIntegrationSaml2ProviderCustom, nil
	default:
		return "", fmt.Errorf("invalid Saml2SecurityIntegrationSaml2ProviderOption: %s", s)
	}
}

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

type ScimSecurityIntegrationScimClientOption string

const (
	ScimSecurityIntegrationScimClientOkta    ScimSecurityIntegrationScimClientOption = "OKTA"
	ScimSecurityIntegrationScimClientAzure   ScimSecurityIntegrationScimClientOption = "AZURE"
	ScimSecurityIntegrationScimClientGeneric ScimSecurityIntegrationScimClientOption = "GENERIC"
)

var AllScimSecurityIntegrationScimClients = []ScimSecurityIntegrationScimClientOption{
	ScimSecurityIntegrationScimClientOkta,
	ScimSecurityIntegrationScimClientAzure,
	ScimSecurityIntegrationScimClientGeneric,
}

func ToScimSecurityIntegrationScimClientOption(s string) (ScimSecurityIntegrationScimClientOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ScimSecurityIntegrationScimClientOkta):
		return ScimSecurityIntegrationScimClientOkta, nil
	case string(ScimSecurityIntegrationScimClientAzure):
		return ScimSecurityIntegrationScimClientAzure, nil
	case string(ScimSecurityIntegrationScimClientGeneric):
		return ScimSecurityIntegrationScimClientGeneric, nil
	default:
		return "", fmt.Errorf("invalid ScimSecurityIntegrationScimClientOption: %s", s)
	}
}

type ScimSecurityIntegrationRunAsRoleOption string

const (
	ScimSecurityIntegrationRunAsRoleOktaProvisioner        ScimSecurityIntegrationRunAsRoleOption = "OKTA_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleAadProvisioner         ScimSecurityIntegrationRunAsRoleOption = "AAD_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleGenericScimProvisioner ScimSecurityIntegrationRunAsRoleOption = "GENERIC_SCIM_PROVISIONER"
)

var AllScimSecurityIntegrationRunAsRoles = []ScimSecurityIntegrationRunAsRoleOption{
	ScimSecurityIntegrationRunAsRoleOktaProvisioner,
	ScimSecurityIntegrationRunAsRoleAadProvisioner,
	ScimSecurityIntegrationRunAsRoleGenericScimProvisioner,
}

func ToScimSecurityIntegrationRunAsRoleOption(s string) (ScimSecurityIntegrationRunAsRoleOption, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(ScimSecurityIntegrationRunAsRoleOktaProvisioner):
		return ScimSecurityIntegrationRunAsRoleOktaProvisioner, nil
	case string(ScimSecurityIntegrationRunAsRoleAadProvisioner):
		return ScimSecurityIntegrationRunAsRoleAadProvisioner, nil
	case string(ScimSecurityIntegrationRunAsRoleGenericScimProvisioner):
		return ScimSecurityIntegrationRunAsRoleGenericScimProvisioner, nil
	default:
		return "", fmt.Errorf("invalid ScimSecurityIntegrationRunAsRoleOption: %s", s)
	}
}
