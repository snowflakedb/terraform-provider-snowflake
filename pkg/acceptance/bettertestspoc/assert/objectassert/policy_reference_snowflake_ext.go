package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (p *PolicyReferenceAssert) HasNoRefColumnName() *PolicyReferenceAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PolicyReference) error {
		t.Helper()
		if o.RefColumnName != nil {
			return fmt.Errorf("expected ref column name to be nil; got: %s", *o.RefColumnName)
		}
		return nil
	})
	return p
}
