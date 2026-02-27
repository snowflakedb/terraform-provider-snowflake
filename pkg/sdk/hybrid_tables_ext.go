package sdk

import "context"

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
		TableName:    r.Table,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.IsUnique.Valid {
		v := r.IsUnique.String == "Y"
		idx.IsUnique = &v
	}
	if r.Columns.Valid {
		idx.Columns = &r.Columns.String
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
