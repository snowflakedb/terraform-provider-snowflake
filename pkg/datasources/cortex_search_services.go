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

var cortexSearchServicesSchema = map[string]*schema.Schema{
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of cortex search services.",
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
		Description: "Filters the output with **case-sensitive** characters indicating the beginning of the object name.",
	},
	"limit": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Limits the number of rows returned. If the `limit.from` is set, then the limit will start from the first element matched by the expression. The expression is only used to match with the first element, later on the elements are not matched by the prefix, but you can enforce a certain pattern with `starts_with` or `like`.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rows": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The maximum number of rows to return.",
				},
				"from": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies a **case-sensitive** pattern that is used to match object name. After the first match, the limit on the number of rows will be applied.",
				},
			},
		},
	},
	"cortex_search_services": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the output of SHOW CORTEX SEARCH SERVICES.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"created_on": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date and time when the cortex search service was created.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the cortex search service.",
				},
				"database_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Database in which the cortex search service is stored.",
				},
				"schema_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Schema in which the cortex search service is stored.",
				},
				"comment": {
					Type:        schema.TypeString,
					Description: "Comment for the cortex search service.",
					Computed:    true,
				},
			},
		},
	},
}

// CortexSearchServices Snowflake Cortex search services resource.
func CortexSearchServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.CortexSearchServicesDatasource), TrackingReadWrapper(datasources.CortexSearchServices, ReadCortexSearchServices)),
		Schema:      cortexSearchServicesSchema,
	}
}

// ReadCortexSearchServices Reads the cortex search services metadata information.
func ReadCortexSearchServices(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	request := sdk.NewShowCortexSearchServiceRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		request.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
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

	dts, err := client.CortexSearchServices.Show(context.Background(), request)
	if err != nil {
		log.Printf("[DEBUG] snowflake_cortex_search_services.go: %v", err)
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId("cortex_search_services")
	records := make([]map[string]any, 0, len(dts))
	for _, dt := range dts {
		record := map[string]any{}
		record["created_on"] = dt.CreatedOn.String()
		record["name"] = dt.Name
		record["database_name"] = dt.DatabaseName
		record["schema_name"] = dt.SchemaName
		record["comment"] = dt.Comment
		records = append(records, record)
	}
	if err := d.Set("cortex_search_services", records); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
