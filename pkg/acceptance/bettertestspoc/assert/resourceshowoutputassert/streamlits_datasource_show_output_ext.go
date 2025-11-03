package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func StreamlitsDatasourceShowOutput(t *testing.T, datasourceReference string) *StreamlitShowOutputAssert {
	t.Helper()

	s := StreamlitShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "streamlits.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
