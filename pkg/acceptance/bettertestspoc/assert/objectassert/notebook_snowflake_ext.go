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

func (n *NotebookAssert) HasNoComment() *NotebookAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.Notebook) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty")
		}
		return nil
	})
	return n
}

func (n *NotebookAssert) HasNoQueryWarehouse() *NotebookAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.Notebook) error {
		t.Helper()
		if o.QueryWarehouse != nil {
			return fmt.Errorf("expected query_warehouse to be empty")
		}
		return nil
	})
	return n
}
