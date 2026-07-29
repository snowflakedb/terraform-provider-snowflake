package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeIcebergTableSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"source_iceberg_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_nullable": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"default": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"primary_key": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"unique_key": {
		Type:     schema.TypeBool,
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
	"name_mapping": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"write_default": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func IcebergTableDetailsToSchema(details []sdk.IcebergTableDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i, d := range details {
		row := map[string]any{
			"name":                d.Name,
			"type":                d.Type.ToSql(),
			"source_iceberg_type": d.SourceIcebergType,
			"kind":                d.Kind,
			"is_nullable":         d.IsNullable,
			"primary_key":         d.PrimaryKey,
			"unique_key":          d.UniqueKey,
		}
		if d.Default != nil {
			row["default"] = *d.Default
		}
		if d.Check != nil {
			row["check"] = *d.Check
		}
		if d.Expression != nil {
			row["expression"] = *d.Expression
		}
		if d.Comment != nil {
			row["comment"] = *d.Comment
		}
		if d.PolicyName != nil {
			row["policy_name"] = d.PolicyName.FullyQualifiedName()
		}
		if d.PrivacyDomain != nil {
			row["privacy_domain"] = *d.PrivacyDomain
		}
		if d.NameMapping != nil {
			row["name_mapping"] = *d.NameMapping
		}
		if d.WriteDefault != nil {
			row["write_default"] = *d.WriteDefault
		}
		result[i] = row
	}
	return result
}
