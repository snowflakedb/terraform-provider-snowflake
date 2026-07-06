package resourceshowoutputassert

func (s *StorageLifecyclePolicyShowOutputAssert) HasCreatedOnNotEmpty() *StorageLifecyclePolicyShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}
