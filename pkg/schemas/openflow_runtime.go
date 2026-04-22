package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowOpenflowRuntimeSchema = map[string]*schema.Schema{
	"name": {Type: schema.TypeString, Computed: true},
	"status": {Type: schema.TypeString, Computed: true},
	"database_name": {Type: schema.TypeString, Computed: true},
	"schema_name": {Type: schema.TypeString, Computed: true},
	"deployment": {Type: schema.TypeString, Computed: true},
	"node_type": {Type: schema.TypeString, Computed: true},
	"min_nodes": {Type: schema.TypeInt, Computed: true},
	"max_nodes": {Type: schema.TypeInt, Computed: true},
	"execute_as_role": {Type: schema.TypeString, Computed: true},
	"external_access_integrations": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"display_name": {Type: schema.TypeString, Computed: true},
	"comment": {Type: schema.TypeString, Computed: true},
	"owner": {Type: schema.TypeString, Computed: true},
	"created_on": {Type: schema.TypeString, Computed: true},
}

var _ = ShowOpenflowRuntimeSchema

func OpenflowRuntimeToSchema(r *sdk.OpenflowRuntime) map[string]any {
	m := make(map[string]any)
	m["name"] = r.Name
	m["status"] = string(r.Status)
	m["database_name"] = r.DatabaseName
	m["schema_name"] = r.SchemaName
	m["deployment"] = r.Deployment
	m["node_type"] = string(r.NodeType)
	m["min_nodes"] = r.MinNodes
	m["max_nodes"] = r.MaxNodes
	m["execute_as_role"] = r.ExecuteAsRole
	m["external_access_integrations"] = r.ExternalAccessIntegrations
	m["owner"] = r.Owner
	m["created_on"] = r.CreatedOn.String()
	if r.DisplayName != nil {
		m["display_name"] = *r.DisplayName
	}
	if r.Comment != nil {
		m["comment"] = *r.Comment
	}
	return m
}

var DescribeOpenflowRuntimeSchema = map[string]*schema.Schema{
	"name": {Type: schema.TypeString, Computed: true},
	"status": {Type: schema.TypeString, Computed: true},
	"database_name": {Type: schema.TypeString, Computed: true},
	"schema_name": {Type: schema.TypeString, Computed: true},
	"deployment": {Type: schema.TypeString, Computed: true},
	"node_type": {Type: schema.TypeString, Computed: true},
	"min_nodes": {Type: schema.TypeInt, Computed: true},
	"max_nodes": {Type: schema.TypeInt, Computed: true},
	"execute_as_role": {Type: schema.TypeString, Computed: true},
	"external_access_integrations": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"display_name": {Type: schema.TypeString, Computed: true},
	"comment": {Type: schema.TypeString, Computed: true},
	"owner": {Type: schema.TypeString, Computed: true},
	"created_on": {Type: schema.TypeString, Computed: true},
	"error_code": {Type: schema.TypeString, Computed: true},
	"status_message": {Type: schema.TypeString, Computed: true},
}

func OpenflowRuntimeDetailsToSchema(r sdk.OpenflowRuntimeDetails) map[string]any {
	m := make(map[string]any)
	m["name"] = r.Name
	m["status"] = string(r.Status)
	m["database_name"] = r.DatabaseName
	m["schema_name"] = r.SchemaName
	m["deployment"] = r.Deployment
	m["node_type"] = string(r.NodeType)
	m["min_nodes"] = r.MinNodes
	m["max_nodes"] = r.MaxNodes
	m["execute_as_role"] = r.ExecuteAsRole
	m["external_access_integrations"] = r.ExternalAccessIntegrations
	m["owner"] = r.Owner
	m["created_on"] = r.CreatedOn.String()
	if r.DisplayName != nil {
		m["display_name"] = *r.DisplayName
	}
	if r.Comment != nil {
		m["comment"] = *r.Comment
	}
	if r.ErrorCode != nil {
		m["error_code"] = *r.ErrorCode
	}
	if r.StatusMessage != nil {
		m["status_message"] = *r.StatusMessage
	}
	return m
}
