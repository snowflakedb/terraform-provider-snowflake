package sdk

import "context"

// TableContact represents a CONTACT <purpose> = <contact_name> assignment.
type TableContact struct {
	Purpose string `ddl:"keyword"`
	Contact string `ddl:"parameter,no_equals,single_quotes"`
}

// HybridTableColumnsConstraintsAndIndexes is the parenthesized body of CREATE HYBRID TABLE,
// containing column definitions, out-of-line constraints, and out-of-line indexes.
type HybridTableColumnsConstraintsAndIndexes struct {
	Columns             []HybridTableColumn              `ddl:"keyword"`
	OutOfLineConstraint []HybridTableOutOfLineConstraint `ddl:"keyword"`
	OutOfLineIndex      []HybridTableOutOfLineIndex      `ddl:"keyword"`
}

// HybridTableColumn defines a single column in a hybrid table.
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
type HybridTableColumnInlineConstraint struct {
	Name       *string              `ddl:"parameter,no_equals" sql:"CONSTRAINT"`
	Type       ColumnConstraintType `ddl:"keyword"`
	ForeignKey *InlineForeignKey    `ddl:"keyword" sql:"REFERENCES"`
}

// HybridTableOutOfLineConstraint reuses OutOfLineConstraint from tables.go.
type HybridTableOutOfLineConstraint = OutOfLineConstraint

// HybridTableOutOfLineIndex defines a secondary index in CREATE HYBRID TABLE.
type HybridTableOutOfLineIndex struct {
	index          bool     `ddl:"static" sql:"INDEX"`
	Name           string   `ddl:"keyword"`
	Columns        []string `ddl:"keyword,parentheses"`
	IncludeColumns []string `ddl:"keyword,parentheses" sql:"INCLUDE"`
}

func (r hybridTableRow) convert() (*HybridTable, error) {
	ht := &HybridTable{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.Rows.Valid {
		v := int(r.Rows.Int64)
		ht.Rows = &v
	}
	if r.Bytes.Valid {
		v := int(r.Bytes.Int64)
		ht.Bytes = &v
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

func (r hybridTableIndexRow) convert() (*HybridTableIndex, error) {
	idx := &HybridTableIndex{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		IsUnique:     r.IsUnique == "Y",
		Columns:      r.Columns,
		TableName:    r.Table,
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

// Standalone index operations

func (v *hybridTables) CreateIndex(ctx context.Context, request *CreateIndexHybridTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *hybridTables) DropIndex(ctx context.Context, request *DropIndexHybridTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *hybridTables) ShowIndexes(ctx context.Context, request *ShowIndexesHybridTableRequest) ([]HybridTableIndex, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[hybridTableIndexRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[hybridTableIndexRow, HybridTableIndex](dbRows)
}
