package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func SchemasDatasourceShowOutput(t *testing.T, datasourceReference string) *SchemaShowOutputAssert {
	t.Helper()

	s := SchemaShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "schemas.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
