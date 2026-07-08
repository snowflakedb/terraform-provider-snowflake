package sdk

import (
	"context"
	"fmt"
)

func (v *sessions) ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, opts)
}

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

func (v *sessions) UseSecondaryRoles(ctx context.Context, opt SecondaryRoleOption) error {
	sql := fmt.Sprintf(`USE SECONDARY ROLES %s`, opt)
	_, err := v.client.exec(ctx, sql)
	return err
}
