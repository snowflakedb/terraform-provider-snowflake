package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func HybridTableDetailsListToSchema(details []sdk.HybridTableDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i := range details {
		result[i] = HybridTableDetailsToSchema(&details[i])
	}
	return result
}
