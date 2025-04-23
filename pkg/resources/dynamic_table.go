package resources

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var refreshModePattern = regexp.MustCompile(`refresh_mode = '(\w+)'`)

var dynamicTableSchema = map[string]*schema.Schema{
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
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the dynamic table.",
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
	},
	"query": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the query to use to populate the dynamic table.",
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the dynamic table.",
	},
	"refresh_mode": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      sdk.DynamicTableRefreshModeAuto,
		Description:  "INCREMENTAL to use incremental refreshes, FULL to recompute the whole table on every refresh, or AUTO to let Snowflake decide.",
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllDynamicRefreshModes), true),
		ForceNew:     true,
	},
	"initialize": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      sdk.DynamicTableInitializeOnCreate,
		Description:  "Initialize trigger for the dynamic table. Can only be set on creation. Available options are ON_CREATE and ON_SCHEDULE.",
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllDynamicTableInitializes), true),
		ForceNew:     true,
	},
	"created_on": {
		Type:        schema.TypeString,
		Description: "Time when this dynamic table was created.",
		Computed:    true,
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
		Description: "Displays ACTIVE for dynamic tables that are actively scheduling refreshes and SUSPENDED for suspended dynamic tables.",
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
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func DynamicTable() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.DynamicTables.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.DynamicTableResource), TrackingCreateWrapper(resources.DynamicTable, CreateDynamicTable)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.DynamicTableResource), TrackingReadWrapper(resources.DynamicTable, ReadDynamicTable)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.DynamicTableResource), TrackingUpdateWrapper(resources.DynamicTable, UpdateDynamicTable)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.DynamicTableResource), TrackingDeleteWrapper(resources.DynamicTable, deleteFunc)),

		Schema: dynamicTableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// ReadDynamicTable implements schema.ReadFunc.
func ReadDynamicTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	dynamicTable, err := client.DynamicTables.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query dynamic table. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Dynamic table id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", dynamicTable.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database", dynamicTable.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", dynamicTable.SchemaName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse", dynamicTable.Warehouse); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", dynamicTable.Comment); err != nil {
		return diag.FromErr(err)
	}
	tl := map[string]interface{}{}
	if dynamicTable.TargetLag == "DOWNSTREAM" {
		tl["downstream"] = true
		if err := d.Set("target_lag", []interface{}{tl}); err != nil {
			return diag.FromErr(err)
		}
	} else {
		tl["maximum_duration"] = dynamicTable.TargetLag
		if err := d.Set("target_lag", []interface{}{tl}); err != nil {
			return diag.FromErr(err)
		}
	}
	if strings.Contains(dynamicTable.Text, "OR REPLACE") {
		if err := d.Set("or_replace", true); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("or_replace", false); err != nil {
			return diag.FromErr(err)
		}
	}
	if strings.Contains(dynamicTable.Text, "initialize = 'ON_CREATE'") {
		if err := d.Set("initialize", "ON_CREATE"); err != nil {
			return diag.FromErr(err)
		}
	} else if strings.Contains(dynamicTable.Text, "initialize = 'ON_SCHEDULE'") {
		if err := d.Set("initialize", "ON_SCHEDULE"); err != nil {
			return diag.FromErr(err)
		}
	}
	m := refreshModePattern.FindStringSubmatch(dynamicTable.Text)
	if len(m) > 1 {
		if err := d.Set("refresh_mode", m[1]); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_on", dynamicTable.CreatedOn.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_by", dynamicTable.ClusterBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rows", dynamicTable.Rows); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("bytes", dynamicTable.Bytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner", dynamicTable.Owner); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("refresh_mode_reason", dynamicTable.RefreshModeReason); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("automatic_clustering", dynamicTable.AutomaticClustering); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scheduling_state", string(dynamicTable.SchedulingState)); err != nil {
		return diag.FromErr(err)
	}
	/*
		guides on time formatting
		https://docs.snowflake.com/en/user-guide/date-time-input-output
		https://pkg.go.dev/time
		note: format may depend on what the account parameter for TIMESTAMP_OUTPUT_FORMAT is set to. Perhaps we should return this as a string rather than a time.Time?
	*/
	if err := d.Set("last_suspended_on", dynamicTable.LastSuspendedOn.Format("2006-01-02T16:04:05.000 -0700")); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_clone", dynamicTable.IsClone); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_replica", dynamicTable.IsReplica); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_timestamp", dynamicTable.DataTimestamp.Format("2006-01-02T16:04:05.000 -0700")); err != nil {
		return diag.FromErr(err)
	}

	extractor := snowflake.NewViewSelectStatementExtractor(dynamicTable.Text)
	query, err := extractor.ExtractDynamicTable()
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("query", query); err != nil {
		return diag.FromErr(err)
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
func CreateDynamicTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

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
	if v, ok := d.GetOk("refresh_mode"); ok {
		request.WithRefreshMode(sdk.DynamicTableRefreshMode(v.(string)))
	}
	if v, ok := d.GetOk("initialize"); ok {
		request.WithInitialize(sdk.DynamicTableInitialize(v.(string)))
	}
	if err := client.DynamicTables.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadDynamicTable(ctx, d, meta)
}

// UpdateDynamicTable implements schema.UpdateFunc.
func UpdateDynamicTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	request := sdk.NewAlterDynamicTableRequest(id)

	runSet := false
	set := sdk.NewDynamicTableSetRequest()
	if d.HasChange("target_lag") {
		tl := parseTargetLag(d.Get("target_lag"))
		set.WithTargetLag(tl)
		runSet = true
	}

	if d.HasChange("warehouse") {
		warehouseName := d.Get("warehouse").(string)
		set.WithWarehouse(sdk.NewAccountObjectIdentifier(warehouseName))
		runSet = true
	}

	if runSet {
		request.WithSet(set)
		if err := client.DynamicTables.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("comment") {
		err := client.Comments.Set(ctx, &sdk.SetCommentOptions{
			ObjectType: sdk.ObjectTypeDynamicTable,
			ObjectName: id,
			Value:      sdk.String(d.Get("comment").(string)),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadDynamicTable(ctx, d, meta)
}
