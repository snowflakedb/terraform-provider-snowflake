package sdk

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (v *postgresMirrors) DropSafely(ctx context.Context, request *DropPostgresMirrorRequest) error {
	if err := v.Drop(ctx, request); err != nil && !postgresMirrorNotFound(err) {
		return err
	}
	return nil
}

// postgresMirrorNotFound reports the procedure's missing-mirror failure. The message does not match the phrases
// decodeDriverError rewrites, so DropSafely recognizes it here.
func postgresMirrorNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "[NOT_FOUND] Mirror does not exist")
}

func (s *CreatePostgresMirrorRequest) GetName() string {
	return s.arguments.MirrorName
}

func (v *CreatePostgresMirrorArguments) additionalValidations() error {
	if v.ExternalVolume != nil && (v.TargetTableType == nil || *v.TargetTableType != PostgresMirrorTargetTableTypeIceberg) {
		return NewError("CreatePostgresMirrorOptions.arguments field: ExternalVolume can only be set when TargetTableType is ICEBERG")
	}
	return nil
}

// IsMembershipKnown reports whether DESCRIBE_MIRROR could read PostgresTables and PostgresSchemas from the publication
// on the source instance. Both are nil, not empty, when the source could not be read (or the SNOWFLAKE application
// predates these columns), so "unknown" can be told apart from "no tables". Entries use Postgres quoting.
func (v *PostgresMirrorDetails) IsMembershipKnown() bool {
	return v.PostgresTables != nil && v.PostgresSchemas != nil
}

// additionalConvert fills the fields marked WithManualConvert in postgres_mirrors_def.go.
func (r postgresMirrorDetailsRow) additionalConvert(result *PostgresMirrorDetails) error {
	// An unrecognized status is kept. ToPostgresMirrorStatus would reject it and fail Describe.
	result.Status = PostgresMirrorStatus(r.Status.String)
	result.RecentQueryDurations = compactPostgresMirrorJson(r.RecentQueryDurations)
	result.PostgresStatus = compactPostgresMirrorJson(r.PostgresStatus)

	var errs []error
	var err error
	if r.TargetTableType.Valid {
		if result.TargetTableType, err = ToPostgresMirrorTargetTableType(r.TargetTableType.String); err != nil {
			errs = append(errs, fmt.Errorf("TARGET_TABLE_TYPE: %w", err))
		}
	}
	if result.PostgresTables, err = parsePostgresMirrorNames(r.PostgresTables); err != nil {
		errs = append(errs, fmt.Errorf("POSTGRES_TABLES: %w", err))
	}
	if result.PostgresSchemas, err = parsePostgresMirrorNames(r.PostgresSchemas); err != nil {
		errs = append(errs, fmt.Errorf("POSTGRES_SCHEMAS: %w", err))
	}
	return errors.Join(errs...)
}

func parsePostgresMirrorNames(value sql.NullString) ([]string, error) {
	if !value.Valid {
		return nil, nil
	}
	names := make([]string, 0)
	if err := json.Unmarshal([]byte(value.String), &names); err != nil {
		return nil, fmt.Errorf("parsing %q as a JSON array of strings: %w", value.String, err)
	}
	return names, nil
}

// compactPostgresMirrorJson drops the indentation the driver adds to ARRAY and OBJECT values.
func compactPostgresMirrorJson(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	var compacted bytes.Buffer
	if err := json.Compact(&compacted, []byte(value.String)); err != nil {
		return String(value.String)
	}
	return String(compacted.String())
}
