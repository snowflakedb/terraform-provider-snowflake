package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v092ToWarehouseSize(s string) (sdk.WarehouseSize, error) {
	s = strings.ToUpper(s)
	switch s {
	case "XSMALL", "X-SMALL":
		return sdk.WarehouseSizeXSmall, nil
	case "SMALL":
		return sdk.WarehouseSizeSmall, nil
	case "MEDIUM":
		return sdk.WarehouseSizeMedium, nil
	case "LARGE":
		return sdk.WarehouseSizeLarge, nil
	case "XLARGE", "X-LARGE":
		return sdk.WarehouseSizeXLarge, nil
	case "XXLARGE", "X2LARGE", "2X-LARGE", "2XLARGE":
		return sdk.WarehouseSizeXXLarge, nil
	case "XXXLARGE", "X3LARGE", "3X-LARGE", "3XLARGE":
		return sdk.WarehouseSizeXXXLarge, nil
	case "X4LARGE", "4X-LARGE", "4XLARGE":
		return sdk.WarehouseSizeX4Large, nil
	case "X5LARGE", "5X-LARGE", "5XLARGE":
		return sdk.WarehouseSizeX5Large, nil
	case "X6LARGE", "6X-LARGE", "6XLARGE":
		return sdk.WarehouseSizeX6Large, nil
	default:
		return "", fmt.Errorf("invalid warehouse size: %s", s)
	}
}

// v092WarehouseSizeStateUpgrader is needed because:
// - we are removing incorrect mapped values from sdk.ToWarehouseSize (like 2XLARGE, 3XLARGE, ...); result of:
//   - https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/1873
//   - https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/1946
//   - https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1889#issuecomment-1631149585
//
// - deprecated wait_for_provisioning attribute was removed
// - clear the old resource monitor representation
// - set query_acceleration_max_scale_factor
func v092WarehouseSizeStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldWarehouseSize := rawState["warehouse_size"].(string)
	if oldWarehouseSize != "" {
		warehouseSize, err := v092ToWarehouseSize(oldWarehouseSize)
		if err != nil {
			return nil, err
		}
		rawState["warehouse_size"] = string(warehouseSize)
	}

	// remove deprecated attribute
	delete(rawState, "wait_for_provisioning")

	// clear the old resource monitor representation
	oldResourceMonitor := rawState["resource_monitor"].(string)
	if oldResourceMonitor == "null" {
		delete(rawState, "resource_monitor")
	}

	// get the warehouse from Snowflake
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(rawState["id"].(string)).(sdk.AccountObjectIdentifier)

	w, err := client.Warehouses.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// fill out query_acceleration_max_scale_factor in state if it was disabled before (old coupling logic that was removed)
	// - if config have no value, then we should have a UNSET after migration
	// - if config have the same value, then we should have a no-op after migration
	// - if config have different value, then we will have SET after migration
	previousEnableQueryAcceleration := rawState["enable_query_acceleration"].(bool)
	if !previousEnableQueryAcceleration {
		rawState["query_acceleration_max_scale_factor"] = w.QueryAccelerationMaxScaleFactor
	}

	return rawState, nil
}

// v2_6_0_WarehouseResourceConstraintUpgrader is needed because:
// - we are adding resource_constraint to the show output and to the config
// - this means that the old output in handleExternalChangesToObject would be empty after migration
// - that caused the plan to show a diff for resource_constraint
// - this upgrade fixes that by setting the default value for resource_constraint
func v2_6_0_WarehouseResourceConstraintUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldShowOutputRaw, ok := rawState[ShowOutputAttributeName].([]any)
	if !ok || len(oldShowOutputRaw) != 1 {
		return nil, fmt.Errorf("expected exactly one warehouse show output; got %d", len(oldShowOutputRaw))
	}
	oldShowOutput, ok := oldShowOutputRaw[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cannot read warehouse show output from the state")
	}

	warehouseName, ok := rawState["name"].(string)
	if !ok {
		return nil, fmt.Errorf("cannot read warehouse name from the state")
	}

	client := meta.(*provider.Context).Client
	warehouseInSnowflake, err := client.Warehouses.ShowByIDSafely(ctx, sdk.NewAccountObjectIdentifier(warehouseName))
	if err != nil {
		return nil, err
	}

	if warehouseInSnowflake.ResourceConstraint != nil {
		switch warehouseInSnowflake.Type {
		case sdk.WarehouseTypeSnowparkOptimized:
			oldShowOutput["resource_constraint"] = string(*warehouseInSnowflake.ResourceConstraint)
		default:
			log.Printf("[DEBUG] handling resource_constraint for warehouse type %s is not supported, ignoring", warehouseInSnowflake.Type)
		}
	}
	rawState[ShowOutputAttributeName] = []any{oldShowOutput}

	return rawState, nil
}
