package poc

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type SdkObjectDef struct {
	name string
	// TODO [next PR]: can be removed?
	file       string
	definition *generator.Interface
}

func GetSdkDefinitions() []*generator.Interface {
	allDefinitions := allSdkObjectDefinitions
	definitions := make([]*generator.Interface, len(allDefinitions))
	for idx, s := range allDefinitions {
		definitions[idx] = s.definition
	}
	return definitions
}

func WithPreamble(i *generator.Interface, preamble *genhelpers.PreambleModel) *generator.Interface {
	i.PreambleModel = preamble
	return i
}

var allSdkObjectDefinitions = []SdkObjectDef{
	{
		name:       "NetworkPolicies",
		file:       "network_policies_def.go",
		definition: sdk.NetworkPoliciesDef,
	},
}
