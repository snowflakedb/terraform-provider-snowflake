package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func StreamOnExternalTableBase(resourceName string, id, externalTableId sdk.SchemaObjectIdentifier) *StreamOnExternalTableModel {
	return StreamOnExternalTable(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), externalTableId.FullyQualifiedName()).WithInsertOnly("true")
}
