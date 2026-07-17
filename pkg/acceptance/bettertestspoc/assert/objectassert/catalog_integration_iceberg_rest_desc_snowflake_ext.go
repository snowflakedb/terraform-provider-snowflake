package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogIntegrationIcebergRestDetailsAssert) HasBearerRestAuthentication() *CatalogIntegrationIcebergRestDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CatalogIntegrationIcebergRestDetails) error {
		if o.BearerRestAuthentication == nil {
			return fmt.Errorf("expected bearer rest authentication to have value; got: nil")
		}
		return nil
	})
	return c
}
