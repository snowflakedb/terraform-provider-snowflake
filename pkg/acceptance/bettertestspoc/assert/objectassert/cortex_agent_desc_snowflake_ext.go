package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CortexAgentDetailsAssert) HasCreatedOnNotEmpty() *CortexAgentDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CortexAgentDetails) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return c
}
