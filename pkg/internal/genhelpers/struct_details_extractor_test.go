package genhelpers

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// TODO [SNOW-1501905]: do we need any of:
//   - (?) slice of pointers to interface
//   - (?) slice of pointers to structs
//   - (?) slice of pointers to basic
//   - (?) pointer to slice
//
// TODO [SNOW-1501905]: test type of slice fields
func Test_ExtractStructDetails(t *testing.T) {
	type testIntEnum int
	type testStruct struct {
		unexportedString     string
		unexportedInt        int
		unexportedBool       bool
		unexportedFloat64    float64
		unexportedStringPtr  *string
		unexportedIntPtr     *int
		unexportedBoolPtr    *bool
		unexportedFloat64Ptr *float64

		unexportedTime    time.Time
		unexportedTimePtr *time.Time

		unexportedStringEnum    sdk.WarehouseType
		unexportedStringEnumPtr *sdk.WarehouseType
		unexportedIntEnum       testIntEnum
		unexportedIntEnumPtr    *testIntEnum

		unexportedAccountIdentifier           sdk.AccountIdentifier
		unexportedExternalObjectIdentifier    sdk.ExternalObjectIdentifier
		unexportedAccountObjectIdentifier     sdk.AccountObjectIdentifier
		unexportedDatabaseObjectIdentifier    sdk.DatabaseObjectIdentifier
		unexportedSchemaObjectIdentifier      sdk.SchemaObjectIdentifier
		unexportedTableColumnIdentifier       sdk.TableColumnIdentifier
		unexportedAccountIdentifierPtr        *sdk.AccountIdentifier
		unexportedExternalObjectIdentifierPtr *sdk.ExternalObjectIdentifier
		unexportedAccountObjectIdentifierPtr  *sdk.AccountObjectIdentifier
		unexportedDatabaseObjectIdentifierPtr *sdk.DatabaseObjectIdentifier
		unexportedSchemaObjectIdentifierPtr   *sdk.SchemaObjectIdentifier
		unexportedTableColumnIdentifierPtr    *sdk.TableColumnIdentifier

		unexportedStringSlice     []string
		unexportedIntSlice        []int
		unexportedStringEnumSlice []sdk.WarehouseType
		unexportedIdentifierSlice []sdk.SchemaObjectIdentifier

		unexportedInterface sdk.ObjectIdentifier
		unexportedStruct    sdk.FileFormatTypeOptions

		ExportedString     string
		ExportedInt        int
		ExportedBool       bool
		ExportedFloat64    float64
		ExportedStringPtr  *string
		ExportedIntPtr     *int
		ExportedBoolPtr    *bool
		ExportedFloat64Ptr *float64
	}

	assertFieldExtracted := func(field Field, expectedName string, expectedConcreteType string, expectedUnderlyingType string) {
		assert.Equal(t, expectedName, field.Name)
		assert.Equal(t, expectedConcreteType, field.ConcreteType)
		assert.Equal(t, expectedUnderlyingType, field.UnderlyingType)
	}

	t.Run("test struct details extraction", func(t *testing.T) {
		structDetails := ExtractStructDetails(testStruct{})

		assert.Equal(t, "genhelpers.testStruct", structDetails.Name)

		assertFieldExtracted(structDetails.Fields[0], "unexportedString", "string", "string")
		assertFieldExtracted(structDetails.Fields[1], "unexportedInt", "int", "int")
		assertFieldExtracted(structDetails.Fields[2], "unexportedBool", "bool", "bool")
		assertFieldExtracted(structDetails.Fields[3], "unexportedFloat64", "float64", "float64")
		assertFieldExtracted(structDetails.Fields[4], "unexportedStringPtr", "*string", "*string")
		assertFieldExtracted(structDetails.Fields[5], "unexportedIntPtr", "*int", "*int")
		assertFieldExtracted(structDetails.Fields[6], "unexportedBoolPtr", "*bool", "*bool")
		assertFieldExtracted(structDetails.Fields[7], "unexportedFloat64Ptr", "*float64", "*float64")

		assertFieldExtracted(structDetails.Fields[8], "unexportedTime", "time.Time", "struct")
		assertFieldExtracted(structDetails.Fields[9], "unexportedTimePtr", "*time.Time", "*struct")

		assertFieldExtracted(structDetails.Fields[10], "unexportedStringEnum", "sdk.WarehouseType", "string")
		assertFieldExtracted(structDetails.Fields[11], "unexportedStringEnumPtr", "*sdk.WarehouseType", "*string")
		assertFieldExtracted(structDetails.Fields[12], "unexportedIntEnum", "genhelpers.testIntEnum", "int")
		assertFieldExtracted(structDetails.Fields[13], "unexportedIntEnumPtr", "*genhelpers.testIntEnum", "*int")

		assertFieldExtracted(structDetails.Fields[14], "unexportedAccountIdentifier", "sdk.AccountIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[15], "unexportedExternalObjectIdentifier", "sdk.ExternalObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[16], "unexportedAccountObjectIdentifier", "sdk.AccountObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[17], "unexportedDatabaseObjectIdentifier", "sdk.DatabaseObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[18], "unexportedSchemaObjectIdentifier", "sdk.SchemaObjectIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[19], "unexportedTableColumnIdentifier", "sdk.TableColumnIdentifier", "struct")
		assertFieldExtracted(structDetails.Fields[20], "unexportedAccountIdentifierPtr", "*sdk.AccountIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[21], "unexportedExternalObjectIdentifierPtr", "*sdk.ExternalObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[22], "unexportedAccountObjectIdentifierPtr", "*sdk.AccountObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[23], "unexportedDatabaseObjectIdentifierPtr", "*sdk.DatabaseObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[24], "unexportedSchemaObjectIdentifierPtr", "*sdk.SchemaObjectIdentifier", "*struct")
		assertFieldExtracted(structDetails.Fields[25], "unexportedTableColumnIdentifierPtr", "*sdk.TableColumnIdentifier", "*struct")

		assertFieldExtracted(structDetails.Fields[26], "unexportedStringSlice", "[]string", "slice")
		assertFieldExtracted(structDetails.Fields[27], "unexportedIntSlice", "[]int", "slice")
		assertFieldExtracted(structDetails.Fields[28], "unexportedStringEnumSlice", "[]sdk.WarehouseType", "slice")
		assertFieldExtracted(structDetails.Fields[29], "unexportedIdentifierSlice", "[]sdk.SchemaObjectIdentifier", "slice")

		assertFieldExtracted(structDetails.Fields[30], "unexportedInterface", "sdk.ObjectIdentifier", "interface")
		assertFieldExtracted(structDetails.Fields[31], "unexportedStruct", "sdk.FileFormatTypeOptions", "struct")
	})
}

func Test_Field_IsSlice(t *testing.T) {
	t.Run("is a slice", func(t *testing.T) {
		field := Field{ConcreteType: "[]any_type"}

		assert.True(t, field.IsSlice())
	})

	t.Run("is not a slice", func(t *testing.T) {
		field := Field{ConcreteType: "*any_type"}

		assert.False(t, field.IsSlice())

		field = Field{ConcreteType: "*[]any_type"}

		assert.False(t, field.IsSlice())

		field = Field{ConcreteType: "any_type"}

		assert.False(t, field.IsSlice())
	})
}

func Test_Field_IsPointer(t *testing.T) {
	t.Run("is a pointer", func(t *testing.T) {
		field := Field{ConcreteType: "*any_type"}

		assert.True(t, field.IsPointer())
	})

	t.Run("is not a pointer", func(t *testing.T) {
		field := Field{ConcreteType: "[]any_type"}

		assert.False(t, field.IsPointer())

		field = Field{ConcreteType: "[]*any_type"}

		assert.False(t, field.IsPointer())

		field = Field{ConcreteType: "any_type"}

		assert.False(t, field.IsPointer())
	})
}
