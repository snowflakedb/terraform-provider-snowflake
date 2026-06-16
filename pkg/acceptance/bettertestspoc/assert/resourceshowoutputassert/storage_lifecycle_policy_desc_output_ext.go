package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLifecyclePolicyDescribeOutputAssert) HasSignature(expected ...sdk.TableColumnSignature) *StorageLifecyclePolicyDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("signature.#", strconv.Itoa(len(expected))))
	for i, signature := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("signature.%d.name", i), signature.Name))
		s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("signature.%d.type", i), signature.Type.ToSql()))
	}
	return s
}
