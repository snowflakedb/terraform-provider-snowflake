package ssot

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// ClassifiedField is one attribute or parameter as a test generator should see it.
// Storage is the Terraform storage kind. Elem is set for a list or set.
type ClassifiedField struct {
	Name     string
	Storage  AttributeType
	Elem     AttributeType
	Required bool
	genhelpers.SdkValueType
}

// IdentityFields returns the required name fields for obj.IdType.
func IdentityFields(obj Object) []ClassifiedField {
	name := stringField("name", true)
	database := stringField("database", true)
	schemaName := stringField("schema", true)
	switch obj.IdType {
	case "AccountObjectIdentifier":
		return []ClassifiedField{name}
	case "DatabaseObjectIdentifier":
		return []ClassifiedField{database, name}
	case "SchemaObjectIdentifier":
		return []ClassifiedField{database, schemaName, name}
	default:
		panic(fmt.Sprintf("object %q: unsupported IdType %q", obj.Name, obj.IdType))
	}
}

// ClassifyAttribute classifies one attribute for test generators.
func ClassifyAttribute(obj Object, attr Attribute) ClassifiedField {
	switch attr.Type {
	case AttrInt:
		return ClassifiedField{Name: attr.Name, Storage: AttrInt, Required: attr.Required}
	case AttrBool:
		if attr.Required {
			return ClassifiedField{Name: attr.Name, Storage: AttrBool, Required: true}
		}
		return ClassifiedField{Name: attr.Name, Storage: AttrString, SdkValueType: genhelpers.SdkValueType{SdkType: "bool", SdkKind: "bool"}}
	case AttrEnum:
		if attr.Enum == nil || attr.Enum.Name == "" {
			panic(fmt.Sprintf("object %q: enum attribute %q has no Enum", obj.Name, attr.Name))
		}
		return ClassifiedField{
			Name:         attr.Name,
			Storage:      AttrString,
			Required:     attr.Required,
			SdkValueType: genhelpers.SdkValueType{SdkType: "sdk." + attr.Enum.Name, SdkKind: "enum"},
		}
	case AttrIdentifier:
		if obj.IdType == "" {
			panic(fmt.Sprintf("object %q: identifier attribute %q has no IdType", obj.Name, attr.Name))
		}
		return ClassifiedField{
			Name:         attr.Name,
			Storage:      AttrString,
			Required:     attr.Required,
			SdkValueType: genhelpers.SdkValueType{SdkType: "sdk." + obj.IdType, SdkKind: "identifier"},
		}
	case AttrString:
		return stringField(attr.Name, attr.Required)
	case AttrList:
		return collectionField(obj, attr, AttrList)
	case AttrSet:
		return collectionField(obj, attr, AttrSet)
	default:
		panic(fmt.Sprintf("object %q: attribute %q has unknown type %q", obj.Name, attr.Name, attr.Type))
	}
}

// ClassifyParameter classifies one catalog parameter by name and kind.
// Primitive kinds stay on the storage type.
func ClassifyParameter(name, kind string) ClassifiedField {
	switch kind {
	case string(AttrBool):
		return ClassifiedField{Name: name, Storage: AttrBool}
	case string(AttrInt):
		return ClassifiedField{Name: name, Storage: AttrInt}
	case string(AttrString), "StringAllowEmpty", string(AttrFloat):
		return stringField(name, false)
	}
	if strings.HasSuffix(kind, "Identifier") {
		return ClassifiedField{
			Name:         name,
			Storage:      AttrString,
			SdkValueType: genhelpers.SdkValueType{SdkType: "sdk." + kind, SdkKind: "identifier"},
		}
	}
	return ClassifiedField{
		Name:         name,
		Storage:      AttrString,
		SdkValueType: genhelpers.SdkValueType{SdkType: "sdk." + kind, SdkKind: "enum"},
	}
}

func collectionField(obj Object, attr Attribute, storage AttributeType) ClassifiedField {
	if attr.Elem == nil {
		panic(fmt.Sprintf("object %q: attribute %q has no Elem", obj.Name, attr.Name))
	}
	var elem AttributeType
	switch attr.Elem.Type {
	case AttrString, AttrEnum, AttrIdentifier:
		elem = AttrString
	case AttrBool:
		elem = AttrBool
	case AttrInt:
		elem = AttrInt
	default:
		panic(fmt.Sprintf("object %q: attribute %q has unsupported elem type %q", obj.Name, attr.Name, attr.Elem.Type))
	}
	return ClassifiedField{Name: attr.Name, Storage: storage, Elem: elem, Required: attr.Required}
}

func stringField(name string, required bool) ClassifiedField {
	return ClassifiedField{Name: name, Storage: AttrString, Required: required}
}
