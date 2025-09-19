package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Secrets                       = (*secrets)(nil)
	_ convertibleRow[Secret]        = new(secretDBRow)
	_ convertibleRow[SecretDetails] = new(secretDetailsDBRow)
)

type secrets struct {
	client *Client
}

func (v *secrets) CreateWithOAuthClientCredentialsFlow(ctx context.Context, request *CreateWithOAuthClientCredentialsFlowSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) CreateWithOAuthAuthorizationCodeFlow(ctx context.Context, request *CreateWithOAuthAuthorizationCodeFlowSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) CreateWithBasicAuthentication(ctx context.Context, request *CreateWithBasicAuthenticationSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) CreateWithGenericString(ctx context.Context, request *CreateWithGenericStringSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) Alter(ctx context.Context, request *AlterSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) Drop(ctx context.Context, request *DropSecretRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *secrets) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropSecretRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *secrets) Show(ctx context.Context, request *ShowSecretRequest) ([]Secret, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[secretDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[secretDBRow, Secret](dbRows)
}

func (v *secrets) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Secret, error) {
	request := NewShowSecretRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithIn(ExtendedIn{In: In{Schema: id.SchemaId()}})
	secrets, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(secrets, func(r Secret) bool { return r.Name == id.Name() })
}

func (v *secrets) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Secret, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *secrets) Describe(ctx context.Context, id SchemaObjectIdentifier) (*SecretDetails, error) {
	opts := &DescribeSecretOptions{
		name: id,
	}
	result, err := validateAndQueryOne[secretDetailsDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return conversionErrorWrapped(result.convert())
}

func (r *CreateWithOAuthClientCredentialsFlowSecretRequest) toOpts() *CreateWithOAuthClientCredentialsFlowSecretOptions {
	opts := &CreateWithOAuthClientCredentialsFlowSecretOptions{
		OrReplace:      r.OrReplace,
		IfNotExists:    r.IfNotExists,
		name:           r.name,
		ApiIntegration: r.ApiIntegration,
		Comment:        r.Comment,
	}

	if r.OauthScopes != nil {
		opts.OauthScopes = &OauthScopesList{
			OauthScopesList: r.OauthScopes.OauthScopesList,
		}
	}

	return opts
}

func (r *CreateWithOAuthAuthorizationCodeFlowSecretRequest) toOpts() *CreateWithOAuthAuthorizationCodeFlowSecretOptions {
	opts := &CreateWithOAuthAuthorizationCodeFlowSecretOptions{
		OrReplace:                   r.OrReplace,
		IfNotExists:                 r.IfNotExists,
		name:                        r.name,
		OauthRefreshToken:           r.OauthRefreshToken,
		OauthRefreshTokenExpiryTime: r.OauthRefreshTokenExpiryTime,
		ApiIntegration:              r.ApiIntegration,
		Comment:                     r.Comment,
	}
	return opts
}

func (r *CreateWithBasicAuthenticationSecretRequest) toOpts() *CreateWithBasicAuthenticationSecretOptions {
	opts := &CreateWithBasicAuthenticationSecretOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Username:    r.Username,
		Password:    r.Password,
		Comment:     r.Comment,
	}
	return opts
}

func (r *CreateWithGenericStringSecretRequest) toOpts() *CreateWithGenericStringSecretOptions {
	opts := &CreateWithGenericStringSecretOptions{
		OrReplace:    r.OrReplace,
		IfNotExists:  r.IfNotExists,
		name:         r.name,
		SecretString: r.SecretString,
		Comment:      r.Comment,
	}
	return opts
}

func (r *AlterSecretRequest) toOpts() *AlterSecretOptions {
	opts := &AlterSecretOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}

	if r.Set != nil {
		opts.Set = &SecretSet{
			Comment: r.Set.Comment,
		}

		if r.Set.SetForFlow != nil {
			opts.Set.SetForFlow = &SetForFlow{}

			if r.Set.SetForFlow.SetForOAuthClientCredentials != nil {
				opts.Set.SetForFlow.SetForOAuthClientCredentials = &SetForOAuthClientCredentials{}

				if r.Set.SetForFlow.SetForOAuthClientCredentials.OauthScopes != nil {
					opts.Set.SetForFlow.SetForOAuthClientCredentials.OauthScopes = &OauthScopesList{
						OauthScopesList: r.Set.SetForFlow.SetForOAuthClientCredentials.OauthScopes.OauthScopesList,
					}
				}
			}

			if r.Set.SetForFlow.SetForOAuthAuthorization != nil {
				opts.Set.SetForFlow.SetForOAuthAuthorization = &SetForOAuthAuthorization{
					OauthRefreshToken:           r.Set.SetForFlow.SetForOAuthAuthorization.OauthRefreshToken,
					OauthRefreshTokenExpiryTime: r.Set.SetForFlow.SetForOAuthAuthorization.OauthRefreshTokenExpiryTime,
				}
			}

			if r.Set.SetForFlow.SetForBasicAuthentication != nil {
				opts.Set.SetForFlow.SetForBasicAuthentication = &SetForBasicAuthentication{
					Username: r.Set.SetForFlow.SetForBasicAuthentication.Username,
					Password: r.Set.SetForFlow.SetForBasicAuthentication.Password,
				}
			}

			if r.Set.SetForFlow.SetForGenericString != nil {
				opts.Set.SetForFlow.SetForGenericString = &SetForGenericString{
					SecretString: r.Set.SetForFlow.SetForGenericString.SecretString,
				}
			}
		}
	}

	if r.Unset != nil {
		opts.Unset = &SecretUnset{
			Comment: r.Unset.Comment,
		}
	}

	return opts
}

func (r *DropSecretRequest) toOpts() *DropSecretOptions {
	opts := &DropSecretOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowSecretRequest) toOpts() *ShowSecretOptions {
	opts := &ShowSecretOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r secretDBRow) convert() (*Secret, error) {
	s := &Secret{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		SchemaName:    r.SchemaName,
		DatabaseName:  r.DatabaseName,
		Owner:         r.Owner,
		SecretType:    r.SecretType,
		OwnerRoleType: r.OwnerRoleType,
	}
	if r.Comment.Valid {
		s.Comment = String(r.Comment.String)
	}
	if r.OauthScopes.Valid {
		s.OauthScopes = ParseCommaSeparatedStringArray(r.OauthScopes.String, false)
	}
	return s, nil
}

func (r *DescribeSecretRequest) toOpts() *DescribeSecretOptions {
	opts := &DescribeSecretOptions{
		name: r.name,
	}
	return opts
}

func (r secretDetailsDBRow) convert() (*SecretDetails, error) {
	s := &SecretDetails{
		CreatedOn:                   r.CreatedOn,
		Name:                        r.Name,
		SchemaName:                  r.SchemaName,
		DatabaseName:                r.DatabaseName,
		Owner:                       r.Owner,
		SecretType:                  r.SecretType,
		OauthAccessTokenExpiryTime:  r.OauthAccessTokenExpiryTime,
		OauthRefreshTokenExpiryTime: r.OauthRefreshTokenExpiryTime,
	}
	if r.Username.Valid {
		s.Username = String(r.Username.String)
	}
	if r.Comment.Valid {
		s.Comment = String(r.Comment.String)
	}
	if r.OauthScopes.Valid {
		s.OauthScopes = ParseCommaSeparatedStringArray(r.OauthScopes.String, false)
	}
	if r.IntegrationName.Valid {
		s.IntegrationName = String(r.IntegrationName.String)
	}
	return s, nil
}
