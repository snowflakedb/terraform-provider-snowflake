package sdk

import (
	"context"
	"strings"
)

func (v *Database) SetTransient(value string) {
	parts := strings.Split(value, ", ")
	for _, part := range parts {
		if part == "TRANSIENT" {
			v.Transient = true
		}
	}
}

// Use is based on https://docs.snowflake.com/en/sql-reference/sql/use-database.
func (v *databases) Use(ctx context.Context, id AccountObjectIdentifier) error {
	// proxy to sessions
	return v.client.Sessions.UseDatabase(ctx, id)
}

func (v *databases) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Database: id,
		},
	})
}

func (v *databases) Describe(ctx context.Context, id AccountObjectIdentifier) (*DatabaseDetails, error) {
	opts := &describeDatabaseOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []DatabaseDetailsRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	details := DatabaseDetails{
		Rows: rows,
	}
	return &details, err
}
