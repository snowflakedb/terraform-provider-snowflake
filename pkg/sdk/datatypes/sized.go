package datatypes

import (
	"fmt"
	"strconv"
	"strings"
)

// sizedDataType is the common implementation of the data types with one size attribute (e.g. VARCHAR(length), TIME(precision)).
// underlyingType keeps the used synonym (e.g. TEXT), legacyType keeps the base type (e.g. VARCHAR).
// sizeKnown is false when the size was not specified and the default was used.
type sizedDataType struct {
	size           int
	underlyingType string
	legacyType     string

	sizeKnown bool
}

func (t *sizedDataType) ToSql() string {
	return fmt.Sprintf("%s(%d)", t.underlyingType, t.size)
}

func (t *sizedDataType) ToLegacyDataTypeSql() string {
	return t.legacyType
}

func (t *sizedDataType) Canonical() string {
	return fmt.Sprintf("%s(%d)", t.legacyType, t.size)
}

func (t *sizedDataType) ToSqlWithoutUnknowns() string {
	switch {
	case t.sizeKnown:
		return fmt.Sprintf("%s(%d)", t.underlyingType, t.size)
	default:
		return t.underlyingType
	}
}

func (t *sizedDataType) asSized() *sizedDataType {
	return t
}

type sized interface {
	asSized() *sizedDataType
}

// parseSize parses the optional size given in parentheses after the matched type (e.g. TIME(3)).
// It returns defaultSize and known=false when no size is given.
func parseSize(raw sanitizedDataTypeRaw, typeName string, sizeName string, defaultSize int) (size int, known bool, err error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" {
		return defaultSize, false, nil
	}
	if !strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")") {
		return 0, false, fmt.Errorf(`%s %s could not be parsed, use "%s(%s)" format`, typeName, raw.raw, raw.matchedByType, sizeName)
	}
	sizeRaw := r[1 : len(r)-1]
	size, err = strconv.Atoi(strings.TrimSpace(sizeRaw))
	if err != nil {
		return 0, false, fmt.Errorf(`could not parse the %s's %s: "%s", err: %w`, typeName, sizeName, sizeRaw, err)
	}
	return size, true, nil
}

func areSizedDataTypesTheSame(a, b *sizedDataType) bool {
	return a.size == b.size
}

func areSizedDataTypesDefinitelyDifferent(a, b *sizedDataType) bool {
	return a.sizeKnown && b.sizeKnown && a.size != b.size
}
