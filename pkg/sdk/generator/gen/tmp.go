package gen

// RemoveValidations removes all validations from the QueryStruct and all nested fields recursively.
// Currently used, to prevent generation of each validation in reused struct like file formats.
// Will be removed when the support for reusable structs is added.
func RemoveValidations(qs *QueryStruct) *QueryStruct {
	qs.validations = nil
	for _, f := range qs.fields {
		removeFieldValidations(f)
	}
	return qs
}

func removeFieldValidations(f *Field) {
	f.Validations = nil
	for i := range f.Fields {
		removeFieldValidations(&f.Fields[i])
	}
}
