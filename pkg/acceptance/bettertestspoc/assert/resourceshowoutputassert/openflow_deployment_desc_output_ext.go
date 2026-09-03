package resourceshowoutputassert

func (o *OpenflowDeploymentDescribeOutputAssert) HasKeyNotEmpty() *OpenflowDeploymentDescribeOutputAssert {
	o.ValuePresent("key")
	return o
}
