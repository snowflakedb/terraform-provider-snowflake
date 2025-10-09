package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ SemanticViews = (*semanticViews)(nil)

type semanticViews struct {
	client *Client
}

func (v *semanticViews) Create(ctx context.Context, request *CreateSemanticViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *semanticViews) Alter(ctx context.Context, request *AlterSemanticViewRequest) error {
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
	return convertRows[semanticViewDetailsRow, SemanticViewDetails](rows)
}

func (v *semanticViews) Show(ctx context.Context, request *ShowSemanticViewRequest) ([]SemanticView, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[semanticViewDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[semanticViewDBRow, SemanticView](dbRows)
}

func (v *semanticViews) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*SemanticView, error) {
	request := NewShowSemanticViewRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithIn(In{Schema: id.SchemaId()})
	semanticViews, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(semanticViews, func(r SemanticView) bool { return r.Name == id.Name() })
}

func (v *semanticViews) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*SemanticView, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateSemanticViewRequest) toOpts() *CreateSemanticViewOptions {
	opts := &CreateSemanticViewOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		// Adjusted manually, removed (invalid conversion, e.g., LogicalTables <- LogicalTablesRequest):
		// logicalTables:             r.logicalTables,
		// semanticViewRelationships: r.semanticViewRelationships,
		// semanticViewFacts:         r.semanticViewFacts,
		// semanticViewDimensions:    r.semanticViewDimensions,
		// semanticViewMetrics:       r.semanticViewMetrics,
		// the mapping is done manually below
		Comment:    r.Comment,
		CopyGrants: r.CopyGrants,
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
				tableNameOrAlias:    &RelationshipTableAlias{RelationshipTableAlias: v.tableNameOrAlias.RelationshipTableAlias},
				refTableNameOrAlias: &RelationshipTableAlias{RelationshipTableAlias: v.refTableNameOrAlias.RelationshipTableAlias},
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
		s := make([]MetricDefinition, len(r.semanticViewMetrics))
		for i, v := range r.semanticViewMetrics {
			s[i] = MetricDefinition{}
			if v.semanticExpression != nil {
				s[i].semanticExpression = &SemanticExpression{
					qualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: v.semanticExpression.qualifiedExpressionName.QualifiedExpressionName},
					sqlExpression:           &SemanticSqlExpression{SqlExpression: v.semanticExpression.sqlExpression.SqlExpression},
					Comment:                 v.semanticExpression.Comment,
				}
				if v.semanticExpression.synonyms != nil {
					s[i].semanticExpression.synonyms = &Synonyms{
						WithSynonyms: v.semanticExpression.synonyms.WithSynonyms,
					}
				}
			}
			if v.windowFunctionMetricDefinition != nil {
				s[i].windowFunctionMetricDefinition = &WindowFunctionMetricDefinition{
					WindowFunction: v.windowFunctionMetricDefinition.WindowFunction,
					Metric:         v.windowFunctionMetricDefinition.Metric,
				}
				if v.windowFunctionMetricDefinition.OverClause != nil {
					s[i].windowFunctionMetricDefinition.OverClause = &WindowFunctionOverClause{}
					if v.windowFunctionMetricDefinition.OverClause.PartitionBy != nil {
						s[i].windowFunctionMetricDefinition.OverClause.PartitionBy = v.windowFunctionMetricDefinition.OverClause.PartitionBy
					}
					if v.windowFunctionMetricDefinition.OverClause.OrderBy != nil {
						s[i].windowFunctionMetricDefinition.OverClause.OrderBy = v.windowFunctionMetricDefinition.OverClause.OrderBy
					}
					if v.windowFunctionMetricDefinition.OverClause.WindowFrameClause != nil {
						s[i].windowFunctionMetricDefinition.OverClause.WindowFrameClause = v.windowFunctionMetricDefinition.OverClause.WindowFrameClause
					}
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

func (r semanticViewDetailsRow) convert() (*SemanticViewDetails, error) {
	semanticViewDescribe := &SemanticViewDetails{
		Property:      r.Property,
		PropertyValue: r.PropertyValue,
	}
	if r.ObjectKind.Valid {
		semanticViewDescribe.ObjectKind = String(r.ObjectKind.String)
	}
	if r.ObjectName.Valid {
		semanticViewDescribe.ObjectName = String(r.ObjectName.String)
	}
	if r.ParentEntity.Valid {
		semanticViewDescribe.ParentEntity = String(r.ParentEntity.String)
	}

	return semanticViewDescribe, nil
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

func (r semanticViewDBRow) convert() (*SemanticView, error) {
	semanticViewShow := &SemanticView{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Owner:         r.Owner,
		OwnerRoleType: r.OwnerRoleType,
	}

	if r.Comment.Valid {
		semanticViewShow.Comment = String(r.Comment.String)
	}

	if r.Extension.Valid {
		semanticViewShow.Extension = String(r.Extension.String)
	}

	return semanticViewShow, nil
}

func (r *AlterSemanticViewRequest) toOpts() *AlterSemanticViewOptions {
	opts := &AlterSemanticViewOptions{
		IfExists:     r.IfExists,
		name:         r.name,
		SetComment:   r.SetComment,
		UnsetComment: r.UnsetComment,
	}
	return opts
}
