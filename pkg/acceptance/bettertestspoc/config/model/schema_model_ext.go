package model

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
)

func SchemaWithImplicitDatabaseDependency(
	resourceName string,
	schemaName string,
	databaseModel *DatabaseModel,
) *SchemaModel {
	return Schema(resourceName, "", schemaName).
		WithDatabaseValue(config.UnquotedWrapperVariable(fmt.Sprintf("%s.name", databaseModel.ResourceReference())))
}
