package sdk

import (
	"context"
	"errors"
	"fmt"
)

var (
	_ validatable = new(AlterSessionOptions)
	_ validatable = new(ShowParametersOptions)
)

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

var _ Sessions = (*sessions)(nil)

type sessions struct {
	client *Client
}

// AlterSessionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-session.
type AlterSessionOptions struct {
	alter   bool          `ddl:"static" sql:"ALTER"`
	session bool          `ddl:"static" sql:"SESSION"`
	Set     *SessionSet   `ddl:"keyword" sql:"SET"`
	Unset   *SessionUnset `ddl:"keyword" sql:"UNSET"`
}

func (opts *AlterSessionOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueNil(opts.Set, opts.Unset) {
		errs = append(errs, errOneOf("AlterSessionOptions", "Set", "Unset"))
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type SessionSet struct {
	SessionParameters *SessionParameters `ddl:"list"`
}

func (v *SessionSet) validate() error {
	if err := v.SessionParameters.validate(); err != nil {
		return err
	}
	return nil
}

type SessionUnset struct {
	SessionParametersUnset *SessionParametersUnset `ddl:"list"`
}

func (v *SessionUnset) validate() error {
	if err := v.SessionParametersUnset.validate(); err != nil {
		return err
	}
	return nil
}

func (v *sessions) AlterSession(ctx context.Context, opts *AlterSessionOptions) error {
	if opts == nil {
		opts = &AlterSessionOptions{}
	}
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *sessions) ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, opts)
}

// Context
func (v *sessions) UseWarehouse(ctx context.Context, warehouse AccountObjectIdentifier) error {
	sql := fmt.Sprintf(`USE WAREHOUSE %s`, warehouse.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}

func (v *sessions) UseDatabase(ctx context.Context, database AccountObjectIdentifier) error {
	sql := fmt.Sprintf(`USE DATABASE %s`, database.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}

func (v *sessions) UseSchema(ctx context.Context, schema DatabaseObjectIdentifier) error {
	sql := fmt.Sprintf(`USE SCHEMA %s`, schema.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}

func (v *sessions) UseRole(ctx context.Context, role AccountObjectIdentifier) error {
	sql := fmt.Sprintf(`USE ROLE %s`, role.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}

// SecondaryRoleOption is based on https://docs.snowflake.com/en/sql-reference/sql/use-secondary-roles.
type SecondaryRoleOption string

const (
	SecondaryRolesAll  SecondaryRoleOption = "ALL"
	SecondaryRolesNone SecondaryRoleOption = "NONE"
)

func (v *sessions) UseSecondaryRoles(ctx context.Context, opt SecondaryRoleOption) error {
	sql := fmt.Sprintf(`USE SECONDARY ROLES %s`, opt)
	_, err := v.client.exec(ctx, sql)
	return err
}
