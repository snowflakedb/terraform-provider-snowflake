package datatypes

// DateDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#date
// It does not have synonyms. It does not have any attributes.
type DateDataType struct {
	plainDataType
}

func NewDateDataType() DataType {
	return newDateDataType(DateLegacyDataType)
}

func newDateDataType(underlyingType string) *DateDataType {
	return &DateDataType{plainDataType{underlyingType, DateLegacyDataType}}
}

var DateDataTypeSynonyms = []string{DateLegacyDataType}

func parseDateDataTypeRaw(raw sanitizedDataTypeRaw) (*DateDataType, error) {
	return newDateDataType(raw.matchedByType), nil
}
