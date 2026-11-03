package sdk

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (v *postgresMirrors) DropSafely(ctx context.Context, request *DropPostgresMirrorRequest) error {
	if err := v.Drop(ctx, request); err != nil && !postgresMirrorNotFound(err) {
		return err
	}
	return nil
}

func (r postgresMirrorRow) additionalConvert(result *PostgresMirror) error {
	if r.Status.Valid {
		status, err := ToPostgresMirrorStatus(r.Status.String)
		if err != nil {
			return fmt.Errorf("STATUS: %w", err)
		}
		result.Status = status
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

// PostgresMirrorRefreshInterval is a refresh interval in one unit. CREATE_MIRROR and ALTER_MIRROR store the phrase they
// receive, and DESCRIBE_MIRROR returns it unchanged.
type PostgresMirrorRefreshInterval struct {
	Hours   int
	Minutes int
	Seconds int
}

// ParsePostgresMirrorRefreshInterval reads a day as 24 hours. It converts DESCRIBE_MIRROR output and checks no range,
// so it must not validate input.
func ParsePostgresMirrorRefreshInterval(interval string) (PostgresMirrorRefreshInterval, error) {
	invalid := fmt.Errorf("invalid postgres mirror refresh interval: %q", interval)
	parts := strings.Fields(strings.ToLower(interval))
	if len(parts) != 2 {
		return PostgresMirrorRefreshInterval{}, invalid
	}
	count, err := strconv.Atoi(parts[0])
	if err != nil {
		return PostgresMirrorRefreshInterval{}, invalid
	}
	switch parts[1] {
	case "hour", "hours":
		return PostgresMirrorRefreshInterval{Hours: count}, nil
	case "minute", "minutes":
		return PostgresMirrorRefreshInterval{Minutes: count}, nil
	case "second", "seconds":
		return PostgresMirrorRefreshInterval{Seconds: count}, nil
	case "day", "days":
		return PostgresMirrorRefreshInterval{Hours: 24 * count}, nil
	default:
		return PostgresMirrorRefreshInterval{}, invalid
	}
}

// String spells the phrase like the Snowflake documentation ("1 minute", "10 minutes"), because describe_output shows it
// as sent.
func (i PostgresMirrorRefreshInterval) String() string {
	switch {
	case i.Hours > 0:
		return postgresMirrorRefreshIntervalPhrase(i.Hours, "hour")
	case i.Minutes > 0:
		return postgresMirrorRefreshIntervalPhrase(i.Minutes, "minute")
	default:
		return postgresMirrorRefreshIntervalPhrase(i.Seconds, "second")
	}
}

func postgresMirrorRefreshIntervalPhrase(count int, unit string) string {
	if count == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", count, unit)
}

// additionalConvert fills the fields marked WithManualConvert in postgres_mirrors_def.go.
func (r postgresMirrorDetailsRow) additionalConvert(result *PostgresMirrorDetails) error {
	result.RecentQueryDurations = compactPostgresMirrorJson(r.RecentQueryDurations)
	result.PostgresStatus = compactPostgresMirrorJson(r.PostgresStatus)

	var errs []error
	var err error
	if r.Status.Valid {
		if result.Status, err = ToPostgresMirrorStatus(r.Status.String); err != nil {
			errs = append(errs, fmt.Errorf("STATUS: %w", err))
		}
	}
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
