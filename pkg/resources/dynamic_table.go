package resources

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var dynamicTableShema = map[string]*schema.Schema{
	"or_replace": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to replace the dynamic table if it already exists.",
		Default:     false,
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier (i.e. name) for the dynamic table; must be unique for the schema in which the dynamic table is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the dynamic table.",
		ForceNew:    true, // todo: can we make this not force new?
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the dynamic table.",
		ForceNew:    true, // todo: can we make this not force new?
	},
	"target_lag": {
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Description: "Specifies the target lag time for the dynamic table.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"maximum_duration": {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{"target_lag.downstream"},
					Description:   "Specifies the maximum target lag time for the dynamic table.",
				},
				"downstream": {
					Type:          schema.TypeBool,
					Optional:      true,
					ConflictsWith: []string{"target_lag.maximum_duration"},
					Description:   "Specifies whether the target lag time is downstream.",
				},
			},
		},
	},
	"warehouse": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse in which to create the dynamic table.",
		ForceNew:    true,
	},
	"query": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the query to use to populate the dynamic table.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the dynamic table.",
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
}

// DynamicTable returns a pointer to the resource representing a dynamic table.
func DynamicTable() *schema.Resource {
	return &schema.Resource{
		Create: CreateDynamicTable,
		Read:   ReadDynamicTable,
		Update: UpdateDynamicTable,
		Delete: DeleteDynamicTable,

		Schema: dynamicTableShema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadDynamicTable implements schema.ReadFunc.
func ReadDynamicTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	dynamicTable, err := client.DynamicTables.ShowByID(context.Background(), id)
	if err != nil {
		log.Printf("[DEBUG] dynamic table (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err := d.Set("name", dynamicTable.Name); err != nil {
		return err
	}
	if err := d.Set("database", dynamicTable.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema", dynamicTable.SchemaName); err != nil {
		return err
	}
	if err := d.Set("warehouse", dynamicTable.Warehouse); err != nil {
		return err
	}
	if err := d.Set("comment", dynamicTable.Comment); err != nil {
		return err
	}
	tl := map[string]interface{}{}
	if dynamicTable.TargetLag == "DOWNSTREAM" {
		tl["downstream"] = true
		if err := d.Set("target_lag", []interface{}{tl}); err != nil {
			return err
		}
	} else {
		tl["maximum_duration"] = dynamicTable.TargetLag
		if err := d.Set("target_lag", []interface{}{tl}); err != nil {
			return err
		}
	}
	text := dynamicTable.Text
	if strings.Contains(text, "OR REPLACE") {
		if err := d.Set("or_replace", true); err != nil {
			return err
		}
	} else {
		if err := d.Set("or_replace", false); err != nil {
			return err
		}
	}
	// trim up to " ..AS" and remove whitespace
	query := strings.TrimSpace(text[strings.Index(text, "AS")+3:])
	if err := d.Set("query", query); err != nil {
		return err
	}
	if err := d.Set("cluster_by", dynamicTable.ClusterBy); err != nil {
		return err
	}
	if err := d.Set("rows", dynamicTable.Rows); err != nil {
		return err
	}
	if err := d.Set("bytes", dynamicTable.Bytes); err != nil {
		return err
	}
	if err := d.Set("owner", dynamicTable.Owner); err != nil {
		return err
	}
	if err := d.Set("refresh_mode", string(dynamicTable.RefreshMode)); err != nil {
		return err
	}
	if err := d.Set("refresh_mode_reason", dynamicTable.RefreshModeReason); err != nil {
		return err
	}
	if err := d.Set("automatic_clustering", dynamicTable.AutomaticClustering); err != nil {
		return err
	}
	if err := d.Set("scheduling_state", string(dynamicTable.SchedulingState)); err != nil {
		return err
	}
	/*
		guides on time formatting
		https://docs.snowflake.com/en/user-guide/date-time-input-output
		https://pkg.go.dev/time
		note: format may depend on what the account parameter for TIMESTAMP_OUTPUT_FORMAT is set to. Perhaps we should return this as a string rather than a time.Time?
	*/
	if err := d.Set("last_suspended_on", dynamicTable.LastSuspendedOn.Format("2006-01-02T16:04:05.000 -0700")); err != nil {
		return err
	}
	if err := d.Set("is_clone", dynamicTable.IsClone); err != nil {
		return err
	}
	if err := d.Set("is_replica", dynamicTable.IsReplica); err != nil {
		return err
	}
	if err := d.Set("data_timestamp", dynamicTable.DataTimestamp.Format("2006-01-02T16:04:05.000 -0700")); err != nil {
		return err
	}
	return nil
}

func parseTargetLag(v interface{}) sdk.TargetLag {
	var result sdk.TargetLag
	tl := v.([]interface{})[0].(map[string]interface{})
	if v, ok := tl["maximum_duration"]; ok {
		result.MaximumDuration = sdk.String(v.(string))
	}
	if v, ok := tl["downstream"]; ok && v.(bool) {
		result.MaximumDuration = nil
		result.Downstream = sdk.Bool(v.(bool))
	}
	return result
}

// CreateDynamicTable implements schema.CreateFunc.
func CreateDynamicTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	warehouse := sdk.NewAccountObjectIdentifier(d.Get("warehouse").(string))
	tl := parseTargetLag(d.Get("target_lag"))
	query := d.Get("query").(string)

	request := sdk.NewCreateDynamicTableRequest(id, warehouse, tl, query)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		request.WithOrReplace(true)
	}
	if err := client.DynamicTables.Create(context.Background(), request); err != nil {
		return err
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadDynamicTable(d, meta)
}

// UpdateDynamicTable implements schema.UpdateFunc.
func UpdateDynamicTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	request := sdk.NewAlterDynamicTableRequest(id)

	set := sdk.NewDynamicTableSetRequest()
	if d.HasChange("target_lag") {
		tl := parseTargetLag(d.Get("target_lag"))
		set.WithTargetLag(tl)
	}

	if d.HasChange("warehouse") {
		warehouseName := d.Get("warehouse").(string)
		set.WithWarehouse(sdk.NewAccountObjectIdentifier(warehouseName))
	}

	request.WithSet(set)
	if err := client.DynamicTables.Alter(ctx, request); err != nil {
		return err
	}

	if d.HasChange("comment") {
		err := client.Comments.Set(ctx, &sdk.SetCommentOptions{
			ObjectType: sdk.ObjectTypeDynamicTable,
			ObjectName: id,
			Value:      sdk.String(d.Get("comment").(string)),
		})
		if err != nil {
			return err
		}
	}
	return ReadDynamicTable(d, meta)
}

// DeleteDynamicTable implements schema.DeleteFunc.
func DeleteDynamicTable(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	if err := client.DynamicTables.Drop(context.Background(), sdk.NewDropDynamicTableRequest(id)); err != nil {
		return err
	}
	d.SetId("")

	return nil
}
