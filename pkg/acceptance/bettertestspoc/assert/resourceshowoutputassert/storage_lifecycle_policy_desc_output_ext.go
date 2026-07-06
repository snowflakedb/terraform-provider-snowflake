package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (s *StorageLifecyclePolicyDescribeOutputAssert) HasReturnType(expected datatypes.DataType) *StorageLifecyclePolicyDescribeOutputAssert {
	s.StringValueSet("return_type", expected.ToSql())
	return s
}

func (s *StorageLifecyclePolicyDescribeOutputAssert) HasSignature(expected ...sdk.TableColumnSignature) *StorageLifecyclePolicyDescribeOutputAssert {
	s.StringValueSet("signature.#", strconv.Itoa(len(expected)))
	for i, signature := range expected {
		s.StringValueSet(fmt.Sprintf("signature.%d.name", i), signature.Name)
		s.StringValueSet(fmt.Sprintf("signature.%d.type", i), signature.Type.ToSql())
	}
	return s
}
