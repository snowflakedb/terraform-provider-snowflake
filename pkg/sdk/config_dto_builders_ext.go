package sdk

func ConfigFileWithDefaultProfile(c *ConfigDTO) *ConfigFile {
	return ConfigFileWithProfile(c, "default")
}

func ConfigFileWithProfile(c *ConfigDTO, profile string) *ConfigFile {
	return NewConfigFile().WithProfiles(map[string]ConfigDTO{
		profile: *c,
	})
}
