package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *StreamlitResourceAssert) HasExternalAccessIntegrations(expected []string) *StreamlitResourceAssert {
	s.AddAssertion(assert.ValueSet("external_access_integrations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, val := range expected {
		s.AddAssertion(assert.SetElem("external_access_integrations", val))
	}
	return s
}
