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
		DataMetricFunctionReferenceDef,
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
		SemanticViewsDef,
		SequencesDef,
	)
	fmt.Println("SDK object definitions:")
	for _, def := range gen.AllSdkObjectDefinitions {
		fmt.Printf(" - %s\n", def.Name)
	}
}
