package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (r *RowAccessPolicyResourceAssert) HasArguments(args []sdk.TableColumnSignature) *RowAccessPolicyResourceAssert {
	r.ValueSet("argument.#", strconv.FormatInt(int64(len(args)), 10))
	for i, v := range args {
		r.ValueSet(fmt.Sprintf("argument.%d.name", i), v.Name)
		r.ValueSet(fmt.Sprintf("argument.%d.type", i), v.Type.ToSql())
	}
	return r
}
