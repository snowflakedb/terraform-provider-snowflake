package generator

type dbStruct struct {
	name   string
	fields []dbField
}

type dbField struct {
	name string
	kind string
}

func DbStruct(name string) *dbStruct {
	return &dbStruct{
		name:   name,
		fields: make([]dbField, 0),
	}
}

func (v *dbStruct) Field(dbName string, kind string) *dbStruct {
	v.fields = append(v.fields, dbField{
		name: dbName,
		kind: kind,
	})
	return v
}

func (v *dbStruct) Text(dbName string) *dbStruct {
	return v.Field(dbName, "string")
}

func (v *dbStruct) Time(dbName string) *dbStruct {
	return v.Field(dbName, "time.Time")
}

func (v *dbStruct) OptionalTime(dbName string) *dbStruct {
	return v.Field(dbName, "sql.NullTime")
}

func (v *dbStruct) OptionalText(dbName string) *dbStruct {
	return v.Field(dbName, "sql.NullString")
}

func (v *dbStruct) Bool(dbName string) *dbStruct {
	return v.Field(dbName, "bool")
}

func (v *dbStruct) OptionalBool(dbName string) *dbStruct {
	return v.Field(dbName, "sql.NullBool")
}

func (v *dbStruct) Number(dbName string) *dbStruct {
	return v.Field(dbName, "int")
}

func (v *dbStruct) OptionalNumber(dbName string) *dbStruct {
	return v.Field(dbName, "sql.NullInt64")
}

func (v *dbStruct) IntoField() *Field {
	f := NewField(v.name, v.name, nil, nil)
	for _, field := range v.fields {
		f.withField(NewField(sqlToFieldName(field.name, true), field.kind, Tags().DB(field.name), nil))
	}
	return f
}
