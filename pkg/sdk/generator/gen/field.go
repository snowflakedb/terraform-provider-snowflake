package gen

import (
	"fmt"
	"slices"
	"strings"
)

// Field defines properties of a single field or struct (by defining Fields)
type Field struct {
	// Parent allows to traverse fields hierarchy more easily, nil for root
	Parent *Field
	// Fields defines children, use for struct fields
	Fields []Field
	// Validations defines validations on given field level (e.g. oneOf for children)
	Validations []*Validation
	// Name is how field is called in parent struct
	Name string
	// Kind is fields type (e.g. string, *bool)
	Kind string
	// Tags should contain ddl and sql tags used for SQL generation
	Tags map[string][]string
	// Required is used to mark fields which are essential (it's used e.g. for DTO builders generation)
	Required bool
}

func NewField(name string, kind string, tagBuilder *TagBuilder, transformer FieldTransformer) *Field {
	var tags map[string][]string
	if tagBuilder != nil {
		tags = tagBuilder.Build()
	} else {
		tags = make(map[string][]string)
	}
	f := &Field{
		Name: name,
		Kind: kind,
		Tags: tags,
	}
	if !IsNil(transformer) {
		return transformer.Transform(f)
	}
	return f
}

func (f *Field) withField(field *Field) *Field {
	f.Fields = append(f.Fields, *field)
	return f
}

func (f *Field) withFields(fields ...*Field) *Field {
	for _, field := range fields {
		f.Fields = append(f.Fields, *field)
	}
	return f
}

func (f *Field) withValidations(validations ...*Validation) *Field {
	f.Validations = validations
	return f
}

// HasAnyValidationInSubtree checks if any validations are present from current field level downwards
func (f *Field) HasAnyValidationInSubtree() bool {
	if len(f.Validations) > 0 {
		return true
	}
	for _, f := range f.Fields {
		if f.HasAnyValidationInSubtree() {
			return true
		}
	}
	return false
}

// TagsPrintable defines how tags are printed in options structs, it ensures the same order of tags for every field
func (f *Field) TagsPrintable() string {
	tagNames := []string{"ddl", "sql", "db"}
	var tagParts []string
	for _, tagName := range tagNames {
		v, ok := f.Tags[tagName]
		if ok {
			tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, tagName, strings.Join(v, ",")))
		}
	}
	if len(tagParts) > 0 {
		return fmt.Sprintf("`%s`", strings.Join(tagParts, " "))
	}
	return ""
}

// KindNoPtr return field's Kind but without pointer and array
func (f *Field) KindNoPtr() string {
	kindWithoutPtr, _ := strings.CutPrefix(f.Kind, "*")
	kindWithoutPtr, _ = strings.CutPrefix(kindWithoutPtr, "[]")
	return kindWithoutPtr
}

// KindNoSlice return field's Kind but without array
func (f *Field) KindNoSlice() string {
	kindWithoutSlice, _ := strings.CutPrefix(f.Kind, "[]")
	return kindWithoutSlice
}

// IsStruct checks if field is a struct
func (f *Field) IsStruct() bool {
	return len(f.Fields) > 0
}

func (f *Field) IsPointer() bool {
	return strings.HasPrefix(f.Kind, "*")
}

func (f *Field) IsSlice() bool {
	return strings.HasPrefix(f.Kind, "[]")
}

// ShouldBeInDto checks if field is not some static SQL field which should not be interacted with by SDK user
// TODO: this is a very naive implementation, consider fixing it with DSL builder connection
func (f *Field) ShouldBeInDto() bool {
	return !slices.Contains(f.Tags["ddl"], "static")
}

// IsRoot checks if field is at the top of field hierarchy, basically it is true for Option structs
func (f *Field) IsRoot() bool {
	return f.Parent == nil
}

// Path returns the way through the tree to the top, with dot separator (e.g. .SomeField.SomeChild)
func (f *Field) Path() string {
	if f.IsRoot() {
		return ""
	} else {
		return fmt.Sprintf("%s.%s", f.Parent.Path(), f.Name)
	}
}

// PathWithRoot returns the way through the tree to the top, with dot separator and root included (e.g. Struct.SomeField.SomeChild)
func (f *Field) PathWithRoot() string {
	if f.IsRoot() {
		return f.Name
	} else {
		return fmt.Sprintf("%s.%s", f.Parent.PathWithRoot(), f.Name)
	}
}

// DtoKind returns what should be fields kind in generated DTO, because it may differ from Kind
func (f *Field) DtoKind() string {
	switch {
	case f.IsRoot():
		withoutSuffix, _ := strings.CutSuffix(f.Kind, "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	case f.IsStruct():
		return fmt.Sprintf("%sRequest", f.Kind)
	default:
		return f.Kind
	}
}

// DtoDecl returns how struct should be declared in generated DTO (e.g. definition is without a pointer)
func (f *Field) DtoDecl() string {
	switch {
	case f.Parent == nil:
		withoutSuffix, _ := strings.CutSuffix(f.KindNoPtr(), "Options")
		return fmt.Sprintf("%sRequest", withoutSuffix)
	case f.IsStruct():
		return fmt.Sprintf("%sRequest", f.KindNoPtr())
	default:
		return f.KindNoPtr()
	}
}
