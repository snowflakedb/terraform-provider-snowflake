package resourceshowoutputassert

func (n *NotebookShowOutputAssert) HasCreatedOnNotEmpty() *NotebookShowOutputAssert {
	n.ValuePresent("created_on")
	return n
}
