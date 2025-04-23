package datatypes

// VariantDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-semistructured#variant
// It does not have synonyms. It does not have any attributes.
type VariantDataType struct {
	underlyingType string
}

func (t *VariantDataType) ToSql() string {
	return t.underlyingType
}

func (t *VariantDataType) ToLegacyDataTypeSql() string {
	return VariantLegacyDataType
}

func (t *VariantDataType) Canonical() string {
	return VariantLegacyDataType
}

func (t *VariantDataType) ToSqlWithoutUnknowns() string {
	return t.ToSql()
}

var VariantDataTypeSynonyms = []string{VariantLegacyDataType}

func parseVariantDataTypeRaw(raw sanitizedDataTypeRaw) (*VariantDataType, error) {
	return &VariantDataType{raw.matchedByType}, nil
}
