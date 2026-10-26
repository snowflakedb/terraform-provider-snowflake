package gen

import "fmt"

// ShowProjection is a method that copies a field subset from Source onto Target.
// The target struct is not emitted by this projection.
type ShowProjection struct {
	Source string
	Method string
	Target string
	Fields []string
}

// ShowProjection records a projection method. It emits the method only.
// Generation panics when a target field is missing on the source.
func (i *Interface) ShowProjection(source *Plain, method string, target *Plain) *Interface {
	if source == nil || source.name == "" {
		panic("show projection source is empty")
	}
	if target == nil || target.name == "" {
		panic(fmt.Sprintf("show projection %s target is empty", method))
	}
	if method == "" {
		panic(fmt.Sprintf("show projection on %s has an empty method name", source.name))
	}
	sourceFields := make(map[string]struct{}, len(source.fields))
	for _, field := range source.fields {
		sourceFields[field.name] = struct{}{}
	}
	fields := make([]string, 0, len(target.fields))
	for _, field := range target.fields {
		if _, ok := sourceFields[field.name]; !ok {
			panic(fmt.Sprintf("projection %s.%s: field %q is missing on %s", source.name, method, field.name, source.name))
		}
		fields = append(fields, field.name)
	}
	i.ShowProjections = append(i.ShowProjections, &ShowProjection{
		Source: source.name,
		Method: method,
		Target: target.name,
		Fields: fields,
	})
	return i
}
