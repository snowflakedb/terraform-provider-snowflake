package sdk

func LegacyConfigFileWithDefaultProfile(c *LegacyConfigDTO) *LegacyConfigFile {
	return LegacyConfigFileWithProfile(c, "default")
}

func LegacyConfigFileWithProfile(c *LegacyConfigDTO, profile string) *LegacyConfigFile {
	return NewLegacyConfigFile().WithProfiles(map[string]LegacyConfigDTO{
		profile: *c,
	})
}
