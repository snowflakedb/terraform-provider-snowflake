package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (f *FileFormatAssert) HasCreatedOnNotEmpty() *FileFormatAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FileFormat) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return f
}
