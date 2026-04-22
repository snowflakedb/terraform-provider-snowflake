package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowOpenflowConnectorSchema = map[string]*schema.Schema{
	"name":                 {Type: schema.TypeString, Computed: true},
	"status":               {Type: schema.TypeString, Computed: true},
	"database_name":        {Type: schema.TypeString, Computed: true},
	"schema_name":          {Type: schema.TypeString, Computed: true},
	"runtime":              {Type: schema.TypeString, Computed: true},
	"connector_definition": {Type: schema.TypeString, Computed: true},
	"started":              {Type: schema.TypeBool, Computed: true},
	"display_name":         {Type: schema.TypeString, Computed: true},
	"comment":              {Type: schema.TypeString, Computed: true},
	"owner":                {Type: schema.TypeString, Computed: true},
	"created_on":           {Type: schema.TypeString, Computed: true},
}

var _ = ShowOpenflowConnectorSchema

func OpenflowConnectorToSchema(c *sdk.OpenflowConnector) map[string]any {
	m := make(map[string]any)
	m["name"] = c.Name
	m["status"] = string(c.Status)
	m["database_name"] = c.DatabaseName
	m["schema_name"] = c.SchemaName
	m["runtime"] = c.Runtime
	m["started"] = c.Started
	m["owner"] = c.Owner
	m["created_on"] = c.CreatedOn.String()
	if c.ConnectorDefinition != nil {
		m["connector_definition"] = *c.ConnectorDefinition
	}
	if c.DisplayName != nil {
		m["display_name"] = *c.DisplayName
	}
	if c.Comment != nil {
		m["comment"] = *c.Comment
	}
	return m
}

var DescribeOpenflowConnectorSchema = map[string]*schema.Schema{
	"name":                 {Type: schema.TypeString, Computed: true},
	"status":               {Type: schema.TypeString, Computed: true},
	"database_name":        {Type: schema.TypeString, Computed: true},
	"schema_name":          {Type: schema.TypeString, Computed: true},
	"runtime":              {Type: schema.TypeString, Computed: true},
	"connector_definition": {Type: schema.TypeString, Computed: true},
	"started":              {Type: schema.TypeBool, Computed: true},
	"display_name":         {Type: schema.TypeString, Computed: true},
	"comment":              {Type: schema.TypeString, Computed: true},
	"owner":                {Type: schema.TypeString, Computed: true},
	"created_on":           {Type: schema.TypeString, Computed: true},
	"error_code":           {Type: schema.TypeString, Computed: true},
	"status_message":       {Type: schema.TypeString, Computed: true},
}

func OpenflowConnectorDetailsToSchema(c sdk.OpenflowConnectorDetails) map[string]any {
	m := make(map[string]any)
	m["name"] = c.Name
	m["status"] = string(c.Status)
	m["database_name"] = c.DatabaseName
	m["schema_name"] = c.SchemaName
	m["runtime"] = c.Runtime
	m["started"] = c.Started
	m["owner"] = c.Owner
	m["created_on"] = c.CreatedOn.String()
	if c.ConnectorDefinition != nil {
		m["connector_definition"] = *c.ConnectorDefinition
	}
	if c.DisplayName != nil {
		m["display_name"] = *c.DisplayName
	}
	if c.Comment != nil {
		m["comment"] = *c.Comment
	}
	if c.ErrorCode != nil {
		m["error_code"] = *c.ErrorCode
	}
	if c.StatusMessage != nil {
		m["status_message"] = *c.StatusMessage
	}
	return m
}
