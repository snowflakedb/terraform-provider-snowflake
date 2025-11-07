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

func (l *LogicalTable) SetLogicalTableAlias(alias string) {
	l.logicalTableAlias = &LogicalTableAlias{LogicalTableAlias: alias}
}

func (l *LogicalTable) SetPrimaryKeys(keys []SemanticViewColumn) {
	l.primaryKeys = &PrimaryKeys{
		PrimaryKey: keys,
	}
}

func (l *LogicalTable) SetUniqueKeys(keys [][]SemanticViewColumn) {
	for _, key := range keys {
		l.uniqueKeys = append(l.uniqueKeys, UniqueKeys{Unique: key})
	}
}

func (l *LogicalTable) SetSynonyms(synonyms []Synonym) {
	l.synonyms = &Synonyms{
		WithSynonyms: synonyms,
	}
}

func (m *MetricDefinition) GetSemanticExpression() *SemanticExpression {
	return m.semanticExpression
}

func (m *MetricDefinition) SetSemanticExpression(semExp *SemanticExpression) {
	m.semanticExpression = semExp
}

func (m *MetricDefinition) GetWindowFunctionMetricDefinition() *WindowFunctionMetricDefinition {
	return m.windowFunctionMetricDefinition
}

func (m *MetricDefinition) SetWindowFunctionMetricDefinition(windowFunc *WindowFunctionMetricDefinition) {
	m.windowFunctionMetricDefinition = windowFunc
}

func (s *SemanticExpression) GetQualifiedExpressionName() *QualifiedExpressionName {
	return s.qualifiedExpressionName
}

func (s *SemanticExpression) SetQualifiedExpressionName(qExName string) {
	s.qualifiedExpressionName = &QualifiedExpressionName{QualifiedExpressionName: qExName}
}

func (s *SemanticExpression) GetSqlExpression() *SemanticSqlExpression {
	return s.sqlExpression
}

func (s *SemanticExpression) SetSqlExpression(sqlExpression string) {
	s.sqlExpression = &SemanticSqlExpression{SqlExpression: sqlExpression}
}

func (s *SemanticExpression) GetSynonyms() *Synonyms {
	return s.synonyms
}

func (s *SemanticExpression) SetSynonyms(synonyms []Synonym) {
	s.synonyms = &Synonyms{
		WithSynonyms: synonyms,
	}
}

func (r *SemanticViewRelationship) GetRelationshipAlias() *RelationshipAlias {
	return r.relationshipAlias
}

func (r *SemanticViewRelationship) SetRelationshipAlias(alias string) {
	r.relationshipAlias = &RelationshipAlias{RelationshipAlias: alias}
}

func (r *SemanticViewRelationship) GetTableNameOrAlias() *RelationshipTableAlias {
	return r.tableNameOrAlias
}

func (r *SemanticViewRelationship) SetTableNameOrAlias(tableNameOrAlias RelationshipTableAlias) {
	r.tableNameOrAlias = &tableNameOrAlias
}

func (r *SemanticViewRelationship) GetRelationshipColumnsNames() []SemanticViewColumn {
	return r.relationshipColumnNames
}

func (r *SemanticViewRelationship) SetRelationshipColumnsNames(keys []SemanticViewColumn) {
	r.relationshipColumnNames = keys
}

func (r *SemanticViewRelationship) GetRefTableNameOrAlias() *RelationshipTableAlias {
	return r.refTableNameOrAlias
}

func (r *SemanticViewRelationship) SetRefTableNameOrAlias(refTableNameOrAlias RelationshipTableAlias) {
	r.refTableNameOrAlias = &refTableNameOrAlias
}

func (r *SemanticViewRelationship) GetRelationshipRefColumnsNames() []SemanticViewColumn {
	return r.relationshipRefColumnNames
}

func (r *SemanticViewRelationship) SetRelationshipRefColumnsNames(keys []SemanticViewColumn) {
	r.relationshipRefColumnNames = keys
}
