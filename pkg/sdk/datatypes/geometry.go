package datatypes

// GeometryDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-geospatial#geometry-data-type
// It does not have synonyms. It does not have any attributes.
type GeometryDataType struct {
	underlyingType string
}

func (t *GeometryDataType) ToSql() string {
	return t.underlyingType
}

func (t *GeometryDataType) ToLegacyDataTypeSql() string {
	return GeometryLegacyDataType
}

func (t *GeometryDataType) Canonical() string {
	return GeometryLegacyDataType
}

func (t *GeometryDataType) ToSqlWithoutUnknowns() string {
	return t.ToSql()
}

var GeometryDataTypeSynonyms = []string{GeometryLegacyDataType}

func parseGeometryDataTypeRaw(raw sanitizedDataTypeRaw) (*GeometryDataType, error) {
	return &GeometryDataType{raw.matchedByType}, nil
}
