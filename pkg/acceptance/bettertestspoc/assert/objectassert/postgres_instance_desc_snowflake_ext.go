package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (p *PostgresInstanceDetailsAssert) HasPostgresVersionNotEmpty() *PostgresInstanceDetailsAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.PostgresInstanceDetails) error {
		t.Helper()
		if o.PostgresVersion == 0 {
			return fmt.Errorf("expected postgres_version to be not empty")
		}
		return nil
	})
	return p
}
