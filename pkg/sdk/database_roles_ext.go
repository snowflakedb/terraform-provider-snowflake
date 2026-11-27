package sdk

import (
	"context"
	"errors"
)

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

// additionalValidations requires IN DATABASE, the only scope SHOW DATABASE ROLES accepts. The generated IN is optional and
// can also name an account, an application, or an application package.
func (opts *ShowDatabaseRoleOptions) additionalValidations() error {
	var errs []error
	if opts.In == nil || !ValidObjectIdentifier(opts.In.Database) || valueSet(opts.In.Account) || valueSet(opts.In.Application) || valueSet(opts.In.ApplicationPackage) {
		errs = append(errs, errInvalidIdentifier("ShowDatabaseRoleOptions", "Database"))
	}
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	return errors.Join(errs...)
}
