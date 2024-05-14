// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateSAML2SecurityIntegrationRequest(
	name AccountObjectIdentifier,
	Enabled bool,
	Saml2Issuer string,
	Saml2SsoUrl string,
	Saml2Provider string,
	Saml2X509Cert string,
) *CreateSAML2SecurityIntegrationRequest {
	s := CreateSAML2SecurityIntegrationRequest{}
	s.name = name
	s.Enabled = Enabled
	s.Saml2Issuer = Saml2Issuer
	s.Saml2SsoUrl = Saml2SsoUrl
	s.Saml2Provider = Saml2Provider
	s.Saml2X509Cert = Saml2X509Cert
	return &s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithOrReplace(OrReplace *bool) *CreateSAML2SecurityIntegrationRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithIfNotExists(IfNotExists *bool) *CreateSAML2SecurityIntegrationRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithAllowedUserDomains(AllowedUserDomains []UserDomain) *CreateSAML2SecurityIntegrationRequest {
	s.AllowedUserDomains = AllowedUserDomains
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithAllowedEmailPatterns(AllowedEmailPatterns []EmailPattern) *CreateSAML2SecurityIntegrationRequest {
	s.AllowedEmailPatterns = AllowedEmailPatterns
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2SpInitiatedLoginPageLabel(Saml2SpInitiatedLoginPageLabel *string) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2SpInitiatedLoginPageLabel = Saml2SpInitiatedLoginPageLabel
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2EnableSpInitiated(Saml2EnableSpInitiated *bool) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2EnableSpInitiated = Saml2EnableSpInitiated
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2SnowflakeX509Cert(Saml2SnowflakeX509Cert *string) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2SnowflakeX509Cert = Saml2SnowflakeX509Cert
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2SignRequest(Saml2SignRequest *bool) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2SignRequest = Saml2SignRequest
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2RequestedNameidFormat(Saml2RequestedNameidFormat *string) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2RequestedNameidFormat = Saml2RequestedNameidFormat
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2PostLogoutRedirectUrl(Saml2PostLogoutRedirectUrl *string) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2PostLogoutRedirectUrl = Saml2PostLogoutRedirectUrl
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2ForceAuthn(Saml2ForceAuthn *bool) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2ForceAuthn = Saml2ForceAuthn
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2SnowflakeIssuerUrl(Saml2SnowflakeIssuerUrl *string) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2SnowflakeIssuerUrl = Saml2SnowflakeIssuerUrl
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithSaml2SnowflakeAcsUrl(Saml2SnowflakeAcsUrl *string) *CreateSAML2SecurityIntegrationRequest {
	s.Saml2SnowflakeAcsUrl = Saml2SnowflakeAcsUrl
	return s
}

func (s *CreateSAML2SecurityIntegrationRequest) WithComment(Comment *string) *CreateSAML2SecurityIntegrationRequest {
	s.Comment = Comment
	return s
}

func NewCreateSCIMSecurityIntegrationRequest(
	name AccountObjectIdentifier,
	Enabled bool,
	ScimClient string,
	RunAsRole string,
) *CreateSCIMSecurityIntegrationRequest {
	s := CreateSCIMSecurityIntegrationRequest{}
	s.name = name
	s.Enabled = Enabled
	s.ScimClient = ScimClient
	s.RunAsRole = RunAsRole
	return &s
}

func (s *CreateSCIMSecurityIntegrationRequest) WithOrReplace(OrReplace *bool) *CreateSCIMSecurityIntegrationRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateSCIMSecurityIntegrationRequest) WithIfNotExists(IfNotExists *bool) *CreateSCIMSecurityIntegrationRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateSCIMSecurityIntegrationRequest) WithNetworkPolicy(NetworkPolicy *AccountObjectIdentifier) *CreateSCIMSecurityIntegrationRequest {
	s.NetworkPolicy = NetworkPolicy
	return s
}

func (s *CreateSCIMSecurityIntegrationRequest) WithSyncPassword(SyncPassword *bool) *CreateSCIMSecurityIntegrationRequest {
	s.SyncPassword = SyncPassword
	return s
}

func (s *CreateSCIMSecurityIntegrationRequest) WithComment(Comment *string) *CreateSCIMSecurityIntegrationRequest {
	s.Comment = Comment
	return s
}

func NewAlterSAML2IntegrationSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterSAML2IntegrationSecurityIntegrationRequest {
	s := AlterSAML2IntegrationSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterSAML2IntegrationSecurityIntegrationRequest) WithIfExists(IfExists *bool) *AlterSAML2IntegrationSecurityIntegrationRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterSAML2IntegrationSecurityIntegrationRequest) WithSetTags(SetTags []TagAssociation) *AlterSAML2IntegrationSecurityIntegrationRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterSAML2IntegrationSecurityIntegrationRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterSAML2IntegrationSecurityIntegrationRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterSAML2IntegrationSecurityIntegrationRequest) WithSet(Set *SAML2IntegrationSetRequest) *AlterSAML2IntegrationSecurityIntegrationRequest {
	s.Set = Set
	return s
}

func (s *AlterSAML2IntegrationSecurityIntegrationRequest) WithUnset(Unset *SAML2IntegrationUnsetRequest) *AlterSAML2IntegrationSecurityIntegrationRequest {
	s.Unset = Unset
	return s
}

func (s *AlterSAML2IntegrationSecurityIntegrationRequest) WithRefreshSaml2SnowflakePrivateKey(RefreshSaml2SnowflakePrivateKey *bool) *AlterSAML2IntegrationSecurityIntegrationRequest {
	s.RefreshSaml2SnowflakePrivateKey = RefreshSaml2SnowflakePrivateKey
	return s
}

func NewSAML2IntegrationSetRequest() *SAML2IntegrationSetRequest {
	return &SAML2IntegrationSetRequest{}
}

func (s *SAML2IntegrationSetRequest) WithEnabled(Enabled *bool) *SAML2IntegrationSetRequest {
	s.Enabled = Enabled
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2Issuer(Saml2Issuer *string) *SAML2IntegrationSetRequest {
	s.Saml2Issuer = Saml2Issuer
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2SsoUrl(Saml2SsoUrl *string) *SAML2IntegrationSetRequest {
	s.Saml2SsoUrl = Saml2SsoUrl
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2Provider(Saml2Provider *string) *SAML2IntegrationSetRequest {
	s.Saml2Provider = Saml2Provider
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2X509Cert(Saml2X509Cert *string) *SAML2IntegrationSetRequest {
	s.Saml2X509Cert = Saml2X509Cert
	return s
}

func (s *SAML2IntegrationSetRequest) WithAllowedUserDomains(AllowedUserDomains []UserDomain) *SAML2IntegrationSetRequest {
	s.AllowedUserDomains = AllowedUserDomains
	return s
}

func (s *SAML2IntegrationSetRequest) WithAllowedEmailPatterns(AllowedEmailPatterns []EmailPattern) *SAML2IntegrationSetRequest {
	s.AllowedEmailPatterns = AllowedEmailPatterns
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2SpInitiatedLoginPageLabel(Saml2SpInitiatedLoginPageLabel *string) *SAML2IntegrationSetRequest {
	s.Saml2SpInitiatedLoginPageLabel = Saml2SpInitiatedLoginPageLabel
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2EnableSpInitiated(Saml2EnableSpInitiated *bool) *SAML2IntegrationSetRequest {
	s.Saml2EnableSpInitiated = Saml2EnableSpInitiated
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2SnowflakeX509Cert(Saml2SnowflakeX509Cert *string) *SAML2IntegrationSetRequest {
	s.Saml2SnowflakeX509Cert = Saml2SnowflakeX509Cert
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2SignRequest(Saml2SignRequest *bool) *SAML2IntegrationSetRequest {
	s.Saml2SignRequest = Saml2SignRequest
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2RequestedNameidFormat(Saml2RequestedNameidFormat *string) *SAML2IntegrationSetRequest {
	s.Saml2RequestedNameidFormat = Saml2RequestedNameidFormat
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2PostLogoutRedirectUrl(Saml2PostLogoutRedirectUrl *string) *SAML2IntegrationSetRequest {
	s.Saml2PostLogoutRedirectUrl = Saml2PostLogoutRedirectUrl
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2ForceAuthn(Saml2ForceAuthn *bool) *SAML2IntegrationSetRequest {
	s.Saml2ForceAuthn = Saml2ForceAuthn
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2SnowflakeIssuerUrl(Saml2SnowflakeIssuerUrl *string) *SAML2IntegrationSetRequest {
	s.Saml2SnowflakeIssuerUrl = Saml2SnowflakeIssuerUrl
	return s
}

func (s *SAML2IntegrationSetRequest) WithSaml2SnowflakeAcsUrl(Saml2SnowflakeAcsUrl *string) *SAML2IntegrationSetRequest {
	s.Saml2SnowflakeAcsUrl = Saml2SnowflakeAcsUrl
	return s
}

func (s *SAML2IntegrationSetRequest) WithComment(Comment *string) *SAML2IntegrationSetRequest {
	s.Comment = Comment
	return s
}

func NewSAML2IntegrationUnsetRequest() *SAML2IntegrationUnsetRequest {
	return &SAML2IntegrationUnsetRequest{}
}

func (s *SAML2IntegrationUnsetRequest) WithEnabled(Enabled *bool) *SAML2IntegrationUnsetRequest {
	s.Enabled = Enabled
	return s
}

func (s *SAML2IntegrationUnsetRequest) WithSaml2ForceAuthn(Saml2ForceAuthn *bool) *SAML2IntegrationUnsetRequest {
	s.Saml2ForceAuthn = Saml2ForceAuthn
	return s
}

func (s *SAML2IntegrationUnsetRequest) WithSaml2RequestedNameidFormat(Saml2RequestedNameidFormat *bool) *SAML2IntegrationUnsetRequest {
	s.Saml2RequestedNameidFormat = Saml2RequestedNameidFormat
	return s
}

func (s *SAML2IntegrationUnsetRequest) WithSaml2PostLogoutRedirectUrl(Saml2PostLogoutRedirectUrl *bool) *SAML2IntegrationUnsetRequest {
	s.Saml2PostLogoutRedirectUrl = Saml2PostLogoutRedirectUrl
	return s
}

func (s *SAML2IntegrationUnsetRequest) WithComment(Comment *bool) *SAML2IntegrationUnsetRequest {
	s.Comment = Comment
	return s
}

func NewAlterSCIMIntegrationSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterSCIMIntegrationSecurityIntegrationRequest {
	s := AlterSCIMIntegrationSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterSCIMIntegrationSecurityIntegrationRequest) WithIfExists(IfExists *bool) *AlterSCIMIntegrationSecurityIntegrationRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterSCIMIntegrationSecurityIntegrationRequest) WithSetTags(SetTags []TagAssociation) *AlterSCIMIntegrationSecurityIntegrationRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterSCIMIntegrationSecurityIntegrationRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterSCIMIntegrationSecurityIntegrationRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterSCIMIntegrationSecurityIntegrationRequest) WithSet(Set *SCIMIntegrationSetRequest) *AlterSCIMIntegrationSecurityIntegrationRequest {
	s.Set = Set
	return s
}

func (s *AlterSCIMIntegrationSecurityIntegrationRequest) WithUnset(Unset *SCIMIntegrationUnsetRequest) *AlterSCIMIntegrationSecurityIntegrationRequest {
	s.Unset = Unset
	return s
}

func NewSCIMIntegrationSetRequest() *SCIMIntegrationSetRequest {
	return &SCIMIntegrationSetRequest{}
}

func (s *SCIMIntegrationSetRequest) WithEnabled(Enabled *bool) *SCIMIntegrationSetRequest {
	s.Enabled = Enabled
	return s
}

func (s *SCIMIntegrationSetRequest) WithNetworkPolicy(NetworkPolicy *AccountObjectIdentifier) *SCIMIntegrationSetRequest {
	s.NetworkPolicy = NetworkPolicy
	return s
}

func (s *SCIMIntegrationSetRequest) WithSyncPassword(SyncPassword *bool) *SCIMIntegrationSetRequest {
	s.SyncPassword = SyncPassword
	return s
}

func (s *SCIMIntegrationSetRequest) WithComment(Comment *string) *SCIMIntegrationSetRequest {
	s.Comment = Comment
	return s
}

func NewSCIMIntegrationUnsetRequest() *SCIMIntegrationUnsetRequest {
	return &SCIMIntegrationUnsetRequest{}
}

func (s *SCIMIntegrationUnsetRequest) WithEnabled(Enabled *bool) *SCIMIntegrationUnsetRequest {
	s.Enabled = Enabled
	return s
}

func (s *SCIMIntegrationUnsetRequest) WithNetworkPolicy(NetworkPolicy *bool) *SCIMIntegrationUnsetRequest {
	s.NetworkPolicy = NetworkPolicy
	return s
}

func (s *SCIMIntegrationUnsetRequest) WithSyncPassword(SyncPassword *bool) *SCIMIntegrationUnsetRequest {
	s.SyncPassword = SyncPassword
	return s
}

func (s *SCIMIntegrationUnsetRequest) WithComment(Comment *bool) *SCIMIntegrationUnsetRequest {
	s.Comment = Comment
	return s
}

func NewDropSecurityIntegrationRequest(
	name AccountObjectIdentifier,
) *DropSecurityIntegrationRequest {
	s := DropSecurityIntegrationRequest{}
	s.name = name
	return &s
}

func (s *DropSecurityIntegrationRequest) WithIfExists(IfExists *bool) *DropSecurityIntegrationRequest {
	s.IfExists = IfExists
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

func (s *ShowSecurityIntegrationRequest) WithLike(Like *Like) *ShowSecurityIntegrationRequest {
	s.Like = Like
	return s
}
