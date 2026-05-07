package sdk

import "context"

func (v *applicationRoles) RevokeSafely(ctx context.Context, request *RevokeApplicationRoleRequest) error {
	return SafeRevokePrivileges(func() error { return v.Revoke(ctx, request) })
}
