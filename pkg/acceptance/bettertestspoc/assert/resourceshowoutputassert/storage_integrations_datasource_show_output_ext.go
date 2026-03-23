package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func StorageIntegrationsDatasourceShowOutput(t *testing.T, datasourceReference string) *StorageIntegrationShowOutputAssert {
	t.Helper()
	s := StorageIntegrationShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "storage_integrations.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
