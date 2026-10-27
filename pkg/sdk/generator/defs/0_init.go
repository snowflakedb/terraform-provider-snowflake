//go:build sdk_generation

package defs

import (
	"fmt"

	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
)

func init() {
	gen.AllSdkObjectDefinitions = append(
		gen.AllSdkObjectDefinitions,
		accountsDef,
		alertsDef,
		apiIntegrationsDef,
		applicationPackagesDef,
		applicationRolesDef,
		applicationsDef,
		authenticationPoliciesDef,
		budgetsDef,
		connectionsDef,
		contextFunctionsDef,
		cortexAgentsDef,
		cortexSearchServicesDef,
		databaseRolesDef,
		databasesDef,
		dynamicTablesDef,
		dataMetricFunctionReferencesDef,
		eventTablesDef,
		externalTablesDef,
		externalAccessIntegrationsDef,
		externalFunctionsDef,
		failoverGroupsDef,
		fileFormatsDef,
		functionsDef,
		gitRepositoriesDef,
		grantsDef,
		hybridTablesDef,
		icebergTablesDef,
		imageRepositoriesDef,
		listingsDef,
		managedAccountsDef,
		maskingPoliciesDef,
		materializedViewsDef,
		mcpServersDef,
		networkPoliciesDef,
		networkRulesDef,
		notebooksDef,
		notificationIntegrationsDef,
		openflowConnectorDefinitionsDef,
		openflowConnectorsDef,
		openflowDeploymentsDef,
		openflowRuntimesDef,
		organizationAccountsDef,
		parametersDef,
		passwordPoliciesDef,
		pipesDef,
		policyReferencesDef,
		postgresInstancesDef,
		postgresMirrorsDef,
		proceduresDef,
		replicationFunctionsDef,
		resourceMonitorsDef,
		rolesDef,
		rowAccessPoliciesDef,
		secretsDef,
		securityIntegrationsDef,
		sessionsDef,
		sharesDef,
		semanticViewsDef,
		sequencesDef,
		servicesDef,
		sessionPoliciesDef,
		snowflakeIntelligencesDef,
		stagesDef,
		storageIntegrationsDef,
		storageLifecyclePoliciesDef,
		streamlitsDef,
		streamsDef,
		systemFunctionsDef,
		tablesDef,
		tagReferencesDef,
		tagsDef,
		tasksDef,
		userProgrammaticAccessTokensDef,
		usersDef,
		viewsDef,
	)
	gen.AllSdkObjectDefinitions = gen.MergeSdkDefinitions(gen.AllSdkObjectDefinitions, gen.SdkInterfacesFromSsot(ssotdefs.All(), sdkInterfaceOverrides))
	fmt.Println("SDK object definitions:")
	for _, def := range gen.AllSdkObjectDefinitions {
		fmt.Printf(" - %s\n", def.Name)
	}
}
