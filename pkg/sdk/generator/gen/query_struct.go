package gen

// TODO For Field abstractions use internal Field representation instead of copying only needed fields, e.g.
//
//	type QueryStruct struct {
//		internalRepresentation *Field
//		...additional fields that are not present in the Field
//	}
type QueryStruct struct {
	name            string
	fields          []*Field
	identifierField *Field
	validations     []*Validation
}

func NewQueryStruct(name string) *QueryStruct {
	return &QueryStruct{
		name:        name,
		fields:      make([]*Field, 0),
		validations: make([]*Validation, 0),
	}
}

func (v *QueryStruct) IntoField() *Field {
	return NewField(v.name, v.name, nil, nil).
		withFields(v.fields...).
		withValidations(v.validations...)
}

func (v *QueryStruct) WithValidation(validationType ValidationType, fieldNames ...string) *QueryStruct {
	v.validations = append(v.validations, NewValidation(validationType, fieldNames...))
	return v
}

func (v *QueryStruct) PredefinedQueryStructField(name string, kind string, transformer FieldTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, kind, Tags(), transformer))
	return v
}

func (v *QueryStruct) QueryStructField(name string, queryStruct *QueryStruct, transformer FieldTransformer) *QueryStruct {
	return v.queryStructField(name, queryStruct, "", transformer)
}

func (v *QueryStruct) ListQueryStructField(name string, queryStruct *QueryStruct, transformer FieldTransformer) *QueryStruct {
	return v.queryStructField(name, queryStruct, "[]", transformer)
}

func (v *QueryStruct) OptionalQueryStructField(name string, queryStruct *QueryStruct, transformer FieldTransformer) *QueryStruct {
	return v.queryStructField(name, queryStruct, "*", transformer)
}

func (v *QueryStruct) queryStructField(name string, queryStruct *QueryStruct, kindPrefix string, transformer FieldTransformer) *QueryStruct {
	qs := queryStruct.IntoField()
	qs.Name = name
	qs.Kind = kindPrefix + qs.Kind
	if transformer != nil {
		qs = transformer.Transform(qs)
	}
	v.fields = append(v.fields, qs)
	return v
}
