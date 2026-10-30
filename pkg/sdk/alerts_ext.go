package sdk

// NewAlertConditionFromString builds the IF (EXISTS (...)) block from a SQL condition string.
func NewAlertConditionFromString(condition string) AlertConditionRequest {
	return *NewAlertConditionRequest().WithCondition([]string{condition})
}

func (r *CreateAlertRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
