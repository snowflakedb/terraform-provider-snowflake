package gen

import "strings"

// SetStateToValueFromConfigFields is the field list for setStateToValuesFromConfig.
// Names follow Attributes order. The helper copies string, int, and sets of string or
// int. A float attribute is stored as a string. An optional bool element is a
// string. comment is left out because SetConfigFields writes it from SHOW.
// A nested block or a list is left out. A create-only optional bool is left out:
// it is a string with IgnoreAfterCreation, and a native create-only bool is TypeBool.
func (m ResourceModel) SetStateToValueFromConfigFields() []string {
	names := make([]string, 0, len(m.Attributes))
	for _, attr := range m.Attributes {
		if setStateField(attr) {
			names = append(names, attr.Name)
		}
	}
	return names
}

func setStateField(attr ResourceAttribute) bool {
	if attr.Required || attr.Name == "comment" || len(attr.Children) > 0 || strings.Contains(attr.DiffSuppressExpr, "IgnoreAfterCreation") {
		return false
	}
	switch attr.TfType {
	case SchemaTypeString, SchemaTypeInt:
		return true
	case SchemaTypeSet:
		if attr.Elem == nil {
			return false
		}
		return attr.Elem.TfType == SchemaTypeString || attr.Elem.TfType == SchemaTypeInt
	default:
		return false
	}
}
