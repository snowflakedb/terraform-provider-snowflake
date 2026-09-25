package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func DatabaseDetailsListToSchema(details *sdk.DatabaseDetails) []map[string]any {
	result := make([]map[string]any, len(details.Rows))
	for i := range details.Rows {
		result[i] = DatabaseDetailsRowToSchema(&details.Rows[i])
	}
	return result
}
