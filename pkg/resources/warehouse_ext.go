package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
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
	warehouseCreateExt = TrackingCreateWrapper(resources.Warehouse, CreateWarehouse)
	warehouseReadExt   = TrackingReadWrapper(resources.Warehouse, ReadWarehouseFunc(true))
	warehouseUpdateExt = TrackingUpdateWrapper(resources.Warehouse, UpdateWarehouse)
	warehouseDeleteExt = TrackingDeleteWrapper(resources.Warehouse, ResourceDeleteContextFunc(
		warehouseParseIdExt,
		warehouseDropSafelyInSdkExt,
	))
	warehouseDescriptionExt = "Resource used to manage warehouse objects. For more information, check [warehouse documentation](https://docs.snowflake.com/en/sql-reference/commands-warehouse)."
	// Listed on the warehouse SSOT object and appended to the generated customdiff.All.
	warehouseForceNewIfSizeRemoved = customdiff.ForceNewIfChange("warehouse_size", func(_ context.Context, old, new, _ any) bool {
		return old.(string) != "" && new.(string) == ""
	})
	warehouseForceNewIfWasInteractive = customdiff.ForceNewIfChange("warehouse_type", func(_ context.Context, oldValue, _, _ any) bool {
		oldRaw, ok := oldValue.(string)
		if !ok || oldRaw == "" {
			return false
		}
		oldType, err := sdk.ToWarehouseType(oldRaw)
		return err == nil && oldType == sdk.WarehouseTypeInteractive
	})
	warehouseImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.Warehouse, ImportWarehouse),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

var warehouseSetStateToValueFromConfigFieldsExt = []string{
	"warehouse_type",
	"warehouse_size",
	"max_cluster_count",
	"min_cluster_count",
	"scaling_policy",
	"auto_suspend",
	"auto_resume",
	"resource_monitor",
	"enable_query_acceleration",
	"query_acceleration_max_scale_factor",
	"resource_constraint",
	"generation",
}

type warehouseAlterRequestsExt struct {
	set   *sdk.WarehouseSetRequest
	unset *sdk.WarehouseUnsetRequest
}

func warehouseParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(id)
}

func warehouseNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateWarehouseRequest, error) {
	return sdk.NewCreateWarehouseRequest(id), nil
}

func warehouseApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateWarehouseRequest) error {
	return errors.Join(
		attributeMappedValueCreateBuilder(d, "warehouse_type", request.WithWarehouseType, sdk.ToWarehouseTypeUserSettable),
		attributeMappedValueCreateBuilder(d, "warehouse_size", request.WithWarehouseSize, sdk.ToWarehouseSize),
		intAttributeCreateBuilder(d, "max_cluster_count", request.WithMaxClusterCount),
		intAttributeCreateBuilder(d, "min_cluster_count", request.WithMinClusterCount),
		attributeMappedValueCreateBuilder(d, "scaling_policy", request.WithScalingPolicy, sdk.ToScalingPolicy),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend", request.WithAutoSuspend),
		booleanStringAttributeCreateBuilder(d, "auto_resume", request.WithAutoResume),
		boolAttributeCreateBuilder(d, "initially_suspended", request.WithInitiallySuspended),
		accountObjectIdentifierAttributeCreateBuilder(d, "resource_monitor", request.WithResourceMonitor),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		booleanStringAttributeCreateBuilder(d, "enable_query_acceleration", request.WithEnableQueryAcceleration),
		intAttributeWithSpecialDefaultCreateBuilder(d, "query_acceleration_max_scale_factor", request.WithQueryAccelerationMaxScaleFactor),
		attributeMappedValueCreateBuilder(d, "resource_constraint", request.WithResourceConstraint, sdk.ToWarehouseResourceConstraint),
		attributeMappedValueCreateBuilder(d, "generation", request.WithGeneration, sdk.ToWarehouseGeneration),
	)
}

func warehouseHandleExternalChangesExt(d *schema.ResourceData, warehouse *sdk.WarehouseRegular, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	var generation, resourceConstraint string
	if warehouse.Generation != nil {
		generation = string(*warehouse.Generation)
	}
	if warehouse.ResourceConstraint != nil {
		resourceConstraint = string(*warehouse.ResourceConstraint)
	}
	sizeVal, sizeStr := optionalStringOutputMapping(warehouse.Size)
	scalingPolicyVal, scalingPolicyStr := optionalStringOutputMapping(warehouse.ScalingPolicy)
	autoResumeVal, autoResumeStr := warehouse.AutoResume, booleanStringFromBool(warehouse.AutoResume)
	enableQAVal, enableQAStr := optionalBooleanStringOutputMapping(warehouse.EnableQueryAcceleration)
	maxClusterCount := optionalIntOutputMapping(warehouse.MaxClusterCount)
	minClusterCount := optionalIntOutputMapping(warehouse.MinClusterCount)
	autoSuspend := optionalIntOutputMappingIntDefault(warehouse.AutoSuspend)
	queryAccelerationMaxScaleFactor := optionalIntOutputMappingIntDefault(warehouse.QueryAccelerationMaxScaleFactor)
	return handleExternalChangesToObjectInShow(d,
		outputMapping{"type", "warehouse_type", string(warehouse.Type), warehouse.Type, nil},
		outputMapping{"size", "warehouse_size", sizeStr, sizeVal, nil},
		outputMapping{"max_cluster_count", "max_cluster_count", maxClusterCount, maxClusterCount, nil},
		outputMapping{"min_cluster_count", "min_cluster_count", minClusterCount, minClusterCount, nil},
		outputMapping{"scaling_policy", "scaling_policy", scalingPolicyStr, scalingPolicyVal, nil},
		outputMapping{"auto_suspend", "auto_suspend", autoSuspend, autoSuspend, nil},
		outputMapping{"auto_resume", "auto_resume", autoResumeVal, autoResumeStr, nil},
		outputMapping{"resource_monitor", "resource_monitor", warehouse.ResourceMonitor.Name(), warehouse.ResourceMonitor.Name(), nil},
		outputMapping{"enable_query_acceleration", "enable_query_acceleration", enableQAVal, enableQAStr, nil},
		outputMapping{"query_acceleration_max_scale_factor", "query_acceleration_max_scale_factor", queryAccelerationMaxScaleFactor, queryAccelerationMaxScaleFactor, nil},
		outputMapping{"generation", "generation", generation, generation, nil},
		outputMapping{"resource_constraint", "resource_constraint", resourceConstraint, resourceConstraint, nil},
	)
}

func warehouseSetConfigFieldsExt(d *schema.ResourceData, warehouse *sdk.WarehouseRegular) error {
	return errors.Join(
		d.Set("comment", warehouse.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseRegularToSchema(warehouse)}),
	)
}

func warehouseHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.AccountObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	newID, err := renameAccountLevelResource(ctx, d, *id, func(ctx context.Context, id, newID sdk.AccountObjectIdentifier) error {
		return client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithRenameTo(newID))
	})
	if err != nil {
		return diag.FromErr(err)
	}
	*id = newID
	return nil
}

func warehouseApplySetUnsetExt(d *schema.ResourceData) (*warehouseAlterRequestsExt, error) {
	set, unset := sdk.NewWarehouseSetRequest(), sdk.NewWarehouseUnsetRequest()
	if err := errors.Join(
		attributeMappedValueUpdate(d, "warehouse_type", &set.WarehouseType, &unset.WarehouseType, sdk.ToWarehouseTypeUserSettable),
		attributeMappedValueUpdateSetOnlyAlways(d, "warehouse_size", &set.WarehouseSize, sdk.ToWarehouseSize),
		intAttributeUpdate(d, "max_cluster_count", &set.MaxClusterCount, &unset.MaxClusterCount),
		intAttributeUpdate(d, "min_cluster_count", &set.MinClusterCount, &unset.MinClusterCount),
		attributeMappedValueUpdate(d, "scaling_policy", &set.ScalingPolicy, &unset.ScalingPolicy, sdk.ToScalingPolicy),
		// TODO [SNOW-1473453]: UNSET of auto suspend works incorrectly, so removal falls back to 600.
		intAttributeWithSpecialDefaultUnsetFallbackUpdate(d, "auto_suspend", &set.AutoSuspend, 600),
		booleanStringAttributeUpdate(d, "auto_resume", &set.AutoResume, &unset.AutoResume),
		accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &set.ResourceMonitor, &unset.ResourceMonitor),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		booleanStringAttributeUpdate(d, "enable_query_acceleration", &set.EnableQueryAcceleration, &unset.EnableQueryAcceleration),
		intAttributeWithSpecialDefaultUpdate(d, "query_acceleration_max_scale_factor", &set.QueryAccelerationMaxScaleFactor, &unset.QueryAccelerationMaxScaleFactor),
		warehouseResourceConstraintUpdate(d, set, unset),
		warehouseGenerationUpdate(d, set, unset),
	); err != nil {
		return nil, err
	}
	if set.WarehouseSize != nil {
		// For now, we always want to wait for the resize completion. In the future, we may parametrize it.
		set.WithWaitForCompletion(true)
	}
	return &warehouseAlterRequestsExt{set: set, unset: unset}, nil
}

func warehouseSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, warehouse *sdk.WarehouseRegular, _ bool) error {
	if err := d.Set("name", id.Name()); err != nil {
		return err
	}
	return errors.Join(
		d.Set("warehouse_type", warehouse.Type),
		setOptionalValueWithMapping(d, "warehouse_size", warehouse.Size, func(v *sdk.WarehouseSize) string { return string(*v) }),
		d.Set("max_cluster_count", warehouse.MaxClusterCount),
		d.Set("min_cluster_count", warehouse.MinClusterCount),
		setOptionalValueWithMapping(d, "scaling_policy", warehouse.ScalingPolicy, func(v *sdk.ScalingPolicy) string { return string(*v) }),
		d.Set("auto_suspend", warehouse.AutoSuspend),
		d.Set("auto_resume", booleanStringFromBool(warehouse.AutoResume)),
		d.Set("resource_monitor", warehouse.ResourceMonitor.Name()),
		d.Set("comment", warehouse.Comment),
		setOptionalValueWithMapping(d, "enable_query_acceleration", warehouse.EnableQueryAcceleration, func(v *bool) string { return booleanStringFromBool(*v) }),
		d.Set("query_acceleration_max_scale_factor", warehouse.QueryAccelerationMaxScaleFactor),
		setOptionalFromPtr(d, "resource_constraint", warehouse.ResourceConstraint),
		setOptionalFromPtr(d, "generation", warehouse.Generation),
	)
}

func warehouseResourceConstraintUpdate(d *schema.ResourceData, set *sdk.WarehouseSetRequest, unset *sdk.WarehouseUnsetRequest) error {
	if !d.HasChange("resource_constraint") {
		return nil
	}
	warehouseTypeRaw := d.Get("warehouse_type").(string)
	// Resource constraint is only supported for SNOWPARK-OPTIMIZED warehouses.
	// Ignore the resource constraint if the warehouse type is standard or is not set.
	if warehouseTypeRaw == "" {
		return nil
	}
	warehouseType, err := sdk.ToWarehouseType(warehouseTypeRaw)
	if err != nil {
		return err
	}
	if warehouseType != sdk.WarehouseTypeSnowparkOptimized {
		log.Printf("[DEBUG] resource constraint is not supported for %s warehouses, ignoring", warehouseType)
		return nil
	}
	return attributeMappedValueUpdate(d, "resource_constraint", &set.ResourceConstraint, &unset.ResourceConstraint, sdk.ToWarehouseResourceConstraint)
}

func warehouseGenerationUpdate(d *schema.ResourceData, set *sdk.WarehouseSetRequest, unset *sdk.WarehouseUnsetRequest) error {
	if !d.HasChange("generation") {
		return nil
	}
	warehouseTypeRaw := d.Get("warehouse_type").(string)
	setGeneration := warehouseTypeRaw == ""
	if !setGeneration {
		warehouseType, err := sdk.ToWarehouseType(warehouseTypeRaw)
		if err != nil {
			return err
		}
		setGeneration = warehouseType == sdk.WarehouseTypeStandard
	}
	if !setGeneration {
		log.Printf("[DEBUG] generation is not supported for %s warehouses, ignoring", warehouseTypeRaw)
		return nil
	}
	return attributeMappedValueUpdate(d, "generation", &set.Generation, &unset.Generation, sdk.ToWarehouseGeneration)
}

// =============================================================================
// SDK incision points
// =============================================================================

func warehouseCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateWarehouseRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.Warehouses.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func warehouseShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.WarehouseRegular, diag.Diagnostics) {
	warehouse, err := showWarehouseRegularByIDSafely(ctx, meta.(*provider.Context), id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query warehouse. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Warehouse id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return warehouse, nil
}

func warehouseShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.WarehouseRegular, error) {
	return showWarehouseRegularByID(ctx, meta.(*provider.Context), id)
}

func warehouseAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *warehouseAlterRequestsExt) error {
	if *reqs.set != *sdk.NewWarehouseSetRequest() {
		client := meta.(*provider.Context).Client
		return client.Warehouses.AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(*reqs.set))
	}
	return nil
}

func warehouseAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, reqs *warehouseAlterRequestsExt) error {
	if *reqs.unset != *sdk.NewWarehouseUnsetRequest() {
		client := meta.(*provider.Context).Client
		return client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithUnset(*reqs.unset))
	}
	return nil
}

func warehouseDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
	return client.Warehouses.DropSafely
}

func showWarehouseRegularByID(ctx context.Context, providerCtx *provider.Context, id sdk.AccountObjectIdentifier) (*sdk.WarehouseRegular, error) {
	if providerCtx.Experiments.IsEnabled(experimentalfeatures.WarehouseShowImprovedPerformance) {
		warehouse, err := providerCtx.Client.Warehouses.ShowByIDExperimental(ctx, id)
		if err != nil {
			return nil, err
		}
		return warehouse.AsRegular(), nil
	}
	return providerCtx.Client.Warehouses.ShowByIDRegular(ctx, id)
}

func showWarehouseRegularByIDSafely(ctx context.Context, providerCtx *provider.Context, id sdk.AccountObjectIdentifier) (*sdk.WarehouseRegular, error) {
	if providerCtx.Experiments.IsEnabled(experimentalfeatures.WarehouseShowImprovedPerformance) {
		warehouse, err := providerCtx.Client.Warehouses.ShowByIDExperimentalSafely(ctx, id)
		if err != nil {
			return nil, err
		}
		return warehouse.AsRegular(), nil
	}
	return providerCtx.Client.Warehouses.ShowByIDSafelyRegular(ctx, id)
}

func showWarehouseAdaptiveByID(ctx context.Context, providerCtx *provider.Context, id sdk.AccountObjectIdentifier) (*sdk.WarehouseAdaptive, error) {
	if providerCtx.Experiments.IsEnabled(experimentalfeatures.WarehouseShowImprovedPerformance) {
		warehouse, err := providerCtx.Client.Warehouses.ShowByIDExperimental(ctx, id)
		if err != nil {
			return nil, err
		}
		return warehouse.AsAdaptive(), nil
	}
	return providerCtx.Client.Warehouses.ShowByIDAdaptive(ctx, id)
}

func showWarehouseAdaptiveByIDSafely(ctx context.Context, providerCtx *provider.Context, id sdk.AccountObjectIdentifier) (*sdk.WarehouseAdaptive, error) {
	if providerCtx.Experiments.IsEnabled(experimentalfeatures.WarehouseShowImprovedPerformance) {
		warehouse, err := providerCtx.Client.Warehouses.ShowByIDExperimentalSafely(ctx, id)
		if err != nil {
			return nil, err
		}
		return warehouse.AsAdaptive(), nil
	}
	return providerCtx.Client.Warehouses.ShowByIDSafelyAdaptive(ctx, id)
}

func showWarehouseInteractiveByID(ctx context.Context, providerCtx *provider.Context, id sdk.AccountObjectIdentifier) (*sdk.WarehouseInteractive, error) {
	if providerCtx.Experiments.IsEnabled(experimentalfeatures.WarehouseShowImprovedPerformance) {
		warehouse, err := providerCtx.Client.Warehouses.ShowByIDExperimental(ctx, id)
		if err != nil {
			return nil, err
		}
		return warehouse.AsInteractive(), nil
	}
	return providerCtx.Client.Warehouses.ShowByIDInteractive(ctx, id)
}

func showWarehouseInteractiveByIDSafely(ctx context.Context, providerCtx *provider.Context, id sdk.AccountObjectIdentifier) (*sdk.WarehouseInteractive, error) {
	if providerCtx.Experiments.IsEnabled(experimentalfeatures.WarehouseShowImprovedPerformance) {
		warehouse, err := providerCtx.Client.Warehouses.ShowByIDExperimentalSafely(ctx, id)
		if err != nil {
			return nil, err
		}
		return warehouse.AsInteractive(), nil
	}
	return providerCtx.Client.Warehouses.ShowByIDSafelyInteractive(ctx, id)
}
