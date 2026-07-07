package resourceshowoutputassert

func (e *ExternalVolumeShowOutputAssert) HasCommentEmpty() *ExternalVolumeShowOutputAssert {
	e.StringValueSet("comment", "")
	return e
}
