package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CortexAgentAssert) HasCreatedOnNotEmpty() *CortexAgentAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgent) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return c
}
