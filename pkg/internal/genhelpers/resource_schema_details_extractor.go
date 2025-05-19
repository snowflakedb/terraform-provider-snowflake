package genhelpers

import (
	"slices"

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
}

var preferredFrontOrdering = []string{"database", "schema", "name"}

// TODO: test
func ExtractResourceSchemaDetails(name string, resourceSchema map[string]*schema.Schema) ResourceSchemaDetails {
	orderedAttributeNames := make([]string, 0)
	for key := range resourceSchema {
		orderedAttributeNames = append(orderedAttributeNames, key)
	}
	slices.Sort(orderedAttributeNames)

	for _, v := range slices.Backward(preferredFrontOrdering) {
		if idx := slices.Index(orderedAttributeNames, v); idx != -1 {
			orderedAttributeNames = append(orderedAttributeNames[:idx], orderedAttributeNames[idx+1:]...)
			orderedAttributeNames = append([]string{v}, orderedAttributeNames...)
		}
	}

	attributes := make([]SchemaAttribute, 0)
	for _, k := range orderedAttributeNames {
		s := resourceSchema[k]
		subtype := getComplexAttributeSubType(s)
		attributes = append(attributes, SchemaAttribute{
			Name:             k,
			AttributeType:    s.Type,
			AttributeSubType: subtype,
			Required:         s.Required,
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
