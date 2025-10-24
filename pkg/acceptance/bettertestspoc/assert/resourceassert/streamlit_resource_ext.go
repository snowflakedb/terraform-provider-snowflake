package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// HasExternalAccessIntegrations checks that the external_access_integrations field contains the expected values
func (s *StreamlitResourceAssert) HasExternalAccessIntegrations(expected []string) *StreamlitResourceAssert {
	s.AddAssertion(assert.ValueSet("external_access_integrations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i, val := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("external_access_integrations.%d", i), val))
	}
	return s
}
