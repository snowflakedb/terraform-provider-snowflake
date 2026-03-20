package sdk

import (
	"context"
	"errors"
	"fmt"
)

type SemanticViewTableDetails struct {
	TableNameOrAlias      string
	BaseTableDatabaseName string
	BaseTableSchemaName   string
	BaseTableName         string
	PrimaryKeys           string
	UniqueKeys            string
	Synonyms              string
	Comment               string
}

type SemanticViewRelationshipDetails struct {
	RelationshipAlias   string
	TableNameOrAlias    string
	ForeignKeys         string
	RefTableNameOrAlias string
	RefKeys             string
	ParentEntity        string
}

type SemanticViewDimensionDetails struct {
	DimensionAlias   string
	TableNameOrAlias string
	Expression       string
	DataType         string
	AccessModifier   string
	Synonyms         string
	Comment          string
	ParentEntity     string
}

type SemanticViewFactDetails struct {
	FactAlias        string
	TableNameOrAlias string
	Expression       string
	DataType         string
	AccessModifier   string
	Synonyms         string
	Comment          string
	ParentEntity     string
}

type SemanticViewMetricDetails struct {
	MetricAlias      string
	TableNameOrAlias string
	Expression       string
	DataType         string
	AccessModifier   string
	Synonyms         string
	Comment          string
	ParentEntity     string
}

type SemanticViewDescribeDetails struct {
	Id               SchemaObjectIdentifier
	Tables           []SemanticViewTableDetails
	Relationships    []SemanticViewRelationshipDetails
	Dimensions       []SemanticViewDimensionDetails
	Facts            []SemanticViewFactDetails
	Metrics          []SemanticViewMetricDetails
	Comment          string
	DescribeRowCount int
}

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

func (m *MetricDefinition) CheckIsPrivateIsNotNil() bool {
	return m.isPrivate != nil
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

func (f *FactDefinition) GetIsPrivate() bool {
	if f.isPrivate == nil {
		return false
	}
	return *f.isPrivate
}

func (f *FactDefinition) CheckIsPrivateIsNotNil() bool {
	return f.isPrivate != nil
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

func (v *semanticViews) DescribeSemanticViewDetails(ctx context.Context, id SchemaObjectIdentifier) (*SemanticViewDescribeDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseSemanticViewDescribeOutput(properties, id)
}

func parseSemanticViewDescribeOutput(properties []SemanticViewDetails, id SchemaObjectIdentifier) (*SemanticViewDescribeDetails, error) {
	details := &SemanticViewDescribeDetails{
		Id:               id,
		Tables:           []SemanticViewTableDetails{},
		Relationships:    []SemanticViewRelationshipDetails{},
		Dimensions:       []SemanticViewDimensionDetails{},
		Facts:            []SemanticViewFactDetails{},
		Metrics:          []SemanticViewMetricDetails{},
		DescribeRowCount: 0,
	}
	//tableDetails := make([]SemanticViewTableDetails, 0)
	//relationshipDetails := make([]SemanticViewRelationshipDetails, 0)
	//dimensionDetails := make([]SemanticViewDimensionDetails, 0)
	//factDetails := make([]SemanticViewFactDetails, 0)
	//metricDetails := make([]SemanticViewMetricDetails, 0)
	var errs []error
	for _, prop := range properties {
		details.DescribeRowCount++
		if prop.ObjectKind == nil {
			if prop.Property == "COMMENT" {
				details.Comment = prop.PropertyValue
			} else {
				err := errors.New(fmt.Sprintf("Unknown property in DESCRIBE %s", prop.Property))
				errs = append(errs, err)
			}
			continue
		}

		switch *prop.ObjectKind {
		case "TABLE":
			var currentTable *SemanticViewTableDetails
			for i := range details.Tables {
				if details.Tables[i].TableNameOrAlias == *prop.ObjectName {
					currentTable = &details.Tables[i]
					break
				}
			}
			if currentTable == nil {
				details.Tables = append(details.Tables, SemanticViewTableDetails{
					TableNameOrAlias: *prop.ObjectName,
				})
				currentTable = &details.Tables[len(details.Tables)-1]
			}
			switch prop.Property {
			case "BASE_TABLE_DATABASE_NAME":
				currentTable.BaseTableDatabaseName = prop.PropertyValue
			case "BASE_TABLE_SCHEMA_NAME":
				currentTable.BaseTableSchemaName = prop.PropertyValue
			case "BASE_TABLE_NAME":
				currentTable.BaseTableName = prop.PropertyValue
			case "PRIMARY_KEY":
				currentTable.PrimaryKeys = prop.PropertyValue
			case "UNIQUE_KEY":
				currentTable.UniqueKeys = prop.PropertyValue
			case "SYNONYMS":
				currentTable.Synonyms = prop.PropertyValue
			case "COMMENT":
				currentTable.Comment = prop.PropertyValue
			}
		case "RELATIONSHIP":
			var currentRelationship *SemanticViewRelationshipDetails
			for i := range details.Relationships {
				if details.Relationships[i].RelationshipAlias == *prop.ObjectName {
					currentRelationship = &details.Relationships[i]
					break
				}
			}
			if currentRelationship == nil {
				details.Relationships = append(details.Relationships, SemanticViewRelationshipDetails{
					RelationshipAlias: *prop.ObjectName,
					ParentEntity:      *prop.ParentEntity,
				})
				currentRelationship = &details.Relationships[len(details.Relationships)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentRelationship.TableNameOrAlias = prop.PropertyValue
			case "FOREIGN_KEY":
				currentRelationship.ForeignKeys = prop.PropertyValue
			case "REF_TABLE":
				currentRelationship.RefTableNameOrAlias = prop.PropertyValue
			case "REF_KEY":
				currentRelationship.RefKeys = prop.PropertyValue
			}
		case "DIMENSION":
			var currentDimension *SemanticViewDimensionDetails
			for i := range details.Dimensions {
				if details.Dimensions[i].DimensionAlias == *prop.ObjectName {
					currentDimension = &details.Dimensions[i]
					break
				}
			}
			if currentDimension == nil {
				details.Dimensions = append(details.Dimensions, SemanticViewDimensionDetails{
					DimensionAlias: *prop.ObjectName,
					ParentEntity:   *prop.ParentEntity,
				})
				currentDimension = &details.Dimensions[len(details.Dimensions)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentDimension.TableNameOrAlias = prop.PropertyValue
			case "EXPRESSION":
				currentDimension.Expression = prop.PropertyValue
			case "DATA_TYPE":
				currentDimension.DataType = prop.PropertyValue
			case "ACCESS_MODIFIER":
				currentDimension.AccessModifier = prop.PropertyValue
			case "SYNONYMS":
				currentDimension.Synonyms = prop.PropertyValue
			case "COMMENT":
				currentDimension.Comment = prop.PropertyValue
			}
		case "FACT":
			var currentFact *SemanticViewFactDetails
			for i := range details.Facts {
				if details.Facts[i].FactAlias == *prop.ObjectName {
					currentFact = &details.Facts[i]
					break
				}
			}
			if currentFact == nil {
				details.Facts = append(details.Facts, SemanticViewFactDetails{
					FactAlias:    *prop.ObjectName,
					ParentEntity: *prop.ParentEntity,
				})
				currentFact = &details.Facts[len(details.Facts)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentFact.TableNameOrAlias = prop.PropertyValue
			case "EXPRESSION":
				currentFact.Expression = prop.PropertyValue
			case "DATA_TYPE":
				currentFact.DataType = prop.PropertyValue
			case "ACCESS_MODIFIER":
				currentFact.AccessModifier = prop.PropertyValue
			case "SYNONYMS":
				currentFact.Synonyms = prop.PropertyValue
			case "COMMENT":
				currentFact.Comment = prop.PropertyValue
			}
		case "METRIC":
			var currentMetric *SemanticViewMetricDetails
			for i := range details.Metrics {
				if details.Metrics[i].MetricAlias == *prop.ObjectName {
					currentMetric = &details.Metrics[i]
					break
				}
			}
			if currentMetric == nil {
				details.Metrics = append(details.Metrics, SemanticViewMetricDetails{
					MetricAlias:  *prop.ObjectName,
					ParentEntity: *prop.ParentEntity,
				})
				currentMetric = &details.Metrics[len(details.Metrics)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentMetric.TableNameOrAlias = prop.PropertyValue
			case "EXPRESSION":
				currentMetric.Expression = prop.PropertyValue
			case "DATA_TYPE":
				currentMetric.DataType = prop.PropertyValue
			case "ACCESS_MODIFIER":
				currentMetric.AccessModifier = prop.PropertyValue
			case "SYNONYMS":
				currentMetric.Synonyms = prop.PropertyValue
			case "COMMENT":
				currentMetric.Comment = prop.PropertyValue
			}
		}
	}

	return details, errors.Join(errs...)
}
