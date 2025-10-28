package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1501905]: generalize this type of assertion
type SecurityIntegrationNonExistenceCheck struct {
	id sdk.AccountObjectIdentifier
}

func (w *SecurityIntegrationNonExistenceCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		_, err := testClient.SecurityIntegration.Show(t, w.id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				return nil
			}
			return err
		}
		return fmt.Errorf("expected database %s to be missing, but it exists", w.id.FullyQualifiedName())
	}
}

func SecurityIntegrationDoesNotExist(t *testing.T, id sdk.AccountObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return &SecurityIntegrationNonExistenceCheck{id: id}
}
