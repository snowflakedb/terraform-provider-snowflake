package generator

func (v *QueryStruct) assignment(name string, kind string, transformer *ParameterTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, kind, Tags().Parameter(), transformer))
	return v
}

func (v *QueryStruct) Assignment(sqlPrefix string, kind string, transformer *ParameterTransformer) *QueryStruct {
	if transformer != nil {
		transformer = transformer.SQL(sqlPrefix)
	} else {
		transformer = ParameterOptions().SQL(sqlPrefix)
	}
	return v.assignment(sqlToFieldName(sqlPrefix, true), kind, transformer)
}

func (v *QueryStruct) OptionalAssignment(sqlPrefix string, kind string, transformer *ParameterTransformer) *QueryStruct {
	if len(kind) > 0 && kind[0] != '*' {
		kind = KindOfPointer(kind)
	}
	return v.Assignment(sqlPrefix, kind, transformer)
}

func (v *QueryStruct) ListAssignment(sqlPrefix string, listItemKind string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, KindOfSlice(listItemKind), transformer)
}

func (v *QueryStruct) OptionalListAssignment(sqlPrefix string, listItemKind string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, "*[]"+listItemKind, transformer)
}

func (v *QueryStruct) NumberAssignment(sqlPrefix string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, "int", transformer)
}

func (v *QueryStruct) OptionalNumberAssignment(sqlPrefix string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, "*int", transformer)
}

func (v *QueryStruct) TextAssignment(sqlPrefix string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, "string", transformer.Required())
}

func (v *QueryStruct) OptionalTextAssignment(sqlPrefix string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, "*string", transformer)
}

func (v *QueryStruct) BooleanAssignment(sqlPrefix string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, "bool", transformer)
}

func (v *QueryStruct) OptionalBooleanAssignment(sqlPrefix string, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, "*bool", transformer)
}

func (v *QueryStruct) OptionalIdentifierAssignment(sqlPrefix string, identifierKind string, transformer *ParameterTransformer) *QueryStruct {
	return v.OptionalAssignment(sqlPrefix, identifierKind, transformer)
}

func (v *QueryStruct) OptionalComment() *QueryStruct {
	return v.OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes())
}

func (v *QueryStruct) SetComment() *QueryStruct {
	return v.OptionalTextAssignment("SET COMMENT", ParameterOptions().SingleQuotes())
}
