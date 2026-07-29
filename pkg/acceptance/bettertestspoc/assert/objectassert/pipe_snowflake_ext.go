package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (p *PipeAssert) HasNotEmptyCreatedOn() *PipeAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.Pipe) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return p
}
