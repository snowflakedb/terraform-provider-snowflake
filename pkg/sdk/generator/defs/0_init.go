//go:build sdk_generation

package defs

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
)

func init() {
	gen.AllSdkObjectDefinitions = append(gen.AllSdkObjectDefinitions,
		apiIntegrationsDef,
		applicationPackagesDef,
		applicationRolesDef,
		applicationsDef,
		authenticationPoliciesDef,
		catalogIntegrationsDef,
		computePoolsDef,
		connectionsDef,
		cortexSearchServicesDef,
		dataMetricFunctionReferencesDef,
		eventTablesDef,
		externalFunctionsDef,
		externalVolumesDef,
		fileFormatsDef,
		functionsDef,
		gitRepositoriesDef,
		hybridTablesDef,
		imageRepositoriesDef,
		listingsDef,
		managedAccountsDef,
		materializedViewsDef,
		networkPoliciesDef,
		networkRulesDef,
		notebooksDef,
		notificationIntegrationsDef,
		organizationAccountsDef,
		proceduresDef,
		rowAccessPoliciesDef,
		secretsDef,
		securityIntegrationsDef,
		semanticViewsDef,
		sequencesDef,
		servicesDef,
		sessionPoliciesDef,
		stagesDef,
		storageIntegrationsDef,
		streamlitsDef,
		streamsDef,
		tasksDef,
		userProgrammaticAccessTokensDef,
		viewsDef,
	)
	fmt.Println("SDK object definitions:")
	for _, def := range gen.AllSdkObjectDefinitions {
		fmt.Printf(" - %s\n", def.Name)
	}
}
