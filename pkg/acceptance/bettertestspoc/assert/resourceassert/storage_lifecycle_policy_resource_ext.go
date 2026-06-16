package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLifecyclePolicyResourceAssert) HasArguments(args []sdk.TableColumnSignature) *StorageLifecyclePolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("argument.#", strconv.Itoa(len(args))))
	for i, v := range args {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("argument.%d.name", i), v.Name))
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("argument.%d.type", i), v.Type.ToSqlWithoutUnknowns()))
	}
	return s
}
