package resourceparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// UsersDatasourceParameters is a temporary workaround to have better parameter assertions in data source acceptance tests.
func UsersDatasourceParameters(t *testing.T, datasourceReference string) *UserResourceParametersAssert {
	t.Helper()

	u := UserResourceParametersAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "parameters", "users.0."),
	}
	u.AddAssertion(assert.ValueSet("parameters.#", "1"))
	return &u
}
