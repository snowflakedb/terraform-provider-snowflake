package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// SecretsDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func SecretsDatasourceShowOutput(t *testing.T, datasourceReference string) *SecretShowOutputAssert {
	t.Helper()

	s := SecretShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "secrets.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
