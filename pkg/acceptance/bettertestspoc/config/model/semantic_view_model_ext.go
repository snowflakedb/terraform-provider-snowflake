package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SemanticViewModel) WithTables(tables []sdk.LogicalTable) *SemanticViewModel {
	maps := make([]tfconfig.Variable, len(tables))
	for i, v := range tables {
		m := map[string]tfconfig.Variable{
			"table_name": tfconfig.StringVariable(v.TableName.FullyQualifiedName()),
		}
		if v.Comment != nil {
			m["comment"] = tfconfig.StringVariable(*v.Comment)
		}
		logicalTableAlias := v.GetLogicalTableAlias()
		if logicalTableAlias != nil {
			m["table_alias"] = tfconfig.StringVariable(logicalTableAlias.LogicalTableAlias)
		}
		primaryKeys := v.GetPrimaryKeys()
		if primaryKeys != nil {
			keys := make([]tfconfig.Variable, len(primaryKeys.PrimaryKey))
			for j, key := range primaryKeys.PrimaryKey {
				keys[j] = tfconfig.StringVariable(key.Name)
			}
			m["primary_key"] = tfconfig.ListVariable(keys...)
		}
		uniqueKeys := v.GetUniqueKeys()
		if uniqueKeys != nil {
			keys := make([]tfconfig.Variable, len(uniqueKeys))
			for j, key := range uniqueKeys {
				uniKeys := make([]tfconfig.Variable, len(key.Unique))
				for k, uniKey := range key.Unique {
					uniKeys[k] = tfconfig.StringVariable(uniKey.Name)
				}
				keys[j] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"values": tfconfig.ListVariable(uniKeys...),
				})
			}
			m["unique"] = tfconfig.ListVariable(keys...)
		}
		synonyms := v.GetSynonyms()
		if synonyms != nil {
			syns := make([]tfconfig.Variable, len(synonyms.WithSynonyms))
			for j, synonym := range synonyms.WithSynonyms {
				syns[j] = tfconfig.StringVariable(synonym.Synonym)
			}
			m["synonym"] = tfconfig.SetVariable(syns...)
		}
		maps[i] = tfconfig.ObjectVariable(m)
	}
	s.Tables = tfconfig.ListVariable(maps...)
	return s
}

func (s *SemanticViewModel) WithRelationships(relationships []sdk.SemanticViewRelationship) *SemanticViewModel {
	maps := make([]tfconfig.Variable, len(relationships))
	for i, v := range relationships {
		m := map[string]tfconfig.Variable{}
		tableNameOrAlias := v.GetTableNameOrAlias()
		if tableNameOrAlias.RelationshipTableName != nil {
			m["table_name"] = tfconfig.StringVariable(tableNameOrAlias.RelationshipTableName.FullyQualifiedName())
		}
		if tableNameOrAlias.RelationshipTableAlias != nil {
			m["table_alias"] = tfconfig.StringVariable(*tableNameOrAlias.RelationshipTableAlias)
		}

		refTableNameOrAlias := v.GetRefTableNameOrAlias()
		if refTableNameOrAlias.RelationshipTableName != nil {
			m["referenced_table_name"] = tfconfig.StringVariable(refTableNameOrAlias.RelationshipTableName.FullyQualifiedName())
		}
		if refTableNameOrAlias.RelationshipTableAlias != nil {
			m["referenced_table_alias"] = tfconfig.StringVariable(*refTableNameOrAlias.RelationshipTableAlias)
		}

		relColumnNames := v.GetRelationshipColumnsNames()
		if relColumnNames != nil {
			keys := make([]tfconfig.Variable, len(relColumnNames))
			for j, key := range relColumnNames {
				keys[j] = tfconfig.StringVariable(key.Name)
			}
			m["relationship_columns"] = tfconfig.ListVariable(keys...)
		}

		relRefColumnNames := v.GetRelationshipRefColumnsNames()
		if relRefColumnNames != nil {
			keys := make([]tfconfig.Variable, len(relRefColumnNames))
			for j, key := range relRefColumnNames {
				keys[j] = tfconfig.StringVariable(key.Name)
			}
			m["referenced_relationship_columns"] = tfconfig.ListVariable(keys...)
		}

		relationshipAlias := v.GetRelationshipAlias()
		if relationshipAlias != nil {
			m["relationship_alias"] = tfconfig.StringVariable(relationshipAlias.RelationshipAlias)
		}

		maps[i] = tfconfig.ObjectVariable(m)
	}
	s.Relationships = tfconfig.ListVariable(maps...)
	return s
}

func (s *SemanticViewModel) WithMetrics(metrics []sdk.MetricDefinition) *SemanticViewModel {
	maps := make([]tfconfig.Variable, len(metrics))
	for i, v := range metrics {
		m := map[string]tfconfig.Variable{}
		semExp := v.GetSemanticExpression()
		windFunc := v.GetWindowFunctionMetricDefinition()
		if semExp != nil {
			semExpVar := map[string]tfconfig.Variable{}
			if semExp.Comment != nil {
				semExpVar["comment"] = tfconfig.StringVariable(*semExp.Comment)
			}
			qExpName := semExp.GetQualifiedExpressionName()
			if qExpName != nil {
				semExpVar["qualified_expression_name"] = tfconfig.StringVariable(qExpName.QualifiedExpressionName)
			}
			sqlExp := semExp.GetSqlExpression()
			if sqlExp != nil {
				semExpVar["sql_expression"] = tfconfig.StringVariable(sqlExp.SqlExpression)
			}
			synonyms := semExp.GetSynonyms()
			if synonyms != nil {
				syns := make([]tfconfig.Variable, len(synonyms.WithSynonyms))
				for j, synonym := range synonyms.WithSynonyms {
					syns[j] = tfconfig.StringVariable(synonym.Synonym)
				}
				semExpVar["synonym"] = tfconfig.SetVariable(syns...)
			}
			m["semantic_expression"] = tfconfig.ListVariable(tfconfig.ObjectVariable(semExpVar))
		} else if windFunc != nil {
			windFuncVar := map[string]tfconfig.Variable{
				"window_function": tfconfig.StringVariable(windFunc.WindowFunction),
				"metric":          tfconfig.StringVariable(windFunc.Metric),
			}
			if windFunc.OverClause != nil {
				overClauseVar := map[string]tfconfig.Variable{}
				if windFunc.OverClause.PartitionBy != nil {
					overClauseVar["partition_by"] = tfconfig.StringVariable(*windFunc.OverClause.PartitionBy)
				}
				if windFunc.OverClause.OrderBy != nil {
					overClauseVar["order_by"] = tfconfig.StringVariable(*windFunc.OverClause.OrderBy)
				}
				if windFunc.OverClause.WindowFrameClause != nil {
					overClauseVar["window_frame_clause"] = tfconfig.StringVariable(*windFunc.OverClause.WindowFrameClause)
				}
				windFuncVar["over_clause"] = tfconfig.ObjectVariable(overClauseVar)
			}
			m["window_function"] = tfconfig.ListVariable(tfconfig.ObjectVariable(windFuncVar))
		}
		maps[i] = tfconfig.ObjectVariable(m)
	}
	s.Metrics = tfconfig.ListVariable(maps...)
	return s
}

func (s *SemanticViewModel) WithFacts(facts []sdk.SemanticExpression) *SemanticViewModel {
	maps := make([]tfconfig.Variable, len(facts))
	for i, semExp := range facts {
		m := map[string]tfconfig.Variable{}
		if semExp.Comment != nil {
			m["comment"] = tfconfig.StringVariable(*semExp.Comment)
		}
		qExpName := semExp.GetQualifiedExpressionName()
		if qExpName != nil {
			m["qualified_expression_name"] = tfconfig.StringVariable(qExpName.QualifiedExpressionName)
		}
		sqlExp := semExp.GetSqlExpression()
		if sqlExp != nil {
			m["sql_expression"] = tfconfig.StringVariable(sqlExp.SqlExpression)
		}
		synonyms := semExp.GetSynonyms()
		if synonyms != nil {
			syns := make([]tfconfig.Variable, len(synonyms.WithSynonyms))
			for j, synonym := range synonyms.WithSynonyms {
				syns[j] = tfconfig.StringVariable(synonym.Synonym)
			}
			m["synonym"] = tfconfig.SetVariable(syns...)
		}
		maps[i] = tfconfig.ObjectVariable(m)
	}
	s.Facts = tfconfig.ListVariable(maps...)
	return s
}

func (s *SemanticViewModel) WithDimensions(dimensions []sdk.SemanticExpression) *SemanticViewModel {
	maps := make([]tfconfig.Variable, len(dimensions))
	for i, semExp := range dimensions {
		m := map[string]tfconfig.Variable{}
		if semExp.Comment != nil {
			m["comment"] = tfconfig.StringVariable(*semExp.Comment)
		}
		qExpName := semExp.GetQualifiedExpressionName()
		if qExpName != nil {
			m["qualified_expression_name"] = tfconfig.StringVariable(qExpName.QualifiedExpressionName)
		}
		sqlExp := semExp.GetSqlExpression()
		if sqlExp != nil {
			m["sql_expression"] = tfconfig.StringVariable(sqlExp.SqlExpression)
		}
		synonyms := semExp.GetSynonyms()
		if synonyms != nil {
			syns := make([]tfconfig.Variable, len(synonyms.WithSynonyms))
			for j, synonym := range synonyms.WithSynonyms {
				syns[j] = tfconfig.StringVariable(synonym.Synonym)
			}
			m["synonym"] = tfconfig.SetVariable(syns...)
		}
		maps[i] = tfconfig.ObjectVariable(m)
	}
	s.Dimensions = tfconfig.ListVariable(maps...)
	return s
}

func LogicalTableWithProps(
	alias string,
	tableName sdk.SchemaObjectIdentifier,
	primaryKeys []sdk.SemanticViewColumn,
	uniqueKeys [][]sdk.SemanticViewColumn,
	synonyms []sdk.Synonym,
	comment string,
) *sdk.LogicalTable {
	table := &sdk.LogicalTable{
		TableName: tableName,
		Comment:   &comment,
	}
	if alias != "" {
		table.SetLogicalTableAlias(alias)
	}
	if primaryKeys != nil {
		table.SetPrimaryKeys(primaryKeys)
	}
	if uniqueKeys != nil {
		table.SetUniqueKeys(uniqueKeys)
	}
	if synonyms != nil {
		table.SetSynonyms(synonyms)
	}
	return table
}

func SemanticExpressionWithProps(
	qualifiedExpressionName string,
	sqlExpression string,
	synonyms []sdk.Synonym,
	comment string,
) *sdk.SemanticExpression {
	semanticExpression := &sdk.SemanticExpression{
		Comment: &comment,
	}
	if qualifiedExpressionName != "" {
		semanticExpression.SetQualifiedExpressionName(qualifiedExpressionName)
	}
	if sqlExpression != "" {
		semanticExpression.SetSqlExpression(sqlExpression)
	}
	if synonyms != nil {
		semanticExpression.SetSynonyms(synonyms)
	}

	return semanticExpression
}

func WindowFunctionMetricDefinitionWithProps(
	windowFunction string,
	metric string,
	overClause sdk.WindowFunctionOverClause,
) *sdk.WindowFunctionMetricDefinition {
	windowFunctionMetricDefinition := &sdk.WindowFunctionMetricDefinition{
		WindowFunction: windowFunction,
		Metric:         metric,
		OverClause:     &overClause,
	}

	return windowFunctionMetricDefinition
}

func MetricDefinitionWithProps(semExp *sdk.SemanticExpression, windowFunc *sdk.WindowFunctionMetricDefinition) *sdk.MetricDefinition {
	metric := &sdk.MetricDefinition{}
	if semExp != nil {
		metric.SetSemanticExpression(semExp)
	} else if windowFunc != nil {
		metric.SetWindowFunctionMetricDefinition(windowFunc)
	}

	return metric
}

func RelationshipTableAliasWithProps(
	alias string,
	tableName sdk.SchemaObjectIdentifier,
) *sdk.RelationshipTableAlias {
	res := &sdk.RelationshipTableAlias{
		RelationshipTableName:  &tableName,
		RelationshipTableAlias: &alias,
	}

	return res
}

func RelationshipWithProps(
	relationshipAlias string,
	tableNameOrAlias sdk.RelationshipTableAlias,
	relColumnNames []sdk.SemanticViewColumn,
	refTableNameOrAlias sdk.RelationshipTableAlias,
	relRefColumnNames []sdk.SemanticViewColumn,
) *sdk.SemanticViewRelationship {
	rel := &sdk.SemanticViewRelationship{}
	rel.SetTableNameOrAlias(tableNameOrAlias)
	rel.SetRelationshipColumnsNames(relColumnNames)
	rel.SetRefTableNameOrAlias(refTableNameOrAlias)
	if relRefColumnNames != nil {
		rel.SetRelationshipRefColumnsNames(relRefColumnNames)
	}
	if relationshipAlias != "" {
		rel.SetRelationshipAlias(relationshipAlias)
	}

	return rel
}
