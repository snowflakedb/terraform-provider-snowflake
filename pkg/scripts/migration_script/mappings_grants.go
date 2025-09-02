package main

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func HandleGrants(config *Config, csvInput [][]string) (string, error) {
	grants, err := ConvertCsvInput[GrantCsvRow, sdk.Grant](csvInput)
	if err != nil {
		return "", err
	}

	groupedGrants := GroupGrants(grants)

	resourceModels := make([]accconfig.ResourceModel, 0)
	importModels := make([]ImportModel, 0)

	for _, grantGroup := range groupedGrants {
		mappedModel, importModel, err := MapGrantToModel(grantGroup)
		if err != nil {
			log.Printf("Error converting grant group: %+v to model: %v. Skipping grant and continuing with other mappings.", grantGroup, err)
		} else {
			resourceModels = append(resourceModels, mappedModel)
			importModels = append(importModels, *importModel)
		}
	}

	mappedModels, err := collections.MapErr(resourceModels, ResourceFromModel)
	if err != nil {
		return "", fmt.Errorf("errors from resource model to HCL conversion: %w", err)
	}

	mappedImports, err := collections.MapErr(importModels, func(importModel ImportModel) (string, error) {
		return TransformImportModel(config, importModel)
	})
	if err != nil {
		return "", fmt.Errorf("errors during import transformations: %v", err)
	}

	outputBuilder := new(strings.Builder)
	outputBuilder.WriteString(collections.JoinStrings(mappedModels, "\n"))
	outputBuilder.WriteString(collections.JoinStrings(mappedImports, ""))

	return outputBuilder.String(), nil
}

func GroupGrants(grants []sdk.Grant) map[string][]sdk.Grant {
	return collections.GroupByProperty(grants, func(grant sdk.Grant) string {
		return strings.Join([]string{
			grant.GrantedOn.String(),
			grant.Name.FullyQualifiedName(),
			grant.GrantedTo.String(),
			grant.GranteeName.FullyQualifiedName(),
			strconv.FormatBool(grant.GrantOption),
		}, "_")
	})
}

func MapGrantToModel(grantGroup []sdk.Grant) (accconfig.ResourceModel, *ImportModel, error) {
	// Assuming all grants in the group are the same type and only differ by privileges
	grant := grantGroup[0]

	privileges := collections.Map(grantGroup, func(grant sdk.Grant) string { return grant.Privilege })
	privileges = slices.DeleteFunc(privileges, func(s string) bool { return s == "" })

	// Remove duplicates
	slices.Sort(privileges)
	privileges = slices.Compact(privileges)

	switch {
	case grant.Privilege == "USAGE" &&
		grant.GrantedOn == sdk.ObjectTypeRole && grant.Name != nil &&
		(grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeUser):
		return MapToGrantAccountRole(grant)
	case grant.Privilege == "USAGE" &&
		grant.GrantedOn == sdk.ObjectTypeDatabaseRole && grant.Name != nil &&
		(grant.GrantedTo == sdk.ObjectTypeDatabaseRole || grant.GrantedTo == sdk.ObjectTypeRole):
		return MapToGrantDatabaseRole(grant)
	case grant.GrantedTo == sdk.ObjectTypeRole:
		return MapToGrantPrivilegesToAccountRole(grant, privileges)
	case grant.GrantedTo == sdk.ObjectTypeDatabaseRole:
		return MapToGrantPrivilegesToDatabaseRole(grant, privileges)
	default:
		return nil, nil, fmt.Errorf("unsupported grant mapping")
	}
}

func MapToGrantAccountRole(grant sdk.Grant) (accconfig.ResourceModel, *ImportModel, error) {
	switch grant.GrantedTo {
	case sdk.ObjectTypeRole:
		resourceId := NormalizeResourceId(fmt.Sprintf("grant_%s_to_role_%s", grant.Name.Name(), grant.GranteeName.Name()))
		resourceModel := model.GrantAccountRole(resourceId, grant.Name.Name()).WithParentRoleName(grant.GranteeName.Name())

		stateResourceId := fmt.Sprintf("%s|%s|%s", grant.Name.FullyQualifiedName(), sdk.ObjectTypeRole.String(), grant.GranteeName.FullyQualifiedName())

		return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId), nil
	case sdk.ObjectTypeUser:
		resourceId := NormalizeResourceId(fmt.Sprintf("grant_%s_to_user_%s", grant.Name.Name(), grant.GranteeName.Name()))
		resourceModel := model.GrantAccountRole(resourceId, grant.Name.Name()).WithUserName(grant.GranteeName.Name())

		stateResourceId := fmt.Sprintf("%s|%s|%s", grant.Name.FullyQualifiedName(), sdk.ObjectTypeUser.String(), grant.GranteeName.FullyQualifiedName())

		return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId), nil
	default:
		return nil, nil, fmt.Errorf("unsupported grant account role mapping")
	}
}

func MapToGrantDatabaseRole(grant sdk.Grant) (accconfig.ResourceModel, *ImportModel, error) {
	databaseRoleName, err := sdk.ParseDatabaseObjectIdentifier(grant.Name.FullyQualifiedName())
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing database object identifier from grant name %s: %w", grant.Name.FullyQualifiedName(), err)
	}

	switch grant.GrantedTo {
	case sdk.ObjectTypeDatabaseRole:
		granteeName := sdk.NewDatabaseObjectIdentifier(databaseRoleName.DatabaseName(), grant.GranteeName.Name())

		resourceId := NormalizeResourceId(fmt.Sprintf("grant_%s_to_database_role_%s", databaseRoleName.FullyQualifiedName(), granteeName.FullyQualifiedName()))
		resourceModel := model.GrantDatabaseRole(resourceId, databaseRoleName.FullyQualifiedName()).WithParentDatabaseRoleName(granteeName.FullyQualifiedName())

		stateResourceId := fmt.Sprintf("%s|%s|%s", databaseRoleName.FullyQualifiedName(), sdk.ObjectTypeDatabaseRole.String(), granteeName.FullyQualifiedName())

		return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId), nil
	case sdk.ObjectTypeRole:
		resourceId := NormalizeResourceId(fmt.Sprintf("grant_%s_to_role_%s", databaseRoleName.FullyQualifiedName(), grant.GranteeName.Name()))
		resourceModel := model.GrantDatabaseRole(resourceId, databaseRoleName.FullyQualifiedName()).WithParentRoleName(grant.GranteeName.Name())

		stateResourceId := fmt.Sprintf("%s|%s|%s", databaseRoleName.FullyQualifiedName(), sdk.ObjectTypeRole.String(), grant.GranteeName.FullyQualifiedName())

		return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId), nil
	default:
		return nil, nil, fmt.Errorf("unsupported grant database role mapping")
	}
}

func MapToGrantPrivilegesToAccountRole(grant sdk.Grant, privileges []string) (accconfig.ResourceModel, *ImportModel, error) {
	var resourceModel accconfig.ResourceModel
	var stateResourceId resources.GrantPrivilegesToAccountRoleId

	switch {
	case grant.GrantedOn == sdk.ObjectTypeAccount:
		resourceId := CreateGrantPrivilegesToAccountRoleResourceIdOnAccount(grant)

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnAccount(true).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId = resources.NewGrantPrivilegesToAccountRoleIdOnAccount(
			grant.GranteeName.(sdk.AccountObjectIdentifier),
			grant.GrantOption,
			false,
			false,
			privileges...,
		)
	case slices.Contains(sdk.ValidGrantToAccountObjectTypesString, string(grant.GrantedOn)):
		resourceId := CreateGrantPrivilegesToAccountRoleResourceIdOnAccountObject(grant)

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.Name()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId = resources.NewGrantPrivilegesToAccountRoleIdOnAccountObject(
			grant.GranteeName.(sdk.AccountObjectIdentifier),
			grant.GrantOption,
			false,
			false,
			grant.GrantedOn,
			grant.Name.(sdk.AccountObjectIdentifier),
			privileges...,
		)
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		resourceId := CreateGrantPrivilegesToAccountRoleResourceIdOnSchema(grant)

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId = resources.NewGrantPrivilegesToAccountRoleIdOnSchemaOnSchema(
			grant.GranteeName.(sdk.AccountObjectIdentifier),
			grant.GrantOption,
			false,
			false,
			grant.Name.(sdk.DatabaseObjectIdentifier),
			privileges...,
		)
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		resourceId := CreateGrantPrivilegesToAccountRoleResourceIdOnSchemaObject(grant)

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId = resources.NewGrantPrivilegesToAccountRoleIdOnSchemaObjectOnObject(
			grant.GranteeName.(sdk.AccountObjectIdentifier),
			grant.GrantOption,
			false,
			false,
			grant.GrantedOn,
			grant.Name.(sdk.SchemaObjectIdentifier),
			privileges...,
		)
	default:
		return nil, nil, fmt.Errorf("unsupported grant privileges to account role mapping")
	}

	return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId.String()), nil
}

func MapToGrantPrivilegesToDatabaseRole(grant sdk.Grant, privileges []string) (accconfig.ResourceModel, *ImportModel, error) {
	var resourceModel accconfig.ResourceModel
	var stateResourceId resources.GrantPrivilegesToDatabaseRoleId

	withGrantOption := "without"
	if grant.GrantOption {
		withGrantOption = "with"
	}

	switch {
	case grant.GrantedOn == sdk.ObjectTypeDatabase:
		databaseRoleName := grant.GranteeName.Name()
		granteeName := sdk.NewDatabaseObjectIdentifier(grant.Name.Name(), databaseRoleName)

		resourceId := NormalizeResourceId(fmt.Sprintf("grant_on_database_%s_to_%s_%s_grant_option", grant.Name.Name(), granteeName.FullyQualifiedName(), withGrantOption))

		resourceModel = model.GrantPrivilegesToDatabaseRole(resourceId, granteeName.FullyQualifiedName()).
			WithPrivileges(privileges).
			WithOnDatabase(grant.Name.Name()).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId = resources.NewGrantPrivilegesToDatabaseRoleIdOnDatabase(
			granteeName,
			grant.GrantOption,
			false,
			false,
			grant.Name.(sdk.AccountObjectIdentifier),
			privileges...,
		)
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		databaseRoleName := grant.GranteeName.Name()
		granteeName := sdk.NewDatabaseObjectIdentifier(grant.Name.(sdk.DatabaseObjectIdentifier).DatabaseName(), databaseRoleName)

		resourceId := NormalizeResourceId(fmt.Sprintf("grant_on_schema_%s_to_%s_%s_grant_option", grant.Name.FullyQualifiedName(), granteeName.FullyQualifiedName(), withGrantOption))

		resourceModel = model.GrantPrivilegesToDatabaseRole(resourceId, granteeName.FullyQualifiedName()).
			WithPrivileges(privileges).
			WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId = resources.NewGrantPrivilegesToDatabaseRoleIdOnSchemaOnSchema(
			granteeName,
			grant.GrantOption,
			false,
			false,
			grant.Name.(sdk.DatabaseObjectIdentifier),
			privileges...,
		)
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		databaseRoleName := grant.GranteeName.Name()
		granteeName := sdk.NewDatabaseObjectIdentifier(grant.Name.(sdk.SchemaObjectIdentifier).DatabaseName(), databaseRoleName)

		resourceId := NormalizeResourceId(fmt.Sprintf("grant_on_%s_%s_to_%s_%s_grant_option", grant.GrantedOn, grant.Name.FullyQualifiedName(), granteeName.FullyQualifiedName(), withGrantOption))

		resourceModel = model.GrantPrivilegesToDatabaseRole(resourceId, granteeName.FullyQualifiedName()).
			WithPrivileges(privileges).
			WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId = resources.NewGrantPrivilegesToDatabaseRoleIdOnSchemaObjectOnObject(
			granteeName,
			grant.GrantOption,
			false,
			false,
			grant.GrantedOn,
			grant.Name.(sdk.SchemaObjectIdentifier),
			privileges...,
		)
	default:
		return nil, nil, fmt.Errorf("unsupported grant privileges to database role mapping")
	}

	return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId.String()), nil
}
