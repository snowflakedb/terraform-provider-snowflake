package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func TableColumnDescribeOutputOnIdx(t *testing.T, name string, idx int) *TableColumnDescribeOutputAssert {
	t.Helper()
	return &TableColumnDescribeOutputAssert{
		ResourceAssert: assert.NewResourceDescribeOutputAssertAtRow(name, idx),
	}
}

func (t *TableColumnDescribeOutputAssert) HasType(expected datatypes.DataType) *TableColumnDescribeOutputAssert {
	t.StringValueSet("type", expected.ToSql())
	return t
}
