package datatypes

// FloatDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-numeric#data-types-for-floating-point-numbers
// It does have synonyms. It does not have any attributes.
type FloatDataType struct {
	plainDataType
}

func NewFloatDataType() DataType {
	return newFloatDataType(FloatLegacyDataType)
}

func newFloatDataType(underlyingType string) *FloatDataType {
	return &FloatDataType{plainDataType{underlyingType, FloatLegacyDataType}}
}

var FloatDataTypeSynonyms = []string{"FLOAT8", "FLOAT4", FloatLegacyDataType, "DOUBLE PRECISION", "DOUBLE", "REAL"}

func parseFloatDataTypeRaw(raw sanitizedDataTypeRaw) (*FloatDataType, error) {
	return newFloatDataType(raw.matchedByType), nil
}
