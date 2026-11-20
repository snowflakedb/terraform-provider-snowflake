package datatypes

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// ObjectDataType is based on:
// - https://docs.snowflake.com/en/sql-reference/data-types-semistructured#object (semi-structured OBJECT),
// - https://docs.snowflake.com/en/sql-reference/data-types-structured#label-structured-types-specifying-object (structured OBJECT).
// It does not have synonyms.
// Semi-structured OBJECT does not have any attributes.
// Structured OBJECT consists of an ordered list of fields (key, value type, and optional NOT NULL constraint); may be empty.
// Structured OBJECT without fields (OBJECT()) is different from the semi-structured OBJECT.
type ObjectDataType struct {
	fields         []objectDataTypeField
	structured     bool
	underlyingType string
}

type objectDataTypeField struct {
	name     string
	dataType DataType
	notNull  bool
}

func (f *objectDataTypeField) format(formatType func(DataType) string) string {
	return formatWithNotNull(fmt.Sprintf("%s %s", f.name, formatType(f.dataType)), f.notNull)
}

func (t *ObjectDataType) ToSql() string {
	return t.format(t.underlyingType, func(dt DataType) string { return dt.ToSql() })
}

func (t *ObjectDataType) ToLegacyDataTypeSql() string {
	return ObjectLegacyDataType
}

func (t *ObjectDataType) Canonical() string {
	return t.format(ObjectLegacyDataType, func(dt DataType) string { return dt.Canonical() })
}

func (t *ObjectDataType) ToSqlWithoutUnknowns() string {
	return t.format(t.underlyingType, func(dt DataType) string { return dt.ToSqlWithoutUnknowns() })
}

func (t *ObjectDataType) format(objectType string, formatType func(DataType) string) string {
	if !t.structured {
		return objectType
	}
	fields := strings.Join(collections.Map(t.fields, func(f objectDataTypeField) string {
		return f.format(formatType)
	}), ", ")
	return fmt.Sprintf("%s(%s)", objectType, fields)
}

var ObjectDataTypeSynonyms = []string{ObjectLegacyDataType}

// parseObjectDataTypeRaw returns the semi-structured OBJECT when no parentheses are given
// and the structured OBJECT when the (possibly empty) list of fields is given in parentheses.
func parseObjectDataTypeRaw(raw sanitizedDataTypeRaw) (*ObjectDataType, error) {
	args, hasArguments, ok := structuredTypeArguments(raw)
	if !ok {
		return nil, fmt.Errorf(`object %s could not be parsed, use "%s" or "%s(key value_type [ NOT NULL ], ...)" format`, raw.raw, raw.matchedByType, raw.matchedByType)
	}
	if !hasArguments {
		return &ObjectDataType{underlyingType: raw.matchedByType}, nil
	}
	if args == "" {
		return &ObjectDataType{
			fields:         make([]objectDataTypeField, 0),
			structured:     true,
			underlyingType: raw.matchedByType,
		}, nil
	}
	fields, err := collections.MapErr(splitColumnDefs(args), parseObjectDataTypeField)
	if err != nil {
		return nil, err
	}
	return &ObjectDataType{
		fields:         fields,
		structured:     true,
		underlyingType: raw.matchedByType,
	}, nil
}

func parseObjectDataTypeField(field string) (objectDataTypeField, error) {
	trimmed, notNull := trimNotNullSuffix(strings.TrimSpace(field))
	parts := strings.Fields(trimmed)
	if len(parts) < 2 {
		return objectDataTypeField{}, fmt.Errorf("could not parse object field: %s, it should contain the following format `<key> <value_type> [ NOT NULL ]`", field)
	}
	name := parts[0]
	valueTypeRaw := strings.TrimSpace(strings.TrimPrefix(trimmed, name))
	valueType, err := parseNestedDataType(valueTypeRaw)
	if err != nil {
		return objectDataTypeField{}, fmt.Errorf(`could not parse the object's value type for key "%s": "%s", err: %w`, name, valueTypeRaw, err)
	}
	return objectDataTypeField{
		name:     name,
		dataType: valueType,
		notNull:  notNull,
	}, nil
}

func areObjectDataTypesTheSame(a, b *ObjectDataType) bool {
	if a.structured != b.structured || len(a.fields) != len(b.fields) {
		return false
	}
	for i := range a.fields {
		aField, bField := a.fields[i], b.fields[i]
		if aField.name != bField.name || aField.notNull != bField.notNull || !AreTheSame(aField.dataType, bField.dataType) {
			return false
		}
	}
	return true
}

// objects are different if:
// - one is structured and the other is semi-structured
// - they have different numbers of fields
// - key or NOT NULL constraint differs for at least one field (the order of keys matters)
// - value type is definitely different for at least one field
func areObjectDataTypesDefinitelyDifferent(a, b *ObjectDataType) bool {
	if a.structured != b.structured || len(a.fields) != len(b.fields) {
		return true
	}
	for i := range a.fields {
		aField, bField := a.fields[i], b.fields[i]
		if aField.name != bField.name || aField.notNull != bField.notNull || AreDefinitelyDifferent(aField.dataType, bField.dataType) {
			return true
		}
	}
	return false
}
