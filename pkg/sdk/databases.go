package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ validatable = new(CreateDatabaseOptions)
	_ validatable = new(CreateSharedDatabaseOptions)
	_ validatable = new(CreateSecondaryDatabaseOptions)
	_ validatable = new(AlterDatabaseOptions)
	_ validatable = new(AlterDatabaseReplicationOptions)
	_ validatable = new(AlterDatabaseFailoverOptions)
	_ validatable = new(DropDatabaseOptions)
	_ validatable = new(undropDatabaseOptions)
	_ validatable = new(ShowDatabasesOptions)
	_ validatable = new(describeDatabaseOptions)
)

type Databases interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error
	CreateShared(ctx context.Context, id AccountObjectIdentifier, shareID ExternalObjectIdentifier, opts *CreateSharedDatabaseOptions) error
	CreateSecondary(ctx context.Context, id AccountObjectIdentifier, primaryID ExternalObjectIdentifier, opts *CreateSecondaryDatabaseOptions) error
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseOptions) error
	AlterReplication(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseReplicationOptions) error
	AlterFailover(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseFailoverOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Undrop(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowDatabasesOptions) ([]Database, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Database, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Database, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) (*DatabaseDetails, error)
	Use(ctx context.Context, id AccountObjectIdentifier) error
	ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error)
}

var _ Databases = (*databases)(nil)

type databases struct {
	client *Client
}

type Database struct {
	CreatedOn     time.Time
	Name          string
	IsDefault     bool
	IsCurrent     bool
	Origin        *ExternalObjectIdentifier
	Owner         string
	Comment       string
	Options       string
	RetentionTime int
	ResourceGroup string
	DroppedOn     time.Time
	Transient     bool
	Kind          string
	OwnerRoleType string
}

func (v *Database) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Database) ObjectType() ObjectType {
	return ObjectTypeDatabase
}

type databaseRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	IsDefault     sql.NullString `db:"is_default"`
	IsCurrent     sql.NullString `db:"is_current"`
	Origin        sql.NullString `db:"origin"`
	Owner         sql.NullString `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime sql.NullString `db:"retention_time"`
	ResourceGroup sql.NullString `db:"resource_group"`
	DroppedOn     sql.NullTime   `db:"dropped_on"`
	Kind          sql.NullString `db:"kind"`
	OwnerRoleType sql.NullString `db:"owner_role_type"`
}

func (row databaseRow) convert() *Database {
	database := &Database{
		CreatedOn: row.CreatedOn,
		Name:      row.Name,
	}
	if row.IsDefault.Valid {
		database.IsDefault = row.IsDefault.String == "Y"
	}
	if row.IsCurrent.Valid {
		database.IsCurrent = row.IsCurrent.String == "Y"
	}
	if row.Origin.Valid && row.Origin.String != "" {
		originId, err := ParseExternalObjectIdentifier(row.Origin.String)
		if err != nil {
			// TODO(SNOW-1561641): Return error
			log.Printf("[DEBUG] unable to parse origin ID: %v", row.Origin.String)
		} else {
			database.Origin = &originId
		}
	}
	if row.Owner.Valid {
		database.Owner = row.Owner.String
	}
	if row.Comment.Valid {
		database.Comment = row.Comment.String
	}
	if row.Options.Valid {
		database.Options = row.Options.String
	}
	if row.RetentionTime.Valid {
		retentionTimeInt, err := strconv.Atoi(row.RetentionTime.String)
		if err != nil {
			database.RetentionTime = 0
		}
		database.RetentionTime = retentionTimeInt
	}
	if row.ResourceGroup.Valid {
		database.ResourceGroup = row.ResourceGroup.String
	}
	if row.DroppedOn.Valid {
		database.DroppedOn = row.DroppedOn.Time
	}
	if row.Options.Valid {
		parts := strings.Split(row.Options.String, ", ")
		for _, part := range parts {
			if part == "TRANSIENT" {
				database.Transient = true
			}
		}
	}
	if row.Kind.Valid {
		database.Kind = row.Kind.String
	}
	if row.OwnerRoleType.Valid {
		database.OwnerRoleType = row.OwnerRoleType.String
	}
	return database
}

type StorageSerializationPolicy string

const (
	StorageSerializationPolicyCompatible StorageSerializationPolicy = "COMPATIBLE"
	StorageSerializationPolicyOptimized  StorageSerializationPolicy = "OPTIMIZED"
)

func ToStorageSerializationPolicy(value string) (StorageSerializationPolicy, error) {
	switch strings.ToUpper(value) {
	case string(StorageSerializationPolicyCompatible):
		return StorageSerializationPolicyCompatible, nil
	case string(StorageSerializationPolicyOptimized):
		return StorageSerializationPolicyOptimized, nil
	default:
		return "", fmt.Errorf("unknown storage serialization policy: %s", value)
	}
}

var AllStorageSerializationPolicies = []StorageSerializationPolicy{
	StorageSerializationPolicyCompatible,
	StorageSerializationPolicyOptimized,
}

// CreateDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-database.
type CreateDatabaseOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	Transient   *bool                   `ddl:"keyword" sql:"TRANSIENT"`
	database    bool                    `ddl:"static" sql:"DATABASE"`
	IfNotExists *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`
	Clone       *Clone                  `ddl:"-"`

	// Parameters
	DataRetentionTimeInDays                 *int                        `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays              *int                        `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ExternalVolume                          *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"CATALOG"`
	ReplaceInvalidCharacters                *bool                       `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	DefaultDDLCollation                     *string                     `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	StorageSerializationPolicy              *StorageSerializationPolicy `ddl:"parameter" sql:"STORAGE_SERIALIZATION_POLICY"`
	LogLevel                                *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	TraceLevel                              *TraceLevel                 `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	SuspendTaskAfterNumFailures             *int                        `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	TaskAutoRetryAttempts                   *int                        `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskManagedInitialWarehouseSize     *WarehouseSize              `ddl:"parameter" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskTimeoutMs                       *int                        `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	UserTaskMinimumTriggerIntervalInSeconds *int                        `ddl:"parameter" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	QuotedIdentifiersIgnoreCase             *bool                       `ddl:"parameter" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	EnableConsoleOutput                     *bool                       `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`

	Comment *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag     []TagAssociation `ddl:"keyword,parentheses" sql:"TAG"`
}

func (opts *CreateDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Clone) {
		if err := opts.Clone.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateDatabaseOptions", "OrReplace", "IfNotExists"))
	}
	if opts.ExternalVolume != nil && !ValidObjectIdentifier(opts.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("CreateDatabaseOptions", "ExternalVolume"))
	}
	if opts.Catalog != nil && !ValidObjectIdentifier(opts.Catalog) {
		errs = append(errs, errInvalidIdentifier("CreateDatabaseOptions", "Catalog"))
	}
	return errors.Join(errs...)
}

func (v *databases) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error {
	if opts == nil {
		opts = &CreateDatabaseOptions{}
	}
	opts.name = id
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

// CreateSharedDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-database.
type CreateSharedDatabaseOptions struct {
	create      bool                     `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	Transient   *bool                    `ddl:"keyword" sql:"TRANSIENT"`
	database    bool                     `ddl:"static" sql:"DATABASE"`
	IfNotExists *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier  `ddl:"identifier"`
	fromShare   ExternalObjectIdentifier `ddl:"identifier" sql:"FROM SHARE"`

	// Parameters
	ExternalVolume                          *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"CATALOG"`
	ReplaceInvalidCharacters                *bool                       `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	DefaultDDLCollation                     *string                     `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	StorageSerializationPolicy              *StorageSerializationPolicy `ddl:"parameter" sql:"STORAGE_SERIALIZATION_POLICY"`
	LogLevel                                *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	TraceLevel                              *TraceLevel                 `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	SuspendTaskAfterNumFailures             *int                        `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	TaskAutoRetryAttempts                   *int                        `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskManagedInitialWarehouseSize     *WarehouseSize              `ddl:"parameter" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskTimeoutMs                       *int                        `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	UserTaskMinimumTriggerIntervalInSeconds *int                        `ddl:"parameter" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	QuotedIdentifiersIgnoreCase             *bool                       `ddl:"parameter" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	EnableConsoleOutput                     *bool                       `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`

	Comment *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag     []TagAssociation `ddl:"keyword,parentheses" sql:"TAG"`
}

func (opts *CreateSharedDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSharedDatabaseOptions", "OrReplace", "IfNotExists"))
	}
	if opts.ExternalVolume != nil && !ValidObjectIdentifier(opts.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("CreateSharedDatabaseOptions", "ExternalVolume"))
	}
	if opts.Catalog != nil && !ValidObjectIdentifier(opts.Catalog) {
		errs = append(errs, errInvalidIdentifier("CreateSharedDatabaseOptions", "Catalog"))
	}
	if !ValidObjectIdentifier(opts.fromShare) {
		errs = append(errs, errInvalidIdentifier("CreateSharedDatabaseOptions", "fromShare"))
	}
	return errors.Join(errs...)
}

func (v *databases) CreateShared(ctx context.Context, id AccountObjectIdentifier, shareID ExternalObjectIdentifier, opts *CreateSharedDatabaseOptions) error {
	if opts == nil {
		opts = &CreateSharedDatabaseOptions{}
	}

	opts.name = id
	opts.fromShare = shareID

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

// CreateSecondaryDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-database.
type CreateSecondaryDatabaseOptions struct {
	create          bool                     `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	Transient       *bool                    `ddl:"keyword" sql:"TRANSIENT"`
	database        bool                     `ddl:"static" sql:"DATABASE"`
	IfNotExists     *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            AccountObjectIdentifier  `ddl:"identifier"`
	primaryDatabase ExternalObjectIdentifier `ddl:"identifier" sql:"AS REPLICA OF"`

	// Parameters
	DataRetentionTimeInDays                 *int                        `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays              *int                        `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ExternalVolume                          *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"CATALOG"`
	ReplaceInvalidCharacters                *bool                       `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	DefaultDDLCollation                     *string                     `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	StorageSerializationPolicy              *StorageSerializationPolicy `ddl:"parameter" sql:"STORAGE_SERIALIZATION_POLICY"`
	LogLevel                                *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	TraceLevel                              *TraceLevel                 `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	SuspendTaskAfterNumFailures             *int                        `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	TaskAutoRetryAttempts                   *int                        `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskManagedInitialWarehouseSize     *WarehouseSize              `ddl:"parameter" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskTimeoutMs                       *int                        `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	UserTaskMinimumTriggerIntervalInSeconds *int                        `ddl:"parameter" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	QuotedIdentifiersIgnoreCase             *bool                       `ddl:"parameter" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	EnableConsoleOutput                     *bool                       `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`

	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateSecondaryDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.primaryDatabase) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "primaryDatabase"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateSecondaryDatabaseOptions", "OrReplace", "IfNotExists"))
	}
	if opts.ExternalVolume != nil && !ValidObjectIdentifier(opts.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "ExternalVolume"))
	}
	if opts.Catalog != nil && !ValidObjectIdentifier(opts.Catalog) {
		errs = append(errs, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "Catalog"))
	}
	return errors.Join(errs...)
}

func (v *databases) CreateSecondary(ctx context.Context, id AccountObjectIdentifier, primaryID ExternalObjectIdentifier, opts *CreateSecondaryDatabaseOptions) error {
	if opts == nil {
		opts = &CreateSecondaryDatabaseOptions{}
	}
	opts.name = id
	opts.primaryDatabase = primaryID
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

// AlterDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-database.
type AlterDatabaseOptions struct {
	alter    bool                     `ddl:"static" sql:"ALTER"`
	database bool                     `ddl:"static" sql:"DATABASE"`
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier  `ddl:"identifier"`
	NewName  *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SwapWith *AccountObjectIdentifier `ddl:"identifier" sql:"SWAP WITH"`
	Set      *DatabaseSet             `ddl:"list,no_parentheses" sql:"SET"`
	Unset    *DatabaseUnset           `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTag   []TagAssociation         `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier       `ddl:"keyword" sql:"UNSET TAG"`
}

func (opts *AlterDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.NewName != nil && !ValidObjectIdentifier(opts.NewName) {
		errs = append(errs, errInvalidIdentifier("AlterDatabaseOptions", "NewName"))
	}
	if opts.SwapWith != nil && !ValidObjectIdentifier(opts.SwapWith) {
		errs = append(errs, errInvalidIdentifier("AlterDatabaseOptions", "SwapWith"))
	}
	if !exactlyOneValueSet(opts.NewName, opts.Set, opts.Unset, opts.SwapWith, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterDatabaseOptions", "NewName", "Set", "Unset", "SwapWith", "SetTag", "UnsetTag"))
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

type DatabaseSet struct {
	// Parameters
	DataRetentionTimeInDays                 *int                        `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays              *int                        `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ExternalVolume                          *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"CATALOG"`
	ReplaceInvalidCharacters                *bool                       `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	DefaultDDLCollation                     *string                     `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	StorageSerializationPolicy              *StorageSerializationPolicy `ddl:"parameter" sql:"STORAGE_SERIALIZATION_POLICY"`
	LogLevel                                *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	TraceLevel                              *TraceLevel                 `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	SuspendTaskAfterNumFailures             *int                        `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	TaskAutoRetryAttempts                   *int                        `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskManagedInitialWarehouseSize     *WarehouseSize              `ddl:"parameter" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskTimeoutMs                       *int                        `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	UserTaskMinimumTriggerIntervalInSeconds *int                        `ddl:"parameter" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	QuotedIdentifiersIgnoreCase             *bool                       `ddl:"parameter" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	EnableConsoleOutput                     *bool                       `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`

	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (v *DatabaseSet) validate() error {
	var errs []error
	if v.ExternalVolume != nil && !ValidObjectIdentifier(v.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("DatabaseSet", "ExternalVolume"))
	}
	if v.Catalog != nil && !ValidObjectIdentifier(v.Catalog) {
		errs = append(errs, errInvalidIdentifier("DatabaseSet", "Catalog"))
	}
	if !anyValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ExternalVolume,
		v.Catalog,
		v.ReplaceInvalidCharacters,
		v.DefaultDDLCollation,
		v.StorageSerializationPolicy,
		v.LogLevel,
		v.TraceLevel,
		v.SuspendTaskAfterNumFailures,
		v.TaskAutoRetryAttempts,
		v.UserTaskManagedInitialWarehouseSize,
		v.UserTaskTimeoutMs,
		v.UserTaskMinimumTriggerIntervalInSeconds,
		v.QuotedIdentifiersIgnoreCase,
		v.EnableConsoleOutput,
		v.Comment,
	) {
		errs = append(errs, errAtLeastOneOf(
			"DatabaseSet",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"ReplaceInvalidCharacters",
			"DefaultDDLCollation",
			"StorageSerializationPolicy",
			"LogLevel",
			"TraceLevel",
			"SuspendTaskAfterNumFailures",
			"TaskAutoRetryAttempts",
			"UserTaskManagedInitialWarehouseSize",
			"UserTaskTimeoutMs",
			"UserTaskMinimumTriggerIntervalInSeconds",
			"QuotedIdentifiersIgnoreCase",
			"EnableConsoleOutput",
			"Comment",
		))
	}
	return errors.Join(errs...)
}

type DatabaseUnset struct {
	// Parameters
	DataRetentionTimeInDays                 *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays              *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ExternalVolume                          *bool `ddl:"keyword" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *bool `ddl:"keyword" sql:"CATALOG"`
	ReplaceInvalidCharacters                *bool `ddl:"keyword" sql:"REPLACE_INVALID_CHARACTERS"`
	DefaultDDLCollation                     *bool `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	StorageSerializationPolicy              *bool `ddl:"keyword" sql:"STORAGE_SERIALIZATION_POLICY"`
	LogLevel                                *bool `ddl:"keyword" sql:"LOG_LEVEL"`
	TraceLevel                              *bool `ddl:"keyword" sql:"TRACE_LEVEL"`
	SuspendTaskAfterNumFailures             *bool `ddl:"keyword" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	TaskAutoRetryAttempts                   *bool `ddl:"keyword" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskManagedInitialWarehouseSize     *bool `ddl:"keyword" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	UserTaskTimeoutMs                       *bool `ddl:"keyword" sql:"USER_TASK_TIMEOUT_MS"`
	UserTaskMinimumTriggerIntervalInSeconds *bool `ddl:"keyword" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	QuotedIdentifiersIgnoreCase             *bool `ddl:"keyword" sql:"QUOTED_IDENTIFIERS_IGNORE_CASE"`
	EnableConsoleOutput                     *bool `ddl:"keyword" sql:"ENABLE_CONSOLE_OUTPUT"`

	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

func (v *DatabaseUnset) validate() error {
	var errs []error
	if !anyValueSet(
		v.DataRetentionTimeInDays,
		v.MaxDataExtensionTimeInDays,
		v.ExternalVolume,
		v.Catalog,
		v.ReplaceInvalidCharacters,
		v.DefaultDDLCollation,
		v.StorageSerializationPolicy,
		v.LogLevel,
		v.TraceLevel,
		v.SuspendTaskAfterNumFailures,
		v.TaskAutoRetryAttempts,
		v.UserTaskManagedInitialWarehouseSize,
		v.UserTaskTimeoutMs,
		v.UserTaskMinimumTriggerIntervalInSeconds,
		v.QuotedIdentifiersIgnoreCase,
		v.EnableConsoleOutput,
		v.Comment,
	) {
		errs = append(errs, errAtLeastOneOf(
			"DatabaseUnset",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"ReplaceInvalidCharacters",
			"DefaultDDLCollation",
			"StorageSerializationPolicy",
			"LogLevel",
			"TraceLevel",
			"SuspendTaskAfterNumFailures",
			"TaskAutoRetryAttempts",
			"UserTaskManagedInitialWarehouseSize",
			"UserTaskTimeoutMs",
			"UserTaskMinimumTriggerIntervalInSeconds",
			"QuotedIdentifiersIgnoreCase",
			"EnableConsoleOutput",
			"Comment",
		))
	}
	return errors.Join(errs...)
}

func (v *databases) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseOptions) error {
	if opts == nil {
		opts = &AlterDatabaseOptions{}
	}
	opts.name = id
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

// AlterDatabaseReplicationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-database.
type AlterDatabaseReplicationOptions struct {
	alter              bool                    `ddl:"static" sql:"ALTER"`
	database           bool                    `ddl:"static" sql:"DATABASE"`
	name               AccountObjectIdentifier `ddl:"identifier"`
	EnableReplication  *EnableReplication      `ddl:"keyword" sql:"ENABLE REPLICATION"`
	DisableReplication *DisableReplication     `ddl:"keyword" sql:"DISABLE REPLICATION"`
	Refresh            *bool                   `ddl:"keyword" sql:"REFRESH"`
}

func (opts *AlterDatabaseReplicationOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.EnableReplication, opts.DisableReplication, opts.Refresh) {
		errs = append(errs, errExactlyOneOf("AlterDatabaseReplicationOptions", "EnableReplication", "DisableReplication", "Refresh"))
	}
	return errors.Join(errs...)
}

type EnableReplication struct {
	ToAccounts         []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
	IgnoreEditionCheck *bool               `ddl:"keyword" sql:"IGNORE EDITION CHECK"`
}

type DisableReplication struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
}

func (v *databases) AlterReplication(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseReplicationOptions) error {
	if opts == nil {
		opts = &AlterDatabaseReplicationOptions{}
	}
	opts.name = id
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

// AlterDatabaseFailoverOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-database.
type AlterDatabaseFailoverOptions struct {
	alter           bool                    `ddl:"static" sql:"ALTER"`
	database        bool                    `ddl:"static" sql:"DATABASE"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	EnableFailover  *EnableFailover         `ddl:"keyword" sql:"ENABLE FAILOVER"`
	DisableFailover *DisableFailover        `ddl:"keyword" sql:"DISABLE FAILOVER"`
	Primary         *bool                   `ddl:"keyword" sql:"PRIMARY"`
}

func (opts *AlterDatabaseFailoverOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.EnableFailover, opts.DisableFailover, opts.Primary) {
		errs = append(errs, errExactlyOneOf("AlterDatabaseFailoverOptions", "EnableFailover", "DisableFailover", "Primary"))
	}
	return errors.Join(errs...)
}

type EnableFailover struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
}

type DisableFailover struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
}

func (v *databases) AlterFailover(ctx context.Context, id AccountObjectIdentifier, opts *AlterDatabaseFailoverOptions) error {
	if opts == nil {
		opts = &AlterDatabaseFailoverOptions{}
	}
	opts.name = id
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

// DropDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-database.
type DropDatabaseOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`
	database bool                    `ddl:"static" sql:"DATABASE"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
	Cascade  *bool                   `ddl:"keyword" sql:"CASCADE"`
	Restrict *bool                   `ddl:"keyword" sql:"RESTRICT"`
}

func (opts *DropDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Cascade, opts.Restrict) {
		errs = append(errs, errOneOf("DropDatabaseOptions", "Cascade", "Restrict"))
	}
	return JoinErrors(errs...)
}

func (v *databases) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropDatabaseOptions) error {
	if opts == nil {
		opts = &DropDatabaseOptions{}
	}
	opts.name = id
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

func (v *databases) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropDatabaseOptions{IfExists: Bool(true)}) }, ctx, id)
}

// undropDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-database.
type undropDatabaseOptions struct {
	undrop   bool                    `ddl:"static" sql:"UNDROP"`
	database bool                    `ddl:"static" sql:"DATABASE"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *undropDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *databases) Undrop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &undropDatabaseOptions{
		name: id,
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

// ShowDatabasesOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-databases.
type ShowDatabasesOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	databases  bool       `ddl:"static" sql:"DATABASES"`
	History    *bool      `ddl:"keyword" sql:"HISTORY"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom  *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowDatabasesOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (v *databases) Show(ctx context.Context, opts *ShowDatabasesOptions) ([]Database, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[databaseRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[databaseRow, Database](dbRows)
}

func (v *databases) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Database, error) {
	databases, err := v.Show(ctx, &ShowDatabasesOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(databases, func(r Database) bool { return r.Name == id.Name() })
}

func (v *databases) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Database, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

type DatabaseDetails struct {
	Rows []DatabaseDetailsRow
}

type DatabaseDetailsRow struct {
	CreatedOn time.Time
	Name      string
	Kind      string
}

// describeDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-database.
type describeDatabaseOptions struct {
	describe bool                    `ddl:"static" sql:"DESCRIBE"`
	database bool                    `ddl:"static" sql:"DATABASE"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *describeDatabaseOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
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
