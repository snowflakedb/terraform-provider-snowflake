package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SchemaField struct {
	Name                    string
	SchemaType              schema.ValueType
	OriginalName            string
	IsOriginalTypePointer   bool
	IsOriginalTypeInterface bool
	Mapper                  genhelpers.Mapper
}

// TODO [SNOW-1501905]: handle other basic type variants
// TODO [SNOW-1501905]: handle any other interface (error)
// TODO [SNOW-1501905]: handle slices
// TODO [SNOW-1501905]: handle structs (chosen one or all)
func MapToSchemaField(field genhelpers.Field) SchemaField {
	isPointer := field.IsPointer()
	isInterface := field.IsInterface()
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	name := genhelpers.ToSnakeCase(field.Name)
	switch concreteTypeWithoutPtr {
	case "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.Identity}
	case "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, isInterface, genhelpers.Identity}
	case "float64":
		return SchemaField{name, schema.TypeFloat, field.Name, isPointer, isInterface, genhelpers.Identity}
	case "bool":
		return SchemaField{name, schema.TypeBool, field.Name, isPointer, isInterface, genhelpers.Identity}
	case "time.Time":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.ToString}
	case "sdk.AccountObjectIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.Name}
	case "sdk.AccountIdentifier", "sdk.ExternalObjectIdentifier", "sdk.DatabaseObjectIdentifier",
		"sdk.SchemaObjectIdentifier", "sdk.TableColumnIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.FullyQualifiedName}
	case "sdk.ObjectIdentifier":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.FullyQualifiedName}
	}

	underlyingTypeWithoutPtr, _ := strings.CutPrefix(field.UnderlyingType, "*")
	isSdkDeclaredObject := strings.HasPrefix(concreteTypeWithoutPtr, "sdk.")
	switch {
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "string":
		return SchemaField{name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.CastToString}
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "int":
		return SchemaField{name, schema.TypeInt, field.Name, isPointer, isInterface, genhelpers.CastToInt}
	}
	return SchemaField{name, schema.TypeInvalid, field.Name, isPointer, isInterface, genhelpers.Identity}
}
