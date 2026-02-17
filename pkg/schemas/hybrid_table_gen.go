// Code generated for hybrid table schema mappings; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowHybridTableSchema represents output of SHOW query for a single HybridTable.
var ShowHybridTableSchema = map[string]*schema.Schema{
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
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"rows": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"bytes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowHybridTableSchema

func HybridTableToSchema(hybridTable *sdk.HybridTable) map[string]any {
	hybridTableSchema := make(map[string]any)
	hybridTableSchema["created_on"] = hybridTable.CreatedOn.String()
	hybridTableSchema["name"] = hybridTable.Name
	hybridTableSchema["database_name"] = hybridTable.DatabaseName
	hybridTableSchema["schema_name"] = hybridTable.SchemaName
	hybridTableSchema["owner"] = hybridTable.Owner
	hybridTableSchema["rows"] = hybridTable.Rows
	hybridTableSchema["bytes"] = hybridTable.Bytes
	hybridTableSchema["comment"] = hybridTable.Comment
	hybridTableSchema["owner_role_type"] = hybridTable.OwnerRoleType
	return hybridTableSchema
}

var _ = HybridTableToSchema

// HybridTableDescribeSchema represents output of DESCRIBE query for HybridTable.
var HybridTableDescribeSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_nullable": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"primary_key": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"unique_key": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"check": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"expression": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"policy_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"privacy_domain": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_evolution_record": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = HybridTableDescribeSchema

func HybridTableDetailsToSchema(details *sdk.HybridTableDetails) map[string]any {
	detailsSchema := make(map[string]any)
	detailsSchema["name"] = details.Name
	detailsSchema["type"] = details.Type
	detailsSchema["kind"] = details.Kind
	detailsSchema["is_nullable"] = details.IsNullable
	detailsSchema["default"] = details.Default
	detailsSchema["primary_key"] = details.PrimaryKey
	detailsSchema["unique_key"] = details.UniqueKey
	detailsSchema["check"] = details.Check
	detailsSchema["expression"] = details.Expression
	detailsSchema["comment"] = details.Comment
	detailsSchema["policy_name"] = details.PolicyName
	detailsSchema["privacy_domain"] = details.PrivacyDomain
	detailsSchema["schema_evolution_record"] = details.SchemaEvolutionRecord
	return detailsSchema
}

var _ = HybridTableDetailsToSchema
