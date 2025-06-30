// Content of this file should be moved to production files after proceeding with Terraform Plugin Framework.

package testfunctional

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func booleanAttributeCreate(boolAttribute types.Bool, createField **bool) {
	if !boolAttribute.IsNull() {
		*createField = sdk.Bool(boolAttribute.ValueBool())
	}
}

func int64AttributeCreate(int64Attribute types.Int64, createField **int) {
	if !int64Attribute.IsNull() {
		*createField = sdk.Int(int(int64Attribute.ValueInt64()))
	}
}

func stringAttributeCreate(stringAttribute types.String, createField **string) {
	if !stringAttribute.IsNull() {
		*createField = sdk.String(stringAttribute.ValueString())
	}
}
