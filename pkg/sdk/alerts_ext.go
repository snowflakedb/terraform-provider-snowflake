package sdk

// NewAlertConditionFromString creates a single-element AlertConditionRequest slice from a SQL condition string.
func NewAlertConditionFromString(condition string) []AlertConditionRequest {
	return []AlertConditionRequest{*NewAlertConditionRequest().WithCondition([]string{condition})}
}

func (r *CreateAlertRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// ID returns the schema identifier of the described object.
func (v *AlertDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}
