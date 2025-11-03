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

	// TODO [SNOW-1501905]: test client passed here temporarily to be able to check secondary (by default our assertions use the default one)
	testClient *helpers.TestClient
}

func (s *snowflakeParameterUpdate) ToTerraformTestCheckFunc(t *testing.T, _ *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if s.testClient == nil {
			return errors.New("testClient must not be nil")
		}

		revertParameter := s.testClient.Parameter.UpdateAccountParameterTemporarily(t, s.parameter, s.newValue)
		t.Cleanup(revertParameter)

		return nil
	}
}

func UpdateAccountParameterTemporarily(t *testing.T, parameter sdk.AccountParameter, newValue string, testClient *helpers.TestClient) assert.TestCheckFuncProvider {
	t.Helper()
	return &snowflakeParameterUpdate{
		parameter:  parameter,
		newValue:   newValue,
		testClient: testClient,
	}
}
