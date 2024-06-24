package resources

import (
	"context"
	"fmt"
	"strings"

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
func v092WarehouseSizeStateUpgrader(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
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

	return rawState, nil
}
