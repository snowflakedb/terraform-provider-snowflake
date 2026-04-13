package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (e *ExternalVolumeShowOutputAssert) HasCommentEmpty() *ExternalVolumeShowOutputAssert {
	e.AddAssertion(assert.ResourceShowOutputValueSet("comment", ""))
	return e
}

func ExternalVolumesDatasourceShowOutput(t *testing.T, datasourceReference string) *ExternalVolumeShowOutputAssert {
	t.Helper()
	s := ExternalVolumeShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "external_volumes.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
