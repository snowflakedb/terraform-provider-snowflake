package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func StageWithId(resourceName string, id sdk.SchemaObjectIdentifier) *StageModel {
	return Stage(resourceName, id.DatabaseName(), id.SchemaName(), id.Name())
}
