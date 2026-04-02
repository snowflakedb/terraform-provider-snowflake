package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (n *NotebookShowOutputAssert) HasCreatedOnNotEmpty() *NotebookShowOutputAssert {
	n.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return n
}

func NotebooksDatasourceShowOutput(t *testing.T, name string) *NotebookShowOutputAssert {
	t.Helper()

	n := NotebookShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "notebooks.0."),
	}
	n.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &n
}
