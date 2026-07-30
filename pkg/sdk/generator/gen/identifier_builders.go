package gen

// InterfaceIdentifierKind is a sentinel kind value meaning "the interface's own identifier type".
const InterfaceIdentifierKind = "<interface identifier type>"

// InterfaceIdentifierPointerKind is a sentinel kind value meaning "pointer to the interface's own identifier type".
const InterfaceIdentifierPointerKind = "<interface identifier pointer type>"

// Name adds identifier with field name "name" and type inferred from the interface definition.
func (v *QueryStruct) Name() *QueryStruct {
	identifier := NewField("name", InterfaceIdentifierKind, Tags().Identifier(), IdentifierOptions().Required())
	v.fields = append(v.fields, identifier)
	return v
}

// RenameTo adds an optional "RenameTo" identifier field with SQL "RENAME TO".
// The type is inferred from the interface's identifier kind (pointer), same mechanism as Name().
func (v *QueryStruct) RenameTo() *QueryStruct {
	field := NewField("RenameTo", InterfaceIdentifierPointerKind, Tags().Identifier(), IdentifierOptions().SQL("RENAME TO"))
	v.fields = append(v.fields, field)
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
