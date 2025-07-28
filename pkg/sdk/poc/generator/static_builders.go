package generator

func (v *QueryStruct) SQL(sql string) *QueryStruct {
	return v.SQLWithCustomFieldName(sqlToFieldName(sql, false), sql)
}

func (v *QueryStruct) SQLWithCustomFieldName(name, sql string) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "bool", Tags().Static().SQL(sql), nil))
	return v
}

func (v *QueryStruct) Create() *QueryStruct {
	return v.SQL("CREATE")
}

func (v *QueryStruct) CreateOrAlter() *QueryStruct {
	return v.SQL("CREATE OR ALTER")
}

func (v *QueryStruct) Alter() *QueryStruct {
	return v.SQL("ALTER")
}

func (v *QueryStruct) Drop() *QueryStruct {
	return v.SQL("DROP")
}

func (v *QueryStruct) Show() *QueryStruct {
	return v.SQL("SHOW")
}

func (v *QueryStruct) Describe() *QueryStruct {
	return v.SQL("DESCRIBE")
}

func (v *QueryStruct) Grant() *QueryStruct {
	return v.SQL("GRANT")
}

func (v *QueryStruct) Revoke() *QueryStruct {
	return v.SQL("REVOKE")
}
