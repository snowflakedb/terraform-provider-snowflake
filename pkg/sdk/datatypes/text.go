package datatypes

import (
	"fmt"
	"slices"
)

const (
	DefaultVarcharLength = 16 * 1024 * 1024
	DefaultCharLength    = 1

	MaxVarcharLength = 128 * 1024 * 1024
)

// TextDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-text#data-types-for-text-strings
// It does have synonyms that allow specifying length.
// Length can be known or unknown.
// It does have synonyms that allow specifying length but differ with the default length when length is omitted; here called subtypes.
type TextDataType struct {
	sizedDataType
}

func NewVarcharDataType() DataType {
	return newTextDataType(DefaultVarcharLength, VarcharLegacyDataType, false)
}

func NewVarcharDataTypeWithLength(length int) DataType {
	return newTextDataType(length, VarcharLegacyDataType, true)
}

func NewCharDataType() DataType {
	return newTextDataType(DefaultCharLength, "CHAR", false)
}

func NewCharDataTypeWithLength(length int) DataType {
	return newTextDataType(length, "CHAR", true)
}

func newTextDataType(length int, underlyingType string, lengthKnown bool) *TextDataType {
	return &TextDataType{sizedDataType{length, underlyingType, VarcharLegacyDataType, lengthKnown}}
}

var (
	TextDataTypeSynonyms = []string{VarcharLegacyDataType, "STRING", "TEXT", "NVARCHAR2", "NVARCHAR", "CHAR VARYING", "NCHAR VARYING"}
	TextDataTypeSubtypes = []string{"CHARACTER", "CHAR", "NCHAR"}
	AllTextDataTypes     = append(TextDataTypeSynonyms, TextDataTypeSubtypes...)
)

// parseTextDataTypeRaw extracts length from the raw text data type input.
// It returns default if it can't parse arguments, data type is different, or no length argument was provided.
func parseTextDataTypeRaw(raw sanitizedDataTypeRaw) (*TextDataType, error) {
	var defaultLength int
	switch {
	case slices.Contains(TextDataTypeSynonyms, raw.matchedByType):
		defaultLength = DefaultVarcharLength
	case slices.Contains(TextDataTypeSubtypes, raw.matchedByType):
		defaultLength = DefaultCharLength
	default:
		return nil, fmt.Errorf("unknown text data type: %s", raw.raw)
	}
	length, known, err := parseSize(raw, "varchar", "length", defaultLength)
	if err != nil {
		return nil, err
	}
	return newTextDataType(length, raw.matchedByType, known), nil
}
