package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type SemanticViewTableDetails struct {
	TableNameOrAlias string
	BaseTable        SchemaObjectIdentifier
	PrimaryKeys      []string
	UniqueKeys       [][]string
	Synonyms         []string
	Comment          string
}

type SemanticViewRelationshipDetails struct {
	RelationshipAlias   string
	TableNameOrAlias    string
	ForeignKeys         []string
	RefTableNameOrAlias string
	RefKeys             []string
	ParentEntity        string
}

type CommonProperties struct {
	TableNameOrAlias string
	Expression       string
	DataType         string
	AccessModifier   string
}

type SemanticViewDimensionDetails struct {
	DimensionAlias string
	Properties     CommonProperties
	Synonyms       []string
	Comment        string
	ParentEntity   string
}

type SemanticViewFactDetails struct {
	FactAlias    string
	Properties   CommonProperties
	Synonyms     []string
	Comment      string
	ParentEntity string
}

type SemanticViewMetricDetails struct {
	MetricAlias  string
	Properties   CommonProperties
	Synonyms     []string
	Comment      string
	ParentEntity string
}

func (d *SemanticViewDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

type SemanticViewDetails struct {
	Id               SchemaObjectIdentifier
	Tables           []SemanticViewTableDetails
	Relationships    []SemanticViewRelationshipDetails
	Dimensions       []SemanticViewDimensionDetails
	Facts            []SemanticViewFactDetails
	Metrics          []SemanticViewMetricDetails
	Comment          string
	DescribeRowCount int
}

func (opts *CreateSemanticViewOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.SemanticViewRelationships) {
		for _, v := range opts.SemanticViewRelationships {
			if !exactlyOneValueSet(v.TableNameOrAlias.RelationshipTableName, v.TableNameOrAlias.RelationshipTableAlias) {
				errs = append(errs, errExactlyOneOf("CreateSemanticViewOptions.semanticViewRelationships.tableNameOrAlias", "RelationshipTableName", "RelationshipTableAlias"))
			}
			if !exactlyOneValueSet(v.RefTableNameOrAlias.RelationshipTableName, v.RefTableNameOrAlias.RelationshipTableAlias) {
				errs = append(errs, errExactlyOneOf("CreateSemanticViewOptions.semanticViewRelationships.refTableNameOrAlias", "RelationshipTableName", "RelationshipTableAlias"))
			}
		}
	}
	return JoinErrors(errs...)
}

func (s *CreateSemanticViewRequest) GetName() SchemaObjectIdentifier {
	return s.name
}

func (l *LogicalTable) GetLogicalTableAlias() *LogicalTableAlias {
	return l.LogicalTableAlias
}

func (l *LogicalTable) GetPrimaryKeys() *PrimaryKeys {
	return l.PrimaryKeys
}

func (l *LogicalTable) GetUniqueKeys() []UniqueKeys {
	return l.UniqueKeys
}

func (l *LogicalTable) GetSynonyms() *Synonyms {
	return l.Synonyms
}

func (l *LogicalTable) WithLogicalTableAlias(alias string) *LogicalTable {
	l.LogicalTableAlias = &LogicalTableAlias{LogicalTableAlias: alias}

	return l
}

func (l *LogicalTable) WithTableName(tableName SchemaObjectIdentifier) *LogicalTable {
	l.TableName = tableName

	return l
}

func (l *LogicalTable) WithPrimaryKeys(keys []SemanticViewColumn) *LogicalTable {
	l.PrimaryKeys = &PrimaryKeys{
		PrimaryKey: keys,
	}

	return l
}

func (l *LogicalTable) WithUniqueKeys(keys [][]SemanticViewColumn) *LogicalTable {
	for _, key := range keys {
		l.UniqueKeys = append(l.UniqueKeys, UniqueKeys{Unique: key})
	}
	return l
}

func (l *LogicalTable) WithSynonyms(synonyms []Synonym) *LogicalTable {
	l.Synonyms = &Synonyms{
		WithSynonyms: synonyms,
	}

	return l
}

func (l *LogicalTable) WithComment(comment string) *LogicalTable {
	l.Comment = &comment

	return l
}

func (m *MetricDefinition) GetIsPrivate() bool {
	if m.IsPrivate == nil {
		return false
	}
	return *m.IsPrivate
}

func (m *MetricDefinition) CheckIsPrivateIsNotNil() bool {
	return m.IsPrivate != nil
}

func (m *MetricDefinition) GetSemanticExpression() *SemanticExpression {
	return m.SemanticExpression
}

func (m *MetricDefinition) GetWindowFunctionMetricDefinition() *WindowFunctionMetricDefinition {
	return m.WindowFunctionMetricDefinition
}

func (m *MetricDefinition) WithIsPrivate(isPrivate bool) *MetricDefinition {
	m.IsPrivate = &isPrivate

	return m
}

func (m *MetricDefinition) WithSemanticExpression(semExp *SemanticExpression) *MetricDefinition {
	m.SemanticExpression = semExp

	return m
}

func (m *MetricDefinition) WithWindowFunctionMetricDefinition(windowFunc *WindowFunctionMetricDefinition) *MetricDefinition {
	m.WindowFunctionMetricDefinition = windowFunc

	return m
}

func (f *FactDefinition) GetIsPrivate() bool {
	if f.IsPrivate == nil {
		return false
	}
	return *f.IsPrivate
}

func (f *FactDefinition) CheckIsPrivateIsNotNil() bool {
	return f.IsPrivate != nil
}

func (f *FactDefinition) WithIsPrivate(isPrivate bool) *FactDefinition {
	f.IsPrivate = &isPrivate

	return f
}

func (f *FactDefinition) GetSemanticExpression() *SemanticExpression {
	return f.SemanticExpression
}

func (f *FactDefinition) WithSemanticExpression(semExp *SemanticExpression) *FactDefinition {
	f.SemanticExpression = semExp

	return f
}

func (d *DimensionDefinition) GetSemanticExpression() *SemanticExpression {
	return d.SemanticExpression
}

func (d *DimensionDefinition) WithSemanticExpression(semExp *SemanticExpression) *DimensionDefinition {
	d.SemanticExpression = semExp

	return d
}

func (s *SemanticExpression) GetQualifiedExpressionName() *QualifiedExpressionName {
	return s.QualifiedExpressionName
}

func (s *SemanticExpression) WithQualifiedExpressionName(qExName string) *SemanticExpression {
	s.QualifiedExpressionName = &QualifiedExpressionName{QualifiedExpressionName: qExName}

	return s
}

func (s *SemanticExpression) WithSqlExpression(sqlExpression string) *SemanticExpression {
	s.SqlExpression = &SemanticSqlExpression{SqlExpression: sqlExpression}

	return s
}

func (s *SemanticExpression) WithSynonyms(synonyms []Synonym) *SemanticExpression {
	s.Synonyms = &Synonyms{
		WithSynonyms: synonyms,
	}

	return s
}

func (s *SemanticExpression) WithComment(comment string) *SemanticExpression {
	s.Comment = &comment

	return s
}

func (s *SemanticExpression) GetSqlExpression() *SemanticSqlExpression {
	return s.SqlExpression
}

func (s *SemanticExpression) GetSynonyms() *Synonyms {
	return s.Synonyms
}

func (w *WindowFunctionMetricDefinition) GetQualifiedExpressionName() *QualifiedExpressionName {
	return w.QualifiedExpressionName
}

func (w *WindowFunctionMetricDefinition) WithQualifiedExpressionName(qExName string) *WindowFunctionMetricDefinition {
	w.QualifiedExpressionName = &QualifiedExpressionName{QualifiedExpressionName: qExName}

	return w
}

func (w *WindowFunctionMetricDefinition) GetSqlExpression() *SemanticSqlExpression {
	return w.SqlExpression
}

func (w *WindowFunctionMetricDefinition) WithSqlExpression(sqlExpression string) *WindowFunctionMetricDefinition {
	w.SqlExpression = &SemanticSqlExpression{SqlExpression: sqlExpression}

	return w
}

func (w *WindowFunctionMetricDefinition) WithOverClause(overClause WindowFunctionOverClause) *WindowFunctionMetricDefinition {
	w.OverClause = &overClause

	return w
}

func (r *SemanticViewRelationship) GetRelationshipAlias() *RelationshipAlias {
	return r.RelationshipAlias
}

func (r *SemanticViewRelationship) WithRelationshipAlias(alias string) *SemanticViewRelationship {
	r.RelationshipAlias = &RelationshipAlias{RelationshipAlias: alias}

	return r
}

func (r *SemanticViewRelationship) GetTableNameOrAlias() *RelationshipTableAlias {
	return r.TableNameOrAlias
}

func (r *SemanticViewRelationship) WithTableNameOrAlias(tableNameOrAlias RelationshipTableAlias) *SemanticViewRelationship {
	r.TableNameOrAlias = &tableNameOrAlias

	return r
}

func (r *SemanticViewRelationship) GetRelationshipColumnsNames() []SemanticViewColumn {
	return r.RelationshipColumnNames
}

func (r *SemanticViewRelationship) WithRelationshipColumnsNames(keys []SemanticViewColumn) *SemanticViewRelationship {
	r.RelationshipColumnNames = keys

	return r
}

func (r *SemanticViewRelationship) GetRefTableNameOrAlias() *RelationshipTableAlias {
	return r.RefTableNameOrAlias
}

func (r *SemanticViewRelationship) WithRefTableNameOrAlias(refTableNameOrAlias RelationshipTableAlias) *SemanticViewRelationship {
	r.RefTableNameOrAlias = &refTableNameOrAlias

	return r
}

func (r *SemanticViewRelationship) GetRelationshipRefColumnsNames() []SemanticViewColumn {
	return r.RelationshipRefColumnNames
}

func (r *SemanticViewRelationship) WithRelationshipRefColumnsNames(keys []SemanticViewColumn) *SemanticViewRelationship {
	r.RelationshipRefColumnNames = keys

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

func (v *semanticViews) DescribeSemanticViewDetails(ctx context.Context, id SchemaObjectIdentifier) (*SemanticViewDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseSemanticViewDescribeOutput(properties, id)
}

func parseSemanticViewDescribeOutput(properties []SemanticViewDetail, id SchemaObjectIdentifier) (*SemanticViewDetails, error) {
	details := &SemanticViewDetails{
		Id:               id,
		Tables:           []SemanticViewTableDetails{},
		Relationships:    []SemanticViewRelationshipDetails{},
		Dimensions:       []SemanticViewDimensionDetails{},
		Facts:            []SemanticViewFactDetails{},
		Metrics:          []SemanticViewMetricDetails{},
		DescribeRowCount: 0,
	}
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
			// loop over all the tables until we either find one with a matching name, or exhaust the set
			// this property should be added to the corresponding table object
			for i := range details.Tables {
				if details.Tables[i].TableNameOrAlias == *prop.ObjectName {
					currentTable = &details.Tables[i]
					break
				}
			}
			// couldn't find a matching table in the details list
			// create a new one with the name and add the property
			if currentTable == nil {
				details.Tables = append(details.Tables, SemanticViewTableDetails{
					TableNameOrAlias: *prop.ObjectName,
				})
				currentTable = &details.Tables[len(details.Tables)-1]
			}
			switch prop.Property {
			case "BASE_TABLE_DATABASE_NAME":
				currentTable.BaseTable.databaseName = prop.PropertyValue
			case "BASE_TABLE_SCHEMA_NAME":
				currentTable.BaseTable.schemaName = prop.PropertyValue
			case "BASE_TABLE_NAME":
				currentTable.BaseTable.name = prop.PropertyValue
			case "PRIMARY_KEY":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentTable.PrimaryKeys)
				if err != nil {
					errs = append(errs, err)
				}
			case "UNIQUE_KEY":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentTable.UniqueKeys)
				if err != nil {
					errs = append(errs, err)
				}
			case "SYNONYMS":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentTable.Synonyms)
				if err != nil {
					errs = append(errs, err)
				}
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
				objectName := ""
				if prop.ObjectName != nil {
					objectName = *prop.ObjectName
				}
				parentEntity := ""
				if prop.ParentEntity != nil {
					parentEntity = *prop.ParentEntity
				}
				details.Relationships = append(details.Relationships, SemanticViewRelationshipDetails{
					RelationshipAlias: objectName,
					ParentEntity:      parentEntity,
				})
				currentRelationship = &details.Relationships[len(details.Relationships)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentRelationship.TableNameOrAlias = prop.PropertyValue
			case "FOREIGN_KEY":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentRelationship.ForeignKeys)
				if err != nil {
					errs = append(errs, err)
				}
			case "REF_TABLE":
				currentRelationship.RefTableNameOrAlias = prop.PropertyValue
			case "REF_KEY":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentRelationship.RefKeys)
				if err != nil {
					errs = append(errs, err)
				}
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
				objectName := ""
				if prop.ObjectName != nil {
					objectName = *prop.ObjectName
				}
				parentEntity := ""
				if prop.ParentEntity != nil {
					parentEntity = *prop.ParentEntity
				}
				details.Dimensions = append(details.Dimensions, SemanticViewDimensionDetails{
					DimensionAlias: objectName,
					ParentEntity:   parentEntity,
				})
				currentDimension = &details.Dimensions[len(details.Dimensions)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentDimension.Properties.TableNameOrAlias = prop.PropertyValue
			case "EXPRESSION":
				currentDimension.Properties.Expression = prop.PropertyValue
			case "DATA_TYPE":
				currentDimension.Properties.DataType = prop.PropertyValue
			case "ACCESS_MODIFIER":
				currentDimension.Properties.AccessModifier = prop.PropertyValue
			case "SYNONYMS":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentDimension.Synonyms)
				if err != nil {
					errs = append(errs, err)
				}
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
				objectName := ""
				if prop.ObjectName != nil {
					objectName = *prop.ObjectName
				}
				parentEntity := ""
				if prop.ParentEntity != nil {
					parentEntity = *prop.ParentEntity
				}
				details.Facts = append(details.Facts, SemanticViewFactDetails{
					FactAlias:    objectName,
					ParentEntity: parentEntity,
				})
				currentFact = &details.Facts[len(details.Facts)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentFact.Properties.TableNameOrAlias = prop.PropertyValue
			case "EXPRESSION":
				currentFact.Properties.Expression = prop.PropertyValue
			case "DATA_TYPE":
				currentFact.Properties.DataType = prop.PropertyValue
			case "ACCESS_MODIFIER":
				currentFact.Properties.AccessModifier = prop.PropertyValue
			case "SYNONYMS":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentFact.Synonyms)
				if err != nil {
					errs = append(errs, err)
				}
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
				objectName := ""
				if prop.ObjectName != nil {
					objectName = *prop.ObjectName
				}
				parentEntity := ""
				if prop.ParentEntity != nil {
					parentEntity = *prop.ParentEntity
				}
				details.Metrics = append(details.Metrics, SemanticViewMetricDetails{
					MetricAlias:  objectName,
					ParentEntity: parentEntity,
				})
				currentMetric = &details.Metrics[len(details.Metrics)-1]
			}
			switch prop.Property {
			case "TABLE":
				currentMetric.Properties.TableNameOrAlias = prop.PropertyValue
			case "EXPRESSION":
				currentMetric.Properties.Expression = prop.PropertyValue
			case "DATA_TYPE":
				currentMetric.Properties.DataType = prop.PropertyValue
			case "ACCESS_MODIFIER":
				currentMetric.Properties.AccessModifier = prop.PropertyValue
			case "SYNONYMS":
				err := json.Unmarshal([]byte(prop.PropertyValue), &currentMetric.Synonyms)
				if err != nil {
					errs = append(errs, err)
				}
			case "COMMENT":
				currentMetric.Comment = prop.PropertyValue
			}
		}
	}

	return details, errors.Join(errs...)
}
