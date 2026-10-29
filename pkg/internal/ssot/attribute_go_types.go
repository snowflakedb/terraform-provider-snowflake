package ssot

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// AttributeGoTypes returns the concrete Go type and underlying kind for one column,
// without the OptionalInShow pointer. ColumnField applies that pointer.
func AttributeGoTypes(obj Object, attr Attribute) (concrete, underlying string) {
	switch attr.Type {
	case AttrString:
		return "string", "string"
	case AttrInt:
		return "int", "int"
	case AttrBool:
		return "bool", "bool"
	case AttrFloat:
		return "float64", "float64"
	case AttrTime:
		return "time.Time", "struct"
	case AttrIdentifier:
		return "sdk." + attr.IdentifierType, "struct"
	case AttrEnum:
		return "sdk." + attr.ColumnEnum().Name, "string"
	case AttrList, AttrSet:
		if attr.Elem == nil {
			panic(fmt.Sprintf("object %q: collection attribute %q has no Elem", obj.Name, attr.Name))
		}
		if attr.Elem.Type == AttrEnum {
			return "[]string", "slice"
		}
		elemConcrete, _ := AttributeGoTypes(obj, *attr.Elem)
		return "[]" + elemConcrete, "slice"
	case AttrObject:
		// ColumnField and object assertions need the sdk package. The SDK paired-struct path strips the prefix.
		return "sdk." + NestedObjectDetails(obj, attr), "struct"
	default:
		panic(fmt.Sprintf("object %q: attribute %q has unknown type %q", obj.Name, attr.Name, attr.Type))
	}
}

// NestedObjectDetails is the struct name of a column that is an object, in the details struct:
// the object Name, the column key in camel case, and Details. Its fields are the Children of the column.
// CatalogIntegrationOpenCatalog with rest_config is CatalogIntegrationOpenCatalogRestConfigDetails.
func NestedObjectDetails(obj Object, attr Attribute) string {
	return obj.Name + genhelpers.SnakeCaseToCamel(attr.OutputKey()) + "Details"
}

// IsOptionalColumn reports whether the column is a pointer in the details: it is OptionalInShow, or it is an object
// that is not required.
func (a Attribute) IsOptionalColumn() bool {
	return a.OptionalInShow || (a.Type == AttrObject && !a.Required)
}

// NestedDetailsElem is the element type name of a nested list or set column in the details struct:
// the object Name, the column key in camel case without the plural s, and Details.
// ExternalVolume with storage_locations is ExternalVolumeStorageLocationDetails.
func NestedDetailsElem(obj Object, attr *Attribute) string {
	if attr.Type != AttrList && attr.Type != AttrSet {
		panic(fmt.Sprintf("object %q: nested column %q is not a list or set", obj.Name, attr.Name))
	}
	return obj.Name + strings.TrimSuffix(genhelpers.SnakeCaseToCamel(attr.OutputKey()), "s") + "Details"
}

// ColumnFieldName is the Go field name of a column. It comes from OutputKey.
func ColumnFieldName(attr Attribute) string {
	return genhelpers.SnakeCaseToCamel(attr.OutputKey())
}

// ColumnField is the field ExtractStructDetails would report for one column.
// The name comes from OutputKey. OptionalInShow prefixes "*" on both type strings.
func ColumnField(obj Object, attr Attribute) genhelpers.Field {
	concrete, underlying := AttributeGoTypes(obj, attr)
	if attr.IsOptionalColumn() {
		concrete = "*" + concrete
		underlying = "*" + underlying
	}
	return genhelpers.Field{
		Name:           ColumnFieldName(attr),
		ConcreteType:   concrete,
		UnderlyingType: underlying,
	}
}
