package datatypes

// FloatDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-numeric#data-types-for-floating-point-numbers
// It does have synonyms. It does not have any attributes.
type FloatDataType struct {
	underlyingType string
}

func (t *FloatDataType) ToSql() string {
	return t.underlyingType
}

func (t *FloatDataType) ToLegacyDataTypeSql() string {
	return FloatLegacyDataType
}

func (t *FloatDataType) Canonical() string {
	return FloatLegacyDataType
}

func (t *FloatDataType) ToSqlWithoutUnknowns() string {
	return t.ToSql()
}

var FloatDataTypeSynonyms = []string{"FLOAT8", "FLOAT4", FloatLegacyDataType, "DOUBLE PRECISION", "DOUBLE", "REAL"}

func parseFloatDataTypeRaw(raw sanitizedDataTypeRaw) (*FloatDataType, error) {
	return &FloatDataType{raw.matchedByType}, nil
}
