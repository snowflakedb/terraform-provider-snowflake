package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (o *OpenflowConnectorVersionAssert) HasCommentNotEmpty() *OpenflowConnectorVersionAssert {
	o.AddAssertion(func(t *testing.T, o *sdk.OpenflowConnectorVersion) error {
		t.Helper()
		if o.Comment == nil || *o.Comment == "" {
			return fmt.Errorf("expected comment to be not empty")
		}
		return nil
	})
	return o
}
