package datatypes

const (
	DefaultBinarySize = 8 * 1024 * 1024

	MaxBinarySize = 64 * 1024 * 1024
)

// BinaryDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-text#data-types-for-binary-strings
// It does have synonyms that allow specifying size.
// Size can be known or unknown.
type BinaryDataType struct {
	sizedDataType
}

func NewBinaryDataType() DataType {
	return newBinaryDataType(DefaultBinarySize, BinaryLegacyDataType, false)
}

func newBinaryDataType(size int, underlyingType string, sizeKnown bool) *BinaryDataType {
	return &BinaryDataType{sizedDataType{size, underlyingType, BinaryLegacyDataType, sizeKnown}}
}

var BinaryDataTypeSynonyms = []string{BinaryLegacyDataType, "VARBINARY"}

func parseBinaryDataTypeRaw(raw sanitizedDataTypeRaw) (*BinaryDataType, error) {
	size, known, err := parseSize(raw, "binary", "size", DefaultBinarySize)
	if err != nil {
		return nil, err
	}
	return newBinaryDataType(size, raw.matchedByType, known), nil
}
