package resources

import (
	"context"
	"errors"
	"fmt"

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
	warehouseAdaptiveCreateExt        = TrackingCreateWrapper(resources.WarehouseAdaptive, CreateWarehouseAdaptive)
	warehouseAdaptiveReadExt          = TrackingReadWrapper(resources.WarehouseAdaptive, ReadWarehouseAdaptiveFunc(true))
	warehouseAdaptiveUpdateExt        = TrackingUpdateWrapper(resources.WarehouseAdaptive, UpdateWarehouseAdaptive)
	warehouseAdaptiveDescriptionExt   = "Resource used to manage adaptive warehouse objects. Adaptive Compute is a compute service focused on delivering strong performance with effortless operations. It replaces the fixed compute of the Standard Warehouse with a workload-aware one that adapts to your queries automatically. The system decides how to allocate resources for the best performance, eliminating the need for infrastructure tuning."
	warehouseAdaptiveCustomizeDiffExt = TrackingCustomDiffWrapper(resources.WarehouseAdaptive, customdiff.All(
		ComputedIfAnyAttributeChanged(warehouseAdaptiveSchema, ShowOutputAttributeName, "name", "comment", "max_query_performance_level", "query_throughput_multiplier", "resource_monitor"),
		ComputedIfAnyAttributeChanged(
			warehouseAdaptiveParametersAttributesSchema, ParametersAttributeName,
			warehouseAdaptiveParametersFieldNames...,
		),
		ComputedIfAnyAttributeChanged(warehouseAdaptiveSchema, FullyQualifiedNameAttributeName, "name"),
		warehouseAdaptiveParametersCustomDiff,
		HandleWarehouseExternalTypeChange(sdk.WarehouseTypeAdaptive),
	))
	warehouseAdaptiveImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.WarehouseAdaptive, ImportWarehouseAdaptive),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

type warehouseAdaptiveAlterRequestsExt struct {
	set   *sdk.WarehouseAdaptiveSetRequest
	unset *sdk.WarehouseAdaptiveUnsetRequest
}

func warehouseAdaptiveParseIdFromConfigExt(d *schema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
}

func warehouseAdaptiveNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateAdaptiveWarehouseRequest, error) {
	return sdk.NewCreateAdaptiveWarehouseRequest(id), nil
}

func warehouseAdaptiveApplyCreateOptionalsExt(d *schema.ResourceData, opts *sdk.CreateAdaptiveWarehouseRequest) error {
	return errors.Join(
		stringAttributeCreate(d, "comment", &opts.Comment),
		attributeMappedValueCreate(d, "max_query_performance_level", &opts.MaxQueryPerformanceLevel, func(v any) (*sdk.MaxQueryPerformanceLevel, error) {
			s, err := sdk.ToMaxQueryPerformanceLevel(v.(string))
			if err != nil {
				return nil, err
			}
			return &s, nil
		}),
		intAttributeWithSpecialDefaultCreate(d, "query_throughput_multiplier", &opts.QueryThroughputMultiplier),
		accountObjectIdentifierAttributeCreate(d, "resource_monitor", &opts.ResourceMonitor),
	)
}

func warehouseAdaptiveHandleExternalChangesExt(d *schema.ResourceData, w *sdk.WarehouseAdaptive, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	return handleExternalChangesToObjectInShow(d, warehouseAdaptiveOutputMappingsExt(w)...)
}

func warehouseAdaptiveOutputMappingsExt(w *sdk.WarehouseAdaptive) []outputMapping {
	var maxQueryPerformanceLevel string
	if w.MaxQueryPerformanceLevel != nil {
		maxQueryPerformanceLevel = string(*w.MaxQueryPerformanceLevel)
	}
	queryThroughputMultiplier := optionalIntOutputMappingIntDefault(w.QueryThroughputMultiplier)
	return []outputMapping{
		{"max_query_performance_level", "max_query_performance_level", maxQueryPerformanceLevel, maxQueryPerformanceLevel, nil},
		{"query_throughput_multiplier", "query_throughput_multiplier", queryThroughputMultiplier, queryThroughputMultiplier, nil},
		{"resource_monitor", "resource_monitor", w.ResourceMonitor.Name(), w.ResourceMonitor.Name(), nil},
	}
}

func warehouseAdaptiveSetConfigFieldsExt(d *schema.ResourceData, w *sdk.WarehouseAdaptive) error {
	return errors.Join(
		d.Set("name", w.Name),
		d.Set("comment", w.Comment),
		d.Set("warehouse_type", string(w.Type)),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseAdaptiveToSchema(w)}),
	)
}

func warehouseAdaptiveHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.AccountObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	newID, err := renameAccountLevelResource(ctx, d, *id, func(ctx context.Context, id, newId sdk.AccountObjectIdentifier) error {
		return client.Warehouses.AlterAdaptive(ctx, sdk.NewAlterAdaptiveWarehouseRequest(id).WithRenameTo(newId))
	})
	if err != nil {
		return diag.FromErr(err)
	}
	*id = newID
	return nil
}

func warehouseAdaptiveApplySetUnsetExt(d *schema.ResourceData) (*warehouseAdaptiveAlterRequestsExt, error) {
	set := sdk.NewWarehouseAdaptiveSetRequest()
	unset := sdk.NewWarehouseAdaptiveUnsetRequest()

	if d.HasChange("warehouse_type") {
		set.WithWarehouseType(sdk.WarehouseTypeAdaptive)
	}

	if err := errors.Join(
		intAttributeWithSpecialDefaultUpdate(d, "query_throughput_multiplier", &set.QueryThroughputMultiplier, &unset.QueryThroughputMultiplier),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		attributeMappedValueUpdate(d, "max_query_performance_level", &set.MaxQueryPerformanceLevel, &unset.MaxQueryPerformanceLevel, sdk.ToMaxQueryPerformanceLevel),
		accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &set.ResourceMonitor, &unset.ResourceMonitor),
	); err != nil {
		return nil, err
	}
	return &warehouseAdaptiveAlterRequestsExt{set: set, unset: unset}, nil
}

func warehouseAdaptiveSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, w *sdk.WarehouseAdaptive, _ bool) error {
	// Warehouses of other types can be imported and converted to ADAPTIVE with an in-place ALTER during the next update.
	// INTERACTIVE is the exception and Snowflake rejects it.
	if w.Type == sdk.WarehouseTypeInteractive {
		return fmt.Errorf("warehouse %s is an interactive warehouse and cannot be converted to ADAPTIVE; use snowflake_warehouse_interactive instead", id.FullyQualifiedName())
	}

	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("comment", w.Comment),
		setOptionalFromPtr(d, "max_query_performance_level", w.MaxQueryPerformanceLevel),
		setOptionalFromPtr(d, "query_throughput_multiplier", w.QueryThroughputMultiplier),
		d.Set("resource_monitor", w.ResourceMonitor.Name()),
	)
}

// =============================================================================
// SDK incision points
// =============================================================================

func warehouseAdaptiveCreateInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, req *sdk.CreateAdaptiveWarehouseRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.Warehouses.CreateAdaptive(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating adaptive warehouse %s: %w", id.FullyQualifiedName(), err))
	}
	return nil
}

func warehouseAdaptiveShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.WarehouseAdaptive, diag.Diagnostics) {
	w, err := showWarehouseAdaptiveByIDSafely(ctx, meta.(*provider.Context), id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "Failed to query adaptive warehouse. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Adaptive warehouse id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return w, nil
}

func warehouseAdaptiveShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.WarehouseAdaptive, error) {
	return showWarehouseAdaptiveByID(ctx, meta.(*provider.Context), id)
}

func warehouseAdaptiveAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *warehouseAdaptiveAlterRequestsExt) error {
	if *reqs.set != *sdk.NewWarehouseAdaptiveSetRequest() {
		client := meta.(*provider.Context).Client
		if err := client.Warehouses.AlterAdaptive(ctx, sdk.NewAlterAdaptiveWarehouseRequest(id).WithSet(*reqs.set)); err != nil {
			return fmt.Errorf("error setting adaptive warehouse properties: %w", err)
		}
	}
	return nil
}

func warehouseAdaptiveAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, reqs *warehouseAdaptiveAlterRequestsExt) error {
	if *reqs.unset != *sdk.NewWarehouseAdaptiveUnsetRequest() {
		client := meta.(*provider.Context).Client
		if err := client.Warehouses.AlterAdaptive(ctx, sdk.NewAlterAdaptiveWarehouseRequest(id).WithUnset(*reqs.unset)); err != nil {
			return fmt.Errorf("error unsetting adaptive warehouse properties: %w", err)
		}
	}
	return nil
}
