package datatypes

// BooleanDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-logical
// It does not have synonyms. It does not have any attributes.
type BooleanDataType struct {
	plainDataType
}

func NewBooleanDataType() DataType {
	return newBooleanDataType(BooleanLegacyDataType)
}

func newBooleanDataType(underlyingType string) *BooleanDataType {
	return &BooleanDataType{plainDataType{underlyingType, BooleanLegacyDataType}}
}

var BooleanDataTypeSynonyms = []string{BooleanLegacyDataType}

func parseBooleanDataTypeRaw(raw sanitizedDataTypeRaw) (*BooleanDataType, error) {
	return newBooleanDataType(raw.matchedByType), nil
}
