package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func StreamOnTableBase(resourceName string, id, tableId sdk.SchemaObjectIdentifier) *StreamOnTableModel {
	return StreamOnTable(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), tableId.FullyQualifiedName())
}
