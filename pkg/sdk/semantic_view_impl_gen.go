package sdk

import (
	"context"
)

var _ SemanticViews = (*semanticViews)(nil)

type semanticViews struct {
	client *Client
}

func (v *semanticViews) Create(ctx context.Context, request *CreateSemanticViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *semanticViews) Drop(ctx context.Context, request *DropSemanticViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *semanticViews) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropSemanticViewRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *semanticViews) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]SemanticViewDetails, error) {
	opts := &DescribeSemanticViewOptions{
		name: id,
	}
	rows, err := validateAndQuery[semanticViewDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[semanticViewDetailsRow, SemanticViewDetails](rows), nil
}

func (v *semanticViews) Show(ctx context.Context, request *ShowSemanticViewRequest) ([]SemanticView, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[semanticViewDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[semanticViewDBRow, SemanticView](dbRows)
	return resultList, nil
}

func (r *CreateSemanticViewRequest) toOpts() *CreateSemanticViewOptions {
	opts := &CreateSemanticViewOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Comment:     r.Comment,
		CopyGrants:  r.CopyGrants,
	}
	if r.logicalTables != nil {
		s := make([]LogicalTable, len(r.logicalTables))
		for i, v := range r.logicalTables {
			s[i] = LogicalTable{
				TableName: v.TableName,
				Comment:   v.Comment,
			}
			if v.logicalTableAlias != nil {
				s[i].logicalTableAlias = &LogicalTableAlias{
					LogicalTableAlias: v.logicalTableAlias.LogicalTableAlias,
				}
			}
			if v.primaryKeys != nil {
				s[i].primaryKeys = &PrimaryKeys{
					PrimaryKey: v.primaryKeys.PrimaryKey,
				}
			}
			if v.synonyms != nil {
				s[i].synonyms = &Synonyms{
					WithSynonyms: v.synonyms.WithSynonyms,
				}
			}
			if v.uniqueKeys != nil {
				u := make([]UniqueKeys, len(v.uniqueKeys))
				for j, w := range v.uniqueKeys {
					u[j] = UniqueKeys{
						Unique: w.Unique,
					}
				}
				s[i].uniqueKeys = u
			}
		}
		opts.logicalTables = s
	}
	if r.semanticViewRelationships != nil {
		s := make([]SemanticViewRelationship, len(r.semanticViewRelationships))
		for i, v := range r.semanticViewRelationships {
			s[i] = SemanticViewRelationship{
				tableName:    &RelationshipTableAlias{RelationshipTableAlias: v.tableName.RelationshipTableAlias},
				refTableName: &RelationshipTableAlias{RelationshipTableAlias: v.refTableName.RelationshipTableAlias},
			}
			if v.relationshipAlias != nil {
				s[i].relationshipAlias = &RelationshipAlias{
					RelationshipAlias: v.relationshipAlias.RelationshipAlias,
				}
			}
			if v.relationshipColumnNames != nil {
				u := make([]SemanticViewColumn, len(v.relationshipColumnNames))
				for j, w := range v.relationshipColumnNames {
					u[j] = SemanticViewColumn{
						Name: w.Name,
					}
				}
				s[i].relationshipColumnNames = u
			}
			if v.relationshipRefColumnNames != nil {
				u := make([]SemanticViewColumn, len(v.relationshipRefColumnNames))
				for j, w := range v.relationshipRefColumnNames {
					u[j] = SemanticViewColumn{
						Name: w.Name,
					}
				}
				s[i].relationshipRefColumnNames = u
			}
		}
		opts.semanticViewRelationships = s
	}
	if r.semanticViewFacts != nil {
		s := make([]SemanticExpression, len(r.semanticViewFacts))
		for i, v := range r.semanticViewFacts {
			s[i] = SemanticExpression{
				qualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: v.qualifiedExpressionName.QualifiedExpressionName},
				sqlExpression:           &SemanticSqlExpression{SqlExpression: v.sqlExpression.SqlExpression},
				Comment:                 v.Comment,
			}
			if v.synonyms != nil {
				s[i].synonyms = &Synonyms{
					WithSynonyms: v.synonyms.WithSynonyms,
				}
			}
		}
		opts.semanticViewFacts = s
	}
	if r.semanticViewDimensions != nil {
		s := make([]SemanticExpression, len(r.semanticViewDimensions))
		for i, v := range r.semanticViewDimensions {
			s[i] = SemanticExpression{
				qualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: v.qualifiedExpressionName.QualifiedExpressionName},
				sqlExpression:           &SemanticSqlExpression{SqlExpression: v.sqlExpression.SqlExpression},
				Comment:                 v.Comment,
			}
			if v.synonyms != nil {
				s[i].synonyms = &Synonyms{
					WithSynonyms: v.synonyms.WithSynonyms,
				}
			}
		}
		opts.semanticViewDimensions = s
	}
	if r.semanticViewMetrics != nil {
		s := make([]SemanticExpression, len(r.semanticViewMetrics))
		for i, v := range r.semanticViewMetrics {
			s[i] = SemanticExpression{
				qualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: v.qualifiedExpressionName.QualifiedExpressionName},
				sqlExpression:           &SemanticSqlExpression{SqlExpression: v.sqlExpression.SqlExpression},
				Comment:                 v.Comment,
			}
			if v.synonyms != nil {
				s[i].synonyms = &Synonyms{
					WithSynonyms: v.synonyms.WithSynonyms,
				}
			}
		}
		opts.semanticViewMetrics = s
	}
	return opts
}

func (r *DropSemanticViewRequest) toOpts() *DropSemanticViewOptions {
	opts := &DropSemanticViewOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeSemanticViewRequest) toOpts() *DescribeSemanticViewOptions {
	opts := &DescribeSemanticViewOptions{
		name: r.name,
	}
	return opts
}

func (r semanticViewDetailsRow) convert() *SemanticViewDetails {
	// TODO: Mapping
	return &SemanticViewDetails{}
}

func (r *ShowSemanticViewRequest) toOpts() *ShowSemanticViewOptions {
	opts := &ShowSemanticViewOptions{
		Terse:      r.Terse,
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r semanticViewDBRow) convert() *SemanticView {
	// TODO: Mapping
	return &SemanticView{}
}
