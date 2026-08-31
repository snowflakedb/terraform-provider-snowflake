package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowDeploymentAssert) HasKeyNotEmpty() *OpenflowDeploymentAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowDeployment) error {
		t.Helper()
		if o.Key == nil || *o.Key == "" {
			return fmt.Errorf("expected key to be not empty")
		}
		return nil
	})
	return o
}
