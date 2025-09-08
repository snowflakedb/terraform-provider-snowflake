package invokeactionassert

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO [SNOW-1501905]: generalize this type of assertion
type nonExistenceCheck struct {
	id sdk.SchemaObjectIdentifier
}

func (w *nonExistenceCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		_, err := testClient.View.Show(t, w.id)
		if err == nil {
			return errors.New("expected err got nil")
		}
		if !strings.Contains(err.Error(), "object does not exist") {
			return fmt.Errorf("expected `object does not exist` is missing in %w", err)
		}
		return nil
	}
}

func ViewDoesNotExist(t *testing.T, id sdk.SchemaObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return &nonExistenceCheck{id: id}
}
