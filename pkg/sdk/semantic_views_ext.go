package sdk

func (s *CreateSemanticViewRequest) GetName() SchemaObjectIdentifier {
	return s.name
}

func (l *LogicalTable) GetLogicalTableAlias() *LogicalTableAlias {
	return l.logicalTableAlias
}

func (l *LogicalTable) GetPrimaryKeys() *PrimaryKeys {
	return l.primaryKeys
}

func (l *LogicalTable) GetUniqueKeys() []UniqueKeys {
	return l.uniqueKeys
}

func (l *LogicalTable) GetSynonyms() *Synonyms {
	return l.synonyms
}

func (l *LogicalTable) WithLogicalTableAlias(alias string) *LogicalTable {
	l.logicalTableAlias = &LogicalTableAlias{LogicalTableAlias: alias}

	return l
}

func (l *LogicalTable) WithTableName(tableName SchemaObjectIdentifier) *LogicalTable {
	l.TableName = tableName

	return l
}

func (l *LogicalTable) WithPrimaryKeys(keys []SemanticViewColumn) *LogicalTable {
	l.primaryKeys = &PrimaryKeys{
		PrimaryKey: keys,
	}

	return l
}

func (l *LogicalTable) WithUniqueKeys(keys [][]SemanticViewColumn) *LogicalTable {
	for _, key := range keys {
		l.uniqueKeys = append(l.uniqueKeys, UniqueKeys{Unique: key})
	}
	return l
}

func (l *LogicalTable) WithSynonyms(synonyms []Synonym) *LogicalTable {
	l.synonyms = &Synonyms{
		WithSynonyms: synonyms,
	}

	return l
}

func (l *LogicalTable) WithComment(comment string) *LogicalTable {
	l.Comment = &comment

	return l
}

func (m *MetricDefinition) GetIsPrivate() bool {
	if m.isPrivate == nil {
		return false
	}
	return *m.isPrivate
}

func (m *MetricDefinition) GetSemanticExpression() *SemanticExpression {
	return m.semanticExpression
}

func (m *MetricDefinition) GetWindowFunctionMetricDefinition() *WindowFunctionMetricDefinition {
	return m.windowFunctionMetricDefinition
}

func (m *MetricDefinition) WithIsPrivate(isPrivate bool) *MetricDefinition {
	m.isPrivate = &isPrivate

	return m
}

func (m *MetricDefinition) WithSemanticExpression(semExp *SemanticExpression) *MetricDefinition {
	m.semanticExpression = semExp

	return m
}

func (m *MetricDefinition) WithWindowFunctionMetricDefinition(windowFunc *WindowFunctionMetricDefinition) *MetricDefinition {
	m.windowFunctionMetricDefinition = windowFunc

	return m
}

func (f *FactDefinition) GetIsPrivate() *bool {
	return f.isPrivate
}

func (f *FactDefinition) WithIsPrivate(isPrivate bool) *FactDefinition {
	f.isPrivate = &isPrivate

	return f
}

func (f *FactDefinition) GetSemanticExpression() *SemanticExpression {
	return f.semanticExpression
}

func (f *FactDefinition) WithSemanticExpression(semExp *SemanticExpression) *FactDefinition {
	f.semanticExpression = semExp

	return f
}

func (d *DimensionDefinition) GetSemanticExpression() *SemanticExpression {
	return d.semanticExpression
}

func (d *DimensionDefinition) WithSemanticExpression(semExp *SemanticExpression) *DimensionDefinition {
	d.semanticExpression = semExp

	return d
}

func (s *SemanticExpression) GetQualifiedExpressionName() *QualifiedExpressionName {
	return s.qualifiedExpressionName
}

func (s *SemanticExpression) WithQualifiedExpressionName(qExName string) *SemanticExpression {
	s.qualifiedExpressionName = &QualifiedExpressionName{QualifiedExpressionName: qExName}

	return s
}

func (s *SemanticExpression) WithSqlExpression(sqlExpression string) *SemanticExpression {
	s.sqlExpression = &SemanticSqlExpression{SqlExpression: sqlExpression}

	return s
}

func (s *SemanticExpression) WithSynonyms(synonyms []Synonym) *SemanticExpression {
	s.synonyms = &Synonyms{
		WithSynonyms: synonyms,
	}

	return s
}

func (s *SemanticExpression) WithComment(comment string) *SemanticExpression {
	s.Comment = &comment

	return s
}

func (s *SemanticExpression) GetSqlExpression() *SemanticSqlExpression {
	return s.sqlExpression
}

func (s *SemanticExpression) GetSynonyms() *Synonyms {
	return s.synonyms
}

func (w *WindowFunctionMetricDefinition) GetQualifiedExpressionName() *QualifiedExpressionName {
	return w.qualifiedExpressionName
}

func (w *WindowFunctionMetricDefinition) WithQualifiedExpressionName(qExName string) *WindowFunctionMetricDefinition {
	w.qualifiedExpressionName = &QualifiedExpressionName{QualifiedExpressionName: qExName}

	return w
}

func (w *WindowFunctionMetricDefinition) GetSqlExpression() *SemanticSqlExpression {
	return w.sqlExpression
}

func (w *WindowFunctionMetricDefinition) WithSqlExpression(sqlExpression string) *WindowFunctionMetricDefinition {
	w.sqlExpression = &SemanticSqlExpression{SqlExpression: sqlExpression}

	return w
}

func (w *WindowFunctionMetricDefinition) WithOverClause(overClause WindowFunctionOverClause) *WindowFunctionMetricDefinition {
	w.OverClause = &overClause

	return w
}

func (r *SemanticViewRelationship) GetRelationshipAlias() *RelationshipAlias {
	return r.relationshipAlias
}

func (r *SemanticViewRelationship) WithRelationshipAlias(alias string) *SemanticViewRelationship {
	r.relationshipAlias = &RelationshipAlias{RelationshipAlias: alias}

	return r
}

func (r *SemanticViewRelationship) GetTableNameOrAlias() *RelationshipTableAlias {
	return r.tableNameOrAlias
}

func (r *SemanticViewRelationship) WithTableNameOrAlias(tableNameOrAlias RelationshipTableAlias) *SemanticViewRelationship {
	r.tableNameOrAlias = &tableNameOrAlias

	return r
}

func (r *SemanticViewRelationship) GetRelationshipColumnsNames() []SemanticViewColumn {
	return r.relationshipColumnNames
}

func (r *SemanticViewRelationship) WithRelationshipColumnsNames(keys []SemanticViewColumn) *SemanticViewRelationship {
	r.relationshipColumnNames = keys

	return r
}

func (r *SemanticViewRelationship) GetRefTableNameOrAlias() *RelationshipTableAlias {
	return r.refTableNameOrAlias
}

func (r *SemanticViewRelationship) WithRefTableNameOrAlias(refTableNameOrAlias RelationshipTableAlias) *SemanticViewRelationship {
	r.refTableNameOrAlias = &refTableNameOrAlias

	return r
}

func (r *SemanticViewRelationship) GetRelationshipRefColumnsNames() []SemanticViewColumn {
	return r.relationshipRefColumnNames
}

func (r *SemanticViewRelationship) WithRelationshipRefColumnsNames(keys []SemanticViewColumn) *SemanticViewRelationship {
	r.relationshipRefColumnNames = keys

	return r
}

func (r *RelationshipTableAlias) WithRelationshipTableAlias(alias string) *RelationshipTableAlias {
	r.RelationshipTableAlias = &alias

	return r
}

func (r *RelationshipTableAlias) WithRelationshipTableName(tableName SchemaObjectIdentifier) *RelationshipTableAlias {
	r.RelationshipTableName = &tableName

	return r
}
