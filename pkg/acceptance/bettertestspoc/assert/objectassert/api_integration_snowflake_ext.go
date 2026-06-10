package objectassert

func (a *ApiIntegrationAssert) HasApiTypeExternalApi() *ApiIntegrationAssert {
	return a.HasApiType("EXTERNAL_API")
}

func (a *ApiIntegrationAssert) HasCategoryApi() *ApiIntegrationAssert {
	return a.HasCategory("API")
}
