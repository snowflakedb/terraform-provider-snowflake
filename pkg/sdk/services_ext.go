package sdk

func (s *CreateServiceRequest) GetName() SchemaObjectIdentifier {
	return s.name
}

func (s *ExecuteJobServiceRequest) GetName() SchemaObjectIdentifier {
	return s.Name
}
