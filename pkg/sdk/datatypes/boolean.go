package datatypes

// BooleanDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-logical
// It does not have synonyms. It does not have any attributes.
type BooleanDataType struct {
	underlyingType string
}

func (t *BooleanDataType) ToSql() string {
	return t.underlyingType
}

func (t *BooleanDataType) ToLegacyDataTypeSql() string {
	return BooleanLegacyDataType
}

func (t *BooleanDataType) Canonical() string {
	return BooleanLegacyDataType
}

func (t *BooleanDataType) ToSqlWithoutUnknowns() string {
	return t.ToSql()
}

var BooleanDataTypeSynonyms = []string{BooleanLegacyDataType}

func parseBooleanDataTypeRaw(raw sanitizedDataTypeRaw) (*BooleanDataType, error) {
	return &BooleanDataType{raw.matchedByType}, nil
}
