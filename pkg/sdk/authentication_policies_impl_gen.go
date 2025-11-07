package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ AuthenticationPolicies = (*authenticationPolicies)(nil)

var (
	_ convertibleRow[AuthenticationPolicy]            = new(showAuthenticationPolicyDBRow)
	_ convertibleRow[AuthenticationPolicyDescription] = new(describeAuthenticationPolicyDBRow)
)

type authenticationPolicies struct {
	client *Client
}

func (v *authenticationPolicies) Create(ctx context.Context, request *CreateAuthenticationPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *authenticationPolicies) Alter(ctx context.Context, request *AlterAuthenticationPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *authenticationPolicies) Drop(ctx context.Context, request *DropAuthenticationPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *authenticationPolicies) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropAuthenticationPolicyRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *authenticationPolicies) Show(ctx context.Context, request *ShowAuthenticationPolicyRequest) ([]AuthenticationPolicy, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showAuthenticationPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[showAuthenticationPolicyDBRow, AuthenticationPolicy](dbRows)
}

func (v *authenticationPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicy, error) {
	request := NewShowAuthenticationPolicyRequest().
		WithIn(ExtendedIn{In: In{Schema: id.SchemaId()}}).
		WithLike(Like{Pattern: String(id.Name())})
	authenticationPolicies, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(authenticationPolicies, func(r AuthenticationPolicy) bool { return r.Name == id.Name() })
}

func (v *authenticationPolicies) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicy, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *authenticationPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]AuthenticationPolicyDescription, error) {
	opts := &DescribeAuthenticationPolicyOptions{
		name: id,
	}
	rows, err := validateAndQuery[describeAuthenticationPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[describeAuthenticationPolicyDBRow, AuthenticationPolicyDescription](rows)
}

func (r *CreateAuthenticationPolicyRequest) toOpts() *CreateAuthenticationPolicyOptions {
	opts := &CreateAuthenticationPolicyOptions{
		OrReplace:                r.OrReplace,
		IfNotExists:              r.IfNotExists,
		name:                     r.name,
		AuthenticationMethods:    r.AuthenticationMethods,
		MfaAuthenticationMethods: r.MfaAuthenticationMethods,
		MfaEnrollment:            r.MfaEnrollment,

		ClientTypes: r.ClientTypes,

		Comment: r.Comment,
	}
	if r.MfaPolicy != nil {
		opts.MfaPolicy = &AuthenticationPolicyMfaPolicy{
			EnforceMfaOnExternalAuthentication: r.MfaPolicy.EnforceMfaOnExternalAuthentication,
			AllowedMethods:                     r.MfaPolicy.AllowedMethods,
		}
	}
	if r.SecurityIntegrations != nil {
		opts.SecurityIntegrations = &SecurityIntegrationsOption{
			All:                  r.SecurityIntegrations.All,
			SecurityIntegrations: r.SecurityIntegrations.SecurityIntegrations,
		}
	}
	if r.PatPolicy != nil {
		opts.PatPolicy = &AuthenticationPolicyPatPolicy{
			DefaultExpiryInDays:     r.PatPolicy.DefaultExpiryInDays,
			MaxExpiryInDays:         r.PatPolicy.MaxExpiryInDays,
			NetworkPolicyEvaluation: r.PatPolicy.NetworkPolicyEvaluation,
		}
	}
	if r.WorkloadIdentityPolicy != nil {
		opts.WorkloadIdentityPolicy = &AuthenticationPolicyWorkloadIdentityPolicy{
			AllowedProviders:    r.WorkloadIdentityPolicy.AllowedProviders,
			AllowedAwsAccounts:  r.WorkloadIdentityPolicy.AllowedAwsAccounts,
			AllowedAzureIssuers: r.WorkloadIdentityPolicy.AllowedAzureIssuers,
			AllowedOidcIssuers:  r.WorkloadIdentityPolicy.AllowedOidcIssuers,
		}
	}
	return opts
}

func (r *AlterAuthenticationPolicyRequest) toOpts() *AlterAuthenticationPolicyOptions {
	opts := &AlterAuthenticationPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,

		RenameTo: r.RenameTo,
	}
	if r.Set != nil {
		opts.Set = &AuthenticationPolicySet{
			AuthenticationMethods:    r.Set.AuthenticationMethods,
			MfaAuthenticationMethods: r.Set.MfaAuthenticationMethods,
			MfaEnrollment:            r.Set.MfaEnrollment,

			ClientTypes: r.Set.ClientTypes,

			Comment: r.Set.Comment,
		}
		if r.Set.MfaPolicy != nil {
			opts.Set.MfaPolicy = &AuthenticationPolicyMfaPolicy{
				EnforceMfaOnExternalAuthentication: r.Set.MfaPolicy.EnforceMfaOnExternalAuthentication,
				AllowedMethods:                     r.Set.MfaPolicy.AllowedMethods,
			}
		}
		if r.Set.SecurityIntegrations != nil {
			opts.Set.SecurityIntegrations = &SecurityIntegrationsOption{
				All:                  r.Set.SecurityIntegrations.All,
				SecurityIntegrations: r.Set.SecurityIntegrations.SecurityIntegrations,
			}
		}
		if r.Set.PatPolicy != nil {
			opts.Set.PatPolicy = &AuthenticationPolicyPatPolicy{
				DefaultExpiryInDays:     r.Set.PatPolicy.DefaultExpiryInDays,
				MaxExpiryInDays:         r.Set.PatPolicy.MaxExpiryInDays,
				NetworkPolicyEvaluation: r.Set.PatPolicy.NetworkPolicyEvaluation,
			}
		}
		if r.Set.WorkloadIdentityPolicy != nil {
			opts.Set.WorkloadIdentityPolicy = &AuthenticationPolicyWorkloadIdentityPolicy{
				AllowedProviders:    r.Set.WorkloadIdentityPolicy.AllowedProviders,
				AllowedAwsAccounts:  r.Set.WorkloadIdentityPolicy.AllowedAwsAccounts,
				AllowedAzureIssuers: r.Set.WorkloadIdentityPolicy.AllowedAzureIssuers,
				AllowedOidcIssuers:  r.Set.WorkloadIdentityPolicy.AllowedOidcIssuers,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &AuthenticationPolicyUnset{
			ClientTypes:              r.Unset.ClientTypes,
			AuthenticationMethods:    r.Unset.AuthenticationMethods,
			SecurityIntegrations:     r.Unset.SecurityIntegrations,
			MfaAuthenticationMethods: r.Unset.MfaAuthenticationMethods,
			MfaEnrollment:            r.Unset.MfaEnrollment,
			MfaPolicy:                r.Unset.MfaPolicy,
			PatPolicy:                r.Unset.PatPolicy,
			WorkloadIdentityPolicy:   r.Unset.WorkloadIdentityPolicy,
			Comment:                  r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropAuthenticationPolicyRequest) toOpts() *DropAuthenticationPolicyOptions {
	opts := &DropAuthenticationPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowAuthenticationPolicyRequest) toOpts() *ShowAuthenticationPolicyOptions {
	opts := &ShowAuthenticationPolicyOptions{
		Like:       r.Like,
		In:         r.In,
		On:         r.On,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r showAuthenticationPolicyDBRow) convert() (*AuthenticationPolicy, error) {
	return &AuthenticationPolicy{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Kind:          r.Kind,
		Owner:         r.Owner,
		OwnerRoleType: r.OwnerRoleType,
		Options:       r.Options,
		Comment:       r.Comment,
	}, nil
}

func (r *DescribeAuthenticationPolicyRequest) toOpts() *DescribeAuthenticationPolicyOptions {
	opts := &DescribeAuthenticationPolicyOptions{
		name: r.name,
	}
	return opts
}

func (r describeAuthenticationPolicyDBRow) convert() (*AuthenticationPolicyDescription, error) {
	return &AuthenticationPolicyDescription{
		Property:    r.Property,
		Value:       r.Value,
		Default:     r.Default,
		Description: r.Description,
	}, nil
}
