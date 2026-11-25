package datatypes

// VariantDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-semistructured#variant
// It does not have synonyms. It does not have any attributes.
type VariantDataType struct {
	plainDataType
}

func NewVariantDataType() DataType {
	return newVariantDataType(VariantLegacyDataType)
}

func newVariantDataType(underlyingType string) *VariantDataType {
	return &VariantDataType{plainDataType{underlyingType, VariantLegacyDataType}}
}

var VariantDataTypeSynonyms = []string{VariantLegacyDataType}

func parseVariantDataTypeRaw(raw sanitizedDataTypeRaw) (*VariantDataType, error) {
	return newVariantDataType(raw.matchedByType), nil
}
