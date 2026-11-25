package datatypes

const DefaultDecfloatPrecision = 38

// DecfloatDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-numeric#decfloat
// It doesn't have synonyms. It does have optional precision attribute.
// Currently, the only allowed value is 38. To be consistent with other data types we are not validating this value (e.g., it's possible to set number's precision to 40 and it will fail on Snowflake).
// Precision can be known or unknown.
type DecfloatDataType struct {
	sizedDataType
}

func NewDecfloatDataType() DataType {
	return newDecfloatDataType(DefaultDecfloatPrecision, DecfloatLegacyDataType, false)
}

func newDecfloatDataType(precision int, underlyingType string, precisionKnown bool) *DecfloatDataType {
	return &DecfloatDataType{sizedDataType{precision, underlyingType, DecfloatLegacyDataType, precisionKnown}}
}

var DecfloatDataTypeSynonyms = []string{DecfloatLegacyDataType}

func parseDecfloatDataTypeRaw(raw sanitizedDataTypeRaw) (*DecfloatDataType, error) {
	precision, known, err := parseSize(raw, "decfloat", "precision", DefaultDecfloatPrecision)
	if err != nil {
		return nil, err
	}
	return newDecfloatDataType(precision, raw.matchedByType, known), nil
}
