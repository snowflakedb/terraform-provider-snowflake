package main

import (
	"fmt"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"log"
	"slices"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HandleGrants(config *Config, csvInput [][]string) error {
	grants, err := ConvertCsvInput[GrantCsvRow, sdk.Grant](csvInput)
	if err != nil {
		return err
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

	mappedModels := collections.Map(resourceModels, func(resourceModel accconfig.ResourceModel) string {
		return accconfig.ResourceFromModel(&testing.T{}, resourceModel)
	})

	mappedImports, err := collections.MapErr(importModels, func(importModel ImportModel) (string, error) {
		return TransformImportModel(config, importModel)
	})
	if err != nil {
		log.Printf("Errors during import transformations: %v", err)
	}

	fmt.Println(collections.JoinStrings(mappedModels, "\n"))
	fmt.Println(collections.JoinStrings(mappedImports, ""))

	return nil
}

func GroupGrants(grants []sdk.Grant) map[string][]sdk.Grant {
	return GroupByProperty(grants, func(grant sdk.Grant) string {
		return strings.Join([]string{
			grant.GrantOn.String(),
			grant.Name.FullyQualifiedName(),
			grant.GrantedTo.String(),
			grant.GranteeName.FullyQualifiedName(),
		}, "_")
	})
}

func MapGrantToModel(grantGroup []sdk.Grant) (accconfig.ResourceModel, *ImportModel, error) {
	// Assuming all grants in the group are the same type and only differ by privileges
	grant := grantGroup[0]
	privileges := collections.Map(grantGroup, func(grant sdk.Grant) string { return grant.Privilege })
	privilegeListVariable := tfconfig.ListVariable(
		collections.Map(privileges, func(privilege string) tfconfig.Variable {
			return tfconfig.StringVariable(privilege)
		})...,
	)

	switch {
	// TODO: Check how it's returned for SHOW GRANTS OF DATABASE ROLE
	case grant.Role != nil || (grant.GrantedOn == sdk.ObjectTypeRole && (grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeUser)):
		return MapToGrantAccountRole(grant)
	case grant.Role != nil || (grant.GrantedOn == sdk.ObjectTypeDatabaseRole && (grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole)):
		return MapToGrantDatabaseRole(grant)
	case grant.GrantedTo == sdk.ObjectTypeRole:
		return MapToGrantPrivilegesToAccountRole(grant, privileges, privilegeListVariable)
	case grant.GrantedTo == sdk.ObjectTypeDatabaseRole:
		return MapToGrantPrivilegesToDatabaseRole(grant, privileges, privilegeListVariable)
	// TODO: To share and To application role
	default:
		return nil, nil, fmt.Errorf("skipping unsupported grant: %+v", grant)
	}
}

func MapToGrantPrivilegesToAccountRole(grant sdk.Grant, privileges []string, privilegeListVariable tfconfig.Variable) (accconfig.ResourceModel, *ImportModel, error) {
	// TODO: Check other outputs
	switch {
	case grant.GrantedOn == sdk.ObjectTypeAccount:
		resourceModel := model.GrantPrivilegesToAccountRole("test_resource_name_on_account", grant.GranteeName.Name()).
			WithPrivilegesValue(privilegeListVariable).
			WithOnAccount(true).
			WithWithGrantOption(grant.GrantOption)

		return resourceModel, &ImportModel{
			ResourceAddress: resourceModel.ResourceReference(),
			Id: fmt.Sprintf(
				`"%s"|%t|false|%s|OnAccount`,
				grant.GranteeName.Name(),
				grant.GrantOption,
				collections.JoinStrings(privileges, ","),
			),
		}, nil
	case slices.Contains(sdk.ValidGrantToAccountObjectTypesString, string(grant.GrantedOn)):
		resourceModel := model.GrantPrivilegesToAccountRole("test_resource_name_on_account_object", grant.GranteeName.Name()).
			WithPrivilegesValue(privilegeListVariable).
			WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.Name()),
			})).
			WithWithGrantOption(grant.GrantOption)

		return resourceModel, &ImportModel{
			ResourceAddress: resourceModel.ResourceReference(),
			Id: fmt.Sprintf(
				`"%s"|%t|false|%s|OnAccountObject|%s|"%s"`,
				grant.GranteeName.Name(),
				grant.GrantOption,
				collections.JoinStrings(privileges, ","),
				grant.GrantedOn.String(),
				grant.Name.Name(),
			),
		}, nil
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		resourceModel := model.GrantPrivilegesToAccountRole("test_resource_name_on_schema", grant.GranteeName.Name()).
			WithPrivilegesValue(privilegeListVariable).
			WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		return resourceModel, &ImportModel{
			ResourceAddress: resourceModel.ResourceReference(),
			Id: fmt.Sprintf(
				`"%s"|%t|false|%s|OnSchema|OnSchema|%s`,
				grant.GranteeName.Name(),
				grant.GrantOption,
				collections.JoinStrings(privileges, ","),
				grant.Name.FullyQualifiedName(),
			),
		}, nil
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		resourceModel := model.GrantPrivilegesToAccountRole("test_resource_name_on_schema_object", grant.GranteeName.Name()).
			WithPrivilegesValue(privilegeListVariable).
			WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		return resourceModel, &ImportModel{
			ResourceAddress: resourceModel.ResourceReference(),
			Id: fmt.Sprintf(
				`"%s"|%t|false|%s|OnSchemaObject|OnObject|%s|%s`,
				grant.GranteeName.Name(),
				grant.GrantOption,
				collections.JoinStrings(privileges, ","),
				grant.GrantedOn.String(),
				grant.Name.FullyQualifiedName(),
			),
		}, nil
	default:
		return nil, nil, fmt.Errorf("unsupported grant mapping")
	}
}

func MapToGrantPrivilegesToDatabaseRole(grant sdk.Grant, privileges []string, privilegeListVariable tfconfig.Variable) (accconfig.ResourceModel, *ImportModel, error) {
	// TODO: Check other outputs
	switch {
	case grant.GrantedOn == sdk.ObjectTypeDatabase:
		resourceModel := model.GrantPrivilegesToDatabaseRole("test_resource_name_on_schema", grant.GranteeName.Name()).
			WithPrivilegesValue(privilegeListVariable).
			WithOnDatabase(grant.Name.Name()).
			WithWithGrantOption(grant.GrantOption)

		return resourceModel, &ImportModel{
			ResourceAddress: resourceModel.ResourceReference(),
			Id: fmt.Sprintf(
				"%s|%t|false|%s|OnDatabase|%s",
				grant.GranteeName.FullyQualifiedName(),
				grant.GrantOption,
				collections.JoinStrings(privileges, ","),
				grant.Name.FullyQualifiedName(),
			),
		}, nil
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		resourceModel := model.GrantPrivilegesToDatabaseRole("test_resource_name_on_schema", grant.GranteeName.Name()).
			WithPrivilegesValue(privilegeListVariable).
			WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		return resourceModel, &ImportModel{
			ResourceAddress: resourceModel.ResourceReference(),
			Id: fmt.Sprintf(
				"%s|%t|false|%s|OnSchema|OnSchema|%s",
				grant.GranteeName.FullyQualifiedName(),
				grant.GrantOption,
				collections.JoinStrings(privileges, ","),
				grant.Name.FullyQualifiedName(),
			),
		}, nil
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		resourceModel := model.GrantPrivilegesToDatabaseRole("test_resource_name_on_schema_object", grant.GranteeName.Name()).
			WithPrivilegesValue(privilegeListVariable).
			WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		return resourceModel, &ImportModel{
			ResourceAddress: resourceModel.ResourceReference(),
			Id: fmt.Sprintf(
				"%s|%t|false|%s|OnSchemaObject|OnObject|%s|%s",
				grant.GranteeName.FullyQualifiedName(),
				grant.GrantOption,
				collections.JoinStrings(privileges, ","),
				grant.GrantedOn.String(),
				grant.Name.FullyQualifiedName(),
			),
		}, nil
	default:
		return nil, nil, fmt.Errorf("skipping unsupported grant: %+v", grant)
	}
}

func MapToGrantAccountRole(grant sdk.Grant) (accconfig.ResourceModel, *ImportModel, error) {
	var result *model.GrantAccountRoleModel

	if grant.Role != nil {
		// Handle SHOW GRANTS OF ROLE output
		result = model.GrantAccountRole("test_resource_grant_account_role", grant.Role.Name())
	} else {
		// Handle SHOW GRANTS TO X
		result = model.GrantAccountRole("test_resource_grant_account_role", grant.Name.Name())
	}

	if grant.GrantedTo == sdk.ObjectTypeUser {
		result.WithUserName(grant.GranteeName.Name())
	} else if grant.GrantedTo == sdk.ObjectTypeRole {
		result.WithParentRoleName(grant.GranteeName.Name())
	}

	// TODO: Check other outputs
	return result, &ImportModel{
		ResourceAddress: result.ResourceReference(),
		// role_name (string) | grantee_object_type (ROLE|USER) | grantee_name (string)
		Id: fmt.Sprintf(`"%s"|%s|"%s"`, grant.Name.Name(), grant.GrantedTo.String(), grant.GranteeName.Name()),
	}, nil
}

func MapToGrantDatabaseRole(grant sdk.Grant) (accconfig.ResourceModel, *ImportModel, error) {
	var result *model.GrantDatabaseRoleModel

	if grant.Role != nil {
		// Handle SHOW GRANTS OF ROLE output
		result = model.GrantDatabaseRole("test_resource_grant_database_role", grant.Role.FullyQualifiedName())
	} else {
		// Handle SHOW GRANTS TO X
		result = model.GrantDatabaseRole("test_resource_grant_database_role", grant.Name.FullyQualifiedName())
	}

	if grant.GrantedTo == sdk.ObjectTypeDatabaseRole {
		result.WithDatabaseRoleName(grant.GranteeName.Name())
	} else if grant.GrantedTo == sdk.ObjectTypeRole {
		result.WithParentRoleName(grant.GranteeName.Name())
	}

	// TODO: Check other outputs
	return result, &ImportModel{
		ResourceAddress: result.ResourceReference(),
		// TODO: <database_role_identifier>|<object_type (ROLE|DATABASE ROLE|SHARE)>|<object_name>
		Id: fmt.Sprintf("%s|%s|%s", grant.Name.FullyQualifiedName(), grant.GrantedTo.String(), grant.GranteeName.FullyQualifiedName()),
	}, nil
}
