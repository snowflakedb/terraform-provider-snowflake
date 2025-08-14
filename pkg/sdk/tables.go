package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// TODO [SNOW-1007542]: add missing features:
// - show columns (https://docs.snowflake.com/en/sql-reference/sql/show-columns)
// - show primary keys (https://docs.snowflake.com/en/sql-reference/sql/show-primary-keys)
// - describe search optimization (https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization)
// - truncate table (https://docs.snowflake.com/en/sql-reference/sql/truncate-table)
// - undrop table (https://docs.snowflake.com/en/sql-reference/sql/undrop-table)
type Tables interface {
	Create(ctx context.Context, req *CreateTableRequest) error
	CreateAsSelect(ctx context.Context, req *CreateTableAsSelectRequest) error
	CreateUsingTemplate(ctx context.Context, req *CreateTableUsingTemplateRequest) error
	CreateLike(ctx context.Context, req *CreateTableLikeRequest) error
	CreateClone(ctx context.Context, req *CreateTableCloneRequest) error
	Alter(ctx context.Context, req *AlterTableRequest) error
	Drop(ctx context.Context, req *DropTableRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, req *ShowTableRequest) ([]Table, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Table, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Table, error)
	DescribeColumns(ctx context.Context, req *DescribeTableColumnsRequest) ([]TableColumnDetails, error)
	DescribeStage(ctx context.Context, req *DescribeTableStageRequest) ([]TableStageDetails, error)
}

// TODO: check if [...] in the docs (like in https://docs.snowflake.com/en/sql-reference/sql/create-table#create-table-using-template) mean that we can reuse all parameters from "normal" createTableOptions
type createTableAsSelectOptions struct {
	create          bool                   `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table           bool                   `ddl:"static" sql:"TABLE"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	Columns         []TableAsSelectColumn  `ddl:"list,parentheses"`
	ClusterBy       []string               `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	CopyGrants      *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy *TableRowAccessPolicy  `ddl:"keyword"`
	Query           string                 `ddl:"parameter,no_equals" sql:"AS"`
}

type TableAsSelectColumn struct {
	Name          string                            `ddl:"keyword"`
	Type          *DataType                         `ddl:"keyword"`
	MaskingPolicy *TableAsSelectColumnMaskingPolicy `ddl:"keyword"`
}

type TableAsSelectColumnMaskingPolicy struct {
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	Name          SchemaObjectIdentifier `ddl:"identifier"`
}

type createTableUsingTemplateOptions struct {
	create     bool                   `ddl:"static" sql:"CREATE"`
	OrReplace  *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table      bool                   `ddl:"static" sql:"TABLE"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Query      []string               `ddl:"parameter,no_equals,parentheses" sql:"USING TEMPLATE"`
}

type createTableLikeOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table       bool                   `ddl:"static" sql:"TABLE"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	like        bool                   `ddl:"static" sql:"LIKE"`
	SourceTable SchemaObjectIdentifier `ddl:"identifier"`
	ClusterBy   []string               `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	CopyGrants  *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

type createTableCloneOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	table       bool                   `ddl:"static" sql:"TABLE"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	clone       bool                   `ddl:"static" sql:"CLONE"`
	SourceTable SchemaObjectIdentifier `ddl:"identifier"`
	ClonePoint  *ClonePoint            `ddl:"keyword"`
	CopyGrants  *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

type ClonePoint struct {
	Moment CloneMoment `ddl:"parameter,no_equals"`
	At     TimeTravel  `ddl:"list,parentheses,no_comma"`
}

type CloneMoment string

const (
	CloneMomentAt     CloneMoment = "AT"
	CloneMomentBefore CloneMoment = "BEFORE"
)

type createTableOptions struct {
	create                     bool                             `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                            `ddl:"keyword" sql:"OR REPLACE"`
	Scope                      *TableScope                      `ddl:"keyword"`
	Kind                       *TableKind                       `ddl:"keyword"`
	table                      bool                             `ddl:"static" sql:"TABLE"`
	IfNotExists                *bool                            `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier           `ddl:"identifier"`
	ColumnsAndConstraints      CreateTableColumnsAndConstraints `ddl:"list,parentheses"`
	ClusterBy                  []string                         `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	EnableSchemaEvolution      *bool                            `ddl:"parameter" sql:"ENABLE_SCHEMA_EVOLUTION"`
	StageFileFormat            *StageFileFormat                 `ddl:"list,parentheses,no_comma" sql:"STAGE_FILE_FORMAT ="`
	StageCopyOptions           *StageCopyOptions                `ddl:"list,parentheses,no_comma" sql:"STAGE_COPY_OPTIONS ="`
	DataRetentionTimeInDays    *int                             `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int                             `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ChangeTracking             *bool                            `ddl:"parameter" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *string                          `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	CopyGrants                 *bool                            `ddl:"keyword" sql:"COPY GRANTS"`
	RowAccessPolicy            *TableRowAccessPolicy            `ddl:"keyword"`
	Tags                       []TagAssociation                 `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string                          `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type CreateTableColumnsAndConstraints struct {
	Columns             []TableColumn         `ddl:"keyword"`
	OutOfLineConstraint []OutOfLineConstraint `ddl:"list,no_parentheses"`
}

type TableScope string

const (
	GlobalTableScope TableScope = "GLOBAL"
	LocalTableScope  TableScope = "LOCAL"
)

type TableKind string

const (
	TemporaryTableKind TableKind = "TEMPORARY"
	VolatileTableKind  TableKind = "VOLATILE"
	TransientTableKind TableKind = "TRANSIENT"
)

type TableColumn struct {
	Name             string                  `ddl:"keyword"`
	Type             DataType                `ddl:"keyword"`
	InlineConstraint *ColumnInlineConstraint `ddl:"keyword"`
	NotNull          *bool                   `ddl:"keyword" sql:"NOT NULL"`
	Collate          *string                 `ddl:"parameter,no_equals,single_quotes" sql:"COLLATE"`
	DefaultValue     *ColumnDefaultValue     `ddl:"keyword"`
	MaskingPolicy    *ColumnMaskingPolicy    `ddl:"keyword"`
	Tags             []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
	Comment          *string                 `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
}

type ColumnDefaultValue struct {
	// One of
	Expression *string         `ddl:"parameter,no_equals" sql:"DEFAULT"`
	Identity   *ColumnIdentity `ddl:"keyword" sql:"IDENTITY"`
}
type ColumnIdentity struct {
	Start     int   `ddl:"parameter,no_quotes,no_equals" sql:"START"`
	Increment int   `ddl:"parameter,no_quotes,no_equals" sql:"INCREMENT"`
	Order     *bool `ddl:"keyword" sql:"ORDER"`
	Noorder   *bool `ddl:"keyword" sql:"NOORDER"`
}

type ColumnMaskingPolicy struct {
	With          *bool                  `ddl:"keyword" sql:"WITH"`
	maskingPolicy bool                   `ddl:"static" sql:"MASKING POLICY"`
	Name          SchemaObjectIdentifier `ddl:"identifier"`
	Using         []string               `ddl:"keyword,parentheses" sql:"USING"`
}

// OutOfLineConstraint is based on https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint#out-of-line-unique-primary-foreign-key.
type OutOfLineConstraint struct {
	Name       *string              `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	Columns    []string             `ddl:"keyword,parentheses"`
	ForeignKey *OutOfLineForeignKey `ddl:"keyword"`

	// optional
	Enforced           *bool `ddl:"keyword" sql:"ENFORCED"`
	NotEnforced        *bool `ddl:"keyword" sql:"NOT ENFORCED"`
	Deferrable         *bool `ddl:"keyword" sql:"DEFERRABLE"`
	NotDeferrable      *bool `ddl:"keyword" sql:"NOT DEFERRABLE"`
	InitiallyDeferred  *bool `ddl:"keyword" sql:"INITIALLY DEFERRED"`
	InitiallyImmediate *bool `ddl:"keyword" sql:"INITIALLY IMMEDIATE"`
	Enable             *bool `ddl:"keyword" sql:"ENABLE"`
	Disable            *bool `ddl:"keyword" sql:"DISABLE"`
	Validate           *bool `ddl:"keyword" sql:"VALIDATE"`
	NoValidate         *bool `ddl:"keyword" sql:"NOVALIDATE"`
	Rely               *bool `ddl:"keyword" sql:"RELY"`
	NoRely             *bool `ddl:"keyword" sql:"NORELY"`
}

type OutOfLineForeignKey struct {
	references  bool                   `ddl:"static" sql:"REFERENCES"`
	TableName   SchemaObjectIdentifier `ddl:"identifier"`
	ColumnNames []string               `ddl:"parameter,no_equals,parentheses"`
	Match       *MatchType             `ddl:"parameter,no_equals" sql:"MATCH"`
	On          *ForeignKeyOnAction    `ddl:"keyword"`
}

type alterTableOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"`
	table    bool                   `ddl:"static" sql:"TABLE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	NewName                   *SchemaObjectIdentifier         `ddl:"identifier" sql:"RENAME TO"`
	SwapWith                  *SchemaObjectIdentifier         `ddl:"identifier" sql:"SWAP WITH"`
	ClusteringAction          *TableClusteringAction          `ddl:"keyword"`
	ColumnAction              *TableColumnAction              `ddl:"keyword"`
	ConstraintAction          *TableConstraintAction          `ddl:"keyword"`
	ExternalTableAction       *TableExternalTableAction       `ddl:"keyword"`
	SearchOptimizationAction  *TableSearchOptimizationAction  `ddl:"keyword"`
	Set                       *TableSet                       `ddl:"keyword" sql:"SET"`
	SetTags                   []TagAssociation                `ddl:"parameter,no_equals" sql:"SET TAG"`
	UnsetTags                 []ObjectIdentifier              `ddl:"keyword" sql:"UNSET TAG"`
	Unset                     *TableUnset                     `ddl:"keyword" sql:"UNSET"`
	AddRowAccessPolicy        *TableAddRowAccessPolicy        `ddl:"keyword"`
	DropRowAccessPolicy       *TableDropRowAccessPolicy       `ddl:"keyword"`
	DropAndAddRowAccessPolicy *TableDropAndAddRowAccessPolicy `ddl:"list,no_parentheses"`
	DropAllAccessRowPolicies  *bool                           `ddl:"keyword" sql:"DROP ALL ROW ACCESS POLICIES"`
}

type TableClusteringAction struct {
	// one of
	ClusterBy            []string                   `ddl:"keyword,parentheses" sql:"CLUSTER BY"`
	Recluster            *TableReclusterAction      `ddl:"keyword" sql:"RECLUSTER"`
	ChangeReclusterState *TableReclusterChangeState `ddl:"keyword"`
	DropClusteringKey    *bool                      `ddl:"keyword" sql:"DROP CLUSTERING KEY"`
}

type TableReclusterAction struct {
	MaxSize   *int    `ddl:"parameter" sql:"MAX_SIZE"`
	Condition *string `ddl:"parameter,no_equals" sql:"WHERE"`
}

type TableReclusterChangeState struct {
	State     *ReclusterState `ddl:"keyword"`
	recluster bool            `ddl:"static" sql:"RECLUSTER"`
}

type ReclusterState string

const (
	ReclusterStateResume  ReclusterState = "RESUME"
	ReclusterStateSuspend ReclusterState = "SUSPEND"
)

type TableColumnAction struct {
	// One of
	Add                *TableColumnAddAction                     `ddl:"keyword" sql:"ADD"`
	Rename             *TableColumnRenameAction                  `ddl:"keyword"`
	Alter              []TableColumnAlterAction                  `ddl:"keyword" sql:"ALTER"`
	SetMaskingPolicy   *TableColumnAlterSetMaskingPolicyAction   `ddl:"keyword"`
	UnsetMaskingPolicy *TableColumnAlterUnsetMaskingPolicyAction `ddl:"keyword"`
	SetTags            *TableColumnAlterSetTagsAction            `ddl:"keyword"`
	UnsetTags          *TableColumnAlterUnsetTagsAction          `ddl:"keyword"`
	DropColumns        *TableColumnAlterDropColumns              `ddl:"keyword"`
}

type TableColumnAddAction struct {
	column           bool                            `ddl:"static" sql:"COLUMN"`
	IfNotExists      *bool                           `ddl:"keyword" sql:"IF NOT EXISTS"`
	Name             string                          `ddl:"keyword"`
	Type             DataType                        `ddl:"keyword"`
	Collate          *string                         `ddl:"parameter,no_equals,single_quotes" sql:"COLLATE"`
	DefaultValue     *ColumnDefaultValue             `ddl:"keyword"`
	InlineConstraint *TableColumnAddInlineConstraint `ddl:"keyword"`
	MaskingPolicy    *ColumnMaskingPolicy            `ddl:"keyword"`
	Tags             []TagAssociation                `ddl:"keyword,parentheses" sql:"TAG"`
	Comment          *string                         `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
}

type TableColumnAddInlineConstraint struct {
	NotNull    *bool                `ddl:"keyword" sql:"NOT NULL"`
	Name       *string              `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	ForeignKey *ColumnAddForeignKey `ddl:"keyword"`
}

type ColumnAddForeignKey struct {
	TableName  string `ddl:"keyword" sql:"REFERENCES"`
	ColumnName string `ddl:"keyword,parentheses"`
}

type TableColumnRenameAction struct {
	OldName string `ddl:"parameter,no_equals" sql:"RENAME COLUMN"`
	NewName string `ddl:"parameter,no_equals" sql:"TO"`
}

type TableColumnAlterAction struct {
	column bool   `ddl:"static" sql:"COLUMN"`
	Name   string `ddl:"keyword"`

	// One of (except Collate)
	DropDefault       *bool         `ddl:"keyword" sql:"DROP DEFAULT"`
	SetDefault        *SequenceName `ddl:"parameter,no_equals" sql:"SET DEFAULT"`
	NotNullConstraint *TableColumnNotNullConstraint
	Type              *DataType `ddl:"parameter,no_equals" sql:"SET DATA TYPE"`
	Collate           *string   `ddl:"parameter,no_equals,single_quotes" sql:"COLLATE"`
	Comment           *string   `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
	UnsetComment      *bool     `ddl:"keyword" sql:"UNSET COMMENT"`
}

type TableColumnAlterSetMaskingPolicyAction struct {
	alter             bool                   `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName        string                 `ddl:"keyword"`
	setMaskingPolicy  bool                   `ddl:"static" sql:"SET MASKING POLICY"`
	MaskingPolicyName SchemaObjectIdentifier `ddl:"identifier"`
	Using             []string               `ddl:"keyword,parentheses" sql:"USING"`
	Force             *bool                  `ddl:"keyword" sql:"FORCE"`
}

type TableColumnAlterUnsetMaskingPolicyAction struct {
	alter            bool   `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName       string `ddl:"keyword"`
	setMaskingPolicy bool   `ddl:"static" sql:"UNSET MASKING POLICY"`
}

type TableColumnAlterSetTagsAction struct {
	alter      bool             `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName string           `ddl:"keyword"`
	set        bool             `ddl:"static" sql:"SET"`
	Tags       []TagAssociation `ddl:"keyword" sql:"TAG"`
}

type TableColumnAlterUnsetTagsAction struct {
	alter      bool               `ddl:"static" sql:"ALTER COLUMN"`
	ColumnName string             `ddl:"keyword"`
	unset      bool               `ddl:"static" sql:"UNSET"`
	Tags       []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

type TableColumnAlterDropColumns struct {
	dropColumn bool     `ddl:"static" sql:"DROP COLUMN"`
	IfExists   *bool    `ddl:"keyword" sql:"IF EXISTS"`
	Columns    []string `ddl:"keyword"`
}

type TableColumnAlterSequenceName interface {
	String() string
}

type SequenceName string

func (sn SequenceName) String() string {
	return fmt.Sprintf("%s.NEXTVAL", string(sn))
}

type TableColumnNotNullConstraint struct {
	Set  *bool `ddl:"keyword" sql:"SET NOT NULL"`
	Drop *bool `ddl:"keyword" sql:"DROP NOT NULL"`
}

type TableConstraintAction struct {
	Add    *OutOfLineConstraint         `ddl:"keyword" sql:"ADD"`
	Rename *TableConstraintRenameAction `ddl:"keyword" sql:"RENAME CONSTRAINT"`
	Alter  *TableConstraintAlterAction  `ddl:"keyword" sql:"ALTER"`
	Drop   *TableConstraintDropAction   `ddl:"keyword" sql:"DROP"`
}

type TableConstraintRenameAction struct {
	OldName string `ddl:"keyword"`
	NewName string `ddl:"parameter,no_equals" sql:"TO"`
}

type TableConstraintAlterAction struct {
	// One of
	ConstraintName *string `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	PrimaryKey     *bool   `ddl:"keyword" sql:"PRIMARY KEY"`
	Unique         *bool   `ddl:"keyword" sql:"UNIQUE"`
	ForeignKey     *bool   `ddl:"keyword" sql:"FOREIGN KEY"`

	// Optional
	Columns     []string `ddl:"keyword,parentheses"`
	Enforced    *bool    `ddl:"keyword" sql:"ENFORCED"`
	NotEnforced *bool    `ddl:"keyword" sql:"NOT ENFORCED"`
	Validate    *bool    `ddl:"keyword" sql:"VALIDATE"`
	NoValidate  *bool    `ddl:"keyword" sql:"NOVALIDATE"`
	Rely        *bool    `ddl:"keyword" sql:"RELY"`
	NoRely      *bool    `ddl:"keyword" sql:"NORELY"`
}

type TableConstraintDropAction struct {
	// One of
	ConstraintName *string `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	PrimaryKey     *bool   `ddl:"keyword" sql:"PRIMARY KEY"`
	Unique         *bool   `ddl:"keyword" sql:"UNIQUE"`
	ForeignKey     *bool   `ddl:"keyword" sql:"FOREIGN KEY"`

	// Optional
	Columns  []string `ddl:"keyword,parentheses"`
	Cascade  *bool    `ddl:"keyword" sql:"CASCADE"`
	Restrict *bool    `ddl:"keyword" sql:"RESTRICT"`
}

type TableUnsetTags struct {
	Tag []ObjectIdentifier `ddl:"keyword"`
}

type TableExternalTableAction struct {
	// One of
	Add    *TableExternalTableColumnAddAction    `ddl:"keyword"`
	Rename *TableExternalTableColumnRenameAction `ddl:"keyword"`
	Drop   *TableExternalTableColumnDropAction   `ddl:"keyword"`
}

type TableExternalTableColumnAddAction struct {
	addColumn   bool     `ddl:"static" sql:"ADD COLUMN"`
	IfNotExists *bool    `ddl:"keyword" sql:"IF NOT EXISTS"`
	Name        string   `ddl:"keyword"`
	Type        DataType `ddl:"keyword"`
	Expression  []string `ddl:"parameter,no_equals,parentheses" sql:"AS"`
	Comment     *string  `ddl:"parameter,no_equals,single_quotes" sql:"COMMENT"`
}

type TableExternalTableColumnRenameAction struct {
	OldName string `ddl:"parameter,no_equals" sql:"RENAME COLUMN"`
	NewName string `ddl:"parameter,no_equals" sql:"TO"`
}

type TableExternalTableColumnDropAction struct {
	columns  bool     `ddl:"static" sql:"DROP COLUMN"`
	IfExists *bool    `ddl:"keyword" sql:"IF EXISTS"`
	Names    []string `ddl:"keyword"`
}

type TableSearchOptimizationAction struct {
	// One of
	Add  *AddSearchOptimization  `ddl:"keyword"`
	Drop *DropSearchOptimization `ddl:"keyword"`
}

type AddSearchOptimization struct {
	addSearchOptimization bool `ddl:"static" sql:"ADD SEARCH OPTIMIZATION"`
	// Optional
	On []string `ddl:"keyword" sql:"ON"`
}

type DropSearchOptimization struct {
	dropSearchOptimization bool `ddl:"static" sql:"DROP SEARCH OPTIMIZATION"`
	// Optional
	On []string `ddl:"keyword" sql:"ON"`
}

type TableSet struct {
	// Optional
	EnableSchemaEvolution      *bool             `ddl:"parameter" sql:"ENABLE_SCHEMA_EVOLUTION"`
	StageFileFormat            *StageFileFormat  `ddl:"list,parentheses" sql:"STAGE_FILE_FORMAT ="`
	StageCopyOptions           *StageCopyOptions `ddl:"list,parentheses" sql:"STAGE_COPY_OPTIONS ="`
	DataRetentionTimeInDays    *int              `ddl:"parameter" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *int              `ddl:"parameter" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ChangeTracking             *bool             `ddl:"parameter" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *string           `ddl:"parameter,single_quotes" sql:"DEFAULT_DDL_COLLATION"`
	Comment                    *string           `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type TableUnset struct {
	DataRetentionTimeInDays    *bool `ddl:"keyword" sql:"DATA_RETENTION_TIME_IN_DAYS"`
	MaxDataExtensionTimeInDays *bool `ddl:"keyword" sql:"MAX_DATA_EXTENSION_TIME_IN_DAYS"`
	ChangeTracking             *bool `ddl:"keyword" sql:"CHANGE_TRACKING"`
	DefaultDDLCollation        *bool `ddl:"keyword" sql:"DEFAULT_DDL_COLLATION"`
	EnableSchemaEvolution      *bool `ddl:"keyword" sql:"ENABLE_SCHEMA_EVOLUTION"`
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
}

type TableAddRowAccessPolicy struct {
	add             bool                   `ddl:"static" sql:"ADD"`
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
	On              []string               `ddl:"keyword,parentheses" sql:"ON"`
}

type TableDropRowAccessPolicy struct {
	drop            bool                   `ddl:"static" sql:"DROP"`
	RowAccessPolicy SchemaObjectIdentifier `ddl:"identifier" sql:"ROW ACCESS POLICY"`
}

type TableDropAndAddRowAccessPolicy struct {
	Drop TableDropRowAccessPolicy `ddl:"keyword"`
	Add  TableAddRowAccessPolicy  `ddl:"keyword"`
}

// dropTableOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-table
type dropTableOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	table    bool                   `ddl:"static" sql:"TABLE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	Cascade  *bool `ddl:"keyword" sql:"CASCADE"`
	Restrict *bool `ddl:"keyword" sql:"RESTRICT"`
}

type showTableOptions struct {
	show       bool        `ddl:"static" sql:"SHOW"`
	Terse      *bool       `ddl:"keyword" sql:"TERSE"`
	tables     bool        `ddl:"static" sql:"TABLES"`
	History    *bool       `ddl:"keyword" sql:"HISTORY"`
	Like       *Like       `ddl:"keyword" sql:"LIKE"`
	In         *ExtendedIn `ddl:"keyword" sql:"IN"`
	StartsWith *string     `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom  *LimitFrom  `ddl:"keyword" sql:"LIMIT"`
}

type tableDBRow struct {
	CreatedOn                  string         `db:"created_on"`
	Name                       string         `db:"name"`
	SchemaName                 string         `db:"schema_name"`
	DatabaseName               string         `db:"database_name"`
	Kind                       string         `db:"kind"`
	Comment                    sql.NullString `db:"comment"`
	ClusterBy                  sql.NullString `db:"cluster_by"`
	Rows                       sql.NullInt64  `db:"rows"`
	Bytes                      sql.NullInt64  `db:"bytes"`
	Owner                      string         `db:"owner"`
	RetentionTime              sql.NullInt64  `db:"retention_time"`
	DroppedOn                  sql.NullString `db:"dropped_on"`
	AutomaticClustering        sql.NullString `db:"automatic_clustering"`
	ChangeTracking             sql.NullString `db:"change_tracking"`
	SearchOptimization         sql.NullString `db:"search_optimization"`
	SearchOptimizationProgress sql.NullString `db:"search_optimization_progress"`
	SearchOptimizationBytes    sql.NullInt64  `db:"search_optimization_bytes"`
	IsExternal                 sql.NullString `db:"is_external"`
	EnableSchemaEvolution      sql.NullString `db:"enable_schema_evolution"`
	OwnerRoleType              sql.NullString `db:"owner_role_type"`
	IsEvent                    sql.NullString `db:"is_event"`
	Budget                     sql.NullString `db:"budget"`
}

type Table struct {
	CreatedOn                  string
	Name                       string
	DatabaseName               string
	SchemaName                 string
	Kind                       string
	Comment                    string
	ClusterBy                  string
	Rows                       int
	Bytes                      *int
	Owner                      string
	RetentionTime              int
	DroppedOn                  *string
	AutomaticClustering        bool
	ChangeTracking             bool
	SearchOptimization         bool
	SearchOptimizationProgress string
	SearchOptimizationBytes    *int
	IsExternal                 bool
	EnableSchemaEvolution      bool
	OwnerRoleType              string
	IsEvent                    bool
	Budget                     *string
}

// GetClusterByKeys converts the SHOW TABLES result for ClusterBy and converts it to list of keys.
func (v *Table) GetClusterByKeys() []string {
	if v.ClusterBy == "" {
		return nil
	}

	statementWithoutLinear := strings.TrimSuffix(strings.Replace(v.ClusterBy, "LINEAR(", "", 1), ")")
	return splitClusterBy(statementWithoutLinear)
}

func splitClusterBy(statementWithoutLinear string) []string {
	keysClean := make([]string, 0)

	var current string
	var open int
	for _, character := range statementWithoutLinear {
		strChar := string(character)
		switch strChar {
		case "(":
			open++
		case ")":
			open--
		case ",":
			if open == 0 {
				keysClean = append(keysClean, strings.TrimSpace(current))
				current = ""
				continue
			}
		}
		current += strChar
	}
	keysClean = append(keysClean, strings.TrimSpace(current))

	return keysClean
}

func (row tableDBRow) convert() (*Table, error) {
	table := Table{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		SchemaName:   row.SchemaName,
		DatabaseName: row.DatabaseName,
		Owner:        row.Owner,
		Kind:         row.Kind,
	}
	if row.Rows.Valid {
		table.Rows = int(row.Rows.Int64)
	}
	if row.Bytes.Valid {
		table.Bytes = Int(int(row.Bytes.Int64))
	}
	if row.RetentionTime.Valid {
		table.RetentionTime = int(row.RetentionTime.Int64)
	}
	if row.DroppedOn.Valid {
		table.DroppedOn = String(row.DroppedOn.String)
	}
	if row.AutomaticClustering.Valid {
		table.AutomaticClustering = row.AutomaticClustering.String == "ON"
	}
	if row.ChangeTracking.Valid {
		table.ChangeTracking = row.ChangeTracking.String == "ON"
	}
	if row.SearchOptimization.Valid {
		table.SearchOptimization = row.SearchOptimization.String == "ON"
	}
	if row.SearchOptimizationProgress.Valid {
		table.SearchOptimizationProgress = row.SearchOptimizationProgress.String
	}
	if row.SearchOptimizationBytes.Valid {
		table.SearchOptimizationBytes = Int(int(row.SearchOptimizationBytes.Int64))
	}
	if row.IsExternal.Valid {
		table.IsExternal = row.IsExternal.String == "Y"
	}
	if row.IsEvent.Valid {
		table.IsEvent = row.IsEvent.String == "Y"
	}
	if row.EnableSchemaEvolution.Valid {
		table.EnableSchemaEvolution = row.EnableSchemaEvolution.String == "Y"
	}
	if row.Comment.Valid {
		table.Comment = row.Comment.String
	}
	if row.ClusterBy.Valid {
		table.ClusterBy = row.ClusterBy.String
	}
	if row.OwnerRoleType.Valid {
		table.OwnerRoleType = row.OwnerRoleType.String
	}
	if row.Budget.Valid {
		table.Budget = String(row.Budget.String)
	}
	return &table, nil
}

func (v *Table) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Table) ObjectType() ObjectType {
	return ObjectTypeTable
}

// describeExternalTableColumnsOptions based on https://docs.snowflake.com/en/sql-reference/sql/desc-table
type describeTableColumnsOptions struct {
	describeTable bool                   `ddl:"static" sql:"DESCRIBE TABLE"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	columnsType   bool                   `ddl:"static" sql:"TYPE = COLUMNS"`
}

type TableColumnDetails struct {
	Name                  string
	Type                  DataType
	Kind                  string
	IsNullable            bool
	Default               *string
	IsPrimary             bool
	IsUnique              bool
	Check                 *bool
	Expression            *string
	Comment               *string
	PolicyName            *string
	Collation             *string
	SchemaEvolutionRecord *string
}

// tableColumnDetailsRow based on https://docs.snowflake.com/en/sql-reference/sql/desc-table
type tableColumnDetailsRow struct {
	Name                  string         `db:"name"`
	Type                  DataType       `db:"type"`
	Kind                  string         `db:"kind"`
	IsNullable            string         `db:"null?"`
	Default               sql.NullString `db:"default"`
	IsPrimary             string         `db:"primary key"`
	IsUnique              string         `db:"unique key"`
	Check                 sql.NullString `db:"check"`
	Expression            sql.NullString `db:"expression"`
	Comment               sql.NullString `db:"comment"`
	PolicyName            sql.NullString `db:"policy name"`
	SchemaEvolutionRecord sql.NullString `db:"schema evolution record"`
}

func (r tableColumnDetailsRow) convert() (*TableColumnDetails, error) {
	type_, collation := r.splitTypeAndCollation()

	details := &TableColumnDetails{
		Name:       r.Name,
		Type:       type_,
		Kind:       r.Kind,
		IsNullable: r.IsNullable == "Y",
		IsPrimary:  r.IsPrimary == "Y",
		IsUnique:   r.IsUnique == "Y",
		Collation:  collation,
	}
	if r.Default.Valid {
		details.Default = String(r.Default.String)
	}
	if r.Check.Valid {
		details.Check = Bool(r.Check.String == "Y")
	}
	if r.Expression.Valid {
		details.Expression = String(r.Expression.String)
	}
	if r.Comment.Valid {
		details.Comment = String(r.Comment.String)
	}
	if r.PolicyName.Valid {
		details.PolicyName = String(r.PolicyName.String)
	}
	if r.SchemaEvolutionRecord.Valid {
		details.SchemaEvolutionRecord = String(r.SchemaEvolutionRecord.String)
	}
	return details, nil
}

func (r tableColumnDetailsRow) splitTypeAndCollation() (DataType, *string) {
	collateRegexp := regexp.MustCompile(`COLLATE +'([a-zA-Z0-9_-]*)'`)
	matches := collateRegexp.FindStringSubmatch(string(r.Type))

	if len(matches) == 2 {
		collation := matches[1]
		type_ := DataType(strings.TrimSpace(collateRegexp.ReplaceAllString(string(r.Type), "")))
		return type_, &collation
	}
	return r.Type, nil
}

type describeTableStageOptions struct {
	describeTable bool                   `ddl:"static" sql:"DESCRIBE TABLE"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	stageType     bool                   `ddl:"static" sql:"TYPE = STAGE"`
}

type TableStageDetails struct {
	ParentProperty  string
	Property        string
	PropertyType    string
	PropertyValue   string
	PropertyDefault string
}

type tableStageDetailsRow struct {
	ParentProperty  string `db:"parent_property"`
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

func (r tableStageDetailsRow) convert() (*TableStageDetails, error) {
	return &TableStageDetails{
		ParentProperty:  r.ParentProperty,
		Property:        r.Property,
		PropertyType:    r.PropertyType,
		PropertyValue:   r.PropertyValue,
		PropertyDefault: r.PropertyDefault,
	}, nil
}
