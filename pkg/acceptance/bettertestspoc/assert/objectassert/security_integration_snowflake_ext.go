package objectassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"

func (s *SecurityIntegrationAssert) HasEnabledSnowflakeDefault() *SecurityIntegrationAssert {
	env := testenvs.GetSnowflakeEnvironmentWithProdDefault()
	enabled := false
	if env == testenvs.SnowflakeNonProdEnvironment {
		enabled = true
	}
	s.HasEnabled(enabled)
	return s
}
