package gen

import (
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func Test_MapToSchemaField(t *testing.T) {
	type expectedValues struct {
		name       string
		schemaType schema.ValueType
		isPointer  bool
		mapper     genhelpers.Mapper
	}

	testCases := []struct {
		field    genhelpers.Field
		expected expectedValues
	}{
		{
			field:    genhelpers.Field{Name: "unexportedString", ConcreteType: "string", UnderlyingType: "string"},
			expected: expectedValues{"unexported_string", schema.TypeString, false, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedInt", ConcreteType: "int", UnderlyingType: "int"},
			expected: expectedValues{"unexported_int", schema.TypeInt, false, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedBool", ConcreteType: "bool", UnderlyingType: "bool"},
			expected: expectedValues{"unexported_bool", schema.TypeBool, false, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedFloat64", ConcreteType: "float64", UnderlyingType: "float64"},
			expected: expectedValues{"unexported_float64", schema.TypeFloat, false, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedStringPtr", ConcreteType: "*string", UnderlyingType: "*string"},
			expected: expectedValues{"unexported_string_ptr", schema.TypeString, true, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedIntPtr", ConcreteType: "*int", UnderlyingType: "*int"},
			expected: expectedValues{"unexported_int_ptr", schema.TypeInt, true, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedBoolPtr", ConcreteType: "*bool", UnderlyingType: "*bool"},
			expected: expectedValues{"unexported_bool_ptr", schema.TypeBool, true, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedFloat64Ptr", ConcreteType: "*float64", UnderlyingType: "*float64"},
			expected: expectedValues{"unexported_float64_ptr", schema.TypeFloat, true, genhelpers.Identity},
		},
		{
			field:    genhelpers.Field{Name: "unexportedTime", ConcreteType: "time.Time", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_time", schema.TypeString, false, genhelpers.ToString},
		},
		{
			field:    genhelpers.Field{Name: "unexportedTimePtr", ConcreteType: "*time.Time", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_time_ptr", schema.TypeString, true, genhelpers.ToString},
		},
		{
			field:    genhelpers.Field{Name: "unexportedStringEnum", ConcreteType: "sdk.WarehouseType", UnderlyingType: "string"},
			expected: expectedValues{"unexported_string_enum", schema.TypeString, false, genhelpers.CastToString},
		},
		{
			field:    genhelpers.Field{Name: "unexportedStringEnumPtr", ConcreteType: "*sdk.WarehouseType", UnderlyingType: "*string"},
			expected: expectedValues{"unexported_string_enum_ptr", schema.TypeString, true, genhelpers.CastToString},
		},
		{
			field:    genhelpers.Field{Name: "unexportedIntEnum", ConcreteType: "sdk.ResourceMonitorLevel", UnderlyingType: "int"},
			expected: expectedValues{"unexported_int_enum", schema.TypeInt, false, genhelpers.CastToInt},
		},
		{
			field:    genhelpers.Field{Name: "unexportedIntEnumPtr", ConcreteType: "*sdk.ResourceMonitorLevel", UnderlyingType: "*int"},
			expected: expectedValues{"unexported_int_enum_ptr", schema.TypeInt, true, genhelpers.CastToInt},
		},
		{
			field:    genhelpers.Field{Name: "unexportedAccountIdentifier", ConcreteType: "sdk.AccountIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_account_identifier", schema.TypeString, false, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedExternalObjectIdentifier", ConcreteType: "sdk.ExternalObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_external_object_identifier", schema.TypeString, false, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedAccountObjectIdentifier", ConcreteType: "sdk.AccountObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_account_object_identifier", schema.TypeString, false, genhelpers.Name},
		},
		{
			field:    genhelpers.Field{Name: "unexportedDatabaseObjectIdentifier", ConcreteType: "sdk.DatabaseObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_database_object_identifier", schema.TypeString, false, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedSchemaObjectIdentifier", ConcreteType: "sdk.SchemaObjectIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_schema_object_identifier", schema.TypeString, false, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedTableColumnIdentifier", ConcreteType: "sdk.TableColumnIdentifier", UnderlyingType: "struct"},
			expected: expectedValues{"unexported_table_column_identifier", schema.TypeString, false, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedAccountIdentifierPtr", ConcreteType: "*sdk.AccountIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_account_identifier_ptr", schema.TypeString, true, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedExternalObjectIdentifierPtr", ConcreteType: "*sdk.ExternalObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_external_object_identifier_ptr", schema.TypeString, true, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedAccountObjectIdentifierPtr", ConcreteType: "*sdk.AccountObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_account_object_identifier_ptr", schema.TypeString, true, genhelpers.Name},
		},
		{
			field:    genhelpers.Field{Name: "unexportedDatabaseObjectIdentifierPtr", ConcreteType: "*sdk.DatabaseObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_database_object_identifier_ptr", schema.TypeString, true, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedSchemaObjectIdentifierPtr", ConcreteType: "*sdk.SchemaObjectIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_schema_object_identifier_ptr", schema.TypeString, true, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedTableColumnIdentifierPtr", ConcreteType: "*sdk.TableColumnIdentifier", UnderlyingType: "*struct"},
			expected: expectedValues{"unexported_table_column_identifier_ptr", schema.TypeString, true, genhelpers.FullyQualifiedName},
		},
		{
			field:    genhelpers.Field{Name: "unexportedInterface", ConcreteType: "sdk.ObjectIdentifier", UnderlyingType: "interface"},
			expected: expectedValues{"unexported_interface", schema.TypeString, false, genhelpers.FullyQualifiedName},
		},
	}

	assertSchemaFieldMapped := func(schemaField SchemaField, originalField genhelpers.Field, expected expectedValues) {
		assert.Equal(t, expected.name, schemaField.Name)
		assert.Equal(t, expected.schemaType, schemaField.SchemaType)
		assert.Equal(t, originalField.Name, schemaField.OriginalName)
		assert.Equal(t, expected.isPointer, schemaField.IsOriginalTypePointer)
		// TODO [SNOW-1501905]: ugly comparison of functions with the current implementation of mapper
		assert.Equal(t, reflect.ValueOf(expected.mapper).Pointer(), reflect.ValueOf(schemaField.Mapper).Pointer())
	}

	for _, tc := range testCases {
		t.Run(tc.field.Name, func(t *testing.T) {
			schemaField := MapToSchemaField(tc.field)

			assertSchemaFieldMapped(schemaField, tc.field, tc.expected)
		})
	}
}
