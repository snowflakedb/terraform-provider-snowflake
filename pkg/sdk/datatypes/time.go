package datatypes

const DefaultTimePrecision = 9

// TimeDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#time
// It does not have synonyms. It does have optional precision attribute.
// Precision can be known or unknown.
type TimeDataType struct {
	sizedDataType
}

func NewTimeDataType() DataType {
	return newTimeDataType(DefaultTimePrecision, TimeLegacyDataType, false)
}

func newTimeDataType(precision int, underlyingType string, precisionKnown bool) *TimeDataType {
	return &TimeDataType{sizedDataType{precision, underlyingType, TimeLegacyDataType, precisionKnown}}
}

var TimeDataTypeSynonyms = []string{TimeLegacyDataType}

func parseTimeDataTypeRaw(raw sanitizedDataTypeRaw) (*TimeDataType, error) {
	precision, known, err := parseSize(raw, "time", "precision", DefaultTimePrecision)
	if err != nil {
		return nil, err
	}
	return newTimeDataType(precision, raw.matchedByType, known), nil
}
