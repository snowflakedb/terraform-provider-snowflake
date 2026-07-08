package sdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// describeDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-database.
type describeDatabaseOptions struct {
	describe bool                    `ddl:"static" sql:"DESCRIBE"`
	database bool                    `ddl:"static" sql:"DATABASE"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *describeDatabaseOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	if !ValidObjectIdentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

type DatabaseDetails struct {
	Rows []DatabaseDetailsRow
}

type DatabaseDetailsRow struct {
	CreatedOn time.Time
	Name      string
	Kind      string
}

func (v *Database) SetTransient(value string) {
	parts := strings.SplitSeq(value, ", ")
	for part := range parts {
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

// additionalConvert handles manual field conversions for fields declared with WithManualConvert() in databases_def.go.
// Called by the generated convert() in databases_impl_gen.go.
func (r databaseRow) additionalConvert(db *Database) error {
	if r.Origin.Valid && r.Origin.String != "" && r.Origin.String != "<revoked>" {
		originId, err := ParseObjectIdentifierString(r.Origin.String)
		if err != nil {
			return fmt.Errorf("unable to parse origin ID: %w", err)
		}
		db.Origin = originId
	}
	if r.RetentionTime.Valid {
		retentionTimeInt, err := strconv.Atoi(r.RetentionTime.String)
		if err != nil {
			return fmt.Errorf("unable to parse retention time: %w", err)
		}
		db.RetentionTime = retentionTimeInt
	}
	if r.Options.Valid {
		db.SetTransient(r.Options.String)
	}
	return nil
}

func (opts *CloneDatabaseOptions) additionalValidations() error {
	return opts.Clone.validate()
}

func (opts *CreateFromListingDatabaseOptions) additionalValidations() error {
	if opts.FromListing == "" {
		return fmt.Errorf("CreateFromListingDatabaseOptions: listing global name must not be empty")
	}
	return nil
}

func (s *CreateDatabaseRequest) ID() AccountObjectIdentifier {
	return s.name
}
