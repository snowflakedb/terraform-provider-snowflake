package gen

// Name adds identifier with field name "name" and type will be inferred from interface definition
func (v *QueryStruct) Name() *QueryStruct {
	identifier := NewField("name", "<will be replaced>", Tags().Identifier(), IdentifierOptions().Required())
	v.identifierField = identifier
	v.fields = append(v.fields, identifier)
	return v
}

func (v *QueryStruct) Identifier(fieldName string, kind string, transformer *IdentifierTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(fieldName, kind, Tags().Identifier(), transformer))
	return v
}

func (v *QueryStruct) OptionalIdentifier(name string, kind string, transformer *IdentifierTransformer) *QueryStruct {
	if len(kind) > 0 && kind[0] != '*' {
		kind = KindOfPointer(kind)
	}
	v.fields = append(v.fields, NewField(name, kind, Tags().Identifier(), transformer))
	return v
}
