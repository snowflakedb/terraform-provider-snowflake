package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var _ SecurityIntegrations = (*securityIntegrations)(nil)

type securityIntegrations struct {
	client *Client
}

func (v *securityIntegrations) CreateApiAuthenticationWithClientCredentialsFlow(ctx context.Context, request *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateApiAuthenticationWithAuthorizationCodeGrantFlow(ctx context.Context, request *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateApiAuthenticationWithJwtBearerFlow(ctx context.Context, request *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateExternalOauth(ctx context.Context, request *CreateExternalOauthSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateOauthForPartnerApplications(ctx context.Context, request *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateOauthForCustomClients(ctx context.Context, request *CreateOauthForCustomClientsSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateSaml2(ctx context.Context, request *CreateSaml2SecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) CreateScim(ctx context.Context, request *CreateScimSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterApiAuthenticationWithClientCredentialsFlow(ctx context.Context, request *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterApiAuthenticationWithAuthorizationCodeGrantFlow(ctx context.Context, request *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterApiAuthenticationWithJwtBearerFlow(ctx context.Context, request *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterExternalOauth(ctx context.Context, request *AlterExternalOauthSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterOauthForPartnerApplications(ctx context.Context, request *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterOauthForCustomClients(ctx context.Context, request *AlterOauthForCustomClientsSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterSaml2(ctx context.Context, request *AlterSaml2SecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) AlterScim(ctx context.Context, request *AlterScimSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) Drop(ctx context.Context, request *DropSecurityIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *securityIntegrations) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropSecurityIntegrationRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *securityIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]SecurityIntegrationProperty, error) {
	opts := &DescribeSecurityIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[securityIntegrationDescRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[securityIntegrationDescRow, SecurityIntegrationProperty](rows), nil
}

func (v *securityIntegrations) Show(ctx context.Context, request *ShowSecurityIntegrationRequest) ([]SecurityIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[securityIntegrationShowRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[securityIntegrationShowRow, SecurityIntegration](dbRows)
	return resultList, nil
}

func (v *securityIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*SecurityIntegration, error) {
	request := NewShowSecurityIntegrationRequest().
		WithLike(Like{Pattern: String(id.Name())})
	securityIntegrations, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(securityIntegrations, func(r SecurityIntegration) bool { return r.Name == id.Name() })
}

func (v *securityIntegrations) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*SecurityIntegration, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) toOpts() *CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions {
	opts := &CreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions{
		OrReplace:                   r.OrReplace,
		IfNotExists:                 r.IfNotExists,
		name:                        r.name,
		Enabled:                     r.Enabled,
		OauthTokenEndpoint:          r.OauthTokenEndpoint,
		OauthClientAuthMethod:       r.OauthClientAuthMethod,
		OauthClientId:               r.OauthClientId,
		OauthClientSecret:           r.OauthClientSecret,
		OauthGrantClientCredentials: r.OauthGrantClientCredentials,
		OauthAccessTokenValidity:    r.OauthAccessTokenValidity,
		OauthRefreshTokenValidity:   r.OauthRefreshTokenValidity,
		OauthAllowedScopes:          r.OauthAllowedScopes,
		Comment:                     r.Comment,
	}
	return opts
}

func (r *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest) toOpts() *CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions {
	opts := &CreateApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions{
		OrReplace:                   r.OrReplace,
		IfNotExists:                 r.IfNotExists,
		name:                        r.name,
		Enabled:                     r.Enabled,
		OauthAuthorizationEndpoint:  r.OauthAuthorizationEndpoint,
		OauthTokenEndpoint:          r.OauthTokenEndpoint,
		OauthClientAuthMethod:       r.OauthClientAuthMethod,
		OauthClientId:               r.OauthClientId,
		OauthClientSecret:           r.OauthClientSecret,
		OauthGrantAuthorizationCode: r.OauthGrantAuthorizationCode,
		OauthAccessTokenValidity:    r.OauthAccessTokenValidity,
		OauthRefreshTokenValidity:   r.OauthRefreshTokenValidity,
		OauthAllowedScopes:          r.OauthAllowedScopes,
		Comment:                     r.Comment,
	}
	return opts
}

func (r *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest) toOpts() *CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions {
	opts := &CreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions{
		OrReplace:                  r.OrReplace,
		IfNotExists:                r.IfNotExists,
		name:                       r.name,
		Enabled:                    r.Enabled,
		OauthAssertionIssuer:       r.OauthAssertionIssuer,
		OauthAuthorizationEndpoint: r.OauthAuthorizationEndpoint,
		OauthTokenEndpoint:         r.OauthTokenEndpoint,
		OauthClientAuthMethod:      r.OauthClientAuthMethod,
		OauthClientId:              r.OauthClientId,
		OauthClientSecret:          r.OauthClientSecret,
		OauthGrantJwtBearer:        r.OauthGrantJwtBearer,
		OauthAccessTokenValidity:   r.OauthAccessTokenValidity,
		OauthRefreshTokenValidity:  r.OauthRefreshTokenValidity,
		Comment:                    r.Comment,
	}
	return opts
}

func (r *CreateExternalOauthSecurityIntegrationRequest) toOpts() *CreateExternalOauthSecurityIntegrationOptions {
	opts := &CreateExternalOauthSecurityIntegrationOptions{
		OrReplace:                          r.OrReplace,
		IfNotExists:                        r.IfNotExists,
		name:                               r.name,
		Enabled:                            r.Enabled,
		ExternalOauthType:                  r.ExternalOauthType,
		ExternalOauthIssuer:                r.ExternalOauthIssuer,
		ExternalOauthTokenUserMappingClaim: r.ExternalOauthTokenUserMappingClaim,
		ExternalOauthSnowflakeUserMappingAttribute: r.ExternalOauthSnowflakeUserMappingAttribute,
		ExternalOauthJwsKeysUrl:                    r.ExternalOauthJwsKeysUrl,

		ExternalOauthRsaPublicKey:  r.ExternalOauthRsaPublicKey,
		ExternalOauthRsaPublicKey2: r.ExternalOauthRsaPublicKey2,

		ExternalOauthAnyRoleMode:           r.ExternalOauthAnyRoleMode,
		ExternalOauthScopeDelimiter:        r.ExternalOauthScopeDelimiter,
		ExternalOauthScopeMappingAttribute: r.ExternalOauthScopeMappingAttribute,
		Comment:                            r.Comment,
	}

	if r.ExternalOauthBlockedRolesList != nil {
		opts.ExternalOauthBlockedRolesList = &BlockedRolesList{
			BlockedRolesList: r.ExternalOauthBlockedRolesList.BlockedRolesList,
		}
	}

	if r.ExternalOauthAllowedRolesList != nil {
		opts.ExternalOauthAllowedRolesList = &AllowedRolesList{
			AllowedRolesList: r.ExternalOauthAllowedRolesList.AllowedRolesList,
		}
	}

	if r.ExternalOauthAudienceList != nil {
		opts.ExternalOauthAudienceList = &AudienceList{
			AudienceList: r.ExternalOauthAudienceList.AudienceList,
		}
	}

	return opts
}

func (r *CreateOauthForPartnerApplicationsSecurityIntegrationRequest) toOpts() *CreateOauthForPartnerApplicationsSecurityIntegrationOptions {
	opts := &CreateOauthForPartnerApplicationsSecurityIntegrationOptions{
		OrReplace:                 r.OrReplace,
		IfNotExists:               r.IfNotExists,
		name:                      r.name,
		OauthClient:               r.OauthClient,
		OauthRedirectUri:          r.OauthRedirectUri,
		Enabled:                   r.Enabled,
		OauthIssueRefreshTokens:   r.OauthIssueRefreshTokens,
		OauthRefreshTokenValidity: r.OauthRefreshTokenValidity,
		OauthUseSecondaryRoles:    r.OauthUseSecondaryRoles,

		Comment: r.Comment,
	}

	if r.BlockedRolesList != nil {
		opts.BlockedRolesList = &BlockedRolesList{
			BlockedRolesList: r.BlockedRolesList.BlockedRolesList,
		}
	}

	return opts
}

func (r *CreateOauthForCustomClientsSecurityIntegrationRequest) toOpts() *CreateOauthForCustomClientsSecurityIntegrationOptions {
	opts := &CreateOauthForCustomClientsSecurityIntegrationOptions{
		OrReplace:                   r.OrReplace,
		IfNotExists:                 r.IfNotExists,
		name:                        r.name,
		OauthClientType:             r.OauthClientType,
		OauthRedirectUri:            r.OauthRedirectUri,
		Enabled:                     r.Enabled,
		OauthAllowNonTlsRedirectUri: r.OauthAllowNonTlsRedirectUri,
		OauthEnforcePkce:            r.OauthEnforcePkce,
		OauthUseSecondaryRoles:      r.OauthUseSecondaryRoles,

		OauthIssueRefreshTokens:   r.OauthIssueRefreshTokens,
		OauthRefreshTokenValidity: r.OauthRefreshTokenValidity,
		NetworkPolicy:             r.NetworkPolicy,
		OauthClientRsaPublicKey:   r.OauthClientRsaPublicKey,
		OauthClientRsaPublicKey2:  r.OauthClientRsaPublicKey2,
		Comment:                   r.Comment,
	}

	if r.PreAuthorizedRolesList != nil {
		opts.PreAuthorizedRolesList = &PreAuthorizedRolesList{
			PreAuthorizedRolesList: r.PreAuthorizedRolesList.PreAuthorizedRolesList,
		}
	}

	if r.BlockedRolesList != nil {
		opts.BlockedRolesList = &BlockedRolesList{
			BlockedRolesList: r.BlockedRolesList.BlockedRolesList,
		}
	}

	return opts
}

func (r *CreateSaml2SecurityIntegrationRequest) toOpts() *CreateSaml2SecurityIntegrationOptions {
	opts := &CreateSaml2SecurityIntegrationOptions{
		OrReplace:                      r.OrReplace,
		IfNotExists:                    r.IfNotExists,
		name:                           r.name,
		Enabled:                        r.Enabled,
		Saml2Issuer:                    r.Saml2Issuer,
		Saml2SsoUrl:                    r.Saml2SsoUrl,
		Saml2Provider:                  r.Saml2Provider,
		Saml2X509Cert:                  r.Saml2X509Cert,
		AllowedUserDomains:             r.AllowedUserDomains,
		AllowedEmailPatterns:           r.AllowedEmailPatterns,
		Saml2SpInitiatedLoginPageLabel: r.Saml2SpInitiatedLoginPageLabel,
		Saml2EnableSpInitiated:         r.Saml2EnableSpInitiated,
		Saml2SnowflakeX509Cert:         r.Saml2SnowflakeX509Cert,
		Saml2SignRequest:               r.Saml2SignRequest,
		Saml2RequestedNameidFormat:     r.Saml2RequestedNameidFormat,
		Saml2PostLogoutRedirectUrl:     r.Saml2PostLogoutRedirectUrl,
		Saml2ForceAuthn:                r.Saml2ForceAuthn,
		Saml2SnowflakeIssuerUrl:        r.Saml2SnowflakeIssuerUrl,
		Saml2SnowflakeAcsUrl:           r.Saml2SnowflakeAcsUrl,
		Comment:                        r.Comment,
	}
	return opts
}

func (r *CreateScimSecurityIntegrationRequest) toOpts() *CreateScimSecurityIntegrationOptions {
	opts := &CreateScimSecurityIntegrationOptions{
		OrReplace:     r.OrReplace,
		IfNotExists:   r.IfNotExists,
		name:          r.name,
		Enabled:       r.Enabled,
		ScimClient:    r.ScimClient,
		RunAsRole:     r.RunAsRole,
		NetworkPolicy: r.NetworkPolicy,
		SyncPassword:  r.SyncPassword,
		Comment:       r.Comment,
	}
	return opts
}

func (r *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest) toOpts() *AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions {
	opts := &AlterApiAuthenticationWithClientCredentialsFlowSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}

	if r.Set != nil {
		opts.Set = &ApiAuthenticationWithClientCredentialsFlowIntegrationSet{
			Enabled:                     r.Set.Enabled,
			OauthTokenEndpoint:          r.Set.OauthTokenEndpoint,
			OauthClientAuthMethod:       r.Set.OauthClientAuthMethod,
			OauthClientId:               r.Set.OauthClientId,
			OauthClientSecret:           r.Set.OauthClientSecret,
			OauthGrantClientCredentials: r.Set.OauthGrantClientCredentials,
			OauthAccessTokenValidity:    r.Set.OauthAccessTokenValidity,
			OauthRefreshTokenValidity:   r.Set.OauthRefreshTokenValidity,
			OauthAllowedScopes:          r.Set.OauthAllowedScopes,
			Comment:                     r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &ApiAuthenticationWithClientCredentialsFlowIntegrationUnset{
			Enabled: r.Unset.Enabled,
			Comment: r.Unset.Comment,
		}
	}

	return opts
}

func (r *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationRequest) toOpts() *AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions {
	opts := &AlterApiAuthenticationWithAuthorizationCodeGrantFlowSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}

	if r.Set != nil {
		opts.Set = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationSet{
			Enabled:                     r.Set.Enabled,
			OauthAuthorizationEndpoint:  r.Set.OauthAuthorizationEndpoint,
			OauthTokenEndpoint:          r.Set.OauthTokenEndpoint,
			OauthClientAuthMethod:       r.Set.OauthClientAuthMethod,
			OauthClientId:               r.Set.OauthClientId,
			OauthClientSecret:           r.Set.OauthClientSecret,
			OauthGrantAuthorizationCode: r.Set.OauthGrantAuthorizationCode,
			OauthAccessTokenValidity:    r.Set.OauthAccessTokenValidity,
			OauthRefreshTokenValidity:   r.Set.OauthRefreshTokenValidity,
			OauthAllowedScopes:          r.Set.OauthAllowedScopes,
			Comment:                     r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &ApiAuthenticationWithAuthorizationCodeGrantFlowIntegrationUnset{
			Enabled: r.Unset.Enabled,
			Comment: r.Unset.Comment,
		}
	}

	return opts
}

func (r *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest) toOpts() *AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions {
	opts := &AlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}

	if r.Set != nil {
		opts.Set = &ApiAuthenticationWithJwtBearerFlowIntegrationSet{
			Enabled:                    r.Set.Enabled,
			OauthAuthorizationEndpoint: r.Set.OauthAuthorizationEndpoint,
			OauthTokenEndpoint:         r.Set.OauthTokenEndpoint,
			OauthClientAuthMethod:      r.Set.OauthClientAuthMethod,
			OauthClientId:              r.Set.OauthClientId,
			OauthClientSecret:          r.Set.OauthClientSecret,
			OauthGrantJwtBearer:        r.Set.OauthGrantJwtBearer,
			OauthAccessTokenValidity:   r.Set.OauthAccessTokenValidity,
			OauthRefreshTokenValidity:  r.Set.OauthRefreshTokenValidity,
			Comment:                    r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &ApiAuthenticationWithJwtBearerFlowIntegrationUnset{
			Enabled: r.Unset.Enabled,
			Comment: r.Unset.Comment,
		}
	}

	return opts
}

func (r *AlterExternalOauthSecurityIntegrationRequest) toOpts() *AlterExternalOauthSecurityIntegrationOptions {
	opts := &AlterExternalOauthSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}

	if r.Set != nil {
		opts.Set = &ExternalOauthIntegrationSet{
			Enabled:                            r.Set.Enabled,
			ExternalOauthType:                  r.Set.ExternalOauthType,
			ExternalOauthIssuer:                r.Set.ExternalOauthIssuer,
			ExternalOauthTokenUserMappingClaim: r.Set.ExternalOauthTokenUserMappingClaim,
			ExternalOauthSnowflakeUserMappingAttribute: r.Set.ExternalOauthSnowflakeUserMappingAttribute,
			ExternalOauthJwsKeysUrl:                    r.Set.ExternalOauthJwsKeysUrl,

			ExternalOauthRsaPublicKey:  r.Set.ExternalOauthRsaPublicKey,
			ExternalOauthRsaPublicKey2: r.Set.ExternalOauthRsaPublicKey2,

			ExternalOauthAnyRoleMode:           r.Set.ExternalOauthAnyRoleMode,
			ExternalOauthScopeDelimiter:        r.Set.ExternalOauthScopeDelimiter,
			ExternalOauthScopeMappingAttribute: r.Set.ExternalOauthScopeMappingAttribute,
			Comment:                            r.Set.Comment,
		}

		if r.Set.ExternalOauthBlockedRolesList != nil {
			opts.Set.ExternalOauthBlockedRolesList = &BlockedRolesList{
				BlockedRolesList: r.Set.ExternalOauthBlockedRolesList.BlockedRolesList,
			}
		}

		if r.Set.ExternalOauthAllowedRolesList != nil {
			opts.Set.ExternalOauthAllowedRolesList = &AllowedRolesList{
				AllowedRolesList: r.Set.ExternalOauthAllowedRolesList.AllowedRolesList,
			}
		}

		if r.Set.ExternalOauthAudienceList != nil {
			opts.Set.ExternalOauthAudienceList = &AudienceList{
				AudienceList: r.Set.ExternalOauthAudienceList.AudienceList,
			}
		}
	}

	if r.Unset != nil {
		opts.Unset = &ExternalOauthIntegrationUnset{
			Enabled:                   r.Unset.Enabled,
			ExternalOauthAudienceList: r.Unset.ExternalOauthAudienceList,
		}
	}

	return opts
}

func (r *AlterOauthForPartnerApplicationsSecurityIntegrationRequest) toOpts() *AlterOauthForPartnerApplicationsSecurityIntegrationOptions {
	opts := &AlterOauthForPartnerApplicationsSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}

	if r.Set != nil {
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{
			Enabled:                   r.Set.Enabled,
			OauthIssueRefreshTokens:   r.Set.OauthIssueRefreshTokens,
			OauthRedirectUri:          r.Set.OauthRedirectUri,
			OauthRefreshTokenValidity: r.Set.OauthRefreshTokenValidity,
			OauthUseSecondaryRoles:    r.Set.OauthUseSecondaryRoles,
		}
		if r.Set.Comment != nil {
			opts.Set.Comment = &StringAllowEmpty{*r.Set.Comment}
		}

		if r.Set.BlockedRolesList != nil {
			opts.Set.BlockedRolesList = &BlockedRolesList{
				BlockedRolesList: r.Set.BlockedRolesList.BlockedRolesList,
			}
		}
	}

	if r.Unset != nil {
		opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{
			Enabled:                r.Unset.Enabled,
			OauthUseSecondaryRoles: r.Unset.OauthUseSecondaryRoles,
		}
	}

	return opts
}

func (r *AlterOauthForCustomClientsSecurityIntegrationRequest) toOpts() *AlterOauthForCustomClientsSecurityIntegrationOptions {
	opts := &AlterOauthForCustomClientsSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}

	if r.Set != nil {
		opts.Set = &OauthForCustomClientsIntegrationSet{
			Enabled:                     r.Set.Enabled,
			OauthRedirectUri:            r.Set.OauthRedirectUri,
			OauthAllowNonTlsRedirectUri: r.Set.OauthAllowNonTlsRedirectUri,
			OauthEnforcePkce:            r.Set.OauthEnforcePkce,

			OauthIssueRefreshTokens:   r.Set.OauthIssueRefreshTokens,
			OauthRefreshTokenValidity: r.Set.OauthRefreshTokenValidity,
			OauthUseSecondaryRoles:    r.Set.OauthUseSecondaryRoles,
			NetworkPolicy:             r.Set.NetworkPolicy,
			OauthClientRsaPublicKey:   r.Set.OauthClientRsaPublicKey,
			OauthClientRsaPublicKey2:  r.Set.OauthClientRsaPublicKey2,
			Comment:                   r.Set.Comment,
		}

		if r.Set.PreAuthorizedRolesList != nil {
			opts.Set.PreAuthorizedRolesList = &PreAuthorizedRolesList{
				PreAuthorizedRolesList: r.Set.PreAuthorizedRolesList.PreAuthorizedRolesList,
			}
		}

		if r.Set.BlockedRolesList != nil {
			opts.Set.BlockedRolesList = &BlockedRolesList{
				BlockedRolesList: r.Set.BlockedRolesList.BlockedRolesList,
			}
		}
	}

	if r.Unset != nil {
		opts.Unset = &OauthForCustomClientsIntegrationUnset{
			Enabled:                  r.Unset.Enabled,
			NetworkPolicy:            r.Unset.NetworkPolicy,
			OauthClientRsaPublicKey:  r.Unset.OauthClientRsaPublicKey,
			OauthClientRsaPublicKey2: r.Unset.OauthClientRsaPublicKey2,
			OauthUseSecondaryRoles:   r.Unset.OauthUseSecondaryRoles,
		}
	}

	return opts
}

func (r *AlterSaml2SecurityIntegrationRequest) toOpts() *AlterSaml2SecurityIntegrationOptions {
	opts := &AlterSaml2SecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,

		RefreshSaml2SnowflakePrivateKey: r.RefreshSaml2SnowflakePrivateKey,
	}

	if r.Set != nil {
		opts.Set = &Saml2IntegrationSet{
			Enabled:                        r.Set.Enabled,
			Saml2Issuer:                    r.Set.Saml2Issuer,
			Saml2SsoUrl:                    r.Set.Saml2SsoUrl,
			Saml2Provider:                  r.Set.Saml2Provider,
			Saml2X509Cert:                  r.Set.Saml2X509Cert,
			AllowedUserDomains:             r.Set.AllowedUserDomains,
			AllowedEmailPatterns:           r.Set.AllowedEmailPatterns,
			Saml2SpInitiatedLoginPageLabel: r.Set.Saml2SpInitiatedLoginPageLabel,
			Saml2EnableSpInitiated:         r.Set.Saml2EnableSpInitiated,
			Saml2SnowflakeX509Cert:         r.Set.Saml2SnowflakeX509Cert,
			Saml2SignRequest:               r.Set.Saml2SignRequest,
			Saml2RequestedNameidFormat:     r.Set.Saml2RequestedNameidFormat,
			Saml2PostLogoutRedirectUrl:     r.Set.Saml2PostLogoutRedirectUrl,
			Saml2ForceAuthn:                r.Set.Saml2ForceAuthn,
			Saml2SnowflakeIssuerUrl:        r.Set.Saml2SnowflakeIssuerUrl,
			Saml2SnowflakeAcsUrl:           r.Set.Saml2SnowflakeAcsUrl,
			Comment:                        r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &Saml2IntegrationUnset{
			Saml2ForceAuthn:            r.Unset.Saml2ForceAuthn,
			Saml2RequestedNameidFormat: r.Unset.Saml2RequestedNameidFormat,
			Saml2PostLogoutRedirectUrl: r.Unset.Saml2PostLogoutRedirectUrl,
			Comment:                    r.Unset.Comment,
		}
	}

	return opts
}

func (r *AlterScimSecurityIntegrationRequest) toOpts() *AlterScimSecurityIntegrationOptions {
	opts := &AlterScimSecurityIntegrationOptions{
		IfExists:  r.IfExists,
		name:      r.name,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}

	if r.Set != nil {
		opts.Set = &ScimIntegrationSet{
			Enabled:       r.Set.Enabled,
			NetworkPolicy: r.Set.NetworkPolicy,
			SyncPassword:  r.Set.SyncPassword,
			Comment:       r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &ScimIntegrationUnset{
			Enabled:       r.Unset.Enabled,
			NetworkPolicy: r.Unset.NetworkPolicy,
			SyncPassword:  r.Unset.SyncPassword,
		}
	}

	return opts
}

func (r *DropSecurityIntegrationRequest) toOpts() *DropSecurityIntegrationOptions {
	opts := &DropSecurityIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeSecurityIntegrationRequest) toOpts() *DescribeSecurityIntegrationOptions {
	opts := &DescribeSecurityIntegrationOptions{
		name: r.name,
	}
	return opts
}

func (r securityIntegrationDescRow) convert() *SecurityIntegrationProperty {
	return &SecurityIntegrationProperty{
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
}

func (r *ShowSecurityIntegrationRequest) toOpts() *ShowSecurityIntegrationOptions {
	opts := &ShowSecurityIntegrationOptions{
		Like: r.Like,
	}
	return opts
}

func (r securityIntegrationShowRow) convert() *SecurityIntegration {
	s := &SecurityIntegration{
		Name:            r.Name,
		IntegrationType: r.Type,
		Enabled:         r.Enabled,
		CreatedOn:       r.CreatedOn,
		Category:        r.Category,
	}
	if r.Comment.Valid {
		s.Comment = r.Comment.String
	}
	return s
}
