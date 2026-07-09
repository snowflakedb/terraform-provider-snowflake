//go:build sdk_generation

package defs

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
)

func init() {
	gen.AllSdkObjectDefinitions = append(
		gen.AllSdkObjectDefinitions,
		alertsDef,
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
		databaseRolesDef,
		databasesDef,
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
		pipesDef,
		postgresInstancesDef,
		proceduresDef,
		resourceMonitorsDef,
		rolesDef,
		rowAccessPoliciesDef,
		secretsDef,
		schemasDef,
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
		tablesDef,
		tagReferencesDef,
		tagsDef,
		tasksDef,
		userProgrammaticAccessTokensDef,
		usersDef,
		viewsDef,
		warehousesDef,
	)
	fmt.Println("SDK object definitions:")
	for _, def := range gen.AllSdkObjectDefinitions {
		fmt.Printf(" - %s\n", def.Name)
	}
}
