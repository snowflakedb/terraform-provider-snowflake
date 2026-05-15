package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeHybridTableSchema = map[string]*schema.Schema{
	"name":                    {Type: schema.TypeString, Computed: true},
	"type":                    {Type: schema.TypeString, Computed: true},
	"collation":               {Type: schema.TypeString, Computed: true},
	"kind":                    {Type: schema.TypeString, Computed: true},
	"is_nullable":             {Type: schema.TypeBool, Computed: true},
	"default":                 {Type: schema.TypeString, Computed: true},
	"primary_key":             {Type: schema.TypeBool, Computed: true},
	"unique_key":              {Type: schema.TypeBool, Computed: true},
	"check":                   {Type: schema.TypeString, Computed: true},
	"expression":              {Type: schema.TypeString, Computed: true},
	"comment":                 {Type: schema.TypeString, Computed: true},
	"policy_name":             {Type: schema.TypeString, Computed: true},
	"privacy_domain":          {Type: schema.TypeString, Computed: true},
	"schema_evolution_record": {Type: schema.TypeString, Computed: true},
}

func HybridTableDetailsToSchema(detail sdk.HybridTableDetails) map[string]any {
	collation := ""
	if detail.Collation != nil {
		collation = *detail.Collation
	}
	return map[string]any{
		"name":                    detail.Name,
		"type":                    detail.Type,
		"collation":               collation,
		"kind":                    detail.Kind,
		"is_nullable":             detail.IsNullable,
		"default":                 detail.Default,
		"primary_key":             detail.PrimaryKey,
		"unique_key":              detail.UniqueKey,
		"check":                   detail.Check,
		"expression":              detail.Expression,
		"comment":                 detail.Comment,
		"policy_name":             detail.PolicyName,
		"privacy_domain":          detail.PrivacyDomain,
		"schema_evolution_record": detail.SchemaEvolutionRecord,
	}
}

func HybridTableDetailsListToSchema(details []sdk.HybridTableDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i, d := range details {
		result[i] = HybridTableDetailsToSchema(d)
	}
	return result
}
