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
	Skipped                 bool
	Manual                  bool
}

func (f SchemaField) IsListOrSet() bool {
	return f.SchemaType == schema.TypeList || f.SchemaType == schema.TypeSet
}

var (
	mapSchemaObjectIdentifierSlice  = genhelpers.MapSlice("sdk.SchemaObjectIdentifier", "FullyQualifiedName")
	mapAccountIdentifierSlice       = genhelpers.MapSlice("sdk.AccountIdentifier", "Name")
	mapAccountObjectIdentifierSlice = genhelpers.MapSlice("sdk.AccountObjectIdentifier", "Name")
)

// TODO [SNOW-1501905]: handle other basic type variants
// TODO [SNOW-1501905]: handle any other interface (error)
// TODO [SNOW-1501905]: handle remaining slices (enums, other numerics)
// TODO [SNOW-1501905]: handle structs (chosen one or all)
func MapToSchemaField(field genhelpers.Field) SchemaField {
	isPointer := field.IsPointer()
	isInterface := field.IsInterface()
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	name := genhelpers.ToSnakeCase(field.Name)
	switch concreteTypeWithoutPtr {
	case "string":
		return schemaField(name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.Identity)
	case "[]string":
		return schemaField(name, schema.TypeList, field.Name, isPointer, isInterface, genhelpers.Identity)
	case "int":
		return schemaField(name, schema.TypeInt, field.Name, isPointer, isInterface, genhelpers.Identity)
	case "float64":
		return schemaField(name, schema.TypeFloat, field.Name, isPointer, isInterface, genhelpers.Identity)
	case "bool":
		return schemaField(name, schema.TypeBool, field.Name, isPointer, isInterface, genhelpers.Identity)
	case "time.Time":
		return schemaField(name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.ToString)
	case "sdk.AccountObjectIdentifier":
		return schemaField(name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.Name)
	case "sdk.AccountIdentifier", "sdk.ExternalObjectIdentifier", "sdk.DatabaseObjectIdentifier",
		"sdk.SchemaObjectIdentifier", "sdk.TableColumnIdentifier":
		return schemaField(name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.FullyQualifiedName)
	case "sdk.ObjectIdentifier":
		return schemaField(name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.FullyQualifiedName)
	case "datatypes.DataType":
		return schemaField(name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.ToSql)
	case "[]sdk.SchemaObjectIdentifier":
		return schemaField(name, schema.TypeList, field.Name, isPointer, isInterface, mapSchemaObjectIdentifierSlice)
	case "[]sdk.AccountIdentifier":
		return schemaField(name, schema.TypeList, field.Name, isPointer, isInterface, mapAccountIdentifierSlice)
	case "[]sdk.AccountObjectIdentifier":
		return schemaField(name, schema.TypeList, field.Name, isPointer, isInterface, mapAccountObjectIdentifierSlice)
	}

	underlyingTypeWithoutPtr, _ := strings.CutPrefix(field.UnderlyingType, "*")
	isSdkDeclaredObject := strings.HasPrefix(concreteTypeWithoutPtr, "sdk.")
	switch {
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "string":
		return schemaField(name, schema.TypeString, field.Name, isPointer, isInterface, genhelpers.CastToString)
	case isSdkDeclaredObject && underlyingTypeWithoutPtr == "int":
		return schemaField(name, schema.TypeInt, field.Name, isPointer, isInterface, genhelpers.CastToInt)
	}
	return schemaField(name, schema.TypeInvalid, field.Name, isPointer, isInterface, genhelpers.Identity)
}

func schemaField(name string, schemaType schema.ValueType, originalName string, isPointer, isInterface bool, mapper genhelpers.Mapper) SchemaField {
	return SchemaField{
		Name:                    name,
		SchemaType:              schemaType,
		OriginalName:            originalName,
		IsOriginalTypePointer:   isPointer,
		IsOriginalTypeInterface: isInterface,
		Mapper:                  mapper,
	}
}
