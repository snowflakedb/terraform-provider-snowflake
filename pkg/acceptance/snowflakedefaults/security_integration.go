package snowflakedefaults

func EnabledValueForSnowflakeOauthSecurityIntegration() bool {
	if getSnowflakeEnvironmentWithProdDefault() == SnowflakeNonProdEnvironment {
		return true
	}
	return false
}
