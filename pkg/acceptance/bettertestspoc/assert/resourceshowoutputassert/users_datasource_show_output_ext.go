package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// UsersDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func UsersDatasourceShowOutput(t *testing.T, datasourceReference string) *UserShowOutputAssert {
	t.Helper()

	u := UserShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "users.0."),
	}
	u.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &u
}
