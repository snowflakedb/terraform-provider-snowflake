package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeOpenflowRuntimeSchema represents output of DESCRIBE query for the single OpenflowRuntime.
// DESCRIBE returns the SHOW columns minus created_on, updated_on, database_name and schema_name, plus
// server_url and node_type_tier.
var DescribeOpenflowRuntimeSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"deployment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"min_nodes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"max_nodes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"node_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"display_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"external_access_integrations": {
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"initially_suspended": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"execute_as_role": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"key": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"server_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"node_type_tier": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = DescribeOpenflowRuntimeSchema

func OpenflowRuntimeDetailsToSchema(runtime sdk.OpenflowRuntimeDetails) map[string]any {
	runtimeSchema := make(map[string]any)
	runtimeSchema["name"] = runtime.Name
	runtimeSchema["status"] = string(runtime.Status)
	runtimeSchema["deployment"] = runtime.Deployment
	runtimeSchema["min_nodes"] = runtime.MinNodes
	runtimeSchema["max_nodes"] = runtime.MaxNodes
	runtimeSchema["node_type"] = string(runtime.NodeType)
	if runtime.DisplayName != nil {
		runtimeSchema["display_name"] = runtime.DisplayName
	}
	runtimeSchema["external_access_integrations"] = collections.Map(runtime.ExternalAccessIntegrations, sdk.AccountObjectIdentifier.Name)
	runtimeSchema["initially_suspended"] = runtime.InitiallySuspended
	if runtime.ExecuteAsRole != nil {
		runtimeSchema["execute_as_role"] = runtime.ExecuteAsRole
	}
	if runtime.Key != nil {
		runtimeSchema["key"] = runtime.Key
	}
	runtimeSchema["owner"] = runtime.Owner
	if runtime.Comment != nil {
		runtimeSchema["comment"] = runtime.Comment
	}
	if runtime.ServerUrl != nil {
		runtimeSchema["server_url"] = runtime.ServerUrl
	}
	if runtime.NodeTypeTier != nil {
		runtimeSchema["node_type_tier"] = runtime.NodeTypeTier
	}
	return runtimeSchema
}
