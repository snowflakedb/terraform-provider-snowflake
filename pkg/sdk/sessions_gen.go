package sdk

import "context"

type Sessions interface {
	// Parameters
	AlterSession(ctx context.Context, opts *AlterSessionOptions) error
	ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error)
	// Context
	UseWarehouse(ctx context.Context, warehouse AccountObjectIdentifier) error
	UseDatabase(ctx context.Context, database AccountObjectIdentifier) error
	UseSchema(ctx context.Context, schema DatabaseObjectIdentifier) error
	UseRole(ctx context.Context, role AccountObjectIdentifier) error
	UseSecondaryRoles(ctx context.Context, opt SecondaryRoleOption) error
}

// AlterSessionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-session.
type AlterSessionOptions struct {
	alter   bool          `ddl:"static" sql:"ALTER"`
	session bool          `ddl:"static" sql:"SESSION"`
	Set     *SessionSet   `ddl:"keyword" sql:"SET"`
	Unset   *SessionUnset `ddl:"keyword" sql:"UNSET"`
}

type SessionSet struct {
	SessionParameters *SessionParameters `ddl:"list"`
}

type SessionUnset struct {
	SessionParametersUnset *SessionParametersUnset `ddl:"list"`
}
