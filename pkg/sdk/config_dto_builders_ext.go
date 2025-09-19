package sdk

func ConfigFileWithDefaultProfile(c *ConfigDTO) *ConfigFile {
	return ConfigFileWithProfile(c, "default")
}

func ConfigFileWithProfile(c *ConfigDTO, profile string) *ConfigFile {
	return NewConfigFile().WithProfiles(map[string]ConfigDTO{
		profile: *c,
	})
}

func ConfigForSnowflakeAuth(accountIdentifier AccountIdentifier, userId AccountObjectIdentifier, pass string, roleId AccountObjectIdentifier, warehouseId AccountObjectIdentifier) *ConfigDTO {
	return NewConfigDTO().
		WithOrganizationName(accountIdentifier.OrganizationName()).
		WithAccountName(accountIdentifier.AccountName()).
		WithUser(userId.Name()).
		WithPassword(pass).
		WithAuthenticator(string(AuthenticationTypeSnowflake)).
		WithRole(roleId.Name()).
		WithWarehouse(warehouseId.Name())
}

func (c *ConfigDTO) WithAuthenticatorNil() *ConfigDTO {
	c.Authenticator = nil
	return c
}
