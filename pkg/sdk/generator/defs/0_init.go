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
		budgetsDef,
		catalogIntegrationsDef,
		computePoolsDef,
		connectionsDef,
		cortexAgentsDef,
		cortexSearchServicesDef,
		dataMetricFunctionReferencesDef,
		eventTablesDef,
		externalFunctionsDef,
		externalVolumesDef,
		fileFormatsDef,
		functionsDef,
		gitRepositoriesDef,
		hybridTablesDef,
		icebergTablesDef,
		imageRepositoriesDef,
		listingsDef,
		managedAccountsDef,
		materializedViewsDef,
		networkPoliciesDef,
		networkRulesDef,
		notebooksDef,
		notificationIntegrationsDef,
		openflowConnectorsDef,
		openflowDeploymentsDef,
		openflowRuntimesDef,
		organizationAccountsDef,
		passwordPoliciesDef,
		postgresInstancesDef,
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
		storageLifecyclePoliciesDef,
		streamlitsDef,
		streamsDef,
		tagReferencesDef,
		tasksDef,
		userProgrammaticAccessTokensDef,
		viewsDef,
	)
	fmt.Println("SDK object definitions:")
	for _, def := range gen.AllSdkObjectDefinitions {
		fmt.Printf(" - %s\n", def.Name)
	}
}
