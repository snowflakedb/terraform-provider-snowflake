package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (n *NotebookAssert) HasCreatedOnNotEmpty() *NotebookAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.Notebook) error {
		t.Helper()
		if o.CreatedOn.Equal((time.Time{})) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return n
}
