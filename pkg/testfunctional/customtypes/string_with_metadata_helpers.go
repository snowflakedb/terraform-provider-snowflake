package customtypes

func StringWithMetadataAttributeCreate(v StringWithMetadataValue, createField **string) {
	if !v.IsNull() {
		*createField = new(v.ValueString())
	}
}
