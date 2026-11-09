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
	warehouseCreateExt      = TrackingCreateWrapper(resources.Warehouse, CreateWarehouse)
	warehouseReadExt        = TrackingReadWrapper(resources.Warehouse, ReadWarehouseFunc(true))
	warehouseUpdateExt      = TrackingUpdateWrapper(resources.Warehouse, UpdateWarehouse)
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

// warehouseTypeCreateBuilder sends STANDARD when warehouse_type is not set,
// so the type does not depend on the Snowflake default (which is changing to ADAPTIVE).
func warehouseTypeCreateBuilder(d *schema.ResourceData, request *sdk.CreateWarehouseRequest) error {
	if _, ok := d.GetOk("warehouse_type"); !ok {
		request.WithWarehouseType(sdk.WarehouseTypeStandard)
		return nil
	}
	return attributeMappedValueCreateBuilder(d, "warehouse_type", request.WithWarehouseType, sdk.ToWarehouseTypeUserSettable)
}

// warehouseTypeUpdate sets STANDARD when warehouse_type is removed from the configuration instead of calling UNSET,
// because UNSET resets the type to the Snowflake default (which is changing to ADAPTIVE). The generated update calls it
// through CustomUpdate, which also passes the UNSET request; it is never used.
func warehouseTypeUpdate(d *schema.ResourceData, set *sdk.WarehouseSetRequest, _ *sdk.WarehouseUnsetRequest) error {
	if !d.HasChange("warehouse_type") {
		return nil
	}
	if _, ok := d.GetOk("warehouse_type"); !ok {
		set.WithWarehouseType(sdk.WarehouseTypeStandard)
		return nil
	}
	var unused *bool
	return attributeMappedValueUpdate(d, "warehouse_type", &set.WarehouseType, &unused, sdk.ToWarehouseTypeUserSettable)
}

func warehouseNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateWarehouseRequest, error) {
	return sdk.NewCreateWarehouseRequest(id), nil
}

func warehouseApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateWarehouseRequest) error {
	return errors.Join(
		warehouseTypeCreateBuilder(d, request),
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
		outputMapping{"enable_query_acceleration", "enable_query_acceleration", enableQAVal, enableQAStr, nil},
		outputMapping{"query_acceleration_max_scale_factor", "query_acceleration_max_scale_factor", queryAccelerationMaxScaleFactor, queryAccelerationMaxScaleFactor, nil},
		outputMapping{"generation", "generation", generation, generation, nil},
		outputMapping{"resource_constraint", "resource_constraint", resourceConstraint, resourceConstraint, nil},
	)
}

func warehouseSetConfigFieldsExt(_ context.Context, d *schema.ResourceData, _ any, warehouse *sdk.WarehouseRegular) error {
	return errors.Join(
		d.Set("comment", warehouse.Comment),
		d.Set("resource_monitor", warehouse.ResourceMonitor.Name()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseRegularToSchema(warehouse)}),
	)
}

func warehouseSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, warehouse *sdk.WarehouseRegular, importPermissiveMode bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		importStringAttribute(d, "warehouse_type", importPermissiveMode, warehouse.Type),
		importOptionalStringAttribute(d, "warehouse_size", importPermissiveMode, warehouse.Size),
		importOptionalIntAttribute(d, "max_cluster_count", importPermissiveMode, warehouse.MaxClusterCount),
		importOptionalIntAttribute(d, "min_cluster_count", importPermissiveMode, warehouse.MinClusterCount),
		importOptionalStringAttribute(d, "scaling_policy", importPermissiveMode, warehouse.ScalingPolicy),
		importOptionalIntWithIntDefaultAttribute(d, "auto_suspend", importPermissiveMode, warehouse.AutoSuspend),
		importBooleanStringAttribute(d, "auto_resume", importPermissiveMode, warehouse.AutoResume),
		importOptionalBooleanStringAttribute(d, "enable_query_acceleration", importPermissiveMode, warehouse.EnableQueryAcceleration),
		importOptionalIntWithIntDefaultAttribute(d, "query_acceleration_max_scale_factor", importPermissiveMode, warehouse.QueryAccelerationMaxScaleFactor),
		importOptionalStringAttribute(d, "resource_constraint", importPermissiveMode, warehouse.ResourceConstraint),
		importOptionalStringAttribute(d, "generation", importPermissiveMode, warehouse.Generation),
	)
}

// warehouseSizeUpdate sets the size, which is never unset, and waits for the resize to complete, for now always.
// In the future, we may parametrize it. The size cannot be unset, so the UNSET request is never used.
func warehouseSizeUpdate(d *schema.ResourceData, set *sdk.WarehouseSetRequest, _ *sdk.WarehouseUnsetRequest) error {
	if err := attributeMappedValueUpdateSetOnlyAlways(d, "warehouse_size", &set.WarehouseSize, sdk.ToWarehouseSize); err != nil {
		return err
	}
	if set.WarehouseSize != nil {
		set.WithWaitForCompletion(true)
	}
	return nil
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

func warehouseAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *warehouseAlterRequests) error {
	if *reqs.set != *sdk.NewWarehouseSetRequest() {
		client := meta.(*provider.Context).Client
		return client.Warehouses.AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(*reqs.set))
	}
	return nil
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
