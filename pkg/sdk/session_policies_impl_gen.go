package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ SessionPolicies = (*sessionPolicies)(nil)
var _ convertibleRow[SessionPolicy] = new(showSessionPolicyDBRow)
var _ convertibleRow[SessionPolicyDescription] = new(describeSessionPolicyDBRow)

type sessionPolicies struct {
	client *Client
}

func (v *sessionPolicies) Create(ctx context.Context, request *CreateSessionPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *sessionPolicies) Alter(ctx context.Context, request *AlterSessionPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *sessionPolicies) Drop(ctx context.Context, request *DropSessionPolicyRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *sessionPolicies) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropSessionPolicyRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *sessionPolicies) Show(ctx context.Context, request *ShowSessionPolicyRequest) ([]SessionPolicy, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showSessionPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[showSessionPolicyDBRow, SessionPolicy](dbRows)
}

func (v *sessionPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicy, error) {
	request := NewShowSessionPolicyRequest()
	sessionPolicies, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(sessionPolicies, func(r SessionPolicy) bool { return r.Name == id.Name() })
}

func (v *sessionPolicies) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicy, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *sessionPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicyDescription, error) {
	opts := &DescribeSessionPolicyOptions{
		name: id,
	}
	result, err := validateAndQueryOne[describeSessionPolicyDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return conversionErrorWrapped(result.convert())
}

func (r *CreateSessionPolicyRequest) toOpts() *CreateSessionPolicyOptions {
	opts := &CreateSessionPolicyOptions{
		OrReplace:                r.OrReplace,
		IfNotExists:              r.IfNotExists,
		name:                     r.name,
		SessionIdleTimeoutMins:   r.SessionIdleTimeoutMins,
		SessionUiIdleTimeoutMins: r.SessionUiIdleTimeoutMins,
		Comment:                  r.Comment,
	}
	return opts
}

func (r *AlterSessionPolicyRequest) toOpts() *AlterSessionPolicyOptions {
	opts := &AlterSessionPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,
		RenameTo: r.RenameTo,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &SessionPolicySet{
			SessionIdleTimeoutMins:   r.Set.SessionIdleTimeoutMins,
			SessionUiIdleTimeoutMins: r.Set.SessionUiIdleTimeoutMins,
			Comment:                  r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &SessionPolicyUnset{
			SessionIdleTimeoutMins:   r.Unset.SessionIdleTimeoutMins,
			SessionUiIdleTimeoutMins: r.Unset.SessionUiIdleTimeoutMins,
			Comment:                  r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropSessionPolicyRequest) toOpts() *DropSessionPolicyOptions {
	opts := &DropSessionPolicyOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowSessionPolicyRequest) toOpts() *ShowSessionPolicyOptions {
	opts := &ShowSessionPolicyOptions{}
	return opts
}

func (r showSessionPolicyDBRow) convert() (*SessionPolicy, error) {
	return &SessionPolicy{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Kind:          r.Kind,
		Owner:         r.Owner,
		Comment:       r.Comment,
		Options:       r.Options,
		OwnerRoleType: r.OwnerRoleType,
	}, nil
}

func (r *DescribeSessionPolicyRequest) toOpts() *DescribeSessionPolicyOptions {
	opts := &DescribeSessionPolicyOptions{
		name: r.name,
	}
	return opts
}

func (r describeSessionPolicyDBRow) convert() (*SessionPolicyDescription, error) {
	sessionPolicyDescription := SessionPolicyDescription{
		CreatedOn:                r.CreatedOn,
		Name:                     r.Name,
		SessionIdleTimeoutMins:   r.SessionIdleTimeoutMins,
		SessionUIIdleTimeoutMins: r.SessionUiIdleTimeoutMins,
	}
	if r.Comment.Valid {
		sessionPolicyDescription.Comment = r.Comment.String
	}
	return &sessionPolicyDescription, nil
}
