//go:build !account_level_tests

package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func assertThat(t *testing.T, fs ...assert.TestCheckFuncProvider) resource.TestCheckFunc {
	t.Helper()
	return assert.AssertThat(t, acc.TestClient(), fs...)
}
