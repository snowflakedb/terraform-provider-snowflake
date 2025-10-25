package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SecretAssert) HasCreatedOnNotEmpty() *SecretAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Secret) error {
		t.Helper()
		if o.CreatedOn.IsZero() {
			return fmt.Errorf("expected created_on to be not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return s
}

func (s *SecretAssert) HasCommentEmpty() *SecretAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Secret) error {
		t.Helper()
		if o.Comment != nil && *o.Comment != "" {
			return fmt.Errorf("expected comment to be empty; got: %v", *o.Comment)
		}
		return nil
	})
	return s
}
