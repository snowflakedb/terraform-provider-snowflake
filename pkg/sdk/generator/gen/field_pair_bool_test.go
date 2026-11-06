package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldPair_textBool(t *testing.T) {
	y := FieldPair{DbKind: "string", PlainKind: "bool"}
	require.Equal(t, AssignmentKindStringToBool, y.AssignmentKind())

	explicitY := FieldPair{DbKind: "string", PlainKind: "bool", BoolTrueValue: "Y"}
	require.Equal(t, AssignmentKindStringToBool, explicitY.AssignmentKind())

	parsed := FieldPair{DbKind: "string", PlainKind: "bool", BoolParsed: true}
	require.Equal(t, AssignmentKindStringToBoolParsed, parsed.AssignmentKind())

	on := FieldPair{DbKind: "string", PlainKind: "bool", BoolTrueValue: "ON"}
	require.Equal(t, AssignmentKindStringToBoolValue, on.AssignmentKind())

	nullable := FieldPair{DbKind: "sql.NullString", PlainKind: "*bool"}
	require.Equal(t, AssignmentKindNullableStringToNullableBool, nullable.AssignmentKind())

	nullableParsed := FieldPair{DbKind: "sql.NullString", PlainKind: "bool", BoolParsed: true}
	require.Equal(t, AssignmentKindNullableStringToRequiredBoolParsed, nullableParsed.AssignmentKind())

	nullableYes := FieldPair{DbKind: "sql.NullString", PlainKind: "bool", BoolTrueValue: "YES"}
	require.Equal(t, AssignmentKindNullableStringToRequiredBoolValue, nullableYes.AssignmentKind())
}
