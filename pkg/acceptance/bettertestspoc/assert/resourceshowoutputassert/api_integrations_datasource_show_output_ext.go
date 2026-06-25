package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func ApiIntegrationsDatasourceShowOutput(t *testing.T, datasourceReference string) *ApiIntegrationShowOutputAssert {
	t.Helper()
	s := ApiIntegrationShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "api_integrations.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
