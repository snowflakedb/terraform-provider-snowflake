package invokeactionassert

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type sessionCurrentWarehouseCheck struct {
	client   func() *sdk.Client
	expected func() string
}

func (w *sessionCurrentWarehouseCheck) ToTerraformTestCheckFunc(t *testing.T, _ *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		current, err := w.client().ContextFunctions.CurrentWarehouse(context.Background())
		if err != nil {
			return err
		}
		expected := w.expected()
		if current.Value != expected {
			return fmt.Errorf("expected session's current warehouse to be %q, got %q", expected, current.Value)
		}
		return nil
	}
}

// SessionCurrentWarehouseEquals asserts that the session backing client() has expected() selected as its
// current warehouse (CURRENT_WAREHOUSE()). client and expected are both evaluated lazily at check time
// (not at construction time), so expected can reference a value captured earlier in the same test, e.g. in
// a TestStep's PreConfig, and client can reference a provider that is not configured yet when this is
// constructed.
func SessionCurrentWarehouseEquals(client func() *sdk.Client, expected func() string) assert.TestCheckFuncProvider {
	return &sessionCurrentWarehouseCheck{client: client, expected: expected}
}
