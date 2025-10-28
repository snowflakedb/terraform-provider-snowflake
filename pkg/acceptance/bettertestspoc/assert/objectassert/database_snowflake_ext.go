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

type databaseNonExistenceCheck struct {
	id sdk.AccountObjectIdentifier
}

func (w *databaseNonExistenceCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		_, err := testClient.Database.Show(t, w.id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				return nil
			}
			return err
		}
		return fmt.Errorf("expected database %s to be missing, but it exists", w.id.FullyQualifiedName())
	}
}

func DatabaseDoesNotExist(t *testing.T, id sdk.AccountObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return &databaseNonExistenceCheck{id: id}
}
