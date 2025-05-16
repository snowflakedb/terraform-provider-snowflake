package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var _ ApplicationRoles = (*applicationRoles)(nil)

type applicationRoles struct {
	client *Client
}

func (v *applicationRoles) Grant(ctx context.Context, request *GrantApplicationRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationRoles) Revoke(ctx context.Context, request *RevokeApplicationRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *applicationRoles) Show(ctx context.Context, request *ShowApplicationRoleRequest) ([]ApplicationRole, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[applicationRoleDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[applicationRoleDbRow, ApplicationRole](dbRows)
	return resultList, nil
}

func (v *applicationRoles) ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*ApplicationRole, error) {
	request := NewShowApplicationRoleRequest().
		WithApplicationName(id.DatabaseId())
	applicationRoles, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(applicationRoles, func(r ApplicationRole) bool { return r.Name == id.Name() })
}

func (v *applicationRoles) ShowByIDSafely(ctx context.Context, id DatabaseObjectIdentifier) (*ApplicationRole, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *GrantApplicationRoleRequest) toOpts() *GrantApplicationRoleOptions {
	opts := &GrantApplicationRoleOptions{
		name: r.name,
	}
	opts.To = KindOfRole{
		RoleName:            r.To.RoleName,
		ApplicationRoleName: r.To.ApplicationRoleName,
		ApplicationName:     r.To.ApplicationName,
	}
	return opts
}

func (r *RevokeApplicationRoleRequest) toOpts() *RevokeApplicationRoleOptions {
	opts := &RevokeApplicationRoleOptions{
		name: r.name,
	}
	opts.From = KindOfRole{
		RoleName:            r.From.RoleName,
		ApplicationRoleName: r.From.ApplicationRoleName,
		ApplicationName:     r.From.ApplicationName,
	}
	return opts
}

func (r *ShowApplicationRoleRequest) toOpts() *ShowApplicationRoleOptions {
	opts := &ShowApplicationRoleOptions{
		ApplicationName: r.ApplicationName,
		Limit:           r.Limit,
	}
	return opts
}

func (r applicationRoleDbRow) convert() *ApplicationRole {
	return &ApplicationRole{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		Owner:         r.Owner,
		Comment:       r.Comment,
		OwnerRoleType: r.OwnerRoleType,
	}
}
