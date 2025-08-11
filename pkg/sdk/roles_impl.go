package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Roles                          = (*roles)(nil)
	_ convertibleRowDeprecated[Role] = (*roleDBRow)(nil)
)

type roles struct {
	client *Client
}

func (v *roles) Create(ctx context.Context, req *CreateRoleRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *roles) Alter(ctx context.Context, req *AlterRoleRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *roles) Drop(ctx context.Context, req *DropRoleRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *roles) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropRoleRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *roles) Show(ctx context.Context, req *ShowRoleRequest) ([]Role, error) {
	dbRows, err := validateAndQuery[roleDBRow](v.client, ctx, req.toOpts())
	if err != nil {
		return nil, err
	}
	resultList := convertRows[roleDBRow, Role](dbRows)
	return resultList, nil
}

func (v *roles) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Role, error) {
	roleList, err := v.client.Roles.Show(ctx, NewShowRoleRequest().WithLike(NewLikeRequest(id.Name())))
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(roleList, func(r Role) bool { return r.ID().name == id.Name() })
}

func (v *roles) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Role, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *roles) Grant(ctx context.Context, req *GrantRoleRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *roles) Revoke(ctx context.Context, req *RevokeRoleRequest) error {
	return validateAndExec(v.client, ctx, req.toOpts())
}

func (v *roles) Use(ctx context.Context, req *UseRoleRequest) error {
	return v.client.Sessions.UseRole(ctx, req.id)
}

func (v *roles) UseSecondary(ctx context.Context, req *UseSecondaryRolesRequest) error {
	return v.client.Sessions.UseSecondaryRoles(ctx, req.option)
}
