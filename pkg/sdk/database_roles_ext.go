package sdk

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *databaseRoles) ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*DatabaseRole, error) {
	request := NewShowDatabaseRoleRequest().WithDatabase(id.DatabaseId()).WithLike(Like{Pointer(id.Name())})
	databaseRoles, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}

	result, err := collections.FindFirst(databaseRoles, func(r DatabaseRole) bool { return r.Name == id.Name() })
	if err != nil {
		return nil, err
	}
	result.DatabaseName = id.DatabaseName()
	return result, nil
}

func (v *databaseRoles) ShowByIDSafely(ctx context.Context, id DatabaseObjectIdentifier) (*DatabaseRole, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *databaseRoles) RevokeSafely(ctx context.Context, request *RevokeDatabaseRoleRequest) error {
	return SafeRevokePrivileges(func() error { return v.Revoke(ctx, request) })
}

func (v *databaseRoles) RevokeFromShareSafely(ctx context.Context, request *RevokeFromShareDatabaseRoleRequest) error {
	return SafeRevokePrivileges(func() error { return v.RevokeFromShare(ctx, request) })
}

// WithAccountRole is a convenience method for granting a database role to an account role.
func (s *GrantDatabaseRoleRequest) WithAccountRole(accountRole AccountObjectIdentifier) *GrantDatabaseRoleRequest {
	return s.WithTo(*NewDatabaseRoleKindOfRoleRequest().WithAccountRoleName(accountRole))
}

// WithDatabaseRole is a convenience method for granting a database role to another database role.
func (s *GrantDatabaseRoleRequest) WithDatabaseRole(databaseRole DatabaseObjectIdentifier) *GrantDatabaseRoleRequest {
	return s.WithTo(*NewDatabaseRoleKindOfRoleRequest().WithDatabaseRoleName(databaseRole))
}

// WithAccountRole is a convenience method for revoking a database role from an account role.
func (s *RevokeDatabaseRoleRequest) WithAccountRole(accountRole AccountObjectIdentifier) *RevokeDatabaseRoleRequest {
	return s.WithFrom(*NewDatabaseRoleKindOfRoleRequest().WithAccountRoleName(accountRole))
}

// WithDatabaseRole is a convenience method for revoking a database role from another database role.
func (s *RevokeDatabaseRoleRequest) WithDatabaseRole(databaseRole DatabaseObjectIdentifier) *RevokeDatabaseRoleRequest {
	return s.WithFrom(*NewDatabaseRoleKindOfRoleRequest().WithDatabaseRoleName(databaseRole))
}

func (opts *AlterDatabaseRoleOptions) additionalValidations() error {
	if opts.RenameTo != nil {
		if opts.name.DatabaseName() != opts.RenameTo.DatabaseName() {
			return errors.Join(ErrDifferentDatabase)
		}
	}
	return nil
}

func (opts *ShowDatabaseRoleOptions) additionalValidations() error {
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return errors.Join(ErrPatternRequiredForLikeKeyword)
	}
	return nil
}

func (r databaseRoleDBRow) additionalConvert(_ *DatabaseRole) error {
	// additionalConvert is generated as DatabaseName is a plain only field.
	// it can't be currently set here, as it is not a returned value, and we can get it only from ID(), which is not passed to convert method
	return nil
}
