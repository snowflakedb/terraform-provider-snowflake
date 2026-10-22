package gen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

const (
	KindBool   = "bool"
	KindInt    = "int"
	KindString = "string"
	KindFloat  = "float64"
)

// KindInfo is the single source for converting a kind between its SQL string form and Go, in both
// directions.
type KindInfo struct {
	GoType      string
	FieldType   string // differs from GoType only for StringAllowEmpty
	ReadParser  string
	WriteParser string
}

// InfoForKind resolves a kind's conversion info. Any bare named type that is neither a primitive nor
// an identifier is assumed to be an SDK enum with a To<Enum> converter, so a missing converter fails
// when compiling the generated code rather than here.
func InfoForKind(kind string) (KindInfo, error) {
	switch kind {
	case KindBool:
		return KindInfo{GoType: KindBool, FieldType: KindBool, ReadParser: "strconv.ParseBool", WriteParser: "strconv.ParseBool"}, nil
	case KindInt:
		return KindInfo{GoType: KindInt, FieldType: KindInt, ReadParser: "strconv.Atoi", WriteParser: "strconv.Atoi"}, nil
	case KindString:
		return KindInfo{GoType: KindString, FieldType: KindString, ReadParser: "identityParse", WriteParser: "identityParse"}, nil
	case KindFloat:
		return KindInfo{GoType: KindFloat, FieldType: KindFloat, ReadParser: "ToFloat64", WriteParser: "ToFloat64"}, nil
	case KindOfT[sdkcommons.StringAllowEmpty]():
		// The wrapper exists so an empty string still renders in SQL; readers get a plain string.
		return KindInfo{GoType: KindString, FieldType: kind, ReadParser: "identityParse", WriteParser: "ToStringAllowEmpty"}, nil
	}
	if !isBareGoTypeName(kind) {
		return KindInfo{}, fmt.Errorf("cannot derive parsers for kind %q; expected a primitive, an identifier, or an enum type name", kind)
	}
	parser := "To" + kind
	if _, err := ToObjectIdentifierKind(kind); err == nil {
		parser = "Parse" + kind
	}
	return KindInfo{GoType: kind, FieldType: kind, ReadParser: parser, WriteParser: parser}, nil
}

func isBareGoTypeName(kind string) bool {
	if kind == "" {
		return false
	}
	return !strings.ContainsAny(kind, "*[]{}. ")
}

func KindOfT[T any]() string {
	t := reflect.TypeFor[T]()
	return t.Name()
}

func KindOfTPointer[T any]() string {
	return KindOfPointer(KindOfT[T]())
}

func KindOfTSlice[T any]() string {
	return KindOfSlice(KindOfT[T]())
}

func KindOfPointer(kind string) string {
	return "*" + kind
}

func KindOfSlice(kind string) string {
	return "[]" + kind
}
