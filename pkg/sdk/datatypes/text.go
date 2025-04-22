package datatypes

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	DefaultVarcharLength = 16777216
	DefaultCharLength    = 1
)

// TextDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-text#data-types-for-text-strings
// It does have synonyms that allow specifying length.
// It does have synonyms that allow specifying length but differ with the default length when length is omitted; here called subtypes.
type TextDataType struct {
	length         int
	underlyingType string

	lengthKnown bool
}

func (t *TextDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.length)
}

func (t *TextDataType) ToLegacyDataTypeSql() string {
	return VarcharLegacyDataType
}

func (t *TextDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", VarcharLegacyDataType, t.length)
}

func (t *TextDataType) ToSqlWithoutUnknowns() string {
	switch {
	case t.lengthKnown:
		return fmt.Sprintf("%s(%d)", t.underlyingType, t.length)
	default:
		return fmt.Sprintf("%s", t.underlyingType)
	}
}

var (
	TextDataTypeSynonyms = []string{VarcharLegacyDataType, "STRING", "TEXT", "NVARCHAR2", "NVARCHAR", "CHAR VARYING", "NCHAR VARYING"}
	TextDataTypeSubtypes = []string{"CHARACTER", "CHAR", "NCHAR"}
	AllTextDataTypes     = append(TextDataTypeSynonyms, TextDataTypeSubtypes...)
)

// parseTextDataTypeRaw extracts length from the raw text data type input.
// It returns default if it can't parse arguments, data type is different, or no length argument was provided.
func parseTextDataTypeRaw(raw sanitizedDataTypeRaw) (*TextDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		switch {
		case slices.Contains(TextDataTypeSynonyms, raw.matchedByType):
			return &TextDataType{DefaultVarcharLength, raw.matchedByType, false}, nil
		case slices.Contains(TextDataTypeSubtypes, raw.matchedByType):
			return &TextDataType{DefaultCharLength, raw.matchedByType, false}, nil
		default:
			return nil, fmt.Errorf("unknown text data type: %s", raw.raw)
		}
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return nil, fmt.Errorf(`text %s could not be parsed, use "%s(length)" format`, raw.raw, raw.matchedByType)
	}
	lengthRaw := r[1 : len(r)-1]
	length, err := strconv.Atoi(strings.TrimSpace(lengthRaw))
	if err != nil {
		return nil, fmt.Errorf(`could not parse the varchar's length: "%s", err: %w`, lengthRaw, err)
	}
	return &TextDataType{length, raw.matchedByType, true}, nil
}

func areTextDataTypesTheSame(a, b *TextDataType) bool {
	return a.length == b.length
}

func areTextDataTypesDefinitelyDifferent(a, b *TextDataType) bool {
	var lengthDefinitelyDifferent bool
	if a.lengthKnown && b.lengthKnown {
		lengthDefinitelyDifferent = a.length != b.length
	}
	return lengthDefinitelyDifferent
}
