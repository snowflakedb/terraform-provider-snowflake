package gen

import (
	"fmt"
	"strings"
)

// requireConsistentStructShapes panics when two operations of one interface define a struct with the same name and a
// different shape. The generator emits a struct once by name, so the second definition would be dropped silently and its
// operation would build different SQL than the definition says.
func requireConsistentStructShapes(definition *Interface) {
	shapes := make(map[string]string)
	var visit func(f *Field)
	visit = func(f *Field) {
		if f == nil || f.IsShared || !f.IsStruct() {
			return
		}
		name := f.KindNoPtr()
		shape := structShape(f)
		if previous, ok := shapes[name]; ok && previous != shape {
			panic(fmt.Sprintf("%s: struct %s is defined twice with a different shape:\n%s\n%s", definition.NameSingular, name, previous, shape))
		}
		shapes[name] = shape
		for idx := range f.Fields {
			visit(&f.Fields[idx])
		}
	}
	for _, op := range definition.Operations {
		if op.OptsField != nil {
			for idx := range op.OptsField.Fields {
				visit(&op.OptsField.Fields[idx])
			}
		}
		for _, helper := range op.HelperStructs {
			visit(helper)
		}
	}
}

// structShape lists the fields of a struct with their kinds, tags and validations. A nested struct is a field here; its own
// shape is compared where it is visited.
func structShape(f *Field) string {
	lines := make([]string, 0, len(f.Fields)+len(f.Validations))
	for _, child := range f.Fields {
		lines = append(lines, fmt.Sprintf("%s %s `%s`", child.Name, child.Kind, child.TagsPrintable()))
	}
	for _, validation := range f.Validations {
		lines = append(lines, fmt.Sprintf("validation %v %v", validation.Type, validation.FieldNames))
	}
	return strings.Join(lines, "\n")
}
