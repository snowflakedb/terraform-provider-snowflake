package sdk

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *databaseRoles) ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*DatabaseRole, error) {
	request := NewShowDatabaseRoleRequest(id.DatabaseId()).WithLike(Like{Pointer(id.Name())})
	databaseRoles, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(databaseRoles, func(r DatabaseRole) bool { return r.Name == id.Name() })
}

func (v *databaseRoles) ShowByIDSafely(ctx context.Context, id DatabaseObjectIdentifier) (*DatabaseRole, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *databaseRoles) RevokeSafely(ctx context.Context, request *RevokeDatabaseRoleRequest) error {
	return SafeRevokePrivileges(func() error { return v.Revoke(ctx, request) })
}

func (v *databaseRoles) RevokeFromShareSafely(ctx context.Context, request *RevokeDatabaseRoleFromShareRequest) error {
	return SafeRevokePrivileges(func() error { return v.RevokeFromShare(ctx, request) })
}

func (r databaseRoleDBRow) additionalConvert(_ *DatabaseRole) error {
	return nil
}

func (opts *alterDatabaseRoleOptions) additionalValidations() error {
	if opts.Rename != nil {
		if opts.name.DatabaseName() != opts.Rename.DatabaseName() {
			return errors.Join(ErrDifferentDatabase)
		}
	}
	return nil
}
