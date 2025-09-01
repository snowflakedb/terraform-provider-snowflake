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
	// Granting a role to a role or user
	// When calling SHOW GRANTS TO ROLE / USER, the USAGE privilege should be shown with the ROLE granted_on field.
	// When calling SHOW GRANTS OF ROLE, the Role field should be populated.
	// The granted_to field should always point to either ROLE or USER.
	case (grant.Role != nil || grant.Privilege == "USAGE" && grant.GrantedOn == sdk.ObjectTypeRole && grant.Name != nil) &&
		(grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeUser):
		return MapToGrantAccountRole(grant)
	//case grant.Role != nil || (grant.GrantedOn == sdk.ObjectTypeDatabaseRole && (grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeDatabaseRole)):
	//	return MapToGrantDatabaseRole(grant)
	case len(privileges) > 0 && grant.GrantedTo == sdk.ObjectTypeRole:
		return MapToGrantPrivilegesToAccountRole(grant, privileges)
	//case grant.GrantedOn == sdk.ObjectTypeDatabaseRole:
	//	return MapToGrantPrivilegesToDatabaseRole(grant, privilegeListVariable)
	//// TODO: To share and To application role
	default:
		return nil, nil, fmt.Errorf("unsupported grant mapping")
	}

}

func MapToGrantAccountRole(grant sdk.Grant) (accconfig.ResourceModel, *ImportModel, error) {
	var roleIdentifier sdk.AccountObjectIdentifier

	switch {
	case grant.Role != nil:
		roleIdentifier = grant.Role.(sdk.AccountObjectIdentifier)
	case grant.Privilege == "USAGE" && grant.GrantedOn == sdk.ObjectTypeRole && grant.Name != nil:
		roleIdentifier = grant.Name.(sdk.AccountObjectIdentifier)
	default:
		return nil, nil, fmt.Errorf("invalid grant account role mapping: missing role information")
	}

	switch {
	case grant.GrantedTo == sdk.ObjectTypeRole:
		resourceId := MapResourceId(fmt.Sprintf("grant_%s_to_role_%s", roleIdentifier.Name(), grant.GranteeName.Name()))
		resourceModel := model.GrantAccountRole(resourceId, roleIdentifier.Name()).WithParentRoleName(grant.GranteeName.Name())

		stateResourceId := fmt.Sprintf("%s|%s|%s", roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeRole.String(), grant.GranteeName.FullyQualifiedName())

		return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId), nil
	case grant.GrantedTo == sdk.ObjectTypeUser:
		resourceId := MapResourceId(fmt.Sprintf("grant_%s_to_user_%s", roleIdentifier.Name(), grant.GranteeName.Name()))
		resourceModel := model.GrantAccountRole(resourceId, roleIdentifier.Name()).WithUserName(grant.GranteeName.Name())

		stateResourceId := fmt.Sprintf("%s|%s|%s", roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeUser.String(), grant.GranteeName.FullyQualifiedName())

		return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId), nil
	default:
		return nil, nil, fmt.Errorf("unsupported grant account role mapping")
	}
}

func MapToGrantPrivilegesToAccountRole(grant sdk.Grant, privileges []string) (accconfig.ResourceModel, *ImportModel, error) {
	var resourceModel accconfig.ResourceModel
	stateResourceId := resources.GrantPrivilegesToAccountRoleId{
		RoleName:        grant.GranteeName.(sdk.AccountObjectIdentifier),
		WithGrantOption: grant.GrantOption,
		Privileges:      privileges,
	}

	withGrantOption := "without"
	if grant.GrantOption {
		withGrantOption = "with"
	}

	switch {
	case grant.GrantedOn == sdk.ObjectTypeAccount:
		resourceId := MapResourceId(fmt.Sprintf("grant_on_account_to_%s_%s_grant_option", grant.GranteeName.Name(), withGrantOption))

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnAccount(true).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId.Kind = resources.OnAccountAccountRoleGrantKind
		stateResourceId.Data = new(resources.OnAccountGrantData)
	case slices.Contains(sdk.ValidGrantToAccountObjectTypesString, string(grant.GrantedOn)):
		resourceId := MapResourceId(fmt.Sprintf("grant_on_%s_%s_to_%s_%s_grant_option", grant.GrantedOn, grant.Name.Name(), grant.GranteeName.Name(), withGrantOption))

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnAccountObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.Name()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId.Kind = resources.OnAccountObjectAccountRoleGrantKind
		stateResourceId.Data = &resources.OnAccountObjectGrantData{
			ObjectType: grant.GrantedOn,
			ObjectName: grant.Name.(sdk.AccountObjectIdentifier),
		}
	case grant.GrantedOn == sdk.ObjectTypeSchema:
		resourceId := MapResourceId(fmt.Sprintf("grant_on_schema_%s_to_%s_%s_grant_option", grant.Name.FullyQualifiedName(), grant.GranteeName.Name(), withGrantOption))

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnSchemaValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"schema_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId.Kind = resources.OnSchemaAccountRoleGrantKind
		stateResourceId.Data = &resources.OnSchemaGrantData{
			Kind:       resources.OnSchemaSchemaGrantKind,
			SchemaName: sdk.Pointer(grant.Name.(sdk.DatabaseObjectIdentifier)),
		}
	case slices.Contains(sdk.ValidGrantToSchemaObjectTypesString, string(grant.GrantedOn)):
		resourceId := MapResourceId(fmt.Sprintf("grant_on_%s_%s_to_%s_%s_grant_option", grant.GrantedOn, grant.Name.FullyQualifiedName(), grant.GranteeName.Name(), withGrantOption))

		resourceModel = model.GrantPrivilegesToAccountRole(resourceId, grant.GranteeName.Name()).
			WithPrivileges(privileges).
			WithOnSchemaObjectValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"object_type": tfconfig.StringVariable(string(grant.GrantedOn)),
				"object_name": tfconfig.StringVariable(grant.Name.FullyQualifiedName()),
			})).
			WithWithGrantOption(grant.GrantOption)

		stateResourceId.Kind = resources.OnSchemaObjectAccountRoleGrantKind
		stateResourceId.Data = &resources.OnSchemaObjectGrantData{
			Kind: resources.OnObjectSchemaObjectGrantKind,
			Object: &sdk.Object{
				ObjectType: grant.GrantedOn,
				Name:       grant.Name.(sdk.SchemaObjectIdentifier),
			},
		}
	default:
		return nil, nil, fmt.Errorf("unsupported grant privileges to account role mapping")
	}

	return resourceModel, NewImportModel(resourceModel.ResourceReference(), stateResourceId.String()), nil
}
