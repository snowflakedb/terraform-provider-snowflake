// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateExternalOauthSecurityIntegrationRequest(
	name AccountObjectIdentifier,
	Enabled bool,
	ExternalOauthType ExternalOauthSecurityIntegrationTypeOption,
	ExternalOauthIssuer string,
	ExternalOauthTokenUserMappingClaim []TokenUserMappingClaim,
	ExternalOauthSnowflakeUserMappingAttribute ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption,
) *CreateExternalOauthSecurityIntegrationRequest {
	s := CreateExternalOauthSecurityIntegrationRequest{}
	s.name = name
	s.Enabled = Enabled
	s.ExternalOauthType = ExternalOauthType
	s.ExternalOauthIssuer = ExternalOauthIssuer
	s.ExternalOauthTokenUserMappingClaim = ExternalOauthTokenUserMappingClaim
	s.ExternalOauthSnowflakeUserMappingAttribute = ExternalOauthSnowflakeUserMappingAttribute
	return &s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithOrReplace(OrReplace bool) *CreateExternalOauthSecurityIntegrationRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithIfNotExists(IfNotExists bool) *CreateExternalOauthSecurityIntegrationRequest {
	s.IfNotExists = &IfNotExists
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthJwsKeysUrl(ExternalOauthJwsKeysUrl []JwsKeysUrl) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthJwsKeysUrl = ExternalOauthJwsKeysUrl
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthBlockedRolesList(ExternalOauthBlockedRolesList BlockedRolesListRequest) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthBlockedRolesList = &ExternalOauthBlockedRolesList
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthAllowedRolesList(ExternalOauthAllowedRolesList AllowedRolesListRequest) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthAllowedRolesList = &ExternalOauthAllowedRolesList
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthRsaPublicKey(ExternalOauthRsaPublicKey string) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthRsaPublicKey = &ExternalOauthRsaPublicKey
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthRsaPublicKey2(ExternalOauthRsaPublicKey2 string) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthRsaPublicKey2 = &ExternalOauthRsaPublicKey2
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthAudienceList(ExternalOauthAudienceList AudienceListRequest) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthAudienceList = &ExternalOauthAudienceList
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthAnyRoleMode(ExternalOauthAnyRoleMode ExternalOauthSecurityIntegrationAnyRoleModeOption) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthAnyRoleMode = &ExternalOauthAnyRoleMode
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthScopeDelimiter(ExternalOauthScopeDelimiter string) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthScopeDelimiter = &ExternalOauthScopeDelimiter
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithExternalOauthScopeMappingAttribute(ExternalOauthScopeMappingAttribute string) *CreateExternalOauthSecurityIntegrationRequest {
	s.ExternalOauthScopeMappingAttribute = &ExternalOauthScopeMappingAttribute
	return s
}

func (s *CreateExternalOauthSecurityIntegrationRequest) WithComment(Comment string) *CreateExternalOauthSecurityIntegrationRequest {
	s.Comment = &Comment
	return s
}

func NewBlockedRolesListRequest() *BlockedRolesListRequest {
	return &BlockedRolesListRequest{}
}

func (s *BlockedRolesListRequest) WithBlockedRolesList(BlockedRolesList []AccountObjectIdentifier) *BlockedRolesListRequest {
	s.BlockedRolesList = BlockedRolesList
	return s
}

func NewAllowedRolesListRequest() *AllowedRolesListRequest {
	return &AllowedRolesListRequest{}
}

func (s *AllowedRolesListRequest) WithAllowedRolesList(AllowedRolesList []AccountObjectIdentifier) *AllowedRolesListRequest {
	s.AllowedRolesList = AllowedRolesList
	return s
}

func NewAudienceListRequest() *AudienceListRequest {
	return &AudienceListRequest{}
}

func (s *AudienceListRequest) WithAudienceList(AudienceList []AudienceListItem) *AudienceListRequest {
	s.AudienceList = AudienceList
	return s
}

func NewCreateOauthForPartnerApplicationsSecurityIntegrationRequest(
	name AccountObjectIdentifier,
	OauthClient OauthSecurityIntegrationClientOption,
) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s := CreateOauthForPartnerApplicationsSecurityIntegrationRequest{}
	s.name = name
	s.OauthClient = OauthClient
	return &s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithOrReplace(OrReplace bool) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithIfNotExists(IfNotExists bool) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.IfNotExists = &IfNotExists
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithOauthRedirectUri(OauthRedirectUri string) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.OauthRedirectUri = &OauthRedirectUri
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithEnabled(Enabled bool) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.Enabled = &Enabled
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithOauthIssueRefreshTokens(OauthIssueRefreshTokens bool) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.OauthIssueRefreshTokens = &OauthIssueRefreshTokens
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithOauthRefreshTokenValidity(OauthRefreshTokenValidity int) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.OauthRefreshTokenValidity = &OauthRefreshTokenValidity
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithOauthUseSecondaryRoles(OauthUseSecondaryRoles OauthSecurityIntegrationUseSecondaryRolesOption) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.OauthUseSecondaryRoles = &OauthUseSecondaryRoles
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithBlockedRolesList(BlockedRolesList BlockedRolesListRequest) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.BlockedRolesList = &BlockedRolesList
	return s
}

func (s *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) WithComment(Comment string) *CreateOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.Comment = &Comment
	return s
}

func NewCreateOauthForCustomClientsSecurityIntegrationRequest(
	name AccountObjectIdentifier,
	OauthClientType OauthSecurityIntegrationClientTypeOption,
	OauthRedirectUri string,
) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s := CreateOauthForCustomClientsSecurityIntegrationRequest{}
	s.name = name
	s.OauthClientType = OauthClientType
	s.OauthRedirectUri = OauthRedirectUri
	return &s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOrReplace(OrReplace bool) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithIfNotExists(IfNotExists bool) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.IfNotExists = &IfNotExists
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithEnabled(Enabled bool) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.Enabled = &Enabled
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOauthAllowNonTlsRedirectUri(OauthAllowNonTlsRedirectUri bool) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OauthAllowNonTlsRedirectUri = &OauthAllowNonTlsRedirectUri
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOauthEnforcePkce(OauthEnforcePkce bool) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OauthEnforcePkce = &OauthEnforcePkce
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOauthUseSecondaryRoles(OauthUseSecondaryRoles OauthSecurityIntegrationUseSecondaryRolesOption) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OauthUseSecondaryRoles = &OauthUseSecondaryRoles
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithPreAuthorizedRolesList(PreAuthorizedRolesList PreAuthorizedRolesListRequest) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.PreAuthorizedRolesList = &PreAuthorizedRolesList
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithBlockedRolesList(BlockedRolesList BlockedRolesListRequest) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.BlockedRolesList = &BlockedRolesList
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOauthIssueRefreshTokens(OauthIssueRefreshTokens bool) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OauthIssueRefreshTokens = &OauthIssueRefreshTokens
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOauthRefreshTokenValidity(OauthRefreshTokenValidity int) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OauthRefreshTokenValidity = &OauthRefreshTokenValidity
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithNetworkPolicy(NetworkPolicy AccountObjectIdentifier) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.NetworkPolicy = &NetworkPolicy
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOauthClientRsaPublicKey(OauthClientRsaPublicKey string) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OauthClientRsaPublicKey = &OauthClientRsaPublicKey
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithOauthClientRsaPublicKey2(OauthClientRsaPublicKey2 string) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.OauthClientRsaPublicKey2 = &OauthClientRsaPublicKey2
	return s
}

func (s *CreateOauthForCustomClientsSecurityIntegrationRequest) WithComment(Comment string) *CreateOauthForCustomClientsSecurityIntegrationRequest {
	s.Comment = &Comment
	return s
}

func NewPreAuthorizedRolesListRequest() *PreAuthorizedRolesListRequest {
	return &PreAuthorizedRolesListRequest{}
}

func (s *PreAuthorizedRolesListRequest) WithPreAuthorizedRolesList(PreAuthorizedRolesList []AccountObjectIdentifier) *PreAuthorizedRolesListRequest {
	s.PreAuthorizedRolesList = PreAuthorizedRolesList
	return s
}

func NewCreateSaml2SecurityIntegrationRequest(
	name AccountObjectIdentifier,
	Enabled bool,
	Saml2Issuer string,
	Saml2SsoUrl string,
	Saml2Provider string,
	Saml2X509Cert string,
) *CreateSaml2SecurityIntegrationRequest {
	s := CreateSaml2SecurityIntegrationRequest{}
	s.name = name
	s.Enabled = Enabled
	s.Saml2Issuer = Saml2Issuer
	s.Saml2SsoUrl = Saml2SsoUrl
	s.Saml2Provider = Saml2Provider
	s.Saml2X509Cert = Saml2X509Cert
	return &s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithOrReplace(OrReplace bool) *CreateSaml2SecurityIntegrationRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithIfNotExists(IfNotExists bool) *CreateSaml2SecurityIntegrationRequest {
	s.IfNotExists = &IfNotExists
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithAllowedUserDomains(AllowedUserDomains []UserDomain) *CreateSaml2SecurityIntegrationRequest {
	s.AllowedUserDomains = AllowedUserDomains
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithAllowedEmailPatterns(AllowedEmailPatterns []EmailPattern) *CreateSaml2SecurityIntegrationRequest {
	s.AllowedEmailPatterns = AllowedEmailPatterns
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2SpInitiatedLoginPageLabel(Saml2SpInitiatedLoginPageLabel string) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2SpInitiatedLoginPageLabel = &Saml2SpInitiatedLoginPageLabel
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2EnableSpInitiated(Saml2EnableSpInitiated bool) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2EnableSpInitiated = &Saml2EnableSpInitiated
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2SnowflakeX509Cert(Saml2SnowflakeX509Cert string) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2SnowflakeX509Cert = &Saml2SnowflakeX509Cert
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2SignRequest(Saml2SignRequest bool) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2SignRequest = &Saml2SignRequest
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2RequestedNameidFormat(Saml2RequestedNameidFormat string) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2RequestedNameidFormat = &Saml2RequestedNameidFormat
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2PostLogoutRedirectUrl(Saml2PostLogoutRedirectUrl string) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2PostLogoutRedirectUrl = &Saml2PostLogoutRedirectUrl
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2ForceAuthn(Saml2ForceAuthn bool) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2ForceAuthn = &Saml2ForceAuthn
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2SnowflakeIssuerUrl(Saml2SnowflakeIssuerUrl string) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2SnowflakeIssuerUrl = &Saml2SnowflakeIssuerUrl
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithSaml2SnowflakeAcsUrl(Saml2SnowflakeAcsUrl string) *CreateSaml2SecurityIntegrationRequest {
	s.Saml2SnowflakeAcsUrl = &Saml2SnowflakeAcsUrl
	return s
}

func (s *CreateSaml2SecurityIntegrationRequest) WithComment(Comment string) *CreateSaml2SecurityIntegrationRequest {
	s.Comment = &Comment
	return s
}

func NewCreateScimSecurityIntegrationRequest(
	name AccountObjectIdentifier,
	Enabled bool,
	ScimClient ScimSecurityIntegrationScimClientOption,
	RunAsRole ScimSecurityIntegrationRunAsRoleOption,
) *CreateScimSecurityIntegrationRequest {
	s := CreateScimSecurityIntegrationRequest{}
	s.name = name
	s.Enabled = Enabled
	s.ScimClient = ScimClient
	s.RunAsRole = RunAsRole
	return &s
}

func (s *CreateScimSecurityIntegrationRequest) WithOrReplace(OrReplace bool) *CreateScimSecurityIntegrationRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateScimSecurityIntegrationRequest) WithIfNotExists(IfNotExists bool) *CreateScimSecurityIntegrationRequest {
	s.IfNotExists = &IfNotExists
	return s
}

func (s *CreateScimSecurityIntegrationRequest) WithNetworkPolicy(NetworkPolicy AccountObjectIdentifier) *CreateScimSecurityIntegrationRequest {
	s.NetworkPolicy = &NetworkPolicy
	return s
}

func (s *CreateScimSecurityIntegrationRequest) WithSyncPassword(SyncPassword bool) *CreateScimSecurityIntegrationRequest {
	s.SyncPassword = &SyncPassword
	return s
}

func (s *CreateScimSecurityIntegrationRequest) WithComment(Comment string) *CreateScimSecurityIntegrationRequest {
	s.Comment = &Comment
	return s
}

func NewAlterExternalOauthSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterExternalOauthSecurityIntegrationRequest {
	s := AlterExternalOauthSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterExternalOauthSecurityIntegrationRequest) WithIfExists(IfExists bool) *AlterExternalOauthSecurityIntegrationRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterExternalOauthSecurityIntegrationRequest) WithSetTags(SetTags []TagAssociation) *AlterExternalOauthSecurityIntegrationRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterExternalOauthSecurityIntegrationRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterExternalOauthSecurityIntegrationRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterExternalOauthSecurityIntegrationRequest) WithSet(Set ExternalOauthIntegrationSetRequest) *AlterExternalOauthSecurityIntegrationRequest {
	s.Set = &Set
	return s
}

func (s *AlterExternalOauthSecurityIntegrationRequest) WithUnset(Unset ExternalOauthIntegrationUnsetRequest) *AlterExternalOauthSecurityIntegrationRequest {
	s.Unset = &Unset
	return s
}

func NewExternalOauthIntegrationSetRequest() *ExternalOauthIntegrationSetRequest {
	return &ExternalOauthIntegrationSetRequest{}
}

func (s *ExternalOauthIntegrationSetRequest) WithEnabled(Enabled bool) *ExternalOauthIntegrationSetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthType(ExternalOauthType ExternalOauthSecurityIntegrationTypeOption) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthType = &ExternalOauthType
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthIssuer(ExternalOauthIssuer string) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthIssuer = &ExternalOauthIssuer
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthTokenUserMappingClaim(ExternalOauthTokenUserMappingClaim []TokenUserMappingClaim) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthTokenUserMappingClaim = ExternalOauthTokenUserMappingClaim
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthSnowflakeUserMappingAttribute(ExternalOauthSnowflakeUserMappingAttribute ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthSnowflakeUserMappingAttribute = &ExternalOauthSnowflakeUserMappingAttribute
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthJwsKeysUrl(ExternalOauthJwsKeysUrl []JwsKeysUrl) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthJwsKeysUrl = ExternalOauthJwsKeysUrl
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthBlockedRolesList(ExternalOauthBlockedRolesList BlockedRolesListRequest) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthBlockedRolesList = &ExternalOauthBlockedRolesList
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthAllowedRolesList(ExternalOauthAllowedRolesList AllowedRolesListRequest) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthAllowedRolesList = &ExternalOauthAllowedRolesList
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthRsaPublicKey(ExternalOauthRsaPublicKey string) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthRsaPublicKey = &ExternalOauthRsaPublicKey
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthRsaPublicKey2(ExternalOauthRsaPublicKey2 string) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthRsaPublicKey2 = &ExternalOauthRsaPublicKey2
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthAudienceList(ExternalOauthAudienceList AudienceListRequest) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthAudienceList = &ExternalOauthAudienceList
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthAnyRoleMode(ExternalOauthAnyRoleMode ExternalOauthSecurityIntegrationAnyRoleModeOption) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthAnyRoleMode = &ExternalOauthAnyRoleMode
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithExternalOauthScopeDelimiter(ExternalOauthScopeDelimiter string) *ExternalOauthIntegrationSetRequest {
	s.ExternalOauthScopeDelimiter = &ExternalOauthScopeDelimiter
	return s
}

func (s *ExternalOauthIntegrationSetRequest) WithComment(Comment string) *ExternalOauthIntegrationSetRequest {
	s.Comment = &Comment
	return s
}

func NewExternalOauthIntegrationUnsetRequest() *ExternalOauthIntegrationUnsetRequest {
	return &ExternalOauthIntegrationUnsetRequest{}
}

func (s *ExternalOauthIntegrationUnsetRequest) WithEnabled(Enabled bool) *ExternalOauthIntegrationUnsetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *ExternalOauthIntegrationUnsetRequest) WithExternalOauthAudienceList(ExternalOauthAudienceList bool) *ExternalOauthIntegrationUnsetRequest {
	s.ExternalOauthAudienceList = &ExternalOauthAudienceList
	return s
}

func NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterOauthForPartnerApplicationsSecurityIntegrationRequest {
	s := AlterOauthForPartnerApplicationsSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) WithIfExists(IfExists bool) *AlterOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) WithSetTags(SetTags []TagAssociation) *AlterOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) WithSet(Set OauthForPartnerApplicationsIntegrationSetRequest) *AlterOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.Set = &Set
	return s
}

func (s *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) WithUnset(Unset OauthForPartnerApplicationsIntegrationUnsetRequest) *AlterOauthForPartnerApplicationsSecurityIntegrationRequest {
	s.Unset = &Unset
	return s
}

func NewOauthForPartnerApplicationsIntegrationSetRequest() *OauthForPartnerApplicationsIntegrationSetRequest {
	return &OauthForPartnerApplicationsIntegrationSetRequest{}
}

func (s *OauthForPartnerApplicationsIntegrationSetRequest) WithEnabled(Enabled bool) *OauthForPartnerApplicationsIntegrationSetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *OauthForPartnerApplicationsIntegrationSetRequest) WithOauthIssueRefreshTokens(OauthIssueRefreshTokens bool) *OauthForPartnerApplicationsIntegrationSetRequest {
	s.OauthIssueRefreshTokens = &OauthIssueRefreshTokens
	return s
}

func (s *OauthForPartnerApplicationsIntegrationSetRequest) WithOauthRedirectUri(OauthRedirectUri string) *OauthForPartnerApplicationsIntegrationSetRequest {
	s.OauthRedirectUri = &OauthRedirectUri
	return s
}

func (s *OauthForPartnerApplicationsIntegrationSetRequest) WithOauthRefreshTokenValidity(OauthRefreshTokenValidity int) *OauthForPartnerApplicationsIntegrationSetRequest {
	s.OauthRefreshTokenValidity = &OauthRefreshTokenValidity
	return s
}

func (s *OauthForPartnerApplicationsIntegrationSetRequest) WithOauthUseSecondaryRoles(OauthUseSecondaryRoles OauthSecurityIntegrationUseSecondaryRolesOption) *OauthForPartnerApplicationsIntegrationSetRequest {
	s.OauthUseSecondaryRoles = &OauthUseSecondaryRoles
	return s
}

func (s *OauthForPartnerApplicationsIntegrationSetRequest) WithBlockedRolesList(BlockedRolesList BlockedRolesListRequest) *OauthForPartnerApplicationsIntegrationSetRequest {
	s.BlockedRolesList = &BlockedRolesList
	return s
}

func (s *OauthForPartnerApplicationsIntegrationSetRequest) WithComment(Comment string) *OauthForPartnerApplicationsIntegrationSetRequest {
	s.Comment = &Comment
	return s
}

func NewOauthForPartnerApplicationsIntegrationUnsetRequest() *OauthForPartnerApplicationsIntegrationUnsetRequest {
	return &OauthForPartnerApplicationsIntegrationUnsetRequest{}
}

func (s *OauthForPartnerApplicationsIntegrationUnsetRequest) WithEnabled(Enabled bool) *OauthForPartnerApplicationsIntegrationUnsetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *OauthForPartnerApplicationsIntegrationUnsetRequest) WithOauthUseSecondaryRoles(OauthUseSecondaryRoles bool) *OauthForPartnerApplicationsIntegrationUnsetRequest {
	s.OauthUseSecondaryRoles = &OauthUseSecondaryRoles
	return s
}

func NewAlterOauthForCustomClientsSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterOauthForCustomClientsSecurityIntegrationRequest {
	s := AlterOauthForCustomClientsSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterOauthForCustomClientsSecurityIntegrationRequest) WithIfExists(IfExists bool) *AlterOauthForCustomClientsSecurityIntegrationRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterOauthForCustomClientsSecurityIntegrationRequest) WithSetTags(SetTags []TagAssociation) *AlterOauthForCustomClientsSecurityIntegrationRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterOauthForCustomClientsSecurityIntegrationRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterOauthForCustomClientsSecurityIntegrationRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterOauthForCustomClientsSecurityIntegrationRequest) WithSet(Set OauthForCustomClientsIntegrationSetRequest) *AlterOauthForCustomClientsSecurityIntegrationRequest {
	s.Set = &Set
	return s
}

func (s *AlterOauthForCustomClientsSecurityIntegrationRequest) WithUnset(Unset OauthForCustomClientsIntegrationUnsetRequest) *AlterOauthForCustomClientsSecurityIntegrationRequest {
	s.Unset = &Unset
	return s
}

func NewOauthForCustomClientsIntegrationSetRequest() *OauthForCustomClientsIntegrationSetRequest {
	return &OauthForCustomClientsIntegrationSetRequest{}
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithEnabled(Enabled bool) *OauthForCustomClientsIntegrationSetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthRedirectUri(OauthRedirectUri string) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthRedirectUri = &OauthRedirectUri
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthAllowNonTlsRedirectUri(OauthAllowNonTlsRedirectUri bool) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthAllowNonTlsRedirectUri = &OauthAllowNonTlsRedirectUri
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthEnforcePkce(OauthEnforcePkce bool) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthEnforcePkce = &OauthEnforcePkce
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithPreAuthorizedRolesList(PreAuthorizedRolesList PreAuthorizedRolesListRequest) *OauthForCustomClientsIntegrationSetRequest {
	s.PreAuthorizedRolesList = &PreAuthorizedRolesList
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithBlockedRolesList(BlockedRolesList BlockedRolesListRequest) *OauthForCustomClientsIntegrationSetRequest {
	s.BlockedRolesList = &BlockedRolesList
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthIssueRefreshTokens(OauthIssueRefreshTokens bool) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthIssueRefreshTokens = &OauthIssueRefreshTokens
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthRefreshTokenValidity(OauthRefreshTokenValidity int) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthRefreshTokenValidity = &OauthRefreshTokenValidity
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthUseSecondaryRoles(OauthUseSecondaryRoles OauthSecurityIntegrationUseSecondaryRolesOption) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthUseSecondaryRoles = &OauthUseSecondaryRoles
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithNetworkPolicy(NetworkPolicy AccountObjectIdentifier) *OauthForCustomClientsIntegrationSetRequest {
	s.NetworkPolicy = &NetworkPolicy
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthClientRsaPublicKey(OauthClientRsaPublicKey string) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthClientRsaPublicKey = &OauthClientRsaPublicKey
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithOauthClientRsaPublicKey2(OauthClientRsaPublicKey2 string) *OauthForCustomClientsIntegrationSetRequest {
	s.OauthClientRsaPublicKey2 = &OauthClientRsaPublicKey2
	return s
}

func (s *OauthForCustomClientsIntegrationSetRequest) WithComment(Comment string) *OauthForCustomClientsIntegrationSetRequest {
	s.Comment = &Comment
	return s
}

func NewOauthForCustomClientsIntegrationUnsetRequest() *OauthForCustomClientsIntegrationUnsetRequest {
	return &OauthForCustomClientsIntegrationUnsetRequest{}
}

func (s *OauthForCustomClientsIntegrationUnsetRequest) WithEnabled(Enabled bool) *OauthForCustomClientsIntegrationUnsetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *OauthForCustomClientsIntegrationUnsetRequest) WithNetworkPolicy(NetworkPolicy bool) *OauthForCustomClientsIntegrationUnsetRequest {
	s.NetworkPolicy = &NetworkPolicy
	return s
}

func (s *OauthForCustomClientsIntegrationUnsetRequest) WithOauthClientRsaPublicKey(OauthClientRsaPublicKey bool) *OauthForCustomClientsIntegrationUnsetRequest {
	s.OauthClientRsaPublicKey = &OauthClientRsaPublicKey
	return s
}

func (s *OauthForCustomClientsIntegrationUnsetRequest) WithOauthClientRsaPublicKey2(OauthClientRsaPublicKey2 bool) *OauthForCustomClientsIntegrationUnsetRequest {
	s.OauthClientRsaPublicKey2 = &OauthClientRsaPublicKey2
	return s
}

func (s *OauthForCustomClientsIntegrationUnsetRequest) WithOauthUseSecondaryRoles(OauthUseSecondaryRoles bool) *OauthForCustomClientsIntegrationUnsetRequest {
	s.OauthUseSecondaryRoles = &OauthUseSecondaryRoles
	return s
}

func NewAlterSaml2SecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterSaml2SecurityIntegrationRequest {
	s := AlterSaml2SecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterSaml2SecurityIntegrationRequest) WithIfExists(IfExists bool) *AlterSaml2SecurityIntegrationRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterSaml2SecurityIntegrationRequest) WithSetTags(SetTags []TagAssociation) *AlterSaml2SecurityIntegrationRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterSaml2SecurityIntegrationRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterSaml2SecurityIntegrationRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterSaml2SecurityIntegrationRequest) WithSet(Set Saml2IntegrationSetRequest) *AlterSaml2SecurityIntegrationRequest {
	s.Set = &Set
	return s
}

func (s *AlterSaml2SecurityIntegrationRequest) WithUnset(Unset Saml2IntegrationUnsetRequest) *AlterSaml2SecurityIntegrationRequest {
	s.Unset = &Unset
	return s
}

func (s *AlterSaml2SecurityIntegrationRequest) WithRefreshSaml2SnowflakePrivateKey(RefreshSaml2SnowflakePrivateKey bool) *AlterSaml2SecurityIntegrationRequest {
	s.RefreshSaml2SnowflakePrivateKey = &RefreshSaml2SnowflakePrivateKey
	return s
}

func NewSaml2IntegrationSetRequest() *Saml2IntegrationSetRequest {
	return &Saml2IntegrationSetRequest{}
}

func (s *Saml2IntegrationSetRequest) WithEnabled(Enabled bool) *Saml2IntegrationSetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2Issuer(Saml2Issuer string) *Saml2IntegrationSetRequest {
	s.Saml2Issuer = &Saml2Issuer
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2SsoUrl(Saml2SsoUrl string) *Saml2IntegrationSetRequest {
	s.Saml2SsoUrl = &Saml2SsoUrl
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2Provider(Saml2Provider string) *Saml2IntegrationSetRequest {
	s.Saml2Provider = &Saml2Provider
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2X509Cert(Saml2X509Cert string) *Saml2IntegrationSetRequest {
	s.Saml2X509Cert = &Saml2X509Cert
	return s
}

func (s *Saml2IntegrationSetRequest) WithAllowedUserDomains(AllowedUserDomains []UserDomain) *Saml2IntegrationSetRequest {
	s.AllowedUserDomains = AllowedUserDomains
	return s
}

func (s *Saml2IntegrationSetRequest) WithAllowedEmailPatterns(AllowedEmailPatterns []EmailPattern) *Saml2IntegrationSetRequest {
	s.AllowedEmailPatterns = AllowedEmailPatterns
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2SpInitiatedLoginPageLabel(Saml2SpInitiatedLoginPageLabel string) *Saml2IntegrationSetRequest {
	s.Saml2SpInitiatedLoginPageLabel = &Saml2SpInitiatedLoginPageLabel
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2EnableSpInitiated(Saml2EnableSpInitiated bool) *Saml2IntegrationSetRequest {
	s.Saml2EnableSpInitiated = &Saml2EnableSpInitiated
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2SnowflakeX509Cert(Saml2SnowflakeX509Cert string) *Saml2IntegrationSetRequest {
	s.Saml2SnowflakeX509Cert = &Saml2SnowflakeX509Cert
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2SignRequest(Saml2SignRequest bool) *Saml2IntegrationSetRequest {
	s.Saml2SignRequest = &Saml2SignRequest
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2RequestedNameidFormat(Saml2RequestedNameidFormat string) *Saml2IntegrationSetRequest {
	s.Saml2RequestedNameidFormat = &Saml2RequestedNameidFormat
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2PostLogoutRedirectUrl(Saml2PostLogoutRedirectUrl string) *Saml2IntegrationSetRequest {
	s.Saml2PostLogoutRedirectUrl = &Saml2PostLogoutRedirectUrl
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2ForceAuthn(Saml2ForceAuthn bool) *Saml2IntegrationSetRequest {
	s.Saml2ForceAuthn = &Saml2ForceAuthn
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2SnowflakeIssuerUrl(Saml2SnowflakeIssuerUrl string) *Saml2IntegrationSetRequest {
	s.Saml2SnowflakeIssuerUrl = &Saml2SnowflakeIssuerUrl
	return s
}

func (s *Saml2IntegrationSetRequest) WithSaml2SnowflakeAcsUrl(Saml2SnowflakeAcsUrl string) *Saml2IntegrationSetRequest {
	s.Saml2SnowflakeAcsUrl = &Saml2SnowflakeAcsUrl
	return s
}

func (s *Saml2IntegrationSetRequest) WithComment(Comment string) *Saml2IntegrationSetRequest {
	s.Comment = &Comment
	return s
}

func NewSaml2IntegrationUnsetRequest() *Saml2IntegrationUnsetRequest {
	return &Saml2IntegrationUnsetRequest{}
}

func (s *Saml2IntegrationUnsetRequest) WithSaml2ForceAuthn(Saml2ForceAuthn bool) *Saml2IntegrationUnsetRequest {
	s.Saml2ForceAuthn = &Saml2ForceAuthn
	return s
}

func (s *Saml2IntegrationUnsetRequest) WithSaml2RequestedNameidFormat(Saml2RequestedNameidFormat bool) *Saml2IntegrationUnsetRequest {
	s.Saml2RequestedNameidFormat = &Saml2RequestedNameidFormat
	return s
}

func (s *Saml2IntegrationUnsetRequest) WithSaml2PostLogoutRedirectUrl(Saml2PostLogoutRedirectUrl bool) *Saml2IntegrationUnsetRequest {
	s.Saml2PostLogoutRedirectUrl = &Saml2PostLogoutRedirectUrl
	return s
}

func (s *Saml2IntegrationUnsetRequest) WithComment(Comment bool) *Saml2IntegrationUnsetRequest {
	s.Comment = &Comment
	return s
}

func NewAlterScimSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterScimSecurityIntegrationRequest {
	s := AlterScimSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterScimSecurityIntegrationRequest) WithIfExists(IfExists bool) *AlterScimSecurityIntegrationRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterScimSecurityIntegrationRequest) WithSetTags(SetTags []TagAssociation) *AlterScimSecurityIntegrationRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterScimSecurityIntegrationRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterScimSecurityIntegrationRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterScimSecurityIntegrationRequest) WithSet(Set ScimIntegrationSetRequest) *AlterScimSecurityIntegrationRequest {
	s.Set = &Set
	return s
}

func (s *AlterScimSecurityIntegrationRequest) WithUnset(Unset ScimIntegrationUnsetRequest) *AlterScimSecurityIntegrationRequest {
	s.Unset = &Unset
	return s
}

func NewScimIntegrationSetRequest() *ScimIntegrationSetRequest {
	return &ScimIntegrationSetRequest{}
}

func (s *ScimIntegrationSetRequest) WithEnabled(Enabled bool) *ScimIntegrationSetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *ScimIntegrationSetRequest) WithNetworkPolicy(NetworkPolicy AccountObjectIdentifier) *ScimIntegrationSetRequest {
	s.NetworkPolicy = &NetworkPolicy
	return s
}

func (s *ScimIntegrationSetRequest) WithSyncPassword(SyncPassword bool) *ScimIntegrationSetRequest {
	s.SyncPassword = &SyncPassword
	return s
}

func (s *ScimIntegrationSetRequest) WithComment(Comment string) *ScimIntegrationSetRequest {
	s.Comment = &Comment
	return s
}

func NewScimIntegrationUnsetRequest() *ScimIntegrationUnsetRequest {
	return &ScimIntegrationUnsetRequest{}
}

func (s *ScimIntegrationUnsetRequest) WithEnabled(Enabled bool) *ScimIntegrationUnsetRequest {
	s.Enabled = &Enabled
	return s
}

func (s *ScimIntegrationUnsetRequest) WithNetworkPolicy(NetworkPolicy bool) *ScimIntegrationUnsetRequest {
	s.NetworkPolicy = &NetworkPolicy
	return s
}

func (s *ScimIntegrationUnsetRequest) WithSyncPassword(SyncPassword bool) *ScimIntegrationUnsetRequest {
	s.SyncPassword = &SyncPassword
	return s
}

func (s *ScimIntegrationUnsetRequest) WithComment(Comment bool) *ScimIntegrationUnsetRequest {
	s.Comment = &Comment
	return s
}

func NewDropSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *DropSecurityIntegrationRequest {
	s := DropSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *DropSecurityIntegrationRequest) WithIfExists(IfExists bool) *DropSecurityIntegrationRequest {
	s.IfExists = &IfExists
	return s
}

func NewDescribeSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *DescribeSecurityIntegrationRequest {
	s := DescribeSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func NewShowSecurityIntegrationRequest() *ShowSecurityIntegrationRequest {
	return &ShowSecurityIntegrationRequest{}
}

func (s *ShowSecurityIntegrationRequest) WithLike(Like Like) *ShowSecurityIntegrationRequest {
	s.Like = &Like
	return s
}
