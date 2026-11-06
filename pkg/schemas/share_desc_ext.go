package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func ShareDetailsListToSchema(details []sdk.ShareDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i := range details {
		result[i] = ShareDetailsToSchema(&details[i])
	}
	return result
}
