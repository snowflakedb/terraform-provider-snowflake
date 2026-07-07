package sdk

// NewAlertConditionFromString creates a single-element AlertConditionRequest slice from a SQL condition string.
func NewAlertConditionFromString(condition string) []AlertConditionRequest {
	return []AlertConditionRequest{*NewAlertConditionRequest().WithCondition([]string{condition})}
}

func (r *CreateAlertRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
