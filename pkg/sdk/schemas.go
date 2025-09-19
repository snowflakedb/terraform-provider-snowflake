package sdk

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ validatable = new(CreateSchemaOptions)
	_ validatable = new(AlterSchemaOptions)
	_ validatable = new(DropSchemaOptions)
	_ validatable = new(undropSchemaOptions)
	_ validatable = new(describeSchemaOptions)
	_ validatable = new(ShowSchemaOptions)
)

type Schemas interface {
	Create(ctx context.Context, id DatabaseObjectIdentifier, opts *CreateSchemaOptions) error
	Alter(ctx context.Context, id DatabaseObjectIdentifier, opts *AlterSchemaOptions) error
	Drop(ctx context.Context, id DatabaseObjectIdentifier, opts *DropSchemaOptions) error
	DropSafely(ctx context.Context, id DatabaseObjectIdentifier) error
	Undrop(ctx context.Context, id DatabaseObjectIdentifier) error
	Describe(ctx context.Context, id DatabaseObjectIdentifier) ([]SchemaDetails, error)
	Show(ctx context.Context, opts *ShowSchemaOptions) ([]Schema, error)
	ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*Schema, error)
	ShowByIDSafely(ctx context.Context, id DatabaseObjectIdentifier) (*Schema, error)
	Use(ctx context.Context, id DatabaseObjectIdentifier) error
	ShowParameters(ctx context.Context, id DatabaseObjectIdentifier) ([]*Parameter, error)
}

var _ Schemas = (*schemas)(nil)

type schemas struct {
	client *Client
}

type Schema struct {
	CreatedOn     time.Time
	DroppedOn     time.Time
	Name          string
	IsDefault     bool
	IsCurrent     bool
	DatabaseName  string
	Owner         string
	Comment       string
	Options       *string
	RetentionTime string
	OwnerRoleType string
}

func (s *Schema) IsTransient() bool {
	if s.Options == nil {
		return false
	}
	return slices.Contains(ParseCommaSeparatedStringArray(*s.Options, false), "TRANSIENT")
}

func (s *Schema) IsManagedAccess() bool {
	if s.Options == nil {
		return false
	}
	return slices.Contains(ParseCommaSeparatedStringArray(*s.Options, false), "MANAGED ACCESS")
}

func (v *Schema) ID() DatabaseObjectIdentifier {
	return NewDatabaseObjectIdentifier(v.DatabaseName, v.Name)
}

func (v *Schema) ObjectType() ObjectType {
	return ObjectTypeSchema
}

type schemaDBRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	DroppedOn     sql.NullTime   `db:"dropped_on"`
	Name          string         `db:"name"`
	IsDefault     string         `db:"is_default"`
	IsCurrent     string         `db:"is_current"`
	DatabaseName  string         `db:"database_name"`
	Owner         string         `db:"owner"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime string         `db:"retention_time"`
	OwnerRoleType string         `db:"owner_role_type"`
}

func (row schemaDBRow) toSchema() Schema {
	schema := Schema{
		CreatedOn:     row.CreatedOn,
		Name:          row.Name,
		IsDefault:     row.IsDefault == "Y",
		IsCurrent:     row.IsCurrent == "Y",
		DatabaseName:  row.DatabaseName,
		Owner:         row.Owner,
		RetentionTime: row.RetentionTime,
		OwnerRoleType: row.OwnerRoleType,
	}
	if row.Comment.Valid {
		schema.Comment = row.Comment.String
	}
	if row.Options.Valid {
		schema.Options = &row.Options.String
	}
	if row.DroppedOn.Valid {
		schema.DroppedOn = row.DroppedOn.Time
	}
	return schema
}

// CreateSchemaOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-schema
type CreateSchemaOptions struct {
	create            bool                     `ddl:"static" sql:"CREATE"`
	OrReplace         *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	Transient         *bool                    `ddl:"keyword" sql:"TRANSIENT"`
	schema            bool                     `ddl:"static" sql:"SCHEMA"`
	IfNotExists       *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name              DatabaseObjectIdentifier `ddl:"identifier"`
	Clone             *Clone                   `ddl:"-"`
	WithManagedAccess *bool                    `ddl:"keyword" sql:"WITH MANAGED ACCESS"`

	// Parameters
	DataRetentionTimeInDays                 *int                        `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays              *int                        `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ExternalVolume                          *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"CATALOG"`
	PipeExecutionPaused                     *bool                       `ddl:"parameter" sql:"PIPE_EXECUTION_PAUSED"`
	ReplaceInvalidCharacters                *bool                       `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	DefaultDDLCollation                     *StringAllowEmpty           `ddl:"parameter" sql:"DEFAULT_DDL_COLLATION"`
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

func (opts *CreateSchemaOptions) validate() error {
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
		errs = append(errs, errOneOf("CreateSchemaOptions", "IfNotExists", "OrReplace"))
	}
	if opts.ExternalVolume != nil && !ValidObjectIdentifier(opts.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("CreateSchemaOptions", "ExternalVolume"))
	}
	if opts.Catalog != nil && !ValidObjectIdentifier(opts.Catalog) {
		errs = append(errs, errInvalidIdentifier("CreateSchemaOptions", "Catalog"))
	}
	return errors.Join(errs...)
}

func (v *schemas) Create(ctx context.Context, id DatabaseObjectIdentifier, opts *CreateSchemaOptions) error {
	if opts == nil {
		opts = &CreateSchemaOptions{}
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

// AlterSchemaOptions based on https://docs.snowflake.com/en/sql-reference/sql/alter-schema
type AlterSchemaOptions struct {
	alter    bool                      `ddl:"static" sql:"ALTER"`
	schema   bool                      `ddl:"static" sql:"SCHEMA"`
	IfExists *bool                     `ddl:"keyword" sql:"IF EXISTS"`
	name     DatabaseObjectIdentifier  `ddl:"identifier"`
	NewName  *DatabaseObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SwapWith *DatabaseObjectIdentifier `ddl:"identifier" sql:"SWAP WITH"`
	Set      *SchemaSet                `ddl:"list,no_parentheses" sql:"SET"`
	Unset    *SchemaUnset              `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTag   []TagAssociation          `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier        `ddl:"keyword" sql:"UNSET TAG"`
	// One of
	EnableManagedAccess  *bool `ddl:"keyword" sql:"ENABLE MANAGED ACCESS"`
	DisableManagedAccess *bool `ddl:"keyword" sql:"DISABLE MANAGED ACCESS"`
}

func (opts *AlterSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if opts.NewName != nil && !ValidObjectIdentifier(opts.NewName) {
		errs = append(errs, errInvalidIdentifier("AlterSchemaOptions", "NewName"))
	}
	if opts.SwapWith != nil && !ValidObjectIdentifier(opts.SwapWith) {
		errs = append(errs, errInvalidIdentifier("AlterSchemaOptions", "SwapWith"))
	}
	if !exactlyOneValueSet(opts.NewName, opts.SwapWith, opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag, opts.EnableManagedAccess, opts.DisableManagedAccess) {
		errs = append(errs, errExactlyOneOf("AlterSchemaOptions", "NewName", "SwapWith", "Set", "Unset", "SetTag", "UnsetTag", "EnableManagedAccess", "DisableManagedAccess"))
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

type SchemaSet struct {
	// Parameters
	DataRetentionTimeInDays                 *int                        `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays              *int                        `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ExternalVolume                          *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *AccountObjectIdentifier    `ddl:"identifier,equals" sql:"CATALOG"`
	PipeExecutionPaused                     *bool                       `ddl:"parameter" sql:"PIPE_EXECUTION_PAUSED"`
	ReplaceInvalidCharacters                *bool                       `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	DefaultDDLCollation                     *StringAllowEmpty           `ddl:"parameter" sql:"DEFAULT_DDL_COLLATION"`
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

func (v *SchemaSet) validate() error {
	var errs []error
	if v.ExternalVolume != nil && !ValidObjectIdentifier(v.ExternalVolume) {
		errs = append(errs, errInvalidIdentifier("SchemaSet", "ExternalVolume"))
	}
	if v.Catalog != nil && !ValidObjectIdentifier(v.Catalog) {
		errs = append(errs, errInvalidIdentifier("SchemaSet", "Catalog"))
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
		v.PipeExecutionPaused,
		v.Comment,
	) {
		errs = append(errs, errAtLeastOneOf(
			"SchemaSet",
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
			"PipeExecutionPaused",
			"Comment",
		))
	}
	return errors.Join(errs...)
}

type SchemaUnset struct {
	// Parameters
	DataRetentionTimeInDays                 *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays              *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ExternalVolume                          *bool `ddl:"keyword" sql:"EXTERNAL_VOLUME"`
	Catalog                                 *bool `ddl:"keyword" sql:"CATALOG"`
	PipeExecutionPaused                     *bool `ddl:"keyword" sql:"PIPE_EXECUTION_PAUSED"`
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

func (v *SchemaUnset) validate() error {
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
		v.PipeExecutionPaused,
		v.Comment,
	) {
		return errAtLeastOneOf(
			"SchemaUnset",
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
			"PipeExecutionPaused",
			"Comment",
		)
	}
	return nil
}

func (v *schemas) Alter(ctx context.Context, id DatabaseObjectIdentifier, opts *AlterSchemaOptions) error {
	if opts == nil {
		opts = &AlterSchemaOptions{}
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

// DropSchemaOptions Based on https://docs.snowflake.com/en/sql-reference/sql/drop-schema
type DropSchemaOptions struct {
	drop     bool                     `ddl:"static" sql:"DROP"`
	schema   bool                     `ddl:"static" sql:"SCHEMA"`
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`
	// one of
	Cascade  *bool `ddl:"static" sql:"CASCADE"`
	Restrict *bool `ddl:"static" sql:"RESTRICT"`
}

func (opts *DropSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Cascade, opts.Restrict) {
		errs = append(errs, errOneOf("DropSchemaOptions", "Cascade", "Restrict"))
	}
	return errors.Join(errs...)
}

func (v *schemas) Drop(ctx context.Context, id DatabaseObjectIdentifier, opts *DropSchemaOptions) error {
	if opts == nil {
		opts = &DropSchemaOptions{}
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

func (v *schemas) DropSafely(ctx context.Context, id DatabaseObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropSchemaOptions{IfExists: Bool(true)}) }, ctx, id)
}

// undropSchemaOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-schema.
type undropSchemaOptions struct {
	undrop bool                     `ddl:"static" sql:"UNDROP"`
	schema bool                     `ddl:"static" sql:"SCHEMA"`
	name   DatabaseObjectIdentifier `ddl:"identifier"`
}

func (opts *undropSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *schemas) Undrop(ctx context.Context, id DatabaseObjectIdentifier) error {
	opts := &undropSchemaOptions{
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

// describeSchemaOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-schema.
type describeSchemaOptions struct {
	describe bool                     `ddl:"static" sql:"DESCRIBE"`
	database bool                     `ddl:"static" sql:"SCHEMA"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`
}

func (opts *describeSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

type SchemaDetails struct {
	CreatedOn time.Time `db:"created_on"`
	Name      string    `db:"name"`
	Kind      string    `db:"kind"`
}

func (v *schemas) Describe(ctx context.Context, id DatabaseObjectIdentifier) ([]SchemaDetails, error) {
	opts := &describeSchemaOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var details []SchemaDetails
	err = v.client.query(ctx, &details, sql)
	if err != nil {
		return nil, err
	}
	return details, err
}

type SchemaIn struct {
	Account            *bool                   `ddl:"keyword" sql:"ACCOUNT"`
	Database           *bool                   `ddl:"keyword" sql:"DATABASE"`
	Application        *bool                   `ddl:"keyword" sql:"APPLICATION"`
	ApplicationPackage *bool                   `ddl:"keyword" sql:"APPLICATION_PACKAGE"`
	Name               AccountObjectIdentifier `ddl:"identifier"`
}

// ShowSchemaOptions based on https://docs.snowflake.com/en/sql-reference/sql/show-schemas
type ShowSchemaOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	Terse      *bool      `ddl:"keyword" sql:"TERSE"`
	schemas    bool       `ddl:"static" sql:"SCHEMAS"`
	History    *bool      `ddl:"keyword" sql:"HISTORY"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *SchemaIn  `ddl:"keyword" sql:"IN"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom  *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (opts *ShowSchemaOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (v *schemas) Show(ctx context.Context, opts *ShowSchemaOptions) ([]Schema, error) {
	if opts == nil {
		opts = &ShowSchemaOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []schemaDBRow
	err = v.client.query(ctx, &rows, sql)
	schemas := make([]Schema, len(rows))
	for i, row := range rows {
		schemas[i] = row.toSchema()
	}
	return schemas, err
}

func (v *schemas) ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*Schema, error) {
	schemas, err := v.client.Schemas.Show(ctx, &ShowSchemaOptions{
		In: &SchemaIn{
			Database: Bool(true),
			Name:     id.DatabaseId(),
		},
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(schemas, func(r Schema) bool { return r.Name == id.Name() })
}

func (v *schemas) ShowByIDSafely(ctx context.Context, id DatabaseObjectIdentifier) (*Schema, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *schemas) Use(ctx context.Context, id DatabaseObjectIdentifier) error {
	return v.client.Sessions.UseSchema(ctx, id)
}

func (v *schemas) ShowParameters(ctx context.Context, id DatabaseObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Schema: id,
		},
	})
}
