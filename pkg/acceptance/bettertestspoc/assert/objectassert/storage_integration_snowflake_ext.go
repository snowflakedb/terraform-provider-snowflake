package objectassert

func (s *StorageIntegrationAssert) HasStorageTypeExternal() *StorageIntegrationAssert {
	return s.HasStorageType("EXTERNAL_STAGE")
}

func (s *StorageIntegrationAssert) HasCategoryStorage() *StorageIntegrationAssert {
	return s.HasCategory("STORAGE")
}
