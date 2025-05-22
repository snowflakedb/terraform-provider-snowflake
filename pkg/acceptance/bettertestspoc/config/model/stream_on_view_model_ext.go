package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func StreamOnViewBase(resourceName string, id sdk.SchemaObjectIdentifier, viewId sdk.SchemaObjectIdentifier) *StreamOnViewModel {
	return StreamOnView(resourceName, id.DatabaseName(), id.SchemaName(), id.Name(), viewId.FullyQualifiedName())
}
