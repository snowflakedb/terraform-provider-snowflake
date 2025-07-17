package customtypes

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func StringWithMetadataAttributeCreate(v StringWithMetadataValue, createField **string) {
	if !v.IsNull() {
		*createField = sdk.String(v.ValueString())
	}
}
