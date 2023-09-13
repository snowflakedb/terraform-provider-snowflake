package generator

func (f *Field) Assignment(name string, kind string, transformer *ParameterTransformer) *Field {
	f.Fields = append(f.Fields, NewField(name, kind, Tags().Parameter(), transformer))
	return f
}

func (f *Field) TextAssignment(sqlPrefix string, transformer *ParameterTransformer) *Field {
	return f.Assignment(sqlToFieldName(sqlPrefix, true), "*string", transformer)
}

func (f *Field) OptionalTextAssignment(sqlPrefix string, transformer *ParameterTransformer) *Field {
	return f.Assignment(sqlToFieldName(sqlPrefix, true), "*string", transformer)
}
