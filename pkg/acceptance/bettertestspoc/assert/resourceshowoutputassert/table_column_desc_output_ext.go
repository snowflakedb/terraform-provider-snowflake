package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (t *TableColumnDescribeOutputAssert) HasType(expected datatypes.DataType) *TableColumnDescribeOutputAssert {
	t.StringValueSet("type", expected.ToSql())
	return t
}
