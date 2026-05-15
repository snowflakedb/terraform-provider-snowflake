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

// EnumAssignment adds a required enum parameter assignment.
func (v *QueryStruct) EnumAssignment(sqlPrefix string, enum *Enum, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, enum.Kind(), transformer)
}

// OptionalEnumAssignment adds an optional enum parameter assignment (pointer kind).
func (v *QueryStruct) OptionalEnumAssignment(sqlPrefix string, enum *Enum, transformer *ParameterTransformer) *QueryStruct {
	return v.Assignment(sqlPrefix, enum.KindPtr(), transformer)
}

// EnumAssignmentWithFieldName adds a required enum parameter assignment with a custom Go field name.
func (v *QueryStruct) EnumAssignmentWithFieldName(sqlPrefix string, enum *Enum, transformer *ParameterTransformer, fieldName string) *QueryStruct {
	return v.AssignmentWithFieldName(sqlPrefix, enum.Kind(), transformer, fieldName)
}

// Enum adds a required enum field to a plainStruct.
func (v *plainStruct) Enum(name string, enum *Enum) *plainStruct {
	return v.Field(name, enum.Kind())
}
