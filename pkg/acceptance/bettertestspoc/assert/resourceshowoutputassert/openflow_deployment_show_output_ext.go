package resourceshowoutputassert

func (o *OpenflowDeploymentShowOutputAssert) HasCreatedOnNotEmpty() *OpenflowDeploymentShowOutputAssert {
	o.ValuePresent("created_on")
	return o
}

func (o *OpenflowDeploymentShowOutputAssert) HasUpdatedOnNotEmpty() *OpenflowDeploymentShowOutputAssert {
	o.ValuePresent("updated_on")
	return o
}
