package datatypes

// TimestampLtzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does have optional precision attribute.
// Precision can be known or unknown.
type TimestampLtzDataType struct {
	sizedDataType
}

func NewTimestampLtzDataType() DataType {
	return newTimestampLtzDataType(DefaultTimestampPrecision, TimestampLtzLegacyDataType, false)
}

func newTimestampLtzDataType(precision int, underlyingType string, precisionKnown bool) *TimestampLtzDataType {
	return &TimestampLtzDataType{sizedDataType{precision, underlyingType, TimestampLtzLegacyDataType, precisionKnown}}
}

var TimestampLtzDataTypeSynonyms = []string{TimestampLtzLegacyDataType, "TIMESTAMPLTZ", "TIMESTAMP WITH LOCAL TIME ZONE"}

func parseTimestampLtzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampLtzDataType, error) {
	precision, known, err := parseSize(raw, "timestamp ltz", "precision", DefaultTimestampPrecision)
	if err != nil {
		return nil, err
	}
	return newTimestampLtzDataType(precision, raw.matchedByType, known), nil
}
