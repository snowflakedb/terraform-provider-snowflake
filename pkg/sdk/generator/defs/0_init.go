//go:build sdk_generation

package defs

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
)

func init() {
	gen.AllSdkObjectDefinitions = append(gen.AllSdkObjectDefinitions,
		ApiIntegrationsDef,
		ApplicationPackagesDef,
		ApplicationRolesDef,
		ApplicationsDef,
		AuthenticationPoliciesDef,
		ComputePoolsDef,
		ConnectionsDef,
		CortexSearchServicesDef,
		DataMetricFunctionReferencesDef,
		EventTablesDef,
		ExternalFunctionsDef,
		FunctionsDef,
		GitRepositoriesDef,
		ImageRepositoriesDef,
		ListingsDef,
		ManagedAccountsDef,
		MaterializedViewsDef,
		NetworkPoliciesDef,
		NetworkRulesDef,
		NotebooksDef,
		NotificationIntegrationsDef,
		OrganizationAccountsDef,
		ProceduresDef,
		RowAccessPoliciesDef,
		SecretsDef,
		SecurityIntegrationsDef,
		SemanticViewsDef,
		SequencesDef,
		ServicesDef,
		SessionPoliciesDef,
		StagesDef,
		StorageIntegrationsDef,
		StreamlitsDef,
		StreamsDef,
		TasksDef,
		UserProgrammaticAccessTokensDef,
	)
	fmt.Println("SDK object definitions:")
	for _, def := range gen.AllSdkObjectDefinitions {
		fmt.Printf(" - %s\n", def.Name)
	}
}
