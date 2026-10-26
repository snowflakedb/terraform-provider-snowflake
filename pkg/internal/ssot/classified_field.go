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
	ForceNew bool // set on identity fields; name uses IdentityNameForceNew, database and schema use IdentityHierarchyForceNew
	genhelpers.SdkValueType
}

// IdentityFields returns the required name fields for obj.IdType,
// in New<IdType> argument order.
func IdentityFields(obj Object) []ClassifiedField {
	name := stringField("name", true)
	name.ForceNew = obj.IdentityNameForceNew
	database := stringField("database", true)
	database.ForceNew = obj.IdentityHierarchyForceNew
	schemaName := stringField("schema", true)
	schemaName.ForceNew = obj.IdentityHierarchyForceNew
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
		if attr.Required || attr.NativeBool {
			return ClassifiedField{Name: attr.Name, Storage: AttrBool, Required: attr.Required}
		}
		return ClassifiedField{Name: attr.Name, Storage: AttrString, SdkValueType: genhelpers.SdkValueType{SdkType: "bool", SdkKind: "bool"}}
	case AttrEnum:
		return ClassifiedField{
			Name:         attr.Name,
			Storage:      AttrString,
			Required:     attr.Required,
			SdkValueType: genhelpers.SdkValueType{SdkType: "sdk." + attr.Enum.sdkTypeName(), SdkKind: "enum"},
		}
	case AttrIdentifier:
		return ClassifiedField{
			Name:         attr.Name,
			Storage:      AttrString,
			Required:     attr.Required,
			SdkValueType: genhelpers.SdkValueType{SdkType: "sdk." + attr.IdentifierType, SdkKind: "identifier"},
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
	if attr.Elem == nil || attr.IsObjectList() {
		if attr.IsObjectList() {
			// Nested block (list/set of objects): no scalar Elem to classify. The generator consuming
			// this (model builder, resource assertion) decides what a Storage with no Elem means —
			// e.g. an unsupported-type stub, or a Required constructor using a generator-level type
			// override (ModelRequestType) since there is no mechanical mapping from a nested block to
			// a concrete Go request type.
			return ClassifiedField{Name: attr.Name, Storage: storage, Required: attr.Required}
		}
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
