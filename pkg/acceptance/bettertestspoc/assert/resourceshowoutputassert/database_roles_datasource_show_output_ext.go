package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func DatabaseRolesDatasourceShowOutput(t *testing.T, datasourceReference string) *DatabaseRoleShowOutputAssert {
	t.Helper()

	s := DatabaseRoleShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "database_roles.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
