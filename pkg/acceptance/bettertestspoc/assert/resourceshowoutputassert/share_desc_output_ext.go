package resourceshowoutputassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// ShareDescribeOutputRow asserts a single row of describe_output. DESCRIBE SHARE returns a row for every object granted to
// the share (and for the system objects that Snowflake adds to the shared database), so the number of rows is not fixed.
func ShareDescribeOutputRow(t *testing.T, name string, rowIndex int) *ShareDescribeOutputAssert {
	t.Helper()

	return &ShareDescribeOutputAssert{
		ResourceAssert: assert.NewResourceDescribeOutputAssertAtRow(name, rowIndex),
	}
}

func SharesDatasourceDescribeOutputRow(t *testing.T, name string, shareIdx int, rowIndex int) *ShareDescribeOutputAssert {
	t.Helper()

	return &ShareDescribeOutputAssert{
		ResourceAssert: assert.NewResourceShowOutputAssertAtRowWithPath(name, fmt.Sprintf("shares.%d.describe_output", shareIdx), rowIndex),
	}
}
