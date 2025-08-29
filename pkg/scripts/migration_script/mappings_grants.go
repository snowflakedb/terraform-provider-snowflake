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
		log.Printf("Errors during import transformations: %v", err)
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

	// Remove duplicates
	slices.Sort(privileges)
	privileges = slices.Compact(privileges)

	switch {
	case grant.GrantedOn == sdk.ObjectTypeAccount:
		resourceModel := model.GrantPrivilegesToAccountRole("test_resource_name_on_account", grant.GranteeName.Name()).
			WithPrivileges(privileges).
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
			WithPrivileges(privileges).
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
			WithPrivileges(privileges).
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
			WithPrivileges(privileges).
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
