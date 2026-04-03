package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (d *DatabaseRoleAssert) HasNotEmptyCreatedOn() *DatabaseRoleAssert {
	d.AddAssertion(func(t *testing.T, o *sdk.DatabaseRole) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return d
}
