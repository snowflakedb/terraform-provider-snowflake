package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SecurityIntegrationAssert) HasCreatedOnNotEmpty() *SecurityIntegrationAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SecurityIntegration) error {
		t.Helper()
		if o.CreatedOn.IsZero() {
			return fmt.Errorf("expected created_on to be set, but it was zero")
		}
		return nil
	})
	return s
}
