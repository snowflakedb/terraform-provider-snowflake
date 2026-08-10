package resourceshowoutputassert

func (e *ExternalAccessIntegrationShowOutputAssert) HasCreatedOnNotEmpty() *ExternalAccessIntegrationShowOutputAssert {
	e.ValuePresent("created_on")
	return e
}
