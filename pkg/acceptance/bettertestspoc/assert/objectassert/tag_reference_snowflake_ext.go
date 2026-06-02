package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *TagReferenceAssert) HasNoColumnName() *TagReferenceAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.TagReference) error {
		t.Helper()
		if o.ColumnName != nil {
			return fmt.Errorf("expected column name to be nil; got: %v", *o.ColumnName)
		}
		return nil
	})
	return a
}
