package sdk

import (
	"context"
	"database/sql"
	"time"
)

// --- Manually-defined types for CREATE HYBRID TABLE body ---
//
// IMPORTANT: These structs are defined here because the SDK generator does NOT
// recursively generate nested struct types when using QueryStructField or
// ListQueryStructField. The canonical DSL definitions for these structures are in:
//
//   pkg/sdk/generator/defs/hybrid_tables_def.go
//
// When modifying these structures:
//   1. Update the DSL definition in hybrid_tables_def.go FIRST
//   2. Then update the Go struct here to match
//
// The DSL definitions serve as the source of truth and documentation.
// These Go structs are referenced by the generated CreateHybridTableOptions
// via PredefinedQueryStructField "ColumnsAndConstraints".

// HybridTableColumnsConstraintsAndIndexes is the parenthesized body of CREATE HYBRID TABLE,
// containing column definitions, out-of-line constraints, and out-of-line indexes.
// DSL definition: hybridTableColumnsConstraintsAndIndexesDef in hybrid_tables_def.go
type HybridTableColumnsConstraintsAndIndexes struct {
	Columns             []HybridTableColumn              `ddl:"keyword"`
	OutOfLineConstraint []HybridTableOutOfLineConstraint `ddl:"keyword"`
	OutOfLineIndex      []HybridTableOutOfLineIndex      `ddl:"keyword"`
}

// HybridTableColumn defines a single column in a hybrid table.
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table
// DSL definition: hybridTableColumnDef in hybrid_tables_def.go
type HybridTableColumn struct {
	Name             string                             `ddl:"keyword"`
	Type             DataType                           `ddl:"keyword"`
	InlineConstraint *HybridTableColumnInlineConstraint `ddl:"keyword"`
	NotNull          *bool                              `ddl:"keyword" sql:"NOT NULL"`
	DefaultValue     *ColumnDefaultValue                `ddl:"keyword"`
	Collate          *string                            `ddl:"parameter,no_equals,single_quotes" sql:"COLLATE"`
	Comment          *string                            `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
}

// HybridTableColumnInlineConstraint defines inline PRIMARY KEY, UNIQUE, or FOREIGN KEY on a column.
// For hybrid tables, NOT ENFORCED is invalid — all constraints are enforced.
// DSL definition: hybridTableColumnInlineConstraintDef in hybrid_tables_def.go
type HybridTableColumnInlineConstraint struct {
	Name       *string              `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	ForeignKey *InlineForeignKey    `ddl:"keyword" sql:"FOREIGN KEY"`
}

// HybridTableOutOfLineConstraint reuses OutOfLineConstraint from tables.go.
// The first 4 fields (Name, Type, Columns, ForeignKey) are identical in both struct and DDL tags.
// The 12 additional enforcement fields (Enforced, NotEnforced, Deferrable, etc.) in OutOfLineConstraint
// are all *bool pointers — when nil, structToSQL skips them, producing identical SQL output.
// For hybrid tables, enforcement fields are invalid (all constraints are enforced), so they stay nil.
// HybridTableConstraintActionRename is kept separate because its DDL generation context differs
// (static "RENAME CONSTRAINT" prefix vs parent field tag in tables.go).
type HybridTableOutOfLineConstraint = OutOfLineConstraint

// HybridTableOutOfLineIndex defines a secondary index in CREATE HYBRID TABLE.
// Syntax: INDEX <name> (<cols>) [INCLUDE (<cols>)]
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table
// DSL definition: hybridTableOutOfLineIndexDef in hybrid_tables_def.go
type HybridTableOutOfLineIndex struct {
	index          bool     `ddl:"static" sql:"INDEX"`
	Name           string   `ddl:"keyword"`
	Columns        []string `ddl:"keyword,parentheses"`
	IncludeColumns []string `ddl:"keyword,parentheses" sql:"INCLUDE"`
}

// --- Structs requiring manual DDL tag adjustments ---
// These structs are moved from hybrid_tables_gen.go because the generator cannot produce
// the correct DDL tags for certain Snowflake SQL syntax patterns.
// The canonical DSL definitions remain in hybrid_tables_def.go for documentation.

// HybridTableConstraintActionRename mirrors the structure of TableConstraintRenameAction in tables.go.
// These types cannot be shared directly because the DDL generation context differs:
// in tables.go the RENAME CONSTRAINT keyword is on the parent field tag, while here
// it is embedded as a static SQL prefix within the struct itself.
type HybridTableConstraintActionRename struct {
	renameConstraint bool   `ddl:"static" sql:"RENAME CONSTRAINT"`
	OldName          string `ddl:"keyword"`
	// Manually adjusted: generator produces ddl:"keyword" sql:"TO" which doesn't emit TO keyword.
	// See tables.go:368 TableConstraintRenameAction for the correct pattern.
	NewName string `ddl:"parameter,no_equals" sql:"TO"`
}

type HybridTableAlterColumnAction struct {
	alterColumn bool   `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName  string `ddl:"keyword"`
	// Manually adjusted: ALTER COLUMN syntax requires "COMMENT 'value'" not "COMMENT = 'value'"
	// See: https://docs.snowflake.com/en/sql-reference/sql/alter-table
	Comment      *string `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
	UnsetComment *bool   `ddl:"keyword" sql:"UNSET COMMENT"`
}

// HybridTableModifyColumnAction is an alias for ALTER COLUMN.
// MODIFY is an alias for ALTER in Snowflake when working with columns.
type HybridTableModifyColumnAction struct {
	modifyColumn bool    `ddl:"static" sql:"MODIFY COLUMN"`
	ColumnName   string  `ddl:"keyword"`
	Comment      *string `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
	UnsetComment *bool   `ddl:"keyword" sql:"UNSET COMMENT"`
}

type HybridTableDropColumnAction struct {
	dropColumn bool     `ddl:"static" sql:"DROP COLUMN"`
	IfExists   *bool    `ddl:"keyword" sql:"IF EXISTS"`
	Columns    []string `ddl:"keyword"`
}

type HybridTableDropIndexAction struct {
	dropIndex bool   `ddl:"static" sql:"DROP INDEX"`
	IfExists  *bool  `ddl:"keyword" sql:"IF EXISTS"`
	IndexName string `ddl:"keyword"`
}

type hybridTableDetailsRow struct {
	Name string `db:"name"`
	Type string `db:"type"`
	Kind string `db:"kind"`
	// Manually adjusted: Snowflake DESCRIBE TABLE returns column named "null?" with question mark.
	// The generator cannot produce "?" in tag names, so this requires manual adjustment.
	Null                  string         `db:"null?"`
	Default               sql.NullString `db:"default"`
	PrimaryKey            string         `db:"primary key"`
	UniqueKey             string         `db:"unique key"`
	Check                 sql.NullString `db:"check"`
	Expression            sql.NullString `db:"expression"`
	Comment               sql.NullString `db:"comment"`
	PolicyName            sql.NullString `db:"policy name"`
	PrivacyDomain         sql.NullString `db:"privacy domain"`
	SchemaEvolutionRecord sql.NullString `db:"schema_evolution_record"`
}

type HybridTableDetails struct {
	Name                  string
	Type                  string
	Kind                  string
	IsNullable            string
	Default               string
	PrimaryKey            string
	UniqueKey             string
	Check                 string
	Expression            string
	Comment               string
	PolicyName            string
	PrivacyDomain         string
	SchemaEvolutionRecord string
}

// --- Standalone CREATE INDEX / DROP INDEX commands ---
// These are standalone SQL commands, not part of ALTER TABLE.
// https://docs.snowflake.com/en/sql-reference/sql/create-index
// https://docs.snowflake.com/en/sql-reference/sql/drop-index

// CreateHybridTableIndexOptions represents the standalone CREATE INDEX command for hybrid tables.
// Syntax: CREATE [OR REPLACE] INDEX [IF NOT EXISTS] <name> ON <table> (<cols>) [INCLUDE (<cols>)]
type CreateHybridTableIndexOptions struct {
	create         bool                   `ddl:"static" sql:"CREATE"`
	OrReplace      *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	index          bool                   `ddl:"static" sql:"INDEX"`
	IfNotExists    *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name           SchemaObjectIdentifier `ddl:"identifier"`
	on             bool                   `ddl:"static" sql:"ON"`
	TableName      SchemaObjectIdentifier `ddl:"identifier"`
	Columns        []string               `ddl:"keyword,parentheses"`
	IncludeColumns []string               `ddl:"keyword,parentheses" sql:"INCLUDE"`
}

func (opts *CreateHybridTableIndexOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !ValidObjectIdentifier(opts.TableName) {
		errs = append(errs, errInvalidIdentifier("CreateHybridTableIndexOptions", "TableName"))
	}
	if len(opts.Columns) == 0 {
		errs = append(errs, errNotSet("CreateHybridTableIndexOptions", "Columns"))
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateHybridTableIndexOptions", "OrReplace", "IfNotExists"))
	}
	return JoinErrors(errs...)
}

// CreateHybridTableIndexRequest is the user-facing request DTO for CREATE INDEX.
type CreateHybridTableIndexRequest struct {
	OrReplace      *bool
	IfNotExists    *bool
	name           SchemaObjectIdentifier // required — the index name
	TableName      SchemaObjectIdentifier // required — the hybrid table
	Columns        []string               // required
	IncludeColumns []string
}

func NewCreateHybridTableIndexRequest(
	name SchemaObjectIdentifier,
	tableName SchemaObjectIdentifier,
	columns []string,
) *CreateHybridTableIndexRequest {
	return &CreateHybridTableIndexRequest{
		name:      name,
		TableName: tableName,
		Columns:   columns,
	}
}

func (s *CreateHybridTableIndexRequest) WithOrReplace(orReplace bool) *CreateHybridTableIndexRequest {
	s.OrReplace = &orReplace
	return s
}

func (s *CreateHybridTableIndexRequest) WithIfNotExists(ifNotExists bool) *CreateHybridTableIndexRequest {
	s.IfNotExists = &ifNotExists
	return s
}

func (s *CreateHybridTableIndexRequest) WithIncludeColumns(includeColumns []string) *CreateHybridTableIndexRequest {
	s.IncludeColumns = includeColumns
	return s
}

func (r *CreateHybridTableIndexRequest) toOpts() *CreateHybridTableIndexOptions {
	return &CreateHybridTableIndexOptions{
		OrReplace:      r.OrReplace,
		IfNotExists:    r.IfNotExists,
		name:           r.name,
		TableName:      r.TableName,
		Columns:        r.Columns,
		IncludeColumns: r.IncludeColumns,
	}
}

// DropHybridTableIndexOptions represents the standalone DROP INDEX command for hybrid tables.
// Syntax: DROP INDEX [IF EXISTS] <table_name>.<index_name>
// TODO [SNOW-XXXXXX]: Snowflake requires <table_name>.<index_name> as a dotted identifier,
// which doesn't fit the standard SchemaObjectIdentifier model. This workaround constructs
// the identifier manually. Revisit once Snowflake supports qualified index names natively.
type DropHybridTableIndexOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	index    bool                   `ddl:"static" sql:"INDEX"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropHybridTableIndexOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

// DropHybridTableIndexRequest is the user-facing request DTO for DROP INDEX.
type DropHybridTableIndexRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required — use <table_name>.<index_name> format
}

// NewDropHybridTableIndexRequest creates a DROP INDEX request.
// TODO [SNOW-XXXXXX]: The name should be constructed as a SchemaObjectIdentifier where the object name
// is "<table_name>.<index_name>" following Snowflake's dotted notation requirement.
// Revisit once Snowflake supports qualified index names natively.
func NewDropHybridTableIndexRequest(
	name SchemaObjectIdentifier,
) *DropHybridTableIndexRequest {
	return &DropHybridTableIndexRequest{
		name: name,
	}
}

func (s *DropHybridTableIndexRequest) WithIfExists(ifExists bool) *DropHybridTableIndexRequest {
	s.IfExists = &ifExists
	return s
}

func (r *DropHybridTableIndexRequest) toOpts() *DropHybridTableIndexOptions {
	return &DropHybridTableIndexOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
}

// --- Interface methods for standalone index commands ---
// manually added — standalone SQL commands, not part of standard CRUD operations

// CreateIndex creates a secondary index on a hybrid table.
// This is a standalone CREATE INDEX command, not part of ALTER TABLE.
// https://docs.snowflake.com/en/sql-reference/sql/create-index
func (v *hybridTables) CreateIndex(ctx context.Context, request *CreateHybridTableIndexRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

// DropIndex drops a secondary index on a hybrid table.
// This is a standalone DROP INDEX command, not part of ALTER TABLE.
// https://docs.snowflake.com/en/sql-reference/sql/drop-index
func (v *hybridTables) DropIndex(ctx context.Context, request *DropHybridTableIndexRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

// --- SHOW INDEXES command ---
// Standalone SHOW INDEXES command for hybrid tables.
// https://docs.snowflake.com/en/sql-reference/sql/show-indexes

// ShowHybridTableIndexesOptions represents SHOW INDEXES [IN <table>].
type ShowHybridTableIndexesOptions struct {
	show    bool                    `ddl:"static" sql:"SHOW"`
	indexes bool                    `ddl:"static" sql:"INDEXES"`
	In      *ShowHybridTableIndexIn `ddl:"keyword" sql:"IN"`
}

type ShowHybridTableIndexIn struct {
	table bool                    `ddl:"static" sql:"TABLE"`
	Table *SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *ShowHybridTableIndexesOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if opts.In != nil && opts.In.Table == nil {
		errs = append(errs, errNotSet("ShowHybridTableIndexesOptions", "In.Table"))
	}
	return JoinErrors(errs...)
}

// hybridTableIndexRow maps the SHOW INDEXES result columns.
// Based on https://docs.snowflake.com/en/sql-reference/sql/show-indexes
type hybridTableIndexRow struct {
	CreatedOn       time.Time      `db:"created_on"`
	Name            string         `db:"name"`
	IsUnique        string         `db:"is_unique"`
	Columns         string         `db:"columns"`
	IncludedColumns sql.NullString `db:"included_columns"`
	TableName       string         `db:"table"`
	DatabaseName    string         `db:"database_name"`
	SchemaName      string         `db:"schema_name"`
	Owner           sql.NullString `db:"owner"`
	OwnerRoleType   sql.NullString `db:"owner_role_type"`
}

// HybridTableIndex is the public representation of a hybrid table index.
type HybridTableIndex struct {
	CreatedOn       time.Time
	Name            string
	IsUnique        bool
	Columns         string
	IncludedColumns string
	TableName       string
	DatabaseName    string
	SchemaName      string
	Owner           string
	OwnerRoleType   string
}

func (r hybridTableIndexRow) convert() (*HybridTableIndex, error) {
	idx := &HybridTableIndex{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		IsUnique:     r.IsUnique == "Y",
		Columns:      r.Columns,
		TableName:    r.TableName,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.IncludedColumns.Valid {
		idx.IncludedColumns = r.IncludedColumns.String
	}
	if r.Owner.Valid {
		idx.Owner = r.Owner.String
	}
	if r.OwnerRoleType.Valid {
		idx.OwnerRoleType = r.OwnerRoleType.String
	}
	return idx, nil
}

var _ convertibleRow[HybridTableIndex] = new(hybridTableIndexRow)

// ShowHybridTableIndexesRequest is the user-facing request DTO for SHOW INDEXES.
type ShowHybridTableIndexesRequest struct {
	In *ShowHybridTableIndexIn
}

func NewShowHybridTableIndexesRequest() *ShowHybridTableIndexesRequest {
	return &ShowHybridTableIndexesRequest{}
}

func (s *ShowHybridTableIndexesRequest) WithIn(in ShowHybridTableIndexIn) *ShowHybridTableIndexesRequest {
	s.In = &in
	return s
}

func (r *ShowHybridTableIndexesRequest) toOpts() *ShowHybridTableIndexesOptions {
	return &ShowHybridTableIndexesOptions{
		In: r.In,
	}
}

// ShowIndexes lists indexes on hybrid tables.
// https://docs.snowflake.com/en/sql-reference/sql/show-indexes
func (v *hybridTables) ShowIndexes(ctx context.Context, request *ShowHybridTableIndexesRequest) ([]HybridTableIndex, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[hybridTableIndexRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[hybridTableIndexRow, HybridTableIndex](dbRows)
}

// --- convert() implementations for generated row structs ---
// These convert() functions map database row structs to public domain objects.
// They are manually implemented here instead of in generated files to allow
// customization without being overwritten by the generator.

// convert maps hybridTableRow (SHOW HYBRID TABLES result) to HybridTable.
func (r hybridTableRow) convert() (*HybridTable, error) {
	ht := &HybridTable{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.Rows.Valid {
		ht.Rows = int(r.Rows.Int64)
	}
	if r.Bytes.Valid {
		ht.Bytes = int(r.Bytes.Int64)
	}
	if r.Owner.Valid {
		ht.Owner = r.Owner.String
	}
	if r.Comment.Valid {
		ht.Comment = r.Comment.String
	}
	if r.OwnerRoleType.Valid {
		ht.OwnerRoleType = r.OwnerRoleType.String
	}
	return ht, nil
}

// convert maps hybridTableDetailsRow (DESCRIBE TABLE result) to HybridTableDetails.
func (r hybridTableDetailsRow) convert() (*HybridTableDetails, error) {
	details := &HybridTableDetails{
		Name:       r.Name,
		Type:       r.Type,
		Kind:       r.Kind,
		IsNullable: r.Null,
		PrimaryKey: r.PrimaryKey,
		UniqueKey:  r.UniqueKey,
	}
	if r.Default.Valid {
		details.Default = r.Default.String
	}
	if r.Check.Valid {
		details.Check = r.Check.String
	}
	if r.Expression.Valid {
		details.Expression = r.Expression.String
	}
	if r.Comment.Valid {
		details.Comment = r.Comment.String
	}
	if r.PolicyName.Valid {
		details.PolicyName = r.PolicyName.String
	}
	if r.PrivacyDomain.Valid {
		details.PrivacyDomain = r.PrivacyDomain.String
	}
	if r.SchemaEvolutionRecord.Valid {
		details.SchemaEvolutionRecord = r.SchemaEvolutionRecord.String
	}
	return details, nil
}
