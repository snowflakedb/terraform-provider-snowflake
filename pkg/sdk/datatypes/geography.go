package datatypes

// GeographyDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-geospatial#geography-data-type
// It does not have synonyms. It does not have any attributes.
type GeographyDataType struct {
	underlyingType string
}

func (t *GeographyDataType) ToSql() string {
	return t.underlyingType
}

func (t *GeographyDataType) ToLegacyDataTypeSql() string {
	return GeographyLegacyDataType
}

func (t *GeographyDataType) Canonical() string {
	return GeographyLegacyDataType
}

func (t *GeographyDataType) ToSqlWithoutUnknowns() string {
	return t.ToSql()
}

var GeographyDataTypeSynonyms = []string{GeographyLegacyDataType}

func parseGeographyDataTypeRaw(raw sanitizedDataTypeRaw) (*GeographyDataType, error) {
	return &GeographyDataType{raw.matchedByType}, nil
}
