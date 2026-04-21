package resourceshowoutputassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func CatalogIntegrationsDatasourceShowOutput(t *testing.T, datasourceReference string, index int) *CatalogIntegrationShowOutputAssert {
	t.Helper()

	s := CatalogIntegrationShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", fmt.Sprintf("catalog_integrations.%d.", index)),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
