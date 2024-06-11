package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var warehouseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the virtual warehouse; must be unique for your account.",
	},
	"warehouse_type": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(sdk.ValidWarehouseTypesString, true),
		Description:  fmt.Sprintf("Specifies warehouse type. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.ValidWarehouseTypesString)),
	},
	// TODO: handle forceNew instead of update
	"warehouse_size": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToWarehouseSize),
		Description:      fmt.Sprintf("Specifies the size of the virtual warehouse. Valid values are (case-insensitive): %s. Consult [warehouse documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for the details.", possibleValuesListed(sdk.ValidWarehouseSizesString)),
	},
	"max_cluster_count": {
		Type:         schema.TypeInt,
		Description:  "Specifies the maximum number of server clusters for the warehouse.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"min_cluster_count": {
		Type:         schema.TypeInt,
		Description:  "Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"scaling_policy": {
		Type:         schema.TypeString,
		Description:  fmt.Sprintf("Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.ValidWarehouseScalingPoliciesString)),
		Optional:     true,
		ValidateFunc: validation.StringInSlice(sdk.ValidWarehouseScalingPoliciesString, true),
	},
	"auto_suspend": {
		Type:         schema.TypeInt,
		Description:  "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Default:      -1,
	},
	"auto_resume": {
		Type:        schema.TypeBool,
		Description: "Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it.",
		Optional:    true,
	},
	"initially_suspended": {
		Type:        schema.TypeBool,
		Description: "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
		Optional:    true,
		ForceNew:    true,
	},
	"resource_monitor": {
		Type:             schema.TypeString,
		Description:      "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	// TODO: test setting empty comment
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the warehouse.",
	},
	"enable_query_acceleration": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to enable the query acceleration service for queries that rely on this warehouse for compute resources.",
	},
	"query_acceleration_max_scale_factor": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 100),
		Description:  "Specifies the maximum scale factor for leasing compute resources for query acceleration. The scale factor is used as a multiplier based on warehouse size.",
		Default:      -1,
	},
	strings.ToLower(string(sdk.ObjectParameterMaxConcurrencyLevel)): {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
		Description:  "Object parameter that specifies the concurrency level for SQL statements (i.e. queries and DML) executed by a warehouse.",
		Default:      -1,
	},
	strings.ToLower(string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds)): {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Description:  "Object parameter that specifies the time, in seconds, a SQL statement (query, DDL, DML, etc.) can be queued on a warehouse before it is canceled by the system.",
		Default:      -1,
	},
	strings.ToLower(string(sdk.ObjectParameterStatementTimeoutInSeconds)): {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 604800),
		Description:  "Specifies the time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system",
		Default:      -1,
	},
	showOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSE` for the given warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseSchema,
		},
	},
	parametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseParametersSchema,
		},
	},
}

// Warehouse returns a pointer to the resource representing a warehouse.
func Warehouse() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateWarehouse,
		UpdateContext: UpdateWarehouse,
		ReadContext:   GetReadWarehouseFunc(true),
		DeleteContext: DeleteWarehouse,
		Description:   "Resource used to manage warehouse objects. For more information, check [warehouse documentation](https://docs.snowflake.com/en/sql-reference/commands-warehouse).",

		Schema: warehouseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportWarehouse,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAttributeChanged(showOutputAttributeName, "warehouse_type"),
			ComputedIfAttributeChanged(showOutputAttributeName, "warehouse_size"),
			ComputedIfAttributeChanged(showOutputAttributeName, "max_cluster_count"),
			ComputedIfAttributeChanged(showOutputAttributeName, "min_cluster_count"),
			ComputedIfAttributeChanged(showOutputAttributeName, "scaling_policy"),
			ComputedIfAttributeChanged(showOutputAttributeName, "auto_suspend"),
			ComputedIfAttributeChanged(showOutputAttributeName, "auto_resume"),
			ComputedIfAttributeChanged(showOutputAttributeName, "initially_suspended"),
			ComputedIfAttributeChanged(showOutputAttributeName, "resource_monitor"),
			ComputedIfAttributeChanged(showOutputAttributeName, "comment"),
			ComputedIfAttributeChanged(showOutputAttributeName, "enable_query_acceleration"),
			ComputedIfAttributeChanged(showOutputAttributeName, "query_acceleration_max_scale_factor"),
			ComputedIfAttributeChanged(parametersAttributeName, strings.ToLower(string(sdk.ObjectParameterMaxConcurrencyLevel))),
			ComputedIfAttributeChanged(parametersAttributeName, strings.ToLower(string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds))),
			ComputedIfAttributeChanged(parametersAttributeName, strings.ToLower(string(sdk.ObjectParameterStatementTimeoutInSeconds))),
		),

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v092WarehouseSizeStateUpgrader,
			},
		},
	}
}

func ImportWarehouse(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting warehouse import")
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	w, err := client.Warehouses.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", w.Name); err != nil {
		return nil, err
	}
	if err = d.Set("warehouse_type", w.Type); err != nil {
		return nil, err
	}
	if err = d.Set("warehouse_size", w.Size); err != nil {
		return nil, err
	}
	if err = d.Set("max_cluster_count", w.MaxClusterCount); err != nil {
		return nil, err
	}
	if err = d.Set("min_cluster_count", w.MinClusterCount); err != nil {
		return nil, err
	}
	if err = d.Set("scaling_policy", w.ScalingPolicy); err != nil {
		return nil, err
	}
	if err = d.Set("auto_suspend", w.AutoSuspend); err != nil {
		return nil, err
	}
	if err = d.Set("auto_resume", w.AutoResume); err != nil {
		return nil, err
	}
	if err = d.Set("resource_monitor", w.ResourceMonitor); err != nil {
		return nil, err
	}
	if err = d.Set("comment", w.Comment); err != nil {
		return nil, err
	}
	if err = d.Set("enable_query_acceleration", w.EnableQueryAcceleration); err != nil {
		return nil, err
	}
	if err = d.Set("query_acceleration_max_scale_factor", w.QueryAccelerationMaxScaleFactor); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// CreateWarehouse implements schema.CreateFunc.
func CreateWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	createOptions := &sdk.CreateWarehouseOptions{}

	if v, ok := d.GetOk("warehouse_type"); ok {
		warehouseType, err := sdk.ToWarehouseType(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOptions.WarehouseType = &warehouseType
	}
	if v, ok := d.GetOk("warehouse_size"); ok {
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOptions.WarehouseSize = &size
	}
	if v, ok := d.GetOk("max_cluster_count"); ok {
		createOptions.MaxClusterCount = sdk.Int(v.(int))
	}
	if v, ok := d.GetOk("min_cluster_count"); ok {
		createOptions.MinClusterCount = sdk.Int(v.(int))
	}
	if v, ok := d.GetOk("scaling_policy"); ok {
		scalingPolicy, err := sdk.ToScalingPolicy(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOptions.ScalingPolicy = &scalingPolicy
	}
	if v := d.Get("auto_suspend").(int); v != -1 {
		createOptions.AutoSuspend = sdk.Int(v)
	}
	if v, ok := d.GetOk("auto_resume"); ok {
		createOptions.AutoResume = sdk.Bool(v.(bool))
	}
	if v, ok := d.GetOk("initially_suspended"); ok {
		createOptions.InitiallySuspended = sdk.Bool(v.(bool))
	}
	if v, ok := d.GetOk("resource_monitor"); ok {
		createOptions.ResourceMonitor = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	if v, ok := d.GetOk("comment"); ok {
		createOptions.Comment = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("enable_query_acceleration"); ok {
		createOptions.EnableQueryAcceleration = sdk.Bool(v.(bool))
	}
	if v := d.Get("query_acceleration_max_scale_factor").(int); v != -1 {
		createOptions.QueryAccelerationMaxScaleFactor = sdk.Int(v)
	}
	if v := d.Get("max_concurrency_level").(int); v != -1 {
		createOptions.MaxConcurrencyLevel = sdk.Int(v)
	}
	if v := d.Get("statement_queued_timeout_in_seconds").(int); v != -1 {
		createOptions.StatementQueuedTimeoutInSeconds = sdk.Int(v)
	}
	if v := d.Get("statement_timeout_in_seconds").(int); v != -1 {
		createOptions.StatementTimeoutInSeconds = sdk.Int(v)
	}

	err := client.Warehouses.Create(ctx, id, createOptions)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return GetReadWarehouseFunc(false)(ctx, d, meta)
}

func GetReadWarehouseFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

		w, err := client.Warehouses.ShowByID(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		warehouseParameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
			In: &sdk.ParametersIn{
				Warehouse: id,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObject(d, func(result map[string]any) error {
				if result["type"].(string) != string(w.Type) {
					if err = d.Set("warehouse_type", w.Type); err != nil {
						return err
					}
				}
				if result["size"].(string) != string(w.Size) {
					if err = d.Set("warehouse_size", w.Size); err != nil {
						return err
					}
				}
				if result["max_cluster_count"].(int) != w.MaxClusterCount {
					if err = d.Set("max_cluster_count", w.MaxClusterCount); err != nil {
						return err
					}
				}
				if result["min_cluster_count"].(int) != w.MinClusterCount {
					if err = d.Set("min_cluster_count", w.MinClusterCount); err != nil {
						return err
					}
				}
				if result["scaling_policy"].(string) != string(w.ScalingPolicy) {
					if err = d.Set("scaling_policy", w.ScalingPolicy); err != nil {
						return err
					}
				}
				if result["auto_suspend"].(int) != w.AutoSuspend {
					if err = d.Set("auto_suspend", w.AutoSuspend); err != nil {
						return err
					}
				}
				if result["auto_resume"].(bool) != w.AutoResume {
					if err = d.Set("auto_resume", w.AutoResume); err != nil {
						return err
					}
				}
				if sdk.NewAccountIdentifierFromFullyQualifiedName(result["resource_monitor"].(string)).FullyQualifiedName() != sdk.NewAccountIdentifierFromFullyQualifiedName(w.ResourceMonitor).FullyQualifiedName() {
					if err = d.Set("resource_monitor", w.ResourceMonitor); err != nil {
						return err
					}
				}
				if result["comment"].(string) != w.Comment {
					if err = d.Set("comment", w.Comment); err != nil {
						return err
					}
				}
				if result["enable_query_acceleration"].(bool) != w.EnableQueryAcceleration {
					if err = d.Set("enable_query_acceleration", w.EnableQueryAcceleration); err != nil {
						return err
					}
				}
				if result["query_acceleration_max_scale_factor"].(int) != w.QueryAccelerationMaxScaleFactor {
					if err = d.Set("query_acceleration_max_scale_factor", w.QueryAccelerationMaxScaleFactor); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				return diag.FromErr(err)
			}

			if err = markChangedParameters(sdk.WarehouseParameters, warehouseParameters, d, sdk.ParameterTypeWarehouse); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = d.Set(showOutputAttributeName, []map[string]any{schemas.WarehouseToSchema(w)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(parametersAttributeName, []map[string]any{schemas.WarehouseParametersToSchema(warehouseParameters)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

// UpdateWarehouse implements schema.UpdateFunc.
func UpdateWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	// Change name separately
	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			NewName: &newId,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	// Batch SET operations and UNSET operations
	set := sdk.WarehouseSet{}
	unset := sdk.WarehouseUnset{}
	if d.HasChange("warehouse_type") {
		if v, ok := d.GetOk("warehouse_type"); ok {
			warehouseType, err := sdk.ToWarehouseType(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WarehouseType = &warehouseType
		} else {
			// TODO [SNOW-1473453]: UNSET of type does not work
			//unset.WarehouseType = sdk.Bool(true)
			set.WarehouseType = &sdk.WarehouseTypeStandard
		}
	}
	if d.HasChange("warehouse_size") {
		n := d.Get("warehouse_size").(string)
		// TODO [this PR]: get rid of that part (replace with force new for this parameter)
		if n == "" {
			n = string(sdk.WarehouseSizeXSmall)
		}
		size, err := sdk.ToWarehouseSize(n)
		if err != nil {
			return diag.FromErr(err)
		}
		set.WarehouseSize = &size
	}
	if d.HasChange("max_cluster_count") {
		if v, ok := d.GetOk("max_cluster_count"); ok {
			set.MaxClusterCount = sdk.Int(v.(int))
		} else {
			unset.MaxClusterCount = sdk.Bool(true)
		}
	}
	if d.HasChange("min_cluster_count") {
		if v, ok := d.GetOk("min_cluster_count"); ok {
			set.MinClusterCount = sdk.Int(v.(int))
		} else {
			unset.MinClusterCount = sdk.Bool(true)
		}
	}
	if d.HasChange("scaling_policy") {
		if v, ok := d.GetOk("scaling_policy"); ok {
			scalingPolicy, err := sdk.ToScalingPolicy(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.ScalingPolicy = &scalingPolicy
		} else {
			unset.ScalingPolicy = sdk.Bool(true)
		}
	}
	if d.HasChange("auto_suspend") {
		if v := d.Get("auto_suspend").(int); v != -1 {
			set.AutoSuspend = sdk.Int(v)
		} else {
			unset.AutoSuspend = sdk.Bool(true)
		}
	}
	if d.HasChange("auto_resume") {
		if v, ok := d.GetOk("auto_resume"); ok {
			set.AutoResume = sdk.Bool(v.(bool))
		} else {
			unset.AutoResume = sdk.Bool(true)
		}
	}
	if d.HasChange("resource_monitor") {
		if v, ok := d.GetOk("resource_monitor"); ok {
			set.ResourceMonitor = sdk.NewAccountObjectIdentifier(v.(string))
		} else {
			unset.ResourceMonitor = sdk.Bool(true)
		}
	}
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.Comment = sdk.String(v.(string))
		} else {
			unset.Comment = sdk.Bool(true)
		}
	}
	if d.HasChange("enable_query_acceleration") {
		if v, ok := d.GetOk("enable_query_acceleration"); ok {
			set.EnableQueryAcceleration = sdk.Bool(v.(bool))
		} else {
			unset.EnableQueryAcceleration = sdk.Bool(true)
		}
	}
	if d.HasChange("query_acceleration_max_scale_factor") {
		if v := d.Get("query_acceleration_max_scale_factor").(int); v != -1 {
			set.QueryAccelerationMaxScaleFactor = sdk.Int(v)
		} else {
			unset.QueryAccelerationMaxScaleFactor = sdk.Bool(true)
		}
	}
	if d.HasChange("max_concurrency_level") {
		if v := d.Get("max_concurrency_level").(int); v != -1 {
			set.MaxConcurrencyLevel = sdk.Int(v)
		} else {
			unset.MaxConcurrencyLevel = sdk.Bool(true)
		}
	}
	if d.HasChange("statement_queued_timeout_in_seconds") {
		if v := d.Get("statement_queued_timeout_in_seconds").(int); v != -1 {
			set.StatementQueuedTimeoutInSeconds = sdk.Int(v)
		} else {
			unset.StatementQueuedTimeoutInSeconds = sdk.Bool(true)
		}
	}
	if d.HasChange("statement_timeout_in_seconds") {
		if v := d.Get("statement_timeout_in_seconds").(int); v != -1 {
			set.StatementTimeoutInSeconds = sdk.Int(v)
		} else {
			unset.StatementTimeoutInSeconds = sdk.Bool(true)
		}
	}

	// Apply SET and UNSET changes
	if (set != sdk.WarehouseSet{}) {
		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Set: &set,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if (unset != sdk.WarehouseUnset{}) {
		err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Unset: &unset,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return GetReadWarehouseFunc(false)(ctx, d, meta)
}

// DeleteWarehouse implements schema.DeleteFunc.
func DeleteWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.Warehouses.Drop(ctx, id, &sdk.DropWarehouseOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
