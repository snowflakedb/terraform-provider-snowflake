package sdk

func (s *CreateSemanticViewRequest) GetName() SchemaObjectIdentifier {
	return s.name
}

func NewSemanticViewDetails(
	ObjectKind string,
	ObjectName string,
	ParentEntity *string,
	Property string,
	PropertyValue string,
) SemanticViewDetails {
	details := SemanticViewDetails{
		ObjectKind:    ObjectKind,
		ObjectName:    ObjectName,
		Property:      Property,
		PropertyValue: PropertyValue,
	}
	if ParentEntity != nil {
		details.ParentEntity = ParentEntity
	}

	return details
}
