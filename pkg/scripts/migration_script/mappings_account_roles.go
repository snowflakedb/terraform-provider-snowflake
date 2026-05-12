package main

import (
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HandleAccountRoles(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[AccountRoleCsvRow, AccountRoleRepresentation](config, csvInput, MapAccountRoleToModel)
}

func MapAccountRoleToModel(role AccountRoleRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	roleId := sdk.NewAccountObjectIdentifier(role.Name)
	resourceId := ResourceId(resources.AccountRole, roleId.FullyQualifiedName())
	resourceModel := model.AccountRole(resourceId, role.Name)

	handleIfNotEmpty(role.Comment, resourceModel.WithComment)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		roleId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
