package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (t *TableSearchOptimizationDetailsAssert) HasTargetDataTypeSql(expected datatypes.DataType) *TableSearchOptimizationDetailsAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.TableSearchOptimizationDetails) error {
		t.Helper()
		if o.TargetDataType.ToSql() != expected.ToSql() {
			return fmt.Errorf("expected target data type: %v; got: %v", expected.ToSql(), o.TargetDataType.ToSql())
		}
		return nil
	})
	return t
}
