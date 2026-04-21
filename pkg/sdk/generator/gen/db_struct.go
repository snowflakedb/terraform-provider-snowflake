package gen

type dbStruct struct {
	name   string
	fields []dbField
}

type dbField struct {
	dbColumn  string
	kind      string
	fieldName string
}

func DbStruct(name string) *dbStruct {
	return &dbStruct{
		name:   name,
		fields: make([]dbField, 0),
	}
}

func (v *dbStruct) Field(dbColumn string, kind string) *dbStruct {
	v.fields = append(v.fields, dbField{
		dbColumn:  dbColumn,
		kind:      kind,
		fieldName: sqlToFieldName(dbColumn, true),
	})
	return v
}

func (v *dbStruct) FieldWithName(dbColumn string, kind string, name string) *dbStruct {
	v.fields = append(v.fields, dbField{
		dbColumn:  dbColumn,
		kind:      kind,
		fieldName: name,
	})
	return v
}

func (v *dbStruct) Text(dbName string) *dbStruct {
	return v.Field(dbName, "string")
}

func (v *dbStruct) OptionalText(dbName string) *dbStruct {
	return v.Field(dbName, "sql.NullString")
}

func (v *dbStruct) Time(dbName string) *dbStruct {
	return v.Field(dbName, "time.Time")
}

func (v *dbStruct) OptionalTime(dbName string) *dbStruct {
	return v.Field(dbName, "sql.NullTime")
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
		f.withField(NewField(field.fieldName, field.kind, Tags().DB(field.dbColumn), nil))
	}
	return f
}
