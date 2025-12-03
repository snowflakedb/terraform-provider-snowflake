package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

func (s *SecurityIntegrationShowOutputAssert) HasCreatedOnNotEmpty() *SecurityIntegrationShowOutputAssert {
	s.AddAssertion(assert.ValuePresent("show_output.0.created_on"))
	return s
}

func (s *SecurityIntegrationShowOutputAssert) HasEnabledSnowflakeDefault() *SecurityIntegrationShowOutputAssert {
	env := testenvs.GetSnowflakeEnvironmentWithProdDefault()
	enabled := false
	if env == testenvs.SnowflakeNonProdEnvironment {
		enabled = true
	}
	s.AddAssertion(assert.ResourceShowOutputBoolValueSet("enabled", enabled))
	return s
}
