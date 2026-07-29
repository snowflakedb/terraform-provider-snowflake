package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *NotebookAssert) HasCreatedOnNotEmpty() *NotebookAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Notebook) error {
		t.Helper()
		if o.CreatedOn.Equal((time.Time{})) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (s *NotebookAssert) HasNoQueryWarehouse() *NotebookAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Notebook) error {
		t.Helper()
		if o.QueryWarehouse != nil {
			return fmt.Errorf("expected query_warehouse to be empty")
		}
		return nil
	})
	return s
}
