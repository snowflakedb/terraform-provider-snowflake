package generator

type plainStruct struct {
	name   string
	fields []plainField
}

type plainField struct {
	name string
	kind string
}

func PlainStruct(name string) *plainStruct {
	return &plainStruct{
		name:   name,
		fields: make([]plainField, 0),
	}
}

func (v *plainStruct) Field(name string, kind string) *plainStruct {
	v.fields = append(v.fields, plainField{
		name: name,
		kind: kind,
	})
	return v
}

func (v *plainStruct) Text(name string) *plainStruct {
	return v.Field(name, "string")
}

func (v *plainStruct) Time(name string) *plainStruct {
	return v.Field(name, "time.Time")
}

func (v *plainStruct) OptionalTime(name string) *plainStruct {
	return v.Field(name, "*time.Time")
}

func (v *plainStruct) OptionalText(name string) *plainStruct {
	return v.Field(name, "*string")
}

func (v *plainStruct) Bool(name string) *plainStruct {
	return v.Field(name, "bool")
}

func (v *plainStruct) OptionalBool(name string) *plainStruct {
	return v.Field(name, "*bool")
}

func (v *plainStruct) Number(dbName string) *plainStruct {
	return v.Field(dbName, "int")
}

func (v *plainStruct) OptionalNumber(dbName string) *plainStruct {
	return v.Field(dbName, "*int")
}

func (v *plainStruct) IntoField() *Field {
	f := NewField(v.name, v.name, nil, nil)
	for _, field := range v.fields {
		f.withField(NewField(field.name, field.kind, nil, nil))
	}
	return f
}
