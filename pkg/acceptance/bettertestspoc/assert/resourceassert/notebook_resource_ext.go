package resourceassert

func (n *NotebookResourceAssert) HasFromPathAndStage(expectedPath string, expectedStageId string) *NotebookResourceAssert {
	n.ValueSet("from.0.path", expectedPath)
	n.ValueSet("from.0.stage", expectedStageId)
	return n
}
