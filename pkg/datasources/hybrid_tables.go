package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var hybridTablesSchema = map[string]*schema.Schema{
	"like": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "LIKE clause to filter the list of hybrid tables.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"pattern": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Filters the command output by object name. The filter uses case-insensitive pattern matching with support for SQL wildcard characters (% and _).",
				},
			},
		},
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of hybrid tables.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account": {
					Type:         schema.TypeBool,
					Optional:     true,
					Description:  "Returns records for the entire account.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
				"database": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current database in use or for a specified database (db_name).",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
				"schema": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current schema in use or a specified schema (schema_name).",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
			},
		},
	},
	"starts_with": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Optionally filters the command output based on the characters that appear at the beginning of the object name. The string is case-sensitive.",
	},
	"limit": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Optionally limits the maximum number of rows returned, while also enabling \"pagination\" of the results. Note that the actual number of rows returned might be less than the specified limit (e.g. the number of existing objects is less than the specified limit).",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rows": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Specifies the maximum number of rows to return.",
				},
				"from": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "The optional FROM 'name_string' subclause effectively serves as a \"cursor\" for the results. This enables fetching the specified number of rows following the first row whose object name matches the specified string",
					RequiredWith: []string{"limit.0.rows"},
				},
			},
		},
	},
	"hybrid_tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The list of hybrid tables.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"created_on": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date and time when the hybrid table was created.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the hybrid table.",
				},
				"database_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Database in which the hybrid table is stored.",
				},
				"schema_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Schema in which the hybrid table is stored.",
				},
				"owner": {
					Type:        schema.TypeString,
					Description: "Role that owns the hybrid table.",
					Computed:    true,
				},
				"rows": {
					Type:        schema.TypeInt,
					Description: "Number of rows in the table.",
					Computed:    true,
				},
				"bytes": {
					Type:        schema.TypeInt,
					Description: "Number of bytes that will be scanned if the entire hybrid table is scanned in a query.",
					Computed:    true,
				},
				"comment": {
					Type:        schema.TypeString,
					Description: "Comment for the hybrid table.",
					Computed:    true,
				},
				"owner_role_type": {
					Type:        schema.TypeString,
					Description: "The type of role that owns the object, either ROLE or DATABASE_ROLE.",
					Computed:    true,
				},
			},
		},
	},
}

// HybridTables Snowflake Hybrid Tables resource.
func HybridTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.HybridTablesDatasource), TrackingReadWrapper(datasources.HybridTables, ReadHybridTables)),
		Schema:      hybridTablesSchema,
	}
}

// ReadHybridTables Reads the hybrid tables metadata information.
func ReadHybridTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	request := sdk.NewShowHybridTableRequest()
	if v, ok := d.GetOk("like"); ok {
		like := v.([]interface{})[0].(map[string]interface{})
		pattern := like["pattern"].(string)
		request.WithLike(sdk.Like{Pattern: sdk.String(pattern)})
	}

	if v, ok := d.GetOk("in"); ok {
		in := v.([]interface{})[0].(map[string]interface{})
		if v, ok := in["account"]; ok {
			account := v.(bool)
			if account {
				request.WithIn(sdk.In{Account: sdk.Bool(account)})
			}
		}
		if v, ok := in["database"]; ok {
			database := v.(string)
			if database != "" {
				request.WithIn(sdk.In{Database: sdk.NewAccountObjectIdentifier(database)})
			}
		}
		if v, ok := in["schema"]; ok {
			schema := v.(string)
			if schema != "" {
				request.WithIn(sdk.In{Schema: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schema)})
			}
		}
	}
	if v, ok := d.GetOk("starts_with"); ok {
		startsWith := v.(string)
		request.WithStartsWith(startsWith)
	}
	if v, ok := d.GetOk("limit"); ok {
		l := v.([]interface{})[0].(map[string]interface{})
		limit := sdk.LimitFrom{}
		if v, ok := l["rows"]; ok {
			rows := v.(int)
			limit.Rows = sdk.Int(rows)
		}
		if v, ok := l["from"]; ok {
			from := v.(string)
			limit.From = sdk.String(from)
		}
		request.WithLimit(limit)
	}

	hts, err := client.HybridTables.Show(ctx, request)
	if err != nil {
		log.Printf("[DEBUG] snowflake_hybrid_tables.go: %v", err)
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId("hybrid_tables")
	records := make([]map[string]any, 0, len(hts))
	for _, ht := range hts {
		record := map[string]any{}
		record["created_on"] = ht.CreatedOn.Format("2006-01-02T16:04:05.000 -0700")
		record["name"] = ht.Name
		record["database_name"] = ht.DatabaseName
		record["schema_name"] = ht.SchemaName
		record["owner"] = ht.Owner
		record["rows"] = ht.Rows
		record["bytes"] = ht.Bytes
		record["comment"] = ht.Comment
		record["owner_role_type"] = ht.OwnerRoleType
		records = append(records, record)
	}
	if err := d.Set("hybrid_tables", records); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
