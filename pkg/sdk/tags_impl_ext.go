package sdk

func createAllowedValues(values []string) *AllowedValues {
	items := make([]AllowedValue, 0, len(values))
	for _, value := range values {
		items = append(items, AllowedValue{
			Value: value,
		})
	}
	return &AllowedValues{
		Values: items,
	}
}

func createTagMaskingPolicies(maskingPolicies []SchemaObjectIdentifier) []TagMaskingPolicy {
	items := make([]TagMaskingPolicy, 0, len(maskingPolicies))
	for _, value := range maskingPolicies {
		items = append(items, TagMaskingPolicy{
			Name: value,
		})
	}
	return items
}
