package resourceshowoutputassert

func (h *HybridTableShowOutputAssert) HasCreatedOnNotEmpty() *HybridTableShowOutputAssert {
	h.ValuePresent("created_on")
	return h
}
