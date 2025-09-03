//go:build account_level_tests

package invokeactionassert

import (
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1501905]: generalize this type of assertion
type snowflakeParameterUpdate struct {
	parameter sdk.AccountParameter
	newValue  string
}

func (s *snowflakeParameterUpdate) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}

		revertParameter := testClient.Parameter.UpdateAccountParameterTemporarily(t, s.parameter, s.newValue)
		t.Cleanup(revertParameter)

		return nil
	}
}

func UpdateAccountParameterTemporarily(t *testing.T, parameter sdk.AccountParameter, newValue string) assert.TestCheckFuncProvider {
	t.Helper()
	return &snowflakeParameterUpdate{
		parameter: parameter,
		newValue:  newValue,
	}
}
