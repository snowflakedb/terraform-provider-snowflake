package datatypes

// TimestampTzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does have optional precision attribute.
// Precision can be known or unknown.
type TimestampTzDataType struct {
	sizedDataType
}

func NewTimestampTzDataType() DataType {
	return newTimestampTzDataType(DefaultTimestampPrecision, TimestampTzLegacyDataType, false)
}

func newTimestampTzDataType(precision int, underlyingType string, precisionKnown bool) *TimestampTzDataType {
	return &TimestampTzDataType{sizedDataType{precision, underlyingType, TimestampTzLegacyDataType, precisionKnown}}
}

var TimestampTzDataTypeSynonyms = []string{TimestampTzLegacyDataType, "TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE"}

func parseTimestampTzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampTzDataType, error) {
	precision, known, err := parseSize(raw, "timestamp tz", "precision", DefaultTimestampPrecision)
	if err != nil {
		return nil, err
	}
	return newTimestampTzDataType(precision, raw.matchedByType, known), nil
}
