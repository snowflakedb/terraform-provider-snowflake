package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// FieldPair holds the pairing between a single db row struct field and the corresponding plain SDK struct field.
// It is used by the convert.tmpl template to generate type-safe conversion code.
type FieldPair struct {
	DbFieldName    string
	PlainFieldName string
	DbKind         string
	PlainKind      string
	IsEnum         bool
	IsJson         bool
}

// AssignmentKind is the conversion strategy discriminator used in convert.tmpl.
type AssignmentKind string

const (
	AssignmentKindDirect                          AssignmentKind = "Direct"
	AssignmentKindStringToBool                    AssignmentKind = "StringToBool"
	AssignmentKindStringToStringArray             AssignmentKind = "StringToStringArray"
	AssignmentKindStringToEnum                    AssignmentKind = "StringToEnum"
	AssignmentKindStringToJson                    AssignmentKind = "StringToJson"
	AssignmentKindStringToIdentifier              AssignmentKind = "StringToIdentifier"
	AssignmentKindNullableToNullable              AssignmentKind = "NullableToNullable"
	AssignmentKindNullableToRequired              AssignmentKind = "NullableToRequired"
	AssignmentKindNullableToIdentifier            AssignmentKind = "NullableToIdentifier"
	AssignmentKindNullableToEnum                  AssignmentKind = "NullableToEnum"
	AssignmentKindNullableToStringArray           AssignmentKind = "NullableToStringArray"
	AssignmentKindNullableStringToNullableBool    AssignmentKind = "NullableStringToNullableBool"
	AssignmentKindNullableStringToIdentifierArray AssignmentKind = "NullableStringToIdentifierArray"
	AssignmentKindUnsupported                     AssignmentKind = "Unsupported"
)

func (fp FieldPair) AssignmentKindDirect() bool {
	return fp.AssignmentKind() == AssignmentKindDirect
}

func (fp FieldPair) AssignmentKindStringToBool() bool {
	return fp.AssignmentKind() == AssignmentKindStringToBool
}

func (fp FieldPair) AssignmentKindStringToStringArray() bool {
	return fp.AssignmentKind() == AssignmentKindStringToStringArray
}

func (fp FieldPair) AssignmentKindStringToEnum() bool {
	return fp.AssignmentKind() == AssignmentKindStringToEnum
}

func (fp FieldPair) AssignmentKindStringToJson() bool {
	return fp.AssignmentKind() == AssignmentKindStringToJson
}

func (fp FieldPair) AssignmentKindStringToIdentifier() bool {
	return fp.AssignmentKind() == AssignmentKindStringToIdentifier
}

func (fp FieldPair) AssignmentKindNullableToNullable() bool {
	return fp.AssignmentKind() == AssignmentKindNullableToNullable
}

func (fp FieldPair) AssignmentKindNullableToRequired() bool {
	return fp.AssignmentKind() == AssignmentKindNullableToRequired
}

func (fp FieldPair) AssignmentKindNullableToIdentifier() bool {
	return fp.AssignmentKind() == AssignmentKindNullableToIdentifier
}

func (fp FieldPair) AssignmentKindNullableToEnum() bool {
	return fp.AssignmentKind() == AssignmentKindNullableToEnum
}

func (fp FieldPair) AssignmentKindNullableToStringArray() bool {
	return fp.AssignmentKind() == AssignmentKindNullableToStringArray
}

func (fp FieldPair) AssignmentKindNullableStringToNullableBool() bool {
	return fp.AssignmentKind() == AssignmentKindNullableStringToNullableBool
}

func (fp FieldPair) AssignmentKindNullableStringToIdentifierArray() bool {
	return fp.AssignmentKind() == AssignmentKindNullableStringToIdentifierArray
}

// IdentifierArrayElementType returns the element type name for NullableStringToIdentifierArray fields.
// E.g., for []SchemaObjectIdentifier it returns "SchemaObjectIdentifier".
func (fp FieldPair) IdentifierArrayElementType() string {
	return genhelpers.TypeWithoutPointer(strings.TrimPrefix(fp.PlainKind, "[]"))
}

// AssignmentKind returns the conversion strategy for this field pair.
// The returned value is used as a discriminator via the boolean predicate methods below.
func (fp FieldPair) AssignmentKind() AssignmentKind {
	if fp.DbKind == fp.PlainKind {
		return AssignmentKindDirect
	}

	switch fp.DbKind {
	case "string":
		switch {
		case fp.PlainKind == "bool":
			return AssignmentKindStringToBool
		case fp.PlainKind == "[]string":
			return AssignmentKindStringToStringArray
		case fp.IsEnum:
			return AssignmentKindStringToEnum
		case fp.IsJson:
			return AssignmentKindStringToJson
		case genhelpers.IsIdentifierType(fp.PlainKind):
			return AssignmentKindStringToIdentifier
		}

	case "sql.NullString":
		switch {
		case fp.PlainKind == "*string":
			return AssignmentKindNullableToNullable
		case fp.PlainKind == "string":
			return AssignmentKindNullableToRequired
		case fp.PlainKind == "[]string":
			return AssignmentKindNullableToStringArray
		case fp.PlainKind == "*bool":
			return AssignmentKindNullableStringToNullableBool
		case strings.HasPrefix(fp.PlainKind, "[]") && genhelpers.IsIdentifierType(strings.TrimPrefix(fp.PlainKind, "[]")):
			return AssignmentKindNullableStringToIdentifierArray
		case fp.IsEnum:
			return AssignmentKindNullableToEnum
		case genhelpers.IsIdentifierType(fp.PlainKind):
			return AssignmentKindNullableToIdentifier
		}

	case "sql.NullBool":
		switch fp.PlainKind {
		case "*bool":
			return AssignmentKindNullableToNullable
		case "bool":
			return AssignmentKindNullableToRequired
		}

	case "sql.NullInt64":
		switch fp.PlainKind {
		case "*int":
			return AssignmentKindNullableToNullable
		case "int":
			return AssignmentKindNullableToRequired
		}

	case "sql.NullTime":
		switch fp.PlainKind {
		case "*time.Time":
			return AssignmentKindNullableToNullable
		case "time.Time":
			return AssignmentKindNullableToRequired
		}
	}

	return AssignmentKindUnsupported
}

// nullTypeBaseNames maps a sql.Null* type to the base name used in mapNull* helper function names.
var nullTypeBaseNames = map[string]string{
	"sql.NullString": "String",
	"sql.NullBool":   "Bool",
	"sql.NullInt64":  "Int",
	"sql.NullTime":   "Time",
}

// NullMapFunc returns the mapNullX helper function name for nullable field pairs.
func (fp FieldPair) NullMapFunc() string {
	if base, ok := nullTypeBaseNames[fp.DbKind]; ok {
		return "mapNull" + base
	}
	return ""
}

// NullRequiredMapFunc returns the mapNullX...ToNonNullableField helper function name (for Null* → non-pointer plain field).
// Used analogously to NullMapFunc but for required (non-pointer) targets.
func (fp FieldPair) NullRequiredMapFunc() string {
	if base, ok := nullTypeBaseNames[fp.DbKind]; ok {
		return "mapNull" + base + "ToNonNullableField"
	}
	return ""
}
