package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *TableConstraintDetailsAssert) HasNoComment() *TableConstraintDetailsAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.TableConstraintDetails) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty; got: %v", *o.Comment)
		}
		return nil
	})
	return t
}
