package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func SecurityIntegrationsDatasourceShowOutput(t *testing.T, datasourceReference string) *SecurityIntegrationShowOutputAssert {
	t.Helper()
	s := SecurityIntegrationShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "security_integrations.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
