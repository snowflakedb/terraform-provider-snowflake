package datatypes

import (
	"fmt"
)

// MapDataType is based on https://docs.snowflake.com/en/sql-reference/data-types-structured#label-structured-types-specifying-map.
// It does not have synonyms.
// It does not have a semi-structured variant.
// It has the required key and value types, and the optional NOT NULL constraint for the values.
type MapDataType struct {
	keyType        DataType
	valueType      DataType
	valueNotNull   bool
	underlyingType string
}

func (t *MapDataType) ToSql() string {
	return t.format(t.underlyingType, func(dt DataType) string { return dt.ToSql() })
}

func (t *MapDataType) ToLegacyDataTypeSql() string {
	return MapLegacyDataType
}

func (t *MapDataType) Canonical() string {
	return t.format(MapLegacyDataType, func(dt DataType) string { return dt.Canonical() })
}

func (t *MapDataType) ToSqlWithoutUnknowns() string {
	return t.format(t.underlyingType, func(dt DataType) string { return dt.ToSqlWithoutUnknowns() })
}

func (t *MapDataType) format(mapType string, formatType func(DataType) string) string {
	return fmt.Sprintf("%s(%s, %s)", mapType, formatType(t.keyType), formatWithNotNull(formatType(t.valueType), t.valueNotNull))
}

var MapDataTypeSynonyms = []string{MapLegacyDataType}

func parseMapDataTypeRaw(raw sanitizedDataTypeRaw) (*MapDataType, error) {
	args, hasArguments, ok := structuredTypeArguments(raw)
	if !ok || !hasArguments {
		return nil, fmt.Errorf(`map %s could not be parsed, use "%s(key_type, value_type [ NOT NULL ])" format`, raw.raw, raw.matchedByType)
	}
	parts := splitColumnDefs(args)
	if len(parts) != 2 {
		return nil, fmt.Errorf(`map cannot have %d arguments: "%s"; use "%s(key_type, value_type [ NOT NULL ])" format`, len(parts), args, raw.matchedByType)
	}
	keyTypeRaw := parts[0]
	keyType, err := parseNestedDataType(keyTypeRaw)
	if err != nil {
		return nil, fmt.Errorf(`could not parse the map's key type: "%s", err: %w`, keyTypeRaw, err)
	}
	valueTypeRaw, valueNotNull := trimNotNullSuffix(parts[1])
	valueType, err := parseNestedDataType(valueTypeRaw)
	if err != nil {
		return nil, fmt.Errorf(`could not parse the map's value type: "%s", err: %w`, valueTypeRaw, err)
	}
	return &MapDataType{
		keyType:        keyType,
		valueType:      valueType,
		valueNotNull:   valueNotNull,
		underlyingType: raw.matchedByType,
	}, nil
}

func areMapDataTypesTheSame(a, b *MapDataType) bool {
	return a.valueNotNull == b.valueNotNull && AreTheSame(a.keyType, b.keyType) && AreTheSame(a.valueType, b.valueType)
}

// maps are different if:
// - NOT NULL constraint of the values differs
// - key types are definitely different
// - value types are definitely different
func areMapDataTypesDefinitelyDifferent(a, b *MapDataType) bool {
	return a.valueNotNull != b.valueNotNull || AreDefinitelyDifferent(a.keyType, b.keyType) || AreDefinitelyDifferent(a.valueType, b.valueType)
}
