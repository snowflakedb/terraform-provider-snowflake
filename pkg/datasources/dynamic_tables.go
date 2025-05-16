package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var dynamicTablesSchema = map[string]*schema.Schema{
	"like": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "LIKE clause to filter the list of dynamic tables.",
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
		Description: "IN clause to filter the list of dynamic tables.",
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
		Description: "Optionally limits the maximum number of rows returned, while also enabling “pagination” of the results. Note that the actual number of rows returned might be less than the specified limit (e.g. the number of existing objects is less than the specified limit).",
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
					Description:  "The optional FROM 'name_string' subclause effectively serves as a “cursor” for the results. This enables fetching the specified number of rows following the first row whose object name matches the specified string",
					RequiredWith: []string{"limit.0.rows"},
				},
			},
		},
	},
	"records": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The list of dynamic tables.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"created_on": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Date and time when the dynamic table was created.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name of the dynamic table.",
				},
				"database_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Database in which the dynamic table is stored.",
				},
				"schema_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Schema in which the dynamic table is stored.",
				},
				"cluster_by": {
					Type:        schema.TypeString,
					Description: "The clustering key for the dynamic table.",
					Computed:    true,
				},
				"rows": {
					Type:        schema.TypeInt,
					Description: "Number of rows in the table.",
					Computed:    true,
				},
				"bytes": {
					Type:        schema.TypeInt,
					Description: "Number of bytes that will be scanned if the entire dynamic table is scanned in a query.",
					Computed:    true,
				},
				"owner": {
					Type:        schema.TypeString,
					Description: "Role that owns the dynamic table.",
					Computed:    true,
				},
				"target_lag": {
					Type:        schema.TypeString,
					Description: "The maximum duration that the dynamic table’s content should lag behind real time.",
					Computed:    true,
				},
				"refresh_mode": {
					Type:        schema.TypeString,
					Description: "INCREMENTAL if the dynamic table will use incremental refreshes, or FULL if it will recompute the whole table on every refresh.",
					Computed:    true,
				},
				"refresh_mode_reason": {
					Type:        schema.TypeString,
					Description: "Explanation for why FULL refresh mode was chosen. NULL if refresh mode is not FULL.",
					Computed:    true,
				},
				"warehouse": {
					Type:        schema.TypeString,
					Description: "Warehouse that provides the required resources to perform the incremental refreshes.",
					Computed:    true,
				},
				"comment": {
					Type:        schema.TypeString,
					Description: "Comment for the dynamic table.",
					Computed:    true,
				},
				"text": {
					Type:        schema.TypeString,
					Description: "The text of the command that created this dynamic table (e.g. CREATE DYNAMIC TABLE ...).",
					Computed:    true,
				},
				"automatic_clustering": {
					Type:        schema.TypeBool,
					Description: "Whether auto-clustering is enabled on the dynamic table. Not currently supported for dynamic tables.",
					Computed:    true,
				},
				"scheduling_state": {
					Type:        schema.TypeString,
					Description: "Displays RUNNING for dynamic tables that are actively scheduling refreshes and SUSPENDED for suspended dynamic tables.",
					Computed:    true,
				},
				"last_suspended_on": {
					Type:        schema.TypeString,
					Description: "Timestamp of last suspension.",
					Computed:    true,
				},
				"is_clone": {
					Type:        schema.TypeBool,
					Description: "TRUE if the dynamic table has been cloned, else FALSE.",
					Computed:    true,
				},
				"is_replica": {
					Type:        schema.TypeBool,
					Description: "TRUE if the dynamic table is a replica. else FALSE.",
					Computed:    true,
				},
				"data_timestamp": {
					Type:        schema.TypeString,
					Description: "Timestamp of the data in the base object(s) that is included in the dynamic table.",
					Computed:    true,
				},
			},
		},
	},
}

// DynamicTables Snowflake Dynamic Tables resource.
func DynamicTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.DynamicTablesDatasource), TrackingReadWrapper(datasources.DynamicTables, ReadDynamicTables)),
		Schema:      dynamicTablesSchema,
	}
}

// ReadDynamicTables Reads the dynamic tables metadata information.
func ReadDynamicTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	request := sdk.NewShowDynamicTableRequest()
	if v, ok := d.GetOk("like"); ok {
		like := v.([]interface{})[0].(map[string]interface{})
		pattern := like["pattern"].(string)
		request.WithLike(&sdk.Like{Pattern: sdk.String(pattern)})
	}

	if v, ok := d.GetOk("in"); ok {
		in := v.([]interface{})[0].(map[string]interface{})
		if v, ok := in["account"]; ok {
			account := v.(bool)
			if account {
				request.WithIn(&sdk.In{Account: sdk.Bool(account)})
			}
		}
		if v, ok := in["database"]; ok {
			database := v.(string)
			if database != "" {
				request.WithIn(&sdk.In{Database: sdk.NewAccountObjectIdentifier(database)})
			}
		}
		if v, ok := in["schema"]; ok {
			schema := v.(string)
			if schema != "" {
				request.WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schema)})
			}
		}
	}
	if v, ok := d.GetOk("starts_with"); ok {
		startsWith := v.(string)
		request.WithStartsWith(sdk.String(startsWith))
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
		request.WithLimit(&limit)
	}

	dts, err := client.DynamicTables.Show(context.Background(), request)
	if err != nil {
		log.Printf("[DEBUG] snowflake_dynamic_tables.go: %v", err)
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId("dynamic_tables")
	records := make([]map[string]any, 0, len(dts))
	for _, dt := range dts {
		record := map[string]any{}
		/*
			guides on time formatting
			https://docs.snowflake.com/en/user-guide/date-time-input-output
			https://pkg.go.dev/time
			note: format may depend on what the account parameter for TIMESTAMP_OUTPUT_FORMAT is set to. Perhaps we should return this as a string rather than a time.Time?
		*/
		record["created_on"] = dt.CreatedOn.Format("2006-01-02T16:04:05.000 -0700")
		record["name"] = dt.Name
		record["database_name"] = dt.DatabaseName
		record["schema_name"] = dt.SchemaName
		record["cluster_by"] = dt.ClusterBy
		record["rows"] = dt.Rows
		record["bytes"] = dt.Bytes
		record["owner"] = dt.Owner
		record["target_lag"] = dt.TargetLag
		record["refresh_mode"] = string(dt.RefreshMode)
		record["refresh_mode_reason"] = dt.RefreshModeReason
		record["warehouse"] = dt.Warehouse
		record["comment"] = dt.Comment
		record["text"] = dt.Text
		record["automatic_clustering"] = dt.AutomaticClustering
		record["scheduling_state"] = string(dt.SchedulingState)
		record["last_suspended_on"] = dt.LastSuspendedOn.Format("2006-01-02T16:04:05.000 -0700")
		record["is_clone"] = dt.IsClone
		record["is_replica"] = dt.IsReplica
		record["data_timestamp"] = dt.DataTimestamp.Format("2006-01-02T16:04:05.000 -0700")
		records = append(records, record)
	}
	if err := d.Set("records", records); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
