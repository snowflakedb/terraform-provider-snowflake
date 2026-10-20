package genhelpers

import (
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceSchemaDetails struct {
	Name       string
	Attributes []SchemaAttribute
}

func (s ResourceSchemaDetails) ObjectName() string {
	return s.Name
}

type SchemaAttribute struct {
	Name             string
	AttributeType    schema.ValueType
	AttributeSubType schema.ValueType
	Required         bool
	WriteOnly        bool
}

var preferredFrontOrdering = []string{"database", "schema", "name"}

// OrderByNameWithPreferredFront sorts items by name, then moves database, schema, and name to the front when present.
func OrderByNameWithPreferredFront[T any](items []T, name func(T) string) []T {
	slices.SortFunc(items, func(a, b T) int {
		return strings.Compare(name(a), name(b))
	})
	for _, front := range slices.Backward(preferredFrontOrdering) {
		idx := slices.IndexFunc(items, func(item T) bool { return name(item) == front })
		if idx == -1 {
			continue
		}
		item := items[idx]
		items = append(items[:idx], items[idx+1:]...)
		items = append([]T{item}, items...)
	}
	return items
}

// TODO: test
func ExtractResourceSchemaDetails(name string, resourceSchema map[string]*schema.Schema) ResourceSchemaDetails {
	orderedAttributeNames := make([]string, 0, len(resourceSchema))
	for key := range resourceSchema {
		orderedAttributeNames = append(orderedAttributeNames, key)
	}
	orderedAttributeNames = OrderByNameWithPreferredFront(orderedAttributeNames, func(name string) string { return name })

	attributes := make([]SchemaAttribute, 0, len(orderedAttributeNames))
	for _, k := range orderedAttributeNames {
		s := resourceSchema[k]
		subtype := getComplexAttributeSubType(s)
		attributes = append(attributes, SchemaAttribute{
			Name:             k,
			AttributeType:    s.Type,
			AttributeSubType: subtype,
			Required:         s.Required,
			WriteOnly:        s.WriteOnly,
		})
	}

	return ResourceSchemaDetails{
		Name:       name,
		Attributes: attributes,
	}
}

// getComplexAttributeSubType currently handles list/set of simple values and list/set of complex objects:
// - simple: Elem: &schema.Schema{Type: schema.TypeString};
// - complex: Elem: &schema.Resource{...};
func getComplexAttributeSubType(s *schema.Schema) schema.ValueType {
	if s.Type == schema.TypeList || s.Type == schema.TypeSet {
		switch v := s.Elem.(type) {
		case *schema.Schema:
			return v.Type
		case *schema.Resource:
			return schema.TypeMap
		default:
			return schema.TypeInvalid
		}
	}
	return s.Type
}

// TODO [SNOW-1501905]: create a type-safe type abstraction and reuse inside generators
func ResourceSchemaTypeToString(t schema.ValueType) (stringType string) {
	switch t {
	case schema.TypeInvalid:
		stringType = "Invalid"
	case schema.TypeBool:
		stringType = "Bool"
	case schema.TypeInt:
		stringType = "Int"
	case schema.TypeFloat:
		stringType = "Float"
	case schema.TypeString:
		stringType = "String"
	case schema.TypeList:
		stringType = "List"
	case schema.TypeMap:
		stringType = "Map"
	case schema.TypeSet:
		stringType = "Set"
	}
	return
}
