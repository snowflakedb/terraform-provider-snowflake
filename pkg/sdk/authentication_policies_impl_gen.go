package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var _ AuthenticationPolicies = (*authenticationPolicies)(nil)

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
	resultList := convertRows[showAuthenticationPolicyDBRow, AuthenticationPolicy](dbRows)
	return resultList, nil
}

func (v *authenticationPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicy, error) {
	request := NewShowAuthenticationPolicyRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithIn(In{Schema: id.SchemaId()})
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
	return convertRows[describeAuthenticationPolicyDBRow, AuthenticationPolicyDescription](rows), nil
}

func (r *CreateAuthenticationPolicyRequest) toOpts() *CreateAuthenticationPolicyOptions {
	opts := &CreateAuthenticationPolicyOptions{
		OrReplace:                r.OrReplace,
		IfNotExists:              r.IfNotExists,
		name:                     r.name,
		AuthenticationMethods:    r.AuthenticationMethods,
		MfaAuthenticationMethods: r.MfaAuthenticationMethods,
		MfaEnrollment:            r.MfaEnrollment,
		ClientTypes:              r.ClientTypes,
		SecurityIntegrations:     r.SecurityIntegrations,
		Comment:                  r.Comment,
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
			ClientTypes:              r.Set.ClientTypes,
			SecurityIntegrations:     r.Set.SecurityIntegrations,
			Comment:                  r.Set.Comment,
		}
	}

	if r.Unset != nil {
		opts.Unset = &AuthenticationPolicyUnset{
			ClientTypes:              r.Unset.ClientTypes,
			AuthenticationMethods:    r.Unset.AuthenticationMethods,
			SecurityIntegrations:     r.Unset.SecurityIntegrations,
			MfaAuthenticationMethods: r.Unset.MfaAuthenticationMethods,
			MfaEnrollment:            r.Unset.MfaEnrollment,
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
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r showAuthenticationPolicyDBRow) convert() *AuthenticationPolicy {
	return &AuthenticationPolicy{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Owner:         r.Owner,
		OwnerRoleType: r.OwnerRoleType,
		Options:       r.Options,
		Comment:       r.Comment,
	}
}

func (r *DescribeAuthenticationPolicyRequest) toOpts() *DescribeAuthenticationPolicyOptions {
	opts := &DescribeAuthenticationPolicyOptions{
		name: r.name,
	}
	return opts
}

func (r describeAuthenticationPolicyDBRow) convert() *AuthenticationPolicyDescription {
	return &AuthenticationPolicyDescription{
		Property:    r.Property,
		Value:       r.Value,
		Default:     r.Default,
		Description: r.Description,
	}
}
