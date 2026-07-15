package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var warehouseInteractiveSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the interactive warehouse; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"warehouse_size": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          string(sdk.WarehouseSizeXSmall),
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseSize), IgnoreChangeToCurrentSnowflakeValueInShow("size")),
		Description:      fmt.Sprintf("Specifies the size of the interactive warehouse. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.ValidWarehouseSizesString)),
	},
	"max_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count"),
		Description:      "Specifies the maximum number of server clusters for the interactive warehouse. Snowflake always assigns a value, so this field is computed when not set.",
	},
	"min_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("min_cluster_count"),
		Description:      "Specifies the minimum number of server clusters for the interactive warehouse (only applies to multi-cluster warehouses). Snowflake always assigns a value, so this field is computed when not set.",
	},
	"auto_suspend": {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend"),
		Description:      "Specifies the number of seconds of inactivity after which an interactive warehouse is automatically suspended. Snowflake always assigns a value, so this field is computed when not set.",
	},
	"auto_resume": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume"),
		Description:      booleanStringFieldDescription("Specifies whether to automatically resume an interactive warehouse when a SQL statement (e.g. query) is submitted to it."),
		Default:          BooleanDefault,
	},
	"initially_suspended": {
		Type:             schema.TypeBool,
		Optional:         true,
		DiffSuppressFunc: IgnoreAfterCreation,
		Description:      "Specifies whether the interactive warehouse is created initially in the ‘Suspended’ state.",
	},
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor")),
		Description:      relatedResourceDescription("Specifies the name of a resource monitor that is explicitly assigned to the interactive warehouse.", resources.ResourceMonitor),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the interactive warehouse.",
	},
	"tables": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("tables"),
		Description:      "Specifies the fully qualified names of the tables associated with the interactive warehouse. Changes are applied incrementally (ADD TABLES / DROP TABLES) rather than by full re-association.",
	},
	"fallback_warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      relatedResourceDescription("Specifies the name of the fallback warehouse for the interactive warehouse.", resources.Warehouse),
	},
	"warehouse_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the type for the interactive warehouse. This field is used for checking external changes and recreating the resource if needed.",
	},
	strings.ToLower(string(sdk.WarehouseParameterMaxConcurrencyLevel)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Object parameter that specifies the concurrency level for SQL statements (i.e. queries and DML) executed by an interactive warehouse.",
	},
	strings.ToLower(string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Object parameter that specifies the time, in seconds, a SQL statement (query, DDL, DML, etc.) can be queued on an interactive warehouse before it is canceled by the system.",
	},
	strings.ToLower(string(sdk.WarehouseParameterStatementTimeoutInSeconds)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 604800)),
		Description:      "Specifies the time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given interactive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseSchemaInteractive,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given interactive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func WarehouseInteractive() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Warehouses.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.WarehouseInteractiveResource), TrackingCreateWrapper(resources.WarehouseInteractive, CreateWarehouseInteractive)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.WarehouseInteractiveResource), TrackingReadWrapper(resources.WarehouseInteractive, ReadWarehouseInteractiveFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.WarehouseInteractiveResource), TrackingUpdateWrapper(resources.WarehouseInteractive, UpdateWarehouseInteractive)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.WarehouseInteractiveResource), TrackingDeleteWrapper(resources.WarehouseInteractive, deleteFunc)),
		Description:   "Resource used to manage interactive warehouse objects. Interactive warehouses are optimized for low-latency, high-concurrency queries against a defined set of tables. For more information, check [interactive warehouse documentation](https://docs.snowflake.com/en/user-guide/warehouses-interactive).",

		Schema: warehouseInteractiveSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.WarehouseInteractive, ImportWarehouseInteractive),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.WarehouseInteractive, customdiff.All(
			ComputedIfAnyAttributeChanged(warehouseInteractiveSchema, ShowOutputAttributeName, "name", "warehouse_size", "max_cluster_count", "min_cluster_count", "auto_suspend", "auto_resume", "resource_monitor", "comment", "tables"),
			ComputedIfAnyAttributeChanged(
				warehouseInteractiveSchema, ParametersAttributeName,
				strings.ToLower(string(sdk.WarehouseParameterMaxConcurrencyLevel)),
				strings.ToLower(string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)),
				strings.ToLower(string(sdk.WarehouseParameterStatementTimeoutInSeconds)),
			),
			ComputedIfAnyAttributeChanged(warehouseInteractiveSchema, FullyQualifiedNameAttributeName, "name"),
			ParametersCustomDiff(
				warehouseParametersProvider,
				parameter[sdk.AccountParameter]{sdk.AccountParameterMaxConcurrencyLevel, valueTypeInt, sdk.ParameterTypeWarehouse},
				parameter[sdk.AccountParameter]{sdk.AccountParameterStatementQueuedTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeWarehouse},
				parameter[sdk.AccountParameter]{sdk.AccountParameterStatementTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeWarehouse},
			),
			// Snowflake does not allow changing WAREHOUSE_TYPE via ALTER (to or from INTERACTIVE),
			// so if the underlying object is no longer interactive the only way to reconcile is to
			// recreate it.
			RecreateWhenResourceTypeChangedExternally("warehouse_type", sdk.WarehouseTypeInteractive, sdk.ToWarehouseType),
		)),
		Timeouts: defaultTimeouts,
	}
}

// parseTablesSet converts a *schema.Set of fully-qualified table identifier strings into
// a slice of sdk.SchemaObjectIdentifier.
func parseTablesSet(raw *schema.Set) ([]sdk.SchemaObjectIdentifier, error) {
	tables := make([]sdk.SchemaObjectIdentifier, 0, raw.Len())
	for _, v := range raw.List() {
		id, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return nil, err
		}
		tables = append(tables, id)
	}
	return tables, nil
}

// fallbackWarehouseFromParameters returns the FALLBACK_WAREHOUSE value from SHOW PARAMETERS output.
// FALLBACK_WAREHOUSE is exposed as a warehouse parameter (not a SHOW WAREHOUSES column); the value
// is the fallback warehouse name, or empty when it is not set.
func fallbackWarehouseFromParameters(parameters []*sdk.Parameter) string {
	for _, parameter := range parameters {
		if parameter.Key == "FALLBACK_WAREHOUSE" {
			return parameter.Value
		}
	}
	return ""
}

func ImportWarehouseInteractive(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	w, err := client.Warehouses.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !w.IsInteractiveWarehouse() {
		return nil, fmt.Errorf("warehouse %s is not an interactive warehouse; use snowflake_warehouse or snowflake_warehouse_adaptive instead", id.FullyQualifiedName())
	}

	tables := make([]string, len(w.Tables))
	for i, table := range w.Tables {
		tables[i] = table.FullyQualifiedName()
	}

	errs := errors.Join(
		d.Set("name", id.Name()),
		d.Set("comment", w.Comment),
		setOptionalFromPtr(d, "auto_suspend", w.AutoSuspend),
		setOptionalFromPtr(d, "max_cluster_count", w.MaxClusterCount),
		setOptionalFromPtr(d, "min_cluster_count", w.MinClusterCount),
		d.Set("tables", tables),
	)
	if w.Size != nil {
		errs = errors.Join(errs, d.Set("warehouse_size", string(*w.Size)))
	}

	fallbackWarehouse, err := client.Warehouses.ShowParameters(ctx, id)
	if err != nil {
		return nil, err
	}
	if fw := fallbackWarehouseFromParameters(fallbackWarehouse); fw != "" {
		errs = errors.Join(errs, d.Set("fallback_warehouse", fw))
	}
	if errs != nil {
		return nil, errs
	}

	return []*schema.ResourceData{d}, nil
}

func CreateWarehouseInteractive(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateInteractiveWarehouseRequest(id)

	if v, ok := d.GetOk("tables"); ok {
		tables, err := parseTablesSet(v.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithTables(tables)
	}
	if v := d.Get("warehouse_size").(string); v != "" {
		size, err := sdk.ToWarehouseSize(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithWarehouseSize(size)
	}
	if v, ok := d.GetOk("max_cluster_count"); ok {
		req.WithMaxClusterCount(v.(int))
	}
	if v, ok := d.GetOk("min_cluster_count"); ok {
		req.WithMinClusterCount(v.(int))
	}
	if v, ok := d.GetOk("auto_suspend"); ok {
		req.WithAutoSuspend(v.(int))
	}
	if v := d.Get("auto_resume").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithAutoResume(parsed)
	}
	if v, ok := d.GetOk("initially_suspended"); ok {
		req.WithInitiallySuspended(v.(bool))
	}
	if v, ok := d.GetOk("resource_monitor"); ok {
		req.WithResourceMonitor(sdk.NewAccountObjectIdentifier(v.(string)))
	}
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}
	if v := GetConfigPropertyAsPointerAllowingZeroValue[int](d, "max_concurrency_level"); v != nil {
		req.WithMaxConcurrencyLevel(*v)
	}
	if v := GetConfigPropertyAsPointerAllowingZeroValue[int](d, "statement_queued_timeout_in_seconds"); v != nil {
		req.WithStatementQueuedTimeoutInSeconds(*v)
	}
	if v := GetConfigPropertyAsPointerAllowingZeroValue[int](d, "statement_timeout_in_seconds"); v != nil {
		req.WithStatementTimeoutInSeconds(*v)
	}

	if err := client.Warehouses.CreateInteractive(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating interactive warehouse %s: %w", id.FullyQualifiedName(), err))
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// FALLBACK_WAREHOUSE is not a CREATE property; set it via a follow-up ALTER.
	if v, ok := d.GetOk("fallback_warehouse"); ok {
		set := sdk.NewWarehouseSetRequest().WithFallbackWarehouse(sdk.NewAccountObjectIdentifier(v.(string)))
		if err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting fallback warehouse for interactive warehouse %s: %w", id.FullyQualifiedName(), err))
		}
	}

	return ReadWarehouseInteractiveFunc(false)(ctx, d, meta)
}

func ReadWarehouseInteractiveFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		providerCtx := meta.(*provider.Context)
		client := providerCtx.Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		w, err := client.Warehouses.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Failed to query interactive warehouse. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Interactive warehouse id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		warehouseParameters, err := client.Warehouses.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		// Snowflake reports interactive warehouses through the type column (type = INTERACTIVE)
		effectiveType := string(sdk.WarehouseTypeStandard)
		if w.IsInteractiveWarehouse() {
			effectiveType = string(sdk.WarehouseTypeInteractive)
		}

		tables := make([]string, len(w.Tables))
		for i, table := range w.Tables {
			tables[i] = table.FullyQualifiedName()
		}

		// FALLBACK_WAREHOUSE is exposed as a warehouse parameter, not a SHOW WAREHOUSES column.
		fallbackWarehouse := fallbackWarehouseFromParameters(warehouseParameters)

		if withExternalChangesMarking {
			sizeVal, sizeStr := optionalStringOutputMapping(w.Size)
			autoResumeVal, autoResumeStr := w.AutoResume, booleanStringFromBool(w.AutoResume)
			if err = handleExternalChangesToObjectInShow(
				d,
				outputMapping{"size", "warehouse_size", sizeStr, sizeVal, nil},
				outputMapping{"auto_resume", "auto_resume", autoResumeVal, autoResumeStr, nil},
				outputMapping{"resource_monitor", "resource_monitor", w.ResourceMonitor.Name(), w.ResourceMonitor.Name(), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, warehouseInteractiveSchema, []string{
			"warehouse_size",
			"auto_resume",
			"resource_monitor",
		}); err != nil {
			return diag.FromErr(err)
		}

		// auto_suspend, max_cluster_count and min_cluster_count are always assigned by Snowflake for
		// interactive warehouses, so they are computed and driven directly from the SHOW output.
		autoSuspendValue := 0
		if w.AutoSuspend != nil {
			autoSuspendValue = *w.AutoSuspend
		}
		maxClusterCountValue := 0
		if w.MaxClusterCount != nil {
			maxClusterCountValue = *w.MaxClusterCount
		}
		minClusterCountValue := 0
		if w.MinClusterCount != nil {
			minClusterCountValue = *w.MinClusterCount
		}

		errs := errors.Join(
			d.Set("name", w.Name),
			d.Set("comment", w.Comment),
			d.Set("warehouse_type", effectiveType),
			d.Set("auto_suspend", autoSuspendValue),
			d.Set("max_cluster_count", maxClusterCountValue),
			d.Set("min_cluster_count", minClusterCountValue),
			d.Set("fallback_warehouse", fallbackWarehouse),
			d.Set("tables", tables),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseInteractiveToSchema(w)}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.WarehouseParametersToSchema(warehouseParameters, providerCtx)}),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		if diags := handleWarehouseParameterRead(d, warehouseParameters); diags != nil {
			return diags
		}

		return nil
	}
}

func UpdateWarehouseInteractive(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		if err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// Tables are applied incrementally: compute the ADD/DROP delta from the change set.
	if d.HasChange("tables") {
		o, n := d.GetChange("tables")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		toAdd, err := parseTablesSet(newSet.Difference(oldSet))
		if err != nil {
			return diag.FromErr(err)
		}
		toDrop, err := parseTablesSet(oldSet.Difference(newSet))
		if err != nil {
			return diag.FromErr(err)
		}

		if len(toAdd) > 0 {
			if err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithAddTables(toAdd)); err != nil {
				return diag.FromErr(fmt.Errorf("error adding tables to interactive warehouse %s: %w", id.FullyQualifiedName(), err))
			}
		}
		if len(toDrop) > 0 {
			if err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithDropTables(toDrop)); err != nil {
				return diag.FromErr(fmt.Errorf("error dropping tables from interactive warehouse %s: %w", id.FullyQualifiedName(), err))
			}
		}
	}

	set := sdk.NewWarehouseSetRequest()
	unset := sdk.NewWarehouseUnsetRequest()

	if d.HasChange("warehouse_size") {
		size, err := sdk.ToWarehouseSize(d.Get("warehouse_size").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		set.WithWarehouseSize(size)
		set.WithWaitForCompletion(true)
	}
	if d.HasChange("max_cluster_count") {
		if v, ok := d.GetOk("max_cluster_count"); ok {
			set.WithMaxClusterCount(v.(int))
		} else {
			unset.WithMaxClusterCount(true)
		}
	}
	if d.HasChange("min_cluster_count") {
		if v, ok := d.GetOk("min_cluster_count"); ok {
			set.WithMinClusterCount(v.(int))
		} else {
			unset.WithMinClusterCount(true)
		}
	}
	if d.HasChange("auto_suspend") {
		if v, ok := d.GetOk("auto_suspend"); ok {
			set.WithAutoSuspend(v.(int))
		}
	}
	if d.HasChange("auto_resume") {
		if v := d.Get("auto_resume").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithAutoResume(parsed)
		} else {
			unset.WithAutoResume(true)
		}
	}
	if d.HasChange("resource_monitor") {
		if v, ok := d.GetOk("resource_monitor"); ok {
			set.WithResourceMonitor(sdk.NewAccountObjectIdentifier(v.(string)))
		} else {
			unset.WithResourceMonitor(true)
		}
	}
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.WithComment(v.(string))
		} else {
			unset.WithComment(true)
		}
	}
	if d.HasChange("fallback_warehouse") {
		if v, ok := d.GetOk("fallback_warehouse"); ok {
			set.WithFallbackWarehouse(sdk.NewAccountObjectIdentifier(v.(string)))
		} else {
			unset.WithFallbackWarehouse(true)
		}
	}

	if diags := JoinDiags(
		handleParameterUpdate(d, sdk.WarehouseParameterMaxConcurrencyLevel, &set.MaxConcurrencyLevel, &unset.MaxConcurrencyLevel),
		handleParameterUpdate(d, sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, sdk.WarehouseParameterStatementTimeoutInSeconds, &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
	); diags != nil {
		return diags
	}

	if *set != *sdk.NewWarehouseSetRequest() {
		if err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting interactive warehouse properties: %w", err))
		}
	}
	if *unset != *sdk.NewWarehouseUnsetRequest() {
		if err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting interactive warehouse properties: %w", err))
		}
	}

	return ReadWarehouseInteractiveFunc(false)(ctx, d, meta)
}
