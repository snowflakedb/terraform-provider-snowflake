package gen

// Enum adds a required enum field using a generator-defined Enum.
func (v *QueryStruct) Enum(name string, enum *Enum, transformer FieldTransformer) *QueryStruct {
	return v.PredefinedQueryStructField(name, enum.Kind(), transformer)
}

// OptionalEnum adds an optional enum field using a generator-defined Enum.
func (v *QueryStruct) OptionalEnum(name string, enum *Enum, transformer FieldTransformer) *QueryStruct {
	return v.PredefinedQueryStructField(name, enum.KindPtr(), transformer)
}

// EnumLegacy returns a required Field for a non-generated enum T.
// Use with WithField: .WithField(g.EnumLegacy[sdkcommons.T]("Name", transformer))
func EnumLegacy[T any](name string, transformer FieldTransformer) *Field {
	return NewField(name, KindOfT[T](), Tags(), transformer)
}

// OptionalEnumLegacy returns an optional Field for a non-generated enum T.
// Use with WithField: .WithField(g.OptionalEnumLegacy[sdkcommons.T]("Name", transformer))
func OptionalEnumLegacy[T any](name string, transformer FieldTransformer) *Field {
	return NewField(name, KindOfTPointer[T](), Tags(), transformer)
}

// WithField appends a pre-built *Field to the QueryStruct.
func (v *QueryStruct) WithField(field *Field) *QueryStruct {
	v.fields = append(v.fields, field)
	return v
}
