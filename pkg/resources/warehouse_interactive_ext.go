package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	warehouseInteractiveCreateExt        = TrackingCreateWrapper(resources.WarehouseInteractive, CreateWarehouseInteractive)
	warehouseInteractiveReadExt          = TrackingReadWrapper(resources.WarehouseInteractive, ReadWarehouseInteractiveFunc(true))
	warehouseInteractiveUpdateExt        = TrackingUpdateWrapper(resources.WarehouseInteractive, UpdateWarehouseInteractive)
	warehouseInteractiveDescriptionExt   = "Resource used to manage interactive warehouse objects. Interactive warehouses are optimized for low-latency, high-concurrency queries against a defined set of tables. For more information, check [interactive warehouse documentation](https://docs.snowflake.com/en/user-guide/warehouses-interactive)."
	warehouseInteractiveCustomizeDiffExt = TrackingCustomDiffWrapper(resources.WarehouseInteractive, customdiff.All(
		ComputedIfAnyAttributeChanged(warehouseInteractiveSchema, ShowOutputAttributeName, "name", "warehouse_size", "max_cluster_count", "min_cluster_count", "auto_suspend", "auto_resume", "resource_monitor", "comment", "tables"),
		ComputedIfAnyAttributeChanged(
			warehouseInteractiveParametersAttributesSchema, ParametersAttributeName,
			warehouseInteractiveParametersFieldNames...,
		),
		ComputedIfAnyAttributeChanged(warehouseInteractiveSchema, FullyQualifiedNameAttributeName, "name"),
		// A resize between two concrete sizes is applied in place via suspend/alter/resume (see
		// UpdateWarehouseInteractive), mirroring how a regular warehouse handles it. Only removing
		// the size forces replacement: Snowflake has no UNSET WAREHOUSE_SIZE, so recreation is the
		// only way to fall back to the account default.
		customdiff.ForceNewIfChange("warehouse_size", func(ctx context.Context, old, new, meta any) bool {
			return old.(string) != "" && new.(string) == ""
		}),
		warehouseInteractiveParametersCustomDiff,
		HandleWarehouseExternalTypeChange(sdk.WarehouseTypeInteractive),
	))
	warehouseInteractiveImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.WarehouseInteractive, ImportWarehouseInteractive),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

// interactiveWarehouseDefaultAutoSuspend is the AUTO_SUSPEND value Snowflake applies when an
// interactive warehouse is created without one. When auto_suspend is removed from config we re-apply
// it with SET (rather than UNSET) so that removal matches CREATE-time behavior: UNSET AUTO_SUSPEND
// yields NULL and additionally misbehaves (see SNOW-1473453). This mirrors the regular warehouse.
const interactiveWarehouseDefaultAutoSuspend = 86400

type warehouseInteractiveAlterRequestsExt struct {
	set   *sdk.WarehouseInteractiveSetRequest
	unset *sdk.WarehouseInteractiveUnsetRequest
}

func warehouseInteractiveParseIdFromConfigExt(d *schema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
}

func warehouseInteractiveNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateInteractiveWarehouseRequest, error) {
	return sdk.NewCreateInteractiveWarehouseRequest(id), nil
}

func warehouseInteractiveApplyCreateOptionalsExt(d *schema.ResourceData, req *sdk.CreateInteractiveWarehouseRequest) error {
	return errors.Join(
		setValueCreate(d, "tables", &req.Tables, interactiveWarehouseTable),
		boolAttributeCreateBuilder(d, "initially_suspended", req.WithInitiallySuspended),
		accountObjectIdentifierAttributeCreateBuilder(d, "resource_monitor", req.WithResourceMonitor),
		attributeMappedValueCreateBuilder(d, "warehouse_size", req.WithWarehouseSize, sdk.ToWarehouseSize),
		intAttributeCreateBuilder(d, "max_cluster_count", req.WithMaxClusterCount),
		intAttributeCreateBuilder(d, "min_cluster_count", req.WithMinClusterCount),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend", req.WithAutoSuspend),
		booleanStringAttributeCreateBuilder(d, "auto_resume", req.WithAutoResume),
		stringAttributeCreateBuilder(d, "comment", req.WithComment),
	)
}

func warehouseInteractiveHandleExternalChangesExt(d *schema.ResourceData, w *sdk.WarehouseInteractive, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	return handleExternalChangesToObjectInShow(d, warehouseInteractiveOutputMappingsExt(w)...)
}

func warehouseInteractiveOutputMappingsExt(w *sdk.WarehouseInteractive) []outputMapping {
	sizeVal, sizeStr := optionalStringOutputMapping(w.Size)
	autoResumeVal, autoResumeStr := w.AutoResume, booleanStringFromBool(w.AutoResume)
	maxClusterCount := optionalIntOutputMapping(w.MaxClusterCount)
	minClusterCount := optionalIntOutputMapping(w.MinClusterCount)
	autoSuspend := optionalIntOutputMappingIntDefault(w.AutoSuspend)
	return []outputMapping{
		{"size", "warehouse_size", sizeStr, sizeVal, nil},
		{"max_cluster_count", "max_cluster_count", maxClusterCount, maxClusterCount, nil},
		{"min_cluster_count", "min_cluster_count", minClusterCount, minClusterCount, nil},
		{"auto_suspend", "auto_suspend", autoSuspend, autoSuspend, nil},
		{"auto_resume", "auto_resume", autoResumeVal, autoResumeStr, nil},
		{"resource_monitor", "resource_monitor", w.ResourceMonitor.Name(), w.ResourceMonitor.Name(), nil},
	}
}

func warehouseInteractiveSetConfigFieldsExt(d *schema.ResourceData, w *sdk.WarehouseInteractive) error {
	// Snowflake reports interactive warehouses through the type column (type = INTERACTIVE).
	effectiveType := string(sdk.WarehouseTypeStandard)
	if w.Type == sdk.WarehouseTypeInteractive {
		effectiveType = string(sdk.WarehouseTypeInteractive)
	}

	tables := make([]string, len(w.Tables))
	for i, table := range w.Tables {
		tables[i] = table.FullyQualifiedName()
	}

	return errors.Join(
		d.Set("name", w.Name),
		d.Set("comment", w.Comment),
		d.Set("warehouse_type", effectiveType),
		d.Set("tables", tables),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseInteractiveToSchema(w)}),
	)
}

func warehouseInteractiveHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.AccountObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	newId, err := renameAccountLevelResource(ctx, d, *id, func(ctx context.Context, id, newId sdk.AccountObjectIdentifier) error {
		return client.Warehouses.AlterInteractive(ctx, sdk.NewAlterInteractiveWarehouseRequest(id).WithRenameTo(newId))
	})
	if err != nil {
		return diag.FromErr(err)
	}
	*id = newId
	return nil
}

func warehouseInteractiveBeforeAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	// Tables are applied incrementally: compute the ADD/DROP delta from the change set.
	if d.HasChange("tables") {
		o, n := d.GetChange("tables")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		toAdd, err := collections.MapErr(newSet.Difference(oldSet).List(), interactiveWarehouseTable)
		if err != nil {
			return diag.FromErr(err)
		}
		toDrop, err := collections.MapErr(oldSet.Difference(newSet).List(), interactiveWarehouseTable)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(toAdd) > 0 {
			if err := client.Warehouses.AlterInteractive(ctx, sdk.NewAlterInteractiveWarehouseRequest(id).WithAddTables(toAdd)); err != nil {
				return diag.FromErr(fmt.Errorf("error adding tables to interactive warehouse %s: %w", id.FullyQualifiedName(), err))
			}
		}
		if len(toDrop) > 0 {
			if err := client.Warehouses.AlterInteractive(ctx, sdk.NewAlterInteractiveWarehouseRequest(id).WithDropTables(toDrop)); err != nil {
				return diag.FromErr(fmt.Errorf("error dropping tables from interactive warehouse %s: %w", id.FullyQualifiedName(), err))
			}
		}
	}
	return nil
}

func warehouseInteractiveApplySetUnsetExt(d *schema.ResourceData) (*warehouseInteractiveAlterRequestsExt, error) {
	set := sdk.NewWarehouseInteractiveSetRequest()
	unset := sdk.NewWarehouseInteractiveUnsetRequest()

	if err := errors.Join(
		attributeMappedValueUpdateSetOnlyAlways(d, "warehouse_size", &set.WarehouseSize, sdk.ToWarehouseSize),
		intAttributeUpdate(d, "max_cluster_count", &set.MaxClusterCount, &unset.MaxClusterCount),
		intAttributeUpdate(d, "min_cluster_count", &set.MinClusterCount, &unset.MinClusterCount),
		intAttributeWithSpecialDefaultUnsetFallbackUpdate(d, "auto_suspend", &set.AutoSuspend, interactiveWarehouseDefaultAutoSuspend),
		booleanStringAttributeUpdate(d, "auto_resume", &set.AutoResume, &unset.AutoResume),
		accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &set.ResourceMonitor, &unset.ResourceMonitor),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); err != nil {
		return nil, err
	}
	return &warehouseInteractiveAlterRequestsExt{set: set, unset: unset}, nil
}

func warehouseInteractiveSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, w *sdk.WarehouseInteractive, _ bool) error {
	if w.Type != sdk.WarehouseTypeInteractive {
		return fmt.Errorf("warehouse %s is not an interactive warehouse; use snowflake_warehouse or snowflake_warehouse_adaptive instead", id.FullyQualifiedName())
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
	return errs
}

// =============================================================================
// SDK incision points
// =============================================================================

func warehouseInteractiveCreateInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, req *sdk.CreateInteractiveWarehouseRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.Warehouses.CreateInteractivePreservingSession(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating interactive warehouse %s: %w", id.FullyQualifiedName(), err))
	}
	return nil
}

func warehouseInteractiveShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.WarehouseInteractive, diag.Diagnostics) {
	w, err := showWarehouseInteractiveByIDSafely(ctx, meta.(*provider.Context), id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "Failed to query interactive warehouse. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Interactive warehouse id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return w, nil
}

func warehouseInteractiveShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.WarehouseInteractive, error) {
	return showWarehouseInteractiveByID(ctx, meta.(*provider.Context), id)
}

func warehouseInteractiveAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *warehouseInteractiveAlterRequestsExt) error {
	if *reqs.set != *sdk.NewWarehouseInteractiveSetRequest() {
		// AlterInteractiveWithSuspend suspends and resumes the warehouse around the alter when it changes the
		// size, because an interactive warehouse cannot be resized while running. Other SET changes
		// are applied without a suspend.
		client := meta.(*provider.Context).Client
		if err := client.Warehouses.AlterInteractiveWithSuspend(ctx, sdk.NewAlterInteractiveWarehouseRequest(id).WithSet(*reqs.set)); err != nil {
			return fmt.Errorf("error setting interactive warehouse properties: %w", err)
		}
	}
	return nil
}

func warehouseInteractiveAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, reqs *warehouseInteractiveAlterRequestsExt) error {
	if *reqs.unset != *sdk.NewWarehouseInteractiveUnsetRequest() {
		client := meta.(*provider.Context).Client
		if err := client.Warehouses.AlterInteractive(ctx, sdk.NewAlterInteractiveWarehouseRequest(id).WithUnset(*reqs.unset)); err != nil {
			return fmt.Errorf("error unsetting interactive warehouse properties: %w", err)
		}
	}
	return nil
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func interactiveWarehouseTable(v any) (sdk.SchemaObjectIdentifier, error) {
	return sdk.ParseSchemaObjectIdentifier(v.(string))
}
