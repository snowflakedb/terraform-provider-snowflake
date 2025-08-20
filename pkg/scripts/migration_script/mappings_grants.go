package main

import (
	"fmt"
	"log"
	"slices"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
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

func MapGrantToModel(grantGroup []sdk.Grant) config.ResourceModel {
	// Assuming all grants in the group are the same type and only differ by privileges
	grant := grantGroup[0]
	privileges := collections.Map(grantGroup, func(grant sdk.Grant) string { return grant.Privilege })
	privilegeListVariable := tfconfig.ListVariable(
		collections.Map(privileges, func(privilege string) tfconfig.Variable {
			return tfconfig.StringVariable(privilege)
		})...,
	)
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
			WithWithGrantOption(grant.GrantOption), nil
	default:
		return nil, fmt.Errorf("unsupported grant mapping")
	}
}
