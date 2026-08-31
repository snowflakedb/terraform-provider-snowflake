package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowDeploymentDetailsAssert) HasKeyNotEmpty() *OpenflowDeploymentDetailsAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowDeploymentDetails) error {
		t.Helper()
		if o.Key == nil || *o.Key == "" {
			return fmt.Errorf("expected key to be not empty")
		}
		return nil
	})
	return o
}
