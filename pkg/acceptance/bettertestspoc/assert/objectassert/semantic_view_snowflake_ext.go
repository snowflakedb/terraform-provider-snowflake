package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SemanticViewAssert) HasCreatedOnNotEmpty() *SemanticViewAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SemanticView) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return s
}

func (s *SemanticViewAssert) HasNoComment() *SemanticViewAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SemanticView) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty")
		}
		return nil
	})
	return s
}
