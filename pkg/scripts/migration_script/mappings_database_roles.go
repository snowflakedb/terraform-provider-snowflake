package main

import (
	"fmt"
	"strings"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HandleDatabaseRoles(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[DatabaseRoleCsvRow, DatabaseRoleRepresentation](config, csvInput, MapDatabaseRoleToModel)
}

func MapDatabaseRoleToModel(role DatabaseRoleRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	roleId := sdk.NewDatabaseObjectIdentifier(role.DatabaseName, role.Name)
	resourceId := NormalizeResourceId(fmt.Sprintf("%s_%s", strings.TrimPrefix(string(resources.DatabaseRole), "snowflake_"), roleId.FullyQualifiedName()))
	resourceModel := model.DatabaseRole(resourceId, role.DatabaseName, role.Name)

	handleIfNotEmpty(role.Comment, resourceModel.WithComment)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		roleId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
