package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowSemanticViewSchema represents output of SHOW query for the single SemanticView.
var ShowSemanticViewSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"extension": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func SemanticViewToSchema(semanticView *sdk.SemanticView) map[string]any {
	semanticViewSchema := make(map[string]any)
	semanticViewSchema["created_on"] = semanticView.CreatedOn.String()
	semanticViewSchema["name"] = semanticView.Name
	semanticViewSchema["database_name"] = semanticView.DatabaseName
	semanticViewSchema["schema_name"] = semanticView.SchemaName
	if semanticView.Comment != nil {
		semanticViewSchema["comment"] = semanticView.Comment
	}
	semanticViewSchema["owner"] = semanticView.Owner
	semanticViewSchema["owner_role_type"] = semanticView.OwnerRoleType
	if semanticView.Extension != nil {
		semanticViewSchema["extension"] = semanticView.Extension
	}
	return semanticViewSchema
}
