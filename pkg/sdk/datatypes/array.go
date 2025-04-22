package datatypes

// ArrayDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-semistructured#array
// It does not have synonyms. It does not have any attributes.
type ArrayDataType struct {
	underlyingType string
}

func (t *ArrayDataType) ToSql() string {
	return t.underlyingType
}

func (t *ArrayDataType) ToLegacyDataTypeSql() string {
	return ArrayLegacyDataType
}

func (t *ArrayDataType) Canonical() string {
	return ArrayLegacyDataType
}

func (t *ArrayDataType) ToSqlWithoutUnknowns() string {
	return t.ToSql()
}

var ArrayDataTypeSynonyms = []string{ArrayLegacyDataType}

func parseArrayDataTypeRaw(raw sanitizedDataTypeRaw) (*ArrayDataType, error) {
	return &ArrayDataType{raw.matchedByType}, nil
}
