package datatypes

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	DefaultBinarySize = 8 * 1024 * 1024

	MaxBinarySize = 64 * 1024 * 1024
)

// BinaryDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-text#data-types-for-binary-strings
// It does have synonyms that allow specifying size.
// Size can be known or unknown.
type BinaryDataType struct {
	size           int
	underlyingType string

	sizeKnown bool
}

func (t *BinaryDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.size)
}

func (t *BinaryDataType) ToLegacyDataTypeSql() string {
	return BinaryLegacyDataType
}

func (t *BinaryDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", BinaryLegacyDataType, t.size)
}

func (t *BinaryDataType) ToSqlWithoutUnknowns() string {
	switch {
	case t.sizeKnown:
		return fmt.Sprintf("%s(%d)", t.underlyingType, t.size)
	default:
		return fmt.Sprintf("%s", t.underlyingType)
	}
}

var BinaryDataTypeSynonyms = []string{BinaryLegacyDataType, "VARBINARY"}

func parseBinaryDataTypeRaw(raw sanitizedDataTypeRaw) (*BinaryDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		return &BinaryDataType{DefaultBinarySize, raw.matchedByType, false}, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return nil, fmt.Errorf(`binary %s could not be parsed, use "%s(size)" format`, raw.raw, raw.matchedByType)
	}
	sizeRaw := r[1 : len(r)-1]
	size, err := strconv.Atoi(strings.TrimSpace(sizeRaw))
	if err != nil {
		return nil, fmt.Errorf(`could not parse the binary's size: "%s", err: %w`, sizeRaw, err)
	}
	return &BinaryDataType{size, raw.matchedByType, true}, nil
}

func areBinaryDataTypesTheSame(a, b *BinaryDataType) bool {
	return a.size == b.size
}

func areBinaryDataTypesDefinitelyDifferent(a, b *BinaryDataType) bool {
	var sizeDefinitelyDifferent bool
	if a.sizeKnown && b.sizeKnown {
		sizeDefinitelyDifferent = a.size != b.size
	}
	return sizeDefinitelyDifferent
}
