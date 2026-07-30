package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func WarehouseInteractiveWithId(id sdk.AccountObjectIdentifier) *WarehouseInteractiveModel {
	return WarehouseInteractiveWithDefaultMeta(id.Name())
}

func (w *WarehouseInteractiveModel) WithTables(tables ...string) *WarehouseInteractiveModel {
	variables := make([]tfconfig.Variable, len(tables))
	for i, table := range tables {
		variables[i] = tfconfig.StringVariable(table)
	}
	w.Tables = tfconfig.SetVariable(variables...)
	return w
}
