package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowConnectorDefinitionAssert) HasCategoriesNotEmpty() *OpenflowConnectorDefinitionAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorDefinition) error {
		t.Helper()
		if len(o.Categories) == 0 {
			return fmt.Errorf("expected categories to be not empty")
		}
		return nil
	})
	return o
}

// HasMaxNodeCountPositive covers both a missing column and one that parsed to zero, either of which a wrong
// db tag produces.
func (o *OpenflowConnectorDefinitionAssert) HasMaxNodeCountPositive() *OpenflowConnectorDefinitionAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorDefinition) error {
		t.Helper()
		if o.MaxNodeCount == nil {
			return fmt.Errorf("expected max node count to have value; got: nil")
		}
		if *o.MaxNodeCount <= 0 {
			return fmt.Errorf("expected max node count to be positive; got: %d", *o.MaxNodeCount)
		}
		return nil
	})
	return o
}
