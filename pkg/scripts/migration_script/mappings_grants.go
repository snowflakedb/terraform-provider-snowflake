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

func HandleGrants(csvInput [][]string) error {
	grants, err := ConvertCsvInput[GrantCsvRow, sdk.Grant](csvInput)
	if err != nil {
		return err
	}

	groupedGrants := GroupGrants(grants)

	resourceModels := make([]accconfig.ResourceModel, 0)
	for _, grantGroup := range groupedGrants {
		mappedModel, err := MapGrantToModel(grantGroup)
		if err != nil {
			log.Printf("Error converting grant group: %+v to model: %v. Skipping grant and continuing with other mappings.", grantGroup, err)
		} else {
			resourceModels = append(resourceModels, mappedModel)
		}
	}

	mappedModels := collections.Map(resourceModels, func(resourceModel accconfig.ResourceModel) string {
		return accconfig.ResourceFromModel(&testing.T{}, resourceModel)
	})
	fmt.Println(collections.JoinStrings(mappedModels, "\n"))

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

func MapGrantToModel(grantGroup []sdk.Grant) (accconfig.ResourceModel, error) {
	// Assuming all grants in the group are the same type and only differ by privileges
	grant := grantGroup[0]
	privileges := collections.Map(grantGroup, func(grant sdk.Grant) string { return grant.Privilege })
	privilegeListVariable := tfconfig.ListVariable(
		collections.Map(privileges, func(privilege string) tfconfig.Variable {
			return tfconfig.StringVariable(privilege)
		})...,
	)

	switch {
	//// TODO: Check how it's returned for SHOW GRANTS OF DATABASE ROLE
	case grant.Role != nil || (grant.GrantedOn == sdk.ObjectTypeRole && (grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeUser)):
		return MapToGrantAccountRole(grant)
	case grant.Role != nil || (grant.GrantedOn == sdk.ObjectTypeDatabaseRole && (grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole)):
		return MapToGrantDatabaseRole(grant)
	case grant.GrantedOn == sdk.ObjectTypeRole:
		return MapToGrantPrivilegesToAccountRole(grant, privilegeListVariable)
	case grant.GrantedOn == sdk.ObjectTypeDatabaseRole:
		return MapToGrantPrivilegesToDatabaseRole(grant, privilegeListVariable)
	//	// TODO: To share and To application role
	default:
		return nil, fmt.Errorf("skipping unsupported grant: %+v", grant)
	}
}

func MapToGrantPrivilegesToAccountRole(grant sdk.Grant, privilegeListVariable tfconfig.Variable) (accconfig.ResourceModel, error) {
	// TODO: Check other outputs
	switch {
	case grant.GrantedOn == sdk.ObjectTypeAccount:
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_account", grant.GranteeName.Name()).
				WithPrivilegesValue(privilegeListVariable).
				WithOnAccount(true).
				WithWithGrantOption(grant.GrantOption),
			nil
	case slices.Contains(sdk.ValidGrantToAccountObjectTypesString, string(grant.GrantedOn)):
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_account_object", grant.GranteeName.Name()).
				WithPrivilegesValue(privilegeListVariable).
				WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
					"object_name": tfconfig.StringVariable(grant.Name.Name()),
				})).
				WithWithGrantOption(grant.GrantOption),
			nil
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_schema", grant.GranteeName.Name()).
				WithPrivilegesValue(privilegeListVariable).
				WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
				})).
				WithWithGrantOption(grant.GrantOption),
			nil
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_schema_object", grant.GranteeName.Name()).
				WithPrivilegesValue(privilegeListVariable).
				WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
					"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
				})).
				WithWithGrantOption(grant.GrantOption),
			nil
	default:
		return nil, fmt.Errorf("unsupported grant mapping")
	}
}

func MapToGrantPrivilegesToDatabaseRole(grant sdk.Grant, privilegeListVariable tfconfig.Variable) (accconfig.ResourceModel, error) {
	// TODO: Check other outputs
	switch {
	case grant.GrantedOn == sdk.ObjectTypeDatabase:
		return model.GrantPrivilegesToDatabaseRole("test_resource_name_on_schema", grant.GranteeName.Name()).
				WithPrivilegesValue(privilegeListVariable).
				WithOnDatabase(grant.Name.Name()).
				WithWithGrantOption(grant.GrantOption),
			nil
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		return model.GrantPrivilegesToDatabaseRole("test_resource_name_on_schema", grant.GranteeName.Name()).
				WithPrivilegesValue(privilegeListVariable).
				WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
				})).
				WithWithGrantOption(grant.GrantOption),
			nil
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		return model.GrantPrivilegesToDatabaseRole("test_resource_name_on_schema_object", grant.GranteeName.Name()).
				WithPrivilegesValue(privilegeListVariable).
				WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
					"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
				})).
				WithWithGrantOption(grant.GrantOption),
			nil
	default:
		return nil, fmt.Errorf("skipping unsupported grant: %+v", grant)
	}
}

func MapToGrantAccountRole(grant sdk.Grant) (accconfig.ResourceModel, error) {
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
	return result, nil
}

func MapToGrantDatabaseRole(grant sdk.Grant) (accconfig.ResourceModel, error) {
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
	return result, nil
}
