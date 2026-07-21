package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionAssert) HasConnectionUrlNotEmpty() *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.ConnectionUrl == "" {
			return fmt.Errorf("expected connection url not empty, got: %s", o.ConnectionUrl)
		}
		return nil
	})

	return c
}
