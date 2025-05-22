package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

// UsersDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func UsersDatasourceShowOutput(t *testing.T, name string) *UserShowOutputAssert {
	t.Helper()

	u := UserShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "users.0."),
	}
	u.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &u
}
