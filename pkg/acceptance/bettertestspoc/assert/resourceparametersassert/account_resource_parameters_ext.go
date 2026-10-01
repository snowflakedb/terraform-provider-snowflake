package resourceparametersassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

// parametersAbsentFromShow lists catalog parameters that SHOW PARAMETERS IN ACCOUNT does not return.
// SHARE_RESTRICTIONS can be set on an account, but Snowflake omits it from SHOW PARAMETERS.
var parametersAbsentFromShow = map[string]struct{}{
	"SHARE_RESTRICTIONS": {},
}

// TODO(SNOW-4173331): generate HasAllParametersPresent from the parameter assertion generator once it can assert that every catalog parameter is present without checking values.
// HasAllParametersPresent checks that every account parameter from the catalog is present.
// Values are not checked, because they depend on the connected account.
func (a *AccountResourceParametersAssert) HasAllParametersPresent() *AccountResourceParametersAssert {
	for _, def := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelAccountExt) {
		if _, absent := parametersAbsentFromShow[def.SqlName]; absent {
			a.CollectionLength(def.FieldName(), 0)
			continue
		}
		a.CollectionLength(def.FieldName(), 1)
	}
	return a
}
