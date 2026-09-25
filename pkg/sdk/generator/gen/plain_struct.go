package gen

import "fmt"

// Plain is the builder returned by PlainStruct. Defs helpers take *Plain the same way QueryStruct helpers take *QueryStruct.
type Plain struct {
	name   string
	fields []plainField
}

type plainField struct {
	name string
	kind string
}

func PlainStruct(name string) *Plain {
	return &Plain{
		name:   name,
		fields: make([]plainField, 0),
	}
}

func (v *Plain) Field(name string, kind string) *Plain {
	v.fields = append(v.fields, plainField{
		name: name,
		kind: kind,
	})
	return v
}

func (v *Plain) OptionalField(name string, kind string) *Plain {
	return v.Field(name, fmt.Sprintf("*%s", kind))
}

func (v *Plain) Text(name string) *Plain {
	return v.Field(name, "string")
}

func (v *Plain) OptionalText(name string) *Plain {
	return v.Field(name, "*string")
}

func (v *Plain) Time(name string) *Plain {
	return v.Field(name, "time.Time")
}

func (v *Plain) OptionalTime(name string) *Plain {
	return v.Field(name, "*time.Time")
}

func (v *Plain) Bool(name string) *Plain {
	return v.Field(name, "bool")
}

func (v *Plain) OptionalBool(name string) *Plain {
	return v.Field(name, "*bool")
}

func (v *Plain) Number(dbName string) *Plain {
	return v.Field(dbName, "int")
}

func (v *Plain) OptionalNumber(dbName string) *Plain {
	return v.Field(dbName, "*int")
}

func (v *Plain) StringList(dbName string) *Plain {
	return v.Field(dbName, "[]string")
}

func (v *Plain) AccountObjectIdentifier() *Plain {
	return v.Field("Id", "AccountObjectIdentifier")
}

func (v *Plain) SchemaObjectIdentifier() *Plain {
	return v.Field("Id", "SchemaObjectIdentifier")
}

func (v *Plain) IntoField() *Field {
	f := NewField(v.name, v.name, nil, nil)
	for _, field := range v.fields {
		f.withField(NewField(field.name, field.kind, nil, nil))
	}
	return f
}
