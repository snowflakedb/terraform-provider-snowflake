package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeSemanticViewSchema represents output of DESCRIBE query for the single SemanticView
var DescribeSemanticViewSchema = map[string]*schema.Schema{
	"object_kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"object_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"parent_entity": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"property": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"property_value": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func SemanticViewDetailsToSchema(semanticViewDetails []sdk.SemanticViewDetails) []map[string]any {
	semanticViewSchema := make([]map[string]any, len(semanticViewDetails))

	for i, detail := range semanticViewDetails {
		row := make(map[string]any)
		row["object_kind"] = detail.ObjectKind
		row["object_name"] = detail.ObjectName
		row["property"] = detail.Property
		row["property_value"] = detail.PropertyValue
		if detail.ParentEntity != nil {
			row["parent_entity"] = detail.ParentEntity
		}
		semanticViewSchema[i] = row
	}

	return semanticViewSchema
}
