package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ DatabaseRoles                = (*databaseRoles)(nil)
	_ convertibleRow[DatabaseRole] = new(databaseRoleDBRow)
)

type databaseRoles struct {
	client *Client
}

func (v *databaseRoles) Create(ctx context.Context, request *CreateDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Alter(ctx context.Context, request *AlterDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Drop(ctx context.Context, request *DropDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) DropSafely(ctx context.Context, id DatabaseObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropDatabaseRoleRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *databaseRoles) Show(ctx context.Context, request *ShowDatabaseRoleRequest) ([]DatabaseRole, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[databaseRoleDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	resultList, err := convertRows[databaseRoleDBRow, DatabaseRole](dbRows)
	if err != nil {
		return nil, err
	}

	for i := range resultList {
		resultList[i].DatabaseName = request.database.name
	}

	return resultList, nil
}

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

func (v *databaseRoles) Grant(ctx context.Context, request *GrantDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) Revoke(ctx context.Context, request *RevokeDatabaseRoleRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) GrantToShare(ctx context.Context, request *GrantDatabaseRoleToShareRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *databaseRoles) RevokeFromShare(ctx context.Context, request *RevokeDatabaseRoleFromShareRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (s *CreateDatabaseRoleRequest) toOpts() *createDatabaseRoleOptions {
	return &createDatabaseRoleOptions{
		OrReplace:   Bool(s.orReplace),
		IfNotExists: Bool(s.ifNotExists),
		name:        s.name,
		Comment:     s.comment,
	}
}

func (s *AlterDatabaseRoleRequest) toOpts() *alterDatabaseRoleOptions {
	opts := alterDatabaseRoleOptions{
		IfExists:  Bool(s.ifExists),
		name:      s.name,
		Rename:    s.rename,
		SetTags:   s.setTags,
		UnsetTags: s.unsetTags,
	}
	if s.set != nil {
		opts.Set = &DatabaseRoleSet{s.set.comment}
	}
	if s.unset != nil {
		opts.Unset = &DatabaseRoleUnset{true}
	}
	return &opts
}

func (s *DropDatabaseRoleRequest) toOpts() *dropDatabaseRoleOptions {
	return &dropDatabaseRoleOptions{
		IfExists: Bool(s.ifExists),
		name:     s.name,
	}
}

func (s *ShowDatabaseRoleRequest) toOpts() *showDatabaseRoleOptions {
	return &showDatabaseRoleOptions{
		Like:     s.like,
		Database: s.database,
		Limit:    s.limit,
	}
}

func (s *GrantDatabaseRoleRequest) toOpts() *grantDatabaseRoleOptions {
	opts := grantDatabaseRoleOptions{
		name: s.name,
	}

	grantToRole := grantOrRevokeDatabaseRoleObject{}
	if s.databaseRole != nil {
		grantToRole.DatabaseRoleName = s.databaseRole
	}
	if s.accountRole != nil {
		grantToRole.AccountRoleName = s.accountRole
	}
	opts.ParentRole = grantToRole

	return &opts
}

func (s *RevokeDatabaseRoleRequest) toOpts() *revokeDatabaseRoleOptions {
	opts := revokeDatabaseRoleOptions{
		name: s.name,
	}

	revokeFromRole := grantOrRevokeDatabaseRoleObject{}
	if s.databaseRole != nil {
		revokeFromRole.DatabaseRoleName = s.databaseRole
	}
	if s.accountRole != nil {
		revokeFromRole.AccountRoleName = s.accountRole
	}
	opts.ParentRole = revokeFromRole

	return &opts
}

func (s *GrantDatabaseRoleToShareRequest) toOpts() *grantDatabaseRoleToShareOptions {
	return &grantDatabaseRoleToShareOptions{
		name:  s.name,
		Share: s.share,
	}
}

func (s *RevokeDatabaseRoleFromShareRequest) toOpts() *revokeDatabaseRoleFromShareOptions {
	return &revokeDatabaseRoleFromShareOptions{
		name:  s.name,
		Share: s.share,
	}
}
