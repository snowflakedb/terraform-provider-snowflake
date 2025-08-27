package main

import (
	"fmt"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"log"
	"slices"
)

func HandleGrants(csvInput [][]string) error {
	grants, err := ConvertCsvInput[GrantCsvRow, sdk.Grant](csvInput)
	if err != nil {
		return err
	}

	// TODO(SNOW-2277608): Group same grants with different privileges

	resourceModels := make([]accconfig.ResourceModel, 0)
	for _, grant := range grants {
		mappedModel, err := MapGrantToModel(grant)
		if err != nil {
			log.Printf("Error converting grant %+v to model: %v. Skipping grant and continuing with other mappings.", grant, err)
		} else {
			resourceModels = append(resourceModels, mappedModel)
		}
	}

	mappedModels, err := collections.MapErr(resourceModels, ResourceFromModel)
	if err != nil {
		return fmt.Errorf("errors from resource model to HCL conversion: %w", err)
	}
	fmt.Println(collections.JoinStrings(mappedModels, "\n"))

	return nil
}

// TODO(SNOW-2277608): it should receive []sdk.Grant because there may be a few rows for the same type but different privileges
// TODO(SNOW-2277608): there should be a grouping step before this function.
func MapGrantToModel(grant sdk.Grant) (accconfig.ResourceModel, error) {
	switch {
	case grant.GrantedOn == sdk.ObjectTypeAccount:
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_account", grant.GranteeName.Name()).
				WithPrivilegesValue(tfconfig.ListVariable(
					tfconfig.StringVariable(grant.Privilege),
				)).
				WithOnAccount(true).
				WithWithGrantOption(grant.GrantOption),
			nil
	case slices.Contains(sdk.ValidGrantToAccountObjectTypesString, string(grant.GrantedOn)):
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_account_object", grant.GranteeName.Name()).
				WithPrivilegesValue(tfconfig.ListVariable(
					tfconfig.StringVariable(grant.Privilege),
				)).
				WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
					"object_name": tfconfig.StringVariable(grant.Name.Name()),
				})).
				WithWithGrantOption(grant.GrantOption),
			nil
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_schema", grant.GranteeName.Name()).
				WithPrivilegesValue(tfconfig.ListVariable(
					tfconfig.StringVariable(grant.Privilege),
				)).
				WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
				})).
				WithWithGrantOption(grant.GrantOption),
			nil
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_schema_object", grant.GranteeName.Name()).
			WithPrivilegesValue(tfconfig.ListVariable(
				tfconfig.StringVariable(grant.Privilege),
			)).
			WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption), nil
	default:
		return nil, fmt.Errorf("unsupported grant mapping")
	}
}
