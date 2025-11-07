package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// ListingsDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func ListingsDatasourceShowOutput(t *testing.T, name string) *ListingShowOutputAssert {
	t.Helper()

	l := ListingShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "listings.0."),
	}
	l.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &l
}
