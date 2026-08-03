package objectassert

func (e *ExternalAccessIntegrationAssert) HasTypeExternalAccess() *ExternalAccessIntegrationAssert {
	return e.HasType("EXTERNAL_ACCESS")
}

func (e *ExternalAccessIntegrationAssert) HasCategoryExternalAccess() *ExternalAccessIntegrationAssert {
	return e.HasCategory("EXTERNAL_ACCESS")
}
