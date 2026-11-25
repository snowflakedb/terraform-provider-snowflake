package datatypes

// TimestampNtzDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-datetime#timestamp-ltz-timestamp-ntz-timestamp-tz
// It does have synonyms. It does have optional precision attribute.
// Precision can be known or unknown.
type TimestampNtzDataType struct {
	sizedDataType
}

func NewTimestampNtzDataType() DataType {
	return newTimestampNtzDataType(DefaultTimestampPrecision, TimestampNtzLegacyDataType, false)
}

func NewTimestampNtzDataTypeWithPrecision(precision int) DataType {
	return newTimestampNtzDataType(precision, TimestampNtzLegacyDataType, true)
}

func newTimestampNtzDataType(precision int, underlyingType string, precisionKnown bool) *TimestampNtzDataType {
	return &TimestampNtzDataType{sizedDataType{precision, underlyingType, TimestampNtzLegacyDataType, precisionKnown}}
}

var TimestampNtzDataTypeSynonyms = []string{TimestampNtzLegacyDataType, "TIMESTAMPNTZ", "TIMESTAMP WITHOUT TIME ZONE", "DATETIME"}

func parseTimestampNtzDataTypeRaw(raw sanitizedDataTypeRaw) (*TimestampNtzDataType, error) {
	precision, known, err := parseSize(raw, "timestamp ntz", "precision", DefaultTimestampPrecision)
	if err != nil {
		return nil, err
	}
	return newTimestampNtzDataType(precision, raw.matchedByType, known), nil
}
