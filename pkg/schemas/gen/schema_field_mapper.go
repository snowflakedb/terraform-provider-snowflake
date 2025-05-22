package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SchemaField struct {
	Name                  string
	SchemaType            schema.ValueType
	OriginalName          string
	IsOriginalTypePointer bool
	Mapper                genhelpers.Mapper
}

// TODO [SNOW-1501905]: handle other basic type variants
// TODO [SNOW-1501905]: handle any other interface (error)
// TODO [SNOW-1501905]: handle slices
// TODO [SNOW-1501905]: handle structs (chosen one or all)
func MapToSchemaField(field genhelpers.Field) SchemaField {
	isPointer := field.IsPointer()
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	name := genhelpers.ToSnakeCase(field.Name)
	switch concreteTypeWithoutPtr {
	case "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, genhelpers.Identity}
	case "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, genhelpers.Identity}
	case "float64":
		return SchemaField{name, schema.TypeFloat, field.Name, isPointer, genhelpers.Identity}
	case "bool":
		return SchemaField{name, schema.TypeBool, field.Name, isPointer, genhelpers.Identity}
	case "time.Time":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, genhelpers.ToString}
	case "sdk.AccountObjectIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, genhelpers.Name}
	case "sdk.AccountIdentifier", "sdk.ExternalObjectIdentifier", "sdk.DatabaseObjectIdentifier",
		"sdk.SchemaObjectIdentifier", "sdk.TableColumnIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, genhelpers.FullyQualifiedName}
	case "sdk.ObjectIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, genhelpers.FullyQualifiedName}
	}

	underlyingTypeWithoutPtr, _ := strings.CutPrefix(field.UnderlyingType, "*")
	isSdkDeclaredObject := strings.HasPrefix(concreteTypeWithoutPtr, "sdk.")
	switch {
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, genhelpers.CastToString}
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, genhelpers.CastToInt}
	}
	return SchemaField{name, schema.TypeInvalid, field.Name, isPointer, genhelpers.Identity}
}
