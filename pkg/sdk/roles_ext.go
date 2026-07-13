package sdk

import "context"

// UseRoleRequest and UseSecondaryRolesRequest are custom (no Options counterpart).
type UseRoleRequest struct {
	id AccountObjectIdentifier // required
}

func NewUseRoleRequest(id AccountObjectIdentifier) *UseRoleRequest {
	return &UseRoleRequest{id: id}
}

type UseSecondaryRolesRequest struct {
	option SecondaryRoleOption // required
}

func NewUseSecondaryRolesRequest(option SecondaryRoleOption) *UseSecondaryRolesRequest {
	return &UseSecondaryRolesRequest{option: option}
}

func (s *CreateRoleRequest) GetName() AccountObjectIdentifier {
	return s.name
}

func (opts *ShowRoleOptions) additionalValidations() error {
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return ErrPatternRequiredForLikeKeyword
	}
	return nil
}

func (opts *GrantRoleOptions) additionalValidations() error {
	if opts.Grant.Role != nil && !ValidObjectIdentifier(opts.Grant.Role) {
		return errInvalidIdentifier("GrantRoleOptions.Grant", "Role")
	}
	if opts.Grant.User != nil && !ValidObjectIdentifier(opts.Grant.User) {
		return errInvalidIdentifier("GrantRoleOptions.Grant", "User")
	}
	return nil
}

func (v *roles) RevokeSafely(ctx context.Context, req *RevokeRoleRequest) error {
	return SafeRevokePrivileges(func() error { return v.Revoke(ctx, req) })
}

func (v *roles) Use(ctx context.Context, req *UseRoleRequest) error {
	return v.client.Sessions.UseRole(ctx, NewUseRoleSessionRequest(req.id))
}

func (v *roles) UseSecondary(ctx context.Context, req *UseSecondaryRolesRequest) error {
	return v.client.Sessions.UseSecondaryRoles(ctx, NewUseSecondaryRolesSessionRequest(req.option))
}
