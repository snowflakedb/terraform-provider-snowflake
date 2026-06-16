package sdk

import (
	"context"
	"database/sql"
	"time"
)

var (
	_ validatable = new(CreateDatabaseOptions)
	_ validatable = new(CreateSharedDatabaseOptions)
	_ validatable = new(CreateSecondaryDatabaseOptions)
	_ validatable = new(CreateDatabaseFromListingOptions)
	_ validatable = new(AlterDatabaseOptions)
	_ validatable = new(AlterDatabaseReplicationOptions)
	_ validatable = new(AlterDatabaseFailoverOptions)
	_ validatable = new(DropDatabaseOptions)
	_ validatable = new(undropDatabaseOptions)
	_ validatable = new(ShowDatabasesOptions)
	_ validatable = new(describeDatabaseOptions)

	_ convertibleRow[Database] = new(databaseRow)
)

type Databases interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateDatabaseOptions) error
	CreateShared(ctx context.Context, id AccountObjectIdentifier, shareID ExternalObjectIdentifier, opts *CreateSharedDatabaseOptions) error
	CreateSecondary(ctx context.Context, id AccountObjectIdentifier, primaryID ExternalObjectIdentifier, opts *CreateSecondaryDatabaseOptions) error
	CreateFromListing(ctx context.Context, id AccountObjectIdentifier, listingGlobalName string, opts *CreateDatabaseFromListingOptions) error
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

type Database struct {
	CreatedOn     time.Time
	Name          string
	IsDefault     bool
	IsCurrent     bool
	Origin        ObjectIdentifier
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

func (v *Database) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Database) ObjectType() ObjectType {
	return ObjectTypeDatabase
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
	LogEventLevel                           *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_EVENT_LEVEL"`
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
	LogEventLevel                           *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_EVENT_LEVEL"`
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
	LogEventLevel                           *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_EVENT_LEVEL"`
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

// CreateDatabaseFromListingOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-database.
// Supports both external listings and organization listings via the same FROM LISTING syntax.
// SQL: CREATE DATABASE <name> FROM LISTING '<listing_global_name>'
type CreateDatabaseFromListingOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"`
	database    bool                    `ddl:"static" sql:"DATABASE"`
	name        AccountObjectIdentifier `ddl:"identifier"`
	fromListing string                  `ddl:"parameter,single_quotes,no_equals" sql:"FROM LISTING"`
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
	LogEventLevel                           *LogLevel                   `ddl:"parameter,single_quotes" sql:"LOG_EVENT_LEVEL"`
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
	LogEventLevel                           *bool `ddl:"keyword" sql:"LOG_EVENT_LEVEL"`
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

// AlterDatabaseReplicationOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-database.
type AlterDatabaseReplicationOptions struct {
	alter              bool                    `ddl:"static" sql:"ALTER"`
	database           bool                    `ddl:"static" sql:"DATABASE"`
	name               AccountObjectIdentifier `ddl:"identifier"`
	EnableReplication  *EnableReplication      `ddl:"keyword" sql:"ENABLE REPLICATION"`
	DisableReplication *DisableReplication     `ddl:"keyword" sql:"DISABLE REPLICATION"`
	Refresh            *bool                   `ddl:"keyword" sql:"REFRESH"`
}

type EnableReplication struct {
	ToAccounts         []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
	IgnoreEditionCheck *bool               `ddl:"keyword" sql:"IGNORE EDITION CHECK"`
}

type DisableReplication struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
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

type EnableFailover struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
}

type DisableFailover struct {
	ToAccounts []AccountIdentifier `ddl:"keyword,no_parentheses" sql:"TO ACCOUNTS"`
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

// undropDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-database.
type undropDatabaseOptions struct {
	undrop   bool                    `ddl:"static" sql:"UNDROP"`
	database bool                    `ddl:"static" sql:"DATABASE"`
	name     AccountObjectIdentifier `ddl:"identifier"`
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

// describeDatabaseOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-database.
type describeDatabaseOptions struct {
	describe bool                    `ddl:"static" sql:"DESCRIBE"`
	database bool                    `ddl:"static" sql:"DATABASE"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

type DatabaseDetails struct {
	Rows []DatabaseDetailsRow
}

type DatabaseDetailsRow struct {
	CreatedOn time.Time
	Name      string
	Kind      string
}
