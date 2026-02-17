package sdk

import (
	"context"
	"database/sql"
	"time"
)

// --- Manually-defined types for CREATE HYBRID TABLE body (rule 13) ---
// The column, constraint, and index structures in a CREATE HYBRID TABLE statement
// are too complex for the generator DSL. These types are referenced by the generated
// CreateHybridTableOptions via the PredefinedQueryStructField "ColumnsAndConstraints".

// HybridTableColumnsConstraintsAndIndexes is the parenthesized body of CREATE HYBRID TABLE,
// containing column definitions, out-of-line constraints, and out-of-line indexes.
type HybridTableColumnsConstraintsAndIndexes struct {
	Columns             []HybridTableColumn              `ddl:"keyword"`
	OutOfLineConstraint []HybridTableOutOfLineConstraint `ddl:"keyword"`
	OutOfLineIndex      []HybridTableOutOfLineIndex      `ddl:"keyword"`
}

// HybridTableColumn defines a single column in a hybrid table.
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table
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
type HybridTableColumnInlineConstraint struct {
	Name       *string              `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	ForeignKey *InlineForeignKey    `ddl:"keyword" sql:"FOREIGN KEY"`
}

// HybridTableOutOfLineConstraint defines out-of-line PRIMARY KEY, UNIQUE, or FOREIGN KEY.
// For hybrid tables, NOT ENFORCED is invalid — all constraints are enforced.
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table
type HybridTableOutOfLineConstraint struct {
	Name       *string              `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	Columns    []string             `ddl:"keyword,parentheses"`
	ForeignKey *OutOfLineForeignKey `ddl:"keyword"`
}

// HybridTableOutOfLineIndex defines a secondary index in CREATE HYBRID TABLE.
// Syntax: INDEX <name> (<cols>) [INCLUDE (<cols>)]
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table
type HybridTableOutOfLineIndex struct {
	index          bool     `ddl:"static" sql:"INDEX"`
	Name           string   `ddl:"keyword"`
	Columns        []string `ddl:"keyword,parentheses"`
	IncludeColumns []string `ddl:"keyword,parentheses" sql:"INCLUDE"`
}

// --- Standalone CREATE INDEX / DROP INDEX commands (rule 13) ---
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
// Note: Snowflake requires <table_name>.<index_name> as a dotted identifier.
// We model this by constructing the SQL manually via structToSQL.
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
// The name should be constructed as a SchemaObjectIdentifier where the object name
// is "<table_name>.<index_name>" following Snowflake's dotted notation requirement.
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

// --- SHOW INDEXES command (rule 13) ---
// Standalone SHOW INDEXES command for hybrid tables.
// https://docs.snowflake.com/en/sql-reference/sql/show-indexes
// Note: This is a hybrid-table-specific command; output columns validated against
// official Snowflake documentation for SHOW INDEXES.

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
	return nil
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

// --- Rule 9 documentation: Discrepancies between Snowflake docs and actual behavior ---
//
// 1. DROP HYBRID TABLE vs DROP TABLE:
//    Codebase insights reference "DROP HYBRID TABLE" syntax, but official Snowflake docs
//    (https://docs.snowflake.com/en/sql-reference/sql/drop-table) use "DROP TABLE".
//    The SDK implements "DROP TABLE" per official docs. This needs validation via
//    integration tests to confirm which syntax is accepted.
//
// 2. CREATE INDEX and UNIQUE:
//    Codebase insights say "Supports UNIQUE indexes", but official Snowflake docs
//    (https://docs.snowflake.com/en/sql-reference/sql/create-index) explicitly state
//    "The CREATE INDEX command cannot be used to add a foreign, primary, or unique key
//    constraint." UNIQUE indexes are created via UNIQUE constraints (ALTER TABLE ADD UNIQUE
//    or inline UNIQUE), not via CREATE INDEX. The SDK follows official docs.
//
// 3. ALTER TABLE ... DROP INDEX / BUILD INDEX:
//    Codebase insights show these as ALTER TABLE sub-commands. Official docs do not
//    document these. They are included in the SDK based on codebase insights and will
//    be validated via integration tests.
