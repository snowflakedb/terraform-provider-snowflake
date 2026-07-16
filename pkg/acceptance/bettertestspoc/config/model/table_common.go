package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func columnsVariable(columns []sdk.Column) tfconfig.Variable {
	return tfconfig.ListVariable(collections.Map(columns, func(c sdk.Column) tfconfig.Variable { return tfconfig.StringVariable(c.Value) })...)
}

func setStringIfNotNil(m map[string]tfconfig.Variable, key string, value *string) {
	if value != nil {
		m[key] = tfconfig.StringVariable(*value)
	}
}

func setBoolPairIfNotNil(m map[string]tfconfig.Variable, key string, trueVal, falseVal *bool) {
	if trueVal != nil {
		m[key] = tfconfig.StringVariable(resources.BooleanTrue)
	}
	if falseVal != nil {
		m[key] = tfconfig.StringVariable(resources.BooleanFalse)
	}
}
