package gen

type dbStruct struct {
	name   string
	fields []dbField
}

type dbField struct {
	name string
	kind string
	// goName is an optional Go field name override. When empty, sqlToFieldName(name, true) is used.
	goName string
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

// FieldWithGoName adds a field where the Go identifier in the generated struct is explicitly
// set to goName instead of being derived from dbName via sqlToFieldName. The db: tag still
// uses dbName as the column name.
func (v *dbStruct) FieldWithGoName(dbName string, goName string, kind string) *dbStruct {
	v.fields = append(v.fields, dbField{
		name:   dbName,
		kind:   kind,
		goName: goName,
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
		goFieldName := field.goName
		if goFieldName == "" {
			goFieldName = sqlToFieldName(field.name, true)
		}
		f.withField(NewField(goFieldName, field.kind, Tags().DB(field.name), nil))
	}
	return f
}
