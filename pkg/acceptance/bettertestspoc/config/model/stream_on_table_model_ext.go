package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func StreamOnTableBase(resourceName string, id, tableId sdk.SchemaObjectIdentifier) *StreamOnTableModel {
	return StreamOnTable(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), tableId.FullyQualifiedName())
}
