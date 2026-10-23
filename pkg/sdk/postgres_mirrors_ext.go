package sdk

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
)

func (v *postgresMirrors) DropSafely(ctx context.Context, request *DropPostgresMirrorRequest) error {
	if err := v.Drop(ctx, request); err != nil && !postgresMirrorNotFound(err) {
		return err
	}
	return nil
}

const postgresMirrorPollInterval = 10 * time.Second

func (v *postgresMirrors) WaitForInitialSnapshot(ctx context.Context, postgresInstance string, name string, minimumTables int) error {
	if minimumTables < 1 {
		return fmt.Errorf("minimumTables must be at least 1: LIST_MIRRORS reports 0 tables remaining before any table is registered")
	}
	normalizedName := strings.ToLower(name)
	_, err := util.PollUntil(ctx, postgresMirrorPollInterval, func() (*PostgresMirror, error) {
		mirrors, err := v.List(ctx, NewListPostgresMirrorRequest(*NewListPostgresMirrorArgumentsRequest().WithPostgresInstance(postgresInstance)))
		if err != nil {
			return nil, err
		}
		mirror, err := collections.FindFirst(mirrors, func(mirror PostgresMirror) bool {
			return mirror.Name == normalizedName
		})
		if err != nil {
			return nil, fmt.Errorf("%w: postgres mirror %s on instance %s", ErrObjectNotFound, name, postgresInstance)
		}
		// A suspended apply task never finishes the snapshot; the procedures suspend it after repeated failures.
		if mirror.Status == PostgresMirrorStatusSuspended {
			suspended := fmt.Errorf("postgres mirror %s is suspended, so its initial snapshot cannot complete", name)
			if mirror.ErrorMessage != nil && *mirror.ErrorMessage != "" {
				suspended = fmt.Errorf("%w: %s", suspended, *mirror.ErrorMessage)
			}
			return nil, suspended
		}
		return mirror, nil
	}, func(mirror *PostgresMirror) bool {
		return mirror.InitialSnapshotCompleted(minimumTables)
	}, fmt.Sprintf("initial snapshot of postgres mirror %s did not complete", name))
	return err
}

// InitialSnapshotCompleted needs more than "0 tables remaining": LIST_MIRRORS also reports 0 before the apply task has
// registered any table, so the mirror must be ACTIVE and have registered at least minimumTables tables.
// Callers must pass minimumTables >= 1; WaitForInitialSnapshot enforces that before calling this.
// INITIAL_SNAPSHOT_TABLES_REMAINING counts only the first copy of each table, not later re-snapshots.
func (v *PostgresMirror) InitialSnapshotCompleted(minimumTables int) bool {
	return v.Status == PostgresMirrorStatusActive &&
		v.InitialSnapshotTablesRemaining != nil && *v.InitialSnapshotTablesRemaining == 0 &&
		v.TableCount != nil && *v.TableCount >= minimumTables
}

func (r postgresMirrorRow) additionalConvert(result *PostgresMirror) error {
	// Normalized so WaitForInitialSnapshot's comparisons against the uppercase status constants still match.
	// An unrecognized status is kept (uppercased). ToPostgresMirrorStatus would reject it and fail List/Wait.
	result.Status = PostgresMirrorStatus(strings.ToUpper(r.Status.String))
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

// ALTER_MIRROR rejects a call that changes nothing. UnsetWarehouse counts only when true: false is omitted from the SQL,
// so it neither satisfies this check nor conflicts with Warehouse.
func (v *AlterPostgresMirrorArguments) additionalValidations() error {
	unsetWarehouse := v.UnsetWarehouse != nil && *v.UnsetWarehouse
	if valueSet(v.Warehouse) && unsetWarehouse {
		return errOneOf("AlterPostgresMirrorOptions.arguments", "Warehouse", "UnsetWarehouse")
	}
	if !unsetWarehouse && !anyValueSet(v.RefreshInterval, v.AddTables, v.RemoveTables, v.AddSchemas, v.RemoveSchemas, v.Warehouse) {
		return errAtLeastOneOf("AlterPostgresMirrorOptions.arguments", "RefreshInterval", "AddTables", "RemoveTables", "AddSchemas", "RemoveSchemas", "Warehouse", "UnsetWarehouse")
	}
	return nil
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
