package main

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"slices"
	"testing"
	_ "unsafe"
)

//go:linkname convertGrantRow github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk.(*GrantRow).convert
func convertGrantRow(row *sdk.GrantRow) *sdk.Grant

func HandleGrants(csvInput [][]string) {
	grants := ConvertCsvInput[sdk.GrantRow, sdk.Grant](
		csvInput,
		WithAdditionalConvertMapping(func(row sdk.GrantRow, convertedValue *sdk.Grant) {
			if row.GranteeName != "" {
				convertedValue.GranteeName = sdk.NewAccountObjectIdentifier(row.GranteeName)
			}
		}),
	)

	// TODO: Group same grants with different privileges

	resourceModels := make([]config.ResourceModel, 0)
	for _, grant := range grants {
		resourceModels = append(resourceModels, MapGrantToModel(grant))
	}

	models := collections.Map(resourceModels, func(resourceModel config.ResourceModel) any { return any(resourceModel) })
	mappedModels := config.FromModels(&testing.T{}, models...)
	fmt.Println(mappedModels)
}

// TODO: it should receive []sdk.Grant because there may be a few rows for the same type but different privileges
// TODO: there should be a grouping step before this function.
func MapGrantToModel(grant sdk.Grant) config.ResourceModel {
	switch {
	case grant.GrantedOn == sdk.ObjectTypeAccount:
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_account", grant.GranteeName.Name()).
			WithPrivilegesValue(tfconfig.ListVariable(
				tfconfig.StringVariable(grant.Privilege),
			)).
			WithOnAccount(true).
			WithWithGrantOption(grant.GrantOption)
	case slices.Contains(sdk.ValidGrantToAccountObjectTypesString, string(grant.GrantedOn)):
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_account_object", grant.GranteeName.Name()).
			WithPrivilegesValue(tfconfig.ListVariable(
				tfconfig.StringVariable(grant.Privilege),
			)).
			WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.Name()),
			})).
			WithWithGrantOption(grant.GrantOption)
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_schema", grant.GranteeName.Name()).
			WithPrivilegesValue(tfconfig.ListVariable(
				tfconfig.StringVariable(grant.Privilege),
			)).
			WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		return model.GrantPrivilegesToAccountRole("test_resource_name_on_schema_object", grant.GranteeName.Name()).
			WithPrivilegesValue(tfconfig.ListVariable(
				tfconfig.StringVariable(grant.Privilege),
			)).
			WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)
	default:
		return model.GrantPrivilegesToAccountRole("test_resource_name", "test_account_role")
	}
}
