package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (a *ImageRepositoryAssert) HasCreatedOnNotEmpty() *ImageRepositoryAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *ImageRepositoryAssert) HasRepositoryUrlNotEmpty() *ImageRepositoryAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.RepositoryUrl == "" {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}
