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
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseSize), IgnoreChangeToCurrentSnowflakeValueInShow("size")),
		Description:      fmt.Sprintf("Specifies the size of the interactive warehouse. Valid values are (case-insensitive): %s. Note: removing the size from config will result in the resource recreation.", possibleValuesListed(sdk.ValidWarehouseSizesString)),
	},
	"max_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count"),
		Description:      "Specifies the maximum number of server clusters for the interactive warehouse.",
	},
	"min_cluster_count": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("min_cluster_count"),
		Description:      "Specifies the minimum number of server clusters for the interactive warehouse (only applies to multi-cluster warehouses).",
	},
	"auto_suspend": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend"),
		Description:      "Specifies the number of seconds of inactivity after which an interactive warehouse is automatically suspended.",
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
			Schema: schemas.ShowWarehouseParametersSchemaInteractive,
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
			customdiff.ForceNewIfChange("warehouse_size", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) != "" && new.(string) == ""
			}),
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
// FALLBACK_WAREHOUSE is a warehouse parameter (exposed in SHOW PARAMETERS IN WAREHOUSE, not as a
// SHOW WAREHOUSES column); unlike the int/bool warehouse parameters its value is the fallback
// warehouse name (an account-object identifier), or empty when it is not set.
func fallbackWarehouseFromParameters(parameters []*sdk.Parameter) string {
	for _, parameter := range parameters {
		if parameter.Key == string(sdk.WarehouseParameterFallbackWarehouse) {
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

	// Only the fields reconciled through handleExternalChangesToObjectInShow need seeding here: Read leaves
	// them unset on a fresh import because there is no previous SHOW output to compare against. Every other
	// attribute (name, comment, tables, fallback_warehouse, ...) is populated by the Read that Terraform
	// runs immediately after this importer, so setting them here too would be redundant.
	errs := errors.Join(
		d.Set("auto_resume", booleanStringFromBool(w.AutoResume)),
		setOptionalFromPtr(d, "auto_suspend", w.AutoSuspend),
		setOptionalFromPtr(d, "max_cluster_count", w.MaxClusterCount),
		setOptionalFromPtr(d, "min_cluster_count", w.MinClusterCount),
	)
	if w.Size != nil {
		errs = errors.Join(errs, d.Set("warehouse_size", string(*w.Size)))
	}
	if rm := w.ResourceMonitor.Name(); rm != "" {
		errs = errors.Join(errs, d.Set("resource_monitor", rm))
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
	if v, ok := d.GetOk("initially_suspended"); ok {
		req.WithInitiallySuspended(v.(bool))
	}
	if v, ok := d.GetOk("resource_monitor"); ok {
		req.WithResourceMonitor(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "warehouse_size", req.WithWarehouseSize, sdk.ToWarehouseSize),
		intAttributeCreateBuilder(d, "max_cluster_count", req.WithMaxClusterCount),
		intAttributeCreateBuilder(d, "min_cluster_count", req.WithMinClusterCount),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend", req.WithAutoSuspend),
		booleanStringAttributeCreateBuilder(d, "auto_resume", req.WithAutoResume),
		stringAttributeCreateBuilder(d, "comment", req.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
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

	// CREATE WAREHOUSE makes the newly created warehouse the current warehouse for the session
	// connection. For an interactive warehouse this forces a 5-second statement timeout, which would
	// then apply to the follow-up ALTER, the post-create Read, and any later statement that reuses the
	// same pooled connection (e.g. a subsequent apply or import). Capture the current warehouse first
	// and restore it right after creating, so DDL/metadata operations keep running on the original
	// warehouse rather than the interactive one.
	previousWarehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading current warehouse before creating interactive warehouse %s: %w", id.FullyQualifiedName(), err))
	}

	if err := client.Warehouses.CreateInteractive(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating interactive warehouse %s: %w", id.FullyQualifiedName(), err))
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if previousWarehouse != "" {
		if err := client.Sessions.UseWarehouse(ctx, sdk.NewUseWarehouseSessionRequest(sdk.NewAccountObjectIdentifier(previousWarehouse))); err != nil {
			return diag.FromErr(fmt.Errorf("error restoring warehouse %q after creating interactive warehouse %s: %w", previousWarehouse, id.FullyQualifiedName(), err))
		}
	}

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
			maxClusterCount := optionalIntOutputMapping(w.MaxClusterCount)
			minClusterCount := optionalIntOutputMapping(w.MinClusterCount)
			autoSuspend := optionalIntOutputMappingIntDefault(w.AutoSuspend)
			if err = handleExternalChangesToObjectInShow(
				d,
				outputMapping{"size", "warehouse_size", sizeStr, sizeVal, nil},
				outputMapping{"max_cluster_count", "max_cluster_count", maxClusterCount, maxClusterCount, nil},
				outputMapping{"min_cluster_count", "min_cluster_count", minClusterCount, minClusterCount, nil},
				outputMapping{"auto_suspend", "auto_suspend", autoSuspend, autoSuspend, nil},
				outputMapping{"auto_resume", "auto_resume", autoResumeVal, autoResumeStr, nil},
				outputMapping{"resource_monitor", "resource_monitor", w.ResourceMonitor.Name(), w.ResourceMonitor.Name(), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set("name", w.Name),
			d.Set("comment", w.Comment),
			d.Set("warehouse_type", effectiveType),
			d.Set("fallback_warehouse", fallbackWarehouse),
			d.Set("tables", tables),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseInteractiveToSchema(w)}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.WarehouseInteractiveParametersToSchema(warehouseParameters, providerCtx)}),
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
	}

	if err := errors.Join(
		intAttributeUpdate(d, "max_cluster_count", &set.MaxClusterCount, &unset.MaxClusterCount),
		intAttributeUpdate(d, "min_cluster_count", &set.MinClusterCount, &unset.MinClusterCount),
		intAttributeWithSpecialDefaultUpdate(d, "auto_suspend", &set.AutoSuspend, &unset.AutoSuspend),
		booleanStringAttributeUpdate(d, "auto_resume", &set.AutoResume, &unset.AutoResume),
		accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &set.ResourceMonitor, &unset.ResourceMonitor),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		accountObjectIdentifierAttributeUpdate(d, "fallback_warehouse", &set.FallbackWarehouse, &unset.FallbackWarehouse),
	); err != nil {
		return diag.FromErr(err)
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
