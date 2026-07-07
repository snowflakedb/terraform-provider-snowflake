package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (i *IcebergTableDescribeOutputAssert) HasType(expected datatypes.DataType) *IcebergTableDescribeOutputAssert {
	i.StringValueSet("type", expected.ToSql())
	return i
}
