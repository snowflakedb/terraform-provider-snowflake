package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ApiIntegrations = (*apiIntegrations)(nil)

type apiIntegrations struct {
	client *Client
}

func (v *apiIntegrations) Create(ctx context.Context, request *CreateApiIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *apiIntegrations) Alter(ctx context.Context, request *AlterApiIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *apiIntegrations) Drop(ctx context.Context, request *DropApiIntegrationRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *apiIntegrations) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropApiIntegrationRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *apiIntegrations) Show(ctx context.Context, request *ShowApiIntegrationRequest) ([]ApiIntegration, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showApiIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showApiIntegrationsDbRow, ApiIntegration](dbRows)
	return resultList, nil
}

func (v *apiIntegrations) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegration, error) {
	request := NewShowApiIntegrationRequest().
		WithLike(Like{Pattern: String(id.Name())})
	apiIntegrations, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(apiIntegrations, func(r ApiIntegration) bool { return r.Name == id.Name() })
}

func (v *apiIntegrations) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ApiIntegration, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *apiIntegrations) Describe(ctx context.Context, id AccountObjectIdentifier) ([]ApiIntegrationProperty, error) {
	opts := &DescribeApiIntegrationOptions{
		name: id,
	}
	rows, err := validateAndQuery[descApiIntegrationsDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[descApiIntegrationsDbRow, ApiIntegrationProperty](rows), nil
}

func (r *CreateApiIntegrationRequest) toOpts() *CreateApiIntegrationOptions {
	opts := &CreateApiIntegrationOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		ApiAllowedPrefixes: r.ApiAllowedPrefixes,
		ApiBlockedPrefixes: r.ApiBlockedPrefixes,
		Enabled:            r.Enabled,
		Comment:            r.Comment,
	}
	if r.AwsApiProviderParams != nil {
		opts.AwsApiProviderParams = &AwsApiParams{
			ApiProvider:   r.AwsApiProviderParams.ApiProvider,
			ApiAwsRoleArn: r.AwsApiProviderParams.ApiAwsRoleArn,
			ApiKey:        r.AwsApiProviderParams.ApiKey,
		}
	}
	if r.AzureApiProviderParams != nil {
		opts.AzureApiProviderParams = &AzureApiParams{
			AzureTenantId:        r.AzureApiProviderParams.AzureTenantId,
			AzureAdApplicationId: r.AzureApiProviderParams.AzureAdApplicationId,
			ApiKey:               r.AzureApiProviderParams.ApiKey,
		}
	}
	if r.GoogleApiProviderParams != nil {
		opts.GoogleApiProviderParams = &GoogleApiParams{
			GoogleAudience: r.GoogleApiProviderParams.GoogleAudience,
		}
	}
	if r.GitApiProviderParams != nil {
		opts.GitApiProviderParams = &GitApiParams{
			AllowedAuthenticationSecrets: r.GitApiProviderParams.AllowedAuthenticationSecrets,
		}
	}
	return opts
}

func (r *AlterApiIntegrationRequest) toOpts() *AlterApiIntegrationOptions {
	opts := &AlterApiIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &ApiIntegrationSet{
			Enabled:            r.Set.Enabled,
			ApiAllowedPrefixes: r.Set.ApiAllowedPrefixes,
			ApiBlockedPrefixes: r.Set.ApiBlockedPrefixes,
			Comment:            r.Set.Comment,
		}
		if r.Set.AwsParams != nil {
			opts.Set.AwsParams = &SetAwsApiParams{
				ApiAwsRoleArn: r.Set.AwsParams.ApiAwsRoleArn,
				ApiKey:        r.Set.AwsParams.ApiKey,
			}
		}
		if r.Set.AzureParams != nil {
			opts.Set.AzureParams = &SetAzureApiParams{
				AzureTenantId:        r.Set.AzureParams.AzureTenantId,
				AzureAdApplicationId: r.Set.AzureParams.AzureAdApplicationId,
				ApiKey:               r.Set.AzureParams.ApiKey,
			}
		}
		if r.Set.GoogleParams != nil {
			opts.Set.GoogleParams = &SetGoogleApiParams{
				GoogleAudience: r.Set.GoogleParams.GoogleAudience,
			}
		}
		if r.Set.GitParams != nil {
			opts.Set.GitParams = &SetGitApiParams{
				AllowedAuthenticationSecrets: r.Set.GitParams.AllowedAuthenticationSecrets,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &ApiIntegrationUnset{
			ApiKey:             r.Unset.ApiKey,
			Enabled:            r.Unset.Enabled,
			ApiBlockedPrefixes: r.Unset.ApiBlockedPrefixes,
			Comment:            r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropApiIntegrationRequest) toOpts() *DropApiIntegrationOptions {
	opts := &DropApiIntegrationOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowApiIntegrationRequest) toOpts() *ShowApiIntegrationOptions {
	opts := &ShowApiIntegrationOptions{
		Like: r.Like,
	}
	return opts
}

func (r showApiIntegrationsDbRow) convert() *ApiIntegration {
	s := &ApiIntegration{
		Name:      r.Name,
		ApiType:   r.Type,
		Category:  r.Category,
		Enabled:   r.Enabled,
		CreatedOn: r.CreatedOn,
	}
	if r.Comment.Valid {
		s.Comment = r.Comment.String
	}
	return s
}

func (r *DescribeApiIntegrationRequest) toOpts() *DescribeApiIntegrationOptions {
	opts := &DescribeApiIntegrationOptions{
		name: r.name,
	}
	return opts
}

func (r descApiIntegrationsDbRow) convert() *ApiIntegrationProperty {
	return &ApiIntegrationProperty{
		Name:    r.Property,
		Type:    r.PropertyType,
		Value:   r.PropertyValue,
		Default: r.PropertyDefault,
	}
}
