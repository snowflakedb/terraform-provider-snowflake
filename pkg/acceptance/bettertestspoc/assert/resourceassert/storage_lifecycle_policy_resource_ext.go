package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLifecyclePolicyResourceAssert) HasArguments(args []sdk.TableColumnSignature) *StorageLifecyclePolicyResourceAssert {
	s.ValueSet("argument.#", strconv.Itoa(len(args)))
	for i, v := range args {
		s.ValueSet(fmt.Sprintf("argument.%d.name", i), v.Name)
		s.ValueSet(fmt.Sprintf("argument.%d.type", i), v.Type.ToSqlWithoutUnknowns())
	}
	return s
}
