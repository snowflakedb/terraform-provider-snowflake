package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (i *IcebergTableDetailsAssert) HasTypeEmpty() *IcebergTableDetailsAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.IcebergTableDetails) error {
		t.Helper()
		if o.Type != nil {
			return fmt.Errorf("expected type to be empty")
		}
		return nil
	})
	return i
}
